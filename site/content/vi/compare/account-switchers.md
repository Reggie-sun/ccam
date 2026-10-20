---
title: ccam và các công cụ đổi tài khoản Claude Code khác
description: "So sánh ccam với claude-swap, CCS, clauth, CAAM, cux và các công cụ đổi tài khoản Claude Code khác, và cách tiếp cận nào hợp với nhu cầu nào."
updated: 2026-10-20
order: 2
translation: machine
sourceHash: d058da937218
---

Các công cụ đổi tài khoản Claude Code đi theo một trong hai cách. Phần lớn hoán đổi thông tin đăng nhập bên trong một thư mục `~/.claude` duy nhất, nên mỗi lúc chỉ có một tài khoản hoạt động. Một số khác, trong đó có ccam, cho mỗi tài khoản một thư mục riêng bằng `CLAUDE_CONFIG_DIR`, nên mọi tài khoản đều giữ trạng thái đăng nhập và có thể chạy nhiều tài khoản cùng lúc. ccam còn bổ sung việc dùng chung cài đặt và máy chủ MCP, và tiếp tục một cuộc hội thoại trong thư mục của tài khoản khác.

## Hoán đổi một đăng nhập hay dùng thư mục riêng?

**Hoán đổi đăng nhập** chỉ giữ một bộ cài đặt, máy chủ MCP và lịch sử, vì chỉ có một thư mục. Việc chuyển đổi nhanh, và một phiên đang chạy có thể nhận đăng nhập mới. Nhưng mỗi lúc chỉ có một tài khoản hoạt động, cho mọi terminal và trình soạn thảo cùng lúc.

Để hoán đổi, công cụ giữ một bản sao thông tin đăng nhập của mọi tài khoản khác và ghi bản đó vào đúng chỗ. Mỗi lần gia hạn đăng nhập, Claude Code lại thay refresh token bên trong, nên mỗi bản sao đã lưu phải luôn được giữ khớp với các phiên đang chạy. Tài khoản bị đăng xuất hoặc phải `/login` lại sau khi chuyển là một trong những lỗi mà các dự án này phải sửa thường xuyên nhất.

