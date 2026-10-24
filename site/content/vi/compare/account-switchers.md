---
title: ccam và các công cụ đổi tài khoản Claude Code khác
description: "So sánh ccam với claude-swap, CCS, clausona, clauth, CAAM, cux và các công cụ đổi tài khoản Claude Code khác, và cách tiếp cận nào hợp với nhu cầu nào."
updated: 2026-10-24
order: 2
translation: machine
sourceHash: "6d08203f1ea6"
---

Các công cụ đổi tài khoản Claude Code đi theo một trong hai cách. Phần lớn hoán đổi thông tin đăng nhập bên trong một thư mục `~/.claude` duy nhất, nên mỗi lúc chỉ có một tài khoản hoạt động. Một số khác, trong đó có ccam, cho mỗi tài khoản một thư mục riêng bằng `CLAUDE_CONFIG_DIR`, nên mọi tài khoản đều giữ trạng thái đăng nhập và có thể chạy nhiều tài khoản cùng lúc. ccam còn bổ sung việc dùng chung cài đặt và máy chủ MCP, và tiếp tục một cuộc hội thoại trong thư mục của tài khoản khác.

## Hoán đổi một đăng nhập hay dùng thư mục riêng?

**Hoán đổi đăng nhập** chỉ giữ một bộ cài đặt, máy chủ MCP và lịch sử, vì chỉ có một thư mục. Việc chuyển đổi nhanh, và một phiên đang chạy có thể nhận đăng nhập mới. Nhưng mỗi lúc chỉ có một tài khoản hoạt động, cho mọi terminal và trình soạn thảo cùng lúc.

Để hoán đổi, công cụ giữ một bản sao thông tin đăng nhập của mọi tài khoản khác và ghi bản đó vào đúng chỗ. Mỗi lần gia hạn đăng nhập, Claude Code lại thay refresh token bên trong, nên mỗi bản sao đã lưu phải luôn được giữ khớp với các phiên đang chạy. Tài khoản bị đăng xuất hoặc phải `/login` lại sau khi chuyển là một trong những lỗi mà các dự án này phải sửa thường xuyên nhất.

