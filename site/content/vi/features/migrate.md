---
title: Sao chép thiết lập hiện có
description: "Sao chép cài đặt, CLAUDE.md, tiện ích mở rộng và máy chủ MCP của một ~/.claude hiện có vào một nhóm của ccam bằng ccam migrate."
updated: 2026-10-26
order: 6
translation: machine
sourceHash: f835fbe84ade
---

Nếu bạn đã dùng Claude Code trước khi có ccam, thiết lập của bạn nằm trong `~/.claude` và `~/.claude.json`. `ccam migrate` sao chép nó vào một nhóm dùng chung, để mọi tài khoản dùng chung các mục đó qua nhóm đều bắt đầu với thiết lập bạn đã có. Thư mục nguồn chỉ được đọc, và `claude` thông thường vẫn hoạt động như trước. Hãy chạy lệnh này một lần sau `ccam add` đầu tiên; khi `~/.claude.json` có máy chủ MCP, lần `ccam add` đầu tiên sẽ nêu tên chúng trong lời gợi ý của nó.

## Cách dùng

```sh
ccam migrate                         # ~/.claude into group default, checklist on a terminal
ccam --dry-run migrate               # show the plan, copy nothing
ccam migrate --settings --mcp        # only these items, no checklist
ccam migrate --history --yes         # conversation history only, no questions
ccam migrate ~/.claude-work --group work
```

Trên terminal, khi không có cờ chọn mục và không có `--yes`, ccam trước tiên hiển thị một danh sách kiểm (checklist):

```console
 Copy into group default

 From /home/me/.claude
 Choose what to copy (Space to toggle)
 ❯ [✔] MCP servers
   [✔] settings.json
   [✔] global CLAUDE.md
   [✔] skills, agents, commands, plugins
   [ ] conversation history

 space to toggle · enter to show the plan · esc to cancel
```

Space (hoặc `x`) bật/tắt mục dưới con trỏ, các phím mũi tên hoặc `j`/`k` để di chuyển, Enter hiển thị kế hoạch, còn Esc, `q` hoặc Ctrl+C in ra `Cancelled.`. Những mục mà nguồn không có gì được ghi là `not found`. Nếu bạn bỏ chọn mọi thứ, ccam in ra `Nothing selected.` và dừng lại.

Sau đó ccam in ra kế hoạch và hỏi trước khi sao chép:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        not copied (a sign-in works for one account only); sign in with /mcp: linear in personal, work
  settings    settings.json   new: 61 B
              not found, skipped: keybindings.json
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/
MCP sign-ins are not copied (a sign-in works for one account only). Sign in with /mcp to linear in personal, work.
Group default: personal, work
```

Một nhóm chỉ tác động đến các tài khoản của chính nó, và chỉ với các mục nó dùng chung, nên phần tóm tắt liệt kê các tài khoản của nhóm. Khi nhóm chưa dùng chung một trong các mục đã sao chép, hoặc có tài khoản chưa thuộc nhóm, ccam in ra các lệnh `ccam group set <group> --add <items>` và `ccam group join <group> <accounts>` để thay đổi điều đó.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `[folder]` | Thư mục Claude Code để sao chép từ đó | `~/.claude` |
| `--settings` | Sao chép `settings.json` | tắt |
| `--memory` | Sao chép `CLAUDE.md` toàn cục | tắt |
| `--extensions` | Sao chép `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | tắt |
| `--mcp` | Sao chép máy chủ MCP người dùng (`mcpServers` trong `.claude.json`) | tắt |
| `--history` | Sao chép lịch sử hội thoại (`projects/`, `history.jsonl`) | tắt |
| `-g`, `--group <group>` | Nhóm dùng chung để sao chép vào | `default` |
| `-y`, `--yes` | Sao chép mà không hỏi; khi xung đột, bản của nhóm được giữ | tắt |
| `--dry-run` (toàn cục) | In ra kế hoạch và dừng trước mọi câu hỏi hay việc sao chép | tắt |

