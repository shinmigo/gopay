package payment

type GoodsDetail struct {
	GoodsId        string `json:"goods_id"`
	AliPayGoodsId  string `json:"alipay_goods_id,omitempty"`
	GoodsName      string `json:"goods_name"`
	Quantity       string `json:"quantity"`
	Price          string `json:"price"`
	GoodsCategory  string `json:"goods_category,omitempty"`
	CategoriesTree string `json:"categories_tree,omitempty"`
	Body           string `json:"body,omitempty"`
	ShowURL        string `json:"show_url,omitempty"`
}

type Trade struct {
	Subject           string         `json:"subject"`                   //商品的标题/交易标题/订单标题/订单关键字等。
	OutTradeNo        string         `json:"out_trade_no"`              //商户网站唯一订单号
	TotalAmount       string         `json:"total_amount"`              //订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]
	TimeoutExpress    string         `json:"timeout_express,omitempty"` //该笔订单允许的最晚付款时间，逾期将关闭交易。取值范围：5m～15d。m-分钟，h-小时，d-天，1c-当天（1c-当天的情况下，无论交易何时创建，都在0点关闭）。 该参数数值不接受小数点， 如 1.5h，可转换为 90m。
	TimeExpire        string         `json:"time_expire,omitempty"`     //绝对超时时间，格式为yyyy-MM-dd HH:mm。
	AuthToken         string         `json:"auth_token,omitempty"`
	GoodsType         string         `json:"goods_type,omitempty"` // 商品主类型：0—虚拟类商品，1—实物类商品 注：虚拟类商品不支持使用花呗渠道
	QuitUrl           string         `json:"quit_url,omitempty"`
	Body              string         `json:"body,omitempty"`            //对一笔交易的具体描述信息。如果是多种商品，请将商品描述字符串累加传给body
	PromoParams       string         `json:"promo_params,omitempty"`    // 优惠参数 注：仅与支付宝协商后可用
	PassbackParams    string         `json:"passback_params,omitempty"` //公用回传参数，如果请求时传递了该参数，则返回给商户时会回传该参数
	GoodsDetail       []*GoodsDetail `json:"goods_detail,omitempty"`
	EnablePayChannels string         `json:"enable_pay_channels,omitempty"` //可用渠道，用户只能在指定渠道范围内支付 当有多个渠道时用“,”分隔 注，与disable_pay_channels互斥
	StoreId           string         `json:"store_id,omitempty"`            // 商户门店编号。该参数用于请求参数中以区分各门店，非必传项
	SpecifiedChannel  string         `json:"specified_channel,omitempty"`   // 指定渠道，目前仅支持传入pcredit  若由于用户原因渠道不可用，用户可选择是否用其他渠道支付。  注：该参数不可与花呗分期参数同时传入
	BusinessParams    string         `json:"business_params,omitempty"`     // 商户传入业务信息，具体值要和支付宝约定，应用于安全，营销等参数直传场景，格式为json格式
}

type App struct {
	Trade
	NotifyUrl   string `json:"-"`                      //异步通知地址
	ReturnUrl   string `json:"-"`                      //支付返回地址
	ProductCode string `json:"product_code,omitempty"` //销售产品码，商家和支付宝签约的产品码
}

func (m *App) GetAliPayMethod() string {
	return "alipay.trade.app.pay"
}

type Wap struct {
	Trade
	NotifyUrl   string `json:"-"`                      //异步通知地址
	ReturnUrl   string `json:"-"`                      //支付返回地址
	ProductCode string `json:"product_code,omitempty"` //销售产品码，商家和支付宝签约的产品码
}

func (m *Wap) GetAliPayMethod() string {
	return "alipay.trade.wap.pay"
}

type Page struct {
	Trade
	NotifyUrl   string `json:"-"`                      //异步通知地址
	ReturnUrl   string `json:"-"`                      //支付返回地址
	ProductCode string `json:"product_code,omitempty"` //销售产品码，商家和支付宝签约的产品码
}