**Thư mục riêng** là cách mà [tài liệu của Claude Code mô tả](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Mọi tài khoản đều giữ trạng thái đăng nhập, và bạn có thể chạy song song một phiên công việc và một phiên cá nhân. Thông tin đăng nhập của mỗi tài khoản nằm yên trong thư mục riêng, nơi Claude Code gia hạn nó như bình thường. Cái giá phải trả là mỗi thư mục bắt đầu với cài đặt, máy chủ MCP và lịch sử riêng, trừ khi có công cụ dùng chung chúng.

## So sánh các công cụ

Tổng hợp từ tài liệu và mã nguồn của từng dự án vào ngày 2026-10-20 (clausona: v0.12.0-beta). "Không có trong tài liệu" nghĩa là dự án không nói đến, không có nghĩa là thiếu tính năng.

| Công cụ | Cách tiếp cận | Nhiều tài khoản cùng lúc | Mức sử dụng 5 giờ và hàng tuần | Khi chạm giới hạn | Nền tảng |
| --- | --- | --- | --- | --- | --- |
| ccam | Thư mục riêng | Có | Có, `ccam usage` | Một thông báo trong phiên; `ccam next --auto --continue` hoặc `/ccam:switch` chuyển cuộc hội thoại sang tài khoản còn nhiều hạn mức nhất, dù có dùng chung lịch sử hay không; `ccam --fallback` tự chạy lại các tác vụ `-p` ở đó, và `ccam pick` phân bổ các worker | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | Hoán đổi đăng nhập; thư mục riêng với `cswap run` | Có, với `cswap run` | Có | `cswap auto` đổi đăng nhập khi chạm ngưỡng; phiên đang chạy vẫn tiếp tục | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Thư mục riêng | Có | Một phần | Nhóm tài khoản chuyển yêu cầu sang tài khoản kế tiếp qua một proxy cục bộ | Linux, macOS, Windows |
| [clausona](https://github.com/larcane97/clausona) | Thư mục riêng | Có | Có | Các route chọn một tài khoản còn hạn mức trước mỗi lần chạy, có thể chờ đặt lại, và giao các tài khoản khác nhau cho các worker | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Hoán đổi đăng nhập; thư mục riêng với `clauth start` | Có, với `clauth start` | Có | Chuyển lần lượt theo chuỗi tài khoản dự phòng | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Hoán đổi đăng nhập, hoặc hồ sơ tách biệt | Có, với hồ sơ tách biệt | Có | `caam run` thử lại các lệnh không tương tác trên tài khoản khác | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Hoán đổi đăng nhập và bọc lệnh `claude` | Không, chỉ một tài khoản hoạt động | Có | Hoán đổi đăng nhập và tiếp tục cùng cuộc hội thoại | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Hoán đổi đăng nhập | Không | Có | Tùy chọn đổi khi chạm ngưỡng | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Hoán đổi đăng nhập | Không | Có | Tùy chọn đổi khi đạt 100% | macOS |

ccswap là một bản fork của claude-swap, quản lý thêm cả đăng nhập Codex. [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) và [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) cũng cho mỗi tài khoản một thư mục riêng. Một số công cụ trong số này cũng hỗ trợ các coding agent khác như Codex hoặc Gemini.

## Cách chúng được xây dựng

| Công cụ | Viết bằng | Trong lúc Claude Code chạy | Sao chép thông tin đăng nhập giữa các nơi | Những gì các tài khoản có thể dùng chung |
| --- | --- | --- | --- | --- |
| ccam | Go, một tệp nhị phân tĩnh | Không có gì: ccam tự thay thế mình bằng `claude` (trên Windows, nơi không làm được điều này, nó chờ tiến trình đó kết thúc; một relay tùy chọn chạy bên cạnh) | Không | Nhóm: mỗi nhóm nêu các mục nó dùng chung và các tài khoản của nó; một tài khoản có thể thuộc nhiều nhóm |
| claude-swap | Python 3.12+, cài bằng uv hoặc pipx | Không có gì với `cswap run`; `cswap auto` và ứng dụng thanh menu vẫn chạy | Có, khi hoán đổi đăng nhập | Với `cswap run`: cài đặt, `CLAUDE.md` và tiện ích mở rộng đi cùng nhau, lịch sử tùy chọn; máy chủ MCP được sao chép mỗi lần khởi chạy |
| CCS | Node.js 18+ hoặc Bun | CCS vẫn là tiến trình cha; nhóm tài khoản thêm một proxy cục bộ | Không với các thư mục tài khoản của nó; nhóm tài khoản giữ thông tin đăng nhập riêng | Dùng chung hoặc cục bộ theo từng tài khoản, lịch sử trong các nhóm ngữ cảnh; máy chủ MCP được sao chép |
| clausona | TypeScript, Node.js 20+ | `claude` thông thường chạy qua một hàm shell; `clausona run` giữ Node làm tiến trình cha; một tiến trình Node chạy nền sau mỗi lần thoát | Không, nhưng nó tự gia hạn token của các tài khoản đang rảnh để đọc mức sử dụng của chúng | Mọi thứ trừ các tệp riêng của từng tài khoản, từ một tài khoản chính; lịch sử tùy chọn; máy chủ thêm bằng `claude mcp add` không được dùng chung |
| clauth | Rust | clauth vẫn là tiến trình cha | Có, khi hoán đổi đăng nhập | Toàn bộ `~/.claude`, hoặc không gì cả |
| CAAM | Go | Tùy chế độ; có một daemon tùy chọn | Có, ở chế độ kho (vault) | Cố định theo chế độ; máy chủ MCP được sao chép mỗi lần khởi chạy |
| cux | Go | Một lớp bọc nằm quanh `claude` | Có | Mọi thứ: một thư mục |
| CCSwitcher, Claude Switcher | Ứng dụng macOS (Swift, Python) | Ứng dụng thanh menu | Có | Mọi thứ: một thư mục |

## Điểm mạnh của ccam

### Thông tin đăng nhập nằm nguyên nơi Claude Code đặt

ccam không bao giờ sao chép thông tin đăng nhập Claude từ nơi này sang nơi khác. Thông tin đăng nhập của mỗi tài khoản nằm trong thư mục riêng của nó, hoặc trên macOS là trong mục Chuỗi khóa (Keychain) riêng của thư mục đó, và Claude Code là chương trình gia hạn nó. Không có bản sao đã lưu nào có thể bị lệch. Một `~/.claude` có sẵn được đăng ký bằng `ccam import` vẫn giữ mục Chuỗi khóa mà lệnh `claude` thông thường dùng, nên nó vẫn giữ đăng nhập dù bạn chạy theo cách nào.

ccam chỉ tự gia hạn token đã hết hạn khi bạn chạy `ccam usage --refresh`. Nó khóa trước và từ chối nếu tài khoản đó đang chạy, bất kể Claude Code được khởi động theo cách nào (kể cả phiên nền và daemon của chúng), và khi chính Claude Code đang gia hạn đăng nhập đó, nên không bao giờ thay token ngay dưới một phiên đang hoạt động.

### Không có gì chạy giữa bạn và Claude Code

ccam đọc một tệp, kiểm tra các liên kết của nó rồi tự thay thế mình bằng `claude` (`exec`). Trong lúc bạn làm việc, không có tiến trình ccam, lớp bọc, proxy hay daemon nào: đầu vào chuẩn, tín hiệu và mã thoát thuộc về Claude Code, và các yêu cầu của nó đi thẳng tới Anthropic. Windows không thể thay thế một tiến trình, nên ở đó ccam khởi động Claude Code với cùng console và chờ mã thoát của nó, ngoài ra không làm gì khác.

Khởi động một tài khoản theo cách này chỉ tốn thêm khoảng 4 ms, dưới 6 ms ở phân vị thứ 95, với bộ nhớ đỉnh 8 MB ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) với hyperfine, 300 lần chạy). Trên cùng chiếc laptop đó, chỉ riêng việc khởi động trình thông dịch Python đã mất 13 ms và Node.js mất 31 ms, trước khi chương trình làm bất cứ điều gì.

