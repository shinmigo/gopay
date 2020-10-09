package kernel

import "errors"

const (
	/**
	 * 请求地址相关UR
	 */
	AliPayBoxURL  = "https://openapi.alipaydev.com/gateway.do"
	AliPayProdURL = "https://openapi.alipay.com/gateway.do"

	/**
	 * 默认的签名算法，统一固定使用RSA2签名算法（即SHA_256_WITH_RSA）
	 */
	AliPaySignType = "RSA"

	AliPaySignType2 = "RSA2"

	/**
	 * 默认字符集编码，统一固定使用UTF-8编码
	 */
	AliPayCharset = "UTF-8"

	/**
	 * 支付宝调用的接口版本
	 */
	AliPayVersion = "1.0"

	/**
	 * 支付宝调用的接口数据类型
	 */
	AliPayFormat = "json"

	/**
	 * 支付宝请求数据类型
	 */
	ContentType = "application/x-www-form-urlencoded;charset=utf-8"

	/*
	 * 支付宝时间戳
	 */
	AliPayTimeFormat = "2006-01-02 15:04:05"

	/*
	 * 支付宝公钥Key SN码
	 */
	AliPayPublicKeySN = "alipay-public-key"

	/*
	 * 支付宝根证书结束符
	 */
	AliPayRootCertEnd = "-----END CERTIFICATE-----"

	/*
	 * 支付宝PKCS1|PKCS8格式符号
	 */
	AliPayPKCS1Prefix     = "-----BEGIN RSA PRIVATE KEY-----"
	AliPayPKCS1Suffix     = "-----END RSA PRIVATE KEY-----"
	AliPayPKCS8Prefix     = "-----BEGIN PRIVATE KEY-----"
	AliPayPKCS8Suffix     = "-----END PRIVATE KEY-----"
	AliPayPublicKeyPrefix = "-----BEGIN PUBLIC KEY-----"
	AliPayPublicKeySuffix = "-----END PUBLIC KEY-----"

	/*
	 * 错误信息
	 */
	MerchantPrivateKeyEmpty = "merchant private key cannot be empty"
	PrivateKeyWrongFormat   = "incorrect private key format"
	PublicKeyWrongFormat    = "incorrect public key format"
	CertEmpty               = "the certificate cannot be empty"
	CertWrongFormat         = "certificate format error"
	InitializeDataErr       = "please initialize the data"
	TimeZoneErr             = "time zone error"

	/*
	 * 响应信息
	 */
	AliPaySignNodeName     = "sign"
	AliPaySignTypeNodeName = "sign_type"
	AliPayCertSNNodeName   = "alipay_cert_sn"
	AliPayErrorResponse    = "error_response"

	CodeSuccess string = "10000" // 接口调用成功
)

var (
	SignNotFound            = errors.New("alipay: sign content not found")
	AliPayPublicKeyNotFound = errors.New("alipay: alipay public key not found")
)
