---
title: Chuyển sang tài khoản tiếp theo
description: "ccam next mở Claude Code bằng tài khoản khác; --auto chọn tài khoản còn nhiều hạn mức nhất, --continue chuyển hội thoại sang, --fallback chạy lại tác vụ -p."
updated: 2026-10-21
order: 8
translation: machine
sourceHash: "afee5cce02cd"
---

`ccam next` mở Claude Code bằng một tài khoản khác mà bạn không cần nêu tên. Hãy dùng nó khi một tài khoản chạm giới hạn sử dụng: `ccam next` không kèm tùy chọn chọn tài khoản đứng sau tài khoản bạn dùng gần nhất, `--auto` chọn tài khoản còn nhiều hạn mức nhất, và `--continue` chuyển cuộc hội thoại gần nhất của thư mục này sang tài khoản mới, dù hai tài khoản có dùng chung lịch sử hay không. Bên trong Claude Code, [plugin ccam](#inside-claude-code) báo cho bạn khi một giới hạn dừng một phiên, và `/ccam:switch` chuyển cuộc hội thoại đi tiếp. Với các script chạy `claude -p`, `ccam --fallback` tự thực hiện việc đổi (xem [Chạy lại khi chạm giới hạn](#run-again-at-the-limit)).

## Cách dùng

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the most room
ccam next --continue            # carry this folder's last conversation on as the next account
ccam next -c --auto             # both: the conversation, on the account with the most room
ccam next --auto --max-usage 80 # prefer accounts below 80%
ccam next --auto --strategy balance   # the weekly quota that would go unused first
ccam next --auto --cross-org    # also accounts of other organizations
ccam next --auto --wait         # nobody has room: wait for the first reset (up to 5 hours)
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options go after --
ccam --dry-run next --auto      # print the choice and the command, start nothing
```

```console
$ ccam --dry-run next --auto --continue
→ personal → work (12% of its 5 hours used, the least)
  continuing "Fix the flaky tests" (about 182k tokens); work does not share personal's history, so Claude Code copies it into work
CLAUDE_CONFIG_DIR=/home/you/.local/share/ccam/profiles/work /usr/local/bin/claude --resume /home/you/.local/share/ccam/profiles/personal/projects/-home-you-code-app/4f1c….jsonl --fork-session
```

ccam in việc chuyển đổi dưới dạng `→ <from> → <to>` ra đầu ra lỗi chuẩn, kèm lý do của một lựa chọn `--auto`, rồi khởi động Claude Code bằng tài khoản mới theo cùng cách `ccam <name>` làm. `<from>` là tài khoản bạn dùng gần nhất, hoặc với `--continue` là tài khoản sở hữu cuộc hội thoại. Các đối số sau `--` được chuyển nguyên vẹn cho `claude`.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--auto` | Chọn tài khoản đã đăng nhập theo mức sử dụng (xem bên dưới), bỏ qua các tài khoản ở mức 100% | tắt |
| `--max-usage <1-100>` | Với `--auto`: ưu tiên các tài khoản có mức sử dụng dưới giá trị này | `auto_max_usage`, nếu không thì 100 |
| `--strategy headroom\|expiring\|balance` | Với `--auto`: mức sử dụng 5 giờ thấp nhất, khoảng 5 giờ được đặt lại sớm nhất, hoặc hạn mức hàng tuần sẽ bị bỏ phí sớm nhất | `auto_strategy`, nếu không thì `headroom` |
| `--cross-org` | Với `--auto`: chọn cả các tài khoản thuộc tổ chức khác (xem bên dưới) | `cross_org`, nếu không thì tắt |
| `--wait[=<duration>]` | Với `--auto`: khi không tài khoản nào còn hạn mức, chờ tài khoản đầu tiên được đặt lại, tối đa chừng này thời gian | tắt; `5h` khi chỉ ghi cờ |
| `-c`, `--continue` | Chuyển cuộc hội thoại gần nhất của thư mục này sang tài khoản mới (xem bên dưới) | tắt |
| `-- <claude args>` | Các đối số được chuyển cho `claude` sau các đối số của chính ccam | không có |
| `--dry-run` | Được hỗ trợ: in ra `→ <from> → <to>` và dòng lệnh sẽ chạy, mà không khởi động Claude Code hay ghi lại lần khởi động | tắt |

`--auto` và `--continue` dùng được cùng nhau. Với `--continue`, các đối số `--resume` được đặt trước các đối số của bạn.

Các tùy chọn riêng của Claude Code, như `--dangerously-skip-permissions` hoặc `--model`, đặt sau `--`. Khác với `ccam <name>`, vốn chuyển tiếp mọi thứ sau tên, `ccam next` tự đọc các đối số của nó, nên nếu thiếu `--` nó sẽ dừng với `ccam: unknown flag: --dangerously-skip-permissions`.

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

Một tài khoản được đánh giá theo khoảng 5 giờ và khoảng hàng tuần của nó, và khi `--model` nằm trong các đối số cho Claude Code (sau `--`), theo cả khoảng hàng tuần của mô hình đó (`opus`, `sonnet`, `haiku` hoặc `fable` trong tên). `Policy.Choose()` chia các tài khoản thành ba nhóm:

1. **Còn hạn mức**: mọi khoảng đều dưới `auto_max_usage` (`config.toml`, mặc định 100). Chiến lược chọn trong số này: `headroom` (mặc định) lấy mức sử dụng 5 giờ thấp nhất; `expiring` lấy khoảng 5 giờ được đặt lại sớm nhất, để hạn mức còn lại được dùng trước khi mất; `balance` lấy tài khoản còn nhiều hạn mức hàng tuần nhất cho mỗi ngày cho đến khi tuần của nó được đặt lại, có tính cả từng phiên ccam đang chạy trên tài khoản đó, để hạn mức hàng tuần không bị bỏ phí và các phiên song song được phân tán.
2. **Vượt ngưỡng**: dưới 100% nhưng bằng hoặc vượt `auto_max_usage`. Chỉ được dùng khi không tài khoản nào còn hạn mức, ưu tiên mức sử dụng 5 giờ thấp nhất.
3. **Đầy**: có một khoảng ở mức 100%, hoặc chưa rõ mức sử dụng 5 giờ. Không bao giờ được chọn.

```toml
# ~/.config/ccam/config.toml
auto_max_usage = 80        # keep 20% of each account for interactive work
auto_strategy = "expiring"
```

**Tổ chức.** Các lựa chọn tự động chỉ giữ ở những tài khoản giống tài khoản bạn rời đi: cùng tổ chức team hoặc enterprise, hoặc các gói cá nhân (Pro, Max) với nhau, để công việc từ một tài khoản công ty không tiếp tục trên một tài khoản riêng tư hay ngược lại. ccam biết tổ chức của mỗi tài khoản từ thông tin đăng nhập của nó (`ccam login` ghi lại; các lần đăng nhập cũ hơn được đọc từ `.claude.json`); một tài khoản mà nó không xác định được thì hợp với mọi nhóm. Hai thư mục đăng nhập vào cùng một tài khoản Claude dùng chung giới hạn của nó, nên thư mục này không bao giờ thay cho thư mục kia. `--cross-org`, hoặc `cross_org = true` trong `config.toml`, bỏ quy tắc tổ chức này.

Khi không chọn được tài khoản nào, ccam thoát với mã 75 (`EX_TEMPFAIL`), để script phân biệt được "thử lại sau" với một lỗi. Với `--wait`, thay vào đó nó ngủ cho đến khi khoảng đang chặn của tài khoản đầy đầu tiên được đặt lại (30 giây sau thời điểm đặt lại, khi endpoint hiển thị khoảng mới), in ra `ccam: no account has room; waiting until 15:10`, rồi chọn lại; khi không tài khoản nào cho biết lúc nào được đặt lại, nó kiểm tra lại sau 10 phút. Nếu không có tài khoản nào trống ra trước khi hết thời gian chờ, nó thoát với mã 75. Việc tra cứu tuân theo các quy tắc của `ccam usage`: các con số mới hơn 60 giây được dùng lại, mỗi tài khoản có 3 giây, và tài khoản không phản hồi hoặc có đăng nhập đã hết hạn sẽ được đánh giá theo các con số đã lưu. Một khoảng đã lưu mà thời điểm đặt lại đã qua được tính là 0%, nên tài khoản có khoảng giới hạn đã được đặt lại kể từ đó sẽ không bị bỏ qua. `--auto` không bao giờ gia hạn token. Khác với vòng xoay thông thường, `--auto` có thể chọn tài khoản bạn dùng gần nhất nếu nó vẫn còn nhiều hạn mức nhất.

### --continue

1. ccam tìm cuộc hội thoại gần nhất của thư mục này trong mọi tài khoản: `<session>.jsonl` mới nhất trong `projects/<folder>/` (tên thư mục là thư mục làm việc với mọi ký tự không phải chữ cái hay chữ số được đổi thành `-`, theo cách Claude Code đặt tên). Nếu phiên đó để lại một ghi chú (một giới hạn sử dụng hoặc `/ccam:switch`, xem bên dưới) và không có gì mới hơn được ghi kể từ đó, chính phiên đó được dùng.
2. Tài khoản sở hữu cuộc hội thoại là tài khoản bạn chuyển đi, và tài khoản mới được chọn như trên, trong cùng tổ chức.
3. Nếu tài khoản mới dùng chung lịch sử đó (cùng thư mục nhóm), ccam khởi động `claude --resume <session id>`, tức chính cuộc hội thoại đó (`--continue` sẽ mở cuộc hội thoại nào mới nhất trong thư mục).
4. Nếu không, nó khởi động `claude --resume <the conversation's file> --fork-session`: Claude Code sao chép toàn bộ cuộc hội thoại, dưới một id phiên mới, vào lịch sử riêng của tài khoản mới, và bản gốc vẫn giữ nguyên. Từ đó trở đi, `--continue` của tài khoản mới sẽ tìm thấy nó.

ccam cho biết cuộc hội thoại lớn đến đâu: tài khoản mới đọc lại toàn bộ nó với câu trả lời đầu tiên. Bộ nhớ đệm prompt không bao giờ được dùng chung giữa các tổ chức, nên với `--cross-org` câu trả lời đầu tiên đó tốn toàn bộ kích thước. Dùng chung `history` trong một nhóm tránh được việc sao chép, giữ một cuộc hội thoại duy nhất, và còn dùng chung các checkpoint mà `/rewind` dùng để khôi phục mã (`file-history/`).

## Bên trong Claude Code

[Plugin ccam](/docs/install#claude-code-plugin) thêm hai hook (Claude Code 2.1.78 trở lên):

- Khi một lượt dừng lại vì giới hạn sử dụng (`You've hit your session limit`, giới hạn hàng tuần hoặc giới hạn của một mô hình), nó ghi chú cuộc hội thoại và hiện một thông báo trên màn hình ở các terminal hỗ trợ (iTerm2, WezTerm, Windows Terminal, kitty, Ghostty): `ccam: work hit its 5-hour limit (resets 3:45pm). ccam next --continue carries on as another account.` Một lần máy chủ điều tiết ngắn hoặc giới hạn tốc độ của API key không phải là giới hạn sử dụng và không bị đụng tới. Claude Code 2.1.234 trở lên còn chờ trong phiên và tự tiếp tục sau khi khoảng 5 giờ được đặt lại; thông báo cho bạn lựa chọn chuyển tài khoản thay vì chờ.
- `/ccam:switch [account]` ghi chú cuộc hội thoại và kết thúc lượt trước khi có yêu cầu nào được gửi: gõ `/exit`, rồi `ccam next --continue` chuyển nó đi tiếp, bằng tài khoản bạn đã nêu tên hoặc bằng tài khoản kế tiếp.

Ghi chú nằm trong `~/.local/state/ccam/handoff/<account>.json` trong 12 giờ và chỉ được dùng một lần, khi phiên tiếp theo sắp bắt đầu; một lần khởi động thất bại sẽ để nguyên nó cho lần thử sau. Ghi chú thuộc về thư mục nơi phiên bắt đầu, kể cả sau một lệnh `cd` trong shell của Claude Code. Nếu tài khoản mà ghi chú nêu tên đã bị đổi tên hoặc bị xóa kể từ đó, ccam sẽ báo điều này và chọn như bình thường. Với `relay = true` trong `config.toml`, ccam ở lại bên cạnh các phiên tương tác thay vì tự thay thế mình bằng Claude Code, và sau `/ccam:switch` nó tự chuyển cuộc hội thoại của phiên đó đi tiếp (ghi chú do một phiên khác để lại sẽ không được xử lý). Việc này tốn một tiến trình ccam đang chờ cho mỗi phiên, đó là lý do nó tắt theo mặc định. Phiên được chuyển tiếp khởi động bằng `ccam next --continue`, không kèm các đối số mà phiên đầu tiên được khởi động cùng, chẳng hạn `--model`.

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
8. Chỉ các tài khoản thuộc cùng tổ chức được thử, như với `--auto` (`cross_org = true` bỏ giới hạn này).

Tài khoản tiếp theo được chọn theo cùng các quy tắc như `--auto` (`auto_max_usage`, `auto_strategy`, khoảng của `--model`). Khi không còn tài khoản nào, ccam in ra đầu ra của lần thử cuối, báo rằng không còn tài khoản đã đăng nhập nào còn hạn mức, và thoát với mã 75. Với `--wait` đặt trước tên tài khoản (`ccam --fallback --wait=2h work -p …`), thay vào đó nó chờ tài khoản đầu tiên được đặt lại, rồi có thể thử lại mọi tài khoản; nó luôn ngủ ít nhất một lần trước, nên tài khoản vừa chạm giới hạn sẽ không bị chạy lại ngay. `ccam --dry-run --fallback work -p …` in ra lệnh của lần thử đầu tiên.

Các phiên tương tác không tự đổi tài khoản. Thay vào đó, với `ccam statusline install`, dòng trạng thái thêm `→ ccam next --auto --continue` khi mức sử dụng 5 giờ hoặc hàng tuần chạm 90% ([Giới hạn sử dụng](/features/usage)).

## An toàn và các trường hợp đặc biệt

- Tài khoản chưa đăng nhập không bao giờ được chọn. Khi không có tài khoản nào khác đã đăng nhập, `ccam next` thông thường dừng lại với `no logged-in account other than <name>`.
- Với `--continue` và không có cuộc hội thoại nào trong thư mục: `no account has a conversation in <folder> to continue`.
- Khi quy tắc tổ chức không để lại tài khoản nào: `` no logged-in account to pick among <name>'s (<organization>); --cross-org also looks at the others ``.
- Với `--auto` và không có tài khoản nào đã đăng nhập: `no logged-in account to pick`. Khi mọi tài khoản đều đã đầy hoặc chưa rõ mức sử dụng: `` no account with known usage and headroom (see `ccam usage`) ``. Cả hai đều kết thúc với mã thoát 75.
- `--max-usage`, `--strategy`, `--wait` và `--cross-org` chỉ dùng được với `--auto`, và giá trị ngoài khoảng 1-100 hoặc một chiến lược khác sẽ bị từ chối.
- Một bản fork sao chép cuộc hội thoại, không phải trạng thái khác của tài khoản: các checkpoint trong `file-history/` của tài khoản gốc không đi theo, trừ khi lịch sử được dùng chung.
- `--auto` gửi yêu cầu mạng ngay cả với `--dry-run`. Nó chỉ đọc, và không bao giờ ghi vào tệp token.
- Tài khoản được chọn khởi động qua cùng đường đi với `ccam <name>`, nên các kiểm tra tương tự được áp dụng, bao gồm cảnh báo khi một Claude Code khác dùng chung lịch sử đang chạy trong thư mục này.
- `last-used` chỉ được ghi khi Claude Code thực sự khởi động. `--dry-run` không động đến nó, nên chạy thử không đẩy vòng xoay tiến lên.
- `--fallback` giữ đầu ra chuẩn cho đến khi một lần thử kết thúc, nên đầu ra dạng luồng như `--output-format stream-json` đến một lần vào cuối. Các tệp mà lần thử đầu tiên đã thay đổi trước khi chạm giới hạn vẫn giữ nguyên thay đổi; lần thử tiếp theo bắt đầu lại lời nhắc từ đầu.

## Liên quan

- [Giới hạn sử dụng](/features/usage)
- [Giới hạn sử dụng và chuyển đổi tài khoản](/docs/usage-limits)
- [Dùng chung máy chủ MCP, cài đặt và lịch sử](/docs/sharing)
- [Quy trình khi chạm giới hạn sử dụng](/guides/usage-limit-workflow)
