---
title: Dùng chung bằng nhóm
description: "Một nhóm ccam dùng chung các mục đã chọn giữa các tài khoản qua một thư mục; ccam group chuyển dữ liệu theo kế hoạch, có sao lưu và lựa chọn khi xung đột."
updated: 2026-10-20
order: 4
translation: machine
sourceHash: "0c0338362ba2"
---

Mỗi tài khoản có thư mục riêng, nhưng phần lớn nội dung trong đó không cần tách riêng. Một **nhóm** cho biết những gì được dùng chung và do ai: mọi tài khoản trong nhóm liên kết các mục đó tới một thư mục, `~/.local/share/ccam/groups/<group>/`, nên tất cả cùng đọc và ghi một bản. Token đăng nhập không bao giờ được dùng chung.

Một tài khoản có thể thuộc nhiều nhóm, nhưng mỗi mục chỉ đến từ một nhóm. Một thiết lập điển hình là một nhóm cho cài đặt và tiện ích mở rộng, và một nhóm thứ hai cho lịch sử hội thoại giữa hai tài khoản:

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

Tài khoản đầu tiên bạn thêm sẽ tạo `default`, và các tài khoản thêm sau sẽ tham gia nhóm đó trừ khi bạn chỉ định khác.

## Cách dùng

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # every account in default moves together
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # work's CLAUDE.md goes in first
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # leave the group that shares them now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

Mọi lệnh di chuyển dữ liệu đều in ra một kế hoạch, hỏi về từng tệp khác với bản của nhóm, và hỏi `Proceed? [y/N]` trước khi thay đổi bất cứ thứ gì:

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261020-153012
  link     projects, history.jsonl: profiles/work → groups/default
Proceed? [y/N]
```

Các thay đổi tương tự có trên thẻ **Groups** của [bảng điều khiển](/features/dashboard), và `s` trên một tài khoản chọn các nhóm của nó.

## Các mục

| Mục | Bao gồm |
| --- | --- |
| `mcp` | Máy chủ MCP, qua `mcp.json` của nhóm ([Máy chủ MCP dùng chung](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | `CLAUDE.md` toàn cục |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` và `history.jsonl`, để một tài khoản khác có thể `--continue` |

Một mục không được nhóm nào của tài khoản dùng chung sẽ ở lại trong thư mục riêng của tài khoản.

## Lệnh và tùy chọn

| Lệnh | Tác dụng |
| --- | --- |
| `ccam group` hoặc `ccam group list` | liệt kê các nhóm |
| `ccam group show <group>` | những gì nhóm dùng chung, các tài khoản, thư mục và máy chủ MCP của nhóm |
| `ccam group add <group> --items <items> [--members <accounts>]` | tạo một nhóm; các mục của các tài khoản trong nhóm được chuyển vào đó |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | thay đổi những gì nhóm dùng chung, cho mọi tài khoản trong nhóm |
| `ccam group join <group> <account>...` | đưa các tài khoản vào nhóm |
| `ccam group leave <group> <account>...` | đưa các tài khoản ra khỏi nhóm; mỗi tài khoản giữ một bản sao |
| `ccam group rename <group> <new-name>` | chỉ đổi tên trong cấu hình |
| `ccam group remove <group> [--leave]` | xóa một nhóm |

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: tài khoản có bản được đưa vào nhóm trước | tài khoản mặc định, rồi theo số |
| `--replace` | `join`: trước tiên rời các nhóm đang dùng chung cùng các mục đó | tắt |
| `--fresh` | các tài khoản ngừng dùng chung một mục sẽ bắt đầu không có nó thay vì có một bản sao | tắt |
| `--leave` | `remove`: trước tiên đưa các tài khoản của nhóm ra | bắt buộc khi nhóm còn tài khoản |
| `--prefer <group\|mine\|merge>` | cách quyết định xung đột | hỏi; với `--yes`, mặc định của từng tệp |
| `-y`, `--yes` | chạy mà không hỏi; xung đột tuân theo `--prefer` | tắt |
| `--dry-run` (toàn cục) | in ra kế hoạch, không thay đổi gì | tắt |

Tài khoản được chỉ định bằng tên hoặc số; các mục là một danh sách phân tách bằng dấu phẩy. `ccam add` nhận `--group <group>` (có thể lặp lại) hoặc `--no-group`.

## Cách hoạt động

Các tệp liên quan:

- `~/.config/ccam/config.toml`: các mục `[[group]]`: `name`, `items`, `members` và, khi tên thư mục khác tên nhóm, `dir`.
- `~/.local/share/ccam/profiles/<name>/`: `CLAUDE_CONFIG_DIR` của tài khoản.
- `~/.local/share/ccam/groups/<folder>/`: các bản dùng chung.

