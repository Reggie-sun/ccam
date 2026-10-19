---
title: การจัดการบัญชี
description: "เพิ่ม เข้าสู่ระบบ แสดงรายการ เปลี่ยนชื่อ ลบ และนำเข้าบัญชี Claude Code ด้วย ccam add, ccam login, ccam list, ccam rename, ccam remove และ ccam import"
updated: 2026-10-19
order: 2
translation: machine
sourceHash: 99c931574262
---

บัญชี ccam แต่ละบัญชีคือโฟลเดอร์การตั้งค่า Claude Code แยกกันที่มีการเข้าสู่ระบบของตัวเอง คำสั่งเหล่านี้ใช้สร้าง ลงชื่อเข้าใช้ แสดงรายการ เปลี่ยนชื่อ และลบบัญชี รวมถึงลงทะเบียนโฟลเดอร์ที่คุณมีอยู่แล้ว

## การใช้งาน

```sh
ccam add                                   # wizard on a terminal: name, shared items, shortcut, login
ccam add work --shortcut ccwork            # no wizard: log in, then create the ccwork command
ccam add side --no-login --share history=team
ccam login side                            # run claude auth login inside the account
ccam logout side                           # remove the login token
ccam list                                  # table of accounts; ccam ls is the same
ccam list --json                           # machine-readable
ccam use 2                                 # make slot 2 the default account
ccam rename side personal
ccam remove personal                       # asks, then moves the folder to the trash
ccam import                                # find ~/.claude and ~/.claude-* and ask about each
ccam import ~/.claude-ops --name ops
ccam completion zsh > "${fpath[1]}/_ccam"
```

เมื่อคำสั่งรับ `<name|number>` คุณระบุชื่อบัญชีหรือหมายเลขสล็อตก็ได้ ซึ่งก็คือคอลัมน์ `#` ของ `ccam list`

`ccam add` บนเทอร์มินัลโดยไม่ใส่ตัวเลือกจะเปิดตัวช่วยตั้งค่า โดยจะถามชื่อ จากนั้นแสดงรายการตรวจสอบ (checklist) ของรายการที่ใช้ร่วมกัน (เลือกไว้ทุกอย่างยกเว้นประวัติการสนทนา) แล้วถาม `Create a shortcut command for work? [y/N]` และหากตอบ `y` จะถาม `Shortcut name [ccwork]:` ต่อ การใส่ `--share`, `--shortcut` หรือ `--no-login` ตัวใดตัวหนึ่ง หรือการส่งข้อมูลผ่านไปป์แทนเทอร์มินัล จะข้ามตัวช่วยตั้งค่า ในกรณีนั้นต้องระบุชื่อ และจะไม่มีการสร้างคำสั่งลัดเว้นแต่ `--shortcut` จะระบุชื่อไว้

`ccam list` แสดงสล็อต ชื่อ คำสั่งลัด อีเมล รายการที่ใช้ร่วมกัน และเวลาที่ใช้ล่าสุด โดยมี `*` อยู่หน้าบัญชีเริ่มต้น:

```console
$ ccam list
   #  Name      Shortcut  Email             Shared                        Last used
*  1  work      ccwork    user@example.com  mcp settings md ext           -
   2  personal  -         me@side.dev       settings md ext history@team  -
```

เมื่อใช้ `--json` แต่ละบัญชีจะเป็นออบเจ็กต์ที่มี `slot`, `shortcut`, `name`, `email`, `default`, `dir`, `last_used` (ไม่มีหากบัญชีไม่เคยถูกเริ่ม) และ `share` โดย `share` จะจับคู่แต่ละรายการ (`mcp`, `settings`, `memory`, `extensions`, `history`) กับชื่อกลุ่มหรือ `"local"`

## ตัวเลือก

