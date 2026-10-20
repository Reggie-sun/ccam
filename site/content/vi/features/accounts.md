---
title: Quản lý tài khoản
description: "Thêm, đăng nhập, liệt kê, đổi tên, xóa và nhập tài khoản Claude Code bằng ccam add, ccam login, ccam list, ccam rename, ccam remove và ccam import."
updated: 2026-10-20
order: 2
translation: machine
sourceHash: 9eb992214132
---

Mỗi tài khoản ccam là một thư mục cấu hình Claude Code riêng với thông tin đăng nhập riêng. Các lệnh này tạo, đăng nhập, liệt kê, đổi tên và xóa tài khoản, cũng như đăng ký các thư mục bạn đã có.

## Cách dùng

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

Ở bất kỳ lệnh nào nhận `<name|number>`, bạn có thể đưa tên tài khoản hoặc số slot của nó, tức cột `#` của `ccam list`.

`ccam add` chạy trên terminal mà không có tùy chọn nào sẽ mở một trình hướng dẫn. Nó hỏi tên, rồi hiển thị một danh sách kiểm (checklist) các mục dùng chung (mọi thứ trừ lịch sử hội thoại đã được đánh dấu), rồi hỏi `Create a shortcut command for work? [y/N]` và, nếu bạn trả lời `y`, hỏi tiếp `Shortcut name [ccwork]:`. Bất kỳ tùy chọn nào trong `--share`, `--shortcut` hoặc `--no-login`, hoặc một pipe thay cho terminal, đều bỏ qua trình hướng dẫn; khi đó bắt buộc phải có tên và không có lệnh tắt nào được tạo trừ khi `--shortcut` nêu tên một lệnh.

`ccam list` hiển thị slot, tên, lệnh tắt, email, các mục dùng chung và lần dùng gần nhất, với `*` đứng trước tài khoản mặc định:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

Với `--json`, mỗi tài khoản là một đối tượng gồm `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (bị bỏ đi khi tài khoản chưa từng được khởi động) và `share`. `share` ánh xạ từng mục (`mcp`, `settings`, `memory`, `extensions`, `history`) tới tên nhóm của nó hoặc tới `"local"`.

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `ccam add --no-login` | Tạo tài khoản mà không đăng nhập; đăng nhập sau bằng `ccam login <name>` | tắt |
| `ccam add --share <item=group>` | Dùng chung một mục với một nhóm, hoặc `item=local` để giữ riêng. Có thể lặp lại | `mcp`, `settings`, `memory`, `extensions` trong `default`; `history` riêng |
| `ccam add --shortcut <command>` | Tạo thêm lệnh tắt này trong `~/.local/bin` | không có |
| `ccam add --api --key-from <source>` | Tài khoản dùng API key thay cho đăng nhập; xem [Tài khoản API key](/features/api) | tắt |
| `ccam list --json` | In các tài khoản dưới dạng JSON | tắt |
| `ccam remove -y`, `--yes` | Xóa mà không hỏi | tắt |
| `ccam import --name <name>` | Tên tài khoản cho thư mục bạn đưa vào | lấy từ tên thư mục |
| `ccam import -y`, `--yes` | Nhập mọi thư mục tìm thấy mà không hỏi | tắt |
| `--dry-run` | Được `ccam list` và `ccam completion` hỗ trợ, vì chúng không thay đổi gì. Bị `add`, `login`, `logout`, `use`, `rename`, `remove` và `import` từ chối, với thông báo `` --dry-run is not supported by `ccam add`; nothing was changed `` (nêu tên lệnh) | tắt |

`ccam login`, `ccam logout`, `ccam use` và `ccam rename` không có tùy chọn riêng. `ccam completion` nhận `bash`, `zsh` hoặc `powershell`; giá trị khác sẽ thất bại với `only bash, zsh and powershell are supported`.

Tên tài khoản dài từ 1 đến 32 ký tự gồm `a-z`, `0-9`, `-` và `_`, bắt đầu bằng một chữ cái hoặc chữ số. Tên không được chỉ gồm chữ số, vì `ccam 2` nghĩa là slot 2. Tên cũng không được trùng với một lệnh của ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`, `folder`, `pick`, `ext`), nên `ccam <name>` không bao giờ mơ hồ. Danh sách này là `config.Reserved` trong `internal/config`.

## Cách hoạt động

