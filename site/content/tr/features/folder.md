---
title: Klasöre göre hesaplar
description: ccam folder, Claude Code'u hesap adı vermeden başlattığınızda bir klasörün ve altındakilerin tek bir hesap kullanmasını sağlar; örneğin şirket kodu için work.
updated: 2026-10-26
order: 11
translation: machine
sourceHash: 25f82da64580
---

`ccam folder work ~/code/company` komutundan sonra `~/code/company` içinde veya alt klasörlerinde hesap adı verilmeden yapılan her başlatma Claude Code'u `work` olarak başlatır. Bu klasörün dışında varsayılan hesap eskisi gibi başlar. Bir hesabın adını vermek, `ccam personal` ile ya da `ccpersonal` gibi bir kısayolla, her zaman o hesabı başlatır.

## Kullanım

```sh
ccam folder work ~/code/company     # use work in ~/code/company and below
ccam folder personal                # use personal in the current folder
ccam folder                         # list the folders and the account used here
ccam folder --remove ~/code/company # stop using an account there
cd ~/code/company/api && ccam -p "summarise the open issues"   # runs as work
```

```console
$ ccam folder work ~/code/company
✓ /home/me/code/company → work
  ccam without an account name starts work here and in every folder below.
$ cd ~/code/company/api && ccam folder
/home/me/code/company → work

Here (/home/me/code/company/api): work, set for /home/me/code/company
```

Şu başlatmalar klasörün hesabını kullanır:

- `ccam -p "…"` ve ilk argümanı Claude Code için bir seçenek olan diğer tüm başlatmalar
- adsız `ccam run` ve `ccam --dry-run`
- ilk denemesi için `ccam --fallback -p "…"`
- imleci klasörün hesabında açılan kontrol paneli, yani tek başına `ccam`; klasöre sembolik bir bağlantı üzerinden ulaştığınızda bile

## Seçenekler

| Seçenek | Ne yapar | Varsayılan |
| --- | --- | --- |
| `<account>` | Kullanılacak hesabın adı veya numarası | |
| `[folder]` | Göreli veya mutlak klasör | geçerli klasör |
| `--remove` | Tam olarak o klasörün girdisini kaldırır | kapalı |
| `--dry-run` | Desteklenir: değişikliği kaydetmeden yazdırır | kapalı |

## Nasıl çalışır

Her klasör, her başlatmanın zaten okuduğu dosya olan `~/.config/ccam/config.toml` içinde bir `[[folder]]` girdisidir:

```toml
[[folder]]
path = '/home/me/code/company'
account = 'work'
```

1. `ccam folder`, klasörü sembolik bağlantıları çözülmüş mutlak bir yola dönüştürür ve kaydeder. Bir klasörü yeniden ayarlamak hesabını değiştirir.
2. Adsız bir başlatma `launcher.HereProfile` çağırır. Hiç girdi yoksa hemen varsayılan hesabı alır.
3. Aksi halde çalışma dizinini her `path` ile dize olarak karşılaştırır: klasörün kendisi veya altındaki bir klasör eşleşir, yalnızca aynı harflerle başlayan bir klasör asla eşleşmez (`/code/company`, `/code/companyx` klasörünü kapsamaz). En derin girdi kazanır; böylece `~/code` `personal`, `~/code/company` ise `work` olabilir.
4. Çalışma dizini verildiği haliyle hiçbir şeyle eşleşmezse ancak o zaman ccam sembolik bağlantılarını çözer ve yeniden karşılaştırır.
5. Windows ve macOS'ta karşılaştırma, dosya sistemlerinin genellikle yaptığı gibi büyük ve küçük harfi yok sayar.

`config.toml` dışında hiçbir dosya okunmaz; bu yüzden bir klasörün içinde başlatmak ölçülebilir bir süre eklemez (girdilerle veya girdiler olmadan 4,0 ms, `hyperfine`, 300 çalıştırma).

## Güvenlik ve uç durumlar

- `ccam rename` klasör girdilerini yeni ada taşır, `ccam remove` ise kaldırılan hesabın girdilerini siler.
- Hesabı artık bulunmayan bir girdi, o klasörde adsız bir başlatmayı durdurur: `<folder> is set to account "<name>", which does not exist. Change it with ccam folder`. `ccam doctor` bu tür girdileri ve klasörü artık olmayan girdileri bildirir.
- Hesabını bir üst klasörden devralan bir alt klasörde `ccam folder --remove` çalıştırırsanız ccam üst klasörün adını verir ve hiçbir şey kaldırmaz.
- `ccam next` sıraya veya kullanıma göre seçer ve klasörlere bakmaz.

## İlgili

- [Claude Code'u başlatma](/features/launch)
- [Kısayol komutları](/features/shortcuts)
- [Hesapları yönetme](/features/accounts)
- [Komut başvurusu](/docs/commands)
