---
title: Giới hạn sử dụng
description: "ccam usage hiển thị giới hạn 5 giờ và hàng tuần của mọi tài khoản cùng lúc, và ccam statusline đưa tên tài khoản cùng mức sử dụng vào Claude Code."
updated: 2026-10-21
order: 7
translation: machine
sourceHash: "b26800c6fa54"
---

`ccam usage` cho biết mỗi tài khoản đã dùng bao nhiêu trong giới hạn 5 giờ, giới hạn hàng tuần và giới hạn hàng tuần riêng của từng mô hình (như Fable), với một thanh cho mỗi giới hạn cùng thời gian đến lần đặt lại. Hãy dùng nó trước khi bắt đầu một tác vụ dài, hoặc khi Claude Code báo rằng sắp chạm giới hạn, để xem tài khoản nào còn hạn mức. `ccam statusline` là lệnh dòng trạng thái cho Claude Code: nó hiển thị tên tài khoản và mức sử dụng bên trong Claude Code và lưu các con số cho `ccam usage`.

## Cách dùng

```sh
ccam usage                  # every account, once
ccam u                      # the same, shorter
ccam usage --watch          # redraw every 60 seconds until Ctrl+C
ccam usage --watch=10       # redraw every 10 seconds
ccam usage -v               # add where the numbers came from
ccam usage --json           # machine-readable, one object per account
ccam usage --source cache   # saved numbers only, no network
ccam usage --refresh        # also renew expired logins of accounts not running
ccam usage --source claude  # let Claude Code look it up: ccam never touches a token
ccam statusline install     # set the status line in every account's settings.json
```

```console
$ ccam usage
 1 personal  9h 36m ago (cached)
   5-hour   ░░░░░░░░░░░░░░░░░░░░   0%
   Weekly   ██████████████░░░░░░  70%  resets in 2d 16h
   Fable    ██░░░░░░░░░░░░░░░░░░  12%  resets in 2d 16h
 2 work      just now (cached)
   5-hour   █████░░░░░░░░░░░░░░░  24%  resets in 2h 29m
   Weekly   ████████░░░░░░░░░░░░  41%  resets in 3d 11h
   Fable    ███████████░░░░░░░░░  55%  resets in 5d 2h
 3 side      login expired, run `ccam side` once to refresh
→ Most headroom: personal   (open it with ccam next --auto)
```

Mỗi tài khoản có một dòng gồm số thứ tự, tên và độ cũ của các con số; `(cached)` nghĩa là chúng lấy từ một tệp đã lưu, không phải từ endpoint. Bên dưới, mỗi giới hạn có một dòng với thanh 20 ô, tỷ lệ phần trăm và thời gian đến khi khoảng đó được đặt lại: `5-hour`, `Weekly` (mọi mô hình), rồi đến khoảng hàng tuần riêng của từng mô hình, như `Fable`. Khoảng mà tài khoản không có thì hiển thị `-`. Tỷ lệ phần trăm chuyển sang màu vàng từ 80% và màu đỏ ở 100%. Trên terminal, mọi tài khoản xuất hiện ngay cùng lúc dưới dạng `fetching…`, và các dòng của từng tài khoản được vẽ khi nó phản hồi. Dòng cuối nêu tên tài khoản mà [`ccam next --auto`](/features/next) sẽ chọn theo cùng các quy tắc (`auto_max_usage`, `auto_strategy`, các phiên đang chạy, chỉ các tài khoản thuộc tổ chức của tài khoản dùng gần nhất trừ khi đặt `cross_org`, bỏ qua mọi khoảng ở mức 100%, nên một tài khoản đã đầy hạn mức tuần không bao giờ được nêu tên), và chỉ xuất hiện khi bạn có nhiều hơn một tài khoản. Tài khoản chưa đăng nhập hiển thị `not logged in → ccam login <name>`, còn tài khoản dùng khóa API hiển thị `no limits (API)`.

