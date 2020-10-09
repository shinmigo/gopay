package kernel

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Palmer interface {
	GetAliPayMethod() string
}

type AliPayClient struct {
	gatewayHost         string                    //支付宝网关地址
	appId               string                    //商户支付宝应用APPID
	aliPayPublicKeyList map[string]*rsa.PublicKey //支付宝公钥
	merchantPrivateKey  *rsa.PrivateKey           //商户应用私钥
	merchantCertSN      string                    //商户支付宝应用 公钥证书序列号
	aliPayCertSN        string                    //支付宝公钥证书序列号
	aliPayRootCertSN    string                    //支付宝根证书序列号
	notifyUrl           string                    //异步通知地址
	encryptKey          string                    //
	localTimeZone       string                    //时区
	isProd              bool                      //是否为生产环境
	signType            string                    //签名类型
}

/**
 * 支付宝返回错误时的字段
 */
type ErrorRes struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
}

func (m *ErrorRes) Error() string {
	return fmt.Sprintf("%s - %s", m.Code, m.SubMsg)
}

/**
 * 初始化支付宝客户端
 */
func NewAliPayClient(config *Config) (*AliPayClient, error) {
	privateKey, err := ParsePrivateKey(config.MerchantPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	publicKey, err := ParsePublicKey(config.AliPayPublicKeyPath)
	if err != nil {
		return nil, err
	}

	aliPayPublicKeyList := make(map[string]*rsa.PublicKey, 8)
	aliPayPublicKeyList[AliPayPublicKeySN] = publicKey
	client := AliPayClient{
		gatewayHost:         AliPayBoxURL,
		appId:               config.AppId,
		aliPayPublicKeyList: aliPayPublicKeyList,
		merchantPrivateKey:  privateKey,
		aliPayCertSN:        AliPayPublicKeySN,
		notifyUrl:           config.NotifyUrl,
		encryptKey:          config.EncryptKey,
		localTimeZone:       "Asia/Shanghai",
		signType:            AliPaySignType,
	}
	if len(config.SignType) > 0 {
		client.signType = config.SignType
	}
	if config.IsProd {
		client.gatewayHost = AliPayProdURL
	}
	if len(config.LocalTimeZone) > 0 {
		client.localTimeZone = config.LocalTimeZone
	}
	if len(config.MerchantCertPath) > 0 {
		merchantCertSN, err := GetMerchantCertSN(config.MerchantCertPath)
		if err != nil {
			return nil, err
		}

		aliPayCertSN, aliPayPublicKey, err := GetAliPayCertSN(config.AliPayCertPath)
		if err != nil {
			return nil, err
		}

		aliPayRootCertSN, err := GetAliPayRootCertSN(config.AliPayRootCertPath)
		if err != nil {
			return nil, err
		}

		client.merchantCertSN = merchantCertSN
		client.aliPayCertSN = aliPayCertSN
		client.aliPayPublicKeyList[aliPayCertSN] = aliPayPublicKey
		client.aliPayRootCertSN = aliPayRootCertSN
	}

	return &client, nil
}

/**
 * 获取支付宝网关接口
 */
func (m *AliPayClient) GetGatewayHost() string {
	return m.gatewayHost
}

/**
 * 组装支付宝请求参数
 */
func (m *AliPayClient) UrlParams(param Palmer) (url.Values, error) {
	if param == nil {
		return nil, errors.New(InitializeDataErr)
	}
	aliPayRequestMethod := param.GetAliPayMethod()
	if len(aliPayRequestMethod) == 0 {
		return nil, errors.New(InitializeDataErr)
	}

	bizContentBytes, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	timestamp := time.Now()
	location, err := time.LoadLocation(m.localTimeZone)
	if err != nil {
		return nil, errors.New(TimeZoneErr)
	}
	timestamp = timestamp.In(location)
	timestampStr := timestamp.Format(AliPayTimeFormat)
	notifyUrl := m.notifyUrl

	urlMap := url.Values{}
	urlMap.Add("app_id", m.appId)
	urlMap.Add("method", aliPayRequestMethod)
	urlMap.Add("format", AliPayFormat)
	urlMap.Add("charset", AliPayCharset)
	urlMap.Add("sign_type", m.signType)
	urlMap.Add("timestamp", timestampStr)
	urlMap.Add("version", AliPayVersion)
	urlMap.Add("notify_url", notifyUrl)
	urlMap.Add("biz_content", string(bizContentBytes))
	if len(m.merchantCertSN) > 0 {
		urlMap.Add("app_cert_sn", m.merchantCertSN)
		urlMap.Add("alipay_root_cert_sn", m.aliPayRootCertSN)
	}
	sign, err := sign(urlMap, m.merchantPrivateKey, m.signType)
	if err != nil {
		return nil, err
	}
	urlMap.Add("sign", sign)

	return urlMap, nil
}

/**
 * 发送HTTP请求
 */
func (m *AliPayClient) SendRequest(method string, param Palmer, result interface{}) (err error) {
	if param == nil {
		return errors.New(InitializeDataErr)
	}

	urlMap, err := m.UrlParams(param)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(method, m.gatewayHost, strings.NewReader(urlMap.Encode()))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", ContentType)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	responseByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	responseString := string(responseByte)

	//以下是对支付宝返回json结果字符串验证签名
	//首先判断有没有错误
	var sign string
	var certSN string
	var content string
	methodNodeName := strings.ReplaceAll(param.GetAliPayMethod(), ".", "_") + "_response"
	methodNodeNameIndex := strings.LastIndex(responseString, methodNodeName)
	methodErrorIndex := strings.LastIndex(responseString, AliPayErrorResponse)
	if methodNodeNameIndex > 0 {
		content, certSN, sign = m.parseJSONSource(responseString, methodNodeName, methodNodeNameIndex)
		if sign == "" {
			errorRes := &ErrorRes{}
			if err = json.Unmarshal([]byte(content), errorRes); err != nil {
				return err
			}
			if errorRes.Code != CodeSuccess {
				if errorRes != nil {
					return errorRes
				}
				return SignNotFound
			}
		}
	} else if methodErrorIndex > 0 {
		content, certSN, sign = m.parseJSONSource(responseString, AliPayErrorResponse, methodErrorIndex)
		if sign == "" {
			errorRes := &ErrorRes{}
			if err = json.Unmarshal([]byte(content), errorRes); err != nil {
				return err
			}
			return errorRes
		}
	} else {
		return SignNotFound
	}
	if sign != "" {
		aliPayPublicKey, err := m.getAliPayPublicKey(certSN)
		if err != nil {
			return err
		}
		if ok, err := m.verifyData([]byte(content), sign, aliPayPublicKey); ok == false {
			return err
		}
	}
	err = json.Unmarshal(responseByte, result)
	if err != nil {
		return err
	}

	return err
}

/**
 * 异步通知验证
 */
func (m *AliPayClient) NotifyVerify(notifyData url.Values) (bool, error) {
	paramList := make([]string, 0, 16)
	for notifyKey := range notifyData {
		if notifyKey == AliPaySignNodeName || notifyKey == AliPaySignTypeNodeName || notifyKey == AliPayCertSNNodeName {
			continue
		}
		paramValue := notifyData.Get(notifyKey)
		if len(paramValue) == 0 {
			continue
		}

		paramValue = strings.TrimSpace(paramValue)
		paramList = append(paramList, notifyKey+"="+paramValue)
	}
	sort.Strings(paramList)
	notifyParam := strings.Join(paramList, "&")

	return m.verifyData([]byte(notifyParam), notifyData.Get(AliPaySignNodeName), m.aliPayPublicKeyList[m.aliPayCertSN])
}

/**
 * 验证数据
 */
func (m *AliPayClient) verifyData(data []byte, sign string, key *rsa.PublicKey) (bool, error) {
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}

	err = Verify(data, signBytes, key, m.signType)
	if err != nil {
		return false, err
	}
	return true, nil
}

