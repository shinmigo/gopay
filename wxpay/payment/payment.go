package payment

import (
	"encoding/xml"

	"github.com/shinmigo/gopay/wxpay/kernel"
)

type Payment struct {
	Client *kernel.WxClient
}

/**
 * 微信统一下单
 */
func (m *Payment) Pay(param *Trade) (result *TradeRes, err error) {
	if param == nil {
		return nil, nil
	}

	err = m.Client.SendRequest("POST", "pay/unifiedorder", param, &result)
	return
}

/**
 * 微信查询订单
 */
func (m *Payment) Query(param *TradeQuery) (result *TradeQueryRes, err error) {
	if param == nil {
		return nil, nil
	}

	err = m.Client.SendRequest("POST", "pay/orderquery", param, &result)
	return
}

/**
 * 微信关闭订单
 */
func (m *Payment) Close(param *TradeClose) (result *TradeCloseRes, err error) {
	if param == nil {
		return nil, nil
	}

	err = m.Client.SendRequest("POST", "pay/closeorder", param, &result)
	return
}

/**
 * 异步通知验证签名
 */
func (m *Payment) NotifyVerify(reqBody []byte) (res *NotifyRes, err error) {
	err = m.Client.VerifySign(reqBody)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(reqBody, &res)
	return res, err
}