Trên terminal, nếu bạn không đưa cờ chọn mục nào, bạn chọn các mục trong danh sách kiểm. Ngoài terminal, hoặc với `--yes`, ccam sao chép mọi thứ trừ lịch sử hội thoại. Lịch sử bị loại theo mặc định vì `projects/` có thể rất lớn.

`--group on` nghĩa là `default`. `--group off` và `--group local` bị từ chối với `--group takes a group name, not "off"`. `ccam migrate` hỗ trợ `--dry-run`, và danh sách kiểm vẫn xuất hiện khi chạy thử để bạn có thể xem trước một lựa chọn.

## Cách hoạt động

ccam dựng một kế hoạch hoàn chỉnh trước khi sao chép bất cứ thứ gì. Nó chọn các mục (`internal/cli/migrate.go`), rồi đọc nguồn và nhóm và dựng kế hoạch (`migrate.Build`). Không có gì được ghi cho đến `Plan.Apply`, sau khi bạn xác nhận.

1. **Tìm nguồn.** ccam từ chối một đường dẫn không tồn tại, không phải thư mục, hoặc không chứa mục Claude Code nào đã biết (`checkSource`). Chỉ riêng với máy chủ MCP, một `~/.claude` trống hoặc không có vẫn được, miễn là `~/.claude.json` có máy chủ. ccam tìm `.claude.json` thuộc về thư mục (`migrate.ClaudeJSON`): bên trong nó, hoặc với `~/.claude` không có tệp đó bên trong, là `~/.claude.json` nằm cạnh nó. Đó là nơi Claude Code giữ tệp này khi `CLAUDE_CONFIG_DIR` không được đặt.
2. **Ánh xạ mục sang đường dẫn.** Mọi đường dẫn đều đến từ `share.Layout`, cũng là danh sách cho phép mà các nhóm dùng. Các máy chủ MCP được đọc ra khỏi `.claude.json` bằng `mcp.LocalServers`, hàm này dùng bộ quét `profile.TopLevelValue` thay vì phân tích cả tệp.
3. **Sao chép, không bao giờ di chuyển.** Các tệp đi vào `~/.local/share/ccam/groups/<group>/` bằng `fsutil.AtomicWrite` hoặc `fsutil.CopyFile`, hàm này ghi một tệp tạm bên cạnh đích rồi đổi tên nó, nên một lần sao chép bị gián đoạn không bao giờ để lại tệp dở dang. Máy chủ MCP đi vào `mcp.json` của nhóm dưới khóa của nó: một máy chủ mà một tiến trình ccam khác thêm vào trong lúc đó vẫn được giữ, và các bản sao được tạo cho [Remote Control](/guides/remote-control) được cập nhật sau đó.

Mỗi mục ánh xạ một đường dẫn nguồn tới một đường dẫn trong nhóm:

<!-- motion:migrate-items -->

Hai hàng lịch sử có nét đứt vì chúng chỉ được sao chép với `--history`. `.credentials.json` và `.claude.json` được liệt kê trong `share.Never` và không bao giờ đến được nhóm: ccam dừng với một lỗi nội bộ nếu một trong hai xuất hiện trong cấu trúc (`Build`), và việc duyệt thư mục bỏ qua chúng theo tên.

ccam gán cho mỗi đường dẫn một trạng thái, đó là điều mỗi dòng trong kế hoạch mô tả (`buildEntry`):

<!-- motion:migrate-states -->

- **Missing**: không có trong nguồn, hoặc là một tệp ở nơi lẽ ra là thư mục. Được liệt kê một lần cho mỗi mục dưới dạng `not found, skipped`.
- **Linked**: đường dẫn nguồn đã trỏ tới bản của nhóm, ví dụ khi nguồn là một thư mục hồ sơ của ccam.
- **Same**: nhóm đã có nội dung giống hệt.
- **New**: nhóm chưa có, nên được sao chép toàn bộ.
- **Merge**: một thư mục có ở cả hai phía. ccam chỉ thêm các tệp nhóm còn thiếu (`copyMissing`); các tệp khác nhau giữ nguyên như ở nhóm và kế hoạch ghi `group copies kept`. Với `plugins/`, các tệp đăng ký `installed_plugins.json` và `known_marketplaces.json` được gộp, và các đường dẫn plugin tuyệt đối được ghi lại để trỏ tới thư mục của nhóm.
- **Conflict**: một tệp có ở cả hai phía với nội dung khác nhau. Bạn chọn cho từng tệp: `keep group`, `use mine`, cùng với `merge JSON` cho các tệp `.json` (gộp sâu, trong đó nhóm thắng khi giá trị xung đột) hoặc `merge by line` cho `history.jsonl`. Một máy chủ MCP có định nghĩa khác cũng là xung đột, với `keep group` hoặc `use mine`. Lựa chọn đầu tiên được liệt kê là mặc định khi bạn nhấn Enter; `--yes` luôn giữ bản của nhóm.

