---
title: Sao chép thiết lập hiện có
description: "Sao chép cài đặt, CLAUDE.md, tiện ích mở rộng, máy chủ MCP và đăng nhập MCP của một ~/.claude hiện có vào một nhóm của ccam bằng ccam migrate."
updated: 2026-10-19
order: 6
translation: machine
sourceHash: "7cee52e75749"
---

Nếu bạn đã dùng Claude Code trước khi có ccam, thiết lập của bạn nằm trong `~/.claude` và `~/.claude.json`. `ccam migrate` sao chép nó vào một nhóm dùng chung, để mọi tài khoản dùng chung các mục đó qua nhóm đều bắt đầu với thiết lập bạn đã có. Thư mục nguồn chỉ được đọc, và `claude` thông thường vẫn hoạt động như trước. Hãy chạy lệnh này một lần sau `ccam add` đầu tiên, và chạy lại mỗi khi có tài khoản bắt đầu dùng chung máy chủ MCP.

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
 ? Choose what to copy (Space to toggle)
 › [x] MCP servers
   [x] settings.json
   [x] global CLAUDE.md
   [x] agents, commands, skills
   [ ] conversation history

 Enter show the plan · Esc cancel
```

Space (hoặc `x`) bật/tắt mục dưới con trỏ, các phím mũi tên hoặc `j`/`k` để di chuyển, Enter hiển thị kế hoạch, còn Esc, `q` hoặc Ctrl+C in ra `Cancelled.`. Những mục mà nguồn không có gì được ghi là `not found`. Nếu bạn bỏ chọn mọi thứ, ccam in ra `Nothing selected.` và dừng lại.

Sau đó ccam in ra kế hoạch và hỏi trước khi sao chép:

```console
$ ccam migrate --yes
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: fs, linear
              sign-ins        linear → personal; skipped: work (not logged in yet)
  settings    settings.json   new: 61 B
  memory      CLAUDE.md       new: 11 B
  extensions  agents/         add 1 file the group lacks (6 B)
              skills/         add 1 file the group lacks (13 B)
              not found, skipped: commands/, plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, skills/, MCP sign-ins for personal
MCP sign-ins were not copied to: work (not logged in yet). Sign in to those servers there with /mcp.
Accounts that share these through group default:
  mcp         personal, work
  settings    personal, work
  memory      personal, work
  extensions  personal, work