| ตัวเลือก | สิ่งที่ทำ | ค่าเริ่มต้น |
| --- | --- | --- |
| `ccam add --no-login` | สร้างบัญชีโดยไม่เข้าสู่ระบบ แล้วค่อยเข้าสู่ระบบภายหลังด้วย `ccam login <name>` | ปิด |
| `ccam add --share <item=group>` | ใช้รายการร่วมกับกลุ่ม หรือใช้ `item=local` เพื่อเก็บไว้เฉพาะบัญชี ใส่ซ้ำได้หลายครั้ง | `mcp`, `settings`, `memory`, `extensions` อยู่ใน `default` ส่วน `history` เก็บไว้เฉพาะบัญชี |
| `ccam add --shortcut <command>` | สร้างคำสั่งลัดนี้ใน `~/.local/bin` ด้วย | ไม่มี |
| `ccam list --json` | พิมพ์รายการบัญชีเป็น JSON | ปิด |
| `ccam remove -y`, `--yes` | ลบโดยไม่ถาม | ปิด |
| `ccam import --name <name>` | ชื่อบัญชีสำหรับโฟลเดอร์ที่คุณระบุ | ได้มาจากชื่อโฟลเดอร์ |
| `ccam import -y`, `--yes` | นำเข้าทุกโฟลเดอร์ที่พบโดยไม่ถาม | ปิด |
| `--dry-run` | `ccam list` และ `ccam completion` ซึ่งไม่เปลี่ยนแปลงสิ่งใดจะรองรับตัวเลือกนี้ ส่วน `add`, `login`, `logout`, `use`, `rename`, `remove` และ `import` จะปฏิเสธด้วยข้อความ `` --dry-run is not supported by `ccam add`; nothing was changed `` (โดยระบุชื่อคำสั่งนั้น) | ปิด |

`ccam login`, `ccam logout`, `ccam use` และ `ccam rename` ไม่มีตัวเลือกของตัวเอง `ccam completion` รับ `bash` หรือ `zsh` ค่าอื่นจะล้มเหลวด้วย `only bash and zsh are supported`

ชื่อบัญชียาว 1 ถึง 32 อักขระ ประกอบด้วย `a-z`, `0-9`, `-` และ `_` และต้องขึ้นต้นด้วยตัวอักษรหรือตัวเลข ชื่อต้องไม่เป็นตัวเลขล้วน เพราะ `ccam 2` หมายถึงสล็อต 2 และชื่อต้องไม่ซ้ำกับคำสั่งของ ccam (`add`, `list`, `ls`, `run`, `use`, `next`, `usage`, `u`, `share`, `mcp`, `login`, `logout`, `rename`, `remove`, `rm`, `import`, `migrate`, `doctor`, `shims`, `shortcut`, `completion`, `help`, `version`, `statusline`, `update`) เพื่อให้ `ccam <name>` ไม่กำกวม รายการนี้คือ `config.Reserved` ใน `internal/config`

## วิธีการทำงาน

`ccam add work --shortcut ccwork` ทำตามขั้นตอนต่อไปนี้ใน `addProfile` (`internal/cli/profiles.go`):

<!-- motion:accounts-add -->

1. `config.ValidateName` ตรวจสอบชื่อ และ `shim.Check` ตรวจสอบคำสั่งลัด ก่อนที่จะสร้างสิ่งใด
2. `profile.Add` สร้าง `~/.local/share/ccam/profiles/work/` ด้วยโหมด `0700` และให้หมายเลขสล็อตที่ว่างต่ำสุดแก่บัญชี (`Config.NextSlot`) ดังนั้นหมายเลขของบัญชีที่ลบไปแล้วจะถูกนำกลับมาใช้ บัญชีแรกจะเป็นบัญชีเริ่มต้น
3. `applyShare` ลิงก์โฟลเดอร์ใหม่ที่ยังว่างเข้ากับกลุ่มที่ใช้ร่วมกันด้วย `share.Build` แล้วตามด้วย `Apply` เช่น ลิงก์ `settings.json` ไปยัง `groups/default/settings.json` หากขั้นตอนนี้ล้มเหลว โฟลเดอร์ใหม่จะถูกลบ เพราะโฟลเดอร์ที่ค้างอยู่จะทำให้การลองใหม่ทุกครั้งล้มเหลวด้วย `folder already exists`
4. บัญชีจะถูกบันทึกลงใน `~/.config/ccam/config.toml` ผ่าน `fsutil.AtomicWrite`: เขียนไฟล์ชั่วคราว เรียก `fsync` แล้วจึง `rename` ดังนั้นการแครชจะไม่ทิ้งการตั้งค่าไว้ครึ่ง ๆ กลาง ๆ
5. `runLogin` เริ่ม `claude auth login` โดยตั้ง `CLAUDE_CONFIG_DIR` เป็นโฟลเดอร์ของบัญชี หากเข้าสู่ระบบไม่สำเร็จ บัญชีจะยังคงอยู่ และ ccam จะพิมพ์ข้อผิดพลาดพร้อมบอกให้คุณรัน `ccam login work` ภายหลัง
6. `profile.ReadEmail` ดึง `oauthAccount.emailAddress` ออกจาก `profiles/work/.claude.json` ไฟล์นั้นอาจมีขนาดหลายเมกะไบต์ `TopLevelValue` จึงสแกนหาเพียงคีย์เดียวแทนการพาร์สทั้งไฟล์
7. `shim.Install` ลิงก์ `~/.local/bin/ccwork` ไปยังไบนารีของ ccam ซึ่งจะค้นหา `argv[0]` ใน `config.toml` เมื่อเริ่มทำงาน

