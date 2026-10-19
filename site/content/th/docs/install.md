---
title: ติดตั้ง ccam
description: "ติดตั้ง ccam บน Linux หรือ macOS ด้วยคำสั่งเดียว จากแพ็กเกจ .deb หรือจากซอร์สโค้ด และอัปเดตให้เป็นรุ่นล่าสุดเสมอด้วย ccam update โดยต้องติดตั้ง Claude Code ก่อน"
updated: 2026-10-19
order: 1
translation: machine
sourceHash: bde486b82b2e
---

รันคำสั่งนี้ในเทอร์มินัลบน Linux หรือ macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

สคริปต์จะดาวน์โหลดรุ่นสำหรับระบบของคุณ ตรวจสอบ checksum แบบ SHA-256 และวาง `ccam` ไว้ใน `~/.local/bin` ซึ่งเป็นโฟลเดอร์ที่ ccam ใช้สร้างคำสั่งลัดที่คุณเลือกในภายหลังด้วย เช่น `ccwork` หากโฟลเดอร์นั้นไม่อยู่ใน `PATH` ของคุณ สคริปต์จะพิมพ์บรรทัดที่ต้องเพิ่มลงใน `~/.bashrc` หรือ `~/.zshrc` ออกมาให้

ccam เป็นตัวเริ่ม Claude Code ให้คุณ ดังนั้นให้ติดตั้ง [Claude Code](https://docs.claude.com/en/docs/claude-code/overview) ก่อน และตรวจสอบให้แน่ใจว่า `claude --version` ทำงานได้

## ระบบที่รองรับ

- Linux บน x86-64 และ ARM64: Ubuntu, Debian, Fedora, Arch, Alpine และดิสทริบิวชันอื่น ๆ ไบนารีเป็นแบบ static และไม่มี dependency
- macOS บน Apple silicon และ Intel
- ไม่รองรับ Windows โดยตรง ภายใน WSL นั้น ccam ทำงานเหมือนบนระบบ Linux ทั่วไป

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

**รุ่นที่ระบุ** ตั้งค่า `CCAM_VERSION`:

```sh
curl -fsSL https://ccam.dev/install.sh | CCAM_VERSION=v0.1.0 sh
```

## การเติมคำอัตโนมัติใน shell

```sh
ccam completion bash > ~/.local/share/bash-completion/completions/ccam
ccam completion zsh > "${fpath[1]}/_ccam"
```

การเติมคำอัตโนมัติครอบคลุมทุกคำสั่งและตัวเลือกของ `ccam` รวมถึงชื่อบัญชีด้วย

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

ccam ไม่ตรวจหาอัปเดตเองโดยอัตโนมัติ และจะค้นหารุ่นใหม่เฉพาะเมื่อคุณรันคำสั่งนี้เท่านั้น หากต้องการติดตั้งรุ่นที่ระบุ รวมถึงรุ่นที่เก่ากว่า ให้ใส่ `--version` ส่วน `ccam --dry-run update` จะแสดงว่าจะดาวน์โหลดอะไรและแทนที่ไฟล์ใด โดยไม่เปลี่ยนแปลงสิ่งใดจริง

```sh
ccam update --version v0.2.0
ccam --dry-run update
```

การรันสคริปต์ติดตั้งอีกครั้งก็จะอัปเดตสำเนาใน `~/.local/bin` เช่นกัน หากคุณติดตั้ง ccam ด้วย `go install` ให้รันคำสั่งเดิมอีกครั้ง

**แพ็กเกจสำหรับ Debian และ Ubuntu** `ccam update` จะไม่แทนที่ `/usr/bin/ccam` ที่ติดตั้งด้วย apt แต่จะพิมพ์ขั้นตอนต่อไปนี้ออกมาแทน ดาวน์โหลด `ccam_<version>_linux_amd64.deb` รุ่นใหม่จาก[หน้ารุ่นที่เผยแพร่](https://github.com/jhyunwoo/ccam/releases) แล้วติดตั้งทับรุ่นเดิม:

```sh
sudo apt install ./ccam_*_linux_amd64.deb
```

## ถอนการติดตั้ง

ลบคำสั่งลัดของบัญชีก่อน ขณะที่ ccam ยังติดตั้งอยู่ `ccam list` จะแสดงคำสั่งลัดในคอลัมน์ Shortcut ลบทีละรายการด้วย `ccam shortcut <account> --remove` หรือลบลิงก์เหล่านั้นออกจาก `~/.local/bin` ด้วยตัวเอง จากนั้นลบคำสั่ง ccam และลบโฟลเดอร์ข้อมูลหากคุณไม่ต้องการบัญชีเหล่านั้นอีกต่อไป:

```sh
ccam list
ccam shortcut work --remove   # once for each account that has a shortcut
rm ~/.local/bin/ccam
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

การลบ `~/.local/share/ccam` จะลบข้อมูลเข้าสู่ระบบและประวัติการสนทนาที่เก็บไว้ในนั้นด้วย ส่วนโฟลเดอร์ `~/.claude` ที่คุณ[นำเข้า](/docs/commands)ไว้จะยังอยู่ในตำแหน่งเดิม
