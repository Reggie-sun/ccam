---
title: Kiểm tra và sửa chữa
description: "Kiểm tra liên kết dùng chung, quyền tệp, thông tin đăng nhập và phiên bản Claude Code bằng ccam doctor, và sửa những gì an toàn để sửa bằng ccam doctor --fix."
updated: 2026-10-26
order: 9
translation: machine
sourceHash: 282972ac5053
---

`ccam doctor` kiểm tra bản cài đặt ccam của bạn: tệp thực thi Claude Code, thư mục, thông tin đăng nhập và liên kết dùng chung của mọi tài khoản, các máy chủ MCP dùng chung, endpoint mức sử dụng và các lệnh tắt của bạn. `ccam doctor --fix` sửa những gì nó có thể sửa mà không mất dữ liệu, và chuyển mọi tệp nó thay thế vào một thư mục sao lưu trước. Hãy chạy lệnh này khi `ccam <name>` dừng lại vì một vấn đề về liên kết dùng chung, sau khi bạn khôi phục tệp bằng tay, hoặc sau một lần cập nhật Claude Code.

## Cách dùng

```sh
ccam doctor                  # check everything, change nothing
ccam doctor --fix            # check, then repair what is safe to repair
ccam --dry-run doctor --fix  # report only: what --fix would repair and set aside
```

Dòng đầu tiên đánh dấu mọi tài khoản cùng lúc, nên bạn thấy vấn đề nằm ở đâu trước khi đọc bất cứ thứ gì khác. Bên dưới, mỗi tài khoản có một hộp, và các kết quả không thuộc tài khoản nào (Claude Code, biến môi trường, plugin, máy chủ MCP, lệnh tắt, quy tắc thư mục, nhóm và cài đặt, `auto_switch`) nằm trong hộp `General`. Tiêu đề hộp cho thấy ký hiệu nặng nhất của tài khoản và số vấn đề của nó. Mỗi hàng bên trong là một mục đã được kiểm tra, luôn là `Folder`, `Login` và `Share links` theo thứ tự đó, rồi đến các mục khác, kèm một ký hiệu: `✔` ổn, `ℹ` thông tin, `⚠` cảnh báo, `✘` lỗi. Một thư mục không có vấn đề gì hiển thị `✔ Folder  OK`.

```console
$ ccam doctor
 ⚠ General   ✘ personal   ⚠ work

 ╭─ ⚠ General · 1 problem ───────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp              │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✘ personal · 3 problems ─────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700)                │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600)     │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link             │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯

`ccam doctor --fix` repairs what it can.
```

Với `--fix`, mọi kết quả đã được sửa đều kết thúc bằng `→ fixed`, và không còn được tính là vấn đề:

```console
$ ccam doctor --fix
 ✔ General   ✔ personal   ⚠ work

 ╭─ ✔ General ───────────────────────────────────────────────────────────────╮
 │ ✔ Claude Code  Claude Code 2.1.309 (Claude Code)                          │
 │                (/home/me/.local/bin/claude)                               │
 │ ⚠ Shortcuts    missing shortcuts in /home/me/.local/bin: ccp → fixed      │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ✔ personal ──────────────────────────────────────────────────────────────╮
 │ ⚠ Folder       folder permissions are 755 (should be 0700) → fixed        │
 │ ℹ              unknown entries (ccam knows the folder layout of Claude    │
 │                Code 2.1.309; it may have changed): newthing               │
 │ ⚠ Login        .credentials.json permissions are 644 (should be 0600) →   │
 │                fixed                                                      │
 │ ✔              login valid (me@example.com, Claude Team · Acme)           │
 │ ✘ Share links  /home/me/.local/share/ccam/profiles/personal/settings.json │
 │                : set to shared but is a real file, not a link → fixed     │
 │ ℹ              settings.json differed from the group's copy, which it now │
 │                uses; this account's copy is in                            │
 │                /home/me/.local/share/ccam/profiles/personal/.ccam-        │
 │                backup/20261023-212553/settings.json                       │
 │ ✔ Usage        usage lookup OK                                            │
 ╰───────────────────────────────────────────────────────────────────────────╯
 ╭─ ⚠ work · 1 problem ──────────────────────────────────────────────────────╮
 │ ✔ Folder       OK                                                         │
 │ ⚠ Login        not logged in → `ccam login work`                          │
 │ ✔ Share links  share links OK                                             │
 ╰───────────────────────────────────────────────────────────────────────────╯
```

