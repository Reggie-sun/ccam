---
title: Kontrol paneli
description: Terminalde ccam; hesapların, grupların, klasör kurallarının, uzantıların ve denetimlerin tam ekran kontrol panelini açar; her değişiklik önce planını gösterir.
updated: 2026-10-20
order: 15
translation: machine
sourceHash: "decdd5fa54c9"
---

Argümansız `ccam`, standart girdi ve çıktı bir terminal olduğunda tam ekran bir kontrol paneli açar; aksi halde `ccam list` çıktısını yazdırır. Kontrol paneli her hesabı kullanımı ve gruplarıyla birlikte gösterir ve komutların yaptığını yapar: başlatma, ekleme, giriş yapma, grupları değiştirme, klasör kuralları belirleme, denetleme ve onarma. Başlatma paneli kapatır; bir hesap eklemek veya giriş yapmak paneli terminale bırakır ve sonra geri döner.

```console
 ccam v0.9.0 · 3 accounts · 2 groups         ✔ default · ▸ this folder · • running

  Accounts   Groups   Folders   Extensions   Health

    # Account     Email             5-hour         7-day Last  Groups
 ❯ 1. personal ✔   me@example.com    █████▊░░  72%   41% 2h    default
   2. work      ▸• dev@company.com   ██░░░░░░  25%  100% now   default team
   3. side         -                 ✘ not logged in         -     no group

 ✔ personal is the default account
 enter to launch · n for most room · a to add · s for groups · ? for shortcuts
```

İlk kare yalnızca `config.toml` dosyasından, 100 ms'nin çok altında çizilir; kullanım her hesap yanıt verdikçe dolar. Bir yüzde %80'den itibaren sarıya, %100'de kırmızıya döner; başarısız bir sorgudan sonra önbellekten gelen değerler soluk çizilir.

## Sekmeler

| Sekme | Gösterdiği | Tuşlar |
| --- | --- | --- |
| Accounts | numara, ad, e-posta, 5 saatlik çubuk, 7 günlük kullanım, son kullanım, gruplar; `✔` varsayılan, `▸` bu klasörün hesabı, `•` çalışıyor | aşağıda |
| Groups | her grup, neyi paylaştığı ve hesapları, ayrıca hiçbir grupta olmayan hesaplar | `enter` düzenle, `n` yeni, `e` yeniden adlandır, `x` kaldır |
| Folders | [`ccam folder`](/features/folder) kuralları ve burada geçerli olan kural | `n` burada bir hesap kullan, `enter` değiştir, `x` kaldır |
| Extensions | [`ccam ext`](/features/ext) komutunun gösterdikleri, her uzantı bir kez ve geldiği gruplarla | `r` yeniden yükle |
| Health | [`ccam doctor`](/features/doctor) komutunun bildirdikleri | `f` onar (`doctor --fix`) |

`←`/`→` veya `Tab` sekmeler arasında geçer, `?` geçerli sekmenin tuşlarını listeler, `r` kısa önbelleği kullanmadan kullanımı yeniden sorgular, `q` veya `Esc` ise çıkar.

## Accounts sekmesindeki tuşlar

| Tuş | Ne yapar |
| --- | --- |
| `Enter` | Claude Code'u seçili hesapla başlatır |
| `1`–`9` | O numaralı hesaba gider; 12 için `1` `2` tuşlarına hızlıca basın. Yalnızca `Enter` başlatır |
| `↑`/`↓`, `k`/`j` | Hareket eder |
| `n` | [`ccam next --auto`](/features/next) komutunun seçeceği hesabı başlatır |
| `o` | Seçeneklerle başlatır: `--continue`, `--resume`, başka argümanlar ya da bir kuru çalıştırma |
| `u` | Kullanım ayrıntıları: sıfırlanma zamanlarıyla 5 saatlik, 7 günlük ve model pencereleri ve sayıların nereden geldiği |
| `a` | Bir hesap ekler (aşağıdaki sihirbaz) |
| `l` | Yeniden giriş yapar |
| `s` | Hesabın gruplarını seçer |
| `d` | Hesabı varsayılan hesap yapar |
| `f` | Hesabı geçerli klasörde kullanır |
| `e` | Hesabı yeniden adlandırır |
| `c` | Hesabın kısayol komutunu ayarlar veya kaldırır |
| `x` | Hesabı kaldırır (önce, Hayır seçili olarak sorar) |

