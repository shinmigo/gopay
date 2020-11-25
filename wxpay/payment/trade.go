package payment

import (
	"fmt"
	"net/url"
)

const (
	WX_JSAPI  = "JSAPI"
	WX_NATIVE = "NATIVE"
	WX_APP    = "APP"
	WX_MWEB   = "MWEB"
)

/**
 * 微信APP支付、公众号支付、小程序支付
 */
type Trade struct {
	DeviceInfo     string //设备号
	Body           string //商品描述
	Detail         string //商品详情
	Attach         string //附加数据
	OutTradeNo     string //商户订单号
	FeeType        string //货币类型
	TotalFee       uint64 //总金额 分
	SpbillCreateIp string //终端IP 用户的客户端IP
	TimeStart      string //订单生成时间
	TimeExpire     string //订单失效时间
	GoodsTag       string //订单优惠标记
	NotifyUrl      string //异步通知回调地址
	TradeType      string //支付类型
	ProductId      string //商品ID
	LimitPay       string //指定支付方式
	OpenId         string //用户标识
	Receipt        string //开发票入口开放标识
}

func (m *Trade) Params() url.Values {
	paramMap := url.Values{}
	paramMap.Set("device_info", m.DeviceInfo)
	paramMap.Set("body", m.Body)
	paramMap.Set("detail", m.Detail)
	paramMap.Set("attach", m.Attach)
	paramMap.Set("out_trade_no", m.OutTradeNo)
	paramMap.Set("fee_type", m.FeeType)
	if len(m.FeeType) == 0 {
		paramMap.Set("fee_type", "CNY")
	}
	paramMap.Set("total_fee", fmt.Sprintf("%d", m.TotalFee))
	paramMap.Set("spbill_create_ip", m.SpbillCreateIp)
	paramMap.Set("time_start", m.TimeStart)
	paramMap.Set("time_expire", m.TimeExpire)
	paramMap.Set("goods_tag", m.GoodsTag)
	paramMap.Set("notify_url", m.NotifyUrl)
	paramMap.Set("trade_type", m.TradeType)
	paramMap.Set("product_id", m.ProductId)
	paramMap.Set("limit_pay", m.LimitPay)
	paramMap.Set("openid", m.OpenId)
	paramMap.Set("sign_type", "MD5")
	
	return paramMap
}

/**
 * 微信支付的响应结果
 */
type TradeRes struct {
	ReturnCode string `xml:"return_code"`  //返回状态码
	ReturnMsg  string `xml:"return_msg"`   //返回信息
	AppId      string `xml:"app_id"`       //应用APPId
	MchId      string `xml:"mch_id"`       //商户号
	DeviceInfo string `xml:"device_info"`  //设备号
	NonceStr   string `xml:"nonce_str"`    //随机字符串
	Sign       string `xml:"sign"`         //签名
	ResultCode string `xml:"result_code"`  //业务结果
	ResultMsg  string `xml:"result_msg"`   //业务结果描述
	ErrCode    string `xml:"err_code"`     //错误代码
	ErrCodeDes string `xml:"err_code_des"` //错误代码描述
	TradeType  string `xml:"trade_type"`   //交易类型
	PrepayId   string `xml:"prepay_id"`    //预支付交易会话标识
	CodeUrl    string `xml:"code_url"`     //二维码链接
}

/**
 * 微信查询订单
 */
type TradeQuery struct {
	TransactionId string //微信订单号
	OutTradeNo    string //商户订单号
}

func (m *TradeQuery) Params() url.Values {
	paramMap := url.Values{}
	paramMap.Set("transaction_id", m.TransactionId)
	paramMap.Set("out_trade_no", m.OutTradeNo)
	
	return paramMap
}

