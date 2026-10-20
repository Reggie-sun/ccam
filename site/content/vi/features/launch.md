---
title: Khởi động Claude Code
description: "Chạy Claude Code bằng tài khoản bất kỳ với ccam <name|number>, ccam run hoặc bộ chọn của ccam; đối số được chuyển nguyên vẹn cho claude, --dry-run in ra lệnh."
updated: 2026-10-20
order: 1
translation: machine
sourceHash: 401eed57e533
---

`ccam work` khởi động Claude Code bằng tài khoản `work`, và mọi thứ sau tên được chuyển nguyên vẹn cho `claude`. Hãy dùng nó ở bất cứ đâu bạn vẫn gõ `claude`: trong terminal, trong script, hoặc khi có đầu vào được pipe vào.

## Cách dùng

```sh
ccam work                        # Claude Code as the account work
ccam 2                           # the account with number 2 (see ccam list)
ccam work --resume               # claude --resume, as work
ccam work -p "fix the tests" --output-format json
cat log.txt | ccam work -p "find the cause"
ccam -p "summarise"              # no name: the default account
ccam use work                    # make work the default account
ccam run work --continue         # same as ccam work --continue
ccam                             # the picker on a terminal, the list otherwise
ccam --dry-run work -p hi        # print the command instead of running it
```

Hãy đưa tên của một tài khoản hoặc số của nó trong `ccam list`. ccam không đọc các đối số sau tên, nên `ccam work --help` hiển thị trợ giúp của Claude Code. Dấu `--` sau tên là không bắt buộc. Đầu vào chuẩn, mã thoát và tín hiệu hoạt động giống như với `claude`.

Khi đối số đầu tiên bắt đầu bằng `-` và không phải một cờ của chính ccam, ccam khởi động tài khoản mặc định với tất cả các đối số. Tài khoản mặc định là tài khoản được đặt bằng `ccam use`, hoặc tài khoản có số thấp nhất nếu bạn chưa đặt.

### Bộ chọn

`ccam` không kèm đối số sẽ mở bộ chọn khi đầu vào và đầu ra chuẩn là terminal, và in ra `ccam list` trong các trường hợp khác. Bộ chọn hiển thị số, email, mức sử dụng 5 giờ, lần dùng gần nhất và các mục dùng chung của từng tài khoản; mức sử dụng được điền vào khi từng tài khoản phản hồi.

| Phím | Tác dụng |
| --- | --- |
| `Enter` | Khởi động Claude Code bằng tài khoản đang chọn |
| `1` đến `9` | Khởi động tài khoản có số đó |
| `Up`/`Down`, `k`/`j` | Di chuyển lựa chọn |
| `a` | Đóng bộ chọn và chạy trình hướng dẫn `ccam add` |
| `s` | Sửa các mục dùng chung của tài khoản đang chọn (`Space` hoặc `x` để bật/tắt, `Enter` để áp dụng, `Esc` để quay lại) |
| `d` | Đặt tài khoản đang chọn làm mặc định |
| `r` | Lấy lại mức sử dụng, không dùng bộ nhớ đệm ngắn hạn |
| `q`, `Esc`, `Ctrl+C` | Thoát |

