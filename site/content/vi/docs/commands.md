---
title: Tham chiếu lệnh
description: "Mọi lệnh và tùy chọn của ccam, từ ccam add và ccam <name> đến ccam usage, ccam next và ccam doctor."
updated: 2026-10-19
order: 3
translation: machine
sourceHash: 5a5a33b1a2f4
---

Đặt `--dry-run` trước bất kỳ lệnh nào để xem lệnh đó sẽ làm gì mà không thực sự thực hiện, ví dụ `ccam --dry-run work`.

<!-- commands -->

## Khởi động Claude Code

`ccam <name|number> [arguments]` khởi động Claude Code bằng tài khoản đó. ccam không đọc các đối số sau tên tài khoản; nó chuyển chúng cho `claude` theo đúng thứ tự, và đầu vào chuẩn, mã thoát cũng như tín hiệu hoạt động giống hệt như khi chạy chính `claude`. Dấu `--` sau tên tài khoản là không bắt buộc.

`ccam` có đối số nhưng không có tên tài khoản, chẳng hạn `ccam -p "summarise"`, sẽ dùng tài khoản mặc định. Đặt tài khoản mặc định bằng `ccam use <name>`.

Một lệnh tắt như `ccwork` khởi động tài khoản của nó theo cùng cách: `ccwork --resume` làm đúng những gì `ccam work --resume` làm. Xem [lệnh tắt](/docs/commands#shortcuts).

## Thêm và xóa tài khoản

```sh
ccam add                      # wizard
ccam add work --share history=on
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam rename work job
ccam remove side              # folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
```

`ccam import` không kèm thư mục sẽ tìm `~/.claude` và `~/.claude-*` rồi hỏi bạn về từng thư mục. Các thư mục được nhập sẽ được đăng ký ngay tại chỗ; ccam không bao giờ di chuyển hay sao chép thông tin đăng nhập.

## Lệnh tắt

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Lệnh tắt là một lệnh do bạn tự đặt tên, dùng để khởi động Claude Code bằng một tài khoản. Giống như `ccam work`, nó chuyển nguyên vẹn mọi đối số cho Claude Code, nên `ccwork --resume` và `ccwork -p "hi"` hoạt động đúng như bạn mong đợi. Nó là một liên kết đến ccam trong `~/.local/bin`.

Lệnh tắt là tùy chọn. Trình hướng dẫn `ccam add` hỏi bạn có muốn tạo không, `ccam add <name> --shortcut <command>` tạo luôn mà không hỏi, và `ccam list` hiển thị lệnh tắt của từng tài khoản, hoặc `-` khi không có.

Tên dùng các ký tự a-z, 0-9, `.`, `-` và `_`. Tên không được trùng với một lệnh của ccam, lệnh tắt của tài khoản khác, hoặc tên của một chương trình khác trong `PATH`.

`ccam rename` giữ nguyên tên lệnh tắt, còn `ccam remove` xóa lệnh tắt cùng với tài khoản. `ccam shims install` tạo lại liên kết của mọi tài khoản, chẳng hạn sau khi bạn di chuyển ccam.

## Dùng chung

```sh
ccam share work                       # show the current state
ccam share work --history on --mcp team
ccam share work --settings off --fresh
```

Mỗi mục nhận `on` (nhóm `default`), `off`, hoặc tên một nhóm. Xem [dùng chung](/docs/sharing).

## Sao chép thiết lập hiện có

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` sao chép `settings.json`, `CLAUDE.md`, tiện ích mở rộng và máy chủ MCP người dùng từ một thư mục Claude Code vào một nhóm, để mọi tài khoản dùng chung các mục đó đều có chúng. Thư mục chỉ được đọc. Các tùy chọn `--settings`, `--memory`, `--extensions`, `--mcp` và `--history` chỉ sao chép những mục bạn nêu; nếu không có chúng, mọi thứ trừ lịch sử hội thoại đều được sao chép. `-g` hoặc `--group` chọn nhóm, còn `-y` hoặc `--yes` sao chép mà không hỏi và giữ bản của nhóm khi một tệp khác nhau. Xem [sao chép thiết lập hiện có](/docs/sharing#copy-an-existing-setup).

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
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam next --continue
ccam next --auto
```

## Bảo trì

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
```

`ccam update` kiểm tra mã băm SHA-256 của bản phát hành trước khi thay thế ccam, còn `--version v0.2.0` cài một bản phát hành cụ thể. ccam được cài từ gói `.deb` thì cập nhật bằng apt. Xem [cập nhật](/docs/install#update).

## Tệp

| Đường dẫn | Nội dung |
| --- | --- |
| `~/.config/ccam/config.toml` | Các tài khoản và những gì chúng dùng chung |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` của từng tài khoản |
| `~/.local/share/ccam/groups/<group>/` | Cài đặt, máy chủ MCP, tiện ích mở rộng và lịch sử dùng chung |
| `~/.local/state/ccam/` | Bộ nhớ đệm mức sử dụng và thời điểm dùng gần nhất |
| `~/.local/bin/<shortcut>` | Các lệnh tắt như `ccwork`, là liên kết đến ccam |

Các biến `XDG_CONFIG_HOME`, `XDG_DATA_HOME` và `XDG_STATE_HOME` dùng để đổi vị trí của các thư mục này.