### Nhóm quyết định những gì được dùng chung

Máy chủ MCP, `settings.json`, `CLAUDE.md`, tiện ích mở rộng (agent, lệnh, skill, plugin) và lịch sử hội thoại là năm mục riêng biệt. Một nhóm nêu các mục nó dùng chung và các tài khoản thuộc nó, nên hai tài khoản công việc có thể dùng chung lịch sử trong nhóm `work` trong khi tài khoản cá nhân giữ lịch sử riêng, và cả ba dùng chung cùng các máy chủ MCP qua `default`. `ccam group set work --add memory` chuyển mọi tài khoản của nhóm trong một kế hoạch; [bảng điều khiển](/features/dashboard) làm điều tương tự bằng các ô đánh dấu và làm mờ một mục mà tài khoản đã nhận từ một nhóm khác.

- Máy chủ MCP được truyền bằng `--mcp-config` khi khởi chạy, nên theo mặc định một lần khởi chạy không đụng đến `.claude.json`, tệp mà Claude Code ghi lại trong lúc chạy. ccam chỉ ghi vào đó ở những bước có nói rõ, dưới khóa của Claude Code và có bản sao lưu: cho Remote Control sau khi hỏi, khi một tài khoản rời hoặc quay lại việc dùng chung MCP, và cho `ccam mcp import`. `ccam mcp add` thêm một máy chủ một lần cho cả nhóm.
- Các mục còn lại là liên kết tới thư mục nhóm, nên thay đổi ở một tài khoản hiện ra ngay ở các tài khoản khác. Không có gì được sao chép khi khởi chạy.
- Thông tin đăng nhập, `.claude.json` và trạng thái đang chạy không bao giờ được dùng chung.

