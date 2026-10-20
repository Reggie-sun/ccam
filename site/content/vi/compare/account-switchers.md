---
title: ccam và các công cụ đổi tài khoản Claude Code khác
description: "So sánh ccam với claude-swap, CCS, clauth, CAAM, cux và các công cụ đổi tài khoản Claude Code khác, và cách tiếp cận nào hợp với nhu cầu nào."
updated: 2026-10-19
order: 2
translation: machine
sourceHash: dc907c88e2bf
---

Các công cụ đổi tài khoản Claude Code đi theo một trong hai cách. Phần lớn hoán đổi thông tin đăng nhập bên trong một thư mục `~/.claude` duy nhất, nên mỗi lúc chỉ có một tài khoản hoạt động. Một số khác, trong đó có ccam, cho mỗi tài khoản một thư mục riêng bằng `CLAUDE_CONFIG_DIR`, nên mọi tài khoản đều giữ trạng thái đăng nhập và có thể chạy nhiều tài khoản cùng lúc. ccam còn bổ sung việc dùng chung cài đặt và máy chủ MCP, và tiếp tục một cuộc hội thoại trong thư mục của tài khoản khác.

## Hoán đổi một đăng nhập hay dùng thư mục riêng?

**Hoán đổi đăng nhập** chỉ giữ một bộ cài đặt, máy chủ MCP và lịch sử, vì chỉ có một thư mục. Việc chuyển đổi nhanh, và một phiên đang chạy có thể nhận đăng nhập mới. Nhưng mỗi lúc chỉ có một tài khoản hoạt động, cho mọi terminal và trình soạn thảo cùng lúc.

**Thư mục riêng** là cách mà [tài liệu của Claude Code mô tả](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts). Mọi tài khoản đều giữ trạng thái đăng nhập, và bạn có thể chạy song song một phiên công việc và một phiên cá nhân. Cái giá phải trả là mỗi thư mục bắt đầu với cài đặt, máy chủ MCP và lịch sử riêng, trừ khi có công cụ dùng chung chúng.

## So sánh các công cụ

Tổng hợp từ tài liệu của từng dự án vào ngày 2026-10-19. "Không có trong tài liệu" nghĩa là dự án không nói đến, không có nghĩa là thiếu tính năng.

| Công cụ | Cách tiếp cận | Nhiều tài khoản cùng lúc | Mức sử dụng 5 giờ và hàng tuần | Khi chạm giới hạn | Nền tảng |
| --- | --- | --- | --- | --- | --- |
| ccam | Thư mục riêng | Có | Có, `ccam usage` | Bạn chạy `ccam next --auto --continue` để tiếp tục cuộc hội thoại trên tài khoản còn nhiều hạn mức nhất | Linux, macOS, WSL |
| [claude-swap](https://github.com/realiti4/claude-swap) | Hoán đổi đăng nhập; thư mục riêng với `cswap run` | Có, với `cswap run` | Có | `cswap auto` đổi đăng nhập khi chạm ngưỡng; phiên đang chạy vẫn tiếp tục | Linux, macOS, Windows |
| [CCS](https://github.com/kaitranntt/ccs) | Thư mục riêng | Có | Một phần, trong ứng dụng trên thanh menu macOS | Nhóm tài khoản chuyển yêu cầu sang tài khoản kế tiếp qua một proxy cục bộ | Linux, macOS, Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | Hoán đổi đăng nhập; thư mục riêng với `clauth start` | Có, với `clauth start` | Có | Chuyển lần lượt theo chuỗi tài khoản dự phòng | Linux, macOS, Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | Hoán đổi đăng nhập, hoặc hồ sơ tách biệt | Có, với hồ sơ tách biệt | Có | `caam run` thử lại các lệnh không tương tác trên tài khoản khác | Linux, macOS, Windows |
| [cux](https://github.com/inulute/cux) | Hoán đổi đăng nhập và bọc lệnh `claude` | Không có trong tài liệu | Có | Hoán đổi đăng nhập và tiếp tục cùng cuộc hội thoại | Linux, macOS, Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | Hoán đổi đăng nhập | Không | Có | Không có trong tài liệu | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | Hoán đổi đăng nhập | Không | Có | Tùy chọn đổi khi đạt 100% | macOS |

ccswap là một bản fork của claude-swap, quản lý thêm cả đăng nhập Codex. Một số công cụ trong số này cũng hỗ trợ các coding agent khác như Codex hoặc Gemini.

## ccam tập trung vào điều gì

- Mọi tài khoản đều giữ trạng thái đăng nhập. Việc chuyển đổi không bao giờ đăng xuất tài khoản nào, và `ccam work` với `ccam personal` có thể chạy cùng lúc trong hai terminal.
- Máy chủ MCP, `settings.json`, `CLAUDE.md`, agent, skill và plugin được dùng chung giữa các thư mục qua liên kết, theo từng mục và từng tài khoản. Thông tin đăng nhập không bao giờ được dùng chung.
- Khi dùng chung lịch sử, một lệnh duy nhất, `ccam next --auto --continue`, mở cùng cuộc hội thoại trong thư mục của tài khoản còn nhiều hạn mức nhất.
- ccam khởi động Claude Code rồi tự thay thế mình bằng Claude Code, nên không có gì chạy nền. Nó không theo dõi mức sử dụng hay tự đổi tài khoản; bạn quyết định khi nào đổi.

## Những điểm công cụ khác làm được nhiều hơn

- claude-swap, CCS, clauth, CAAM và cux chạy trên Windows mà không cần WSL. ccam thì không.
- claude-swap, clauth và cux tự đổi tài khoản khi một tài khoản vượt ngưỡng sử dụng. ccam để việc đổi cho bạn.
- CCSwitcher, Claude Switcher và claude-swap hiển thị mức sử dụng trên thanh menu macOS.
- CCS, CAAM, ccswap và Claude Switcher cũng quản lý Codex và các coding agent khác.

## Công cụ nào phù hợp?

Hãy chọn công cụ hoán đổi đăng nhập nếu bạn chỉ dùng một tài khoản mỗi lúc, muốn việc đổi diễn ra tự động, hoặc dùng Windows không có WSL. Hãy chọn ccam nếu bạn muốn chạy song song các tài khoản, một bộ thiết lập dùng chung và cùng cuộc hội thoại trên tài khoản kế tiếp. Nếu không cần dùng chung gì, [alias đơn giản của shell](/compare/aliases) là đủ.

Việc bạn có được dùng nhiều tài khoản hay không phụ thuộc vào điều khoản của từng tài khoản. Hãy kiểm tra trước khi dùng.
