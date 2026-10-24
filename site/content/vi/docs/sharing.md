---
title: Dùng chung máy chủ MCP, cài đặt và lịch sử
description: "Cách các nhóm của ccam dùng chung máy chủ MCP, settings.json, CLAUDE.md, skill và lịch sử hội thoại giữa các tài khoản Claude Code, giữ riêng mọi đăng nhập."
updated: 2026-10-23
order: 4
translation: machine
sourceHash: "27e0a0e5abbd"
---

Mỗi tài khoản giữ thông tin đăng nhập riêng. Mọi thứ khác đều có thể dùng chung qua **nhóm**: một nhóm nêu các mục nó dùng chung và các tài khoản thuộc nó, và ccam liên kết tệp của từng tài khoản đó đến thư mục của nhóm, nên thay đổi thực hiện ở một tài khoản sẽ xuất hiện ở mọi tài khoản.

## Những gì có thể dùng chung

| Mục | Nội dung được dùng chung |
| --- | --- |
| `mcp` | Máy chủ MCP: `mcp.json` của nhóm, được truyền cho Claude Code khi khởi động |
| `settings` | `settings.json`: quyền, hook, mô hình, dòng trạng thái |
| `memory` | `CLAUDE.md` toàn cục |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | lịch sử hội thoại: `projects/` và `history.jsonl` |

Tài khoản đầu tiên của bạn tạo nhóm `default`, nhóm này dùng chung mọi thứ trừ lịch sử hội thoại, và các tài khoản thêm sau sẽ tham gia nhóm đó.

Các tệp đăng nhập `.credentials.json` và `.claude.json` không bao giờ được dùng chung. ccam chỉ liên kết các đường dẫn nằm trong một danh sách cố định, và hai tệp này không có trong danh sách đó.

Trang này trình bày các tác vụ hằng ngày. Mọi tùy chọn được mô tả trên các trang tính năng: [dùng chung bằng nhóm](/features/share), [máy chủ MCP dùng chung](/features/mcp) và [sao chép thiết lập hiện có](/features/migrate).

## Nhóm

```sh
ccam group                                   # the groups, what they share, their accounts
ccam group set default --add history         # every account in default shares history
ccam group add team --items history --members work,side
```

Một tài khoản có thể thuộc nhiều nhóm, nhưng mỗi mục chỉ đến từ một nhóm. Vì vậy, để dùng chung lịch sử giữa `work` và `side` mà không dùng chung với các tài khoản khác của bạn, hãy để lịch sử ngoài `default` và cho hai tài khoản đó một nhóm riêng, như trên. `ccam share work` cho thấy một tài khoản dùng chung những gì và qua nhóm nào. [Bảng điều khiển](/features/dashboard) làm điều tương tự trên thẻ Groups.

## Khi một tài khoản bắt đầu dùng chung

Khi một tài khoản bắt đầu dùng chung một mục (nó tham gia một nhóm, hoặc nhóm của nó thêm mục đó), ccam sẽ:

1. in ra kế hoạch gồm mọi tệp nó sẽ di chuyển, sao chép hoặc liên kết cho từng tài khoản, và hỏi bạn trước khi bắt đầu;
2. chuyển các tệp của tài khoản vào thư mục nhóm nếu nhóm vẫn còn trống, hoặc thêm những tệp mà nhóm chưa có;
3. khi cùng một tệp có nội dung khác nhau, hỏi bạn muốn giữ phiên bản của nhóm, dùng phiên bản của bạn, hay gộp JSON;
4. giữ bản gốc của bạn trong `.ccam-backup/<time>/` bên trong thư mục tài khoản;
5. thay bản gốc bằng một liên kết đến nhóm.

Các tệp hội thoại được đặt tên theo ID phiên, nên lịch sử từ nhiều tài khoản được gộp lại mà không xung đột.

ccam không thay đổi một tài khoản khi Claude Code đang chạy bằng tài khoản đó, bất kể Claude Code được khởi động theo cách nào (kể cả phiên nền), vì phiên vẫn tiếp tục ghi vào các tệp đang được di chuyển. Hãy thoát phiên trước. ccam kiểm tra lại ngay trước khi lưu, nên một phiên được khởi động trong lúc kế hoạch đang chờ cũng bị phát hiện.

