---
title: Cập nhật ccam
description: "ccam update thay ccam bằng bản phát hành mới nhất trên GitHub sau khi kiểm tra SHA-256; --check chỉ báo cáo và --version cài một thẻ phát hành cụ thể."
updated: 2026-10-19
order: 10
translation: machine
sourceHash: 83fc20191328
---

`ccam update` tải một bản phát hành ccam từ GitHub, kiểm tra nó với `checksums.txt` của bản phát hành và thay tệp thực thi đang chạy bằng một lần đổi tên. Hãy dùng lệnh này khi bạn cài ccam bằng script cài đặt hoặc từ một tệp nén. ccam không bao giờ tự kiểm tra cập nhật: chỉ lệnh này liên lạc với GitHub, nên mọi lệnh khác vẫn hoạt động như nhau khi không có mạng.

## Cách dùng

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, also an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

Khi không có bản phát hành mới hơn, ccam in ra `ccam 0.4.0 is up to date.` và không thay đổi gì. Với `--version` và đúng phiên bản bạn đang chạy, nó in ra `ccam 0.4.0 is already installed.`

## Tùy chọn

| Tùy chọn | Tác dụng | Mặc định |
| --- | --- | --- |
| `--check` | Chỉ báo có bản phát hành mới hơn hay không | tắt |
| `--version <tag>` | Cài thẻ phát hành này, chẳng hạn `v0.3.0`, mới hơn hoặc cũ hơn bản hiện tại | mới nhất |
| `--dry-run` | Được hỗ trợ: tra thẻ mới nhất (trừ khi `--version` nêu một thẻ), chạy cùng các kiểm tra trên tệp thực thi, và in ra các đường dẫn tải về, mã kiểm tra và thay thế mà không tải gì | tắt |

Không thể dùng `--check` và `--version` cùng nhau. Thẻ phải có dạng `vMAJOR.MINOR.PATCH`, có thể kèm hậu tố pre-release hoặc build; có thể bỏ chữ `v` ở đầu. Mọi dạng khác đều thất bại, ví dụ `--version takes a release tag such as v0.3.0, not "1.2"`.

Biến môi trường `CCAM_RELEASES_URL` thay cho `https://github.com/jhyunwoo/ccam/releases`, dùng cho máy chủ gương hoặc cho kiểm thử. Nó phải được tổ chức theo cùng cách: `<url>/latest` chuyển hướng tới `<url>/tag/<tag>`, và các tệp nằm ở `<url>/download/<tag>/<asset>`.

## Cách hoạt động

Mã nằm trong `internal/update`, gói chỉ được `ccam update` import. Đường khởi động không bao giờ gọi tới nó.

1. `cli.updateCmd` tìm tệp thực thi đang chạy bằng `os.Executable` và phân giải liên kết tượng trưng bằng `filepath.EvalSymlinks`. Đường dẫn thật đó là thứ sẽ bị thay thế.
2. Không có `--version`, `Updater.Latest` gửi yêu cầu tới `<releases>/latest` mà không đi theo chuyển hướng, và đọc thẻ từ header `Location` (`…/tag/v0.5.0`). Cách này dùng trang releases thay cho GitHub API, nên giới hạn 60 yêu cầu mỗi giờ của API không áp dụng. Việc tra cứu có 10 giây.
3. `parseVersion` kiểm tra thẻ theo một mẫu nghiêm ngặt. Thẻ sẽ xuất hiện trong URL tải về và trong một lệnh shell được in ra, nên không có gì chứa `/`, `%`, khoảng trắng hay ký tự shell lọt qua được. `compareVersions` so sánh nó với phiên bản đang chạy; một bản pre-release được xếp trước bản phát hành của nó.
4. `managed()` và `writable()` quyết định ccam có được phép thay tệp thực thi hay không (xem bên dưới).
5. `Updater.Fetch` tải `ccam_<os>_<arch>.tar.gz` và `checksums.txt` cùng lúc, với 2 phút cho cả hai. Mỗi lần tải giới hạn ở 256 MB.
6. SHA-256 của tệp nén phải khớp với dòng của nó trong `checksums.txt` trước khi bất cứ việc gì khác xảy ra. Thiếu dòng hoặc mã kiểm tra khác nhau sẽ dừng việc cập nhật với `nothing was changed`.
7. `extract()` đọc tệp tar và lấy tệp thông thường đầu tiên có tên `ccam`. Thư mục, liên kết và các mục khác bị bỏ qua, và tệp thực thi cũng bị giới hạn ở 256 MB.
8. `Updater.Install` ghi tệp thực thi mới bằng `fsutil.AtomicWrite`: một tệp tạm `.ccam.tmp*` cạnh tệp thực thi, với các bit quyền của tệp cũ, được ghi xuống đĩa rồi đổi tên đè lên tệp cũ.

Toàn bộ quy trình, từ tra thẻ đến đổi tên:

<!-- motion:update-pipeline -->

Tệp nén và `checksums.txt` được tải song song, và không có gì được ghi cho đến khi mã kiểm tra khớp. Việc đổi tên là nguyên tử, nên một sự cố hay đĩa đầy chỉ để lại tệp thực thi cũ hoặc mới, không bao giờ là một nửa. Các lệnh tắt như `~/.local/bin/ccwork` là liên kết tượng trưng tới tệp thực thi ccam, nên chúng chạy phiên bản mới mà không cần thay đổi gì.

### Bản cài qua trình quản lý gói

ccam được cài từ `.deb` thuộc về apt, và thay thế nó mà apt không biết sẽ khiến cơ sở dữ liệu gói bị sai. Trước khi tải về, ccam kiểm tra ai sở hữu tệp thực thi:

<!-- motion:update-managed -->

1. Tệp thực thi trong `/usr/bin`, `/bin`, `/usr/sbin` hoặc `/sbin` được coi là do trình quản lý gói quản lý.
2. Ở nơi khác, nếu `dpkg` đã được cài và `dpkg -S <path>` thành công, nó cũng được coi là do trình quản lý gói quản lý.
3. Nếu không, ccam tạo rồi xóa một tệp tạm `.ccam-update-*` cạnh tệp thực thi để chắc chắn nó ghi được ở đó.

Với tệp thực thi do trình quản lý gói quản lý, ccam dừng lại và in ra các bước với apt:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

Trên các hệ thống không phải Linux, chỉ dòng script cài đặt được in ra. Khi thư mục không ghi được, ccam gợi ý chạy lại script cài đặt, script này cài ccam vào `~/.local/bin`.

## An toàn và các trường hợp đặc biệt

- Bản dựng phát triển báo phiên bản là `dev` và không có bản phát hành nào để so sánh, nên `ccam update` và `ccam update --check` từ chối nó. `ccam update --version vX.Y.Z` cài một bản phát hành đè lên nó.
- Tải về thất bại, thiếu tệp cho nền tảng của bạn, hoặc mã kiểm tra không khớp đều giữ nguyên tệp thực thi đã cài.
- Bản phát hành không tồn tại sẽ thất bại với `there is no release <tag>`. Khi GitHub từ chối với HTTP 403 hoặc 429, ccam cho biết có lẽ đó là giới hạn tần suất và đề nghị bạn thử lại sau.
- `--check` chỉ đọc. Nó không kiểm tra tệp thực thi có do trình quản lý gói quản lý hay có ghi được hay không.
- `--version` với một thẻ cũ hơn là cách để quay về phiên bản trước.

## Liên quan

- [Cài đặt](/docs/install)
- [Lệnh](/docs/commands)
- [Câu hỏi thường gặp](/docs/faq)
