// geetest/click.go
package geetest

import (
	"context"
	"fmt"
)

// ClickService 处理 /click 端点的相关操作。
type ClickService struct {
	client *Client
}

// fillCommonFields 使用客户端的会ushua和代理信息来填充请求。
func (s *ClickService) fillCommonFields(req *commonRequestFields) {
	if req.SessionID == nil {
		req.SessionID = s.client.SessionID
	}
	if req.Proxy == nil {
		req.Proxy = s.client.Proxy
	}
}

// SimpleMatch 方法的修正
func (s *ClickService) SimpleMatch(ctx context.Context, gt, challenge string) (string, error) {
	req := SimpleMatchRequest{GT: gt, Challenge: challenge}
	s.fillCommonFields(&req.commonRequestFields)

	resultPtr, err := post[string](s.client, ctx, "/click/simple_match", req)
	if err != nil {
		return "", err
	}
	if resultPtr == nil {
		return "", fmt.Errorf("API 成功但返回的数据为空")
	}
	return *resultPtr, nil
}

// SimpleMatchRetry 方法的修正
func (s *ClickService) SimpleMatchRetry(ctx context.Context, gt, challenge string) (string, error) {
	req := SimpleMatchRequest{GT: gt, Challenge: challenge}
	s.fillCommonFields(&req.commonRequestFields)

	resultPtr, err := post[string](s.client, ctx, "/click/simple_match_retry", req)
	if err != nil {
		return "", err
	}
	if resultPtr == nil {
		return "", fmt.Errorf("API 成功但返回的数据为空")
	}
	return *resultPtr, nil
}

// RegisterTest 方法是正确的，因为返回值类型匹配
func (s *ClickService) RegisterTest(ctx context.Context, url string) (*TupleResponse2, error) {
	req := RegisterTestRequest{URL: url}
	s.fillCommonFields(&req.commonRequestFields)
	return post[TupleResponse2](s.client, ctx, "/click/register_test", req)
}

// GetCS 方法是正确的，因为返回值类型匹配
func (s *ClickService) GetCS(ctx context.Context, gt, challenge string, w *string) (*CSResponse, error) {
	req := GetCSRequest{GT: gt, Challenge: challenge, W: w}
	s.fillCommonFields(&req.commonRequestFields)
	return post[CSResponse](s.client, ctx, "/click/get_c_s", req)
}

// GetType 方法的修正
func (s *ClickService) GetType(ctx context.Context, gt, challenge string, w *string) (string, error) {
	req := GetTypeRequest{GT: gt, Challenge: challenge, W: w}
	s.fillCommonFields(&req.commonRequestFields)

	resultPtr, err := post[string](s.client, ctx, "/click/get_type", req)
	if err != nil {
		return "", err
	}
	if resultPtr == nil {
		return "", fmt.Errorf("API 成功但返回的数据为空")
	}
	return *resultPtr, nil
}

// Verify 方法是正确的，因为返回值类型匹配
func (s *ClickService) Verify(ctx context.Context, gt, challenge string, w *string) (*TupleResponse2, error) {
	req := VerifyRequest{GT: gt, Challenge: challenge, W: w}
	s.fillCommonFields(&req.commonRequestFields)
	return post[TupleResponse2](s.client, ctx, "/click/verify", req)
}

// GenerateW 方法的修正
func (s *ClickService) GenerateW(ctx context.Context, key, gt, challenge string, c []byte, s_val string) (string, error) {
	req := GenerateWRequest{Key: key, GT: gt, Challenge: challenge, C: c, S: s_val}
	s.fillCommonFields(&req.commonRequestFields)

	resultPtr, err := post[string](s.client, ctx, "/click/generate_w", req)
	if err != nil {
		return "", err
	}
	if resultPtr == nil {
		return "", fmt.Errorf("API 成功但返回的数据为空")
	}
	return *resultPtr, nil
}

// Test 方法的修正
func (s *ClickService) Test(ctx context.Context, url string) (string, error) {
	req := TestRequest{URL: url}
	s.fillCommonFields(&req.commonRequestFields)

	resultPtr, err := post[string](s.client, ctx, "/click/test", req)
	if err != nil {
		return "", err
	}
	if resultPtr == nil {
		return "", fmt.Errorf("API 成功但返回的数据为空")
	}
	return *resultPtr, nil
}