การเข้าสู่ระบบอาจใช้เวลาหลายนาที และ ccam อีกตัวอาจบันทึกการตั้งค่าในระหว่างนั้น คำสั่งที่ต้องรอจึงบันทึกผ่าน `env.update` ซึ่งโหลด `config.toml` ใหม่ ใช้การเปลี่ยนแปลงเพียงรายการเดียว (ในที่นี้คืออีเมล) แล้วบันทึก `ccam remove` ก็โหลดใหม่หลังจากถามด้วยเหตุผลเดียวกัน

`ccam login` ทำขั้นตอนที่ 5 และ 6 สำหรับบัญชีที่มีอยู่แล้ว `ccam logout` รัน `claude auth logout` ในบัญชีนั้นและล้างอีเมลที่บันทึกไว้ ส่วน `ccam use` เพียงเขียนค่า `default` ใน `config.toml` ใหม่

### การลงชื่อเข้าใช้บนเครื่องระยะไกล

เมื่อใช้งานผ่าน SSH ส่วนที่ทำงานในเบราว์เซอร์ของ `claude auth login` อาจใช้ไม่ได้ Claude Code รอการลงชื่อเข้าใช้บนเซิร์ฟเวอร์ callback ที่ `http://localhost:<port>/callback` บนเครื่องเซิร์ฟเวอร์ แต่เบราว์เซอร์ที่ทำการลงชื่อเข้าใช้จนเสร็จทำงานอยู่บนคอมพิวเตอร์ของคุณ ซึ่งที่อยู่นั้นไม่ได้ชี้ไปที่ใด เบราว์เซอร์จึงแสดงข้อผิดพลาด เช่น `ERR_CONNECTION_REFUSED` อย่างไรก็ตาม ที่อยู่นั้นยังมีรหัสแบบใช้ครั้งเดียวอยู่ ccam จึงส่งต่อรหัสให้ได้:

<!-- motion:accounts-relay -->

`runLogin` ไม่ได้ส่งเทอร์มินัลของคุณให้ Claude Code โดยตรง แต่เชื่อมอินพุตมาตรฐานของ Claude Code เข้ากับไปป์ และ `relayLogin` จะคัดลอกแต่ละบรรทัดที่คุณพิมพ์ลงไป ดังนั้นรหัสที่วางในพรอมต์ `Paste code here` ของ Claude Code ก็ยังส่งถึงตามปกติ เมื่อบรรทัดหนึ่งเป็นที่อยู่ callback `callbackURL` จะยอมรับก็ต่อเมื่อโฮสต์เป็น `localhost`, `127.0.0.1` หรือ `[::1]` มีการระบุพอร์ต พาธเป็น `/callback` และมีทั้ง `code` และ `state` โดยจะละ `http://` ไว้ก็ได้ จากนั้น `deliverCallback` จะส่งคำขอไปยังที่อยู่นั้นบนเซิร์ฟเวอร์ โดยจำกัดเวลา 10 วินาทีและไม่ตาม redirect เหมือนที่เบราว์เซอร์จะทำ คำขอจะไปถึงเซิร์ฟเวอร์ callback ของ `claude auth login` ที่กำลังรออยู่ (เส้นประในแผนภาพ) Claude Code ได้รับรหัสและบันทึกการเข้าสู่ระบบ ซึ่งบน Linux จะอยู่ใน `.credentials.json` ในโฟลเดอร์ของบัญชี บรรทัดที่วางไม่มีทางทำให้ ccam ติดต่อโฮสต์อื่นได้

