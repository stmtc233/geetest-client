// geetest/slide.go
package geetest

import (
	"context"
	"fmt"
)

// SlideService 处理 /slide 端点的相关操作。
type SlideService struct {
	client *Client
}

// fillCommonFields 使用客户端的会话和代理信息来填充请求。
func (s *SlideService) fillCommonFields(req *commonRequestFields) {
	if req.SessionID == nil {
		req.SessionID = s.client.SessionID
	}
	if req.Proxy == nil {
		req.Proxy = s.client.Proxy
	}
	if req.User_Agent == nil {
		req.User_Agent = s.client.User_Agent
	}
}

// RegisterTest 方法是正确的
func (s *SlideService) RegisterTest(ctx context.Context, url string) (*TupleResponse2, error) {
	req := RegisterTestRequest{URL: url}
	s.fillCommonFields(&req.commonRequestFields)
	return post[TupleResponse2](s.client, ctx, "/slide/register_test", req)
}

// GetCS 方法是正确的
func (s *SlideService) GetCS(ctx context.Context, gt, challenge string, w *string) (*CSResponse, error) {
	req := GetCSRequest{GT: gt, Challenge: challenge, W: w}
	s.fillCommonFields(&req.commonRequestFields)
	return post[CSResponse](s.client, ctx, "/slide/get_c_s", req)
}

// GetType 方法的修正
func (s *SlideService) GetType(ctx context.Context, gt, challenge string, w *string) (string, error) {
	req := GetTypeRequest{GT: gt, Challenge: challenge, W: w}
	s.fillCommonFields(&req.commonRequestFields)

	resultPtr, err := post[string](s.client, ctx, "/slide/get_type", req)
	if err != nil {
		return "", err
	}
	if resultPtr == nil {
		return "", fmt.Errorf("API 成功但返回的数据为空")
	}
	return *resultPtr, nil
}

// Verify 方法是正确的
func (s *SlideService) Verify(ctx context.Context, gt, challenge string, w *string) (*TupleResponse2, error) {
	req := VerifyRequest{GT: gt, Challenge: challenge, W: w}
	s.fillCommonFields(&req.commonRequestFields)
	return post[TupleResponse2](s.client, ctx, "/slide/verify", req)
}

// GenerateW 方法的修正
func (s *SlideService) GenerateW(ctx context.Context, key, gt, challenge string, c []byte, s_val string) (string, error) {
	req := GenerateWRequest{Key: key, GT: gt, Challenge: challenge, C: c, S: s_val}
	s.fillCommonFields(&req.commonRequestFields)

	resultPtr, err := post[string](s.client, ctx, "/slide/generate_w", req)
	if err != nil {
		return "", err
	}
	if resultPtr == nil {
		return "", fmt.Errorf("API 成功但返回的数据为空")
	}
	return *resultPtr, nil
}

// Test 方法的修正
func (s *SlideService) Test(ctx context.Context, url string) (string, error) {
	req := TestRequest{URL: url}
	s.fillCommonFields(&req.commonRequestFields)

	resultPtr, err := post[string](s.client, ctx, "/slide/test", req)
	if err != nil {
		return "", err
	}
	if resultPtr == nil {
		return "", fmt.Errorf("API 成功但返回的数据为空")
	}
	return *resultPtr, nil
}