Sau `ccam statusline install`, Claude Code hiển thị một dòng như `[work] · Opus · 5h 24% · 7d 41%`.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--json` | In các báo cáo dưới dạng JSON: `name`, `slot`, `status`, `limits`, `source`, `at`, `error` | tắt |
| `--watch[=N]` | Vẽ lại mỗi `N` giây cho đến khi bạn nhấn Ctrl+C. `--watch` không kèm giá trị, hoặc giá trị từ 0 trở xuống, nghĩa là 60 | tắt |
| `-v`, `--verbose` | Thêm một dòng `source:` dưới mỗi tài khoản: nơi các con số của nó được lấy (endpoint mức sử dụng, Claude Code, dòng trạng thái hoặc một tệp đã lưu) | tắt |
| `--no-cache` | Bỏ qua bộ nhớ đệm 60 giây và hỏi lại endpoint | tắt |
| `--source live\|claude\|cache` | `live` hỏi endpoint bằng token của tài khoản; `claude` để Claude Code tra cứu (bên dưới); `cache` chỉ hiển thị các con số đã lưu | `usage_source`, nếu không thì `live` |
| `--refresh` | Gia hạn các access token đã hết hạn, chỉ cho những tài khoản không đang chạy | tắt |
| `--dry-run` | Được chấp nhận: `ccam usage` không thay đổi gì ngoài bộ nhớ đệm của chính nó. Bị từ chối khi đi cùng `--refresh`: `--dry-run cannot be combined with --refresh, which rewrites login tokens` | tắt |

Mọi giá trị `--source` khác đều thất bại với `--source must be live, claude or cache`.

### Mức sử dụng từ chính Claude Code

Với `usage_source = "claude"` trong `config.toml` (hoặc `--source claude`), ccam hoàn toàn không đọc token đăng nhập. Với mỗi tài khoản, nó khởi động `claude -p --input-format stream-json --output-format stream-json --verbose --safe-mode --no-session-persistence` trong thư mục của tài khoản và gửi yêu cầu `get_usage` của Claude Code qua đầu vào chuẩn (Claude Code 2.1.270 trở lên); Claude Code tra cứu mức sử dụng bằng thông tin đăng nhập của chính nó và trả lời với khoảng 5 giờ, khoảng hàng tuần và khoảng theo mô hình. `--safe-mode` ngăn các hook, máy chủ MCP và plugin của tài khoản khởi động, và không có phiên nào được lưu. Mỗi tài khoản cần một lần khởi động Claude Code, khoảng một giây, nên nguồn này chậm hơn `live`; các tài khoản được hỏi cùng lúc, mỗi tài khoản trong vòng 20 giây, và bộ nhớ đệm 60 giây vẫn áp dụng như thường. Claude Code đánh dấu yêu cầu này là thử nghiệm, nên nếu nó ngừng trả lời, ccam hiển thị các con số đã lưu và `ccam doctor` báo điều đó. Một lần `--dry-run` không khởi động Claude Code (vì việc đó ghi vào thư mục của tài khoản) và hiển thị các con số đã lưu. `usage_source` nhận `live` hoặc `claude`; `cache` chỉ là một lựa chọn của `--source`. `ccam usage` không nhận đối số.

`ccam statusline` đọc JSON dòng trạng thái của Claude Code từ đầu vào chuẩn. Tùy chọn duy nhất của nó là `--wrap <command>`: một lệnh dòng trạng thái có sẵn, dòng đầu ra đầu tiên của lệnh đó sẽ được thêm tên tài khoản vào phía trước. `ccam statusline install` không nhận tùy chọn. Cả hai đều từ chối `--dry-run` với `` --dry-run is not supported by `ccam statusline`; nothing was changed ``.

## Cách hoạt động

Các con số đến từ hai nơi: endpoint mức sử dụng không được công bố tài liệu của Anthropic, `GET https://api.anthropic.com/api/oauth/usage`, được gọi bằng thông tin đăng nhập riêng của từng tài khoản, và các giá trị mà Claude Code trao cho lệnh dòng trạng thái của nó. Toàn bộ nằm trong `internal/usage`.

1. `cli.usageCmd` tải `config.toml` và tạo một `usage.Account` cho mỗi tài khoản: thư mục `~/.local/share/ccam/profiles/<name>/` của nó, thư mục phiên của nó (để xem nó có đang chạy không) và việc nó có được nhập tại chỗ hay không.
2. `Fetcher.FetchAll` khởi động một goroutine cho mỗi tài khoản. Tất cả dùng chung một `http.Client`, nên các kết nối tới endpoint được dùng lại. Việc kết nối có thể mất 2 giây, và mỗi tài khoản có tổng cộng 3 giây. Client không đi theo chuyển hướng, vì một lần chuyển hướng sẽ gửi token tới nơi khác, và chỉ đọc tối đa 1 MB của phản hồi.
3. Khi từng tài khoản hoàn tất, hàng của nó được vẽ lại. Tài khoản thất bại hoặc không phản hồi kịp sẽ quay về các con số đã lưu, nên một tài khoản chậm không bao giờ giữ chân các tài khoản khác.