func (m *Page) GetAliPayMethod() string {
	return "alipay.trade.page.pay"
}

/**
 * 统一收单线下交易查询
 */
type TradeQuery struct {
	OutTradeNo   string   `json:"out_trade_no,omitempty"`  // 订单支付时传入的商户订单号, 与 TradeNo 二选一
	TradeNo      string   `json:"trade_no,omitempty"`      // 支付宝交易号
	OrgPid       string   `json:"org_pid,omitempty"`       //银行间联模式下有用，其它场景请不要使用； 双联通过该参数指定需要查询的交易所属收单机构的pid;
	QueryOptions []string `json:"query_options,omitempty"` // 可选 查询选项，商户通过上送该字段来定制查询返回信息 TRADE_SETTLE_INFO
}

func (m *TradeQuery) GetAliPayMethod() string {
	return "alipay.trade.query"
}

type FundBill struct {
	FundChannel string  `json:"fund_channel"`       // 交易使用的资金渠道，详见 支付渠道列表
	Amount      string  `json:"amount"`             // 该支付工具类型所使用的金额
	RealAmount  float64 `json:"real_amount,string"` // 渠道实际付款金额
}

type TradeSettleInfo struct {
	TradeSettleDetailList []*TradeSettleDetail `json:"trade_settle_detail_list"`
}

type TradeSettleDetail struct {
	OperationType     string `json:"operation_type"`
	OperationSerialNo string `json:"operation_serial_no"`
	OperationDate     string `json:"operation_dt"`
	TransOut          string `json:"trans_out"`
	TransIn           string `json:"trans_in"`
	Amount            string `json:"amount"`
}

type TradeQueryResContent struct {
	Code                string           `json:"code"`                        //网关返回码
	Msg                 string           `json:"msg"`                         //网关返回描述
	SubCode             string           `json:"sub_code"`                    //业务返回码
	SubMsg              string           `json:"sub_msg"`                     //业务返回码描述
	TradeNo             string           `json:"trade_no"`                    // 支付宝交易号
	OutTradeNo          string           `json:"out_trade_no"`                // 商家订单号
	BuyerLogonId        string           `json:"buyer_logon_id"`              // 买家支付宝账号
	TradeStatus         string           `json:"trade_status"`                // 交易状态
	TotalAmount         string           `json:"total_amount"`                // 交易的订单金额
	TransCurrency       string           `json:"trans_currency"`              // 标价币种
	SettleCurrency      string           `json:"settle_currency"`             // 订单结算币种
	SettleAmount        string           `json:"settle_amount"`               // 结算币种订单金额
	PayCurrency         string           `json:"pay_currency"`                // 订单支付币种
	PayAmount           string           `json:"pay_amount"`                  // 支付币种订单金额
	SettleTransRate     string           `json:"settle_trans_rate"`           // 结算币种兑换标价币种汇率
	TransPayRate        string           `json:"trans_pay_rate"`              // 标价币种兑换支付币种汇率
	BuyerPayAmount      string           `json:"buyer_pay_amount"`            // 买家实付金额，单位为元，两位小数。
	PointAmount         string           `json:"point_amount"`                // 积分支付的金额，单位为元，两位小数。
	InvoiceAmount       string           `json:"invoice_amount"`              // 交易中用户支付的可开具发票的金额，单位为元，两位小数。
	SendPayDate         string           `json:"send_pay_date"`               // 本次交易打款给卖家的时间
	ReceiptAmount       string           `json:"receipt_amount"`              // 实收金额，单位为元，两位小数
	StoreId             string           `json:"store_id"`                    // 商户门店编号
	TerminalId          string           `json:"terminal_id"`                 // 商户机具终端编号
	FundBillList        []*FundBill      `json:"fund_bill_list"`              // 交易支付使用的资金渠道
	StoreName           string           `json:"store_name"`                  // 请求交易支付中的商户店铺的名称
	BuyerUserId         string           `json:"buyer_user_id"`               // 买家在支付宝的用户id
	ChargeAmount        string           `json:"charge_amount"`               // 该笔交易针对收款方的收费金额；
	ChargeFlags         string           `json:"charge_flags"`                // 费率活动标识，当交易享受活动优惠费率时，返回该活动的标识；
	SettlementId        string           `json:"settlement_id"`               // 支付清算编号，用于清算对账使用；
	TradeSettleInfo     *TradeSettleInfo `json:"trade_settle_info,omitempty"` // 返回的交易结算信息，包含分账、补差等信息
	AuthTradePayMode    string           `json:"auth_trade_pay_mode"`         // 预授权支付模式，该参数仅在信用预授权支付场景下返回。信用预授权支付：CREDIT_PREAUTH_PAY
	BuyerUserType       string           `json:"buyer_user_type"`             // 买家用户类型。CORPORATE:企业用户；PRIVATE:个人用户。
	MdiscountAmount     string           `json:"mdiscount_amount"`            // 商家优惠金额
	DiscountAmount      string           `json:"discount_amount"`             // 平台优惠金额
	Subject             string           `json:"subject"`                     // 订单标题；
	Body                string           `json:"body"`                        // 订单描述;
	AlipaySubMerchantId string           `json:"alipay_sub_merchant_id"`      // 间连商户在支付宝端的商户编号；
	ExtInfos            string           `json:"ext_infos"`                   // 交易额外信息，特殊场景下与支付宝约定返回。
}