Bản thân danh sách kiểm là `tui.MigratePicker`, được dựng trên thành phần `checklist` mà `ccam add` dùng; ccam đánh dấu các mục còn thiếu bằng các lời gọi lstat và danh sách máy chủ MCP (`migrate.Found`).

### Thông tin đăng nhập MCP

Các máy chủ MCP đăng nhập bằng OAuth giữ token trong tệp đăng nhập, `.credentials.json` dưới `mcpOAuth`, không phải trong danh sách máy chủ. ccam không sao chép chúng. Một lần đăng nhập bị sao chép giữ cùng refresh token với bản gốc, và máy chủ xoay vòng token đó ở lần gia hạn đầu tiên: khi đó mọi nơi khác đang giữ nó, kể cả `~/.claude` nguồn, đều bị đăng xuất khỏi máy chủ đó. Vì vậy, với mỗi máy chủ được sao chép mà nguồn đã đăng nhập, dòng `sign-ins` của kế hoạch liệt kê các tài khoản của nhóm phải tự đăng nhập bằng `/mcp`, và phần tóm tắt nhắc lại điều đó. Một mục không có token thì không phải là lần đăng nhập và không được liệt kê, tài khoản đã có lần đăng nhập riêng cũng vậy, và cả máy chủ mà nhóm giữ định nghĩa của mình. Không tệp đăng nhập nào được ghi.

## An toàn và các trường hợp đặc biệt

- Nguồn chỉ được đọc. Không tệp nào trong đó bị di chuyển, liên kết hay thay đổi, và thông tin đăng nhập của nguồn không bao giờ bị ghi, kể cả trên macOS, nơi nó có thể được đọc từ Chuỗi khóa (Keychain).
- Trước khi một tệp của nhóm bị thay thế, do `use mine`, do gộp, hoặc do có máy chủ MCP mới trong một `mcp.json` đã có, bản của nhóm được lưu trong `<group>/.ccam-backup/<time>/`. Tên thư mục là duy nhất cho mỗi lần chạy, nên hai lần chạy trong cùng một giây không ghi đè bản sao lưu của nhau. Phần tóm tắt in ra đường dẫn.
- Với `--yes`, mọi xung đột đều giữ bản của nhóm, và ccam cho biết cách để tự chọn: chạy lại cùng lệnh mà không có `--yes`.
- Nếu không có gì thay đổi, ccam in ra `Nothing to copy.`; nếu mọi xung đột đều giữ bản của nhóm, nó in ra `Nothing to copy: the group keeps its own copies.`
- Nếu việc sao chép thất bại giữa chừng, ccam dừng với `The source folder was not changed. Group files replaced so far are saved in` theo sau là thư mục sao lưu. Chạy lại sẽ sao chép các tệp chưa được sao chép trọn vẹn.
- Khi nhiều nhóm dùng chung MCP, các lệnh mà `ccam doctor`, `ccam ext` và `ccam mcp list` gợi ý có dạng `ccam migrate --mcp --group <a|b>`.
- Nguồn không tồn tại sẽ kết thúc với `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Liên quan

- [Dùng chung máy chủ MCP, cài đặt và lịch sử](/docs/sharing)
- [Kiểm tra và sửa chữa](/features/doctor)
- [Dùng nhiều tài khoản Claude Code trên một máy tính](/guides/multiple-claude-code-accounts)
- [Giải thích về CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