`ccam add work --shortcut ccwork` chạy các bước sau trong `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. `config.ValidateName` kiểm tra tên, và `shim.Check` kiểm tra lệnh tắt trước khi bất cứ thứ gì được tạo.
2. `profile.Add` tạo `~/.local/share/ccam/profiles/work/` với quyền `0700` và gán cho tài khoản số slot trống thấp nhất (`Config.NextSlot`), nên số của một tài khoản đã xóa sẽ được dùng lại. Tài khoản đầu tiên trở thành tài khoản mặc định.
3. `applyShare` liên kết thư mục mới, còn trống vào các nhóm dùng chung của nó bằng `share.Build` rồi `Apply`, ví dụ `settings.json` tới `groups/default/settings.json`. Nếu bước này thất bại, thư mục mới bị xóa, vì một thư mục còn sót lại sẽ chặn mọi lần thử lại với `folder already exists`.
4. Tài khoản được lưu vào `~/.config/ccam/config.toml` qua `fsutil.AtomicWrite`: một tệp tạm, `fsync`, rồi `rename`, nên một sự cố không bao giờ để lại cấu hình dở dang.
5. `runLogin` khởi động `claude auth login` với `CLAUDE_CONFIG_DIR` trỏ tới thư mục tài khoản. Nếu đăng nhập thất bại, tài khoản vẫn được giữ, và ccam in ra lỗi cùng lời nhắc chạy `ccam login work` sau.
6. `profile.ReadEmail` lấy `oauthAccount.emailAddress` ra khỏi `profiles/work/.claude.json`. Tệp đó có thể nặng vài megabyte, nên `TopLevelValue` chỉ quét tìm một khóa đó thay vì phân tích cả tệp.
7. `shim.Install` liên kết `~/.local/bin/ccwork` tới tệp thực thi ccam, và khi khởi động tệp này tra `argv[0]` trong `config.toml`.

Việc đăng nhập có thể mất vài phút, và trong lúc đó một ccam khác có thể lưu cấu hình. Vì vậy các lệnh phải chờ sẽ lưu qua `env.update`, hàm này tải lại `config.toml`, áp dụng đúng một thay đổi (ở đây là email) rồi lưu. `ccam remove` cũng tải lại sau câu hỏi xác nhận vì lý do tương tự.

`ccam login` chạy bước 5 và 6 cho một tài khoản đã có. `ccam logout` chạy `claude auth logout` trong tài khoản và xóa email đã lưu. `ccam use` chỉ ghi lại `default` trong `config.toml`.

### Đăng nhập trên máy từ xa

Qua SSH, phần trình duyệt của `claude auth login` có thể không hoạt động. Claude Code chờ việc đăng nhập trên một máy chủ callback tại `http://localhost:<port>/callback` trên máy chủ, nhưng trình duyệt hoàn tất việc đăng nhập lại chạy trên máy tính của bạn, nơi địa chỉ đó không dẫn đến đâu và trình duyệt hiển thị lỗi như `ERR_CONNECTION_REFUSED`. Địa chỉ đó vẫn mang mã dùng một lần, nên ccam có thể chuyển nó đi:

<!-- motion:accounts-relay -->

`runLogin` không giao terminal của bạn cho Claude Code. Nó nối đầu vào chuẩn của Claude Code với một pipe và `relayLogin` sao chép từng dòng bạn gõ vào đó, nên mã được dán ở lời nhắc `Paste code here` của Claude Code vẫn đến nơi. Khi một dòng là địa chỉ callback, `callbackURL` chỉ chấp nhận nó nếu máy chủ là `localhost`, `127.0.0.1` hoặc `[::1]`, có cổng, đường dẫn là `/callback` và có cả `code` lẫn `state`; có thể bỏ `http://`. Sau đó `deliverCallback` gửi yêu cầu tới địa chỉ đó trên máy chủ với giới hạn 10 giây và không đi theo chuyển hướng, giống như trình duyệt lẽ ra đã làm. Yêu cầu đến máy chủ callback của `claude auth login` đang chờ (đường nét đứt trong sơ đồ). Claude Code nhận mã và lưu thông tin đăng nhập, trên Linux là trong `.credentials.json` ở thư mục tài khoản. Một dòng được dán vào không bao giờ khiến ccam liên lạc với máy chủ khác.

ccam in ra dòng này trước khi bắt đầu đăng nhập:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

Khi thành công, nó in ra `→ Passed the sign-in to Claude Code on this machine.`

Trong lúc việc đăng nhập chờ nhập liệu trên terminal, ccam tắt chế độ echo và chế độ dòng của terminal, và hiển thị một `*` cho mỗi ký tự bạn gõ, nên mã hoặc địa chỉ callback được dán vào sẽ không hiện trên màn hình. `readLine` tự xử lý việc chỉnh sửa dòng: `Enter` gửi dòng, `Backspace` xóa một ký tự, và `Ctrl+D` trên một dòng trống kết thúc việc nhập. Terminal được khôi phục khi đăng nhập kết thúc, và cả khi nhận `Ctrl+C`, `SIGTERM` hoặc `SIGHUP`, những tín hiệu vẫn kết thúc ccam như trước. Khi đầu vào chuẩn và đầu ra lỗi chuẩn không phải là terminal, các dòng được sao chép nguyên như khi chúng đến.