type TradeQueryRes struct {
	Body TradeQueryResContent `json:"alipay_trade_query_response"`
	Sign string               `json:"sign"`
}

/**
 * 统一收单交易关闭
 */
type TradeClose struct {
	NotifyUrl  string `json:"-"`                      //异步通知地址
	ReturnUrl  string `json:"-"`                      //支付返回地址
	TradeNo    string `json:"trade_no,omitempty"`     // 与 OutTradeNo 二选一
	OutTradeNo string `json:"out_trade_no,omitempty"` // 与 TradeNo 二选一
	OperatorId string `json:"operator_id,omitempty"`  // 可选
}

func (m *TradeClose) GetAliPayMethod() string {
	return "alipay.trade.close"
}

type TradeCloseRes struct {
	Body struct {
		Code       string `json:"code"`         //网关返回码
		Msg        string `json:"msg"`          //网关返回描述
		SubCode    string `json:"sub_code"`     //业务返回码
		SubMsg     string `json:"sub_msg"`      //业务返回码描述
		TradeNo    string `json:"trade_no"`     //支付宝交易号
		OutTradeNo string `json:"out_trade_no"` //创建交易传入的商户订单号
	} `json:"alipay_trade_close_response"`
	Sign string `json:"sign"`
}

/**
 * 统一收单交易退款接口
 */
type TradeRefund struct {
	NotifyUrl      string   `json:"-"`                         //异步通知地址
	ReturnUrl      string   `json:"-"`                         //支付返回地址
	OutTradeNo     string   `json:"out_trade_no,omitempty"`    // 与 TradeNo 二选一
	TradeNo        string   `json:"trade_no,omitempty"`        // 与 OutTradeNo 二选一
	RefundAmount   string   `json:"refund_amount"`             // 需要退款的金额，该金额不能大于订单金额,单位为元，支持两位小数
	RefundCurrency string   `json:"refund_currency,omitempty"` // 订单退款币种信息
	RefundReason   string   `json:"refund_reason,omitempty"`   // 退款的原因说明
	OutRequestNo   string   `json:"out_request_no,omitempty"`  // 标识一次退款请求，同一笔交易多次退款需要保证唯一，如需部分退款，则此参数必传。
	OperatorId     string   `json:"operator_id,omitempty"`     // 商户的操作员编号
	StoreId        string   `json:"store_id,omitempty"`        // 商户的门店编号
	TerminalId     string   `json:"terminal_id,omitempty"`     // 商户的终端编号
	QueryOptions   []string `json:"query_options,omitempty"`   //查询选项，商户通过上送该参数来定制同步需要额外返回的信息字段，数组格式。如：["refund_detail_item_list"]
}

