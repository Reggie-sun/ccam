---
title: Dùng chung giữa các tài khoản
description: "Cách ccam share liên kết cài đặt, CLAUDE.md, tiện ích mở rộng và lịch sử vào một thư mục nhóm, với kế hoạch, bản sao lưu và lựa chọn khi xung đột, và cách off sao chép nội dung trở lại."
updated: 2026-10-19
order: 4
translation: machine
sourceHash: b30b7cdb2b15
---

Mỗi tài khoản có thư mục riêng, nhưng phần lớn nội dung trong đó không cần tách riêng. `ccam share` bật hoặc tắt dùng chung cho một tài khoản và một mục mỗi lần: tệp hoặc thư mục của tài khoản được chuyển vào một thư mục nhóm và được thay bằng một liên kết tượng trưng, nên mọi tài khoản trong nhóm đọc và ghi cùng một bản. Hãy dùng nó sau `ccam add` để thay đổi những gì một tài khoản dùng chung, hoặc để đưa một vài tài khoản vào một nhóm riêng.

## Cách dùng

```sh
ccam share work                       # show what work shares, and with which group
ccam share work --history on          # share conversation history with the default group
ccam share side --history team        # share it with the group team instead
ccam share work --settings off        # give work its own copy of settings.json again
ccam share work --memory off --fresh  # stop sharing CLAUDE.md and start with no file
ccam share work --settings on --yes --prefer merge   # no questions; merge JSON on conflicts
ccam --dry-run share work --history on               # print the plan, change nothing
```

Không có tùy chọn chọn mục, `ccam share` in ra trạng thái hiện tại:

```console
$ ccam share work
Sharing for work
  MCP servers                        shared (group default)  --mcp
  settings.json                      shared (group default)  --settings
  global CLAUDE.md                   shared (group default)  --memory
  agents, commands, skills, plugins  shared (group default)  --extensions
  conversation history               local                   --history
```

