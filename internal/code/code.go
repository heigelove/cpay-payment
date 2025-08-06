package code

import (
	_ "embed"

	"github.com/heigelove/cpay-payment/configs"
)

//go:embed code.go
var ByteCodeFile []byte

// Failure 错误时返回结构
type Failure struct {
	Code    int    `json:"code"`    // 业务码
	Message string `json:"message"` // 描述信息
}

const (
	StatsSecret = "6Hw5n3vQdXgcH17l4WDR6wz2t9cNBQ9n" // 统计用的安全令牌

)

const (
	ServerError        = 10101
	TooManyRequests    = 10102
	ParamBindError     = 10103
	AuthorizationError = 10104
	UrlSignError       = 10105
	CacheSetError      = 10106
	CacheGetError      = 10107
	CacheDelError      = 10108
	CacheNotExist      = 10109
	ResubmitError      = 10110
	HashIdsEncodeError = 10111
	HashIdsDecodeError = 10112
	RBACError          = 10113
	RedisConnectError  = 10114
	MySQLConnectError  = 10115
	WriteConfigError   = 10116
	SendEmailError     = 10117
	MySQLExecError     = 10118
	GoVersionError     = 10119
	SocketConnectError = 10120
	SocketSendError    = 10121
	InvalidSecureToken = 10122 // 无效的安全令牌
	MerchantNoRequired = 10123 // 商户号不能为空
	SettleOrderError   = 10124 // 结算订单错误
	OrderNotFound      = 10125 // 订单未找到

	AuthorizedCreateError    = 20101
	AuthorizedListError      = 20102
	AuthorizedDeleteError    = 20103
	AuthorizedUpdateError    = 20104
	AuthorizedDetailError    = 20105
	AuthorizedCreateAPIError = 20106
	AuthorizedListAPIError   = 20107
	AuthorizedDeleteAPIError = 20108

	AdminCreateError             = 20201
	AdminListError               = 20202
	AdminDeleteError             = 20203
	AdminUpdateError             = 20204
	AdminResetPasswordError      = 20205
	AdminLoginError              = 20206
	AdminLogOutError             = 20207
	AdminModifyPasswordError     = 20208
	AdminModifyPersonalInfoError = 20209
	AdminMenuListError           = 20210
	AdminMenuCreateError         = 20211
	AdminOfflineError            = 20212
	AdminDetailError             = 20213
	AdminGoogleAuthError         = 20214
	AdminGoogleAuthEnableError   = 20215
	AdminGoogleAuthDisableError  = 20216

	MenuCreateError       = 20301
	MenuUpdateError       = 20302
	MenuListError         = 20303
	MenuDeleteError       = 20304
	MenuDetailError       = 20305
	MenuCreateActionError = 20306
	MenuListActionError   = 20307
	MenuDeleteActionError = 20308

	CronCreateError  = 20401
	CronUpdateError  = 20402
	CronListError    = 20403
	CronDetailError  = 20404
	CronExecuteError = 20405

	AgentCreateError           = 20501
	AgentUpdateError           = 20502
	AgentListError             = 20503
	AgentDeleteError           = 20504
	AgentDetailError           = 20505
	AgentNameExistError        = 20506
	OperationAccountExistError = 20507
	PayinShareError            = 20508
	PayoutShareError           = 20509
	ExchangeRateError          = 20510
	AgentStatusError           = 20511
	AgentIdError               = 20512

	MerchantCreateError       = 20601
	MerchantUpdateError       = 20602
	MerchantListError         = 20603
	MerchantDeleteError       = 20604
	MerchantDetailError       = 20605
	MerchantIdError           = 20606
	MerchantNameExistError    = 20607
	MerchantStatusError       = 20608
	MerchantTypeError         = 20609
	MerchantSecretError       = 20610
	MerchantPayinShareError   = 20611
	MerchantPayoutShareError  = 20612
	MerchantExchangeRateError = 20613

	AgentSettlementCreateError = 20701
	AgentSettlementUpdateError = 20702
	AgentSettlementListError   = 20703
	AgentSettlementDeleteError = 20704
	AgentSettlementDetailError = 20705

	CountryCreateError    = 20801
	CountryUpdateError    = 20802
	CountryListError      = 20803
	CountryDeleteError    = 20804
	CountryDetailError    = 20805
	CountryIdError        = 20806
	CountryNameExistError = 20807
	CountryStatusError    = 20808

	// MerchantChannel 相关错误码
	MerchantChannelCreateError = 20809
	MerchantChannelUpdateError = 20810
	MerchantChannelDetailError = 20811

	// 其他错误码可以继续添加
	ChannelCreateError    = 20901
	ChannelUpdateError    = 20902
	ChannelListError      = 20903
	ChannelDeleteError    = 20904
	ChannelDetailError    = 20905
	ChannelIdError        = 20906
	ChannelNameExistError = 20907
	ChannelStatusError    = 20908
	ChannelTypeError      = 20909

	ChannelGroupCreateError    = 21001
	ChannelGroupUpdateError    = 21002
	ChannelGroupListError      = 21003
	ChannelGroupDeleteError    = 21004
	ChannelGroupDetailError    = 21005
	ChannelGroupIdError        = 21006
	ChannelGroupNameExistError = 21007
	ChannelGroupStatusError    = 21008

	ChannelGroupConfigCreateError = 21101
	ChannelGroupConfigUpdateError = 21102
	ChannelGroupConfigListError   = 21103
	ChannelGroupConfigDeleteError = 21104
	ChannelGroupConfigDetailError = 21105
	ChannelGroupConfigIdError     = 21106
	ChannelGroupConfigStatusError = 21108

	PayinDailyStatsCreateError        = 21201
	PayinDailyStatsUpdateError        = 21202
	PayinDailyStatsListError          = 21203
	PayinDailyStatsDeleteError        = 21204
	PayinDailyStatsDetailError        = 21205
	PayinDailyStatsIdError            = 21206
	PayinDailyStatsStatsDateError     = 21207
	PayinDailyStatsChannelError       = 21208
	PayinDailyStatsMerchantError      = 21209
	PayinDailyStatsCountryError       = 21210
	PayinDailyStatsExportError        = 21211
	PayinDailyStatsMerchantStatsError = 21212

	PayoutDailyStatsCreateError        = 21300
	PayoutDailyStatsUpdateError        = 21301
	PayoutDailyStatsListError          = 21302
	PayoutDailyStatsDeleteError        = 21303
	PayoutDailyStatsDetailError        = 21304
	PayoutDailyStatsIdError            = 21305
	PayoutDailyStatsStatsDateError     = 21306
	PayoutDailyStatsChannelError       = 21307
	PayoutDailyStatsMerchantError      = 21308
	PayoutDailyStatsCountryError       = 21309
	PayoutDailyStatsExportError        = 21310
	PayoutDailyStatsMerchantStatsError = 21311

	PayinOrderCreateError        = 21400
	PayinOrderUpdateError        = 21401
	PayinOrderListError          = 21402
	PayinOrderDeleteError        = 21403
	PayinOrderDetailError        = 21404
	PayinOrderIdError            = 21405
	PayinOrderMerchantError      = 21406
	PayinOrderChannelError       = 21407
	PayinOrderCountryError       = 21408
	PayinOrderExportError        = 21409
	PayinOrderMerchantStatsError = 21410
	PayinOrderStatusError        = 21411
	PayinOrderNotifyError        = 21412
	PayinOrderNotFoundError      = 21413
	PayinOrderStatisticsError    = 21414
	PayinOrderDailyStatsError    = 21415
	PayinOrderQueryStatusError   = 21416 // 查询代收订单状态错误

	PayoutOrderCreateError        = 21500
	PayoutOrderUpdateError        = 21501
	PayoutOrderListError          = 21502
	PayoutOrderDeleteError        = 21503
	PayoutOrderDetailError        = 21504
	PayoutOrderIdError            = 21505
	PayoutOrderMerchantError      = 21506
	PayoutOrderChannelError       = 21507
	PayoutOrderCountryError       = 21508
	PayoutOrderExportError        = 21509
	PayoutOrderMerchantStatsError = 21510
	PayoutOrderStatusError        = 21511
	PayoutOrderNotifyError        = 21512
	PayoutOrderNotFoundError      = 21513
	PayoutOrderStatisticsError    = 21514
	PayoutOrderQueryStatusError   = 21515 // 查询代付订单状态错误

	DailySummaryCreateError              = 21600
	DailySummaryListError                = 21601
	DailySummaryUpdateError              = 21602
	DailySummaryDeleteError              = 21603
	DailySummaryExportError              = 21604
	DailySummaryIdError                  = 21605
	DailySummaryTypeError                = 21606
	DailySummaryBalanceError             = 21607
	DailySummaryPayinCountError          = 21608
	DailySummaryPayinAmountError         = 21609
	DailySummaryPayinSuccessCountError   = 21610
	DailySummaryPayinSuccessAmountError  = 21611
	DailySummaryPayoutCountError         = 21612
	DailySummaryPayoutAmountError        = 21613
	DailySummaryPayoutSuccessCountError  = 21614
	DailySummaryPayoutSuccessAmountError = 21615
	DailySummaryPayoutFailCountError     = 21616
	DailySummaryPayoutFailAmountError    = 21617
	DailySummaryReportError              = 21416 // 日报表错误

	FinanceRecordCreateError          = 21700
	FinanceRecordListError            = 21701
	FinanceRecordUpdateError          = 21702
	FinanceRecordDeleteError          = 21703
	FinanceRecordDetailError          = 21704
	FinanceRecordIdError              = 21705
	FinanceRecordTypeError            = 21706
	FinanceRecordMerchantError        = 21707
	FinanceRecordAgentError           = 21708
	FinanceRecordExportError          = 21709
	FinanceRecordNotifyError          = 21710
	FinanceRecordNotFoundError        = 21711
	FinanceRecordStatsError           = 21712
	FinanceRecordMerchantBalanceError = 21713

	WalletCreateError    = 21801
	WalletUpdateError    = 21802
	WalletListError      = 21803
	WalletDeleteError    = 21804
	WalletDetailError    = 21805
	WalletIdError        = 21806
	WalletNameExistError = 21807
	WalletStatusError    = 21808

	WithdrawalCreateError       = 21901
	WithdrawalUpdateError       = 21902
	WithdrawalListError         = 21903
	WithdrawalDeleteError       = 21904
	WithdrawalDetailError       = 21905
	WithdrawalIdError           = 21906
	WithdrawalStatisticsError   = 21907
	WithdrawalUpdateStatusError = 21908
	WithdrawalExportError       = 21909
	WithdrawalReviewError       = 21910

	AccountListError   = 22001
	AccountAdjustError = 22002 // 商户账户余额调整错误
	AccountExportError = 22003 // 商户账户导出错误
	AccountDetailError = 22004 // 商户账户详情错误

	UploadProofError = 22101 // 上传凭证错误

	// 结算相关错误码 22201-22299
	SettlementInvalidParams  = 22201 // 结算参数错误
	SettlementOrderNotFound  = 22202 // 结算订单不存在
	SettlementStatusInvalid  = 22203 // 订单状态不允许结算
	SettlementAlreadySettled = 22204 // 订单已结算
	SettlementDBOperation    = 22205 // 结算数据库操作失败

	WhiteListCreateError  = 22301
	WhiteListListError    = 22302
	WhiteListDeleteError  = 22303
	WhiteListIpExistError = 22304
)

func Text(code int) string {
	lang := configs.Get().Language.Local

	if lang == configs.ZhCN {
		return zhCNText[code]
	}

	if lang == configs.EnUS {
		return enUSText[code]
	}

	return zhCNText[code]
}