### Đổi tên, xóa và nhập

`ccam rename` và `ccam remove` tác động lên thư mục của tài khoản và lên những gì ccam lưu ở nơi khác dưới tên tài khoản. `nameState` liệt kê các đường dẫn đó: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` và các bộ nhớ đệm mức sử dụng `~/.local/state/ccam/usage/<name>.json` và `<name>.statusline.json`. `ccam remove work` xử lý chúng theo thứ tự sau:

<!-- motion:accounts-remove -->

`ccam remove` trước tiên từ chối nếu tài khoản đang chạy, rồi mới hỏi. Nếu tài khoản dùng chung tiện ích mở rộng, `RewritePluginPaths` ghi lại `installed_plugins.json` và `known_marketplaces.json` của nhóm, vì Claude Code lưu đường dẫn plugin tuyệt đối có thể đi qua thư mục của tài khoản này. `profile.Trash` chuyển thư mục vào `~/.local/share/ccam/trash/<name>-<time>/`, tài khoản bị gỡ khỏi `config.toml`, các đường dẫn `nameState` bị xóa để một tài khoản cùng tên sau này không thừa hưởng chúng, và liên kết lệnh tắt bị gỡ bỏ. Để khôi phục, hãy chuyển thư mục từ thùng rác về chỗ cũ.

`ccam rename` đổi tên thư mục, ghi lại các đường dẫn plugin trong `plugins/` của nó từ thư mục cũ sang thư mục mới, giữ nguyên trạng thái mặc định, và chuyển các đường dẫn `nameState` sang tên mới. Lệnh tắt giữ nguyên tên. Nếu việc ghi lại plugin thất bại, thư mục được đổi lại tên cũ, vì cấu hình chưa được lưu.

`ccam import` đăng ký một thư mục ngay tại chỗ, dưới dạng `dir` trong `config.toml`. Nó không bao giờ di chuyển hay sao chép thư mục: một refresh token bị sao chép sẽ khiến các bản sao đăng xuất lẫn nhau. Không có đối số thư mục, nó tìm các thư mục `~/.claude` và `~/.claude-*` chưa được đăng ký và hỏi về từng thư mục; tên được lấy từ thư mục (`~/.claude` là `main`, `~/.claude-work` là `work`, `~/.claude-2` là `acct2`), và thư mục không suy ra được tên sẽ được liệt kê kèm lệnh `--name` cần dùng. Khi không có `CLAUDE_CONFIG_DIR`, Claude Code giữ `.claude.json` trong thư mục home của bạn, cạnh `~/.claude`, nên khi bạn nhập `~/.claude` mà nó không có `.claude.json` riêng, ccam tạo liên kết `~/.claude/.claude.json` tới `~/.claude.json`. Khi đó `claude` thông thường và `ccam main` thấy cùng một tài khoản. Tài khoản được nhập không dùng chung gì cho đến khi bạn chạy `ccam share`, và không có lệnh tắt nào được tạo.

`ccam completion bash|zsh` in ra một script tự hoàn thành (completion), script này cũng hoàn thành cả tên tài khoản.

## An toàn và các trường hợp đặc biệt

- `ccam rename` và `ccam remove` từ chối khi tài khoản đang chạy (`refuseRunning`), với `work is running; quit its Claude Code session first`, vì phiên vẫn tiếp tục ghi vào thư mục.
- `ccam remove` trên một tài khoản đã nhập chỉ hủy đăng ký nó: `✓ main unregistered (folder ... left in place)`. `~/.claude` của bạn không bị chuyển vào thùng rác.
- `ccam add` từ chối một tên mà thư mục của nó đã tồn tại trong `profiles/`, và gợi ý dùng `ccam import`.
- `ccam import` từ chối một thư mục đã được đăng ký, một đường dẫn không phải thư mục, và một tên đã có người dùng.
- `ccam logout` và `ccam login` không kiểm tra tài khoản có đang chạy hay không.
- Nếu phần chuyển tiếp không liên lạc được với Claude Code, chẳng hạn sau khi việc đăng nhập đã kết thúc, nó in ra `ccam: could not pass the sign-in to Claude Code:` kèm lý do. Lỗi HTTP thường có nghĩa là mã đã hết hạn; hãy bắt đầu đăng nhập lại.

## Liên quan

- [Khởi động Claude Code](/features/launch)
- [Lệnh tắt](/features/shortcuts)
- [Dùng chung máy chủ MCP, cài đặt và lịch sử](/docs/sharing)
- [Nhiều tài khoản Claude Code](/guides/multiple-claude-code-accounts)