**Thư mục riêng** là cách mà [tài liệu của Claude Code mô tả](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Mọi tài khoản đều giữ trạng thái đăng nhập, và bạn có thể chạy song song một phiên công việc và một phiên cá nhân. Thông tin đăng nhập của mỗi tài khoản nằm yên trong thư mục riêng, nơi Claude Code gia hạn nó như bình thường. Cái giá phải trả là mỗi thư mục bắt đầu với cài đặt, máy chủ MCP và lịch sử riêng, trừ khi có công cụ dùng chung chúng.

## So sánh các công cụ

Tổng hợp từ tài liệu và mã nguồn của từng dự án vào ngày 2026-10-20. "Không có trong tài liệu" nghĩa là dự án không nói đến, không có nghĩa là thiếu tính năng.

| Công cụ | Cách tiếp cận | Nhiều tài khoản cùng lúc | Mức sử dụng 5 giờ và hàng tuần | Khi chạm giới hạn | Nền tảng |
| --- | --- | --- | --- | --- | --- |
| ccam | Thư mục riêng | Có | Có, `ccam usage` | Bạn chạy `ccam next --auto --continue` để tiếp tục cuộc hội thoại trên tài khoản còn nhiều hạn mức nhất; `ccam --fallback` tự chạy lại các tác vụ `-p` ở đó | Linux, macOS, Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | Hoán đổi đăng nhập; thư mục riêng với `cswap run` | Có, với `cswap run` | Có | `cswap auto` đổi đăng nhập khi chạm ngưỡng; phiên đang chạy vẫn tiếp tục | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Thư mục riêng | Có | Một phần | Nhóm tài khoản chuyển yêu cầu sang tài khoản kế tiếp qua một proxy cục bộ | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Hoán đổi đăng nhập; thư mục riêng với `clauth start` | Có, với `clauth start` | Có | Chuyển lần lượt theo chuỗi tài khoản dự phòng | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Hoán đổi đăng nhập, hoặc hồ sơ tách biệt | Có, với hồ sơ tách biệt | Có | `caam run` thử lại các lệnh không tương tác trên tài khoản khác | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Hoán đổi đăng nhập và bọc lệnh `claude` | Không, chỉ một tài khoản hoạt động | Có | Hoán đổi đăng nhập và tiếp tục cùng cuộc hội thoại | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Hoán đổi đăng nhập | Không | Có | Tùy chọn đổi khi chạm ngưỡng | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Hoán đổi đăng nhập | Không | Có | Tùy chọn đổi khi đạt 100% | macOS |

ccswap là một bản fork của claude-swap, quản lý thêm cả đăng nhập Codex. [claude-account](https://github.com/hamzarehmandeveloper/claude-account), [claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) và [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) cũng cho mỗi tài khoản một thư mục riêng. Một số công cụ trong số này cũng hỗ trợ các coding agent khác như Codex hoặc Gemini.

## Cách chúng được xây dựng

| Công cụ | Viết bằng | Trong lúc Claude Code chạy | Sao chép thông tin đăng nhập giữa các nơi | Những gì các tài khoản có thể dùng chung |
| --- | --- | --- | --- | --- |
| ccam | Go, một tệp nhị phân tĩnh | Không có gì: ccam tự thay thế mình bằng `claude` (trên Windows, nơi không làm được điều này, nó chờ tiến trình đó kết thúc) | Không | Từng mục riêng: cục bộ hoặc một nhóm có tên, theo từng tài khoản |
| claude-swap | Python 3.12+, cài bằng uv hoặc pipx | Không có gì với `cswap run`; `cswap auto` và ứng dụng thanh menu vẫn chạy | Có, khi hoán đổi đăng nhập | Với `cswap run`: cài đặt, `CLAUDE.md` và tiện ích mở rộng đi cùng nhau, lịch sử tùy chọn; máy chủ MCP được sao chép mỗi lần khởi chạy |
| CCS | Node.js 18+ hoặc Bun | CCS vẫn là tiến trình cha; nhóm tài khoản thêm một proxy cục bộ | Không với các thư mục tài khoản của nó; nhóm tài khoản giữ thông tin đăng nhập riêng | Dùng chung hoặc cục bộ theo từng tài khoản, lịch sử trong các nhóm ngữ cảnh; máy chủ MCP được sao chép |
| clauth | Rust | clauth vẫn là tiến trình cha | Có, khi hoán đổi đăng nhập | Toàn bộ `~/.claude`, hoặc không gì cả |
| CAAM | Go | Tùy chế độ; có một daemon tùy chọn | Có, ở chế độ kho (vault) | Cố định theo chế độ; máy chủ MCP được sao chép mỗi lần khởi chạy |
| cux | Go | Một lớp bọc nằm quanh `claude` | Có | Mọi thứ: một thư mục |
| CCSwitcher, Claude Switcher | Ứng dụng macOS (Swift, Python) | Ứng dụng thanh menu | Có | Mọi thứ: một thư mục |

## Điểm mạnh của ccam

### Thông tin đăng nhập nằm nguyên nơi Claude Code đặt

ccam không bao giờ sao chép thông tin đăng nhập Claude từ nơi này sang nơi khác. Thông tin đăng nhập của mỗi tài khoản nằm trong thư mục riêng của nó, hoặc trên macOS là trong mục Chuỗi khóa (Keychain) riêng của thư mục đó, và Claude Code là chương trình gia hạn nó. Không có bản sao đã lưu nào có thể bị lệch. Một `~/.claude` có sẵn được đăng ký bằng `ccam import` vẫn giữ mục Chuỗi khóa mà lệnh `claude` thông thường dùng, nên nó vẫn giữ đăng nhập dù bạn chạy theo cách nào.

ccam chỉ tự gia hạn token đã hết hạn khi bạn chạy `ccam usage --refresh`. Nó khóa trước và từ chối nếu tài khoản đó đang chạy, nên không bao giờ thay token ngay dưới một phiên đang hoạt động.

### Không có gì chạy giữa bạn và Claude Code

ccam đọc một tệp, kiểm tra các liên kết của nó rồi tự thay thế mình bằng `claude` (`exec`). Trong lúc bạn làm việc, không có tiến trình ccam, lớp bọc, proxy hay daemon nào: đầu vào chuẩn, tín hiệu và mã thoát thuộc về Claude Code, và các yêu cầu của nó đi thẳng tới Anthropic. Windows không thể thay thế một tiến trình, nên ở đó ccam khởi động Claude Code với cùng console và chờ mã thoát của nó, ngoài ra không làm gì khác.

Khởi động một tài khoản theo cách này chỉ tốn thêm khoảng 4 ms, dưới 6 ms ở phân vị thứ 95, với bộ nhớ đỉnh 8 MB ([`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) với hyperfine, 300 lần chạy). Trên cùng chiếc laptop đó, chỉ riêng việc khởi động trình thông dịch Python đã mất 13 ms và Node.js mất 31 ms, trước khi chương trình làm bất cứ điều gì.

### Dùng chung theo từng mục, từng tài khoản

Máy chủ MCP, `settings.json`, `CLAUDE.md`, tiện ích mở rộng (agent, lệnh, skill, plugin) và lịch sử hội thoại là năm mục riêng biệt. Mỗi tài khoản đặt từng mục là `local` hoặc một nhóm, nên hai tài khoản công việc có thể dùng chung lịch sử trong nhóm `work` trong khi tài khoản cá nhân giữ lịch sử riêng, và cả ba dùng chung cùng các máy chủ MCP.

- Máy chủ MCP được truyền bằng `--mcp-config` khi khởi chạy, nên theo mặc định ccam không bao giờ sửa `.claude.json`, tệp mà Claude Code ghi lại trong lúc chạy. `ccam mcp add` thêm một máy chủ một lần cho cả nhóm.
- Các mục còn lại là liên kết tới thư mục nhóm, nên thay đổi ở một tài khoản hiện ra ngay ở các tài khoản khác. Không có gì được sao chép khi khởi chạy.
- Thông tin đăng nhập, `.claude.json` và trạng thái đang chạy không bao giờ được dùng chung.

### Mọi thay đổi đều được hiển thị trước và sao lưu

`ccam share` và `ccam migrate` in ra kế hoạch và hỏi trước khi di chuyển bất cứ thứ gì. Khi một tệp khác nhau giữa tài khoản và nhóm của nó, bạn chọn: giữ bản của nhóm, dùng bản của bạn, hoặc gộp JSON. Bản gốc được đưa vào `.ccam-backup/<time>/`, thư mục của tài khoản bị xóa được đưa vào một thư mục thùng rác, và tệp được thay bằng một lần đổi tên duy nhất, nên một lần ghi bị gián đoạn vẫn để lại tệp cũ. `--dry-run` cho thấy việc khởi chạy, dùng chung, di chuyển, lối tắt, cập nhật và sửa chữa sẽ làm gì. `ccam doctor --fix` sửa liên kết hỏng và quyền truy cập.

### Cùng cuộc hội thoại trên tài khoản kế tiếp

Khi dùng chung lịch sử, `ccam next --auto --continue` hỏi mức sử dụng của mọi tài khoản cùng lúc, mở tài khoản còn nhiều hạn mức nhất và tiếp tục cuộc hội thoại. Bạn quyết định khi nào đổi; ccam không theo dõi mức sử dụng của bạn ở chế độ nền. Dòng trạng thái gợi ý lệnh này khi một tài khoản chạm 90%.

Script không thể chờ bạn, nên `ccam --fallback work -p "…"` tự đổi thay cho chúng: khi tài khoản chạm giới hạn, cùng lệnh đó chạy lại bằng tài khoản còn nhiều hạn mức nhất, và chỉ đầu ra của lần chạy đã hoàn tất được in ra. Lỗi không phải do giới hạn sử dụng được chuyển tiếp nguyên vẹn.

### Đúng tài khoản cho từng thư mục

`ccam folder work ~/code/company` khiến mọi lần khởi chạy không nêu tên tài khoản bên trong thư mục đó dùng `work`, trong khi các dự án còn lại vẫn dùng tài khoản mặc định. Thiết lập này nằm trong cùng tệp `config.toml` mà mọi lần khởi chạy đều đọc, nên không tốn thêm thời gian.

### Một tệp nhị phân, tải xuống được kiểm tra

ccam là một tệp nhị phân tĩnh khoảng 10 MB, không cần cài thêm môi trường chạy nào, cho Linux, macOS và Windows. Các script cài đặt (`install.sh`, `install.ps1`) và `ccam update` kiểm tra mọi tệp tải xuống với checksum SHA-256 của bản phát hành, và ccam không bao giờ tự kiểm tra cập nhật. Nó không có telemetry: các yêu cầu duy nhất nó gửi là tra cứu mức sử dụng các tài khoản của bạn tới Anthropic, chạy song song với giới hạn 3 giây và bộ nhớ đệm 60 giây, cùng các lượt tải xuống mà bạn yêu cầu.

## Những điểm công cụ khác làm được nhiều hơn

- claude-swap, clauth, cux, CAAM và CCSwitcher có thể tự đổi tài khoản cho một phiên tương tác đang chạy khi một tài khoản vượt ngưỡng sử dụng. ccam đổi các lần chạy `-p` bằng `--fallback` và để các phiên tương tác cho bạn.
- Trên Windows không bật Developer Mode, ccam dùng chung được thư mục nhưng không dùng chung được tệp đơn lẻ như `settings.json`; claude-swap và CCS thì sao chép chúng thay vào đó.
- CCSwitcher, Claude Switcher và claude-swap hiển thị mức sử dụng trên thanh menu macOS.
- CCS, CAAM, clauth, ccswap và Claude Switcher cũng quản lý Codex và các coding agent khác.
- CAAM, CCSwitcher và Claude Switcher cài được bằng Homebrew, còn clauth và CAAM ký checksum của bản phát hành. ccam cài bằng script riêng, gói `.deb` hoặc `go install`, và kiểm tra checksum mà không có chữ ký.

## Công cụ nào phù hợp?

Hãy chọn công cụ hoán đổi đăng nhập nếu bạn chỉ dùng một tài khoản mỗi lúc hoặc muốn các phiên đang chạy tự động đổi tài khoản. Hãy chọn ccam nếu bạn muốn chạy song song các tài khoản, một bộ thiết lập dùng chung, cùng cuộc hội thoại trên tài khoản kế tiếp và không có gì chạy bên cạnh Claude Code, trên Linux, macOS hoặc Windows. Nếu không cần dùng chung gì, [alias đơn giản của shell](/compare/aliases) là đủ.

Việc bạn có được dùng nhiều tài khoản hay không phụ thuộc vào điều khoản của từng tài khoản. Hãy kiểm tra trước khi dùng.