Khi `config.toml` được đọc, các nhóm của mỗi tài khoản được quy về một thư mục cho mỗi mục. Việc này mất khoảng 25 µs với năm tài khoản, nên một lần khởi động vẫn chỉ đọc `config.toml`. Một cấu hình mà trong đó một tài khoản nhận một mục từ hai nhóm là lỗi, và không có gì khởi động cho đến khi nó được sửa.

Một thay đổi được thực hiện trên một bản sao của cấu hình:

1. Lệnh sửa bản sao (`AddGroup`, `SetGroupItems`, `JoinGroup` …), và việc này kiểm tra các quy tắc.
2. `internal/group` so sánh cũ và mới và, với mỗi tài khoản có các mục thay đổi, dựng một kế hoạch bằng `share.Build`. Các kế hoạch được dựng lần lượt, và mỗi kế hoạch thấy được những gì các kế hoạch trước sẽ đưa vào thư mục nhóm, nên chỉ tài khoản đầu tiên chuyển bản của nó vào. Bản đầu tiên đến từ thư mục nhóm nếu thư mục đó đã có, rồi từ `--from`, rồi tài khoản mặc định, rồi các tài khoản khác theo số.
3. ccam từ chối nếu một trong các tài khoản đó đang chạy, và trên Windows không có Developer Mode nếu một tệp cần đến liên kết.
4. Các xung đột được quyết định, qua câu hỏi hoặc qua `--prefer`, kế hoạch được in ra, và ccam hỏi.
5. Cấu hình mới được lưu trước khi bất cứ thứ gì được di chuyển. Nếu một bước dừng giữa chừng, `ccam doctor --fix` và một lần chạy lại sẽ dựng kế hoạch hướng tới các nhóm đã lưu và hoàn tất thay đổi.
6. Các kế hoạch chạy theo thứ tự. Mọi bước đều chuyển dữ liệu vào nhóm hoặc vào thư mục sao lưu; không có gì bị xóa.

`share.Layout` là một danh sách cho phép gồm các đường dẫn mà mỗi mục bao gồm. Bất cứ thứ gì không có trong đó đều không bao giờ đến được nhóm. `share.Never` nêu tên `.credentials.json` và `.claude.json`, và `Build` dừng với lỗi nội bộ nếu một trong hai xuất hiện trong cấu trúc. Mục `mcp` không có đường dẫn: danh sách máy chủ nằm bên trong `.claude.json`, nên [dùng chung MCP](/features/mcp) hoạt động qua `mcp.json` của nhóm.

Đây là điều xảy ra với một thư mục, `projects/`, khi `side` bắt đầu dùng chung lịch sử và nhóm đã có nội dung:

<!-- motion:share-on -->

Nếu nhóm chưa có gì ở đường dẫn đó, tài khoản là tài khoản đầu tiên dùng chung nó và thư mục của nó được chuyển nguyên vào đó. Ngược lại, từng tệp nhóm còn thiếu được chuyển vào nhóm, còn các tệp có ở cả hai phía vẫn nằm trong thư mục của tài khoản, nên phiên bản của nhóm được giữ. Sau đó thư mục của tài khoản được chuyển vào `.ccam-backup/<time>/` bên trong thư mục tài khoản, và một liên kết thế vào chỗ của nó. Các tệp hội thoại được đặt tên theo ID phiên, nên lịch sử từ nhiều tài khoản được gộp lại mà không xung đột. Các tệp đăng ký trong `plugins/` lưu đường dẫn tuyệt đối đi qua thư mục tài khoản; ccam trỏ chúng tới thư mục của nhóm và gộp các mục của tài khoản vào tệp đăng ký mà nhóm đã có.

Rời nhóm sẽ đảo ngược liên kết:

<!-- motion:share-off -->

ccam chỉ gỡ liên kết nếu nó trỏ vào thư mục groups, rồi sao chép nội dung hiện tại của nhóm vào tài khoản. Với `--fresh`, bước sao chép bị bỏ qua. Nhóm vẫn giữ nội dung của nó cho các tài khoản khác. Một tài khoản chuyển một mục sang một nhóm có thư mục vẫn còn trống sẽ mang theo bản sao của những gì nó đã thấy, nên nó giữ được các tệp của mình; máy chủ MCP cũng đi theo cùng cách đó. Một tài khoản ngừng dùng chung MCP sẽ mất các máy chủ mà ccam đã sao chép vào `.claude.json` của nó (có giữ bản sao lưu), vì chúng đến từ nhóm.

Một tệp đơn lẻ có ở cả hai phía với nội dung khác nhau là một xung đột:

<!-- motion:share-conflict -->

