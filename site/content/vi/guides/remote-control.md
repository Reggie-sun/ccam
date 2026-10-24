---
title: Dùng Remote Control với ccam
description: "Khởi động Remote Control bằng bất kỳ tài khoản nào với ccam <name> rc: tùy chọn áp dụng, vì sao ccam đề nghị sao chép máy chủ MCP dùng chung, và cách sửa sự cố."
updated: 2026-10-23
order: 4
translation: machine
sourceHash: "3e17d89f98c8"
---

`ccam work rc` khởi động Remote Control của Claude Code bằng tài khoản `work`, để bạn có thể điều khiển các phiên trên máy tính này từ claude.ai/code hoặc ứng dụng Claude trên điện thoại. Lệnh này hoạt động giống `claude rc` với một khác biệt: nếu tài khoản dùng chung máy chủ MCP, ccam hỏi một lần có sao chép chúng vào `.claude.json` của tài khoản hay không, vì các phiên Remote Control không đọc máy chủ từ nơi nào khác.

## Khởi động

```sh
ccam work rc                          # Remote Control as the account work
ccam 2 rc                             # the account with number 2 (see ccam list)
ccam work remote-control              # the long name of the same command
ccam work rc --name laptop            # options after rc go to Remote Control
ccam work rc --spawn worktree         # each new session in its own git worktree
ccam --dry-run work rc                # print the command instead of running it
```

[Lệnh tắt](/features/shortcuts) cũng hoạt động như vậy: `ccwork rc` chính là `ccam work rc`. Hãy chạy lệnh trong thư mục bạn muốn làm việc, rồi mở claude.ai/code hoặc ứng dụng Claude đã đăng nhập bằng cùng tài khoản Claude với `work`.

## Tùy chọn

Ở đây có hai loại tùy chọn: tùy chọn riêng của ccam, đặt trước tên tài khoản hoặc nằm trong cài đặt của ccam, và tùy chọn của Remote Control, đặt sau `rc` và được chuyển nguyên vẹn đến Claude Code.

### Tùy chọn và cài đặt của ccam

| Tùy chọn hoặc cài đặt | Tác dụng | Mặc định |
| --- | --- | --- |
| `--dry-run` (trước tên tài khoản) | In ra biến môi trường và lệnh, cùng thông báo MCP nếu có. Không hỏi và không sao chép | tắt |
| `Copy them into .claude.json (a backup is kept; servers of the account's own stay)? [y/N]` | Được hỏi khi tài khoản dùng chung máy chủ MCP và `.claude.json` của nó chưa có danh sách hiện tại của nhóm, hoặc một máy chủ đã sao chép bị mất khỏi đó. `y` sao chép chúng, `N` hoặc Enter khởi động mà không có chúng | `N` |
| `mcp_strategy = "sync"` ở đầu `~/.config/ccam/config.toml` | Sao chép các máy chủ của nhóm vào `.claude.json` trước mỗi lần khởi động, dù có phải Remote Control hay không, nên không bao giờ phải hỏi | `"flag"` |
| `ccam group leave <group> work` | Đưa tài khoản ra khỏi nhóm đang dùng chung MCP với nó. Các máy chủ của nhóm được sao chép vào `.claude.json` của nó như máy chủ riêng, nên Remote Control có chúng và không hỏi gì; với `--fresh`, nó khởi động chỉ với các máy chủ riêng của tài khoản | trong `default` |

### Tùy chọn của Remote Control

Mọi thứ sau `rc` thuộc về Claude Code. Đây là các tùy chọn của `claude remote-control` trong Claude Code 2.1.304; `ccam work rc --help` in ra danh sách cho phiên bản của bạn.

