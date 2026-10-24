---
title: Máy chủ MCP dùng chung
description: "Giữ một danh sách máy chủ MCP cho mỗi nhóm bằng ccam mcp add, add-json, list, remove và import, được truyền qua --mcp-config hoặc đồng bộ vào từng tài khoản."
updated: 2026-10-23
order: 5
translation: machine
sourceHash: "fefca7f8d857"
---

Claude Code lưu máy chủ MCP người dùng trong `.claude.json`, cũng là tệp chứa thông tin đăng nhập của tài khoản, nên tệp này không thể dùng chung giữa các tài khoản. ccam lưu các máy chủ dùng chung trong một tệp của nhóm, `mcp.json`, và trao nó cho mọi tài khoản trong nhóm khi khởi động. Hãy dùng `ccam mcp` để thêm một máy chủ một lần cho mọi tài khoản của bạn, và `ccam mcp import` để chuyển một máy chủ bạn đã thêm trong một phiên sang các tài khoản khác.

## Cách dùng

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code   # a local (stdio) server
ccam mcp add ctx --url https://mcp.context7.com/mcp --header 'Authorization: Bearer …'
ccam mcp add docs --url https://example.com/sse --transport sse
ccam mcp add db -e PGHOST=localhost -- pg-mcp                    # with an environment variable
ccam mcp add-json gh '{"type":"http","url":"https://api.githubcopilot.com/mcp/"}'
ccam mcp add tools -g team -- tools-mcp                           # into the group team
ccam mcp list                                                     # every group, and servers only one account or folder has
ccam mcp remove ctx
ccam mcp import work                                              # move work's own servers into its group
ccam mcp import work gh --force                                   # replace the group's different gh
```

Một thay đổi có hiệu lực từ lần khởi động tiếp theo của mỗi tài khoản; các phiên đang chạy giữ các máy chủ mà chúng khởi động cùng:

```console
$ ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
✓ fs added in group default — applies from the next launch
$ ccam mcp list
group default (work, side)
  ctx                  http https://mcp.context7.com/mcp
  fs                   stdio npx -y @modelcontextprotocol/server-filesystem /home/me/code
