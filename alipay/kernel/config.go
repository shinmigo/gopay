package kernel

type Config struct {
	AppId                  string //商户支付宝应用APPID
	AliPayPublicKeyPath    string //支付宝公钥路径
	MerchantPrivateKeyPath string //商户应用私钥路径
	AliPayCertPath         string //支付宝公钥证书路径
	AliPayRootCertPath     string //支付宝根证书文件路径
	MerchantCertPath       string //商户支付宝应用 公钥证书路径
	NotifyUrl              string //异步通知地址
	EncryptKey             string //可设置AES密钥，调用AES加解密相关接口时需要（可选）
	IsProd                 bool   //是否为生产环境
	LocalTimeZone          string //时区
	SignType               string //签名类型
}