Bộ chọn đóng lại trước khi thực hiện thao tác, nên terminal được khôi phục trước khi Claude Code khởi động. Sau `d` hoặc một thay đổi về dùng chung, bộ chọn mở lại. Thay đổi về dùng chung sẽ hiển thị kế hoạch và hỏi trước, như `ccam share`.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--dry-run` | In ra biến môi trường và lệnh mà ccam sẽ chạy, rồi thoát | tắt |
| `-v`, `--version` | In ra `ccam <version>` (chỉ khi là đối số đầu tiên) | |
| `-h`, `--help` | Trợ giúp của ccam, dưới dạng `ccam --help` hoặc `ccam run --help` | |

Quy tắc về đối số:

- ccam chỉ đọc cờ của chính nó ở trước tên tài khoản. Sau tên, mọi đối số đều thuộc về Claude Code, kể cả `--dry-run` và `--help`.
- `--dry-run` hoạt động dưới dạng `ccam --dry-run work`, `ccam run --dry-run work` và `ccam --dry-run` (tài khoản mặc định). Nó được hỗ trợ, không bị từ chối.
- `ccam run <name|number> [claude args]` là dạng đầy đủ của `ccam <name>`. `ccam run` không kèm tên sẽ dùng tài khoản mặc định.
- Tên tài khoản không được trùng với tên lệnh của ccam như `add`, `list`, `run` hoặc `usage`, nên `ccam <word>` không bao giờ mơ hồ.

Lần chạy thử in ra một dòng mà bạn có thể dán vào shell:

```console
$ ccam --dry-run work -p "fix the tests"
CLAUDE_CONFIG_DIR=/home/me/.local/share/ccam/profiles/work /home/me/.local/bin/claude --mcp-config=/home/me/.local/share/ccam/groups/default/mcp.json -p 'fix the tests'
```

## Cách hoạt động

Mỗi lần khởi động đi qua cùng các bước, từ dòng lệnh đến `exec`:

<!-- motion:launch -->

Mỗi ô sáng lên khi bước của nó chạy, và các dòng bên dưới là biến môi trường và đối số mà Claude Code khởi động cùng.

1. `cmd/ccam/main.go` gọi `app.Main`. Nếu ccam được khởi động dưới một tên khác, đó là một [lệnh tắt](/features/shortcuts) và đi thẳng tới trình khởi chạy.
2. `app.parse` xem `os.Args` trước khi cobra tham gia. Chỉ `ccam` đứng một mình (bộ chọn), một tên lệnh của ccam (`config.Reserved`), `-h`/`--help` và `-v`/`--version` mới đi tới `cli.Execute`. Mọi thứ khác là một lần khởi động.
3. `launcher.Launch` đọc `~/.config/ccam/config.toml`, tệp duy nhất nó phân tích.
4. `launcher.Prepare` tìm tài khoản theo tên hoặc số, bỏ một dấu `--` đứng đầu, và tìm tệp thực thi `claude` bằng `ResolveClaude`. Kết quả tra cứu PATH được lưu đệm trong `~/.local/state/ccam/claude-path` cùng với thời điểm sửa đổi của tệp, nên ccam không phải tìm trong `PATH` ở mỗi lần khởi động. Tệp bị thay đổi hoặc không còn sẽ kích hoạt một lần tra cứu mới. Đường dẫn `claude_bin` trong `config.toml` được dùng nguyên như vậy. ccam từ chối một `claude` chính là ccam, so sánh theo inode, vì điều đó sẽ tạo vòng lặp.
5. `checkLinks` so sánh thư mục tài khoản `~/.local/share/ccam/profiles/<name>/` với thiết lập dùng chung của nó bằng `share.Check`, hàm này chỉ gọi `lstat` và `readlink`, một hoặc hai lời gọi hệ thống cho mỗi mục. Liên kết còn thiếu được tạo ngay lập tức. Bất cứ thứ gì có thể làm mất dữ liệu đều dừng việc khởi động.
6. `applyMCP` thêm các máy chủ MCP dùng chung (xem bên dưới).
7. `session.Register` ghi lại phiên. ccam ghi một tệp đánh dấu đặt tên theo PID của nó, chứa thời điểm bắt đầu tiến trình, vào `$XDG_RUNTIME_DIR/ccam/profile-<name>/<pid>` (hoặc `~/.local/state/ccam/locks/` khi không có `XDG_RUNTIME_DIR`). Với tài khoản dùng chung lịch sử hội thoại, `historyCheck` trước tiên tìm một phiên đang chạy trong `history-<group>-<hash>`, trong đó hash là của thư mục hiện tại, và cũng đăng ký ở đó. `exec` giữ nguyên PID, nên khi đó tệp đánh dấu chỉ đúng tiến trình Claude Code. Dùng tệp đánh dấu chứ không giữ khóa, vì các tiến trình do Claude Code khởi động sẽ thừa hưởng khóa và giữ nó sau khi phiên kết thúc.
8. `Plan.Exec` đặt tiêu đề terminal thành `claude · <name>` (khi đầu ra lỗi chuẩn là terminal), ghi thời điểm khởi động vào `~/.local/state/ccam/last-used/<name>` mà không `fsync`, và gọi `syscall.Exec`.

Claude Code chạy với `CLAUDE_CONFIG_DIR` trỏ tới thư mục tài khoản và `CCAM_PROFILE` là tên tài khoản. ccam xóa mọi giá trị của các biến này mà nó thừa hưởng. Với một thư mục `~/.claude` được đăng ký bằng `ccam import` trên macOS, nó còn đặt `CLAUDE_SECURESTORAGE_CONFIG_DIR` thành giá trị rỗng, để tài khoản và `claude` thông thường dùng cùng một thông tin đăng nhập trong Chuỗi khóa (Keychain).

### Vì sao dùng exec

`syscall.Exec` thay tiến trình ccam bằng `claude`, nên ccam không ở lại trong bộ nhớ, phím bấm, tín hiệu và mã thoát là của chính Claude Code, và đối số luôn là một mảng, không bao giờ là chuỗi shell. Công việc sau một phiên, như phát hiện các máy chủ được thêm bằng `claude mcp add`, diễn ra ở lần khởi động tiếp theo: ccam so sánh thời điểm sửa đổi của `.claude.json` với `~/.local/state/ccam/mcp/<name>/seen` và chỉ phân tích tệp khi nó đã thay đổi.

### Cách một dòng lệnh được định tuyến

Ba dòng lệnh, ba đường đi qua `internal/app`:

<!-- motion:launch-route -->

Lệnh tắt được nhận ra qua `argv[0]`, `ccam <name>` qua đối số đầu tiên, và cả hai đều kết thúc ở `syscall.Exec`. Chỉ một lệnh của ccam như `ccam usage` mới đến cobra.

### Máy chủ MCP dùng chung

Với chiến lược mặc định `flag`, ccam truyền các máy chủ MCP của nhóm cho Claude Code dưới dạng một đối số, đặt trước các đối số của bạn:

<!-- motion:launch-mcp -->

`mcp.FlagArgs` chỉ trả về một đối số nếu `mcp.json` của nhóm tồn tại, và `Plan.InsertArgs` đặt nó ngay sau `claude`. Khi đó `--mcp-config` hoặc `--strict-mcp-config` của riêng bạn vẫn được áp dụng chồng lên. Đó là một đối số duy nhất, `--mcp-config=<path>`, vì `--mcp-config` của Claude Code nhận nhiều giá trị: nếu viết thành hai đối số, nó sẽ đọc cả lời nhắc theo sau như một tệp cấu hình.

Với `mcp_strategy = "sync"` trong `config.toml`, ccam ghi các máy chủ vào `.claude.json` của tài khoản thay vào đó, và bỏ qua việc này nếu tài khoản đang chạy.

### Remote Control

`ccam work rc` (hoặc `ccam work remote-control`) khởi động Remote Control của Claude Code bằng tài khoản `work`. `claude remote-control` từ chối `--mcp-config` đặt trước động từ của nó và không nhận cờ này sau động từ, và các phiên của nó chỉ đọc máy chủ MCP từ `.claude.json` của tài khoản. Vì vậy chiến lược flag hoạt động khác ở đây:

<!-- motion:launch-rc -->

1. `remoteControl` coi đối số đầu tiên sau tên tài khoản không bắt đầu bằng `-` là động từ của Claude Code. Khi đó là `remote-control` hoặc `rc`, `mcp.FlagArgs` không thêm gì.
2. `mcp.Synced` so sánh thời điểm sửa đổi của `mcp.json` của nhóm với `sync-stamp` của tài khoản. Nếu `.claude.json` đã chứa các máy chủ hiện tại của nhóm, hoặc nhóm không có máy chủ nào, Remote Control khởi động mà không hỏi.
3. Nếu không, ccam nêu tên các máy chủ mà các phiên sẽ thiếu và hỏi:

   ```console
   ccam: Remote Control sessions only get MCP servers stored in account work's .claude.json, not the shared MCP servers (ctx, fs).
   Copy them into .claude.json (a backup is kept)? [y/N]
   ```

4. Khi chọn `y`, `mcp.Sync` sao chép chúng theo cách của chiến lược `sync`: dưới khóa của Claude Code, sau khi sao lưu, và giữ lại các máy chủ riêng của tài khoản. Khi chọn `N`, hoặc khi không có terminal, Remote Control khởi động mà không có chúng, và ccam gợi ý `mcp_strategy = "sync"`, cách này sao chép chúng trước mỗi lần khởi động.

ccam hỏi ở mỗi lần khởi động Remote Control cho đến khi các máy chủ được sao chép, và sau khi sao chép thì chỉ hỏi một lần khi `mcp.json` của nhóm lại thay đổi. Bản sao trong `.claude.json` được ưu tiên hơn `--mcp-config` ở các lần khởi động bình thường sau đó, nên `ccam mcp add`, `remove` và `import` cập nhật các bản sao mà tài khoản khác đang giữ ([Máy chủ MCP dùng chung](/features/mcp)).

### Ngân sách 10 ms

ccam thêm tối đa 10 ms (p95) trước khi Claude Code khởi động. Đường khởi động không bao giờ dựng cây lệnh cobra, khởi động Bubble Tea hay mở kết nối mạng, và `TestPackageInitIsCheap` thất bại nếu các hàm `init()` của gói mất tổng cộng hơn 5 ms. `scripts/perf-check.sh` đo `ccam --dry-run work`, `ccam work` và một lệnh tắt bằng `hyperfine` và thất bại khi vượt 10 ms, còn `BenchmarkLaunchPath` trong `internal/launcher` chạy các bước trong một vòng lặp. Để xem thời gian của từng bước, hãy đặt `CCAM_TRACE=1`:

```console
$ CCAM_TRACE=1 ccam --dry-run work >/dev/null
ccam trace: start                 2.830ms (total 2.830ms)
ccam trace: load config           0.397ms (total 3.227ms)
ccam trace: prepare               0.112ms (total 3.339ms)
ccam trace: check links           0.238ms (total 3.578ms)
ccam trace: prelaunch             0.083ms (total 3.661ms)
```

## An toàn và các trường hợp đặc biệt

- Một tên hoặc số không xác định sẽ dừng với `no account "<name>"` và chỉ tới `ccam list`. Khi chưa có tài khoản nào, ccam yêu cầu bạn thêm một tài khoản bằng `ccam add <name>`.
- Nếu một đường dẫn dùng chung đã trở thành tệp thông thường hoặc trỏ sai chỗ, ccam dừng với `stopped because of share link problems:`, liệt kê chúng, và yêu cầu bạn chạy `ccam doctor --fix`. Nó không di chuyển dữ liệu trên đường khởi động.
- Nếu một Claude Code khác dùng chung lịch sử hội thoại với tài khoản này (cùng nhóm) đang chạy trong cùng thư mục, ccam cảnh báo rằng lịch sử có thể bị trộn lẫn và hỏi `Run anyway? [y/N]`. Khi không có terminal, nó hiển thị cảnh báo và tiếp tục.
- `--dry-run` bỏ qua việc đồng bộ MCP, thông báo máy chủ mới, kiểm tra lịch sử, tệp đánh dấu phiên, tiêu đề và thời điểm dùng gần nhất, và không exec. Nó vẫn kiểm tra liên kết dùng chung, và tạo liên kết còn thiếu như một lần khởi động bình thường.
- `ccam work rc` khi `work` đang chạy sẽ không động đến `.claude.json`, báo điều đó, và khởi động Remote Control mà không có các máy chủ của nhóm; hãy chạy lại sau khi phiên kết thúc để sao chép chúng. Với `--dry-run`, nó in thông báo nhưng không hỏi và không sao chép.
- Với `mcp_strategy = "sync"`, việc đồng bộ thất bại sẽ in ra `ccam: MCP sync failed (launching anyway): …` và Claude Code khởi động với các máy chủ nó đã có.

## Liên quan

- [Lệnh tắt](/features/shortcuts)
- [Dùng Remote Control với ccam](/guides/remote-control)
- [Quản lý tài khoản](/features/accounts)
- [Tham chiếu lệnh](/docs/commands)
- [Dùng chung máy chủ MCP, cài đặt và lịch sử](/docs/sharing)
- [Cách CLAUDE_CONFIG_DIR hoạt động](/guides/claude-config-dir)