Trên terminal, các hộp vừa với chiều rộng của terminal, thông báo dài được ngắt dòng và màu của chính terminal được dùng (`NO_COLOR` tắt màu). Một hộp không bao giờ hẹp hơn 29 cột, nên không từ nào bị cắt trên terminal hẹp. Khi gửi tới một pipe hoặc tệp, đầu ra không có màu và không thông báo nào bị ngắt dòng hay cắt bớt, nên `ccam doctor | grep settings.json` tìm được cả dòng; khi đó mọi hộp rộng bằng dòng dài nhất. Thẻ Health của bảng điều khiển hiển thị cùng các hộp này.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--fix` | Sửa những gì an toàn để sửa (các tệp bị thay đổi được sao lưu) | tắt |
| `--dry-run` (toàn cục) | Chạy mọi kiểm tra nhưng không sửa gì, kể cả với `--fix` | tắt |

`ccam doctor` không nhận đối số. Nó hỗ trợ `--dry-run`: `ccam --dry-run doctor --fix` in ra các kết quả của `ccam doctor`, thêm những gì `--fix` sẽ để riêng ra, và không thay đổi gì:

```text
ℹ settings.json differs from the group's copy: --fix would use the group's, and move this account's copy to /home/me/.local/share/ccam/profiles/personal/.ccam-backup/20261023-212543/settings.json
```

Mã thoát là `1` khi còn ít nhất một kết quả `✘` chưa được sửa, và `0` trong các trường hợp khác. Cảnh báo và thông tin không bao giờ làm lần chạy thất bại. Khi lần chạy thất bại mà không có `--fix`, ccam kết thúc bằng một dòng trống và `` `ccam doctor --fix` repairs what it can. ``. Một cấu hình không tải được cũng kết thúc với `1` và một lỗi `ccam:`.

## Cách hoạt động

Các kiểm tra chạy theo một thứ tự cố định, và kết quả xuất hiện trong một hộp cho mỗi tài khoản. ccam tải `~/.config/ccam/config.toml` (`internal/cli/doctor.go`) và chạy các kiểm tra, mỗi kết quả nêu tài khoản nó thuộc về và mục đã kiểm tra (`doctor.Run` trong `internal/doctor/doctor.go`). Sau đó nó gom các kết quả vào các hộp (`tui.HealthReport`). `--fix` được truyền qua dưới dạng `Options.Fix`, và được đặt thành false khi chạy thử.

<!-- motion:doctor-checks -->

Khối ở giữa chạy một lần cho mỗi tài khoản. Các kiểm tra, và việc `--fix` làm với từng kiểm tra:

1. **Phiên bản Claude Code** (`checkClaude`). Tìm tệp thực thi `claude` theo cùng cách trình khởi chạy làm và chạy `claude --version`. Thiếu tệp thực thi, một `claude_bin` không phải là tệp thực thi được, hoặc `--version` thất bại là lỗi. Không có cách sửa.
2. **Thông tin đăng nhập trong biến môi trường** (`checkAuthEnv`). Cảnh báo khi `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` hoặc `ANTHROPIC_AUTH_TOKEN` được đặt: mọi tài khoản khởi động từ shell đó đều thừa hưởng nó, và Claude Code ưu tiên nó hơn thông tin đăng nhập riêng của tài khoản. Không có cách sửa; bạn tự bỏ đặt biến đó.
3. **Thư mục** (`checkProfile`). Thiếu thư mục tài khoản là lỗi và kết thúc các kiểm tra của tài khoản đó. Quyền có bất kỳ bit group hoặc other nào được đặt là cảnh báo; `--fix` chạy `chmod 0700`. Một `~/.claude` đã nhập có `.claude.json` riêng bên cạnh `~/.claude.json` là thông tin: khi đó `ccam main` và `claude` thông thường giữ thông tin đăng nhập và máy chủ MCP tách riêng. Windows không có bit quyền, nên kiểm tra này và kiểm tra `0600` bên dưới được bỏ qua ở đó; thư mục trong hồ sơ người dùng của bạn được bảo vệ bằng danh sách kiểm soát truy cập của nó.
4. **Tệp đăng nhập** (`checkCredentials`). `.credentials.json` được kiểm tra bằng `lstat`. Một liên kết tượng trưng là lỗi, vì mỗi tài khoản phải có token riêng, và không được sửa. Quyền khác `0600` là cảnh báo; `--fix` chạy `chmod 0600`.

   Sau đó thông tin đăng nhập được đọc (trên macOS thì đọc từ Chuỗi khóa (Keychain) trước). Thiếu thông tin đăng nhập, hoặc một tệp đăng nhập không chứa thông tin đăng nhập nào, là cảnh báo kèm gợi ý `ccam login <name>`. Refresh token hết hạn, hoặc thông tin đăng nhập mà Claude Code đã xóa sau khi máy chủ từ chối nó, là lỗi `` login expired → `ccam login <name>` ``; access token hết hạn chỉ là thông tin, vì Claude Code sẽ gia hạn nó ở lần khởi động tiếp theo. Thông tin đăng nhập được giữ trong Credential Manager của Windows là thông tin, vì ccam không đọc được nó. Một thông tin đăng nhập hợp lệ hiển thị email mà thư mục đang đăng nhập hiện tại và gói (`login valid (dev@company.com, Claude Team · Acme)`). Không trường hợp nào trong số này được sửa; `ccam login --expired` đăng nhập mọi tài khoản cần đăng nhập lại.
5. **Liên kết dùng chung** (`share.Check`). So sánh mọi đường dẫn trong `share.Layout` với các nhóm của tài khoản, chỉ dùng `lstat` và `readlink`. Nó báo năm loại: thiếu liên kết, liên kết trỏ sang chỗ khác, một tệp thật ở nơi lẽ ra là liên kết, liên kết còn sót lại cho một mục nay đã là riêng, và một liên kết đúng nhưng thư mục nhóm của nó bị thiếu (`the group folder it links to is missing (…), so Claude Code cannot save there (ccam doctor --fix makes it)`; việc khởi chạy không kiểm tra điều này). Thiếu liên kết và thiếu thư mục nhóm là cảnh báo, các loại còn lại là lỗi.

   Một lần sao chép ngược từ nhóm bị dừng giữa chừng cũng được báo: `--fix` đặt một bản sao đã chuẩn bị đầy đủ vào chỗ, nhưng nếu một lần khởi chạy đã tạo lại đường dẫn đó kể từ lúc ấy, nó nêu tên thư mục chứa bản sao đầy đủ để bạn tự xử lý. Các đường dẫn được thêm vào một mục trong phiên bản ccam mới hơn (`file-history/`, `tasks/` và `paste-cache/` cho lịch sử, `keybindings.json` cho cài đặt) chỉ là cảnh báo khi tài khoản vẫn còn bản riêng: `is not shared with the group yet (ccam doctor --fix shares it; if the group has its own copy, that one is kept and this one backed up)`. Việc khởi chạy không bao giờ dừng lại vì chúng.
6. **Mục không xác định** (`checkUnknownEntries`). Liệt kê các tên ở cấp cao nhất trong thư mục tài khoản mà ccam không biết (`share.Known`), như một dấu hiệu cho thấy cấu trúc của Claude Code đã thay đổi. Chỉ là thông tin.
7. **Đường dẫn plugin** (`checkPlugins`). Các tệp đăng ký plugin (`installed_plugins.json`, `known_marketplaces.json`) lưu đường dẫn tuyệt đối. Trong `plugins/` của mỗi nhóm, các đường dẫn đi qua thư mục của một tài khoản sẽ hỏng khi tài khoản đó bị xóa hoặc đổi tên. `--fix` ghi lại chúng thành thư mục nhóm bằng `profile.RewritePluginPaths`, một lần ghi nguyên tử vào tệp đăng ký.
8. **Máy chủ MCP** (`checkMCP`). Tải `mcp.json` của mỗi nhóm một lần. Tệp không đọc được là lỗi; máy chủ không qua được `mcp.Validate` là cảnh báo. Các máy chủ trong `~/.claude.json` của `claude` thông thường mà không tài khoản nào tải là cảnh báo, kèm lệnh `ccam migrate --mcp` để sao chép chúng (với `--group <a|b>` khi nhiều nhóm dùng chung MCP). Các máy chủ mà chỉ một tài khoản trong nhóm có là thông tin, kèm lệnh `ccam mcp import` để dùng chung chúng, và một tài khoản chưa đăng nhập vào các máy chủ MCP mà các tài khoản khác đã đăng nhập sẽ nhận cảnh báo đăng nhập bằng `/mcp` (ccam không bao giờ sao chép các lần đăng nhập). Không có cách sửa.
9. **Định dạng phản hồi mức sử dụng** (`checkUsageFormat`). Hỏi endpoint mức sử dụng, không dùng bộ nhớ đệm, lần lượt cho từng tài khoản cho đến khi có một tài khoản phản hồi, trong một thời gian chờ tối đa duy nhất là 4 giây. Với `usage_source = "claude"`, nó hỏi Claude Code cho một tài khoản thay vào đó. Nếu phản hồi không còn phân tích được, nó cảnh báo rằng `ccam usage` sẽ quay về dùng bộ nhớ đệm của dòng trạng thái; khi không tài khoản nào phản hồi, nó báo điều đó dưới dạng thông tin (`could not check the usage response format: …`). Đây là kiểm tra duy nhất dùng mạng.
10. **Lệnh tắt** (`checkShims`). Chỉ khi có tài khoản có lệnh tắt. Thiếu liên kết trong `~/.local/bin` là cảnh báo, và các lệnh tắt mà ccam đã tạo nhưng vẫn chạy một phiên bản ccam cũ hơn, một bản đã bị thay thế hoặc di chuyển, cũng vậy (`shortcuts in <bin> that run another ccam (one replaced or moved): …`); `--fix` gọi `shim.Install`, hàm này chỉ tạo và làm mới các liên kết mà ccam đã tạo và ghi lại trong `~/.local/state/ccam/shortcuts.json`, và không bao giờ động đến tệp khác. Cảnh báo chỉ kết thúc bằng `→ fixed` khi từng liên kết đó đã được tạo; nếu có liên kết không tạo được, lý do sẽ được thêm vào. Một lệnh tắt chạy một ccam khác vẫn còn được cài thì vẫn hoạt động, nên `ccam doctor` chỉ liệt kê nó dưới dạng thông tin; `--fix` trỏ nó về tệp thực thi ccam này. Thư mục bin không nằm trong `PATH` là cảnh báo và không được sửa.
11. **Thư mục theo tài khoản** (`checkFolders`). Chỉ khi có mục [`ccam folder`](/features/folder). Mục có tài khoản không tồn tại là lỗi, và mục có thư mục đã mất là cảnh báo. Cả hai đều không được sửa, vì xóa một mục sẽ thay đổi tài khoản được khởi động ở đó.
12. **Nhóm và cài đặt** (`checkGroups`). Một thành viên nhóm không phải là tài khoản, một tài khoản có tên về sau đã trở thành lệnh của ccam (khởi động nó bằng `ccam run <name>` hoặc đổi tên nó), một `theme` hoặc `usage_source` không xác định, và các dòng trong `config.toml` mà ccam không đọc đều là cảnh báo. Loại cuối thường là một cài đặt được đặt bên dưới bảng `[[group]]`, `[[profile]]` hoặc `[[folder]]`, khiến nó trở thành một phần của bảng đó (`config.toml line 43: relay is read as part of the [profile.login] table above it; settings for all accounts go above the first [[group]], [[profile]] or [[folder]]`); các lệnh thay đổi cấu hình sẽ từ chối cho đến khi nó được sửa. Không sửa.
13. **Đăng nhập** (`checkIdentities`). Đọc từ `.claude.json` của mỗi thư mục xem thư mục đó đang đăng nhập bằng ai. Một thư mục nay đăng nhập vào tài khoản khác với tài khoản mà `ccam login` đã ghi lại (ví dụ sau `/login` trong một phiên) là cảnh báo; một tài khoản không có bản ghi đăng nhập, chẳng hạn tài khoản được nhập hoặc được tạo bởi một phiên bản ccam cũ hơn, được so sánh theo email đã ghi lại. Một thư mục đăng nhập vào cùng tài khoản Claude với một thư mục khác cũng là cảnh báo: hai thư mục đó dùng chung giới hạn sử dụng. Một thư mục giữ cùng token đăng nhập với một thư mục khác, tức một thư mục được sao chép (kể cả `~/.claude` của `claude` thông thường), là lỗi, vì bên nào gia hạn trước sẽ đăng xuất bên kia; hãy chạy `ccam login <name>` để cho nó một token riêng. Không sửa.
14. **`auto_switch`** (`checkAutoSwitch`). Chỉ khi `auto_switch = true`. Một dòng thông tin trong hộp `General` cho biết `auto_switch` làm gì và liên kết tới [phần mô tả của nó](/features/next#auto-switch). Với mỗi gói có seat trong số các tài khoản (Claude Team, Claude Enterprise), một dòng nêu tên các tài khoản đó, cho biết `auto_switch` chuyển một phiên giữa các seat của cùng một tổ chức, trích quy tắc về seat của Anthropic kèm liên kết tới trung tâm trợ giúp, và để bạn tự quyết định. Các phiên bản Claude Code cũ hơn 2.1.214, vốn không cho hook biết tiến trình của chúng (`CLAUDE_PID`), là cảnh báo ở đó, vì không phiên nào có thể được chuyển.

    Một tài khoản không có plugin ccam, hoặc đã tắt plugin, nhận một cảnh báo trong hộp của chính nó kèm lệnh sửa trong Claude Code: `auto_switch is on, but work does not have the ccam plugin, so its sessions are not moved. In Claude Code as work: /plugin install ccam --marketplace jhyunwoo/ccam` (hoặc `/plugin enable ccam@ccam`). Các tài khoản API key bị bỏ qua. Trên Windows, dòng duy nhất là `auto_switch is ignored on Windows; with relay = true, /ccam:switch followed by /exit still moves a conversation to another account.` Không sửa.

Trước các tài khoản, một **lần thay đổi nhóm** bị dừng giữa chừng là lỗi trong hộp `General`: ``the last group change stopped halfway (ccam doctor --fix finishes it; noted in ~/.local/state/ccam/group-change.json)``. `--fix` hoàn tất nó theo đúng thứ tự đã lên kế hoạch, với các câu trả lời đã đưa ra cho các xung đột của nó; cho đến lúc đó, không lần thay đổi nhóm, đổi tên hay xóa nào khác được bắt đầu. Cho đến khi nó hoàn tất, ccam chỉ tạo các liên kết còn thiếu cho các tài khoản liên quan. Nếu không thể hoàn tất, lỗi sẽ cho biết rằng xóa tệp đó nghĩa là bỏ phần còn lại của lần thay đổi.

`ccam --dry-run doctor --fix` hiển thị thứ tự (`--fix finishes it in its order (b, a, c), with the answers given to its conflicts`) và những gì phần của từng tài khoản sẽ để riêng ra. Một lần thay đổi nhóm mà một tiến trình ccam khác đang thực hiện ngay lúc này là thông tin, và trong khi nó đang diễn ra, `--fix` không sửa liên kết dùng chung nào (`its share links are not repaired while another ccam changes the groups; …`). `--fix` giữ `~/.local/state/ccam/group-change.lock` trong lúc sửa, và đọc lại `config.toml` sau khi hoàn tất một lần thay đổi.

### Sửa liên kết dùng chung

Đường khởi động tạo ngay một liên kết còn thiếu, nhưng dừng lại với bất cứ thứ gì có thể làm mất dữ liệu và chỉ bạn đến đây:

```console
$ ccam work
ccam: stopped because of share link problems:
  /home/me/.local/share/ccam/profiles/work/CLAUDE.md: set to shared but is a real file, not a link
