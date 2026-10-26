---
title: Phiên chạy nền của mọi tài khoản
description: "ccam agents liệt kê các agent chạy nền và phiên đang mở của Claude Code trên mọi tài khoản trong một bảng, còn ccam attach mở một phiên theo ID của nó."
updated: 2026-10-26
order: 16
translation: machine
sourceHash: 06281fb13119
---

Claude Code giữ các phiên chạy nền (`claude --bg`, chế độ xem agent) theo từng tài khoản: `claude agents` chỉ hiển thị các phiên của tài khoản mà nó đang chạy. `ccam agents` hỏi mọi tài khoản cùng lúc và liệt kê tất cả cùng nhau, còn `ccam attach <id|name>` mở một phiên chạy nền bằng đúng tài khoản đã khởi động nó.

## Cách dùng

```sh
ccam work --bg "fix the flaky test"   # start one in the background, as work
ccam agents                           # every account's sessions
ccam agents --all                     # also finished background sessions
ccam agents --json                    # for scripts
ccam attach 49ec0028                  # open it, as the account that started it
ccam attach 49ec                      # a unique prefix of the ID is enough
ccam attach probe                     # or part of the session's name
```

```console
$ ccam agents
Account   ID        Name   State           Folder      Started
work      -         -      open (busy)     ~/code/web  just now
personal  49ec0028  probe  blocked, idle   ~/code/app  12m ago
Open a background session with `ccam attach <id>`.
```

`State` là trạng thái của phiên chạy nền (`working`, `blocked`, `done`, `failed`, `stopped`) và việc tiến trình của nó đang làm (`busy`, `waiting`, `idle`); một phiên tương tác đang mở hiển thị `open`.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--all` | Liệt kê cả các phiên chạy nền đã xong, thất bại hoặc đã dừng | tắt |
| `--json` | In các phiên dưới dạng JSON, mỗi phiên kèm `account` mà nó thuộc về | tắt |

`ccam attach <id|name> [claude args]` chuyển mọi thứ sau ID cho `claude attach` mà không thay đổi, kể cả các cờ: `ccam attach 49ec --dry-run` giao `--dry-run` cho Claude Code. Để in ra lệnh thay vào đó, hãy đặt `--dry-run` trước ID (`ccam --dry-run attach 49ec` hoặc `ccam attach --dry-run 49ec`).

## Cách hoạt động

1. Với mọi tài khoản, ccam chạy `claude agents --json` (kèm `--all` khi được chỉ định) với `CLAUDE_CONFIG_DIR` của tài khoản, chạy song song, mỗi tài khoản với thời gian chờ tối đa 15 giây. Việc này đọc danh sách phiên riêng của tài khoản và không khởi động dịch vụ chạy nền của Claude Code.
2. Mỗi mục được gắn tên tài khoản, và danh sách được sắp xếp theo thời điểm bắt đầu, mới nhất trước. Tài khoản không phản hồi được nêu tên trên đầu ra lỗi chuẩn và bị bỏ ra, rồi ccam thoát với mã 1, để script không nhầm một lỗi thành một danh sách trống.
3. `ccam attach` tìm `jobs/<id>` trong thư mục của từng tài khoản, nơi Claude Code giữ một phiên chạy nền, và chạy `claude attach <id>` bằng tài khoản đó, theo cùng cách `ccam <name>` khởi động Claude Code. Một từ không phải là tiền tố ID sẽ được so khớp với tên phiên, không phân biệt hoa thường; một từ trông giống ID không bao giờ được coi là tên.

## An toàn và các trường hợp đặc biệt

- ID của phiên chạy nền thuộc về một tài khoản: `claude attach` chạy bằng tài khoản khác không tìm thấy chúng, và đó chính là điều `ccam attach` giải quyết.
- Nếu các phiên ở hai tài khoản có chung một tiền tố ID hoặc một từ trong tên, `ccam attach` liệt kê chúng và yêu cầu một ID dài hơn hoặc tên.
- ccam chuyển các máy chủ MCP dùng chung cho Claude Code bằng `--mcp-config=<file>` trong một đối số duy nhất, dạng mà `claude agents` và `claude attach` chấp nhận; các phiên chạy nền tiếp tục dùng tệp đó khi chúng khởi động lại.
- Dịch vụ chạy nền của một tài khoản API key lấy khóa từ shell đã khởi động nó. Hãy khởi động phiên chạy nền đầu tiên của một tài khoản như vậy qua ccam (`ccam gw --bg …`), để dịch vụ có khóa.

## Liên quan

- [Khởi động Claude Code](/features/launch)
- [Chọn tài khoản cho script](/features/pick)
- [Bảng điều khiển](/features/dashboard)
