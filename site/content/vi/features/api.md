---
title: Tài khoản API key
description: ccam add --api thêm tài khoản API key cho API của Anthropic hoặc một gateway; khóa được đọc khi khởi động từ một biến hoặc một lệnh, không bao giờ được lưu.
updated: 2026-10-23
order: 13
translation: machine
sourceHash: "232a71525e42"
---

Tài khoản API key khởi động Claude Code bằng một API key thay cho đăng nhập gói đăng ký Claude: API của Anthropic, hoặc một endpoint dùng cùng định dạng Messages, chẳng hạn một gateway LiteLLM hay OpenRouter. ccam không lưu khóa. Nó đọc khóa mỗi lần tài khoản khởi động, từ một biến môi trường hoặc từ một lệnh như trình quản lý mật khẩu, và chỉ trao khóa cho tiến trình Claude Code đó.

## Cách dùng

```sh
ccam add api --api --key-from env:ANTHROPIC_KEY_WORK
ccam add gw --api --base-url https://llm.example.com --auth bearer \
  --model claude-sonnet-5-5 --key-from 'command:pass show llm/gateway-key'
ccam add mac --api --key-from 'command:security find-generic-password -s anthropic-api -w'
ccam gw -p "summarise the open issues"
ccam list
```

```console
$ ccam list
   #  Name      Shortcut  Email                   Groups   Last used
*  1  personal  -         me@example.com          default  2h ago
   2  gw        -         (API: llm.example.com)  default  -
```

Tài khoản API key tham gia các nhóm như mọi tài khoản khác (theo mặc định là `default`: cài đặt, `CLAUDE.md`, tiện ích mở rộng và máy chủ MCP), và hoạt động với lệnh tắt, `ccam folder` và `ccam ext`.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--api` | Thêm một tài khoản API key thay vì đăng nhập | tắt |
| `--key-from env:NAME` | Đọc khóa từ biến môi trường `NAME` | |
| `--key-from command:COMMAND` | Chạy `COMMAND` (bằng `sh -c`, hoặc trên Windows như được gõ trong `cmd /d /s /c "…"`, giữ nguyên dấu ngoặc kép) và dùng dòng đầu tiên nó in ra | |
| `--base-url <url>` | Endpoint (`ANTHROPIC_BASE_URL`) | API của Anthropic |
| `--auth api-key\|bearer` | Gửi khóa dưới dạng `ANTHROPIC_API_KEY` (`x-api-key`) hoặc `ANTHROPIC_AUTH_TOKEN` (`Authorization: Bearer`, thứ mà gateway thường cần) | `api-key` |
| `--model <model>` | Mô hình (`ANTHROPIC_MODEL`) | Mặc định của Claude Code |

Các thiết lập được lưu trong `config.toml`:

```toml
[[profile]]
name = 'gw'
slot = 2

[profile.api]
base_url = 'https://llm.example.com'
key_from = 'command:pass show llm/gateway-key'
auth = 'bearer'
model = 'claude-sonnet-5-5'
```

## Cách hoạt động

Mỗi lần khởi động tài khoản, sau các bước kiểm tra thông thường và trước khi Claude Code khởi động:

1. ccam đọc khóa: `env:` từ môi trường, `command:` bằng cách chạy lệnh với giới hạn 30 giây. Lệnh có thể hỏi mật khẩu trên terminal; đầu ra lỗi của nó cũng đi tới terminal. Một khóa rỗng sẽ dừng việc khởi động với thông báo cho biết khóa được chờ đợi ở đâu.
2. Nó xóa mọi biến kế thừa quyết định cách Claude Code đăng nhập hoặc nhà cung cấp nào được dùng: các khóa và token (`ANTHROPIC_API_KEY`, `ANTHROPIC_AUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN`, `CLAUDE_CODE_OAUTH_TOKEN_FILE_DESCRIPTOR`, `CLAUDE_CODE_API_KEY_FILE_DESCRIPTOR`), endpoint (`ANTHROPIC_BASE_URL`, `ANTHROPIC_UNIX_SOCKET`, `ANTHROPIC_MODEL`) và các công tắc nhà cung cấp (`CLAUDE_CODE_USE_BEDROCK`, `CLAUDE_CODE_USE_VERTEX`, `CLAUDE_CODE_USE_FOUNDRY`, `CLAUDE_CODE_USE_ANTHROPIC_AWS`, `CLAUDE_CODE_USE_ANTHROPIC_GOOGLE_CLOUD`, `CLAUDE_CODE_USE_MANTLE`, `CLAUDE_CODE_USE_GATEWAY`). Một khóa đặt trong shell cho việc khác không thể tới được tài khoản này, và một đăng nhập gói đăng ký không thể lấn át khóa của nó.
3. Nó đặt `ANTHROPIC_BASE_URL` và `ANTHROPIC_MODEL` khi được cung cấp, và đặt khóa thành `ANTHROPIC_API_KEY` hoặc `ANTHROPIC_AUTH_TOKEN`.

Các tài khoản đăng nhập bằng gói đăng ký vẫn giữ môi trường mà chúng kế thừa, như trước. Trong một phiên tương tác, Claude Code có thể hỏi một lần có dùng khóa nó tìm thấy trong môi trường hay không.

## An toàn và các trường hợp đặc biệt

- Khóa không bao giờ được ghi vào `config.toml`, bộ nhớ đệm hay nhật ký, và `ccam --dry-run gw` không đọc nó: chạy thử không chạy nguồn `command:`. Thay vào đó, lần chạy thử hiển thị endpoint, mô hình và nơi khóa được lấy:

  ```console
  $ ccam --dry-run gw -p hi
  CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/gw CCAM_PROFILE=gw ANTHROPIC_BASE_URL=https://llm.example.com ANTHROPIC_MODEL=claude-sonnet-5-5 ANTHROPIC_AUTH_TOKEN="$(pass show llm/gateway-key | head -n 1)" /home/me/.local/bin/claude -p hi
  ```

  Một biến được kế thừa mà lần khởi chạy gỡ bỏ sẽ hiển thị là `env -u <name>` ở đầu. Với nguồn `env:` có biến trống, lần chạy thử thất bại giống như lần khởi chạy.
- `ccam add --api` kiểm tra nguồn `env:` ngay lập tức và cảnh báo nếu biến rỗng; nguồn `command:` được chạy lần đầu khi khởi động.
- `ccam login` và `ccam logout` từ chối tài khoản API key, vì nó không có đăng nhập. `ccam doctor` hiển thị endpoint và nguồn khóa của nó, và cảnh báo khi một biến `env:` rỗng.
- `ccam usage` hiển thị `no limits (API)` cho nó. `ccam next --auto`, `--fallback` và `ccam pick` chỉ chọn trong các tài khoản gói đăng ký, vì API key không có giới hạn 5 giờ hay hàng tuần để so sánh.
- Một `settings.json` dùng chung có khối `env` hoặc `apiKeyHelper` sẽ áp dụng cho mọi tài khoản trong nhóm, kể cả tài khoản API key. Đừng đặt thiết lập nhà cung cấp trong cài đặt dùng chung, hoặc đưa tài khoản API key vào một nhóm không có `settings`.

## Liên quan

- [Quản lý tài khoản](/features/accounts)
- [Khởi động Claude Code](/features/launch)
- [Tham chiếu lệnh](/docs/commands)
