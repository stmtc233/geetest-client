// geetest/v4.go
package geetest

import (
	"context"
	"fmt"
)

// V4Service handles GeeTest v4 icon challenges.
type V4Service struct {
	client *Client
}

func (s *V4Service) fillCommonFields(req *commonRequestFields) {
	if req.SessionID == nil {
		req.SessionID = s.client.SessionID
	}
	if req.Proxy == nil {
		req.Proxy = s.client.Proxy
	}
	if req.ImageUseProxy == nil {
		req.ImageUseProxy = s.client.ImageUseProxy
	}
	if req.User_Agent == nil {
		req.User_Agent = s.client.User_Agent
	}
	if req.Referer == nil {
		req.Referer = s.client.Referer
	}
}

// Solve exchanges a Bilibili-provided captcha_id/challenge pair for the four
// GeeTest v4 credentials required by the SMS endpoint.
func (s *V4Service) Solve(ctx context.Context, captchaID, challenge string, options ...V4SolveOption) (*V4Solution, error) {
	if captchaID == "" {
		return nil, fmt.Errorf("captcha_id 不能为空")
	}
	req := V4SolveRequest{
		CaptchaID: captchaID,
		Challenge: challenge,
		RiskType:  "icon",
	}
	for _, option := range options {
		option(&req)
	}
	s.fillCommonFields(&req.commonRequestFields)
	return post[V4Solution](s.client, ctx, "/v4/solve", req)
}

// V4SolveOption customizes one v4 solve request.
type V4SolveOption func(*V4SolveRequest)

// WithV4MaxAttempts overrides the server default. The server caps it at 10.
func WithV4MaxAttempts(attempts int) V4SolveOption {
	return func(req *V4SolveRequest) {
		req.MaxAttempts = attempts
	}
}

// WithV4RiskType overrides the challenge type. Bilibili currently uses "icon".
func WithV4RiskType(riskType string) V4SolveOption {
	return func(req *V4SolveRequest) {
		req.RiskType = riskType
	}
}
