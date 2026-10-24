---
title: Tham chiếu lệnh
metaTitle: "Tham chiếu lệnh ccam"
description: "Mọi lệnh và tùy chọn của ccam, từ ccam add và ccam <name> đến ccam usage, ccam next và ccam doctor."
updated: 2026-10-24
order: 3
translation: machine
sourceHash: "36af7e209cd9"
---

Đặt `--dry-run` trước một lệnh để xem lệnh đó sẽ làm gì mà không thực sự thực hiện, ví dụ `ccam --dry-run work`. Tùy chọn này dùng được khi khởi động Claude Code (`ccam <name>`, `ccam run`, `ccam next`, `ccam attach`) và với `group`, `migrate`, `shortcut`, `update` và `doctor --fix`. Các lệnh khác có thay đổi một thứ gì đó sẽ từ chối tùy chọn này thay vì chạy.

<!-- commands -->

## Khởi động Claude Code

`ccam <name|number> [arguments]` khởi động Claude Code bằng tài khoản đó. ccam không đọc các đối số sau tên tài khoản; nó chuyển chúng cho `claude` theo đúng thứ tự, và đầu vào chuẩn, mã thoát cũng như tín hiệu hoạt động giống hệt như khi chạy chính `claude`. Dấu `--` sau tên tài khoản là không bắt buộc.

`ccam` có đối số nhưng không có tên tài khoản, chẳng hạn `ccam -p "summarise"`, sẽ dùng tài khoản được đặt cho thư mục hiện tại, nếu không thì dùng tài khoản mặc định. Đặt chúng bằng `ccam folder <name> [folder]` và `ccam use <name>`.

```sh
ccam folder work ~/code/company   # work inside ~/code/company and below
ccam folder                       # list folders and the account used here
ccam folder --remove ~/code/company
ccam --fallback work -p "fix the tests"   # run again as another account at the limit
```