Khi có tùy chọn, ccam in ra một kế hoạch, hỏi về từng tệp khác với bản của nhóm, và hỏi `Proceed? [y/N]` trước khi thay đổi bất cứ thứ gì. Tài khoản mới dùng chung mọi thứ trừ lịch sử hội thoại; [Dùng chung máy chủ MCP, cài đặt và lịch sử](/docs/sharing) mô tả từng mục từ góc nhìn người dùng.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--mcp <on\|off\|group>` | Máy chủ MCP | không đổi |
| `--settings <on\|off\|group>` | `settings.json` | không đổi |
| `--memory <on\|off\|group>` | `CLAUDE.md` toàn cục | không đổi |
| `--extensions <on\|off\|group>` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | không đổi |
| `--history <on\|off\|group>` | `projects/` và `history.jsonl` | không đổi |
| `--fresh` | khi tắt dùng chung, bắt đầu từ trạng thái trống thay vì sao chép nội dung của nhóm | tắt |
| `--prefer <group\|mine\|merge>` | cách quyết định xung đột | hỏi; với `--yes`, mặc định của từng mục |
| `-y`, `--yes` | chạy mà không hỏi; xung đột tuân theo `--prefer` | tắt |
| `-h`, `--help` | trợ giúp cho `share` | |
| `--dry-run` (toàn cục) | được hỗ trợ: hỏi về xung đột, in ra kế hoạch, không thay đổi gì | tắt |

Giá trị `on` nghĩa là nhóm `default`, `off` nghĩa là riêng (không dùng chung), và mọi giá trị khác là tên nhóm. Tên nhóm trở thành tên thư mục trong `~/.local/share/ccam/groups/`, nên `ccam share side --history ../x` thất bại với `invalid group name: "../x"`. Tài khoản được chỉ định bằng tên hoặc số, và chỉ những mục bạn nêu mới thay đổi.

## Cách hoạt động

Các tệp liên quan:

- `~/.config/ccam/config.toml`: bảng `[profile.share]` của từng tài khoản, mỗi mục một giá trị.
- `~/.local/share/ccam/profiles/<name>/`: `CLAUDE_CONFIG_DIR` của tài khoản.
- `~/.local/share/ccam/groups/<group>/`: các bản dùng chung.

1. `internal/cli/share.go` đọc cấu hình, áp dụng từng giá trị `--<item>` lên một bản sao thiết lập của tài khoản bằng `config.Share.Set`, và từ chối nếu tài khoản đang chạy.
2. `share.Build` so sánh thiết lập hiện tại với thiết lập đích và liệt kê các bước. Nó chỉ đọc đĩa: `os.Lstat` và `os.Readlink` cho mọi đường dẫn trong `share.Layout`. Kế hoạch là một danh sách thao tác: chuyển vào một vị trí còn trống trong nhóm (adopt), chuyển các tệp nhóm còn thiếu (merge dir), ghi một xung đột đã giải quyết vào nhóm, sao lưu, gỡ liên kết, liên kết, sao chép trở lại, ghi lại đường dẫn plugin và gộp tệp đăng ký plugin.
3. Các xung đột được quyết định, qua câu hỏi hoặc qua `--prefer`.
4. `Plan.Describe` in ra kế hoạch, và ccam hỏi.
5. Thiết lập đích được ghi vào `config.toml` trước khi bất cứ thứ gì được di chuyển. Nếu `Plan.Apply` dừng giữa chừng, `ccam doctor --fix` và một lần chạy lại sẽ dựng kế hoạch hướng tới đích đã lưu và hoàn tất thay đổi.
6. `Plan.Apply` chạy các bước theo thứ tự. Mọi bước đều chuyển dữ liệu vào nhóm hoặc vào thư mục sao lưu; không có gì bị xóa.

`share.Layout` là một danh sách cho phép gồm các đường dẫn mà mỗi mục bao gồm. Bất cứ thứ gì không có trong đó đều không bao giờ đến được nhóm. `share.Never` nêu tên `.credentials.json` và `.claude.json`, và `Build` dừng với lỗi nội bộ nếu một trong hai xuất hiện trong cấu trúc. Mục `mcp` không có đường dẫn: danh sách máy chủ nằm bên trong `.claude.json`, nên [dùng chung MCP](/features/mcp) hoạt động qua `mcp.json` của nhóm thay vào đó, và `--mcp` chỉ thay đổi thiết lập.

Đây là điều xảy ra với một thư mục, `projects/`, khi `side` bật dùng chung lịch sử và nhóm đã có nội dung:

<!-- motion:share-on -->

Nếu nhóm chưa có gì ở đường dẫn đó, tài khoản là tài khoản đầu tiên dùng chung nó và thư mục của nó được chuyển nguyên vào đó, không sao lưu. Ngược lại, `mergeDir` chuyển từng tệp nhóm còn thiếu vào nhóm, còn các tệp có ở cả hai phía vẫn nằm trong thư mục của tài khoản, nên phiên bản của nhóm được giữ. Sau đó thư mục của tài khoản được chuyển vào `.ccam-backup/<time>/` bên trong thư mục tài khoản, và `os.Symlink` đặt một liên kết vào chỗ của nó. Các tệp hội thoại được đặt tên theo ID phiên, nên lịch sử từ nhiều tài khoản được gộp lại mà không xung đột.

Việc di chuyển dùng `rename(2)` và quay về sao chép rồi xóa khi đi qua các hệ thống tệp khác nhau. Một liên kết tương đối bên trong cây được di chuyển mà trỏ ra ngoài cây sẽ được ghi lại thành liên kết tuyệt đối (`fsutil.PortableLink`), để nó vẫn trỏ đúng từ vị trí mới. Các tệp đăng ký trong `plugins/` (`installed_plugins.json`, `known_marketplaces.json`) lưu đường dẫn tuyệt đối đi qua thư mục tài khoản; ccam trỏ chúng tới thư mục của nhóm và, khi nhóm đã có tệp đăng ký, gộp các mục của tài khoản vào đó.

Tắt một mục sẽ đảo ngược liên kết:

<!-- motion:share-off -->

ccam chỉ gỡ liên kết nếu nó trỏ vào thư mục groups, rồi sao chép nội dung hiện tại của nhóm vào tài khoản bằng `CopyTree`. Với `--fresh`, bước sao chép bị bỏ qua và tài khoản bắt đầu với không có gì ở đường dẫn đó. Nhóm vẫn giữ nội dung của nó cho các tài khoản khác. Chuyển một mục từ nhóm này sang nhóm khác sẽ gỡ liên kết rồi liên kết nó tới nhóm mới, giữ nguyên nhóm cũ; kế hoạch ghi điều này trong một dòng `Note:`.

Một tệp đơn lẻ có ở cả hai phía với nội dung khác nhau là một xung đột:

<!-- motion:share-conflict -->

`share.NewConflict` liệt kê các lựa chọn phù hợp với tệp. `history.jsonl` có `merge by line`, `keep group` và `use mine`, với gộp đứng đầu. Các tệp `.json` khác có `keep group`, `use mine` và `merge JSON`. Mọi tệp khác có `keep group` và `use mine`. Lựa chọn đầu tiên là mặc định:

```console
history.jsonl differs from the group's copy.
  mine:  ~/.local/share/ccam/profiles/side/history.jsonl
  group: ~/.local/share/ccam/groups/default/history.jsonl
  1) merge by line
  2) keep group
  3) use mine