## Khi một tài khoản ngừng dùng chung

```sh
ccam group leave team work                   # work leaves the group
ccam group set default --remove history      # nobody in default shares history any more
```

Liên kết bị gỡ bỏ và tài khoản nhận một bản sao riêng của nội dung hiện tại trong nhóm; các máy chủ MCP của nhóm được đưa vào `.claude.json` của chính tài khoản, bên cạnh các máy chủ riêng của nó. Thêm `--fresh` để bắt đầu từ trạng thái trống thay vào đó; khi đó kế hoạch liệt kê các máy chủ MCP sẽ ngừng được tải. Một thư mục được nhập như `~/.claude` luôn giữ một bản sao. Nhóm vẫn giữ nội dung của nó cho các tài khoản khác.

## Sao chép thiết lập hiện có

Nếu bạn đã dùng Claude Code trước khi có ccam, cài đặt của bạn nằm trong `~/.claude`. `ccam migrate` sao chép chúng vào một nhóm, để mọi tài khoản dùng chung các mục đó đều có chúng:

```sh
ccam migrate
```

Trên terminal, ccam trước tiên hiển thị một danh sách kiểm (checklist) các mục. Mọi thứ trừ lịch sử hội thoại đều được đánh dấu; nhấn Space để thay đổi một mục và Enter để xem kế hoạch. Những mục mà thư mục không có gì được ghi là `not found`.

Những gì được sao chép:

- `settings.json` và `CLAUDE.md` toàn cục;
- các tiện ích mở rộng `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/` và `rules/`;
- máy chủ MCP người dùng của bạn, tức mục `mcpServers` trong `~/.claude.json`, vào `mcp.json` của nhóm;
- lịch sử hội thoại (`projects/` và `history.jsonl`), nhưng chỉ khi bạn thêm `--history`.

Thư mục nguồn chỉ được đọc: không có gì bên trong bị di chuyển, liên kết hay thay đổi, và lệnh `claude` thông thường vẫn hoạt động như trước. Các tệp đăng nhập `.credentials.json` và `.claude.json` không bao giờ được sao chép vào nhóm.

Các máy chủ MCP đăng nhập bằng OAuth giữ token trong thông tin đăng nhập, không phải trong danh sách máy chủ. ccam không sao chép các lần đăng nhập đó: một lần đăng nhập bị sao chép giữ cùng refresh token, và lần gia hạn đầu tiên sẽ đăng xuất mọi nơi khác đang giữ nó khỏi máy chủ đó, kể cả `~/.claude`. Kế hoạch và phần tóm tắt nêu tên các tài khoản cần tự đăng nhập vào các máy chủ bằng `/mcp`.

ccam in ra kế hoạch và hỏi trước khi sao chép bất cứ thứ gì:

```console
$ ccam migrate
Copy from /home/me/.claude into group default (/home/me/.local/share/ccam/groups/default):
  mcp         ~/.claude.json  new: context7, fs
  settings    settings.json   new: 80 B
  memory      CLAUDE.md       new: 31 B
  extensions  agents/         add 1 file the group lacks (17 B)
              commands/       add 1 file the group lacks (15 B)
              skills/         add 1 file the group lacks (32 B)
              not found, skipped: plugins/, output-styles/, rules/
  Conversation history (projects/, history.jsonl) is not included; add --history to copy it.
Copy these into group default? [y/N] y
✓ Copied into group default: 2 MCP servers, settings.json, CLAUDE.md, agents/, commands/, skills/
Group default: personal
To share them:
  ccam group join default work
```

Cuối cùng, ccam liệt kê các tài khoản của nhóm. Nếu nhóm chưa dùng chung một trong các mục, hoặc có tài khoản chưa thuộc nhóm, nó in ra các lệnh `ccam group` để thay đổi điều đó.

Khi nhóm đã có một bản khác của một tệp, ccam hỏi bạn, giống như khi bật dùng chung: giữ bản của nhóm hay dùng bản của bạn. `settings.json` và các tệp JSON khác còn có thể được hợp nhất, và `history.jsonl` có thể được hợp nhất theo từng dòng. Các bản trước đó của nhóm được lưu trong `.ccam-backup/<time>/` bên trong thư mục nhóm.

