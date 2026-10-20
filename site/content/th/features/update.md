---
title: การอัปเดต ccam
description: "ccam update แทนที่ ccam ด้วยรุ่นล่าสุดบน GitHub หลังตรวจสอบ SHA-256 โดย --check จะรายงานเท่านั้น และ --version จะติดตั้งแท็กรุ่นที่ระบุ"
updated: 2026-10-20
order: 10
translation: machine
sourceHash: d91499b95d3d
---

`ccam update` ดาวน์โหลดรุ่นของ ccam จาก GitHub ตรวจสอบกับ `checksums.txt` ของรุ่นนั้น และแทนที่ไบนารีที่กำลังทำงานด้วยการ rename ครั้งเดียว ใช้คำสั่งนี้เมื่อคุณติดตั้ง ccam ด้วยสคริปต์ติดตั้งหรือจากไฟล์บีบอัด ccam ไม่เคยตรวจหาอัปเดตเอง: มีเพียงคำสั่งนี้ที่ติดต่อ GitHub คำสั่งอื่นทั้งหมดจึงทำงานได้เหมือนเดิมแม้ออฟไลน์

## การใช้งาน

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

เมื่อไม่มีรุ่นที่ใหม่กว่า ccam จะพิมพ์ `ccam 0.4.0 is up to date.` และไม่เปลี่ยนแปลงสิ่งใด หากใช้ `--version` กับเวอร์ชันที่คุณใช้อยู่แล้ว จะพิมพ์ `ccam 0.4.0 is already installed.`

## ตัวเลือก

| ตัวเลือก | สิ่งที่ทำ | ค่าเริ่มต้น |
| --- | --- | --- |
| `--check` | รายงานเพียงว่ามีรุ่นที่ใหม่กว่าหรือไม่ | ปิด |
| `--version <tag>` | ติดตั้งแท็กรุ่นนี้ เช่น `v0.3.0` ไม่ว่าจะใหม่หรือเก่ากว่ารุ่นปัจจุบัน | ล่าสุด |
| `--dry-run` | รองรับ: ค้นหาแท็กล่าสุด (เว้นแต่ `--version` จะระบุไว้) ตรวจไบนารีแบบเดียวกัน และพิมพ์พาธของการดาวน์โหลด checksum และการแทนที่ โดยไม่ดาวน์โหลด | ปิด |

ใช้ `--check` และ `--version` พร้อมกันไม่ได้ แท็กต้องอยู่ในรูป `vMAJOR.MINOR.PATCH` โดยจะมีส่วนต่อท้ายของ pre-release หรือ build ก็ได้ และจะละ `v` ที่นำหน้าไว้ก็ได้ รูปแบบอื่นจะล้มเหลว เช่น `--version takes a release tag such as v0.3.0, not "1.2"`

ตัวแปรสภาพแวดล้อม `CCAM_RELEASES_URL` ใช้แทนที่ `https://github.com/jhyunwoo/ccam/releases` สำหรับมิเรอร์หรือการทดสอบ โดยต้องมีโครงสร้างเหมือนกัน: `<url>/latest` เปลี่ยนเส้นทางไปที่ `<url>/tag/<tag>` และไฟล์อยู่ที่ `<url>/download/<tag>/<asset>`

## วิธีการทำงาน

โค้ดอยู่ใน `internal/update` ซึ่งมีเพียง `ccam update` ที่นำเข้า เส้นทางการเริ่มไม่เคยเรียกใช้

1. `cli.updateCmd` หาไบนารีที่กำลังทำงานด้วย `os.Executable` และหาพาธจริงผ่านลิงก์สัญลักษณ์ด้วย `filepath.EvalSymlinks` พาธจริงนั้นคือสิ่งที่จะถูกแทนที่
2. หากไม่ใส่ `--version` `Updater.Latest` จะส่งคำขอไปที่ `<releases>/latest` โดยไม่ตาม redirect และอ่านแท็กจาก header `Location` (`…/tag/v0.5.0`) วิธีนี้ใช้หน้ารุ่นแทน GitHub API จึงไม่ติดขีดจำกัด 60 คำขอต่อชั่วโมงของ API การค้นหามีเวลา 10 วินาที
3. `parseVersion` ตรวจแท็กกับรูปแบบที่เข้มงวด แท็กจะไปอยู่ใน URL ดาวน์โหลดและในคำสั่งเชลล์ที่พิมพ์ออกมา จึงไม่มีสิ่งใดที่มี `/`, `%` ช่องว่าง หรืออักขระของเชลล์ผ่านไปได้ `compareVersions` เปรียบเทียบแท็กกับเวอร์ชันที่กำลังทำงาน โดย pre-release จะเรียงก่อนรุ่นจริง
4. `managed()` และ `writable()` ตัดสินว่า ccam แทนที่ไบนารีได้หรือไม่ (ดูด้านล่าง)
5. `Updater.Fetch` ดาวน์โหลด `ccam_<os>_<arch>.tar.gz` (บน Windows คือ `ccam_windows_<arch>.zip`) และ `checksums.txt` พร้อมกัน โดยมีเวลา 2 นาทีสำหรับทั้งสองไฟล์ การดาวน์โหลดแต่ละไฟล์จำกัดไว้ที่ 256 MB
6. SHA-256 ของไฟล์บีบอัดต้องตรงกับบรรทัดของไฟล์นั้นใน `checksums.txt` ก่อนที่จะทำสิ่งอื่นใด หากไม่มีบรรทัดนั้นหรือ checksum ต่างกัน การอัปเดตจะหยุดด้วย `nothing was changed`
7. `extract()` อ่านไฟล์ tar และเลือกไฟล์ปกติไฟล์แรกที่ชื่อ `ccam` (`extractZip` เลือก `ccam.exe` จากไฟล์ zip) โฟลเดอร์ ลิงก์ และรายการอื่นจะถูกละเว้น และไบนารีก็จำกัดไว้ที่ 256 MB เช่นกัน
8. `Updater.Install` เขียนไบนารีใหม่ด้วย `fsutil.AtomicWrite`: ไฟล์ชั่วคราว `.ccam.tmp*` ข้างไบนารี โดยใช้บิตสิทธิ์ของไฟล์เดิม บันทึกลงดิสก์ แล้ว rename ทับไฟล์เดิม