### Mọi thay đổi đều được hiển thị trước và sao lưu

`ccam group` và `ccam migrate` in ra kế hoạch và hỏi trước khi di chuyển bất cứ thứ gì. Khi một tệp khác nhau giữa tài khoản và nhóm của nó, bạn chọn: giữ bản của nhóm, dùng bản của bạn, hoặc gộp JSON; bên trong thư mục thì bản của nhóm được giữ lại và kế hoạch nêu tên các tệp khác nhau. Bản gốc được đưa vào `.ccam-backup/<time>/`, tài khoản bị xóa được đăng xuất và thư mục của nó được đưa vào một thư mục thùng rác, và tệp được thay bằng một lần đổi tên duy nhất, nên một lần ghi bị gián đoạn vẫn để lại tệp cũ. `--dry-run` cho thấy việc khởi chạy, dùng chung, di chuyển, lối tắt, cập nhật và sửa chữa sẽ làm gì. `ccam doctor --fix` sửa liên kết hỏng và quyền truy cập.

### Cùng cuộc hội thoại trên tài khoản kế tiếp

`ccam next --auto --continue` hỏi mức sử dụng của mọi tài khoản cùng lúc, mở tài khoản còn nhiều hạn mức nhất và tiếp tục cuộc hội thoại. Khi dùng chung lịch sử, đó là cùng một cuộc hội thoại; nếu không, Claude Code sao chép nó sang tài khoản mới (`--resume <file> --fork-session`), và ccam cho biết câu trả lời đầu tiên phải đọc lại bao nhiêu token. Khi một giới hạn dừng một phiên, hook của plugin ccam hiện một thông báo, và `/ccam:switch` chuyển cuộc hội thoại đi tiếp. Bạn quyết định khi nào đổi; không có gì theo dõi mức sử dụng của bạn ở chế độ nền, và các lựa chọn tự động chỉ nằm trong một tổ chức, nên một cuộc hội thoại của công ty không bao giờ vô tình rơi vào một tài khoản cá nhân. Dòng trạng thái gợi ý lệnh này khi một tài khoản chạm 90%.

Script không thể chờ bạn, nên `ccam --fallback work -p "…"` tự đổi thay cho chúng: khi tài khoản chạm giới hạn, cùng lệnh đó chạy lại bằng tài khoản còn nhiều hạn mức nhất, và chỉ đầu ra của lần chạy đã hoàn tất được in ra. Lỗi không phải do giới hạn sử dụng được chuyển tiếp nguyên vẹn. Với các worker chạy song song, `ccam pick` giao cho mỗi worker một tài khoản khác nhau còn hạn mức, và cả hai đều thoát với mã 75 khi không tài khoản nào còn hạn mức.

**Seat của gói Team và Enterprise.** Lựa chọn tự động (`ccam next`, `--fallback`, `ccam pick`, phím `n` trên bảng điều khiển) cũng chuyển đổi giữa các seat của cùng một tổ chức Claude Team hoặc Enterprise. Trung tâm trợ giúp của Anthropic ghi rằng "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (mỗi seat của gói Team được gán cho đích danh một cá nhân, và chỉ người đó được phép sử dụng) và "Users can't route their work through another user's account to use that user's remaining usage" (người dùng không được chuyển công việc của mình qua tài khoản của một người dùng khác để dùng phần hạn mức còn lại của người đó) ([gói Team](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); tương tự đối với [gói Enterprise](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0)). Vì vậy, dùng lựa chọn tự động giữa các seat có thể vi phạm điều khoản của Anthropic. Hãy xác nhận với tổ chức của bạn và Anthropic trước khi dựa vào nó. Các gói cá nhân (Pro, Max) được đề cập trong [ghi chú về điều khoản ở phần Câu hỏi thường gặp](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts).

