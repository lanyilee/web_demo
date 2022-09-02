package errpkg

// A 表示错误来源于用户，比如参数错误，用户安装版本过低，用户支付超时等问题
// B 表示错误来源于当前系统，往往是业务逻辑出错，或程序健壮性差等问题
// C 表示错误来源于第三方服务，比如 CDN 服务出错，消息投递超时等问题

var (
	OK                      = &Errno{"00000", "成功"}
	UserError               = &Errno{"A0001", "用户端错误"}     //一级宏观错误码
	SystemError             = &Errno{"B0001", "系统执行出错"}    //一级宏观错误码
	InvokeOtherServiceError = &Errno{"C0001", "调用第三方服务出错"} //一级宏观错误码
)

var (
	UserRegistryError   = &Errno{"A0100", "用户注册错误"} //二级宏观错误码
	NotAcceptProtocol   = &Errno{"A0101", "用户未同意隐私协议"}
	RegistryAreaLimit   = &Errno{"A0102", "注册国家或地区受限"}
	UsernameCheckError  = &Errno{"A0110", "用户名校验失败"}
	UsernameExist       = &Errno{"A0111", "用户名已存在"}
	UsernameSensitive   = &Errno{"A0112", "用户名包含敏感词"}
	UsernameSpecialChar = &Errno{"A0113", "用户名包含特殊字符"}
	PasswordCheckError  = &Errno{"A0120", "密码校验失败"}
	PasswordTooShort    = &Errno{"A0121", "密码长度不够"}
	PasswordNotStrength = &Errno{"A0122", "密码强度不够"}
	VerifyCodeError     = &Errno{"A0130", "校验码输入错误"}
	SmsCodeError        = &Errno{"A0131", "短信校验码输入错误"}
	EmailCodeError      = &Errno{"A0132", "邮件校验码输入错误"}
	IvrCodeError        = &Errno{"A0133", "语音校验码输入错误"}
	IdError             = &Errno{"A0140", "用户证件异常"}
	NotSelectIdType     = &Errno{"A0141", "用户证件类型未选择"}
	NidCheckError       = &Errno{"A0142", "大陆身份证编号校验非法"}
	PassportError       = &Errno{"A0143", "护照编号校验非法"}
	MidCard             = &Errno{"A0144", "军官证编号校验非法"}
	BaseInfError        = &Errno{"A0150", "用户基本信息校验失败"}
	CellphoneError      = &Errno{"A0151", "手机格式校验失败"}
	AddressError        = &Errno{"A0152", "地址格式校验失败"}
	MailError           = &Errno{"A0153", "邮箱格式校验失败"}

	LoginError             = &Errno{"A0200", "用户登录异常"} //二级宏观错误码
	UserNotExist           = &Errno{"A0201", "用户账户不存在"}
	UserStop               = &Errno{"A0202", "用户账户被冻结"}
	UserHasCancel          = &Errno{"A0203", "用户账户已作废"}
	UsernamePasswordError  = &Errno{"A0210", "用户密码错误"}
	PasswordExceedsLimit   = &Errno{"A0211", "用户输入密码错误次数超限"}
	UserIdentityCheckError = &Errno{"A0220", "用户身份校验失败"}
	FingerIdError          = &Errno{"A0221", "用户指纹识别失败"}
	FaceIdError            = &Errno{"A0222", "用户面容识别失败"}
	Oauth2NotPermit        = &Errno{"A0223", "用户未获得第三方登录授权"}
	UserLoginExpire        = &Errno{"A0230", "用户登录已过期"}
	UserVerifyCodeError    = &Errno{"A0240", "用户验证码错误"}
	GetUserVerifyCodeLimit = &Errno{"A0241", "用户验证码尝试次数超限"}
	UserOffline            = &Errno{"A0242", "用户不在线"}

	PermitError              = &Errno{"A0300", "访问权限异常"} //二级宏观错误码
	NotPermit                = &Errno{"A0301", "访问未授权"}
	PERMITTING               = &Errno{"A0302", "正在授权中"}
	UserAuthBeRejected       = &Errno{"A0303", "用户授权申请被拒绝"}
	NotAccessPrivateResource = &Errno{"A0310", "因访问对象隐私设置被拦截"}
	PermitExpire             = &Errno{"A0311", "授权已过期"}
	NotPermitApi             = &Errno{"A0312", "无权限使用 API"}
	VisitIntercept           = &Errno{"A0320", "用户访问被拦截"}
	IsUserBlackList          = &Errno{"A0321", "黑名单用户"}
	IllegalIp                = &Errno{"A0323", "非法 IP 地址"}
	GatewayAccessLimit       = &Errno{"A0324", "网关访问受限"}
	IsAreaBlackList          = &Errno{"A0325", "地域黑名单"}
	ServerArrears            = &Errno{"A0330", "服务已欠费"}
	UserSignError            = &Errno{"A0340", "用户签名异常"}
	RsaSignError             = &Errno{"A0341", "RSA 签名错误"}

	ParamsError            = &Errno{"A0400", "用户请求参数错误"} //二级宏观错误码
	HaveSpiteRedirectUrl   = &Errno{"A0401", "包含非法恶意跳转链接"}
	InvalidInput           = &Errno{"A0402", "无效的用户输入"}
	RequireParamsHaveEmpty = &Errno{"A0410", "请求必填参数为空"}
	OrderIdEmpty           = &Errno{"A0411", "用户订单号为空"}
	PurchaseNumIsEmpty     = &Errno{"A0412", "订购数量为空"}
	LackTimestamp          = &Errno{"A0413", "缺少时间戳参数"}
	InvalidTimestamp       = &Errno{"A0414", "非法的时间戳参数"}
	ParamsValueOverLimit   = &Errno{"A0420", "请求参数值超出允许的范围"}
	ParamsTypeNotMatch     = &Errno{"A0421", "参数格式不匹配"}
	NotAtServerArea        = &Errno{"A0422", "地址不在服务范围"}
	NotAtServerTime        = &Errno{"A0423", "时间不在服务范围"}
	AmountOverLimit        = &Errno{"A0424", "金额超出限制"}
	LoginNumLimit          = &Errno{"A0425", "数量超出限制"}
	BatchOverLimit         = &Errno{"A0426", "请求批量处理总个数超出限制"}
	JsonError              = &Errno{"A0427", "请求 JSON 解析失败"}
	IllegalityInput        = &Errno{"A0430", "用户输入内容非法"}
	HaveSensitiveChar      = &Errno{"A0431", "包含违禁敏感词"}
	ImgHaveIllegalityInf   = &Errno{"A0432", "图片包含违禁信息"}
	FilePiracy             = &Errno{"A0433", "文件侵犯版权"}
	UserOperationError     = &Errno{"A0440", "用户操作异常"}
	UserPayOverTime        = &Errno{"A0441", "用户支付超时"}
	ConfirmOrderOverTime   = &Errno{"A0442", "确认订单超时"}
	OrderClosed            = &Errno{"A0443", "订单已关闭"}

	UserRequestError          = &Errno{"A0500", "用户请求服务异常"} //二级宏观错误码
	RequestNumOverLimit       = &Errno{"A0501", "请求次数超出限制"}
	ParallelRequestLimit      = &Errno{"A0502", "请求并发数超出限制"}
	WaitForAction             = &Errno{"A0503", "用户操作请等待"}
	WebsocketConnectionError  = &Errno{"A0504", "WebSocket 连接异常"}
	WebsocketConnectionClosed = &Errno{"A0505", "WebSocket 连接断开"}
	RepeatRequest             = &Errno{"A0506", "用户重复请求"}

	UserResourceError       = &Errno{"A0600", "用户资源异常"} //二级宏观错误码
	BalanceNotSufficient    = &Errno{"A0601", "账户余额不足"}
	UserDiskNotSufficient   = &Errno{"A0602", "用户磁盘空间不足"}
	UserMemoryNotSufficient = &Errno{"A0603", "用户内存空间不足"}
	UserOssNotSufficient    = &Errno{"A0604", "用户 OSS 容量不足"}
	ResourceAllUsed         = &Errno{"A0605", "用户配额已用光"}

	UploadError         = &Errno{"A0700", "用户上传文件异常"} //二级宏观错误码
	UploadFileTypeError = &Errno{"A0701", "用户上传文件类型不匹配"}
	UploadFileTooBig    = &Errno{"A0702", "用户上传文件太大"}
	UploadImgTooBig     = &Errno{"A0703", "用户上传图片太大"}
	UploadVideoTooBig   = &Errno{"A0704", "用户上传视频太大"}
	UploadZipTooBig     = &Errno{"A0705", "用户上传压缩文件太大"}

	VersionError        = &Errno{"A0800", "用户当前版本异常"} //二级宏观错误码
	VersionNotFit       = &Errno{"A0801", "用户安装版本与系统不匹配"}
	VersionTooLow       = &Errno{"A0802", "用户安装版本过低"}
	VersionTooHeight    = &Errno{"A0803", "用户安装版本过高"}
	VersionExpire       = &Errno{"A0804", "用户安装版本已过期"}
	ApiVersionNotFit    = &Errno{"A0805", "用户 API 请求版本不匹配"}
	ApiVersionTooHeight = &Errno{"A0806", "用户 API 请求版本过高"}
	ApiVersionTooLow    = &Errno{"A0807", "用户 API 请求版本过低"}

	PrivacyError         = &Errno{"A0900", "用户隐私未授权"} //二级宏观错误码
	PrivacyNoSign        = &Errno{"A0901", "用户隐私未签署"}
	CameraNotPermit      = &Errno{"A0902", "用户摄像头未授权"}
	PhotoAlbumNotPermit  = &Errno{"A0904", "用户图片库未授权"}
	FileNotPermit        = &Errno{"A0905", "用户文件未授权"}
	LocationNotPermit    = &Errno{"A0906", "用户位置信息未授权"}
	AddressBookNotPermit = &Errno{"A0907", "用户通讯录未授权"}

	EquipmentError        = &Errno{"A1000", "用户设备异常"} //二级宏观错误码
	CameraError           = &Errno{"A1001", "用户相机异常"}
	MicError              = &Errno{"A1002", "用户麦克风异常"}
	TelephoneReceiveError = &Errno{"A1003", "用户听筒异常"}
	LoudspeakerError      = &Errno{"A1004", "用户扬声器异常"}
	GpsError              = &Errno{"A1005", "用户 GPS 定位异常"}

	SystemProcessError    = &Errno{"B0100", "系统执行超时"} //二级宏观错误码
	OrderProcessVoterTime = &Errno{"B0101", "系统订单处理超时"}

	SystemDisasterError       = &Errno{"B0200", "系统容灾功能被触发"} //二级宏观错误码
	SystemCurrentLimiting     = &Errno{"B0210", "系统限流"}
	SystemFunctionDegradation = &Errno{"B0220", "系统功能降级"}

	SystemResourceError       = &Errno{"B0300", "系统资源异常"} //二级宏观错误码
	SystemResourceUseUp       = &Errno{"B0310", "系统资源耗尽"}
	SystemDiskUseUp           = &Errno{"B0311", "系统磁盘空间耗尽"}
	SystemMemoryUseUp         = &Errno{"B0312", "系统内存耗尽"}
	SystemFileHandlerUseUp    = &Errno{"B0313", "文件句柄耗尽"}
	SystemConnectionPoolUseUp = &Errno{"B0314", "系统连接池耗尽"}
	SystemThreadPoolUseUp     = &Errno{"B0315", "系统线程池耗尽"}
	SystemVisitResourceError  = &Errno{"B0320", "系统资源访问异常"}
	SystemVisitDiskFileError  = &Errno{"B0321", "系统读取磁盘文件失败"}

	MiddlewareError              = &Errno{"C0100", "中间件服务出错"} //二级宏观错误码
	RpcServiceError              = &Errno{"C0110", "RPC 服务出错"}
	RpcServiceNotFound           = &Errno{"C0111", "RPC 服务未找到"}
	RpcServiceNotRegistry        = &Errno{"C0112", "RPC 服务未注册"}
	ApiNotExist                  = &Errno{"C0113", "接口不存在"}
	MessageServiceError          = &Errno{"C0120", "消息服务出错"}
	MessageServicePushError      = &Errno{"C0121", "消息投递出错"}
	MessageServiceConsumeError   = &Errno{"C0122", "消息消费出错"}
	MessageServiceSubscribeError = &Errno{"C0123", "消息订阅出错"}
	MessageServiceGourpNotFound  = &Errno{"C0124", "消息分组未查到"}
	CacheServiceError            = &Errno{"C0130", "缓存服务出错"}
	KeyOverLimit                 = &Errno{"C0131", "key 长度超过限制"}
	ValueOverLimit               = &Errno{"C0132", "value 长度超过限制"}
	StorageMemoryUseUp           = &Errno{"C0133", "存储容量已满"}
	NotSupportThisDataType       = &Errno{"C0134", "不支持的数据格式"}
	ConfigServiceError           = &Errno{"C0140", "配置服务出错"}
	NetworkResourceServiceError  = &Errno{"C0150", "网络资源服务出错"}
	VpcServiceError              = &Errno{"C0151", "VPN 服务出错"}
	CdnServiceError              = &Errno{"C0152", "CDN 服务出错"}
	DomainNameServiceError       = &Errno{"C0153", "域名解析服务出错"}
	GatewayError                 = &Errno{"C0154", "网关服务出错"}

	OtherProcessOverTime  = &Errno{"C0200", "第三方系统执行超时"} //二级宏观错误码
	RpcOverTime           = &Errno{"C0210", "RPC 执行超时"}
	MessagePushOverTime   = &Errno{"C0220", "消息投递超时"}
	CacheServiceOverTime  = &Errno{"C0230", "缓存服务超时"}
	ConfigServiceOverTime = &Errno{"C0240", "配置服务超时"}
	DatabaseOverTime      = &Errno{"C0250", "数据库服务超时"}

	DatabaseError                  = &Errno{"C0300", "数据库服务出错"} //二级宏观错误码
	TableNotExist                  = &Errno{"C0311", "表不存在"}
	ColumnNotExist                 = &Errno{"C0312", "列不存在"}
	MultiTableSelectHaveSameColumn = &Errno{"C0321", "多表关联中存在多个相同名称的列"}
	DatabaseDeadlock               = &Errno{"C0331", "数据库死锁"}
	PrimaryKeyConflict             = &Errno{"C0341", "主键冲突"}

	OtherSystemDisaster            = &Errno{"C0400", "第三方容灾系统被触发"} //二级宏观错误码
	OtherSystemCurrentLimiting     = &Errno{"C0401", "第三方系统限流"}
	OtherSystemFunctionDegradation = &Errno{"C0402", "第三方功能降级"}

	NotifyServiceError   = &Errno{"C0500", "通知服务出错"} //二级宏观错误码
	SmsNotifyServiceFail = &Errno{"C0501", "短信提醒服务失败"}
	IvrNotifyServiceFail = &Errno{"C0502", "语音提醒服务失败"}
	MailNotifyFail       = &Errno{"C0503", "邮件提醒服务失败"}
)