ccam พิมพ์ข้อความนี้ก่อนเริ่มเข้าสู่ระบบ:

```text
On a remote machine? If the browser ends on a http://localhost:…/callback page that does not load, paste that address here.
```

เมื่อสำเร็จจะพิมพ์ `→ Passed the sign-in to Claude Code on this machine.`

ขณะที่การเข้าสู่ระบบรออินพุตบนเทอร์มินัล ccam จะปิด echo และโหมดบรรทัดของเทอร์มินัล และแสดง `*` แทนแต่ละตัวอักษรที่คุณพิมพ์ ดังนั้นรหัสหรือที่อยู่ callback ที่วางจะไม่ปรากฏบนหน้าจอ `readLine` จัดการแก้ไขบรรทัดเอง: `Enter` ส่งบรรทัด `Backspace` ลบหนึ่งตัวอักษร และ `Ctrl+D` บนบรรทัดว่างจะจบอินพุต เทอร์มินัลจะถูกคืนค่าเมื่อการเข้าสู่ระบบจบ และเมื่อได้รับ `Ctrl+C`, `SIGTERM` หรือ `SIGHUP` ด้วย ซึ่งยังคงจบ ccam ตามเดิม หากอินพุตมาตรฐานและเอาต์พุตข้อผิดพลาดมาตรฐานไม่ได้เชื่อมกับเทอร์มินัล บรรทัดจะถูกคัดลอกต่อไปตามที่เข้ามา

### การเปลี่ยนชื่อ การลบ และการนำเข้า

`ccam rename` และ `ccam remove` ทำงานกับโฟลเดอร์ของบัญชี และกับสิ่งที่ ccam เก็บไว้ภายใต้ชื่อบัญชีในที่อื่น `nameState` ระบุพาธเหล่านั้น: `~/.local/state/ccam/last-used/<name>`, `~/.local/state/ccam/mcp/<name>` และแคชการใช้งาน `~/.local/state/ccam/usage/<name>.json` และ `<name>.statusline.json` โดย `ccam remove work` จัดการตามลำดับนี้:

<!-- motion:accounts-remove -->

`ccam remove` จะปฏิเสธก่อนหากบัญชีกำลังทำงาน แล้วจึงถาม หากบัญชีใช้ส่วนขยายร่วมกัน `RewritePluginPaths` จะเขียน `installed_plugins.json` และ `known_marketplaces.json` ของกลุ่มใหม่ เพราะ Claude Code เก็บพาธปลั๊กอินแบบสัมบูรณ์ซึ่งอาจผ่านโฟลเดอร์ของบัญชีนี้ `profile.Trash` ย้ายโฟลเดอร์ไปยัง `~/.local/share/ccam/trash/<name>-<time>/` บัญชีจะถูกนำออกจาก `config.toml` พาธใน `nameState` จะถูกลบเพื่อไม่ให้บัญชีที่สร้างภายหลังด้วยชื่อเดียวกันได้รับสืบทอดไป และลิงก์คำสั่งลัดจะถูกลบ หากต้องการกู้คืน ให้ย้ายโฟลเดอร์กลับมาจากถังขยะ

`ccam rename` เปลี่ยนชื่อโฟลเดอร์ เขียนพาธปลั๊กอินใน `plugins/` ของบัญชีใหม่จากโฟลเดอร์เดิมเป็นโฟลเดอร์ใหม่ คงสถานะบัญชีเริ่มต้นไว้ และย้ายพาธใน `nameState` ไปใช้ชื่อใหม่ คำสั่งลัดยังคงชื่อเดิม หากการเขียนพาธปลั๊กอินใหม่ล้มเหลว โฟลเดอร์จะถูกเปลี่ยนชื่อกลับ เนื่องจากยังไม่ได้บันทึกการตั้งค่า