| Tùy chọn | Tác dụng |
| --- | --- |
| `--name <name>` | Tên phiên hiển thị trong claude.ai/code |
| `--remote-control-session-name-prefix <prefix>` | Tiền tố cho tên phiên được tạo tự động (mặc định: tên máy) |
| `-c`, `--continue` | Kết nối lại với phiên mà Remote Control ghi nhận gần nhất cho thư mục này, trong khoảng 4 giờ |
| `--session-id <id>` | Kết nối lại với một phiên cụ thể |
| `--permission-mode <mode>` | Chế độ quyền cho các phiên mới: `acceptEdits`, `auto`, `bypassPermissions`, `default`, `dontAsk` hoặc `plan` |
| `--spawn <mode>` | `same-dir` (mặc định), `worktree` (mỗi phiên một git worktree) hoặc `session` (một phiên, thoát khi phiên kết thúc) |
| `--capacity <N>` | Số phiên tối đa cùng lúc ở chế độ `same-dir` và `worktree` (mặc định 32) |
| `--[no-]create-session-in-dir` | Tạo phiên đầu tiên trong thư mục hiện tại khi khởi động (mặc định bật) |
| `--[no-]chrome` | Claude in Chrome cho các phiên mới (mặc định tắt) |
| `-d`, `--debug`, `--debug-file <path>`, `-v`, `--verbose` | Đầu ra gỡ lỗi và đầu ra chi tiết |

ccam đọc động từ của Claude Code theo cách Claude Code làm, bỏ qua các tùy chọn cùng giá trị của chúng, nên cả `ccam work rc --name laptop` lẫn `ccam work --model opus rc` đều khởi động Remote Control. Một tùy chọn mà ccam không biết được tính là không có giá trị; để chắc chắn, hãy đặt `rc` ngay sau tên tài khoản.

## ccam làm gì

Lần khởi động đi theo đường thông thường ([Khởi động Claude Code](/features/launch)) với một thay đổi ở bước MCP:

1. ccam tìm tài khoản, kiểm tra các liên kết dùng chung của nó và xác định `claude`, như mọi lần khởi động khác.
2. Nó thấy động từ là `rc` hoặc `remote-control` nên bỏ `--mcp-config` đi.
3. Nếu tài khoản dùng chung máy chủ MCP, nó so sánh `mcp.json` của nhóm với `sync-stamp` trong `~/.local/state/ccam/mcp/work/`, và kiểm tra rằng mọi máy chủ của nhóm vẫn còn trong `.claude.json`. Khi `.claude.json` đã chứa các máy chủ hiện tại, hoặc nhóm không có máy chủ nào, nó chuyển thẳng đến bước 6.
4. Nếu không, nó nêu tên các máy chủ mà các phiên sẽ thiếu:

   ```console
   $ ccam work rc
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept; servers of the account's own stay)? [y/N]
   ```

5. Khi chọn `y`, nó lưu một bản sao lưu của `.claude.json`, lấy khóa riêng của Claude Code, thêm các máy chủ của nhóm vào `mcpServers` và giữ lại các máy chủ khác của tài khoản; một máy chủ riêng của tài khoản có tên của một máy chủ nhóm vẫn là của tài khoản, và ccam thông báo điều đó. Khi chọn `N`, khi không có terminal, hoặc khi tài khoản đang chạy, nó không thay đổi gì và thông báo điều đó:

   ```console
         Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
   ```

6. Nó ghi nhận phiên, đặt `CLAUDE_CONFIG_DIR` thành thư mục của tài khoản và tự thay thế mình bằng `claude rc` thông qua `exec`.

Sau đó Remote Control tiếp tục chạy như một máy chủ. Mỗi phiên nó khởi động cho bạn là tiến trình con của tiến trình đó, nên chạy dưới danh nghĩa `work`: với thông tin đăng nhập, cài đặt và lịch sử hội thoại của `work`, và được tính vào giới hạn sử dụng của `work`. Trong khi một máy chủ do ccam khởi động đang chạy, ccam coi tài khoản là đang chạy, và với một máy chủ được khởi động bằng `claude remote-control` thông thường cũng vậy khi nó có phiên; một máy chủ như vậy chưa có phiên nào thì không được nhận ra.

Chạy thử cho thấy sự khác biệt. Một lần khởi động bình thường mang các máy chủ của nhóm dưới dạng đối số, còn Remote Control thì không:

```console
$ ccam --dry-run work
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json
$ ccam --dry-run work rc
ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
      Starting without them. To copy them before every launch instead, set mcp_strategy = "sync" in ~/.config/ccam/config.toml
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work CCAM_PROFILE=work /home/me/.local/bin/claude rc
```

## Vì sao lại hoạt động như vậy