`--fallback` đặt trước tên tài khoản dùng cho các lần chạy `-p`/`--print`. Khi tài khoản chạm giới hạn sử dụng, ccam chạy lại cùng lệnh đó bằng tài khoản đã đăng nhập còn nhiều hạn mức nhất và chỉ in đầu ra của lần chạy đã hoàn tất. Xem [Tài khoản theo thư mục](/features/folder) và [Chạy lại khi chạm giới hạn](/features/next#run-again-at-the-limit).

**Seat của gói Team và Enterprise.** Lựa chọn tự động (`ccam next`, `--fallback`, `ccam pick`, phím `n` trên bảng điều khiển) cũng chuyển đổi giữa các seat của cùng một tổ chức Claude Team hoặc Enterprise. Trung tâm trợ giúp của Anthropic ghi rằng "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (mỗi seat của gói Team được gán cho đích danh một cá nhân, và chỉ người đó được phép sử dụng) và "Users can't route their work through another user's account to use that user's remaining usage" (người dùng không được chuyển công việc của mình qua tài khoản của một người dùng khác để dùng phần hạn mức còn lại của người đó) ([gói Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); tương tự đối với [gói Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Vì vậy, dùng lựa chọn tự động giữa các seat có thể vi phạm điều khoản của Anthropic. Hãy xác nhận với tổ chức của bạn và Anthropic trước khi dựa vào nó. Các gói cá nhân (Pro, Max) được đề cập trong [ghi chú về điều khoản ở phần Câu hỏi thường gặp](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

Một lệnh tắt như `ccwork` khởi động tài khoản của nó theo cùng cách: `ccwork --resume` làm đúng những gì `ccam work --resume` làm. Xem [lệnh tắt](/docs/commands#shortcuts).

## Thêm và xóa tài khoản

```sh
ccam add                      # wizard
ccam add work --group default --group team
ccam add side --shortcut ccside
ccam add ci --no-login        # log in later with: ccam login ci
ccam login --expired          # log in every account whose login no longer works
ccam logout side              # asks first when MCP sign-ins go with the login
ccam rename work job
ccam remove side              # logs out, then the folder moves to ~/.local/share/ccam/trash
ccam import ~/.claude-old --name old
ccam add gw --api --base-url https://llm.example.com --auth bearer --key-from env:GW_KEY
```

`ccam import` không kèm thư mục sẽ tìm `~/.claude` và `~/.claude-*` rồi hỏi bạn về từng thư mục; thư mục nào thất bại sẽ được báo là `Skipped <path>: <reason>` và các thư mục còn lại vẫn được đăng ký, còn thư mục có tên trùng với một tài khoản đã có thì cần `--name`, được đưa ra cùng với thư mục. Các thư mục được nhập sẽ được đăng ký ngay tại chỗ; ccam không bao giờ di chuyển hay sao chép thông tin đăng nhập.

`ccam login --expired` cũng đăng nhập các tài khoản mà Claude Code đã xóa thông tin đăng nhập. `ccam remove` đăng xuất một tài khoản được quản lý trước khi thư mục của nó vào thùng rác (trừ khi một tài khoản khác đang giữ cùng thông tin đăng nhập đó), và `ccam remove` với một tài khoản được nhập sẽ để lại cho các nhóm của nó một bản sao của những gì chúng dùng chung. Trên macOS, `ccam rename` giữ nguyên tên thư mục, vì Keychain lưu thông tin đăng nhập theo đường dẫn của thư mục.

`ccam add --api` thêm một tài khoản dùng API key thay cho đăng nhập Claude; ccam đọc khóa mỗi lần khởi động và không bao giờ lưu nó. Xem [Tài khoản API key](/features/api).

`ccam rename`, `ccam remove` và các lệnh `ccam group` từ chối thay đổi một tài khoản khi Claude Code đang chạy bằng tài khoản đó, bất kể nó được khởi động theo cách nào (kể cả phiên nền). Hãy thoát phiên đó trước. Chúng cũng từ chối khi một lần thay đổi nhóm còn dang dở: ``the last change of groups stopped halfway; finish it with `ccam doctor --fix` first``.

`ccam list --json` in ra mọi tài khoản cùng với slot, lệnh tắt, email và thư mục của nó. Mảng `groups` trong đó liệt kê các nhóm của tài khoản, còn đối tượng `share` ánh xạ từng mục (`mcp`, `settings`, `memory`, `extensions`, `history`) tới nhóm mà mục đó đến từ, hoặc `local`.

## Lệnh tắt

```sh
ccam shortcut work ccwork     # create it, or rename the existing one
ccam shortcut work            # show it
ccam shortcut work --remove
```

Lệnh tắt là một lệnh do bạn tự đặt tên, dùng để khởi động Claude Code bằng một tài khoản. Giống như `ccam work`, nó chuyển nguyên vẹn mọi đối số cho Claude Code, nên `ccwork --resume` và `ccwork -p "hi"` hoạt động đúng như bạn mong đợi. Nó là một liên kết đến ccam trong `~/.local/bin`.

Lệnh tắt là tùy chọn. Trình hướng dẫn `ccam add` hỏi bạn có muốn tạo không, `ccam add <name> --shortcut <command>` tạo luôn mà không hỏi, và `ccam list` hiển thị lệnh tắt của từng tài khoản, hoặc `-` khi không có.

Tên dùng các ký tự a-z, 0-9, `.`, `-` và `_`. Tên không được trùng với một lệnh của ccam, lệnh tắt của tài khoản khác, hoặc tên của một chương trình khác trong `PATH`, cũng không được là một tệp trong `~/.local/bin` mà ccam không tạo ra, chẳng hạn bản sao hay liên kết ccam của riêng bạn. Một lệnh tắt được sửa tay trong `config.toml` mà vi phạm các quy tắc này, là `claude` hoặc `ccam`, hoặc được gán cho hai tài khoản sẽ khiến `config.toml` không tải được, kèm tên của tài khoản đó.

`ccam rename` giữ nguyên tên lệnh tắt, còn `ccam remove` xóa lệnh tắt cùng với tài khoản. `ccam shims install` tạo lại liên kết của mọi tài khoản, chẳng hạn sau khi bạn di chuyển ccam; một tệp chắn đường của một lệnh tắt không ngăn các lệnh tắt khác.

## Dùng chung

```sh
ccam group                                  # the groups, what they share, their accounts
ccam group add team --items history --members work,side
ccam group set default --add history        # every account in default changes together
ccam group join team personal
ccam group leave team side --fresh
ccam share work                             # what work shares, and through which group
```

Một nhóm dùng chung các mục giữa các tài khoản của nó; một tài khoản có thể thuộc nhiều nhóm, nhưng mỗi mục chỉ đến từ một nhóm. Xem [dùng chung bằng nhóm](/features/share).

## Sao chép thiết lập hiện có

```sh
ccam migrate                    # ~/.claude into the default group
ccam --dry-run migrate          # show the plan, copy nothing
ccam migrate ~/.claude-work --group work --settings --mcp
ccam migrate --history --yes    # conversation history only
```

`ccam migrate` sao chép `settings.json`, `CLAUDE.md`, tiện ích mở rộng và máy chủ MCP người dùng từ một thư mục Claude Code vào một nhóm, để mọi tài khoản dùng chung các mục đó đều có chúng. Thư mục chỉ được đọc. Các tùy chọn `--settings`, `--memory`, `--extensions`, `--mcp` và `--history` chỉ sao chép những mục bạn nêu. Nếu không có chúng, trên terminal bạn chọn trong một danh sách kiểm (checklist), với mọi thứ trừ lịch sử hội thoại đã được đánh dấu; trong các trường hợp khác, mọi thứ trừ lịch sử hội thoại đều được sao chép. `-g` hoặc `--group` chọn nhóm, còn `-y` hoặc `--yes` sao chép mà không hỏi và giữ bản của nhóm khi một tệp khác nhau. Các lần đăng nhập MCP không được sao chép, vì một lần đăng nhập bị sao chép sẽ đăng xuất những nơi khác đang giữ nó; kế hoạch nêu tên các tài khoản cần đăng nhập bằng `/mcp`. `--mcp` cũng hoạt động khi `~/.claude` trống hoặc không có nhưng `~/.claude.json` có máy chủ. Xem [sao chép thiết lập hiện có](/docs/sharing#copy-an-existing-setup).

## Máy chủ MCP

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp add docs --url https://mcp.example.com/mcp --header "Authorization: Bearer $TOKEN"
ccam mcp add-json db '{"type":"stdio","command":"db-mcp"}'
ccam mcp list
ccam mcp import work gh
ccam mcp remove fs
```

`ccam mcp import <account>` chuyển các máy chủ chỉ tài khoản đó có vào nhóm của nó, kể cả những máy chủ nó chỉ có cho một thư mục (`claude mcp add` thông thường), và sau đó chúng được tải trong mọi thư mục. Máy chủ có tên mà nhóm đã dùng cho một định nghĩa khác sẽ bị bỏ qua; thêm `--force` để thay bản của nhóm, ccam sẽ sao lưu bản đó trước.

## Mức sử dụng và chuyển đổi tài khoản

```sh
ccam usage           # also: ccam u
ccam usage --watch   # refresh every 60 seconds
ccam usage --json
ccam usage --source claude      # Claude Code looks it up; ccam never touches a token
ccam next --continue            # this folder's last conversation, as the next account
ccam next --auto
ccam next --auto --max-usage 80 --wait
ccam next -c --auto -- --dangerously-skip-permissions   # claude's own options after --
ccam pick                       # the account with the most room, for a script
ccam "$(ccam pick)" -p "…"      # start a worker on it
ccam agents                     # background sessions of every account
ccam attach 49ec0028            # open one, as the account that started it (or by its name)
```

`next --auto`, `--fallback` và `ccam pick` chọn theo cùng một cách: ưu tiên các tài khoản dưới `auto_max_usage` (`config.toml`, mặc định 100), theo `auto_strategy` (`headroom`, `expiring` hoặc `balance`), có tính khoảng hàng tuần của `--model`, không bao giờ chọn tài khoản ở mức 100% (một khoảng mà Claude Code báo là đã đầy được tính là 100% cho đến khi đặt lại), không bao giờ chọn đăng nhập bằng API key, và chỉ trong tổ chức của tài khoản xuất phát (`--cross-org` hoặc `cross_org = true` để tìm rộng hơn). `--wait` nhận giá trị sau dấu `=` (`--wait=2h`); viết với dấu cách thì bị từ chối kèm gợi ý. Khi không chọn được tài khoản nào, chúng thoát với mã 75, hoặc với `--wait` thì chờ tài khoản đầu tiên được đặt lại. `ccam next --continue` chuyển một cuộc hội thoại sang ngay cả giữa các tài khoản không dùng chung lịch sử: khi đó Claude Code sao chép nó vào tài khoản mới. Khác với `ccam <name>`, `ccam next` tự đọc các tùy chọn của nó, nên các tùy chọn của Claude Code đặt sau `--`; thiếu nó, ccam sẽ dừng với `unknown flag`. Bên trong Claude Code, `/ccam:switch [account]` của plugin chuyển một cuộc hội thoại đi tiếp, và `ccam hook` là lệnh mà các hook của plugin chạy. Xem [Chuyển sang tài khoản tiếp theo](/features/next), [Chọn tài khoản cho script](/features/pick) và [Phiên chạy nền của mọi tài khoản](/features/agents).

## Bảo trì

```sh
ccam doctor          # check
ccam doctor --fix    # repair what is safe to repair, with backups
ccam statusline install
ccam ext             # skills, commands, MCP servers and hooks of every account, and problems
ccam shims install   # recreate the shortcut links
ccam update          # install the latest ccam release
ccam update --check  # only check for a newer release
ccam version         # also: ccam -v, ccam --version
ccam completion bash # or zsh, powershell
```

`ccam update` kiểm tra mã băm SHA-256 của bản phát hành trước khi thay thế ccam, còn `--version v0.2.0` cài một bản phát hành cụ thể. ccam được cài từ gói `.deb` thì cập nhật bằng apt. Xem [cập nhật](/docs/install#update).

## Tệp

| Đường dẫn | Nội dung |
| --- | --- |
| `~/.config/ccam/config.toml` | Các tài khoản, nhóm và cài đặt (`theme`, `usage_source`, `relay`, `cross_org` …). Cài đặt đặt ở đầu tệp, phía trên `[[group]]`, `[[profile]]` hoặc `[[folder]]` đầu tiên: cài đặt nằm dưới một bảng thuộc về bảng đó, nên ccam nêu tên nó mỗi lần khởi động và trong `ccam doctor`, và các lệnh lưu cấu hình sẽ từ chối cho đến khi nó được dời lên. Chú thích không được giữ lại khi ccam lưu tệp |
| `~/.local/share/ccam/profiles/<name>/` | `CLAUDE_CONFIG_DIR` của từng tài khoản |
| `~/.local/share/ccam/groups/<group>/` | Cài đặt, máy chủ MCP, tiện ích mở rộng và lịch sử dùng chung |
| `~/.local/state/ccam/` | Bộ nhớ đệm mức sử dụng và bản ghi giới hạn (`usage/<name>.json`, `usage/<name>.limit.json`), thời điểm dùng gần nhất, dấu phiên và khóa (`locks/`), các lệnh tắt mà ccam đã tạo (`shortcuts.json`), một lần thay đổi nhóm chưa xong (`group-change.json`), và các ghi chú mà `/ccam:switch` và hook giới hạn để lại (`handoff/<account>/<session>.json`) |
| `~/.local/bin/<shortcut>` | Các lệnh tắt như `ccwork`, là liên kết đến ccam |

Các biến `XDG_CONFIG_HOME`, `XDG_DATA_HOME` và `XDG_STATE_HOME` dùng để đổi vị trí của các thư mục này. Trên Windows, `~` là `%USERPROFILE%` và các lệnh tắt là `<shortcut>.exe`.
