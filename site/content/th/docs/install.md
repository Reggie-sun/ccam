---
title: ติดตั้ง ccam
description: "ติดตั้ง ccam บน Linux หรือ macOS ด้วยคำสั่งเดียว จากแพ็กเกจ .deb หรือจากซอร์สโค้ด โดยต้องติดตั้ง Claude Code ก่อน"
updated: 2026-10-19
order: 1
translation: machine
sourceHash: c7f92dfa554a
---

รันคำสั่งนี้ในเทอร์มินัลบน Linux หรือ macOS:

```sh
curl -fsSL https://ccam.dev/install.sh | sh
```

สคริปต์จะดาวน์โหลดรุ่นสำหรับระบบของคุณ ตรวจสอบ checksum แบบ SHA-256 และวางคำสั่งสองคำสั่งไว้ใน `~/.local/bin` ได้แก่ `ccam` และ alias แบบสั้นของมันคือ `ccs` หากโฟลเดอร์นั้นไม่อยู่ใน `PATH` ของคุณ สคริปต์จะพิมพ์บรรทัดที่ต้องเพิ่มลงใน `~/.bashrc` หรือ `~/.zshrc` ออกมาให้

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

แพ็กเกจจะติดตั้ง `/usr/bin/ccam` และลิงก์ `/usr/bin/ccs`

**จากซอร์สโค้ด** ใช้ Go 1.26 หรือใหม่กว่า:

```sh
go install github.com/jhyunwoo/ccam/cmd/ccam@latest
ln -s ccam "$(go env GOPATH)/bin/ccs"
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

การเติมคำอัตโนมัติครอบคลุมทั้ง `ccam` และ `ccs` รวมถึงชื่อบัญชีด้วย

## ตรวจสอบการติดตั้ง

```sh
ccam version
ccam doctor
```

`ccam doctor` ตรวจสอบว่าพบ Claude Code และทุกโฟลเดอร์บัญชี ลิงก์ และการเข้าสู่ระบบอยู่ในสภาพเรียบร้อย ทำต่อได้ที่[คู่มือเริ่มต้นใช้งาน](/docs/quick-start)

## ถอนการติดตั้ง

ลบคำสั่งทั้งสองและคำสั่งลัดของบัญชี จากนั้นลบโฟลเดอร์ข้อมูลหากคุณไม่ต้องการบัญชีเหล่านั้นอีกต่อไป:

```sh
rm ~/.local/bin/ccam ~/.local/bin/ccs ~/.local/bin/claude[0-9]*
rm -r ~/.config/ccam ~/.local/share/ccam ~/.local/state/ccam
```

การลบ `~/.local/share/ccam` จะลบข้อมูลเข้าสู่ระบบและประวัติการสนทนาที่เก็บไว้ในนั้นด้วย ส่วนโฟลเดอร์ `~/.claude` ที่คุณ[นำเข้า](/docs/commands)ไว้จะยังอยู่ในตำแหน่งเดิม