**Remote Control không nhận `--mcp-config`.** Với các lần khởi động bình thường, ccam trao các máy chủ của nhóm cho Claude Code dưới dạng `--mcp-config=<path>` và không bao giờ ghi vào tệp của Claude Code. `claude remote-control` từ chối tùy chọn đó khi đặt trước tên động từ của nó và không nhận tùy chọn nào sau đó, nên trước ccam 0.6.0, `ccam work rc` thất bại với mọi tài khoản dùng chung máy chủ MCP. Giờ đây ccam bỏ tùy chọn này đi.

**Các phiên của nó chỉ đọc máy chủ từ `.claude.json`.** Các phiên mà Remote Control khởi động không kế thừa tùy chọn từ dòng lệnh. Nơi duy nhất chúng tìm thấy máy chủ MCP người dùng là `.claude.json` của tài khoản. Tệp đó còn chứa thông tin đăng nhập của tài khoản và trạng thái theo từng thư mục, nên ccam không thể liên kết nó tới một bản dùng chung như cách làm với `settings.json`. Sao chép các máy chủ vào đó là cách duy nhất để các phiên có được danh sách dùng chung.

**ccam hỏi trước khi ghi.** Với chiến lược flag mặc định, ccam không động đến `.claude.json` khi khởi động. Remote Control là lần khởi động duy nhất nó buộc phải ghi vào đó, nên nó hỏi trước, giữ một bản sao lưu và không động đến gì ngoài `mcpServers`. (Ngoài các lần khởi động, việc rời hoặc tham gia dùng chung MCP và `ccam mcp import` cũng ghi vào đó, như kế hoạch của chúng nêu rõ.)

**Nó không ghi khi tài khoản đang chạy.** Claude Code ghi lại `.claude.json` trong suốt phiên, và một thay đổi từ bên ngoài sẽ bị mất hoặc bị ghi đè. Vì thế ccam bỏ qua việc sao chép khi có bất kỳ phiên nào của tài khoản đang chạy, kể cả một máy chủ Remote Control khác và Claude Code được khởi động không qua ccam.

**Nó hỏi một lần cho mỗi thay đổi.** Sau khi sao chép, `sync-stamp` lưu thời điểm sửa đổi của `mcp.json` của nhóm. ccam chỉ hỏi lại khi tệp đó thay đổi, ví dụ sau `ccam mcp add`, hoặc khi một máy chủ đã sao chép bị mất khỏi `.claude.json`. Nếu bạn chọn `N`, vẫn chưa có bản sao nào, nên nó hỏi ở mỗi lần khởi động Remote Control.

**`ccam mcp` giữ cho các bản sao luôn mới.** Một mục trong `.claude.json` được ưu tiên hơn `--mcp-config` khi trùng tên, nên một bản sao cũ cũng sẽ che mất định nghĩa mới hơn ở các lần khởi động bình thường. Vì vậy `ccam mcp add`, `add-json`, `remove` và `import`, cùng `ccam migrate`, cập nhật bản sao của mọi tài khoản đang không chạy, và thông báo khi một bản sao đã thay đổi:

```console
$ ccam mcp add ctx --url https://mcp.context7.com/mcp
✓ ctx added in group default — applies from the next launch
✓ Updated the copy in account work's .claude.json (made for Remote Control)
```

**Không có terminal thì câu trả lời là không.** Remote Control thường chạy từ một script, một dịch vụ hoặc `nohup`, nơi không ai trả lời được. Khi đó ccam khởi động mà không có các máy chủ thay vì chờ đợi.

## Sự cố và cách khắc phục

### Các máy chủ MCP dùng chung không có trong phiên Remote Control

Câu hỏi đã được trả lời `N`, không có terminal, hoặc tài khoản đang chạy. Hãy dừng Remote Control, đóng các phiên khác của tài khoản và chạy lại `ccam work rc` trong terminal; chọn `y`. Để không bao giờ bị hỏi, đặt `mcp_strategy = "sync"` trong `~/.config/ccam/config.toml`. Các máy chủ trong một phiên không thay đổi khi phiên đang chạy, nên hãy bắt đầu phiên mới sau khi sao chép.

### "Account work is running, so its .claude.json is left alone"

Một phiên của `work` vẫn đang mở: một `ccam work` bình thường, một lệnh tắt hoặc một máy chủ Remote Control khác. Thoát phiên đó và chạy lại `ccam work rc`. Các phiên do tài khoản khác chạy không được tính.