| Tùy chọn | Tác dụng |
| --- | --- |
| `--settings`, `--memory`, `--extensions`, `--mcp`, `--history` | Chỉ sao chép những mục bạn nêu, không có danh sách kiểm; nếu không có và ở ngoài terminal, sao chép mọi thứ trừ lịch sử hội thoại |
| `-g`, `--group <group>` | Sao chép vào nhóm này thay vì `default` |
| `-y`, `--yes` | Sao chép mà không hỏi; khi một tệp khác nhau, giữ bản của nhóm |

Đặt `--dry-run` trước lệnh, như `ccam --dry-run migrate`, để xem kế hoạch mà không sao chép gì. Để sao chép từ thư mục khác, hãy nêu tên thư mục đó: `ccam migrate ~/.claude-work --group work`.

## Máy chủ MCP

Claude Code lưu máy chủ MCP trong `.claude.json`, cùng với thông tin đăng nhập, nên tệp đó không thể dùng chung. ccam giữ các máy chủ dùng chung trong `mcp.json` của nhóm và khởi động Claude Code với `--mcp-config=<group>/mcp.json`. Các tùy chọn `--mcp-config` của riêng bạn vẫn được áp dụng thêm bên trên.

```sh
ccam mcp add fs -- npx -y @modelcontextprotocol/server-filesystem ~/code
ccam mcp list
```

Máy chủ bạn thêm trong một phiên bằng `claude mcp add` chỉ được lưu cho tài khoản đó, và cho thư mục bạn đang ở (phạm vi local của Claude Code; `-s user` làm nó thành của tài khoản trong mọi thư mục). Lần tới khi bạn khởi động tài khoản đó, ccam sẽ nêu tên máy chủ, cùng thư mục của nó, và lệnh để chuyển nó sang nhóm:

```sh
ccam mcp import work gh
```

Nếu nhóm đã có một máy chủ khác cùng tên, ccam sẽ bỏ qua và thông báo cho bạn. `--force` thay máy chủ của nhóm và giữ `mcp.json` cũ trong thư mục `.ccam-backup/` của nhóm.

Các máy chủ MCP đăng nhập bằng OAuth lưu token cùng với thông tin đăng nhập của tài khoản, vì vậy hãy cấp quyền cho chúng một lần trong mỗi tài khoản bằng `/mcp`. `ccam ext` và `ccam doctor` nêu tên các máy chủ mà một tài khoản chưa đăng nhập trong khi một tài khoản khác đã đăng nhập.

Nếu bạn muốn các máy chủ xuất hiện trong `claude mcp list` và `/mcp` như máy chủ người dùng thông thường, hãy đặt `mcp_strategy = "sync"` ở đầu `~/.config/ccam/config.toml`, phía trên `[[group]]`, `[[profile]]` hoặc `[[folder]]` đầu tiên. Khi đó ccam ghi chúng vào `.claude.json` của từng tài khoản trước khi khởi động, dưới khóa tệp của chính Claude Code và có bản sao lưu. Một máy chủ riêng của tài khoản có cùng tên vẫn là của tài khoản, và ccam sẽ thông báo điều đó.

## Hai tài khoản trong cùng một thư mục

Các tài khoản dùng chung lịch sử hội thoại sẽ ghi vào cùng các tệp. Nếu bạn khởi động một Claude Code thứ hai trong thư mục mà một tài khoản khác cùng nhóm đang chạy, ccam sẽ cảnh báo và hỏi bạn có muốn tiếp tục không.

## Khi có sự cố

```sh
ccam doctor
ccam doctor --fix
```

Nếu một đường dẫn dùng chung đã trở thành tệp thông thường, ví dụ sau khi bạn tự khôi phục một bản sao lưu, ccam sẽ dừng trước khi khởi động Claude Code và chỉ bạn đến `ccam doctor --fix`, lệnh này liên kết nó với nhóm trở lại: tệp trở thành bản của nhóm khi nhóm chưa có, còn nếu không thì nhóm giữ bản của mình và bản của bạn được đưa vào một bản sao lưu, mà `--fix` sẽ nêu tên. Một lần thay đổi nhóm bị dừng giữa chừng cũng được `ccam doctor --fix` hoàn tất.