Repair them with `ccam doctor --fix`, then run again
```

Với bất kỳ vấn đề liên kết dùng chung nào, `--fix` gọi `share.Build` với thiết lập dùng chung của tài khoản làm cả trạng thái hiện tại lẫn trạng thái đích, rồi `Plan.Apply`. ccam suy ra lại mọi bước từ những gì có trên đĩa (`Build`), tạo ra cùng kế hoạch mà một thay đổi `ccam group` dùng, nên việc sửa tuân theo cùng các quy tắc như khi tham gia một nhóm.

<!-- motion:doctor-repair -->

Ở đây `settings.json` là một tệp thật. Vì nhóm đã có một `settings.json` khác, xung đột nhận lựa chọn mặc định là `keep group` (`KeepGroup` trong sơ đồ), nên bản của nhóm được giữ. Tệp của tài khoản được chuyển vào `.ccam-backup/<time>/settings.json` bên trong thư mục tài khoản, và ccam liên kết đường dẫn tới nhóm (`os.Symlink`). `--fix` thông báo điều đó: `settings.json differed from the group's copy, which it now uses; this account's copy is in …`. Các trường hợp khác:

- Nếu nhóm chưa có bản nào, tệp của tài khoản được chuyển vào nhóm và trở thành nội dung đầu tiên của nhóm.
- Với một thư mục, các tệp nhóm còn thiếu được chuyển vào nhóm. Các tệp giống hệt bản của nhóm sẽ bị xóa, còn các tệp khác nhau được đưa vào bản sao lưu và được nêu tên trong một ghi chú (`projects: 1 file differs from the group's copy (-w/memory/MEMORY.md); the group's is kept, the account's stays in …`).
- Thư mục nhóm bị thiếu sẽ được tạo trước.
- Với `history.jsonl`, mặc định là `merge by line`: các dòng của nhóm, rồi đến các dòng mới của bạn. Bản trước đó của nhóm được lưu trong `.ccam-backup/<time>/group/`.
- Liên kết tới một nhóm khác được thay thế. Liên kết tới một thư mục ngoài ccam được đưa vào bản sao lưu, và nội dung của nó trở thành nội dung ban đầu của nhóm nếu nhóm còn trống.
- Liên kết còn sót lại cho một mục nay đã là riêng sẽ bị gỡ bỏ, và tài khoản nhận một bản sao nội dung của nhóm.