only in account work: pg
only in account side, for /home/me/code/app only: notes
```

Bên dưới các nhóm, `ccam mcp list` nêu tên mọi máy chủ không đến được các tài khoản khác: `only in account work: pg` cho một máy chủ mà một tài khoản tự thêm cho mình, `only in account side, for /home/me/code/app only: notes` cho một máy chủ được thêm cho một thư mục (`claude mcp add` thông thường, xem bên dưới), `account work has its own …, used in its sessions in place of the group's` khi một máy chủ riêng của tài khoản có tên của một máy chủ nhóm, `.mcp.json` của thư mục này cùng các tài khoản đang `waiting for approval in` phê duyệt nó, và các máy chủ `in no account, only in ~/.claude.json (plain claude's)`, kèm lệnh [`ccam migrate --mcp`](/features/migrate) để sao chép chúng.

Thêm một tên mà nhóm đã có sẽ thay định nghĩa của nó và in ra `updated` thay vì `added`. Máy chủ chỉ đến được một tài khoản nếu tài khoản đó thuộc nhóm và nhóm dùng chung MCP (`ccam group set <group> --add mcp`); khi nhóm chưa dùng chung MCP, `add` sẽ báo điều đó.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `-g`, `--group <group>` | nhóm đích cho `add`, `add-json` và `remove` | nhóm duy nhất mà một tài khoản nào đó dùng chung MCP qua, nếu không thì `default` |
| `-- <command> [args…]` (`add`) | chạy lệnh này như một máy chủ stdio | |
| `--url <url>` (`add`) | URL máy chủ từ xa (http/sse) | |
| `--transport <stdio\|http\|sse>` (`add`) | giao thức truyền của máy chủ `--url` | `http` |
| `--header <'Name: value'>` (`add`) | HTTP header, có thể lặp lại | không có |
| `-e`, `--env <KEY=VALUE>` (`add`) | biến môi trường cho máy chủ stdio, có thể lặp lại | không có |
| `--force` (`import`) | thay một máy chủ cùng tên nhưng khác định nghĩa trong nhóm; `mcp.json` cũ được sao lưu | tắt |
| `--keep-local` (`import`) | giữ các mục trong `.claude.json` của tài khoản | tắt |
| `-h`, `--help` | trợ giúp cho từng lệnh | |
| `--dry-run` (toàn cục) | chỉ `mcp list` hỗ trợ; `add`, `add-json`, `remove` và `import` từ chối nó và không thay đổi gì | tắt |

`add` cần một lệnh sau `--` hoặc `--url`; nếu không, nó thất bại với `a command to run (after --) or --url is required`. `--header` áp dụng cho máy chủ `--url`, còn `--env` cho máy chủ dạng lệnh. `add-json` nhận cùng JSON như `claude mcp add-json`, và định nghĩa phải là một đối tượng có `command` hoặc `url`. `import` luôn dùng nhóm mà tài khoản dùng chung MCP qua, còn `list` hiển thị mọi nhóm, nên `-g` không thay đổi chúng.

## Cách hoạt động

Các tệp liên quan:

- `~/.local/share/ccam/groups/<group>/mcp.json`: các máy chủ dùng chung, `{"mcpServers": {…}}`, định dạng mà `--mcp-config` đọc. Tệp được ghi bằng `fsutil.AtomicWrite` với quyền `0600`, vì các mục máy chủ thường chứa khóa API, và chỉ dưới khóa `<group>/.mcp.json.lock`, nên hai lệnh `ccam mcp add` chạy cùng lúc đều giữ được máy chủ của mình.
- `~/.local/share/ccam/profiles/<name>/.claude.json`: tệp riêng của Claude Code chứa các máy chủ của tài khoản. ccam chỉ đọc giá trị `mcpServers` của nó bằng bộ quét `profile.TopLevelValue` và không bao giờ phân tích cả tệp trên đường khởi động.
- `~/.local/state/ccam/mcp/<name>/`: tệp `seen`, `sync-stamp` (kèm hash của từng bản sao ccam đã ghi, để nó phân biệt được bản sao chưa ai thay đổi với bản bạn đã sửa), `sync-checked` (những gì chiến lược sync thấy lần gần nhất), `kept` (các máy chủ mà tài khoản giữ lại khi rời việc dùng chung MCP) và các bản sao lưu `.claude.json`.

Có hai chiến lược, được đặt bằng `mcp_strategy` trong `~/.config/ccam/config.toml`. Mặc định là chiến lược flag:

<!-- motion:mcp-flag -->

1. `launcher.applyMCP` tìm nhóm MCP của tài khoản. `mcp.FlagArgs` kiểm tra bằng một lần `stat` rằng `mcp.json` của nhóm tồn tại, mà không phân tích nó.
2. `Plan.InsertArgs` đặt `--mcp-config=<path>` ngay sau `claude`, trước các đối số của bạn, nên `--mcp-config` hoặc `--strict-mcp-config` của riêng bạn vẫn được áp dụng chồng lên.
3. `mcp.NewLocal` kiểm tra `.claude.json` của tài khoản có thay đổi kể từ lần khởi động trước hay không, và chỉ khi đó mới phân tích các máy chủ của nó.
4. `syscall.Exec` thay ccam bằng Claude Code.

Đường dẫn được truyền dưới dạng một đối số, `--mcp-config=<path>`, chứ không phải hai. `--mcp-config` của Claude Code nhận nhiều giá trị, nên `claude --mcp-config a.json "prompt"` sẽ đọc lời nhắc như một tệp cấu hình thứ hai và thất bại với `MCP config file not found`. Dạng một đối số dừng lại ở đó, và lời nhắc hoặc lệnh con theo sau vẫn hoạt động. `ccam work mcp list`, `get` và `remove` không nhận cờ, vì Claude Code từ chối nó ở đó; thay vào đó ccam nêu tên các máy chủ của nhóm trên đầu ra lỗi chuẩn ([chi tiết](/features/launch#shared-mcp-servers)). Claude Code nạp các máy chủ từ `--mcp-config` bổ sung vào các máy chủ riêng của tài khoản; khi trùng tên, mục trong `.claude.json` được ưu tiên. Chiến lược flag chỉ ghi vào `.claude.json` cho Remote Control, khi một tài khoản rời hoặc quay lại việc dùng chung MCP ([bên dưới](#leaving-and-joining)), và cho `ccam mcp import`.

`claude remote-control` (`rc`) không nhận `--mcp-config`, và các phiên của nó chỉ đọc máy chủ từ `.claude.json`. Với `ccam work rc`, ccam bỏ cờ này đi và, khi có terminal, hỏi có sao chép các máy chủ của nhóm vào `.claude.json` hay không, kèm một bản sao lưu, như chiến lược `sync` vẫn làm. Nó hỏi cho đến khi có bản sao, sau đó chỉ hỏi lại khi `mcp.json` của nhóm thay đổi hoặc một máy chủ đã sao chép bị mất ([Remote Control](/features/launch#remote-control)). Vì bản sao như vậy được ưu tiên hơn `--mcp-config` khi trùng tên, `ccam mcp add`, `add-json`, `remove` và `import` cập nhật các bản sao mà các tài khoản khác trong nhóm đang giữ:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account side's .claude.json (made for Remote Control)
```

