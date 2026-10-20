---
title: การใช้ร่วมกันด้วยกลุ่ม
description: "กลุ่มของ ccam ใช้รายการที่เลือกร่วมกันระหว่างบัญชีในกลุ่มผ่านโฟลเดอร์เดียว ส่วน ccam group ย้ายข้อมูลตามแผน พร้อมการสำรองและตัวเลือกเมื่อไฟล์ขัดแย้งกัน"
updated: 2026-10-20
order: 4
translation: machine
sourceHash: "0c0338362ba2"
---

ทุกบัญชีมีโฟลเดอร์ของตัวเอง แต่สิ่งที่อยู่ในโฟลเดอร์ส่วนใหญ่ไม่จำเป็นต้องแยกกัน **กลุ่ม**เป็นตัวกำหนดว่าจะใช้อะไรร่วมกันและใครใช้: ทุกบัญชีในกลุ่มจะลิงก์รายการเหล่านั้นไปยังโฟลเดอร์เดียว `~/.local/share/ccam/groups/<group>/` ทุกบัญชีจึงอ่านและเขียนสำเนาเดียวกัน โทเค็นการเข้าสู่ระบบจะไม่มีวันถูกใช้ร่วมกัน

บัญชีหนึ่งอยู่ได้หลายกลุ่ม แต่แต่ละรายการมาจากกลุ่มเดียวเท่านั้น การตั้งค่าที่พบบ่อยคือกลุ่มหนึ่งสำหรับการตั้งค่าและส่วนขยาย และอีกกลุ่มหนึ่งสำหรับประวัติการสนทนาระหว่างสองบัญชี:

```toml
[[group]]
name = 'default'
items = ['mcp', 'settings', 'memory', 'extensions']
members = ['personal', 'work', 'side']

[[group]]
name = 'team'
items = ['history']
members = ['work', 'side']
```

บัญชีแรกที่คุณเพิ่มจะสร้าง `default` และบัญชีที่เพิ่มภายหลังจะเข้าร่วมกลุ่มนี้ เว้นแต่คุณจะระบุเป็นอย่างอื่น

## การใช้งาน

```sh
ccam group                                        # the groups, what they share, their accounts
ccam group show default                           # items, accounts, folder and MCP servers
ccam group add team --items history --members work,side
ccam group set default --add history              # every account in default moves together
ccam group set default --remove history --fresh   # stop sharing; accounts start without it
ccam group set team --items history,memory --from work   # work's CLAUDE.md goes in first
ccam group join team personal                     # personal starts sharing team's items
ccam group join team personal --replace           # leave the group that shares them now
ccam group leave team side                        # side keeps a copy of team's items
ccam group rename team pair                       # the folder keeps its name
ccam group remove pair --leave                    # its accounts keep a copy; folder to the trash
ccam share work                                   # what work shares, and through which group
ccam --dry-run group set default --add history    # print the plan, change nothing
```

```console
$ ccam group
Group    Shares                          Accounts
default  mcp settings memory extensions  personal, work, side
team     history                         work, side

$ ccam share work
Sharing for work
  mcp         MCP servers                        group default
  settings    settings.json                      group default
  memory      global CLAUDE.md                   group default
  extensions  agents, commands, skills, plugins  group default
  history     conversation history               group team
Groups: default, team. Change them with `ccam group`.
```

ทุกคำสั่งที่ย้ายข้อมูลจะพิมพ์แผน ถามเกี่ยวกับแต่ละไฟล์ที่ต่างจากสำเนาของกลุ่ม และถาม `Proceed? [y/N]` ก่อนเปลี่ยนแปลงสิ่งใด:

```console
$ ccam group set default --add history
Plan (paths in /home/me/.local/share/ccam):
personal: +history (group default)
  move     projects, history.jsonl: profiles/personal → groups/default (first content of the group)
  link     projects, history.jsonl: profiles/personal → groups/default
work: +history (group default)
  merge    new files in projects: profiles/work → groups/default
  conflict profiles/work/history.jsonl ↔ groups/default/history.jsonl: merge by line
  back up  projects, history.jsonl: profiles/work → profiles/work/.ccam-backup/20261020-153012
  link     projects, history.jsonl: profiles/work → groups/default
Proceed? [y/N]
```