Mỗi tài khoản chạy trong goroutine riêng để gọi endpoint; tài khoản bị quá thời gian sẽ hiển thị các con số đã lưu thay vào đó.

<!-- motion:usage-fetch -->

`personal` và `work` phản hồi kịp và hiển thị số liệu trực tiếp. `side` mất hơn 3 giây, nên `fallback()` đọc hai tệp bộ nhớ đệm của nó và dùng tệp nào được ghi sau cùng. `usage.Table()` vẽ một khối cho mỗi tài khoản và `usage.Best()` chọn tài khoản được nêu tên ở dòng cuối.

Với một tài khoản, `Fetcher.Fetch` thực hiện các kiểm tra sau theo thứ tự:

<!-- motion:usage-cache -->

1. Với `--source cache`, nó đi thẳng tới `fallback()`.
2. Trừ khi đặt `--no-cache`, một kết quả trực tiếp trong `~/.local/state/ccam/usage/<name>.json` mới hơn 60 giây (`CacheTTL`) được trả về nguyên như vậy. Chạy `ccam usage` hai lần liên tiếp chỉ tốn một yêu cầu cho mỗi tài khoản.
3. `readCredentials()` đọc `.credentials.json` (hoặc Chuỗi khóa (Keychain) của macOS) mà không thay đổi nó, và chỉ khóa `claudeAiOauth` được giải mã. Không có thông tin đăng nhập gói đăng ký nghĩa là `not logged in`, hoặc `no limits (API)` khi `.claude.json` chứa một khóa API.
4. Access token đã hết hạn không được dùng. Không có `--refresh`, tài khoản hiển thị các con số đã lưu và được đánh dấu là đã hết hạn.
5. Nếu không, ccam gọi endpoint. Mã 401 hoặc 403 được tính là đăng nhập hết hạn, mọi trạng thái khác hoặc nội dung không đọc được được tính là lỗi, và cả hai đều quay về các con số đã lưu. Mã 429 nghĩa là endpoint đang bị hỏi quá thường xuyên: ccam ghi `<name>.pause` với thời điểm được hỏi lại (`Retry-After`, hoặc 5 phút, tối đa một giờ), hiển thị các con số đã lưu kèm `rate limited by the usage endpoint; asking again after 15:10`, và không gọi endpoint cho tài khoản đó cho đến lúc ấy, kể cả với `--no-cache`. Thử lại ngay chỉ khiến giới hạn kéo dài thêm.
6. Một phản hồi hợp lệ được ghi vào `<name>.json` cùng thời điểm đo.

`fallback()` đọc cả `<name>.json` lẫn `<name>.statusline.json` và giữ tệp mới hơn. Trước khi hiển thị các con số đã lưu, nó kiểm tra thời điểm đặt lại của từng khoảng: một khoảng đã được đặt lại kể từ khi các con số được lưu sẽ hiển thị 0%, vì tỷ lệ phần trăm cũ không còn đúng. Đó là lý do `personal` ở trên hiển thị `0%` mà không có thời gian đặt lại.

Endpoint có hai định dạng phản hồi. Định dạng mới hơn có một mảng `limits` với các mục `session` (5 giờ), `weekly_all` và `weekly_scoped` (theo từng mô hình). Định dạng cũ hơn có các đối tượng `five_hour`, `seven_day`, `seven_day_opus` và `seven_day_sonnet`. `ParseLive` đọc cả hai: mảng `limits` được ưu tiên, và các đối tượng cũ bổ sung những gì mảng đó không có. Một nội dung hoàn toàn không có giới hạn nào được báo là định dạng đã thay đổi, và `ccam doctor` sẽ báo cho bạn điều đó.

### Gia hạn token

Refresh token được xoay vòng: gia hạn một token sẽ cho bạn một refresh token mới và vô hiệu hóa token cũ. Nếu ccam gia hạn thông tin đăng nhập trong lúc Claude Code đang chạy bằng cùng tài khoản đó, phiên đó sẽ bị đăng xuất. Vì vậy `ccam usage` đọc token và không bao giờ ghi chúng, trừ khi bạn truyền `--refresh`. Khi đó `Fetcher.refresh`:

