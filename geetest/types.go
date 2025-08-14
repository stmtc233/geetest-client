// geetest/types.go
package geetest

// ApiResponse 是一个通用包装器，用于封装来自 Rust 服务器的所有 API 响应。
type ApiResponse[T any] struct {
	Success bool    `json:"success"`
	Data    *T      `json:"data"`  // 使用指针以处理 null 数据
	Error   *string `json:"error"` // 使用指针以处理 null 错误
}

// --- 请求结构体 ---

// commonRequestFields 包含了用于会话和代理的通用可选字段。
type commonRequestFields struct {
	SessionID *string `json:"session_id,omitempty"`
	Proxy     *string `json:"proxy,omitempty"`
}

type SimpleMatchRequest struct {
	GT        string `json:"gt"`
	Challenge string `json:"challenge"`
	commonRequestFields
}

type RegisterTestRequest struct {
	URL string `json:"url"`
	commonRequestFields
}

type GetCSRequest struct {
	GT        string  `json:"gt"`
	Challenge string  `json:"challenge"`
	W         *string `json:"w,omitempty"`
	commonRequestFields
}

// GetTypeRequest 和 VerifyRequest 的结构与 GetCSRequest 相同。
type GetTypeRequest = GetCSRequest
type VerifyRequest = GetCSRequest

type GenerateWRequest struct {
	Key       string `json:"key"`
	GT        string `json:"gt"`
	Challenge string `json:"challenge"`
	C         []byte `json:"c"`
	S         string `json:"s"`
	commonRequestFields
}

type TestRequest struct {
	URL string `json:"url"`
	commonRequestFields
}

// --- 响应数据结构体 ---

// TupleResponse2 对应于同名的 Rust 结构体。
type TupleResponse2 struct {
	First  string `json:"first"`
	Second string `json:"second"`
}

// CSResponse 对应于同名的 Rust 结构体。
type CSResponse struct {
	C []byte `json:"c"`
	S string `json:"s"`
}