การเปลี่ยนแปลงแบบเดียวกันนี้อยู่ในแท็บ **Groups** ของ[แดชบอร์ด](/features/dashboard) และการกด `s` บนบัญชีหนึ่งจะใช้เลือกกลุ่มของบัญชีนั้น

## รายการ

| รายการ | สิ่งที่ครอบคลุม |
| --- | --- |
| `mcp` | เซิร์ฟเวอร์ MCP ผ่าน `mcp.json` ของกลุ่ม ([เซิร์ฟเวอร์ MCP ที่ใช้ร่วมกัน](/features/mcp)) |
| `settings` | `settings.json` |
| `memory` | `CLAUDE.md` ระดับ global |
| `extensions` | `agents/`, `commands/`, `skills/`, `plugins/`, `output-styles/`, `rules/` |
| `history` | `projects/` และ `history.jsonl` เพื่อให้อีกบัญชีหนึ่งใช้ `--continue` ได้ |

รายการที่ไม่มีกลุ่มใดของบัญชีใช้ร่วมกันจะอยู่ในโฟลเดอร์ของบัญชีเอง

## คำสั่งและตัวเลือก

| คำสั่ง | สิ่งที่ทำ |
| --- | --- |
| `ccam group` หรือ `ccam group list` | แสดงรายการกลุ่ม |
| `ccam group show <group>` | สิ่งที่กลุ่มใช้ร่วมกัน บัญชีในกลุ่ม โฟลเดอร์ และเซิร์ฟเวอร์ MCP |
| `ccam group add <group> --items <items> [--members <accounts>]` | สร้างกลุ่ม รายการของบัญชีในกลุ่มจะถูกย้ายเข้าไป |
| `ccam group set <group> --add <items> \| --remove <items> \| --items <items>` | เปลี่ยนสิ่งที่กลุ่มใช้ร่วมกัน สำหรับทุกบัญชีในกลุ่ม |
| `ccam group join <group> <account>...` | นำบัญชีเข้ากลุ่ม |
| `ccam group leave <group> <account>...` | นำบัญชีออกจากกลุ่ม โดยแต่ละบัญชีจะเก็บสำเนาไว้ |
| `ccam group rename <group> <new-name>` | เปลี่ยนชื่อเฉพาะในการตั้งค่า |
| `ccam group remove <group> [--leave]` | ลบกลุ่ม |

| ตัวเลือก | สิ่งที่ทำ | ค่าเริ่มต้น |
| --- | --- | --- |
| `--from <account>` | `add`, `set`: บัญชีที่สำเนาของตนจะเข้าไปในกลุ่มก่อน | บัญชีเริ่มต้น แล้วตามหมายเลข |
| `--replace` | `join`: ออกจากกลุ่มที่ใช้รายการเดียวกันร่วมกันอยู่ก่อน | ปิด |
| `--fresh` | บัญชีที่เลิกใช้รายการหนึ่งร่วมกันจะเริ่มโดยไม่มีรายการนั้น แทนที่จะได้สำเนา | ปิด |
| `--leave` | `remove`: นำบัญชีของกลุ่มออกก่อน | จำเป็นเมื่อกลุ่มมีบัญชีอยู่ |
| `--prefer <group\|mine\|merge>` | วิธีตัดสินความขัดแย้ง | ถาม; เมื่อใช้ `--yes` จะใช้ค่าเริ่มต้นของแต่ละไฟล์ |
| `-y`, `--yes` | ทำงานโดยไม่ถาม ความขัดแย้งจะเป็นไปตาม `--prefer` | ปิด |
| `--dry-run` (global) | พิมพ์แผน โดยไม่เปลี่ยนแปลงสิ่งใด | ปิด |

ระบุบัญชีด้วยชื่อหรือหมายเลข และระบุรายการเป็นรายการคั่นด้วยจุลภาค `ccam add` รับ `--group <group>` (ใส่ซ้ำได้) หรือ `--no-group`

## วิธีการทำงาน

ไฟล์ที่เกี่ยวข้อง:

- `~/.config/ccam/config.toml`: รายการ `[[group]]` ได้แก่ `name`, `items`, `members` และ `dir` เมื่อชื่อโฟลเดอร์ต่างจากชื่อกลุ่ม
- `~/.local/share/ccam/profiles/<name>/`: `CLAUDE_CONFIG_DIR` ของบัญชี
- `~/.local/share/ccam/groups/<folder>/`: สำเนาที่ใช้ร่วมกัน

