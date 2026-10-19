---
title: Kiểm tra và sửa chữa
description: "Kiểm tra liên kết dùng chung, quyền tệp, thông tin đăng nhập và phiên bản Claude Code bằng ccam doctor, và sửa những gì an toàn để sửa bằng ccam doctor --fix."
updated: 2026-10-19
order: 9
translation: machine
sourceHash: a1b61efe2b1f
---

`ccam doctor` kiểm tra bản cài đặt ccam của bạn: tệp thực thi Claude Code, thư mục, thông tin đăng nhập và liên kết dùng chung của mọi tài khoản, các máy chủ MCP dùng chung, endpoint mức sử dụng và các lệnh tắt của bạn. `ccam doctor --fix` sửa những gì nó có thể sửa mà không mất dữ liệu, và chuyển mọi tệp nó thay thế vào một thư mục sao lưu trước. Hãy chạy lệnh này khi `ccam <name>` dừng lại vì một vấn đề về liên kết dùng chung, sau khi bạn khôi phục tệp bằng tay, hoặc sau một lần cập nhật Claude Code.

## Cách dùng

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair
```

Mỗi kết quả là một dòng có ký hiệu: `✓` ổn, `·` thông tin, `!` cảnh báo, `✗` lỗi. Kết quả của tài khoản bắt đầu bằng tên tài khoản.

```console
$ ccam doctor
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700)
! personal: .credentials.json permissions are 644 (should be 0600)
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp

`ccam doctor --fix` repairs what it can.
```

Với `--fix`, mọi kết quả đã được sửa đều kết thúc bằng `→ fixed`:

```console
$ ccam doctor --fix
✓ Claude Code 2.1.302 (Claude Code) (/home/me/.local/bin/claude)
! personal: folder permissions are 755 (should be 0700) → fixed
! personal: .credentials.json permissions are 644 (should be 0600) → fixed
✓ personal: login valid (me@example.com)
✗ personal: /home/me/.local/share/ccam/profiles/personal/settings.json: set to shared but is a real file, not a link → fixed
· personal: unknown entries (Claude Code layout may have changed): newthing
! work: not logged in → `ccam login work`
✓ work: share links OK
! missing shortcuts in /home/me/.local/bin: ccp → fixed
```

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--fix` | Sửa những gì an toàn để sửa (các tệp bị thay đổi được sao lưu) | tắt |
| `--dry-run` (toàn cục) | Chạy mọi kiểm tra nhưng không sửa gì, kể cả với `--fix` | tắt |

`ccam doctor` không nhận đối số. Nó hỗ trợ `--dry-run`: `ccam --dry-run doctor --fix` in ra cùng các kết quả như `ccam doctor` và không thay đổi gì.

Mã thoát là `1` khi còn ít nhất một kết quả `✗` chưa được sửa, và `0` trong các trường hợp khác. Cảnh báo và thông tin không bao giờ làm lần chạy thất bại. Khi lần chạy thất bại mà không có `--fix`, ccam kết thúc bằng một dòng trống và `` `ccam doctor --fix` repairs what it can. ``. Một cấu hình không tải được cũng kết thúc với `1` và một lỗi `ccam:`.

## Cách hoạt động

`internal/cli/doctor.go` tải `~/.config/ccam/config.toml` và gọi `doctor.Run` trong `internal/doctor/doctor.go`, hàm này chạy các kiểm tra theo một thứ tự cố định và trả về kết quả. `--fix` được truyền qua dưới dạng `Options.Fix`, và được đặt thành false khi chạy thử.

<!-- motion:doctor-checks -->

Khối ở giữa chạy một lần cho mỗi tài khoản. Các kiểm tra, và việc `--fix` làm với từng kiểm tra:

1. **Phiên bản Claude Code** (`checkClaude`). Tìm tệp thực thi `claude` theo cùng cách trình khởi chạy làm và chạy `claude --version`. Nếu phiên bản không bắt đầu bằng `2.1.302`, phiên bản mà cấu trúc thư mục của ccam đã được kiểm chứng (`doctor.LayoutVersion`), nó thêm `· ccam was verified against the folder layout of Claude Code 2.1.302. Any 'unknown entries' below may mean the layout changed`. Thiếu tệp thực thi hoặc `--version` thất bại là lỗi. Không có cách sửa.
2. **Thông tin đăng nhập trong biến môi trường** (`checkAuthEnv`). Cảnh báo khi `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` hoặc `ANTHROPIC_AUTH_TOKEN` được đặt: mọi tài khoản khởi động từ shell đó đều thừa hưởng nó, và Claude Code ưu tiên nó hơn thông tin đăng nhập riêng của tài khoản. Không có cách sửa; bạn tự bỏ đặt biến đó.
3. **Thư mục** (`checkProfile`). Thiếu thư mục tài khoản là lỗi và kết thúc các kiểm tra của tài khoản đó. Quyền có bất kỳ bit group hoặc other nào được đặt là cảnh báo; `--fix` chạy `chmod 0700`.
4. **Tệp đăng nhập** (`checkCredentials`). `.credentials.json` được kiểm tra bằng `lstat`. Một liên kết tượng trưng là lỗi, vì mỗi tài khoản phải có token riêng, và không được sửa. Quyền khác `0600` là cảnh báo; `--fix` chạy `chmod 0600`. Sau đó thông tin đăng nhập được đọc (trên macOS thì đọc từ Chuỗi khóa (Keychain) trước). Chưa đăng nhập là cảnh báo kèm gợi ý `ccam login <name>`. Refresh token hết hạn là lỗi `login expired`; access token hết hạn chỉ là thông tin, vì Claude Code sẽ gia hạn nó ở lần khởi động tiếp theo. Không trường hợp nào trong số này được sửa.
5. **Liên kết dùng chung** (`share.Check`). So sánh mọi đường dẫn trong `share.Layout` với thiết lập dùng chung của tài khoản, chỉ dùng `lstat` và `readlink`. Nó báo bốn loại: thiếu liên kết, liên kết trỏ sang chỗ khác, một tệp thật ở nơi lẽ ra là liên kết, và liên kết còn sót lại cho một mục nay đã là riêng. Thiếu liên kết là cảnh báo, các loại còn lại là lỗi.
6. **Mục không xác định** (`checkUnknownEntries`). Liệt kê các tên ở cấp cao nhất trong thư mục tài khoản mà ccam không biết (`share.Known`), như một dấu hiệu cho thấy cấu trúc của Claude Code đã thay đổi. Chỉ là thông tin.
7. **Đường dẫn plugin** (`checkPlugins`). Các tệp đăng ký plugin (`installed_plugins.json`, `known_marketplaces.json`) lưu đường dẫn tuyệt đối. Trong `plugins/` của mỗi nhóm, các đường dẫn đi qua thư mục của một tài khoản sẽ hỏng khi tài khoản đó bị xóa hoặc đổi tên. `--fix` ghi lại chúng thành thư mục nhóm bằng `profile.RewritePluginPaths`, một lần ghi nguyên tử vào tệp đăng ký.
8. **Máy chủ MCP** (`checkMCP`). Tải `mcp.json` của mỗi nhóm một lần. Tệp không đọc được là lỗi; máy chủ không qua được `mcp.Validate` là cảnh báo. Không có cách sửa.
9. **Định dạng phản hồi mức sử dụng** (`checkUsageFormat`). Hỏi endpoint mức sử dụng, không dùng bộ nhớ đệm và với giới hạn 3 giây thông thường, lần lượt cho từng tài khoản cho đến khi có một tài khoản phản hồi. Nếu phản hồi không còn phân tích được, nó cảnh báo rằng `ccam usage` sẽ quay về dùng bộ nhớ đệm của dòng trạng thái. Đây là kiểm tra duy nhất dùng mạng.
10. **Lệnh tắt** (`checkShims`). Chỉ khi có tài khoản có lệnh tắt. Thiếu liên kết trong `~/.local/bin` là cảnh báo; `--fix` gọi `shim.Install`, hàm này chỉ tạo các liên kết còn thiếu và không bao giờ động đến tệp khác. Thư mục bin không nằm trong `PATH` là cảnh báo và không được sửa.