กระบวนการทั้งหมด ตั้งแต่ค้นหาแท็กจนถึงการ rename:

<!-- motion:update-pipeline -->

ไฟล์บีบอัดและ `checksums.txt` ถูกดาวน์โหลดคู่กัน และไม่มีสิ่งใดถูกเขียนจนกว่า checksum จะตรงกัน การ rename เป็นแบบ atomic ดังนั้นการแครชหรือดิสก์เต็มจะเหลือไว้เพียงไบนารีเดิมหรือไบนารีใหม่ ไม่มีทางเหลือครึ่ง ๆ กลาง ๆ คำสั่งลัดอย่าง `~/.local/bin/ccwork` เป็นลิงก์สัญลักษณ์ไปยังไบนารี ccam จึงรันเวอร์ชันใหม่ได้โดยไม่ต้องเปลี่ยนแปลงใด ๆ

บน Windows ไฟล์ของโปรแกรมที่กำลังทำงานอยู่ไม่สามารถแทนที่หรือลบได้ ทำได้เพียงเปลี่ยนชื่อ `Updater.Install` จะเขียนไบนารีใหม่ไว้ที่ `ccam.exe.new` เปลี่ยนชื่อ `ccam.exe` ที่กำลังทำงานเป็น `ccam.exe.old` แล้วเปลี่ยนชื่อไฟล์ใหม่เป็น `ccam.exe` หากขั้นตอนสุดท้ายล้มเหลว ไบนารีเดิมจะถูกเปลี่ยนชื่อกลับ การอัปเดตครั้งถัดไปจะลบ `ccam.exe.old` คำสั่งลัดที่เป็นฮาร์ดลิงก์แทนที่จะเป็นลิงก์สัญลักษณ์ ([คำสั่งลัด](/features/shortcuts)) จะถูกต่ออายุทันทีหลังจากนั้น จึงรันเวอร์ชันใหม่ด้วยเช่นกัน

### การติดตั้งที่จัดการโดยตัวจัดการแพ็กเกจ

ccam ที่ติดตั้งจาก `.deb` เป็นของ apt และการแทนที่โดยที่ apt ไม่รู้จะทำให้ฐานข้อมูลแพ็กเกจไม่ถูกต้อง ก่อนดาวน์โหลด ccam จะตรวจว่าใครเป็นเจ้าของไบนารี:

<!-- motion:update-managed -->

1. ไบนารีใน `/usr/bin`, `/bin`, `/usr/sbin` หรือ `/sbin` นับว่าจัดการโดยตัวจัดการแพ็กเกจ
2. ในตำแหน่งอื่น หากติดตั้ง `dpkg` ไว้และ `dpkg -S <path>` สำเร็จ ก็นับว่าจัดการโดยตัวจัดการแพ็กเกจเช่นกัน
3. ในกรณีอื่น ccam จะสร้างและลบไฟล์ชั่วคราว `.ccam-update-*` ข้างไบนารี เพื่อให้แน่ใจว่าเขียนในตำแหน่งนั้นได้

สำหรับไบนารีที่จัดการโดยตัวจัดการแพ็กเกจ ccam จะหยุดและพิมพ์ขั้นตอนของ apt:

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.sh | sh
```

บนระบบอื่นที่ไม่ใช่ Linux จะพิมพ์เฉพาะบรรทัดของสคริปต์ติดตั้ง: `irm https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.ps1 | iex` บน Windows เมื่อโฟลเดอร์เขียนไม่ได้ ccam จะแนะนำให้รันสคริปต์ติดตั้งอีกครั้ง ซึ่งจะติดตั้ง ccam ไว้ใน `~/.local/bin`

## ความปลอดภัยและกรณีพิเศษ

- บิลด์สำหรับการพัฒนาจะรายงานเวอร์ชันเป็น `dev` และไม่มีรุ่นให้เปรียบเทียบ ดังนั้น `ccam update` และ `ccam update --check` จะปฏิเสธ ส่วน `ccam update --version vX.Y.Z` จะติดตั้งรุ่นทับลงไป
- การดาวน์โหลดที่ล้มเหลว ไม่มีไฟล์สำหรับแพลตฟอร์มของคุณ หรือ checksum ไม่ตรงกัน จะไม่เปลี่ยนแปลงไบนารีที่ติดตั้งไว้
- รุ่นที่ไม่มีอยู่จะล้มเหลวด้วย `there is no release <tag>` เมื่อ GitHub ปฏิเสธด้วย HTTP 403 หรือ 429 ccam จะแจ้งว่าน่าจะเป็นเพราะขีดจำกัดอัตราคำขอ และขอให้คุณลองใหม่ภายหลัง
- `--check` อ่านอย่างเดียว ไม่ได้ตรวจว่าไบนารีจัดการโดยตัวจัดการแพ็กเกจหรือเขียนได้หรือไม่
- `--version` กับแท็กที่เก่ากว่าคือวิธีย้อนกลับไปรุ่นก่อนหน้า

## หัวข้อที่เกี่ยวข้อง

- [การติดตั้ง](/docs/install)
- [คำสั่ง](/docs/commands)
- [คำถามที่พบบ่อย](/docs/faq)