เมื่ออ่าน `config.toml` กลุ่มของแต่ละบัญชีจะถูกแปลงเป็นโฟลเดอร์เดียวต่อรายการ ซึ่งใช้เวลาประมาณ 25 µs สำหรับห้าบัญชี การเริ่มบัญชีจึงยังอ่านเพียง `config.toml` การตั้งค่าที่ทำให้บัญชีหนึ่งได้รับรายการเดียวกันจากสองกลุ่มถือเป็นข้อผิดพลาด และจะไม่มีอะไรเริ่มทำงานจนกว่าจะแก้ไข

การเปลี่ยนแปลงทำกับสำเนาของการตั้งค่า:

1. คำสั่งแก้ไขสำเนานั้น (`AddGroup`, `SetGroupItems`, `JoinGroup` …) ซึ่งตรวจสอบกฎไปด้วย
2. `internal/group` เปรียบเทียบของเก่ากับของใหม่ และสร้างแผนด้วย `share.Build` สำหรับแต่ละบัญชีที่รายการเปลี่ยน แผนถูกสร้างทีละแผน และแต่ละแผนจะเห็นสิ่งที่แผนก่อนหน้าจะใส่ลงในโฟลเดอร์ของกลุ่ม จึงมีเพียงบัญชีแรกเท่านั้นที่ย้ายสำเนาของตนเข้าไป สำเนาแรกมาจากโฟลเดอร์ของกลุ่มหากมีอยู่แล้ว ถัดมาคือ `--from` แล้วบัญชีเริ่มต้น แล้วบัญชีอื่นตามหมายเลข
3. ccam จะปฏิเสธหากบัญชีใดในนั้นกำลังทำงาน และบน Windows ที่ไม่ได้เปิด Developer Mode จะปฏิเสธหากไฟล์ใดต้องใช้ลิงก์
4. ความขัดแย้งจะถูกตัดสินผ่านพรอมต์หรือด้วย `--prefer` จากนั้นแผนจะถูกพิมพ์ แล้ว ccam จะถาม
5. การตั้งค่าใหม่จะถูกบันทึกก่อนที่จะย้ายสิ่งใด หากขั้นตอนใดหยุดกลางคัน `ccam doctor --fix` และการรันซ้ำจะสร้างแผนไปสู่กลุ่มที่บันทึกไว้และทำการเปลี่ยนแปลงให้เสร็จ
6. แผนจะรันตามลำดับ ทุกขั้นตอนจะย้ายข้อมูลเข้าไปในกลุ่มหรือโฟลเดอร์สำรอง ไม่มีสิ่งใดถูกลบทิ้ง

`share.Layout` คือรายการที่อนุญาตของพาธที่แต่ละรายการครอบคลุม สิ่งที่ไม่อยู่ในรายการนี้ไม่มีทางไปถึงกลุ่มได้ `share.Never` ระบุ `.credentials.json` และ `.claude.json` และ `Build` จะหยุดพร้อมข้อผิดพลาดภายในหากไฟล์ใดไฟล์หนึ่งปรากฏในโครงสร้าง รายการ `mcp` ไม่มีพาธ: รายการเซิร์ฟเวอร์อยู่ภายใน `.claude.json` ดังนั้น[การใช้ MCP ร่วมกัน](/features/mcp)จึงทำงานผ่าน `mcp.json` ของกลุ่มแทน

นี่คือสิ่งที่เกิดขึ้นกับโฟลเดอร์หนึ่งโฟลเดอร์ `projects/` เมื่อ `side` เริ่มใช้ประวัติร่วมกันและกลุ่มมีเนื้อหาอยู่แล้ว:

<!-- motion:share-on -->

หากกลุ่มยังไม่มีสิ่งใดในพาธนั้น บัญชีนี้จะเป็นบัญชีแรกที่ใช้ร่วมกัน และโฟลเดอร์ของบัญชีจะถูกย้ายไปที่นั่นตามเดิม ในกรณีอื่น แต่ละไฟล์ที่กลุ่มยังไม่มีจะถูกย้ายเข้าไปในกลุ่ม ส่วนไฟล์ที่มีอยู่ทั้งสองฝั่งจะยังคงอยู่ในโฟลเดอร์ของบัญชี เวอร์ชันของกลุ่มจึงถูกเก็บไว้ จากนั้นโฟลเดอร์ของบัญชีจะไปอยู่ที่ `.ccam-backup/<time>/` ภายในโฟลเดอร์ของบัญชี และลิงก์จะถูกวางไว้แทนที่ ไฟล์การสนทนาตั้งชื่อตาม session ID ประวัติจากหลายบัญชีจึงรวมกันได้โดยไม่ชนกัน รีจิสทรีของ `plugins/` เก็บพาธแบบสัมบูรณ์ที่ผ่านโฟลเดอร์ของบัญชี ccam จะเปลี่ยนให้ชี้ไปที่โฟลเดอร์ของกลุ่ม และรวมรายการของบัญชีเข้าไปในรีจิสทรีที่กลุ่มมีอยู่แล้ว