Ngoài việc gỡ liên kết và xóa các tệp giống hệt bản của nhóm, mọi bước đều chuyển dữ liệu vào nhóm hoặc thư mục sao lưu. Không có gì bị mất.

## An toàn và các trường hợp đặc biệt

- `ccam doctor` thông thường không bao giờ thay đổi gì. Hãy dùng `--dry-run` cùng `--fix` để chắc chắn trong script.
- Với một tài khoản đang chạy, bất kể Claude Code được khởi động theo cách nào, `doctor --fix` chỉ tạo các liên kết còn thiếu và các thư mục nhóm còn thiếu, việc không làm mất gì. Khi vẫn còn việc khác, nó thêm `` is running; quit Claude Code there, then run `ccam doctor --fix` again to repair the rest ``, vì phiên vẫn tiếp tục ghi vào các tệp mà việc sửa chữa sẽ di chuyển.
- Nếu việc sửa liên kết dùng chung thất bại, ccam báo `could not repair share links: <error>` cho tài khoản đó và bỏ qua kiểm tra mục không xác định của nó. Các tài khoản khác vẫn được kiểm tra.
- Nếu một thay đổi `ccam group` thất bại giữa chừng, các nhóm mới đã được lưu và lần thay đổi được ghi lại trong `~/.local/state/ccam/group-change.json`; ccam nhắc bạn hoàn tất bằng `ccam doctor --fix`, lệnh này thực hiện phần còn lại của nó, kể cả các bước MCP.
- `.credentials.json` là liên kết tượng trưng, thông tin đăng nhập hết hạn và thông tin đăng nhập trong biến môi trường đều được báo nhưng không bao giờ bị thay đổi.
- Việc ghi lại tệp đăng ký plugin được thực hiện nguyên tử nhưng không được sao vào `.ccam-backup/`.

## Liên quan

- [Dùng chung máy chủ MCP, cài đặt và lịch sử](/docs/sharing)
- [Sao chép thiết lập hiện có](/features/migrate)
- [Giới hạn sử dụng và chuyển đổi tài khoản](/docs/usage-limits)
- [Tham chiếu lệnh](/docs/commands)
