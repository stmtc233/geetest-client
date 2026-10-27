# geetest-client

配合 [biliTicker_gt_http](https://github.com/stmtc233/biliTicker_gt_http) 使用。

```bash
go get github.com/stmtc233/geetest-client
```

服务端同时支持：

- `POST /v4/solve`
- `POST /api/geetest/solve`（兼容旧接入）

## GeeTest v4 图标点选

B 站绑定手机号当前会返回 `recaptcha_type=geetest4`。v4 不是滑块，而是题目给
3 个图标、背景中找出对应图标并点击，成功后返回四项凭据。

```go
client, err := geetest.NewClient(
    "http://127.0.0.1:3000",
    geetest.WithProxy("socks5://127.0.0.1:1080"),
)
if err != nil {
    log.Fatal(err)
}

solution, err := client.V4.Solve(
    context.Background(),
    captchaID,   // B 站 gee_gt
    challenge,   // B 站 gee_challenge
    geetest.WithV4MaxAttempts(6),
)
if err != nil {
    log.Fatal(err)
}

fmt.Println(
    solution.LotNumber,
    solution.CaptchaOutput,
    solution.PassToken,
    solution.GenTime,
)
```

直接把四项凭据用于 `sms/send`：

- `gee_lot_number`
- `gee_captcha_output`
- `gee_pass_token`
- `gee_gen_time`

代理参数会透传给 Rust 服务，使 `/load`、图片下载与 `/verify` 尽可能使用同一出口 IP。
`V4Solution` 还包含 `Attempts`、`ElapsedMS` 和 `CandidateCount`，便于监控识别质量。
