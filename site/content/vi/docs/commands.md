---
title: Tham chiếu lệnh
description: "Mọi lệnh và tùy chọn của ccam, từ ccam add và ccs <name> đến ccam usage, ccam next và ccam doctor."
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 18ecad9a1104
---

`ccam` và `ccs` là cùng một chương trình; bạn thích dùng tên nào cũng được. Đặt `--dry-run` trước bất kỳ lệnh nào để xem lệnh đó sẽ làm gì mà không thực sự thực hiện, ví dụ `ccs --dry-run work`.

<!-- commands -->

## Khởi động Claude Code

`ccs <name|number> [arguments]` khởi động Claude Code bằng tài khoản đó. ccam không đọc các đối số sau tên tài khoản; nó chuyển chúng cho `claude` theo đúng thứ tự, và đầu vào chuẩn, mã thoát cũng như tín hiệu hoạt động giống hệt như khi chạy chính `claude`. Dấu `--` sau tên tài khoản là không bắt buộc.

`ccs` có đối số nhưng không có tên tài khoản, chẳng hạn `ccs -p "summarise"`, sẽ dùng tài khoản mặc định. Đặt tài khoản mặc định bằng `ccam use <name>`.

Các lệnh tắt `claude1`, `claude2`, v.v. là liên kết đến ccam trong `~/.local/bin`. `ccam shims install --names` còn tạo thêm `claude-<name>` cho mọi tài khoản.

## Thêm và xóa tài khoản

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` không kèm thư mục sẽ tìm `~/.claude` và `~/.claude-*` rồi hỏi bạn về từng thư mục. Các thư mục được nhập sẽ được đăng ký ngay tại chỗ; ccam không bao giờ di chuyển hay sao chép thông tin đăng nhập.

## Dùng chung

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Mỗi mục nhận `on` (nhóm `default`), `off`, hoặc tên một nhóm. Xem [dùng chung](/docs/sharing).

## Máy chủ MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

## Mức sử dụng và chuyển đổi tài khoản

```sh
ccs usage            # also: ccs u
ccs usage --watch    # refresh every 60 seconds
ccs usage --json
ccs next --continue
ccs next --auto
```

## Bảo trì

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install
```

## Tệp

| Đường dẫn | Nội dung |
| --- | --- |
| `~/.config/ccam/config.toml` | Các tài khoản và những gì chúng dùng chung |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` của từng tài khoản |
| `~/.local/share/ccam/groups/<group>/` | Cài đặt, máy chủ MCP, tiện ích mở rộng và lịch sử dùng chung |
| `~/.local/state/ccam/` | Bộ nhớ đệm mức sử dụng và thời điểm dùng gần nhất |
| `~/.local/bin/claude1`, `ccs` | Liên kết đến ccam |

Các biến `XDG_CONFIG_HOME`, `XDG_DATA_HOME` và `XDG_STATE_HOME` dùng để đổi vị trí của các thư mục này.