การออกจากกลุ่มจะย้อนลิงก์กลับ:

<!-- motion:share-off -->

ccam จะลบลิงก์เฉพาะเมื่อลิงก์ชี้เข้าไปในโฟลเดอร์ groups จากนั้นคัดลอกเนื้อหาปัจจุบันของกลุ่มเข้าไปในบัญชี เมื่อใช้ `--fresh` จะข้ามการคัดลอก กลุ่มยังคงเก็บเนื้อหาไว้สำหรับบัญชีอื่น บัญชีที่ย้ายรายการหนึ่งไปยังกลุ่มที่โฟลเดอร์ยังว่างอยู่จะนำสำเนาของสิ่งที่ตนเห็นไปด้วย จึงยังเก็บไฟล์ของตนไว้ได้ เซิร์ฟเวอร์ MCP ก็ตามมาในลักษณะเดียวกัน บัญชีที่เลิกใช้ MCP ร่วมกันจะเสียเซิร์ฟเวอร์ที่ ccam คัดลอกไว้ใน `.claude.json` ของบัญชี (โดยมีการสำรองไว้) เพราะเซิร์ฟเวอร์เหล่านั้นมาจากกลุ่ม

ไฟล์เดียวที่มีอยู่ทั้งสองฝั่งแต่เนื้อหาต่างกันถือเป็นความขัดแย้ง:

<!-- motion:share-conflict -->

`history.jsonl` มีตัวเลือก `merge by line`, `keep group` และ `use mine` โดยการรวมมาเป็นอันดับแรก ไฟล์ `.json` อื่น ๆ มีตัวเลือก `keep group`, `use mine` และ `merge JSON` ส่วนไฟล์อื่นมีตัวเลือก `keep group` และ `use mine` `merge JSON` รวมออบเจ็กต์แบบเรียกซ้ำ ต่ออาร์เรย์โดยไม่ให้มีค่าซ้ำ และเก็บค่าของกลุ่มไว้เมื่อทั้งสองฝั่งตั้งคีย์เดียวกันต่างกัน บัญชีอื่นจึงไม่เห็นการเปลี่ยนแปลงใด ทุกตัวเลือกยกเว้น `keep group` จะบันทึกสำเนาของกลุ่มไว้ที่ `.ccam-backup/<time>/group/` ของบัญชีก่อน ไฟล์ของคุณเองจะไปอยู่ในที่สำรองเสมอ

ทุกครั้งที่เริ่ม `share.Check` จะเปรียบเทียบแต่ละพาธในโครงสร้างกับกลุ่มของบัญชีโดยใช้เพียง `lstat` และ `readlink` ลิงก์ที่หายไปจะถูกสร้างทันที ไฟล์จริงในตำแหน่งที่ควรเป็นลิงก์ ลิงก์ที่ชี้ผิดที่ หรือลิงก์ที่ค้างอยู่จะหยุดการเริ่มและแนะนำ `ccam doctor --fix`

### การเปลี่ยนชื่อ การลบ และรูปแบบเดิม