`history.jsonl` có `merge by line`, `keep group` và `use mine`, với gộp đứng đầu. Các tệp `.json` khác có `keep group`, `use mine` và `merge JSON`. Mọi tệp khác có `keep group` và `use mine`. `merge JSON` gộp các đối tượng một cách đệ quy, nối các mảng mà không trùng lặp, và giữ giá trị của nhóm khi hai phía đặt cùng một khóa khác nhau, nên các tài khoản khác không thấy thay đổi gì. Mọi lựa chọn trừ `keep group` trước tiên lưu bản của nhóm vào `.ccam-backup/<time>/group/` của tài khoản. Tệp của riêng bạn luôn được đưa vào bản sao lưu.

Ở mỗi lần khởi động, `share.Check` so sánh từng đường dẫn trong cấu trúc với các nhóm của tài khoản, chỉ dùng `lstat` và `readlink`. Liên kết còn thiếu được tạo ngay lập tức. Một tệp thật ở nơi lẽ ra là liên kết, một liên kết trỏ sai chỗ, hoặc một liên kết lỗi thời sẽ dừng việc khởi động và chỉ tới `ccam doctor --fix`.

### Đổi tên, xóa và định dạng cũ

- `ccam group rename` chỉ thay đổi cấu hình. Thư mục giữ nguyên tên và nhóm ghi lại tên đó dưới dạng `dir`, nên không có gì bị di chuyển và không có gì có thể mất.
- `ccam group remove` từ chối một nhóm còn tài khoản trừ khi bạn thêm `--leave`. Khi đó mỗi tài khoản giữ một bản sao, và thư mục được chuyển vào `~/.local/share/ccam/trash/` trừ khi một nhóm khác đang dùng nó.
- Trước khi có nhóm, `config.toml` lưu một bảng `[profile.share]` cho mỗi tài khoản. ccam đọc định dạng đó như các nhóm mà không di chuyển gì. Ở lần lưu đầu tiên, nó ghi định dạng mới và giữ tệp cũ dưới tên `config.toml.bak`. Các mục được dùng chung bởi cùng các tài khoản trở thành một nhóm mang tên thư mục; các mục còn lại nhận tên như `default-history` và tiếp tục dùng cùng thư mục.

### Trên Windows

Windows chỉ tạo liên kết tượng trưng trong Developer Mode (Settings > System > For developers) hoặc cho quản trị viên. Khi không có quyền đó, một thư mục dùng chung trở thành một directory junction, và hoạt động như nhau với Claude Code. Một tệp dùng chung (`settings.json`, `CLAUDE.md`, `history.jsonl`) chỉ có thể là liên kết tượng trưng, nên một thay đổi nhóm cần đến liên kết như vậy sẽ bị từ chối trước khi có gì được di chuyển.

## An toàn và các trường hợp đặc biệt

- **Tài khoản đang chạy.** Một thay đổi bị từ chối khi Claude Code đang chạy bằng một trong các tài khoản mà nó di chuyển, vì phiên vẫn tiếp tục ghi vào các tệp đó. Việc rời và tham gia nhóm chỉ liên quan đến tài khoản được nêu tên. `--dry-run` vẫn in ra kế hoạch.
- **Bản sao lưu.** Mỗi kế hoạch có thư mục `.ccam-backup/<time>/` riêng trong thư mục tài khoản. ccam không bao giờ xóa chúng.
- **Liên kết của riêng bạn.** Nếu một đường dẫn đã là liên kết do bạn tạo, chẳng hạn tới một kho dotfiles, và nhóm còn trống, ccam sao chép đích của liên kết vào nhóm và chuyển chính liên kết vào bản sao lưu. Đích không bị động đến.
- **Thư mục được nhập.** Một tài khoản được đăng ký tại chỗ bằng `ccam import` (chẳng hạn `~/.claude`) thay đổi ngay tại chỗ, nên lệnh `claude` thông thường cũng thấy thay đổi; kế hoạch có ghi rõ điều này.
- **Hai tài khoản trong một thư mục.** Các tài khoản dùng chung lịch sử ghi vào cùng các tệp. Khi một tài khoản thứ hai khởi động trong một thư mục nơi đã có một tài khoản đang chạy, ccam cảnh báo và hỏi.
- **Sửa bằng tay.** Một thành viên không phải là tài khoản sẽ bị bỏ qua và được `ccam doctor` báo cáo. Các tên `local`, `on`, `off` và `none` không thể làm tên nhóm.

## Liên quan

- [Bảng điều khiển](/features/dashboard)
- [Dùng chung máy chủ MCP, cài đặt và lịch sử](/docs/sharing)
- [Máy chủ MCP dùng chung](/features/mcp)
- [Sao chép thiết lập hiện có](/features/migrate)
- [Tất cả các lệnh](/docs/commands)