### ccam lần nào cũng hỏi

Bạn chưa sao chép các máy chủ. Hãy chọn `y` một lần, hoặc đặt `mcp_strategy = "sync"`. Nếu bạn không bao giờ muốn có các máy chủ đó trong Remote Control, hãy ngừng dùng chung MCP với tài khoản mà không giữ lại chúng: `ccam group leave default work --fresh`, hoặc để MCP trong một nhóm riêng mà tài khoản không thuộc về. Một lệnh `ccam group leave` thông thường sao chép các máy chủ của nhóm vào `.claude.json` của tài khoản, nên Remote Control có chúng.

### `ccam work rc` thất bại ngay với lỗi về `--mcp-config`

ccam của bạn cũ hơn 0.6.0 và truyền `--mcp-config` cho Remote Control. Hãy chạy `ccam update`. Cùng lỗi đó với ccam phiên bản hiện tại có nghĩa là một tùy chọn mà ccam không biết, có giá trị, đã đứng trước `rc`, nên ccam hiểu giá trị đó là động từ; hãy đặt `rc` ngay sau tên tài khoản.

### Thay đổi bằng `ccam mcp` không đến được Remote Control

`ccam mcp add` và `remove` chỉ cập nhật bản sao ở các tài khoản không đang chạy, và có thông báo: `Note: account work is running, so the copy in its .claude.json stays as it was until its next Remote Control launch.` Hãy dừng Remote Control rồi khởi động lại; ccam thấy `mcp.json` đã thay đổi và hỏi lại.

### Một máy chủ riêng của tôi được dùng thay cho máy chủ của nhóm

Nếu tài khoản có máy chủ riêng trùng tên với một máy chủ của nhóm, bản sao để nguyên nó và thông báo điều đó: `account work has its own MCP servers fs, so group default's of those names are not copied there`. Để dùng máy chủ của nhóm, hãy gỡ máy chủ của tài khoản (`claude mcp remove -s user <name>`); để đưa phiên bản của bạn cho mọi tài khoản, hãy chuyển nó vào nhóm bằng `ccam mcp import work <server> --force`. Các bản sao lưu của `.claude.json` nằm trong `~/.local/state/ccam/mcp/work/backups/`, nơi giữ năm bản mới nhất.

### Tôi muốn gỡ các máy chủ đã sao chép

`ccam mcp remove <server>` gỡ một máy chủ khỏi nhóm và khỏi mọi bản sao. Để trở về trạng thái trước khi sao chép, hãy khôi phục bản sao lưu mới nhất từ `~/.local/state/ccam/mcp/work/backups/` khi không có phiên nào của tài khoản đang chạy.

### Remote Control hỏi có tin cậy thư mục không, hoặc thoát khi không có terminal

Trạng thái tin cậy thư mục được lưu riêng cho từng tài khoản trong `.claude.json`, nên một thư mục bạn đã tin cậy với `personal` vẫn là mới với `work`. Hãy chạy `ccam work` một lần trong thư mục đó và chấp nhận, rồi khởi động Remote Control ở đó.

### "You must be logged in with a Claude account that has a subscription"

Remote Control cần đăng nhập Claude có gói đăng ký. Hãy kiểm tra email của tài khoản trong `ccam list`, và đăng nhập lại bằng `ccam login work` nếu cần.

### Phiên không xuất hiện trên điện thoại

Ứng dụng Claude hoặc claude.ai/code đang đăng nhập bằng một tài khoản Claude khác. Hãy so sánh với email mà `ccam list` hiển thị cho `work`.

### ccam cảnh báo rằng một Claude Code khác đang dùng chung lịch sử hội thoại trong thư mục này

Một tài khoản khác trong cùng nhóm lịch sử đang chạy trong cùng thư mục. Các phiên tiếp tục cùng một cuộc hội thoại từ cả hai tài khoản có thể làm lẫn lộn lịch sử. Chỉ chọn `y` nếu chúng sẽ làm việc với các cuộc hội thoại khác nhau, hoặc khởi động Remote Control ở thư mục khác.

## Liên quan

- [Khởi động Claude Code](/features/launch)
- [Máy chủ MCP dùng chung](/features/mcp)
- [Dùng chung giữa các tài khoản](/features/share)
- [Quản lý tài khoản](/features/accounts)