`ccam import` ลงทะเบียนโฟลเดอร์ไว้ในตำแหน่งเดิม โดยเก็บเป็น `dir` ใน `config.toml` ccam ไม่ย้ายหรือคัดลอกโฟลเดอร์เลย เพราะ refresh token ที่ถูกคัดลอกจะทำให้สำเนาแต่ละชุดทำให้กันและกันออกจากระบบ หากไม่ระบุโฟลเดอร์ ccam จะค้นหาโฟลเดอร์ `~/.claude` และ `~/.claude-*` ที่ยังไม่ได้ลงทะเบียนแล้วถามทีละโฟลเดอร์ ชื่อบัญชีได้มาจากชื่อโฟลเดอร์ (`~/.claude` เป็น `main`, `~/.claude-work` เป็น `work`, `~/.claude-2` เป็น `acct2`) และโฟลเดอร์ที่หาชื่อไม่ได้จะแสดงพร้อมคำสั่ง `--name` ที่ต้องใช้ หากไม่มี `CLAUDE_CONFIG_DIR` Claude Code จะเก็บ `.claude.json` ไว้ในโฟลเดอร์ home ข้าง ๆ `~/.claude` ดังนั้นเมื่อคุณนำเข้า `~/.claude` และโฟลเดอร์นั้นไม่มี `.claude.json` ของตัวเอง ccam จะสร้างลิงก์ `~/.claude/.claude.json` ไปยัง `~/.claude.json` จากนั้น `claude` แบบปกติและ `ccam main` จะเห็นบัญชีเดียวกัน บัญชีที่นำเข้าจะไม่ใช้สิ่งใดร่วมกันจนกว่าคุณจะรัน `ccam share` และจะไม่มีการสร้างคำสั่งลัด

`ccam completion bash|zsh` พิมพ์สคริปต์ completion ที่เติมชื่อบัญชีให้ด้วย

## ความปลอดภัยและกรณีพิเศษ

- `ccam rename` และ `ccam remove` จะปฏิเสธขณะที่บัญชีกำลังทำงาน (`refuseRunning`) ด้วยข้อความ `work is running; quit its Claude Code session first` เพราะเซสชันยังคงเขียนลงในโฟลเดอร์
- `ccam remove` กับบัญชีที่นำเข้ามาจะเพียงยกเลิกการลงทะเบียน: `✓ main unregistered (folder ... left in place)` โดย `~/.claude` ของคุณจะไม่ถูกย้ายไปถังขยะ
- `ccam add` จะปฏิเสธชื่อที่มีโฟลเดอร์อยู่แล้วใน `profiles/` และแนะนำให้ใช้ `ccam import`
- `ccam import` จะปฏิเสธโฟลเดอร์ที่ลงทะเบียนไว้แล้ว พาธที่ไม่ใช่โฟลเดอร์ และชื่อที่ถูกใช้ไปแล้ว
- `ccam logout` และ `ccam login` ไม่ตรวจสอบว่าบัญชีกำลังทำงานอยู่หรือไม่
- หากตัวส่งต่อติดต่อ Claude Code ไม่ได้ เช่น หลังจากการเข้าสู่ระบบสิ้นสุดไปแล้ว จะพิมพ์ `ccam: could not pass the sign-in to Claude Code:` พร้อมเหตุผล ข้อผิดพลาด HTTP มักหมายความว่ารหัสหมดอายุแล้ว ให้เริ่มเข้าสู่ระบบใหม่

## หัวข้อที่เกี่ยวข้อง

- [การเริ่ม Claude Code](/features/launch)
- [คำสั่งลัด](/features/shortcuts)
- [ใช้เซิร์ฟเวอร์ MCP การตั้งค่า และประวัติร่วมกัน](/docs/sharing)
- [Claude Code หลายบัญชี](/guides/multiple-claude-code-accounts)
