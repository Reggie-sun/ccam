---
title: scripts साठी खाते निवडणे
description: ccam pick script साठी सर्वाधिक जागा असलेले खाते छापते, एकत्र सुरू झालेल्या workers ना वेगवेगळी खाती देते, आणि जागा नसेल तर 75 ने बाहेर पडते.
updated: 2026-10-25
order: 12
translation: machine
sourceHash: "5c4f74db9714"
---

`ccam pick` सर्वाधिक जागा असलेल्या लॉग इन केलेल्या खात्याचे नाव छापते, जेणेकरून एखादी script त्याने Claude Code सुरू करू शकेल. जागा असलेली खाती पुरेशी असेपर्यंत एकाच वेळी सुरू झालेल्या प्रत्येक worker ला वेगळे खाते मिळते, त्यामुळे `claude -p` jobs चा एक संच एकच खाते संपवण्याऐवजी तुमच्या खात्यांमध्ये विभागला जातो.

**Team आणि Enterprise सीट्स.** आपोआप होणारी निवड (`ccam next`, `--fallback`, `ccam pick`, डॅशबोर्डमधील `n`, आणि तुम्ही चालू केल्यावर [`auto_switch`](/features/next#auto-switch)) एकाच Claude Team किंवा Enterprise organization मधील एका सीटवरून दुसऱ्या सीटवरही जाते. Anthropic च्या मदत केंद्रात म्हटले आहे की "Each seat on a Team plan is assigned to one named individual, and only that person may use it" (Team प्लॅनमधील प्रत्येक सीट नावानिशी एका व्यक्तीला नेमून दिलेली असते, आणि फक्त तीच व्यक्ती ती वापरू शकते) आणि "Users can't route their work through another user's account to use that user's remaining usage" (दुसऱ्या वापरकर्त्याच्या उरलेल्या वापराचा लाभ घेण्यासाठी वापरकर्ते आपले काम त्याच्या खात्यामार्फत चालवू शकत नाहीत) ([Team प्लॅन](https://support.claude.com/en/articles/12004354-purchase-and-manage-seats-on-team-plans#h_b8a7d07129); [Enterprise प्लॅन](https://support.claude.com/en/articles/13393991-purchase-and-manage-seats-on-enterprise-plans#h_fa6b9cade0) साठीही असेच). त्यामुळे सीट्सदरम्यान आपोआप होणारी निवड वापरल्यास Anthropic च्या अटींचे उल्लंघन होऊ शकते. तिथे ती वापरायची की नाही हा तुमचा निर्णय आहे: त्यावर अवलंबून राहण्यापूर्वी तुमच्या संस्थेकडे आणि Anthropic कडे चौकशी करा. वैयक्तिक प्लॅनबद्दल (Pro, Max) [FAQ मधील अटींविषयीची टीप](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts) सांगते.

## वापर

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
ccam pick --strategy balance          # spread over the weekly quota, fewer running sessions first
ccam pick --cross-org                 # also accounts of other organizations
ccam pick --wait=2h                   # nobody has room: wait for a reset
ccam --fallback "$(ccam pick)" -p "…" # also switch if this account hits its limit mid-job
```

```console
$ ccam pick --json
{
  "name": "work",
  "slot": 2,
  "email": "dev@company.com",
  "five_hour": { "percent": 12, "resets_at": "2026-10-20T15:10:00Z", "known": true },
  "weekly": { "percent": 41, "resets_at": "2026-10-24T09:00:00Z", "known": true },
  "source": "live",
  "reason": "12% of its 5 hours used, the least"
}
```

Exit codes: standard output वर नावासह `0`, कोणत्याही खात्यात जागा नसेल किंवा कोणतेही खाते लॉग इन नसेल तर `75` (`EX_TEMPFAIL`), आणि चुकीच्या flag सारख्या इतर त्रुटींसाठी `1`.

## पर्याय

| पर्याय | काय करते | डीफॉल्ट |
| --- | --- | --- |
| `--json` | खाते, त्याचा 5 तासांचा आणि साप्ताहिक वापर, आकडे कुठून आले आणि ते का निवडले ते छापा | बंद |
| `--model <model>` | त्या model ची साप्ताहिक window सुद्धा मोजा, जसे `opus` | काहीही नाही |
| `--max-usage <1-100>` | या वापराखालील खात्यांना प्राधान्य द्या | `auto_max_usage`, नाहीतर 100 |
| `--strategy headroom\|expiring\|balance` | सर्वात कमी 5 तासांचा वापर, सर्वात लवकर होणारा 5 तासांचा reset, किंवा सर्वात आधी न वापरता वाया जाणारा साप्ताहिक quota | `auto_strategy`, नाहीतर `headroom` |
| `--cross-org` | इतर organizations ची खातीसुद्धा निवडा | `cross_org`, नाहीतर बंद |
| `--wait[=<duration>]` | कोणत्याही खात्यात जागा नसेल तर पहिले खाते reset होण्याची वाट पाहा, जास्तीत जास्त इतका वेळ. मूल्य `=` नंतर येते: `--wait 30m` हे `--wait takes its value after '=': --wait=30m` या सूचनेसह नाकारले जाते | बंद; एकटे दिल्यास `5h` |
| `--dry-run` | निवड नोंदवल्याशिवाय निवडा | बंद |

## हे कसे काम करते

1. ccam हा फोल्डर जे खाते सुरू करेल त्या खात्याच्या organization मधील लॉग इन केलेल्या खात्यांची यादी करते (API key खाती आणि Console लॉगिन यांना मर्यादा नसते, त्यामुळे ती वगळली जातात; [organizations](/features/next#auto) पाहा); त्याच Claude खात्यात लॉग इन असलेला फोल्डरसुद्धा वगळला जातो.
2. ते एक lock (`~/.local/state/ccam/picks.lock`) घेते, जेणेकरून एकाच वेळी चालणाऱ्या निवडी एकामागून एक होतात.
3. ते प्रत्येक खात्याचा वापर एकाच वेळी विचारते, `ccam usage` प्रमाणेच, 60 सेकंदांच्या cache सह. Claude Code ने भरलेली म्हणून नोंदवलेली window तिच्या reset पर्यंत 100% मोजली जाते ([मर्यादेची नोंद](/features/usage#the-limit-record)), त्यामुळे नुकतीच मर्यादा गाठलेले खाते, त्याचे साठवलेले आकडे जागा दाखवत असले तरी, बाहेरच राहते.
4. मागील 2 मिनिटांत निवडलेली खाती, `ccam pick` ने किंवा [`auto_switch`](/features/next#auto-switch) ने तिथे सत्र हलवून (`~/.local/state/ccam/picks/<name>`; घड्याळाच्या वेळेनंतरची तारीख असलेली निवड मोजली जात नाही), आधी वगळली जातात, आणि उरलेल्यांतून `ccam next --auto` चे [निवडीचे नियम](/features/next#auto) निवडतात: आधी `auto_max_usage` खालची, strategy नुसार, आणि 100% वरचे खाते कधीच नाही. यामुळे नियम जे खाते सर्वात आधी घेतील ते वगळले गेले, तर `--json` चे `reason` तसे सांगते: `30% of its 5 hours used; work was picked a moment ago`.
5. जागा असलेले प्रत्येक खाते नुकतेच निवडले गेले असेल तर नियम पुन्हा त्या सर्वांतून निवडतात, त्यामुळे तीन खात्यांवर चौथा worker एखादे खाते शेअर करतो.
6. निवड नोंदवली जाते आणि lock सोडला जातो, मग नाव छापले जाते.

2 मिनिटांची window महत्त्वाची आहे कारण नवीन job चा वापर आकड्यांमध्ये दिसायला थोडा वेळ लागतो: तिच्याशिवाय एकत्र सुरू झालेले पाच workers सर्वजण एकाच खात्याला सर्वाधिक जागा असलेले खाते म्हणून पाहतील.

## सुरक्षितता आणि अपवादात्मक स्थिती

- `ccam pick` फक्त वापर वाचते आणि निवडीची वेळ लिहिते; `ccam usage` जे करते त्यापलीकडे ते कधीच tokens नूतनीकरण करत नाही किंवा वाचत नाही.
- निवड खाते राखून ठेवत नाही. 2 मिनिटांनंतर दोन scripts तरीही एकाच खात्यावर सुरू होऊ शकतात; ते खाते मर्यादा गाठेल तेव्हा सावरण्यासाठी `--fallback` वापरा.
- `--wait` ने वाट पाहताना कोणताही lock धरला जात नाही, त्यामुळे दरम्यान इतर निवडी चालू राहतात.
- `ccam pick` [`ccam folder`](/features/folder) कडे दुर्लक्ष करते: ते फक्त वापरावरून निवडते.

## संबंधित

- [पुढच्या खात्यावर जाणे](/features/next)
- [मर्यादेवर पुन्हा चालवणे](/features/next#run-again-at-the-limit)
- [वापर मर्यादा](/features/usage)
- [API key खाती](/features/api)
