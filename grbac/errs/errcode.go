package errs

// 定义公共错误码
var (
	// 常见
	StatusOK            = NewAppError(10000, "请求成功")
	Created             = NewAppError(10201, "创建成功")
	Accepted            = NewAppError(10202, "更新成功")
	BadRequest          = NewAppError(10400, "入参错误")
	Unauthorized        = NewAppError(10401, "凭证无效")
	Forbidden           = NewAppError(10403, "禁止访问")
	NotFound            = NewAppError(10404, "未找到")
	MethodNotAllowed    = NewAppError(10405, "请求方法不支持")
	RequestTimeout      = NewAppError(10408, "请求超时")
	UnprocessableEntity = NewAppError(10422, "无效字段")
	TooManyRequests     = NewAppError(10429, "请求频繁")
	InternalServerError = NewAppError(10500, "服务端错误")

	// 其他
	AlreadyExist = NewAppError(20001, "数据已存在")
)