Choice [1]:
```

`merge JSON` (`MergeJSONBytes`) gộp các đối tượng một cách đệ quy, nối các mảng mà không trùng lặp, và giữ giá trị của nhóm khi hai phía đặt cùng một khóa khác nhau, nên các tài khoản khác không thấy thay đổi gì về hành vi. `merge by line` giữ các dòng của nhóm, rồi thêm các dòng của bạn mà nhóm còn thiếu. Mọi lựa chọn trừ `keep group` trước tiên lưu bản của nhóm vào `.ccam-backup/<time>/group/` trong thư mục tài khoản, rồi ghi kết quả bằng `fsutil.AtomicWrite`. Tệp của riêng bạn luôn được đưa vào bản sao lưu.

`--prefer group`, `mine` hoặc `merge` quyết định mọi xung đột mà không hỏi; `merge` nghĩa là `merge JSON` hoặc `merge by line`, tùy tệp có lựa chọn nào, và tệp không có lựa chọn nào trong hai thì giữ mặc định. Chỉ `--yes` thì lấy mặc định của từng tệp.

Ở mỗi lần khởi động, `share.Check` so sánh từng đường dẫn trong cấu trúc với thiết lập, chỉ dùng `lstat` và `readlink`, nhờ đó đủ nhẹ cho đường khởi động. Liên kết còn thiếu được tạo ngay lập tức, vì không có gì để mất. Một tệp thật ở nơi lẽ ra là liên kết, một liên kết trỏ sai chỗ, hoặc một liên kết lỗi thời sẽ dừng việc khởi động:

```console
ccam: stopped because of share link problems:
  ~/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

`ccam doctor --fix` chạy `share.Build` với thiết lập đã lưu ở cả hai phía, suy ra mọi bước từ những gì có trên đĩa, giữ nội dung của nhóm khi xung đột và sao lưu mọi thứ nó thay thế.

## An toàn và các trường hợp đặc biệt

- **Tài khoản đang chạy.** ccam từ chối thay đổi việc dùng chung khi Claude Code đang chạy bằng tài khoản đó (`account work is running; quit it first, then change its sharing`), vì phiên vẫn tiếp tục ghi vào các tệp đang được di chuyển và các dòng nó ghi thêm trong lúc đó sẽ bị mất. `--dry-run` vẫn in ra kế hoạch.
- **Bản sao lưu.** Mỗi kế hoạch có thư mục `.ccam-backup/<time>/` riêng trong thư mục tài khoản, với hậu tố đánh số khi hai kế hoạch bắt đầu trong cùng một giây. ccam không bao giờ xóa chúng.
- **Thất bại giữa chừng.** Nếu một bước thất bại, ccam in ra lỗi và nơi chứa các bản gốc, và nhắc bạn hoàn tất bằng `ccam doctor --fix`. Đích đã được lưu, nên việc sửa tiếp tục hướng tới đích đó.
- **Liên kết của riêng bạn.** Nếu một đường dẫn đã là liên kết do bạn tạo, chẳng hạn tới một kho dotfiles, và nhóm còn trống, ccam sao chép đích của liên kết vào nhóm và chuyển chính liên kết vào bản sao lưu. Đích không bị động đến.
- **Tệp đăng nhập.** `.credentials.json` và `.claude.json` không bao giờ được liên kết, di chuyển hay sao chép vào nhóm.
- **Hai tài khoản trong một thư mục.** Các tài khoản dùng chung lịch sử ghi vào cùng các tệp. Khi một tài khoản thứ hai của cùng nhóm lịch sử khởi động trong một thư mục nơi đã có một tài khoản đang chạy, ccam cảnh báo và hỏi trước khi tiếp tục.

## Liên quan

- [Dùng chung máy chủ MCP, cài đặt và lịch sử](/docs/sharing)
- [Máy chủ MCP dùng chung](/features/mcp)
- [Sao chép thiết lập hiện có](/features/migrate)
- [Tất cả các lệnh](/docs/commands)
