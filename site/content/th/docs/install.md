---
title: ติดตั้ง ccam
metaTitle: "ติดตั้ง ccam บน Linux, macOS และ Windows"
description: "ติดตั้ง ccam บน Linux, macOS หรือ Windows ด้วยคำสั่งเดียว แพ็กเกจ .deb หรือ go install และอัปเดตด้วย ccam update โดยต้องติดตั้ง Claude Code ก่อน"
updated: 2026-10-28
order: 1
translation: machine
sourceHash: "193df30c9d9a"
---

ccam เป็นตัวเริ่ม Claude Code ให้คุณ ดังนั้นให้ติดตั้ง [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) ก่อน และตรวจสอบให้แน่ใจว่า `claude --version` ทำงานได้

รันคำสั่งนี้ในเทอร์มินัลบน Linux หรือ macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

สคริปต์จะดาวน์โหลดรุ่นสำหรับระบบของคุณ ตรวจสอบ checksum แบบ SHA-256 และวาง `ccam` ไว้ใน `~/.local/bin` ซึ่งเป็นโฟลเดอร์ที่ ccam ใช้สร้างคำสั่งลัดที่คุณเลือกในภายหลังด้วย เช่น `ccwork` หากโฟลเดอร์นั้นไม่อยู่ใน `PATH` ของคุณ สคริปต์จะพิมพ์บรรทัดที่ต้องเพิ่มลงใน `~/.bashrc` หรือ `~/.zshrc` ออกมาให้

บน Windows ให้รันคำสั่งนี้ใน PowerShell:

```powershell
irm https://ccam.dev/install.ps1 | iex
```

สคริปต์จะดาวน์โหลด `ccam_windows_<arch>.zip` ตรวจสอบ checksum แบบ SHA-256 วาง `ccam.exe` ไว้ใน `%USERPROFILE%\.local\bin` ซึ่งเป็นโฟลเดอร์เดียวกับที่ตัวติดตั้งของ Claude Code วาง `claude.exe` ไว้ หากโฟลเดอร์นั้นยังไม่อยู่ใน `PATH` ของผู้ใช้ สคริปต์จะเพิ่มให้ และคงรายการอย่าง `%USERPROFILE%\bin` ไว้ตามเดิม จากนั้นให้เปิดเทอร์มินัลใหม่

## ระบบที่รองรับ