func (m *TradeRefund) GetAliPayMethod() string {
	return "alipay.trade.refund"
}

type RefundDetailItem struct {
	FundChannel string `json:"fund_channel"` // 交易使用的资金渠道，详见 支付渠道列表
	BankCode    string `json:"bank_code"`    //银行卡支付时的银行代码
	Amount      string `json:"amount"`       // 该支付工具类型所使用的金额
	RealAmount  string `json:"real_amount"`  // 渠道实际付款金额
	FundType    string `json:"fund_type"`    //渠道所使用的资金类型,目前只在资金渠道(fund_channel)是银行卡渠道(BANKCARD)的情况下才返回该信息(DEBIT_CARD:借记卡,CREDIT_CARD:信用卡,MIXED_CARD:借贷合一卡)
}

type TradeRefundRes struct {
	Body struct {
		Code                         string              `json:"code"`
		Msg                          string              `json:"msg"`
		SubCode                      string              `json:"sub_code"`
		SubMsg                       string              `json:"sub_msg"`
		TradeNo                      string              `json:"trade_no"`                        // 支付宝交易号
		OutTradeNo                   string              `json:"out_trade_no"`                    // 商户订单号
		BuyerLogonId                 string              `json:"buyer_logon_id"`                  // 用户的登录id
		FundChange                   string              `json:"fund_change"`                     // 本次退款是否发生了资金变化
		RefundFee                    string              `json:"refund_fee"`                      // 退款总金额
		RefundCurrency               string              `json:"refund_currency"`                 // 退款币种信息
		GmtRefundPay                 string              `json:"gmt_refund_pay"`                  // 退款支付时间
		RefundDetailItemList         []*RefundDetailItem `json:"refund_detail_item_list"`         // 退款使用的资金渠道
		StoreName                    string              `json:"store_name"`                      // 交易在支付时候的门店名称
		BuyerUserId                  string              `json:"buyer_user_id"`                   // 买家在支付宝的用户id
		RefundSettlementId           string              `json:"refund_settlement_id"`            // 退款清算编号，用于清算对账使用；只在银行间联交易场景下返回该信息；
		PresentRefundBuyerAmount     string              `json:"present_refund_buyer_amount"`     // 本次退款金额中买家退款金额
		PresentRefundDiscountAmount  string              `json:"present_refund_discount_amount"`  // 本次退款金额中平台优惠退款金额
		PresentRefundMdiscountAmount string              `json:"present_refund_mdiscount_amount"` // 本次退款金额中商家优惠退款金额
		HasDepositBack               string              `json:"has_deposit_back"`                //是否有银行卡冲退
	} `json:"alipay_trade_refund_response"`
	Sign string `json:"sign"`
}

/**
 * 交易退款查询接口
 */
type RefundQuery struct {
	NotifyUrl    string   `json:"-"`                       //异步通知地址
	ReturnUrl    string   `json:"-"`                       //支付返回地址
	TradeNo      string   `json:"trade_no,omitempty"`      // 与 OutTradeNo 二选一
	OutTradeNo   string   `json:"out_trade_no,omitempty"`  // 与 TradeNo 二选一
	OutRequestNo string   `json:"out_request_no"`          // 请求退款接口时，传入的退款请求号，如果在退款请求时未传入，则该值为创建交易时的外部交易号
	OrgPid       string   `json:"org_pid,omitempty"`       //银行间联模式下有用，其它场景请不要使用； 双联通过该参数指定需要查询的交易所属收单机构的pid;
	QueryOptions []string `json:"query_options,omitempty"` // 查询选项，商户通过上送该参数来定制同步需要额外返回的信息字段，数组格式。 refund_detail_item_list
}

