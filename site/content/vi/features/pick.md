---
title: Chọn tài khoản cho script
description: ccam pick in tài khoản còn nhiều hạn mức nhất cho script, giao tài khoản khác nhau cho các worker khởi động cùng lúc, và thoát 75 khi hết hạn mức.
updated: 2026-10-20
order: 12
translation: machine
sourceHash: 65a5a7e179a5
---

`ccam pick` in tên của tài khoản đã đăng nhập còn nhiều hạn mức nhất, để một script dùng nó khởi động Claude Code. Các worker khởi động cùng lúc sẽ mỗi worker nhận một tài khoản khác nhau, miễn là có đủ tài khoản còn hạn mức, nên một loạt tác vụ `claude -p` được trải ra trên các tài khoản của bạn thay vì dùng cạn một tài khoản.

## Cách dùng

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live"
}
```

Mã thoát: `0` kèm một tên trên đầu ra chuẩn, `75` (`EX_TEMPFAIL`) khi không tài khoản nào còn hạn mức hoặc không có tài khoản nào đăng nhập, `1` cho các lỗi khác như một cờ sai.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--json` | In tài khoản, mức sử dụng 5 giờ và hàng tuần của nó, và nguồn của các con số | tắt |
| `--model <model>` | Tính cả khoảng hàng tuần của mô hình đó, chẳng hạn `opus` | không có |
| `--max-usage <1-100>` | Ưu tiên các tài khoản có mức sử dụng dưới giá trị này | `auto_max_usage`, nếu không thì 100 |
| `--strategy headroom\|expiring` | Mức sử dụng 5 giờ thấp nhất, hoặc khoảng 5 giờ được đặt lại sớm nhất | `auto_strategy`, nếu không thì `headroom` |
| `--wait[=<duration>]` | Khi không tài khoản nào còn hạn mức, chờ tài khoản đầu tiên được đặt lại, tối đa chừng này thời gian | tắt; `5h` khi chỉ ghi cờ |
| `--dry-run` | Chọn mà không ghi lại lần chọn | tắt |

## Cách hoạt động

1. ccam liệt kê các tài khoản đã đăng nhập (tài khoản API key không có giới hạn nên bị loại ra).
2. Nó lấy một khóa (`~/.local/state/ccam/picks.lock`), nên các lần chọn chạy cùng lúc sẽ diễn ra lần lượt.
3. Nó tra cứu mức sử dụng của mọi tài khoản cùng lúc, theo cùng cách với `ccam usage`, với bộ nhớ đệm 60 giây.
4. Các tài khoản được chọn trong 2 phút gần nhất (`~/.local/state/ccam/picks/<name>`) bị loại ra trước, và [quy tắc chọn](/features/next#auto) của `ccam next --auto` chọn trong số còn lại: ưu tiên dưới `auto_max_usage`, theo chiến lược, không bao giờ chọn tài khoản ở mức 100%.
5. Nếu mọi tài khoản còn hạn mức đều vừa được chọn gần đây, các quy tắc sẽ chọn lại trong tất cả, nên worker thứ tư trên ba tài khoản sẽ dùng chung một tài khoản.
6. Lần chọn được ghi lại và khóa được nhả ra, rồi tên được in ra.

Khoảng 2 phút này quan trọng vì mức sử dụng của một tác vụ mới cần một lúc mới hiện ra trong các con số: nếu không có nó, năm worker khởi động cùng lúc đều sẽ thấy cùng một tài khoản là tài khoản còn nhiều hạn mức nhất.

## An toàn và các trường hợp đặc biệt

- `ccam pick` chỉ đọc mức sử dụng và ghi thời điểm chọn; nó không bao giờ gia hạn hay đọc token ngoài những gì `ccam usage` làm.
- Một lần chọn không giữ chỗ tài khoản. Sau 2 phút, hai script vẫn có thể khởi động trên cùng một tài khoản; hãy dùng `--fallback` để phục hồi khi tài khoản đó chạm giới hạn.
- Việc chờ với `--wait` không giữ khóa, nên các lần chọn khác vẫn diễn ra trong lúc đó.
- `ccam pick` bỏ qua [`ccam folder`](/features/folder): nó chỉ chọn theo mức sử dụng.

## Liên quan

- [Chuyển sang tài khoản tiếp theo](/features/next)
- [Chạy lại khi chạm giới hạn](/features/next#run-again-at-the-limit)
- [Giới hạn sử dụng](/features/usage)
- [Tài khoản API key](/features/api)
