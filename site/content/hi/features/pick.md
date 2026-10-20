---
title: scripts के लिए अकाउंट चुनना
description: ccam pick script के लिए सबसे ज़्यादा गुंजाइश वाला अकाउंट प्रिंट करता है, साथ शुरू हुए workers को अलग अकाउंट देता है, और गुंजाइश न होने पर 75 देता है।
updated: 2026-10-20
order: 12
translation: machine
sourceHash: 65a5a7e179a5
---

`ccam pick` सबसे ज़्यादा गुंजाइश वाले लॉगिन किए हुए अकाउंट का नाम प्रिंट करता है, ताकि कोई script उससे Claude Code शुरू कर सके। जब तक गुंजाइश वाले अकाउंट काफ़ी हों, एक ही समय पर शुरू हुए हर worker को अलग अकाउंट मिलता है, इसलिए `claude -p` jobs का एक batch किसी एक अकाउंट को ख़ाली करने के बजाय आपके अकाउंट्स में बँट जाता है।

## उपयोग

```sh
ccam "$(ccam pick)" -p "fix the failing tests"
for task in parse lint docs; do
  ccam "$(ccam pick)" -p "do the $task task" > "$task.json" &
done; wait
ccam pick --json                      # the account and its usage
ccam pick --model opus                # also count the Opus weekly limit
ccam pick --max-usage 80 --strategy expiring
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
  "source": "live"
}
```

Exit codes: standard output पर नाम के साथ `0`, जब किसी अकाउंट में गुंजाइश न हो या कोई लॉगिन न हो तो `75` (`EX_TEMPFAIL`), और ग़लत flag जैसी दूसरी त्रुटियों के लिए `1`।

## विकल्प

| विकल्प | क्या करता है | डिफ़ॉल्ट |
| --- | --- | --- |
| `--json` | अकाउंट, उसका 5 घंटे और साप्ताहिक उपयोग, और संख्याएँ कहाँ से आईं, यह प्रिंट करें | बंद |
| `--model <model>` | उस मॉडल की साप्ताहिक विंडो भी गिनें, जैसे `opus` | कोई नहीं |
| `--max-usage <1-100>` | इस उपयोग से नीचे वाले अकाउंट्स को प्राथमिकता दें | `auto_max_usage`, वरना 100 |
| `--strategy headroom\|expiring` | सबसे कम 5 घंटे का उपयोग, या सबसे जल्दी होने वाला 5 घंटे का reset | `auto_strategy`, वरना `headroom` |
| `--wait[=<duration>]` | जब किसी अकाउंट में गुंजाइश न हो, तो पहले अकाउंट के reset होने का इंतज़ार करें, ज़्यादा से ज़्यादा इतनी देर | बंद; अकेले देने पर `5h` |
| `--dry-run` | चुनाव दर्ज किए बिना चुनें | बंद |

## यह कैसे काम करता है

1. ccam लॉगिन किए हुए अकाउंट्स की सूची बनाता है (API key अकाउंट्स की कोई सीमा नहीं होती, इसलिए वे छोड़ दिए जाते हैं)।
2. यह एक lock (`~/.local/state/ccam/picks.lock`) लेता है, ताकि एक साथ चल रहे picks एक के बाद एक हों।
3. यह हर अकाउंट का उपयोग एक साथ पूछता है, उसी तरह जैसे `ccam usage` करता है, 60 सेकंड के cache के साथ।
4. पिछले 2 मिनट में चुने गए अकाउंट (`~/.local/state/ccam/picks/<name>`) पहले छोड़ दिए जाते हैं, और बाक़ी में से `ccam next --auto` के [चुनाव के नियम](/features/next#auto) चुनते हैं: पहले `auto_max_usage` से नीचे वाले, strategy के हिसाब से, और 100% वाला अकाउंट कभी नहीं।
5. अगर गुंजाइश वाला हर अकाउंट हाल ही में चुना गया हो, तो नियम फिर से उन सब में से चुनते हैं, इसलिए तीन अकाउंट्स पर चौथा worker किसी एक को साझा करता है।
6. चुनाव दर्ज होता है और lock छोड़ा जाता है, फिर नाम प्रिंट होता है।

2 मिनट की विंडो इसलिए ज़रूरी है क्योंकि नए job का उपयोग संख्याओं में दिखने में कुछ समय लगता है: इसके बिना, एक साथ शुरू हुए पाँच workers सभी एक ही अकाउंट को सबसे ज़्यादा गुंजाइश वाला मानते।

## सुरक्षा और विशेष स्थितियाँ

- `ccam pick` केवल उपयोग पढ़ता है और चुनाव का समय लिखता है; यह `ccam usage` जितना करता है उससे ज़्यादा tokens कभी renew नहीं करता और न पढ़ता है।
- चुनाव अकाउंट को आरक्षित नहीं करता। 2 मिनट बाद दो scripts फिर भी एक ही अकाउंट पर शुरू हो सकती हैं; जब वह अकाउंट अपनी सीमा तक पहुँचे तो सँभलने के लिए `--fallback` इस्तेमाल करें।
- `--wait` से इंतज़ार करते समय कोई lock नहीं रखा जाता, इसलिए इस बीच दूसरे picks चलते रहते हैं।
- `ccam pick` [`ccam folder`](/features/folder) को नज़रअंदाज़ करता है: यह केवल उपयोग के हिसाब से चुनता है।

## संबंधित

- [अगले अकाउंट पर स्विच करना](/features/next)
- [सीमा पर फिर से चलाना](/features/next#run-again-at-the-limit)
- [उपयोग सीमाएँ](/features/usage)
- [API key अकाउंट](/features/api)
