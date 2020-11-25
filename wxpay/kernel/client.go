package kernel

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type WXPayParam interface {
	// 返回参数列表
	Params() url.Values
}

type WxClient struct {
	appId       string //应用ID
	mchId       string //商户号
	md5Key      string //MD5key
	isProd      bool   //环境
	gatewayHost string //网关地址
}

/**
 * 初始化微信支付参数
 */
func NewWxClient(appId, mchId, md5Key string, isProd bool) *WxClient {
	client := &WxClient{
		appId:       appId,
		mchId:       mchId,
		isProd:      isProd,
		md5Key:      md5Key,
		gatewayHost: "https://api.mch.weixin.qq.com/sandboxnew/",
	}
	if isProd {
		client.gatewayHost = "https://api.mch.weixin.qq.com/"
	}
	
	return client
}

/**
 * 发送微信支付请求
 */
func (m *WxClient) SendRequest(method string, url string, param WXPayParam, result interface{}) (err error) {
	requestParam := m.UrlParams(param)
	requestParamXml := mapToXml(requestParam)
	request, err := http.NewRequest(method, m.gatewayHost+url, strings.NewReader(requestParamXml))
	if err != nil {
		return err
	}
	request.Header.Set("Accept", "application/xml")
	request.Header.Set("Content-Type", "application/xml;charset=utf-8")
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
	err = m.VerifySign(responseByte)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(responseByte, result)
	
	return
}

func (m *WxClient) Jsapi(signType, prepayId, nonceStr string) (param url.Values) {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	param = url.Values{}
	param.Set("appId", m.appId)
	param.Set("timeStamp", timeStamp)
	param.Set("signType", signType)
	param.Set("package", fmt.Sprintf("prepay_id=%s", prepayId))
	param.Set("nonceStr", nonceStr)
	param.Set("paySign", m.sign(param))
	
	return
}

/**
 * 组装微信支付公共参数
 */
func (m *WxClient) UrlParams(param WXPayParam) (requestParam url.Values) {
	requestParam = param.Params()
	requestParam.Set("appid", m.appId)
	requestParam.Set("mch_id", m.mchId)
	requestParam.Set("nonce_str", getNonceStr())
	requestParam.Set("sign", m.sign(requestParam))
	
	return
}

/**
 * 验证微信支付响应结果签名
 */
func (m *WxClient) VerifySign(data []byte) (err error) {
	xmlHandler := make(XmlToMap)
	err = xml.Unmarshal(data, &xmlHandler)
	if err != nil {
		return err
	}
	
	returnCode := xmlHandler.Get("return_code")
	if returnCode == "" {
		return errors.New("解析失败！")
	}
	if returnCode == "FAIL" {
		return errors.New(xmlHandler.Get("return_msg"))
	}
	resultCode := xmlHandler.Get("result_code")
	if resultCode == "" {
		return errors.New("解析失败！")
	}
	if returnCode == "FAIL" {
		return errors.New(xmlHandler.Get("err_code_des"))
	}
	
	srcSign := xmlHandler.Get("sign")
	if srcSign == "" {
		return errors.New("解析失败！")
	}
	delete(xmlHandler, "sign")
	generateSign := m.sign(url.Values(xmlHandler))
	if srcSign == generateSign {
		return nil
	}
	
	return errors.New("签名验证失败")
}

/**
 * 组装微信签名
 */
func (m *WxClient) sign(params url.Values) string {
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
	if len(m.md5Key) > 0 {
		paramList = append(paramList, "key="+m.md5Key)
	}
	requestParam := strings.Join(paramList, "&")
	
	md5Handler := md5.New()
	md5Handler.Write([]byte(requestParam))
	hashByte := md5Handler.Sum(nil)
	
	return strings.ToUpper(hex.EncodeToString(hashByte))
}

/**
 * 微信支付随机字符串，长度要求在32位以内。
 */
func getNonceStr() string {
	srcStr := "abcdefghijklmnopqrstuvwxyz0123456789"
	srcStrLen := len(srcStr) - 1
	
	buffer := bytes.Buffer{}
	randHandler := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 32; i++ {
		index := randHandler.Intn(srcStrLen)
		buffer.WriteString(srcStr[index : index+1])
	}
	
	return buffer.String()
}