- Linux บน x86-64 และ ARM64: Ubuntu, Debian, Fedora, Arch, Alpine และดิสทริบิวชันอื่น ๆ ไบนารีเป็นแบบ static และไม่มี dependency
- macOS บน Apple silicon และ Intel
- Windows 10 และ 11 บน x86-64 และ ARM64 ใน PowerShell, Windows Terminal หรือ Git Bash:
  - โฟลเดอร์ที่ใช้ร่วมกันทำงานได้โดยไม่ต้องเปิด Developer Mode
  - ไฟล์ที่ใช้ร่วมกัน (`settings.json`, `CLAUDE.md`, `history.jsonl`) ต้องเปิด Developer Mode (Settings > System > For developers) เพราะ Windows อนุญาตให้สร้างลิงก์ไปยังไฟล์ได้เฉพาะในโหมดนี้หรือโดยผู้ดูแลระบบ ดู [การใช้ร่วมกันบน Windows](/features/share#on-windows)
  - หากไม่ได้เปิด Developer Mode `ccam add` ครั้งแรกจะสร้างกลุ่ม `default` ที่มีเฉพาะเซิร์ฟเวอร์ MCP และส่วนขยาย และบอกวิธีเพิ่มรายการที่เหลือเมื่อเปิด Developer Mode แล้ว
- ภายใน WSL นั้น ccam ทำงานเหมือนบนระบบ Linux อื่น ๆ ทั่วไป

## วิธีติดตั้งอื่น ๆ

**แพ็กเกจสำหรับ Debian และ Ubuntu** ดาวน์โหลด `ccam_<version>_linux_amd64.deb` จาก[หน้ารุ่นที่เผยแพร่](https://github.com/jhyunwoo/ccam/releases) แล้วติดตั้ง:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

แพ็กเกจจะติดตั้ง `/usr/bin/ccam`

**จากซอร์สโค้ด** ใช้ Go 1.26 หรือใหม่กว่า:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
```

หากคุณติดตั้งจากแท็กรุ่น (`@latest` หรือ `@vX.Y.Z`) ccam จะรายงานรุ่นนั้น และ `ccam update` กับการเสนออัปเดตในแดชบอร์ดจะใช้ได้ ส่วนบิลด์จาก `@main` หรือจากซอร์สในเครื่องที่มีคอมมิตหลังแท็ก จะรายงานเป็น `dev` และ `ccam update` จะปฏิเสธบิลด์นั้น เว้นแต่คุณจะใส่ `--version`

**รุ่นที่ระบุ** ตั้งค่า `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

```powershell
$env:CCAM_VERSION = 'v0.1.0'; irm https://ccam.dev/install.ps1 | iex
```

## การเติมคำอัตโนมัติใน shell

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

```powershell
ccam completion powershell >> $PROFILE
```

การเติมคำอัตโนมัติครอบคลุมทุกคำสั่งและตัวเลือกของ `ccam` รวมถึงชื่อบัญชีด้วย

## ปลั๊กอิน Claude Code

ปลั๊กอิน ccam สอนให้ Claude Code ใช้ ccam ได้เอง: ตรวจว่าบัญชีใดยังเหลือโควตา เริ่ม worker `claude -p` แบบขนานบนบัญชีต่าง ๆ ด้วย `ccam pick` และอธิบายว่าทำไม skill หรือเซิร์ฟเวอร์ MCP จึงหายไปด้วย `ccam ext` ติดตั้งภายใน Claude Code:

```text
/plugin install ccam --marketplace jhyunwoo/ccam
```

ปลั๊กอินนี้เพิ่ม:

- skill หนึ่งตัวที่แนะนำคำสั่งที่เข้าสู่ระบบหรือย้ายไฟล์ แต่ไม่เคยรันคำสั่งเหล่านั้นเอง
- คำสั่ง `/ccam:usage`, `/ccam:ext` และ `/ccam:switch`
- hook สองตัว ([รายละเอียด](/features/next#inside-claude-code)) ตัวแรกแสดงการแจ้งเตือนในเทอร์มินัลเมื่อเซสชันหยุดเพราะถึงขีดจำกัดการใช้งาน บน Linux และ macOS เมื่อเปิด [`auto_switch`](/features/next#auto-switch) ไว้ ccam จะย้ายเซสชันไปยังอีกบัญชีหนึ่งที่ยังมีโควตาเหลือต่อจากนั้น อีกตัวหนึ่งทำให้ `/ccam:switch [account]` ส่งต่อการสนทนาไปยังอีกบัญชีหนึ่งได้

hook เหล่านี้รัน `ccam` ดังนั้น ccam ต้องอยู่ใน `PATH` ที่ Claude Code มองเห็น

## ตรวจสอบการติดตั้ง

```sh
ccam version
ccam doctor
```

`ccam doctor` ตรวจสอบว่าพบ Claude Code และทุกโฟลเดอร์บัญชี ลิงก์ และการเข้าสู่ระบบอยู่ในสภาพเรียบร้อย ทำต่อได้ที่[คู่มือเริ่มต้นใช้งาน](/docs/quick-start)

## การอัปเดต

```sh
ccam update --check
ccam update
```

`ccam update --check` จะบอกเพียงว่ามีรุ่นที่ใหม่กว่าหรือไม่ ส่วน `ccam update` จะดาวน์โหลดรุ่นสำหรับระบบของคุณ ตรวจสอบ checksum แบบ SHA-256 เทียบกับ `checksums.txt` ของรุ่นนั้น แล้วแทนที่ไฟล์ `ccam` ในขั้นตอนเดียว หาก checksum ไม่ตรงกัน จะไม่มีสิ่งใดถูกเปลี่ยนแปลง คำสั่งลัดอย่าง `ccwork` เป็นลิงก์ไปยัง ccam จึงใช้รุ่นใหม่ได้ทันที

```console
$ ccam update --check
ccam 0.3.0 is available (you have 0.2.0). Run `ccam update`.
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.3.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.2.0 → 0.3.0 (/home/me/.local/bin/ccam)
$ ccam update
ccam 0.3.0 is up to date.
```

แดชบอร์ด (`ccam` ที่ไม่มีอาร์กิวเมนต์ ในเทอร์มินัล) ก็ตรวจหารุ่นใหม่ทุกครั้งที่เปิดขึ้นเช่นกัน และเมื่อมีรุ่นใหม่จะเสนอให้ติดตั้ง: กด `U` แล้วยืนยัน ไม่มีคำสั่งอื่นใดตรวจหาการอัปเดต และ ccam ไม่เคยอัปเดตโดยไม่ถามก่อน หากต้องการติดตั้งรุ่นที่ระบุ รวมถึงรุ่นที่เก่ากว่า ให้ใส่ `--version` ส่วน `ccam --dry-run update` จะแสดงว่าจะดาวน์โหลดอะไรและแทนที่ไฟล์ใด โดยไม่เปลี่ยนแปลงสิ่งใดจริง

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

การรันสคริปต์ติดตั้งอีกครั้งก็จะอัปเดตสำเนาใน `~/.local/bin` เช่นกัน หากคุณติดตั้ง ccam ด้วย `go install` จากแท็กรุ่น คุณจะรัน `ccam update` หรือรันคำสั่ง `go install` เดิมอีกครั้งก็ได้

**Windows** `ccam update` จะดาวน์โหลดไฟล์ `.zip` สำหรับ Windows เนื่องจากบน Windows ไม่สามารถแทนที่โปรแกรมที่กำลังทำงานอยู่ได้ ไบนารีเดิมจึงถูกเปลี่ยนชื่อเป็น `ccam.exe.old` ซึ่งการอัปเดตครั้งถัดไปจะลบทิ้ง ขณะที่ยังมีเซสชันทำงานจาก `ccam.exe.old` ก่อนหน้าอยู่ ชื่อจะเป็น `ccam.exe.old-<time>` คำสั่งลัดที่เป็น hard link จะถูกต่ออายุทันทีหลังจากนั้น การรัน `install.ps1` อีกครั้งก็ทำงานแบบเดียวกัน แม้ขณะที่ยังมีเซสชัน Claude Code ทำงานจากไบนารีเดิมอยู่ และจะพิมพ์ `OK 2 shortcuts now run the new ccam` เมื่อต่ออายุคำสั่งลัดที่เป็น hard link

**แพ็กเกจสำหรับ Debian และ Ubuntu** `ccam update` จะไม่แทนที่ `/usr/bin/ccam` ที่ติดตั้งด้วย apt แต่จะพิมพ์ขั้นตอนต่อไปนี้ออกมาแทน ดาวน์โหลด `ccam_<version>_linux_amd64.deb` รุ่นใหม่จาก[หน้ารุ่นที่เผยแพร่](https://github.com/jhyunwoo/ccam/releases) แล้วติดตั้งทับรุ่นเดิม:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## ถอนการติดตั้ง

`ccam uninstall` จะลบ ccam และทุกอย่างที่ ccam สร้างขึ้น โดยจะแสดงรายการสิ่งที่จะถูกลบก่อน แล้วถามยืนยัน:

```sh
ccam --dry-run uninstall   # only show what would be removed
ccam uninstall
```

- ทุกบัญชีจะถูกลบแบบเดียวกับที่ `ccam remove` ทำ บัญชีที่อยู่ในโฟลเดอร์ของ ccam จะถูกออกจากระบบก่อน แล้วลบไปพร้อมกับโฟลเดอร์ รวมถึงประวัติการสนทนาและการตั้งค่าของบัญชีนั้นด้วย
- โฟลเดอร์ `~/.claude` ที่คุณ[นำเข้า](/docs/commands)ไว้จะยังอยู่ในตำแหน่งเดิมเพื่อใช้กับ `claude` แบบปกติ หากโฟลเดอร์นั้นอยู่ในกลุ่ม จะได้สำเนาของทุกอย่างที่เคยใช้ร่วมกันกลับคืนมา แถบสถานะของ ccam ใน `settings.json` ของโฟลเดอร์นั้นจะหลีกทางให้แถบสถานะเดิมที่ ccam ครอบไว้ หากมี และหากโฟลเดอร์นั้นมีปลั๊กอิน ccam อยู่ ccam จะพิมพ์คำสั่ง `/plugin uninstall` ไว้ให้รันใน Claude Code
- คำสั่งลัด, `~/.config/ccam`, `~/.local/share/ccam` (รวมถึงถังขยะของ ccam) และ `~/.local/state/ccam` จะถูกลบ จากนั้นจึงลบคำสั่ง `ccam` เอง

ccam จะปฏิเสธขณะที่มีบัญชีใดกำลังทำงานอยู่ ให้ปิดเซสชัน Claude Code ของบัญชีนั้นก่อน ส่วน `/usr/bin/ccam` ที่ติดตั้งด้วย apt จะยังอยู่ ให้ลบภายหลังด้วย `sudo apt remove ccam` บน Windows โปรแกรมที่กำลังทำงานไม่สามารถลบตัวเองได้ ดังนั้น `ccam uninstall` จะจบด้วยการพิมพ์คำสั่ง `Remove-Item` สำหรับ `ccam.exe`