### Đúng tài khoản cho từng thư mục

`ccam folder work ~/code/company` khiến mọi lần khởi chạy không nêu tên tài khoản bên trong thư mục đó dùng `work`, trong khi các dự án còn lại vẫn dùng tài khoản mặc định. Thiết lập này nằm trong cùng tệp `config.toml` mà mọi lần khởi chạy đều đọc, nên không tốn thêm thời gian.

### Một tệp nhị phân, tải xuống được kiểm tra

ccam là một tệp nhị phân tĩnh khoảng 10 MB, không cần cài thêm môi trường chạy nào, cho Linux, macOS và Windows. Các script cài đặt (`install.sh`, `install.ps1`) và `ccam update` kiểm tra mọi tệp tải xuống với checksum SHA-256 của bản phát hành, và ccam không bao giờ cập nhật mà không hỏi. Nó không có telemetry: các yêu cầu duy nhất nó gửi là tra cứu mức sử dụng các tài khoản của bạn tới Anthropic, chạy song song với giới hạn 3 giây và bộ nhớ đệm 60 giây, câu hỏi bảng điều khiển gửi tới GitHub về bản phát hành mới nhất khi nó mở, cùng các lượt tải xuống mà bạn yêu cầu. Với `usage_source = "claude"`, ngay cả các yêu cầu đó cũng đi qua Claude Code, và ccam không bao giờ chạm vào token.

## Những điểm công cụ khác làm được nhiều hơn

- claude-swap, clauth, cux, CAAM và CCSwitcher có thể tự đổi tài khoản cho một phiên tương tác đang chạy khi một tài khoản vượt ngưỡng sử dụng. ccam đổi các lần chạy `-p` bằng `--fallback` và để các phiên tương tác cho bạn, với một thông báo khi chạm giới hạn và `/ccam:switch`.
- Trên Windows không bật Developer Mode, ccam dùng chung được thư mục nhưng không dùng chung được tệp đơn lẻ như `settings.json` (khi đó lần `ccam add` đầu tiên chỉ dùng chung máy chủ MCP và tiện ích mở rộng); claude-swap và CCS thì sao chép chúng thay vào đó.
- CCSwitcher, Claude Switcher và claude-swap hiển thị mức sử dụng trên thanh menu macOS.
- CCS, clausona, CAAM, clauth, ccswap và Claude Switcher cũng quản lý Codex và các coding agent khác.
- clausona còn bật và tắt tiện ích mở rộng theo từng tài khoản (`ccam ext` của ccam chỉ hiển thị chúng), theo dõi chi phí token cục bộ, và cung cấp các plugin Claude Code chạy hàng loạt worker trong tmux và các terminal khác.
- CAAM, CCSwitcher và Claude Switcher cài được bằng Homebrew, còn clauth và CAAM ký checksum của bản phát hành. ccam cài bằng script riêng, gói `.deb` hoặc `go install`, và kiểm tra checksum mà không có chữ ký.

## Công cụ nào phù hợp?

Hãy chọn công cụ hoán đổi đăng nhập nếu bạn chỉ dùng một tài khoản mỗi lúc hoặc muốn các phiên đang chạy tự động đổi tài khoản. Hãy chọn ccam nếu bạn muốn chạy song song các tài khoản, một bộ thiết lập dùng chung, cùng cuộc hội thoại trên tài khoản kế tiếp và không có gì chạy bên cạnh Claude Code, trên Linux, macOS hoặc Windows. Nếu không cần dùng chung gì, [alias đơn giản của shell](/compare/aliases) là đủ.

Việc bạn có được dùng nhiều tài khoản hay không phụ thuộc vào điều khoản của từng tài khoản. Hãy kiểm tra trước khi dùng.
