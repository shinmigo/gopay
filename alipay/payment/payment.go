package payment

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/shinmigo/gopay/alipay/kernel"
)

type Payment struct {
	Client *kernel.AliPayClient
}

/*
 * APP支付
 */
func (m *Payment) App(param *App) (string, error) {
	if param == nil {
		return "", errors.New(kernel.InitializeDataErr)
	}

	urlMap, err := m.Client.UrlParams(param)
	if err != nil {
		return "", err
	}
	return urlMap.Encode(), nil
}

/*
 * 手机网站支付
 */
func (m *Payment) Wap(param *Wap) (string, error) {
	if param == nil {
		return "", errors.New(kernel.InitializeDataErr)
	}

	param.ProductCode = "QUICK_WAP_WAY"
	urlMap, err := m.Client.UrlParams(param)
	if err != nil {
		return "", err
	}
	return buildForm(m.Client.GetGatewayHost(), urlMap), nil
}

/*
 * PC网站支付
 */
func (m *Payment) Page(param *Page) (string, error) {
	if param == nil {
		return "", errors.New(kernel.InitializeDataErr)
	}

	param.ProductCode = "FAST_INSTANT_TRADE_PAY"
	urlMap, err := m.Client.UrlParams(param)
	if err != nil {
		return "", err
	}
	return buildForm(m.Client.GetGatewayHost(), urlMap), nil
}

/**
 * 统一收单线下交易查询
 */
func (m *Payment) TradeQuery(param *TradeQuery) (result *TradeQueryRes, err error) {
	if param == nil {
		return nil, errors.New(kernel.InitializeDataErr)
	}

	err = m.Client.SendRequest("POST", param, &result)
	return result, err
}

/**
 * 统一收单交易关闭
 */
func (m *Payment) TradeClose(param *TradeClose) (result *TradeCloseRes, err error) {
	if param == nil {
		return nil, errors.New(kernel.InitializeDataErr)
	}

	err = m.Client.SendRequest("POST", param, &result)
	return result, err
}

/**
 * 统一收单交易退款接口
 */
func (m *Payment) TradeRefund(param *TradeRefund) (result *TradeRefundRes, err error) {
	if param == nil {
		return nil, errors.New(kernel.InitializeDataErr)
	}

	err = m.Client.SendRequest("POST", param, &result)
	return result, err
}

/**
 * 交易退款查询接口
 */
func (m *Payment) RefundQuery(param *RefundQuery) (result *RefundQueryRes, err error) {
	if param == nil {
		return nil, errors.New(kernel.InitializeDataErr)
	}

	err = m.Client.SendRequest("POST", param, &result)
	return result, err
}

/*
 *生成页面类请求所需URL或Form表单
 */
func buildForm(actionUrl string, parameters url.Values) string {
	buffers := &bytes.Buffer{}
	buffers.WriteString(fmt.Sprintf("<form id='alipaysubmit' name='alipaysubmit' action='%s?charset=%s' method='POST'>", actionUrl, kernel.AliPayCharset))

	for name := range parameters {
		value := parameters.Get(name)
		if value == "" {
			continue
		}
		value = strings.ReplaceAll(strings.Trim(value, " "), "'", "&apos;")
		buffers.WriteString(fmt.Sprintf("<input type='hidden' name='%s' value='%s'/>", name, value))
	}

	buffers.WriteString("<input type='submit' value='ok' style='display:none;'/></form>")
	buffers.WriteString("<script>document.forms['alipaysubmit'].submit();</script>")
	return buffers.String()
}
