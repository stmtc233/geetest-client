package geetest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV4SolveSendsChallengeAndReturnsCredentials(t *testing.T) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/solve" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var request V4SolveRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if request.CaptchaID != "captcha-id" || request.Challenge != "challenge-value" {
			t.Fatalf("unexpected challenge: %+v", request)
		}
		if request.RiskType != "icon" || request.MaxAttempts != 5 {
			t.Fatalf("unexpected options: %+v", request)
		}
		if request.Proxy == nil || *request.Proxy != "socks5://127.0.0.1:1080" {
			t.Fatalf("client proxy was not forwarded: %+v", request.Proxy)
		}
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"data":{"lot_number":"lot","captcha_output":"output","pass_token":"token","gen_time":"123","attempts":2,"elapsed_ms":2100.5,"candidate_count":4},"error":null}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, WithProxy("socks5://127.0.0.1:1080"))
	if err != nil {
		t.Fatal(err)
	}
	solution, err := client.V4.Solve(
		context.Background(),
		"captcha-id",
		"challenge-value",
		WithV4MaxAttempts(5),
	)
	if err != nil {
		t.Fatal(err)
	}
	if solution.LotNumber != "lot" || solution.CaptchaOutput != "output" ||
		solution.PassToken != "token" || solution.GenTime != "123" {
		t.Fatalf("unexpected solution: %+v", solution)
	}
	if solution.Attempts != 2 || solution.CandidateCount != 4 {
		t.Fatalf("unexpected metrics: %+v", solution)
	}
}

func TestV4SolveRejectsEmptyCaptchaID(t *testing.T) {
	client, err := NewClient("http://127.0.0.1:1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.V4.Solve(context.Background(), "", "challenge"); err == nil {
		t.Fatal("expected an error for an empty captcha_id")
	}
}