1. từ chối các tài khoản được đăng ký tại chỗ bằng `ccam import`, vì `claude` thông thường có thể dùng thông tin đăng nhập đó mà ccam không biết;
2. từ chối thông tin đăng nhập được lưu trong Chuỗi khóa (Keychain) của macOS;
3. lấy một tệp khóa, `$XDG_RUNTIME_DIR/ccam/profile-<name>.refresh.lock` (trong `~/.local/state/ccam/locks/` khi biến đó không được đặt), để hai tiến trình ccam không bao giờ gia hạn cùng một token hai lần;
4. trong khi giữ khóa, từ chối nếu có một phiên của tài khoản đang chạy, và bỏ qua yêu cầu nếu một tiến trình khác đã gia hạn token;
5. gửi refresh token tới endpoint token và ghi các token mới trở lại `.credentials.json` bằng `fsutil.AtomicWrite`, giữ nguyên mọi khóa khác trong tệp.

### Bộ thu thập từ dòng trạng thái

Claude Code chạy lệnh dòng trạng thái và pipe cho nó một tài liệu JSON. ccam đọc tên mô hình và `rate_limits` từ đó:

<!-- motion:usage-statusline -->

1. `runStatusline` tìm tài khoản từ `CCAM_PROFILE`, biến mà ccam đặt khi khởi động Claude Code, hoặc nếu không có thì bằng cách so khớp `CLAUDE_CONFIG_DIR` với các thư mục tài khoản.
2. `ParseStatusLine` đọc `rate_limits.five_hour` và `rate_limits.seven_day`. JSON không có giới hạn tần suất đối với tài khoản dùng khóa API, hoặc trước phản hồi đầu tiên của phiên; khi đó không có gì được lưu.
3. `SaveStatusLine` ghi `~/.local/state/ccam/usage/<name>.statusline.json`. `ccam usage` dùng tệp này khi endpoint không phản hồi.
4. ccam in ra `[name] · model · 5h N% · 7d N%`. Với `--wrap`, nó chạy lệnh của bạn qua `sh -c` với cùng đầu vào và in ra `[name]` theo sau là dòng đầu tiên của lệnh thay vào đó. Trên Windows, nó chạy lệnh theo cách Claude Code làm: bằng Git Bash khi đã được cài, nếu không thì bằng PowerShell.
5. Khi mức sử dụng 5 giờ hoặc hàng tuần ở mức 90% trở lên và có một tài khoản khác, dòng này kết thúc bằng `→ ccam next --auto --continue`, ví dụ `[work] · Opus · 5h 93% · 7d 41% · → ccam next --auto --continue`.

Lệnh được bao bọc có 2 giây. Sau đó nó bị dừng và dòng chỉ hiển thị tên tài khoản, nên một lệnh chậm không thể chặn dòng trạng thái của Claude Code.

`ccam statusline install` ghi `statusLine` vào `settings.json` của từng tài khoản. Khi một tài khoản đã có lệnh dòng trạng thái, ccam giữ nó dưới dạng `--wrap '<old command>'`. Trên Windows, đường dẫn tới ccam được ghi bằng dấu gạch chéo xuôi và `~` cho thư mục home (`~/.local/bin/ccam.exe statusline`), cách viết mà cả Git Bash lẫn PowerShell đều hiểu. Một `settings.json` được dùng chung qua nhóm chỉ được ghi một lần, qua liên kết của nó, nên mọi tài khoản trong nhóm đều nhận được. Chạy lại lệnh sẽ in ra `(already set)` cho các tệp đã dùng ccam.

## An toàn và các trường hợp đặc biệt

- `ccam usage` không bao giờ thay đổi tệp token trừ khi bạn truyền `--refresh`. Khi một thông tin đăng nhập đã hết hạn, hãy khởi động tài khoản một lần bằng `ccam <name>` và Claude Code sẽ gia hạn nó.
- `--refresh` không bao giờ động đến tài khoản đang chạy, tài khoản đã nhập hay thông tin đăng nhập trong Chuỗi khóa (Keychain). Một lần gia hạn bị từ chối sẽ hiển thị các con số đã lưu và lý do.
- Endpoint mức sử dụng không được công bố tài liệu chính thức và có thể thay đổi. Khi nó thay đổi, bạn vẫn thấy các con số đã lưu, và `ccam doctor` báo việc định dạng đã thay đổi.
- Khi không có con số đã lưu và không có phản hồi, dòng của tài khoản đó hiển thị `lookup failed:` kèm lý do, ví dụ `no saved usage` với `--source cache`.
- Các tệp bộ nhớ đệm chỉ chứa tỷ lệ phần trăm và thời điểm đặt lại, không bao giờ chứa token.

## Liên quan

- [Giới hạn sử dụng và chuyển đổi tài khoản](/docs/usage-limits)
- [Chuyển sang tài khoản tiếp theo](/features/next)
- [Quy trình khi chạm giới hạn sử dụng](/guides/usage-limit-workflow)
- [Lệnh](/docs/commands)
