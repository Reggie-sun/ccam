---
title: Chuyển sang tài khoản tiếp theo
description: "ccam next mở Claude Code bằng tài khoản đứng sau tài khoản bạn dùng gần nhất; --auto chọn tài khoản có mức sử dụng 5 giờ thấp nhất và --continue tiếp tục cuộc hội thoại."
updated: 2026-10-19
order: 8
translation: machine
sourceHash: 92352f532e74
---

`ccam next` mở Claude Code bằng một tài khoản khác mà bạn không cần nêu tên. Hãy dùng nó khi một tài khoản chạm giới hạn sử dụng: `ccam next` không kèm tùy chọn chọn tài khoản đứng sau tài khoản bạn dùng gần nhất, `--auto` chọn tài khoản còn nhiều hạn mức nhất trong khoảng 5 giờ, và `--continue` tiếp tục cuộc hội thoại gần nhất trong tài khoản mới.

## Cách dùng

```sh
ccam next                       # the account after the one you used last
ccam next --auto                # the account with the lowest 5-hour usage
ccam next --continue            # the next account that shares history, with claude --continue
ccam next -c --auto             # both: lowest usage among accounts sharing history
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
| `--auto` | Chọn tài khoản đã đăng nhập có mức sử dụng 5 giờ thấp nhất, bỏ qua các tài khoản ở mức 100% | tắt |
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

`--auto` quyết định theo mức sử dụng thay vì thứ tự slot. `cli.pickAuto` giữ lại các tài khoản đã đăng nhập, hỏi mức sử dụng của tất cả cùng lúc bằng chính `Fetcher.FetchAll` mà `ccam usage` dùng, và để `usage.Best()` chọn:

<!-- motion:next-auto -->

`usage.Best()` lấy tỷ lệ phần trăm 5 giờ thấp nhất và bỏ qua các tài khoản ở mức 100% hoặc chưa rõ mức sử dụng 5 giờ. Việc tra cứu tuân theo các quy tắc của `ccam usage`: các con số mới hơn 60 giây được dùng lại, mỗi tài khoản có 3 giây, và tài khoản không phản hồi hoặc có đăng nhập đã hết hạn sẽ được đánh giá theo các con số đã lưu. Một khoảng đã lưu mà thời điểm đặt lại đã qua được tính là 0%, nên tài khoản có khoảng giới hạn đã được đặt lại kể từ đó sẽ không bị bỏ qua. `--auto` không bao giờ gia hạn token. Khác với vòng xoay thông thường, `--auto` có thể chọn tài khoản bạn dùng gần nhất nếu nó vẫn còn nhiều hạn mức nhất.

### --continue

`--continue` của Claude Code mở cuộc hội thoại gần nhất trong thư mục hiện tại, và điều này chỉ hoạt động nếu tài khoản mới nhìn thấy được cuộc hội thoại đó. Vì vậy với `--continue`, ccam chỉ xét các tài khoản trong cùng nhóm `history` với tài khoản bạn dùng gần nhất, rồi đặt `--continue` trước các đối số của bạn. Các tài khoản có lịch sử là `local` không dùng chung với ai, nên `--continue` từ một tài khoản như vậy luôn thất bại. Hãy bật dùng chung trước bằng `ccam share <name> --history on` (xem [dùng chung](/docs/sharing)).

## An toàn và các trường hợp đặc biệt

- Tài khoản chưa đăng nhập không bao giờ được chọn. Khi không có tài khoản nào khác đã đăng nhập, `ccam next` thông thường dừng lại với `no logged-in account other than <name>`.
- Với `--continue` và không có tài khoản nào dùng chung lịch sử: `` no other account shares conversation history with <name>. Turn it on with `ccam share <name> --history on` ``.
- Với `--auto` và không có tài khoản nào đã đăng nhập: `no logged-in account to pick`. Khi mọi tài khoản đều đã đầy hoặc chưa rõ mức sử dụng: `` no account with known usage and headroom (see `ccam usage`) ``.
- `--auto` gửi yêu cầu mạng ngay cả với `--dry-run`. Nó chỉ đọc, và không bao giờ ghi vào tệp token.
- Tài khoản được chọn khởi động qua cùng đường đi với `ccam <name>`, nên các kiểm tra tương tự được áp dụng, bao gồm cảnh báo khi một Claude Code khác dùng chung lịch sử đang chạy trong thư mục này.
- `last-used` chỉ được ghi khi Claude Code thực sự khởi động. `--dry-run` không động đến nó, nên chạy thử không đẩy vòng xoay tiến lên.

## Liên quan

- [Giới hạn sử dụng](/features/usage)
- [Giới hạn sử dụng và chuyển đổi tài khoản](/docs/usage-limits)
- [Dùng chung máy chủ MCP, cài đặt và lịch sử](/docs/sharing)
- [Quy trình khi chạm giới hạn sử dụng](/guides/usage-limit-workflow)
