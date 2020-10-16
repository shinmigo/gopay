# goPay
微信支付|支付宝支付 for golang

## Installation
```go
go get github.com/shinmigo/gopay
```

## 微信支付

### Usage

初始化客户端

```go
//初始化微信支付客户端 初始化一次就可以
wxClient := kernel.NewWxClient("", "", "", false)
```



创建统一支付订单

```go
wxPayment := payment.Payment{Client: wxClient}
//调用微信支付(统一下单接口) 初始化请求参数
wxPayParam := payment.Trade{
   Body:       "",
   Detail:     "",
   OutTradeNo: "",
   TotalFee:   0,
}
payRes, err := wxPayment.Pay(&wxPayParam)
```



查询订单

```go
wxPayment := payment.Payment{Client: wxClient}

//订单查询
queryOrderParam := payment.TradeQuery{
   TransactionId: "",
   OutTradeNo:    "",
}
tradeRes ,err := wxPayment.Query(&queryOrderParam)
```



关闭订单

```go
wxPayment := payment.Payment{Client: wxClient}

//关闭订单
closeOrderParam := payment.TradeClose{
   OutTradeNo:"",
}
tradeRes, err := wxPayment.Close(&closeOrderParam)
```



异步通知验证签名

```go
wxPayment := payment.Payment{Client: wxClient}
//异步验证签名
tradeRes, err := wxPayment.NotifyVerify([]byte(``))
```

## 支付宝支付

### Usage

初始化客户端

```go
config := &kernel.Config{
   AppId:                  "",
   AliPayPublicKeyPath:    "",
   MerchantPrivateKeyPath: "./cert/rsa_private_key.pem",
   AliPayCertPath:         "",
   AliPayRootCertPath:     "",
   MerchantCertPath:       "",
   NotifyUrl:              "",
   EncryptKey:             "",
   IsProd:                 true,
   LocalTimeZone:          "",
}
aliPayClient, err := kernel.NewAliPayClient(config)
```



APP支付

```go
paymentTrade := payment.Payment{Client: aliPayClient}

appPay := payment.App{
   Trade: payment.Trade{
      Subject:     "测试",
      OutTradeNo:  "2020090723897",
      TotalAmount: "100",
   },
}
res, err := paymentTrade.App(&appPay)
```



手机网站支付

```go
paymentTrade := payment.Payment{Client: aliPayClient}

wapPay := payment.Wap{
   Trade: payment.Trade{
      Subject:     "测试",
      OutTradeNo:  "2020090723897",
      TotalAmount: "100",
   },
}
res, err := paymentTrade.Wap(&wapPay)
```



PC网站支付

```go
paymentTrade := payment.Payment{Client: aliPayClient}

pagePay := &payment.Page{
   Trade: payment.Trade{
      Subject:     "测试",
      OutTradeNo:  "2020090723897",
      TotalAmount: "100",
   },
}
res, err := paymentTrade.Page(&pagePay)
```



订单查询

```go
paymentTrade := payment.Payment{Client: aliPayClient}

queryOrderParam := payment.TradeQuery{
   OutTradeNo:   "",
   TradeNo:      "",
   OrgPid:       "",
   QueryOptions: nil,
}
tradeRes ,err := paymentTrade.TradeQuery(&queryOrderParam)
```



订单关闭

```go
paymentTrade := payment.Payment{Client: aliPayClient}

closeOrderParam := payment.TradeClose{
   OutTradeNo:   "",
   TradeNo:      "",
}
tradeRes, err := paymentTrade.TradeClose(&closeOrderParam)
```



退款接口

```go
paymentTrade := payment.Payment{Client: aliPayClient}

refundOrderParam := payment.TradeRefund{
	OutTradeNo: "",
	TradeNo:    "",
}
tradeRes, err := paymentTrade.TradeRefund(&refundOrderParam)
```



退款查询

```go
paymentTrade := payment.Payment{Client: aliPayClient}

refundQueryOrderParam := payment.RefundQuery{
   OutTradeNo: "",
   TradeNo:    "",
}
tradeRes, err := paymentTrade.RefundQuery(&refundQueryOrderParam)
```