Dòng `✓ Updated` chỉ xuất hiện khi một bản sao thay đổi. Một bản sao do ccam cũ hơn tạo ra cũng được cập nhật theo cùng cách, và bản sao sau `ccam migrate` cũng vậy.

Một máy chủ được thêm trong phiên bằng `claude mcp add` chỉ nằm trong `.claude.json` của tài khoản đó. `claude mcp add` thông thường dùng phạm vi local: máy chủ chỉ được tải cho đúng tài khoản đó trong đúng thư mục đó (`-s user` làm nó thành của tài khoản trong mọi thư mục). ccam không thể thấy nó vào lúc đó, vì Claude Code đã thay thế ccam, nên nó kiểm tra ở lần khởi động tiếp theo. Tệp `seen` lưu thời điểm sửa đổi của `.claude.json` và tên các máy chủ đã được báo. Khi thời điểm khác đi, ccam so sánh các tên với nhóm và in ra một gợi ý cho mỗi máy chủ mới một lần, mà không chờ đầu vào, với cả hai chiến lược:

```console
ccam: MCP servers only account work has: pg
      To share them with other accounts: ccam mcp import work pg
ccam: MCP servers account work has for one folder only (claude mcp add without -s user): notes (/home/me/code/app)
      To share them with other accounts, in every folder: ccam mcp import work notes
```

Gợi ý này cũng nêu tên một máy chủ riêng của tài khoản có tên của một máy chủ nhóm, vì máy chủ của tài khoản được ưu tiên trong các phiên của nó.

Với `mcp_strategy = "sync"`, các máy chủ xuất hiện trong `claude mcp list` và `/mcp` như máy chủ người dùng thông thường, và ccam ghi chúng vào `.claude.json` của từng tài khoản trước khi khởi động:

<!-- motion:mcp-sync -->

1. `mcp.Sync` so sánh thời điểm sửa đổi của `mcp.json` của nhóm với `sync-stamp`, và `.claude.json` của tài khoản với những gì nó thấy lần gần nhất. Nếu cả hai đều chưa thay đổi, không có gì khác xảy ra. Nếu không, nó còn kiểm tra rằng mọi máy chủ của nhóm vẫn còn trong `.claude.json`, và đặt lại một bản sao đã mất (`put back group default's MCP servers that were gone from account work's .claude.json: ctx`).
2. Nó lấy khóa riêng của Claude Code, tức thư mục `.claude.json.lock` cạnh tệp, được tạo bằng `mkdir`. Nó chờ tối đa 2 giây và coi một khóa cũ hơn 10 giây là đã hết hiệu lực, giống như thư viện khóa của Claude Code.
3. Nó lưu `.claude.json` hiện tại vào `backups/` trong thư mục trạng thái và giữ 5 bản mới nhất.
4. Nó đặt mọi máy chủ của nhóm vào `mcpServers`, xóa các máy chủ nó đã ghi trước đó mà nhóm không còn có, và ghi tệp một cách nguyên tử với quyền ban đầu của nó, chỉ khi có gì đó thay đổi. Nó chỉ thay đổi hoặc xóa những bản sao nó đã ghi mà chưa ai thay đổi kể từ đó. Một máy chủ riêng của tài khoản có cùng tên với một máy chủ nhóm vẫn là của tài khoản, và ccam báo điều đó một lần: `account work has its own MCP servers fs, so group default's of those names are not copied there`.
5. Nó ghi lại dấu thời gian mới kèm hash của từng bản sao, rồi Claude Code khởi động mà không có `--mcp-config`.

Việc đồng bộ bị bỏ qua, với `ccam: account work is already running, skipping MCP sync`, khi có một phiên của tài khoản đang chạy, bất kể Claude Code được khởi động theo cách nào, vì Claude Code ghi lại `.claude.json` trong lúc chạy. Đồng bộ thất bại sẽ in ra `ccam: MCP sync failed (launching anyway): …` và vẫn khởi động.

### Rời và tham gia

Một tài khoản ngừng dùng chung MCP (`ccam group leave`, `ccam group set --remove mcp`, `ccam group remove --leave`) vẫn giữ các máy chủ của nhóm: kế hoạch sao chép chúng vào `.claude.json` của nó, kèm một bản sao lưu, và các máy chủ riêng cùng tên của nó vẫn được giữ. Với `--fresh`, nó không nhận máy chủ nào, và kế hoạch liệt kê các máy chủ sẽ ngừng được tải. Khi nó lại tham gia một nhóm có các máy chủ đó, các bản sao đã giữ mà chưa ai thay đổi sẽ bị xóa, để các chỉnh sửa sau này của nhóm đến được với nó.