/**
 * 获取支付宝公钥
 */
func (m *AliPayClient) getAliPayPublicKey(certSN string) (key *rsa.PublicKey, err error) {
	if certSN == "" {
		certSN = m.aliPayCertSN
	}
	key = m.aliPayPublicKeyList[certSN]
	if key == nil {
		return nil, AliPayPublicKeyNotFound
	}

	return key, nil
}

/**
 * 解析支付宝响应JSON字符串
 */
func (m *AliPayClient) parseJSONSource(rawData, nodeName string, nodeIndex int) (content, certSN, sign string) {
	var dataEndIndex int
	dataStartIndex := nodeIndex + len(nodeName) + 2
	signIndex := strings.LastIndex(rawData, "\""+AliPaySignNodeName+"\"")
	certIndex := strings.LastIndex(rawData, "\""+AliPayCertSNNodeName+"\"")

	if signIndex > 0 && certIndex > 0 {
		dataEndIndex = int(math.Min(float64(signIndex), float64(certIndex))) - 1
	} else if certIndex > 0 {
		dataEndIndex = certIndex - 1
	} else if signIndex > 0 {
		dataEndIndex = signIndex - 1
	} else {
		dataEndIndex = len(rawData) - 1
	}
	indexLen := dataEndIndex - dataStartIndex
	if indexLen < 0 {
		return "", "", ""
	}

	if certIndex > 0 {
		certStartIndex := certIndex + len(AliPayCertSNNodeName) + 4
		certSN = rawData[certStartIndex:]
		certEndIndex := strings.Index(certSN, "\"")
		certSN = certSN[:certEndIndex]
	}

	if signIndex > 0 {
		signStartIndex := signIndex + len(AliPaySignNodeName) + 4
		sign = rawData[signStartIndex:]
		signEndIndex := strings.LastIndex(sign, "\"")
		sign = sign[:signEndIndex]
	}

	return rawData[dataStartIndex:dataEndIndex], certSN, sign
}

/**
 * 组装支付宝签名
 */
func sign(params url.Values, merchantPrivateKey *rsa.PrivateKey, signType string) (string, error) {
	paramList := make([]string, 0, 16)
	for paramKey := range params {
		paramValue := params.Get(paramKey)
		if len(paramValue) == 0 {
			continue
		}

		paramValue = strings.TrimSpace(paramValue)
		paramList = append(paramList, paramKey+"="+paramValue)
	}
	sort.Strings(paramList)

	requestParam := strings.Join(paramList, "&")
	signByte, err := Sign([]byte(requestParam), merchantPrivateKey, signType)
	if err != nil {
		return "", err
	}
	sign := base64.StdEncoding.EncodeToString(signByte)

	return sign, nil
}