func (m *RefundQuery) GetAliPayMethod() string {
	return "alipay.trade.fastpay.refund.query"
}

type RefundQueryRes struct {
	Body struct {
		Code                 string              `json:"code"`
		Msg                  string              `json:"msg"`
		SubCode              string              `json:"sub_code"`
		SubMsg               string              `json:"sub_msg"`
		TradeNo              string              `json:"trade_no"`                // 支付宝交易号
		OutTradeNo           string              `json:"out_trade_no"`            // 创建交易传入的商户订单号
		OutRequestNo         string              `json:"out_request_no"`          // 本笔退款对应的退款请求号
		RefundReason         string              `json:"refund_reason"`           // 发起退款时，传入的退款原因
		TotalAmount          string              `json:"total_amount"`            // 发该笔退款所对应的交易的订单金额
		RefundAmount         string              `json:"refund_amount"`           // 本次退款请求，对应的退款金额
		RefundDetailItemList []*RefundDetailItem `json:"refund_detail_item_list"` // 本次退款使用的资金渠道；
	} `json:"alipay_trade_fastpay_refund_query_response"`
	Sign string `json:"sign"`
}

/**
 * 异步通知参数
 */
type Notify struct {
	NotifyTime        string  `json:"notify_time"`         //通知的发送时间。格式为yyyy-MM-dd HH:mm:ss
	NotifyType        string  `json:"notify_type"`         //通知的类型
	NotifyId          string  `json:"notify_id"`           //通知校验ID
	Charset           string  `json:"charset"`             //编码格式，如utf-8、gbk、gb2312等
	Version           string  `json:"version"`             //调用的接口版本，固定为：1.0
	SignType          string  `json:"sign_type"`           //签名算法类型，目前支持RSA2和RSA，推荐使用RSA2
	Sign              string  `json:"sign"`                //请参考异步返回结果的验签
	AuthAppId         string  `json:"auth_app_id"`         //授权方的appid，由于本接口暂不开放第三方应用授权，因此auth_app_id=app_id
	TradeNo           string  `json:"trade_no"`            //支付宝交易号
	AppId             string  `json:"app_id"`              //开发者APP_ID
	OutTradeNo        string  `json:"out_trade_no"`        //商户订单号
	OutBizNo          string  `json:"out_biz_no"`          //商户业务号
	BuyerId           string  `json:"buyer_id"`            //买家支付宝用户号
	SellerId          string  `json:"seller_id"`           //卖家支付宝用户号
	TradeStatus       string  `json:"trade_status"`        //交易状态
	TotalAmount       float64 `json:"total_amount"`        //订单金额
	ReceiptAmount     float64 `json:"receipt_amount"`      //实收金额
	InvoiceAmount     float64 `json:"invoice_amount"`      //开票金额
	BuyerPayAmount    float64 `json:"buyer_pay_amount"`    //付款金额
	PointAmount       float64 `json:"point_amount"`        //集分宝金额
	RefundFee         float64 `json:"refund_fee"`          //总退款金额
	Subject           string  `json:"subject"`             //订单标题
	Body              string  `json:"body"`                //商品描述
	GmtCreate         string  `json:"gmt_create"`          //交易创建时间
	GmtPayment        string  `json:"gmt_payment"`         //交易付款时间
	GmtRefund         string  `json:"gmt_refund"`          //交易退款时间
	GmtClose          string  `json:"gmt_close"`           //交易结束时间
	FundBillList      string  `json:"fund_bill_list"`      //支付金额信息
	VoucherDetailList string  `json:"voucher_detail_list"` //优惠券信息
	PassbackParams    string  `json:"passback_params"`     //回传参数
}