type TradeQueryRes struct {
	ReturnCode         string `xml:"return_code"`          //返回状态码
	ReturnMsg          string `xml:"return_msg"`           //返回信息
	AppId              string `xml:"app_id"`               //应用APPId
	MchId              string `xml:"mch_id"`               //商户号
	NonceStr           string `xml:"nonce_str"`            //随机字符串
	Sign               string `xml:"sign"`                 //签名
	ResultCode         string `xml:"result_code"`          //业务结果
	ResultMsg          string `xml:"result_msg"`           //业务结果描述
	ErrCode            string `xml:"err_code"`             //错误代码
	ErrCodeDes         string `xml:"err_code_des"`         //错误代码描述
	DeviceInfo         string `xml:"device_info"`          //设备号
	OpenId             string `xml:"openid"`               //用户标识
	IsSubscribe        string `xml:"is_subscribe"`         //是否关注公众账号
	TradeType          string `xml:"trade_type"`           //交易类型
	TradeState         string `xml:"trade_state"`          //交易状态
	BankType           string `xml:"bank_type"`            //付款银行
	TotalFee           int    `xml:"total_fee"`            //订单总金额 单位分
	SettlementTotalFee int    `xml:"settlement_total_fee"` //应结订单金额
	FeeType            string `xml:"fee_type"`             //货币类型
	CashFee            int    `xml:"cash_fee"`             //现金支付金额
	CashFeeType        string `xml:"cash_fee_type"`        //现金支付币种
	CouponFee          int    `xml:"coupon_fee"`           //代金券金额
	CouponCount        int    `xml:"coupon_count"`         //代金券使用数量
	TransactionId      string `xml:"transaction_id"`       //微信支付订单号
	OutTradeNo         string `xml:"out_trade_no"`         //商户订单号
	Attach             string `xml:"attach"`               //附加数据
	TimeEnd            string `xml:"time_end"`             //订单支付时间
	TradeStateDesc     string `xml:"trade_state_desc"`     //交易状态描述
}

/**
 * 微信关闭订单
 */
type TradeClose struct {
	OutTradeNo string //商户订单号
}

func (m *TradeClose) Params() url.Values {
	paramMap := url.Values{}
	paramMap.Set("out_trade_no", m.OutTradeNo)
	
	return paramMap
}

type TradeCloseRes struct {
	ReturnCode string `xml:"return_code"`  //返回状态码
	ReturnMsg  string `xml:"return_msg"`   //返回信息
	AppId      string `xml:"app_id"`       //应用APPId
	MchId      string `xml:"mch_id"`       //商户号
	NonceStr   string `xml:"nonce_str"`    //随机字符串
	Sign       string `xml:"sign"`         //签名
	ResultCode string `xml:"result_code"`  //业务结果
	ResultMsg  string `xml:"result_msg"`   //业务结果描述
	ErrCode    string `xml:"err_code"`     //错误代码
	ErrCodeDes string `xml:"err_code_des"` //错误代码描述
}

/**
 * 异步通知
 */
type NotifyRes struct {
	ReturnCode         string `xml:"return_code"`
	ReturnMsg          string `xml:"return_msg"`
	AppId              string `xml:"appid"`
	MCHId              string `xml:"mch_id"`
	DeviceInfo         string `xml:"device_info"`
	NonceStr           string `xml:"nonce_str"`
	Sign               string `xml:"sign"`
	SignType           string `xml:"sign_type"`
	ResultCode         string `xml:"result_code"`
	ErrCode            string `xml:"err_code"`
	ErrCodeDes         string `xml:"err_code_des"`
	OpenId             string `xml:"openid"`
	IsSubscribe        string `xml:"is_subscribe"`
	TradeType          string `xml:"trade_type"`
	BankType           string `xml:"bank_type"`
	TotalFee           int    `xml:"total_fee"`
	SettlementTotalFee int    `xml:"settlement_total_fee"`
	FeeType            string `xml:"fee_type"`
	CashFee            int    `xml:"cash_fee"`
	CashFeeType        string `xml:"cash_fee_type"`
	CouponFee          int    `xml:"coupon_fee"`
	CouponCount        int    `xml:"coupon_count"`
	TransactionId      string `xml:"transaction_id"`
	OutTradeNo         string `xml:"out_trade_no"`
	Attach             string `xml:"attach"`
	TimeEnd            string `xml:"time_end"`
}