```

Một nhóm chỉ đến được các tài khoản của chính nó, và chỉ với các mục nó dùng chung, nên phần tóm tắt liệt kê các tài khoản của nhóm. Khi nhóm chưa dùng chung một trong các mục đã sao chép, hoặc có tài khoản chưa thuộc nhóm, ccam in ra các lệnh `ccam group set <group> --add <items>` và `ccam group join <group> <accounts>` để thay đổi điều đó.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `[folder]` | Thư mục Claude Code để sao chép từ đó | `~/.claude` |
| `--settings` | Sao chép `settings.json` | tắt |
| `--memory` | Sao chép `CLAUDE.md` toàn cục | tắt |
| `--extensions` | Sao chép `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` | tắt |
| `--mcp` | Sao chép máy chủ MCP người dùng (`mcpServers` trong `.claude.json`) | tắt |
| `--history` | Sao chép lịch sử hội thoại (`projects/`, `history.jsonl`) | tắt |
| `-g`, `--group string` | Nhóm dùng chung để sao chép vào | `default` |
| `-y`, `--yes` | Sao chép mà không hỏi; khi xung đột, bản của nhóm được giữ | tắt |
| `--dry-run` (toàn cục) | In ra kế hoạch và dừng trước mọi câu hỏi hay việc sao chép | tắt |

Khi không có cờ chọn mục, trên terminal danh sách kiểm sẽ quyết định. Ngoài terminal, hoặc với `--yes`, ccam sao chép mọi thứ trừ lịch sử hội thoại. Lịch sử bị loại theo mặc định vì `projects/` có thể rất lớn.

`--group on` nghĩa là `default`. `--group off` và `--group local` bị từ chối với `--group takes a group name, not "off"`. `ccam migrate` hỗ trợ `--dry-run`, và danh sách kiểm vẫn xuất hiện khi chạy thử để bạn có thể xem trước một lựa chọn.

## Cách hoạt động

`internal/cli/migrate.go` chọn các mục, rồi `migrate.Build` đọc nguồn và nhóm và trả về một kế hoạch. Không có gì được ghi cho đến `Plan.Apply`, sau khi bạn xác nhận.

1. **Tìm nguồn.** `checkSource` từ chối một đường dẫn không tồn tại, không phải thư mục, hoặc không chứa mục Claude Code nào đã biết. `migrate.ClaudeJSON` tìm `.claude.json` thuộc về thư mục: bên trong nó, hoặc với `~/.claude` không có tệp đó bên trong, là `~/.claude.json` nằm cạnh nó. Đó là nơi Claude Code giữ tệp này khi `CLAUDE_CONFIG_DIR` không được đặt.
2. **Ánh xạ mục sang đường dẫn.** Mọi đường dẫn đều đến từ `share.Layout`, cũng là danh sách cho phép mà các nhóm dùng. Các máy chủ MCP được đọc ra khỏi `.claude.json` bằng `mcp.LocalServers`, hàm này dùng bộ quét `profile.TopLevelValue` thay vì phân tích cả tệp.
3. **Sao chép, không bao giờ di chuyển.** Các tệp đi vào `~/.local/share/ccam/groups/<group>/` bằng `fsutil.AtomicWrite` hoặc `fsutil.CopyFile`; máy chủ MCP đi vào `mcp.json` của nhóm.

Mỗi mục ánh xạ một đường dẫn nguồn tới một đường dẫn trong nhóm:

<!-- motion:migrate-items -->

Hai hàng lịch sử có nét đứt vì chúng chỉ được sao chép với `--history`. `.credentials.json` và `.claude.json` được liệt kê trong `share.Never` và không bao giờ đến được nhóm: `Build` trả về lỗi nội bộ nếu một trong hai xuất hiện trong cấu trúc, và việc duyệt thư mục bỏ qua chúng theo tên.

`buildEntry` gán cho mỗi đường dẫn một trạng thái, đó là điều mỗi dòng trong kế hoạch mô tả:

<!-- motion:migrate-states -->

- **Missing**: không có trong nguồn, hoặc là một tệp ở nơi lẽ ra là thư mục. Được liệt kê một lần cho mỗi mục dưới dạng `not found, skipped`.
- **Linked**: đường dẫn nguồn đã trỏ tới bản của nhóm, ví dụ khi nguồn là một thư mục hồ sơ của ccam.
- **Same**: nhóm đã có nội dung giống hệt.
- **New**: nhóm chưa có, nên được sao chép toàn bộ.
- **Merge**: một thư mục có ở cả hai phía. `copyMissing` chỉ thêm các tệp nhóm còn thiếu; các tệp khác nhau giữ nguyên như ở nhóm và kế hoạch ghi `group copies kept`. Với `plugins/`, các tệp đăng ký `installed_plugins.json` và `known_marketplaces.json` được gộp, và các đường dẫn plugin tuyệt đối được ghi lại để trỏ tới thư mục của nhóm.
- **Conflict**: một tệp có ở cả hai phía với nội dung khác nhau. Bạn chọn cho từng tệp: `keep group`, `use mine`, cùng với `merge JSON` cho các tệp `.json` (gộp sâu, trong đó nhóm thắng khi giá trị xung đột) hoặc `merge by line` cho `history.jsonl`. Một máy chủ MCP có định nghĩa khác cũng là xung đột, với `keep group` hoặc `use mine`. Lựa chọn đầu tiên được liệt kê là mặc định khi bạn nhấn Enter; `--yes` luôn giữ bản của nhóm.

Bản thân danh sách kiểm là `tui.MigratePicker`, được dựng trên thành phần `checklist` mà `ccam add` dùng; `migrate.Found` đánh dấu các mục còn thiếu bằng các lời gọi lstat và danh sách máy chủ MCP.

### Thông tin đăng nhập MCP

Các máy chủ MCP đăng nhập bằng OAuth giữ token trong tệp đăng nhập, `.credentials.json` dưới `mcpOAuth`, không phải trong danh sách máy chủ. Claude Code đặt khóa cho mỗi mục là `<name>|` cộng với 16 ký tự đầu của mã băm `sha256` từ `type`, `url` và `headers` của máy chủ. Khóa không phụ thuộc vào thư mục, nên một mục được sao chép vẫn hoạt động ở tài khoản khác miễn là định nghĩa máy chủ giống hệt. Vì vậy `buildLogins` trong `internal/migrate/logins.go` chỉ sao chép các mục của những máy chủ mà cuối cùng nhóm có theo định nghĩa của nguồn: mới, đã giống nhau, hoặc xung đột mà bạn chọn `use mine`. Thông tin đăng nhập của các máy chủ plugin (khóa bắt đầu bằng `plugin:`) được sao chép theo khi bạn sao chép tiện ích mở rộng và nguồn có `plugins/`.

<!-- motion:migrate-signins -->

Các mục được ghi vào `.credentials.json` riêng của từng tài khoản nhận, bằng `fsutil.AtomicWrite` với quyền `0600`. Chúng không bao giờ đi vào nhóm. Một tài khoản chỉ nhận chúng nếu nó dùng chung `mcp` (hoặc, với máy chủ plugin, `extensions`) qua nhóm đích, và chỉ những khóa nó còn thiếu được thêm vào: các mục nó đã có, có thể mới hơn, được giữ nguyên. Dòng `sign-ins` trong kế hoạch liệt kê các máy chủ, các tài khoản nhận và các tài khoản bị bỏ qua. ccam bỏ qua một tài khoản, và nêu tên nó kèm lý do, khi tài khoản đó:

- được đăng ký tại chỗ bằng `ccam import`, vì `claude` thông thường có thể đang dùng tệp đăng nhập của nó;
- đang chạy, vì Claude Code ghi lại tệp đăng nhập của nó (được kiểm tra lại ngay trước khi ghi, và báo là `started meanwhile`);
- chưa đăng nhập;
- giữ thông tin đăng nhập trong Chuỗi khóa (Keychain) của macOS, nơi ccam không ghi vào.

Hãy đăng nhập vào các máy chủ ở tài khoản bị bỏ qua bằng `/mcp`. Tài khoản bắt đầu dùng chung máy chủ MCP về sau sẽ nhận được thông tin đăng nhập khi bạn chạy lại `ccam migrate`; phần tóm tắt nhắc điều này bằng ``Then run `ccam migrate` again to give them the MCP sign-ins too.``

## An toàn và các trường hợp đặc biệt

- Nguồn chỉ được đọc. Không tệp nào trong đó bị di chuyển, liên kết hay thay đổi, và thông tin đăng nhập của nguồn không bao giờ bị ghi, kể cả trên macOS, nơi nó có thể được đọc từ Chuỗi khóa (Keychain).
- Trước khi một tệp của nhóm bị thay thế, do `use mine`, do gộp, hoặc do có máy chủ MCP mới trong một `mcp.json` đã có, bản của nhóm được lưu trong `<group>/.ccam-backup/<time>/`. Tên thư mục là duy nhất cho mỗi lần chạy, nên hai lần chạy trong cùng một giây không ghi đè bản sao lưu của nhau. Phần tóm tắt in ra đường dẫn.
- Với `--yes`, mọi xung đột đều giữ bản của nhóm, và ccam cho biết cách để tự chọn: chạy lại cùng lệnh mà không có `--yes`.
- Nếu không có gì thay đổi, ccam in ra `Nothing to copy.`; nếu mọi xung đột đều giữ bản của nhóm, nó in ra `Nothing to copy: the group keeps its own copies.`
- Nếu việc sao chép thất bại giữa chừng, ccam dừng với `The source folder was not changed. Group files replaced so far are saved in` theo sau là thư mục sao lưu.
- Nguồn không tồn tại sẽ kết thúc với `no Claude Code folder at <path>. Give the folder to copy from: ccam migrate <folder>`.

## Liên quan

- [Dùng chung máy chủ MCP, cài đặt và lịch sử](/docs/sharing)
- [Kiểm tra và sửa chữa](/features/doctor)
- [Nhiều tài khoản Claude Code](/guides/multiple-claude-code-accounts)
- [CLAUDE_CONFIG_DIR](/guides/claude-config-dir)