### Sửa liên kết dùng chung

Đường khởi động tạo ngay một liên kết còn thiếu, nhưng dừng lại với bất cứ thứ gì có thể làm mất dữ liệu và chỉ bạn đến đây:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Với bất kỳ vấn đề liên kết dùng chung nào, `--fix` gọi `share.Build` với thiết lập dùng chung của tài khoản làm cả trạng thái hiện tại lẫn trạng thái đích, rồi `Plan.Apply`. `Build` suy ra lại mọi bước từ những gì có trên đĩa, cũng là kế hoạch mà `ccam share` dùng, nên việc sửa tuân theo cùng các quy tắc như khi bật dùng chung.

<!-- motion:doctor-repair -->

Ở đây `settings.json` là một tệp thật. Vì nhóm đã có một `settings.json` khác, xung đột nhận lựa chọn mặc định là `KeepGroup`, nên bản của nhóm được giữ. Tệp của tài khoản được chuyển vào `.ccam-backup/<time>/settings.json` bên trong thư mục tài khoản, và `os.Symlink` liên kết đường dẫn tới nhóm. Các trường hợp khác:

- Nếu nhóm chưa có bản nào, tệp của tài khoản được chuyển vào nhóm và trở thành nội dung đầu tiên của nhóm.
- Với một thư mục, các tệp nhóm còn thiếu được chuyển vào nhóm; phần còn lại vào bản sao lưu.
- Với `history.jsonl`, mặc định là `merge by line`: các dòng của nhóm, rồi đến các dòng mới của bạn. Bản trước đó của nhóm được lưu trong `.ccam-backup/<time>/group/`.
- Liên kết tới một nhóm khác được thay thế. Liên kết tới một thư mục ngoài ccam được đưa vào bản sao lưu, và nội dung của nó trở thành nội dung ban đầu của nhóm nếu nhóm còn trống.
- Liên kết còn sót lại cho một mục nay đã là riêng sẽ bị gỡ bỏ, và tài khoản nhận một bản sao nội dung của nhóm.

Ngoài việc gỡ liên kết, mọi bước đều chuyển dữ liệu vào nhóm hoặc thư mục sao lưu. Không có gì bị xóa.

## An toàn và các trường hợp đặc biệt

- `ccam doctor` thông thường không bao giờ thay đổi gì. Hãy dùng `--dry-run` cùng `--fix` để chắc chắn trong script.
- `doctor --fix` không kiểm tra tài khoản có đang chạy hay không. Hãy thoát Claude Code của một tài khoản trước khi sửa liên kết dùng chung của nó, vì phiên vẫn tiếp tục ghi vào các tệp đó.
- Nếu việc sửa liên kết dùng chung thất bại, ccam báo `could not repair share links: <error>` cho tài khoản đó và bỏ qua kiểm tra mục không xác định của nó. Các tài khoản khác vẫn được kiểm tra.
- Nếu `ccam share` thất bại giữa chừng, nó lưu thiết lập đích trước và nhắc bạn hoàn tất bằng `ccam doctor --fix`.
- `.credentials.json` là liên kết tượng trưng, thông tin đăng nhập hết hạn và thông tin đăng nhập trong biến môi trường đều được báo nhưng không bao giờ bị thay đổi.
- Việc ghi lại tệp đăng ký plugin được thực hiện nguyên tử nhưng không được sao vào `.ccam-backup/`.

## Liên quan

- [Dùng chung máy chủ MCP, cài đặt và lịch sử](/docs/sharing)
- [Sao chép thiết lập hiện có](/features/migrate)
- [Giới hạn sử dụng và chuyển đổi tài khoản](/docs/usage-limits)
- [Lệnh](/docs/commands)