`ccam mcp import` chuyển các máy chủ riêng của một tài khoản vào nhóm của nó:

<!-- motion:mcp-import -->

ccam đọc các máy chủ của tài khoản, các máy chủ riêng của nó và những máy chủ nó chỉ có cho một thư mục, và so sánh từng máy chủ với máy chủ cùng tên của nhóm bằng `mcp.SameDef`, hàm này bỏ qua thứ tự khóa và khoảng trắng. Một máy chủ của thư mục trở thành máy chủ của nhóm cho mọi thư mục (`notes loaded for one folder of account work only; from the group, every account in it gets them in every folder`); một tên có hai định nghĩa thư mục khác nhau sẽ bị từ chối. Máy chủ của plugin không được đề xuất, vì plugin tự mang chúng theo. Máy chủ mà nhóm chưa có hoặc có với cùng định nghĩa sẽ được sao chép. Máy chủ khác định nghĩa sẽ bị bỏ qua, vì định nghĩa của nhóm là thứ mọi tài khoản khác đang dùng:

```console
Skipped gh: group default already has a different gh (`ccam mcp import work gh --force` replaces it)
✓ Copied to group default: pg
✓ Removed account work's own entries (.claude.json backed up)
```

Với `--force`, máy chủ của nhóm bị thay thế, và `mcp.json` trước đó được lưu trước vào `.ccam-backup/<time>/mcp.json` bên trong thư mục nhóm. Sau đó `mcp.Remove` xóa các mục đã sao chép khỏi `.claude.json` của tài khoản, cả mục riêng lẫn mục của thư mục trong một lần ghi, dưới cùng khóa đó và sau khi sao lưu vào thư mục trạng thái. `mcp.json` của nhóm được giữ khóa từ lúc tải cho đến khi các mục đó bị xóa, để không ai xóa máy chủ khỏi nhóm ở giữa chừng. Việc này là cần thiết vì khi trùng tên, mục của tài khoản sẽ được ưu tiên hơn mục của nhóm. `--keep-local` giữ lại chúng.

## An toàn và các trường hợp đặc biệt

- **Đăng nhập.** Các máy chủ đăng nhập bằng OAuth giữ token trong `.credentials.json` của từng tài khoản (`mcpOAuth`), không phải trong danh sách máy chủ, nên dùng chung một máy chủ không có nghĩa là dùng chung thông tin đăng nhập của nó. Hãy đăng nhập một lần cho mỗi tài khoản bằng `/mcp`. ccam không bao giờ sao chép các lần đăng nhập, kể cả [`ccam migrate`](/features/migrate): một refresh token bị sao chép sẽ đăng xuất mọi nơi khác đang giữ nó. [`ccam ext`](/features/ext) và `ccam doctor` nêu tên các máy chủ mà một tài khoản chưa đăng nhập trong khi một tài khoản khác đã đăng nhập, và `ccam next --continue` cho biết tài khoản mới còn thiếu những máy chủ nào. Nhóm không bao giờ nhận tệp đăng nhập.
- **Tài khoản đang chạy.** `import` sao chép vào nhóm nhưng giữ nguyên các mục riêng của tài khoản đang chạy và báo điều đó; hãy chạy lại sau khi phiên kết thúc. Việc đồng bộ bỏ qua các tài khoản đang chạy. Cả hai đều tính cả Claude Code được khởi động không qua ccam.
- **Bản sao cho Remote Control.** Bản sao của tài khoản đang chạy không bị động đến: `Note: account side is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Với `mcp_strategy = "sync"` thì không cần cập nhật gì, vì mỗi lần khởi động đều đồng bộ.
- **Tệp bị hỏng.** Nếu `.claude.json` không phải JSON hợp lệ, việc đồng bộ không động đến nó và báo lỗi. Nếu không lấy được khóa trong 2 giây, ccam báo `could not lock … (Claude Code is writing the file)`.
- **Tài khoản không dùng chung MCP.** `import` từ chối và yêu cầu bạn đưa tài khoản vào một nhóm dùng chung MCP trước (`ccam group join <group> side`).
- **Tên nhóm.** `-g` chỉ nhận tên nhóm; tên này trở thành một thư mục trong `groups/`, nên các tên như `../x` bị từ chối.

## Liên quan

- [Dùng chung giữa các tài khoản](/features/share)
- [Dùng Remote Control với ccam](/guides/remote-control)
- [Dùng chung máy chủ MCP, cài đặt và lịch sử](/docs/sharing)
- [Sao chép thiết lập hiện có](/features/migrate)
- [Tất cả các lệnh](/docs/commands)