- `ccam group rename` เปลี่ยนเฉพาะการตั้งค่า โฟลเดอร์ยังใช้ชื่อเดิม และกลุ่มจะบันทึกชื่อนั้นไว้เป็น `dir` จึงไม่มีอะไรถูกย้ายและไม่มีอะไรสูญหาย
- `ccam group remove` จะปฏิเสธกลุ่มที่มีบัญชีอยู่ เว้นแต่คุณจะเพิ่ม `--leave` จากนั้นแต่ละบัญชีจะเก็บสำเนาไว้ และโฟลเดอร์จะไปอยู่ที่ `~/.local/share/ccam/trash/` เว้นแต่กลุ่มอื่นใช้โฟลเดอร์นั้นอยู่
- ก่อนมีกลุ่ม `config.toml` เก็บตาราง `[profile.share]` ไว้ต่อบัญชี ccam อ่านรูปแบบนั้นเป็นกลุ่มโดยไม่ย้ายสิ่งใด เมื่อบันทึกครั้งแรกจะเขียนเป็นรูปแบบใหม่และเก็บไฟล์เดิมไว้เป็น `config.toml.bak` รายการที่บัญชีชุดเดียวกันใช้ร่วมกันจะกลายเป็นกลุ่มเดียวที่ตั้งชื่อตามโฟลเดอร์ ส่วนรายการที่เหลือจะได้ชื่ออย่าง `default-history` และยังใช้โฟลเดอร์เดิมต่อไป

### บน Windows

Windows สร้างลิงก์สัญลักษณ์ได้เฉพาะใน Developer Mode (Settings > System > For developers) หรือสำหรับผู้ดูแลระบบ หากไม่มีสิทธิ์นั้น โฟลเดอร์ที่ใช้ร่วมกันจะกลายเป็น directory junction ซึ่งทำงานกับ Claude Code ได้เหมือนกัน ส่วนไฟล์ที่ใช้ร่วมกัน (`settings.json`, `CLAUDE.md`, `history.jsonl`) เป็นได้เพียงลิงก์สัญลักษณ์เท่านั้น การเปลี่ยนแปลงกลุ่มที่ต้องใช้ลิงก์แบบนี้จึงถูกปฏิเสธก่อนที่จะย้ายสิ่งใด

## ความปลอดภัยและกรณีพิเศษ

- **บัญชีที่กำลังทำงาน** การเปลี่ยนแปลงจะถูกปฏิเสธในขณะที่ Claude Code กำลังทำงานด้วยบัญชีใดบัญชีหนึ่งที่การเปลี่ยนแปลงนั้นจะย้าย เพราะเซสชันยังคงเขียนลงไฟล์เหล่านั้น การออกจากกลุ่มและการเข้าร่วมกลุ่มเกี่ยวข้องเฉพาะบัญชีที่ระบุเท่านั้น `--dry-run` ยังคงพิมพ์แผนได้
- **การสำรอง** แต่ละแผนจะได้โฟลเดอร์ `.ccam-backup/<time>/` ของตัวเองในโฟลเดอร์ของบัญชี ccam ไม่เคยลบโฟลเดอร์เหล่านี้
- **ลิงก์ของคุณเอง** หากพาธเป็นลิงก์ที่คุณสร้างไว้อยู่แล้ว เช่น ลิงก์ไปยังรีโพสิทอรี dotfiles และกลุ่มยังว่าง ccam จะคัดลอกปลายทางของลิงก์เข้าไปในกลุ่ม และย้ายตัวลิงก์ไปไว้ในที่สำรอง ปลายทางจะไม่ถูกแตะ
- **โฟลเดอร์ที่นำเข้า** บัญชีที่ลงทะเบียนไว้ในตำแหน่งเดิมด้วย `ccam import` (เช่น `~/.claude`) จะเปลี่ยนแปลงในตำแหน่งเดิม `claude` แบบปกติจึงเห็นการเปลี่ยนแปลงด้วย และแผนจะแจ้งไว้
- **สองบัญชีในโฟลเดอร์เดียว** บัญชีที่ใช้ประวัติร่วมกันจะเขียนไฟล์เดียวกัน เมื่อบัญชีที่สองเริ่มในโฟลเดอร์ที่มีอีกบัญชีหนึ่งทำงานอยู่แล้ว ccam จะเตือนและถาม
- **การแก้ไขด้วยมือ** สมาชิกที่ไม่ใช่บัญชีจะถูกละเว้นและ `ccam doctor` จะรายงานให้ทราบ ชื่อ `local`, `on`, `off` และ `none` ใช้เป็นชื่อกลุ่มไม่ได้

## หัวข้อที่เกี่ยวข้อง

- [แดชบอร์ด](/features/dashboard)
- [ใช้เซิร์ฟเวอร์ MCP การตั้งค่า และประวัติร่วมกัน](/docs/sharing)
- [เซิร์ฟเวอร์ MCP ที่ใช้ร่วมกัน](/features/mcp)
- [คัดลอกการตั้งค่าที่มีอยู่](/features/migrate)
- [คำสั่งทั้งหมด](/docs/commands)
