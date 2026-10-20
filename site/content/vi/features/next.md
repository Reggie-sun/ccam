---
title: Chuyển sang tài khoản tiếp theo
description: "ccam next mở Claude Code bằng tài khoản kế tiếp; --auto chọn tài khoản còn nhiều hạn mức nhất, --continue tiếp tục hội thoại, --fallback chạy lại tác vụ -p."
updated: 2026-10-20
order: 8
translation: machine
sourceHash: 784ae4d1f6d9
---

`ccam next` mở Claude Code bằng một tài khoản khác mà bạn không cần nêu tên. Hãy dùng nó khi một tài khoản chạm giới hạn sử dụng: `ccam next` không kèm tùy chọn chọn tài khoản đứng sau tài khoản bạn dùng gần nhất, `--auto` chọn tài khoản còn nhiều hạn mức nhất trong khoảng 5 giờ, và `--continue` tiếp tục cuộc hội thoại gần nhất trong tài khoản mới. Với các script chạy `claude -p`, `ccam --fallback` tự thực hiện việc đổi đó (xem [Chạy lại khi chạm giới hạn](#run-again-at-the-limit)).

## Cách dùng

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -- -p "summarize"     # pass arguments to claude after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next
→ personal → work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude
```

ccam in việc chuyển đổi dưới dạng `→ <from> → <to>` ra đầu ra lỗi chuẩn, rồi khởi động Claude Code bằng tài khoản mới theo cùng cách `ccam <name>` làm. `<from>` là tài khoản bạn dùng gần nhất. Các đối số sau `--` được chuyển nguyên vẹn cho `claude`.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--auto` | Chọn tài khoản đã đăng nhập theo mức sử dụng (xem bên dưới), bỏ qua các tài khoản ở mức 100% | tắt |
| `--max-usage <1-100>` | Với `--auto`: ưu tiên các tài khoản có mức sử dụng dưới giá trị này | `auto_max_usage`, nếu không thì 100 |
| `--strategy headroom\|expiring` | Với `--auto`: mức sử dụng 5 giờ thấp nhất, hoặc khoảng 5 giờ được đặt lại sớm nhất | `auto_strategy`, nếu không thì `headroom` |
| `--wait[=<duration>]` | Với `--auto`: khi không tài khoản nào còn hạn mức, chờ tài khoản đầu tiên được đặt lại, tối đa chừng này thời gian | tắt; `5h` khi chỉ ghi cờ |
| `-c`, `--continue` | Chỉ xét các tài khoản dùng chung lịch sử hội thoại với tài khoản bạn dùng gần nhất, và khởi động `claude --continue` | tắt |
| `-- <claude args>` | Các đối số được chuyển cho `claude` sau các đối số của chính ccam | không có |
| `--dry-run` | Được hỗ trợ: in ra `→ <from> → <to>` và dòng lệnh sẽ chạy, mà không khởi động Claude Code hay ghi lại lần khởi động | tắt |

`--auto` và `--continue` dùng được cùng nhau. Với `--continue`, `--continue` được đặt trước các đối số của bạn.

## Cách hoạt động

Mỗi lần khởi động qua ccam đều ghi thời điểm vào `~/.local/state/ccam/last-used/<name>`. `ccam next` đọc các tệp này, kiểm tra những tài khoản nào đã đăng nhập, rồi chọn.

1. `cli.candidates` tạo một `next.Candidate` cho mỗi tài khoản từ `config.toml`: tài khoản, thời điểm khởi động gần nhất của nó (`launcher.LastUsed`) và việc nó có thông tin đăng nhập hay không (`creds.Has`: một `.credentials.json` đọc được, hoặc trên macOS là một mục trong Chuỗi khóa (Keychain)).
2. Không có `--auto`, `next.Rotate` sắp xếp các tài khoản theo số slot và tìm tài khoản được khởi động gần nhất. Đó là tài khoản bạn đang chuyển đi. Nếu chưa có tài khoản nào được khởi động, đó là slot thấp nhất.
3. Từ đó nó đi tiếp theo thứ tự slot, quay vòng từ slot cuối về slot đầu, và trả về tài khoản đầu tiên đã đăng nhập.
4. `launcher.LaunchWith` chuẩn bị và khởi động tài khoản được chọn đúng như `ccam <name>`: nó kiểm tra các liên kết của tài khoản, áp dụng việc dùng chung MCP, đăng ký phiên và thay ccam bằng `claude` qua `syscall.Exec`.

Vòng duyệt bắt đầu ở tài khoản dùng gần nhất và dừng ở tài khoản đã đăng nhập đầu tiên sau nó:

<!-- motion:next-rotate -->

`work` được khởi động gần nhất, nên vòng duyệt bắt đầu sau nó. `side` chưa đăng nhập nên bị bỏ qua, và `lab` được chọn. Nếu `lab` cũng chưa đăng nhập, vòng duyệt sẽ quay về `personal` (`% len` trong sơ đồ). Vòng xoay không bao giờ trả về tài khoản bạn bắt đầu: khi không có tài khoản nào khác đã đăng nhập, nó dừng lại với một lỗi.

### --auto

`--auto` quyết định theo mức sử dụng thay vì thứ tự slot. `cli.pickAuto` giữ lại các tài khoản đã đăng nhập, hỏi mức sử dụng của tất cả cùng lúc bằng chính `Fetcher.FetchAll` mà `ccam usage` dùng, và để `usage.Policy.Choose()` chọn:

<!-- motion:next-auto -->

Một tài khoản được đánh giá theo khoảng 5 giờ và khoảng hàng tuần của nó, và khi `--model` nằm trong các đối số cho Claude Code, theo cả khoảng hàng tuần của mô hình đó (`opus`, `sonnet`, `haiku` hoặc `fable` trong tên). `Policy.Choose()` chia các tài khoản thành ba nhóm:

1. **Còn hạn mức**: mọi khoảng đều dưới `auto_max_usage` (`config.toml`, mặc định 100). Chiến lược chọn trong số này: `headroom` (mặc định) lấy mức sử dụng 5 giờ thấp nhất, `expiring` lấy khoảng 5 giờ được đặt lại sớm nhất, để hạn mức còn lại được dùng trước khi mất.
2. **Vượt ngưỡng**: dưới 100% nhưng bằng hoặc vượt `auto_max_usage`. Chỉ được dùng khi không tài khoản nào còn hạn mức, ưu tiên mức sử dụng 5 giờ thấp nhất.
3. **Đầy**: có một khoảng ở mức 100%, hoặc chưa rõ mức sử dụng 5 giờ. Không bao giờ được chọn.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

Khi không chọn được tài khoản nào, ccam thoát với mã 75 (`EX_TEMPFAIL`), để script phân biệt được "thử lại sau" với một lỗi. Với `--wait`, thay vào đó nó ngủ cho đến khi khoảng đang chặn của tài khoản đầy đầu tiên được đặt lại (30 giây sau thời điểm đặt lại, khi endpoint hiển thị khoảng mới), in ra `ccam: no account has room; waiting until 15:10`, rồi chọn lại; khi không tài khoản nào cho biết lúc nào được đặt lại, nó kiểm tra lại sau 10 phút. Nếu không có tài khoản nào trống ra trước khi hết thời gian chờ, nó thoát với mã 75. Việc tra cứu tuân theo các quy tắc của `ccam usage`: các con số mới hơn 60 giây được dùng lại, mỗi tài khoản có 3 giây, và tài khoản không phản hồi hoặc có đăng nhập đã hết hạn sẽ được đánh giá theo các con số đã lưu. Một khoảng đã lưu mà thời điểm đặt lại đã qua được tính là 0%, nên tài khoản có khoảng giới hạn đã được đặt lại kể từ đó sẽ không bị bỏ qua. `--auto` không bao giờ gia hạn token. Khác với vòng xoay thông thường, `--auto` có thể chọn tài khoản bạn dùng gần nhất nếu nó vẫn còn nhiều hạn mức nhất.

### --continue

`--continue` của Claude Code mở cuộc hội thoại gần nhất trong thư mục hiện tại, và điều này chỉ hoạt động nếu tài khoản mới nhìn thấy được cuộc hội thoại đó. Vì vậy với `--continue`, ccam chỉ xét các tài khoản trong cùng nhóm `history` với tài khoản bạn dùng gần nhất, rồi đặt `--continue` trước các đối số của bạn. Các tài khoản có lịch sử là `local` không dùng chung với ai, nên `--continue` từ một tài khoản như vậy luôn thất bại. Hãy bật dùng chung trước bằng `ccam share <name> --history on` (xem [dùng chung](/docs/sharing)).

## Chạy lại khi chạm giới hạn

`ccam next` dành cho một phiên tương tác đã dừng. Một lần chạy không tương tác, `claude -p` trong script hoặc tác vụ định kỳ, không thể chờ bạn, nên `--fallback` đặt trước tên tài khoản cho phép ccam đổi tài khoản thay cho nó:

```sh
ccam --fallback work -p "fix the failing tests" --output-format json
ccam --fallback -p "summarise" < notes.txt      # the folder's account or the default
ccam run --fallback work -p "review the diff"
```

```console
$ ccam --fallback personal -p "summarise" < notes.txt
ccam: personal reached its usage limit; running the command again as work
Summary: …
```

1. `--fallback` chỉ hoạt động cùng với `-p` hoặc `--print`; nếu không, ccam dừng lại và chỉ tới `ccam next --auto --continue`.
2. Đầu vào chuẩn được đọc một lần và đưa cho mọi lần thử.
3. ccam khởi động `claude` như một tiến trình con thay vì `exec`, giữ lại đầu ra chuẩn của nó, và chuyển tiếp đầu ra lỗi chuẩn ngay khi có.
4. Khi `claude` thoát với 0, hoặc thất bại vì lý do khác, đầu ra đã giữ và mã thoát được chuyển tiếp nguyên vẹn.
5. Một lần thất bại được tính là do giới hạn sử dụng khi thông báo của Claude Code nói vậy (`You've hit your limit`, `usage limit reached`), hoặc khi tra cứu mức sử dụng của tài khoản cho thấy khoảng 5 giờ hoặc hàng tuần ở mức 100%.
6. Khi đó đầu ra của lần thử đó bị bỏ đi và cùng lệnh chạy bằng một tài khoản đã đăng nhập khác: tài khoản có mức sử dụng 5 giờ thấp nhất, rồi đến các tài khoản không tra cứu được mức sử dụng, theo thứ tự số. Tài khoản ở mức 100% bị bỏ qua, và mỗi tài khoản được thử nhiều nhất một lần.
7. Khi có `-c`, `--continue`, `-r` hoặc `--resume` trong các đối số, chỉ các tài khoản trong cùng nhóm lịch sử được thử, để có thể tìm thấy cuộc hội thoại.

Tài khoản tiếp theo được chọn theo cùng các quy tắc như `--auto` (`auto_max_usage`, `auto_strategy`, khoảng của `--model`). Khi không còn tài khoản nào, ccam in ra đầu ra của lần thử cuối, báo rằng không còn tài khoản đã đăng nhập nào còn hạn mức, và thoát với mã 75. Với `--wait` đặt trước tên tài khoản (`ccam --fallback --wait=2h work -p …`), thay vào đó nó chờ tài khoản đầu tiên được đặt lại, rồi có thể thử lại mọi tài khoản; nó luôn ngủ ít nhất một lần trước, nên tài khoản vừa chạm giới hạn sẽ không bị chạy lại ngay. `ccam --dry-run --fallback work -p …` in ra lệnh của lần thử đầu tiên.

Các phiên tương tác không tự đổi tài khoản. Thay vào đó, với `ccam statusline install`, dòng trạng thái thêm `→ ccam next --auto --continue` khi mức sử dụng 5 giờ hoặc hàng tuần chạm 90% ([Giới hạn sử dụng](/features/usage)).

## An toàn và các trường hợp đặc biệt

- Tài khoản chưa đăng nhập không bao giờ được chọn. Khi không có tài khoản nào khác đã đăng nhập, `ccam next` thông thường dừng lại với `no logged-in account other than <name>`.
- Với `--continue` và không có tài khoản nào dùng chung lịch sử: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- Với `--auto` và không có tài khoản nào đã đăng nhập: `no logged-in account to pick`. Khi mọi tài khoản đều đã đầy hoặc chưa rõ mức sử dụng: `` no account with known usage and headroom (see `ccam usage`) ``. Cả hai đều kết thúc với mã thoát 75.
- `--max-usage`, `--strategy` và `--wait` chỉ dùng được với `--auto`, và giá trị ngoài khoảng 1-100 hoặc một chiến lược khác sẽ bị từ chối.
- `--auto` gửi yêu cầu mạng ngay cả với `--dry-run`. Nó chỉ đọc, và không bao giờ ghi vào tệp token.
- Tài khoản được chọn khởi động qua cùng đường đi với `ccam <name>`, nên các kiểm tra tương tự được áp dụng, bao gồm cảnh báo khi một Claude Code khác dùng chung lịch sử đang chạy trong thư mục này.
- `last-used` chỉ được ghi khi Claude Code thực sự khởi động. `--dry-run` không động đến nó, nên chạy thử không đẩy vòng xoay tiến lên.
- `--fallback` giữ đầu ra chuẩn cho đến khi một lần thử kết thúc, nên đầu ra dạng luồng như `--output-format stream-json` đến một lần vào cuối. Các tệp mà lần thử đầu tiên đã thay đổi trước khi chạm giới hạn vẫn giữ nguyên thay đổi; lần thử tiếp theo bắt đầu lại lời nhắc từ đầu.

## Liên quan

- [Giới hạn sử dụng](/features/usage)
- [Giới hạn sử dụng và chuyển đổi tài khoản](/docs/usage-limits)
- [Dùng chung máy chủ MCP, cài đặt và lịch sử](/docs/sharing)
- [Quy trình khi chạm giới hạn sử dụng](/guides/usage-limit-workflow)