## Grupları değiştirme

Groups sekmesinde `Enter` düzenleyiciyi açar: solda öğeler, sağda hesaplar; `Tab` aralarında geçer, `Space` açıp kapatır. Hesaplardan birinin zaten başka bir gruptan aldığı bir öğe, nedeniyle birlikte gri gösterilir; bir öğeyi iki kez alacak bir hesap da öyle. Bir öğe yalnızca tek bir gruptan gelebilir. `Enter` planı gösterir:

```console
 ──────────────────────────────────────────────────────────────────────────────
  Groups of side

  ● side  +history (group team)
    ⎿  link     projects, history.jsonl: profiles/side → groups/team

  Do you want to apply this plan?
  ❯ 1. Yes
    2. No
```

Grubun kopyasından farklı olan bir dosya için önce sorulur (grubunkini koru, benimkini kullan ya da birleştir). Ardından değişiklik, [`ccam group`](/features/share) ile aynı kodla arka planda çalışır: önce yapılandırma kaydedilir, her hesabın verisi sırayla taşınır ve orijinaller `.ccam-backup/` klasörüne gider. Çalışan bir hesabın verisini taşıyacak bir değişiklik, hesabın adıyla birlikte reddedilir.

## Hesap ekleme

`a`, `ccam add` komutunun da terminalde kullandığı sihirbazı açar:

1. Hesap adı.
2. Nasıl giriş yapacağı: şimdi bir Claude aboneliği (tarayıcıda), bir [API anahtarı](/features/api) ya da daha sonra bir Claude aboneliği.
3. API anahtarı için: ccam'in anahtarı her başlatmada nereden okuyacağı (`env:NAME` veya `command:COMMAND`), uç nokta, anahtarın nasıl gönderileceği ve model.
4. `default` işaretli olarak katılınacak gruplar; henüz grup yoksa `default` grubunun neyi paylaşacağı.
5. İsterseniz bir kısayol komutu.
6. Onaylamanız için bir özet.

`Esc` bir adım geri gider; ilk adımda ise sihirbazı kapatır.

## Renkler

Renkler Claude Code'un terminal arayüzündeki rolleri izler; böylece ccam onun yanında yabancı durmaz. Grup adları, adlarından türetilen bir renk alır. ccam terminale arka planını sorar ve açık ya da koyu bir tema seçer. Temayı `config.toml` içinde kendiniz ayarlayın:

```toml
theme = 'dark'   # auto (default), dark, light, dark-daltonized, light-daltonized, dark-ansi, light-ansi
```

`-daltonized` temaları kırmızı-yeşil renk körlüğü içindir; `-ansi` temaları ise terminalin kendi 16 rengini kullanır. `NO_COLOR` renkleri kapatır.

## Güvenlik ve uç durumlar

- Kontrol paneli komutlarla aynı işlevleri çağırır; bu yüzden her kural geçerlidir: veri taşınmadan önce planlar, çalışan hesaplar için reddetme, yedekler ve bir hesaptan hiçbir zaman çıkmayan giriş token'ları.
- Kontrol panelinin kuru çalıştırması yoktur: ad verilmeden `ccam --dry-run`, bu klasörün hesabının nasıl başlayacağını yazdırır; bir hesapta `o` ise yalnızca o tek başlatma için bir kuru çalıştırma sunar.
- Terminal, Claude Code başlamadan veya giriş açılmadan önce eski haline getirilir.
- Bir iletişim kutusu için fazla küçük bir ekran önce sekmenin satırlarını gizler; satırlar kaydırılmak yerine `…` ile kesilir.

## İlgili

- [Gruplarla paylaşım](/features/share)
- [Claude Code'u başlatma](/features/launch)
- [Hesapları yönetme](/features/accounts)
- [Kullanım limitleri](/features/usage)
