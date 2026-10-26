---
title: ccam 與其他 Claude Code 帳號切換工具比較
description: ccam 與 claude-swap、CCS、clausona、clauth、CAAM、cux 等 Claude Code 帳號切換工具的比較，以及哪種做法適合哪種用途。
updated: 2026-10-26
order: 2
translation: machine
sourceHash: "b6143bcf3dc7"
---

切換 Claude Code 帳號的工具採用兩種做法之一。大多數工具在同一個 `~/.claude` 資料夾中替換登入資訊，所以同一時間只有一個帳號處於使用中。另一些工具（包括 ccam）則用 `CLAUDE_CONFIG_DIR` 為每個帳號提供獨立的資料夾，因此每個帳號都會保持登入，也能同時執行多個帳號。ccam 還提供共用的設定和 MCP 伺服器，並能在另一個帳號的資料夾中接續對話。

## 替換同一組登入資訊，還是使用獨立的資料夾？

**替換登入資訊**只保留一套設定、MCP 伺服器和對話記錄，因為只有一個資料夾。切換很快，執行中的工作階段也能套用新的登入資訊。但同一時間只有一個帳號處於使用中，而且會同時影響所有終端機和編輯器。

為了替換，工具會為其他每個帳號保存一份登入資訊的副本，並在切換時把它寫回原位。Claude Code 每次更新登入時都會換掉其中的重新整理權杖，所以每份保存的副本都必須和執行中的工作階段保持同步。切換後帳號被登出，或需要重新 `/login`，是這些專案最常修正的問題之一。

**獨立的資料夾**是 [Claude Code 文件中說明的方法](https://code.claude.com/docs/en/authentication#log-in-with-multiple-accounts)。每個帳號都保持登入，你可以讓工作和個人的工作階段同時執行。每份登入資訊都留在各自的資料夾中，由 Claude Code 照常在那裡更新。代價是每個資料夾一開始都有各自的設定、MCP 伺服器和對話記錄，除非有工具幫你共用。

## 各工具比較

資料取自各專案自己的文件與原始碼（2026-10-20；clausona：v0.12.0-beta）。「未說明」表示專案沒有提到，而不是缺少該功能。

| 工具 | 做法 | 同時使用多個帳號 | 5 小時和每週用量 | 達到限制時 | 平台 |
| --- | --- | --- | --- | --- | --- |
| ccam | 獨立的資料夾 | 是 | 是，`ccam usage` | 在工作階段中顯示通知；`ccam next --auto --continue` 或 `/ccam:switch` 會把對話移到剩餘額度最多的帳號，無論是否共用對話記錄；可選的 `auto_switch`（Linux、macOS）會自行移動它；`ccam --fallback` 會重新執行 `-p` 工作；`ccam pick` 會把工作程序分散到各帳號 | Linux、macOS、Windows |
| [claude-swap](https://github.com/realiti4/claude-swap) | 替換登入資訊；`cswap run` 時使用獨立的資料夾 | 是，使用 `cswap run` | 是 | `cswap auto` 在達到門檻時切換登入資訊；執行中的工作階段會繼續 | Linux、macOS、Windows |
| [CCS](https://github.com/kaitranntt/ccs) | 獨立的資料夾 | 是 | 部分支援 | 帳號池透過本機 Proxy 把請求轉給下一個帳號 | Linux、macOS、Windows |
| [clausona](https://github.com/larcane97/clausona) | 獨立的資料夾 | 是 | 是 | 路由會在每次執行前選擇還有配額的帳號，可以等待重設，並為工作程序分配不同的帳號 | Linux、macOS、Windows |
| [clauth](https://github.com/uwuclxdy/clauth) | 替換登入資訊；`clauth start` 時使用獨立的資料夾 | 是，使用 `clauth start` | 是 | 沿著備援鏈依序切換 | Linux、macOS、Windows |
| [CAAM](https://github.com/Dicklesworthstone/coding_agent_account_manager) | 替換登入資訊，或使用隔離的設定檔 | 是，使用隔離的設定檔 | 是 | `caam run` 在另一個帳號上重試非互動式指令 | Linux、macOS、Windows |
| [cux](https://github.com/inulute/cux) | 替換登入資訊並包裝 `claude` | 否，只有一個使用中的帳號 | 是 | 替換登入資訊並恢復同一個對話 | Linux、macOS、Windows |
| [CCSwitcher](https://github.com/XueshiQiao/CCSwitcher) | 替換登入資訊 | 否 | 是 | 可選擇在達到門檻時切換 | macOS |
| [Claude Switcher](https://github.com/Symbioose/claude-account-switcher) | 替換登入資訊 | 否 | 是 | 可選擇在 100% 時切換 | macOS |

ccswap 是 claude-swap 的分支，也能管理 Codex 的登入。[claude-account](https://github.com/hamzarehmandeveloper/claude-account)、[claude-multiprofile](https://github.com/jmdarre-v/claude-multiprofile) 和 [claude-code-profiles](https://github.com/quinnjr/claude-code-profiles) 也為每個帳號提供獨立的資料夾。其中幾個工具也支援其他程式設計代理，例如 Codex 或 Gemini。

## 各工具的實作方式

| 工具 | 撰寫語言 | Claude Code 執行時 | 是否在不同位置之間複製登入資訊 | 帳號之間可以共用什麼 |
| --- | --- | --- | --- | --- |
| ccam | Go，單一靜態執行檔 | 預設什麼都沒有：ccam 會以 `claude` 取代自己（在無法這樣做的 Windows 上，則等待 Claude Code 結束）；開啟 `relay` 或 `auto_switch` 時，ccam 會在互動式工作階段旁持續執行 | 否 | 群組：每個群組指定它共用的項目和它的帳號；一個帳號可以屬於多個群組 |
| claude-swap | Python 3.12+，以 uv 或 pipx 安裝 | 用 `cswap run` 時什麼都沒有；`cswap auto` 和選單列 App 會持續執行 | 是，在替換登入資訊時 | 用 `cswap run` 時：設定、`CLAUDE.md` 和擴充功能一起共用，對話記錄可選；MCP 伺服器在每次啟動時複製 |
| CCS | Node.js 18+ 或 Bun | CCS 以父行程留著；帳號池會加上本機 Proxy | 帳號資料夾不複製；帳號池保存自己的登入資訊 | 每個帳號選擇共用或本機，對話記錄以情境群組共用；MCP 伺服器會被複製 |
| clausona | TypeScript，Node.js 20+ | 一般的 `claude` 透過 Shell 函式執行；`clausona run` 讓 Node 以父行程留著；每次結束後會有一個背景 Node 行程 | 否，但它會自行更新閒置帳號的權杖，以讀取它們的用量 | 除了各帳號專屬的檔案以外的所有內容，來自一個主要帳號；對話記錄可選；以 `claude mcp add` 新增的伺服器不會共用 |
| clauth | Rust | clauth 以父行程留著 | 是，在替換登入資訊時 | 整個 `~/.claude`，或什麼都不共用 |
| CAAM | Go | 視模式而定；可選用的常駐程式 | 是，在其保管庫模式下 | 由模式決定；MCP 伺服器在每次啟動時複製 |
| cux | Go | 包裝程式一直包在 `claude` 外面 | 是 | 全部：只有一個資料夾 |
| CCSwitcher, Claude Switcher | macOS App（Swift、Python） | 選單列 App | 是 | 全部：只有一個資料夾 |

## ccam 的優勢

### 登入資訊留在 Claude Code 存放的位置

ccam 絕不會把 Claude 的登入資訊從一個地方複製到另一個地方。每個帳號的登入資訊都留在自己的資料夾中，在 macOS 上則留在該資料夾專屬的鑰匙圈項目中，由 Claude Code 負責更新。不存在會失去同步的保存副本。以 `ccam import` 註冊的現有 `~/.claude` 會繼續使用一般 `claude` 所用的鑰匙圈項目，所以不論用哪種方式啟動，都會保持登入。

只有在你執行 `ccam usage --refresh` 時，ccam 才會自己更新過期的權杖。它會先取得鎖定，並在該帳號執行中時拒絕更新（不論 Claude Code 是如何啟動的，背景工作階段及其常駐程式也算在內），在 Claude Code 自己正在更新該登入資訊時也會拒絕，所以絕不會在執行中的工作階段底下輪替權杖。

### 預設情況下，你和 Claude Code 之間沒有任何東西在執行

ccam 讀取一個檔案，檢查它的連結，然後以 `claude` 取代自己（`exec`）。預設情況下，你工作時沒有 ccam 行程、包裝程式、Proxy 或常駐程式：標準輸入、訊號和結束代碼都屬於 Claude Code，它的請求直接送往 Anthropic。Windows 無法取代行程，所以在那裡 ccam 會以同一個主控台啟動 Claude Code，並等待它的結束代碼，除此之外什麼都不做。設定 `relay = true` 或 `auto_switch = true` 時，ccam 會以同樣的方式在互動式工作階段旁持續執行，以便在另一個帳號接續對話；Claude Code 的請求仍然直接送往 Anthropic。

以這種方式啟動帳號大約只多花 4 ms，第 95 百分位低於 6 ms，峰值記憶體 8 MB（用 hyperfine 執行 [`scripts/perf-check.sh`](https://github.com/jhyunwoo/ccam/blob/main/scripts/perf-check.sh) 300 次）。在同一台筆電上，光是啟動 Python 直譯器就要 13 ms，Node.js 要 31 ms，這還是在程式做任何事之前。

### 由群組決定共用什麼

MCP 伺服器、`settings.json`、`CLAUDE.md`、擴充功能（agents、commands、skills、外掛）和對話記錄是五個獨立的項目。群組會指定它共用的項目和其中的帳號，所以兩個工作帳號可以在 `work` 群組中共用對話記錄，個人帳號保留自己的對話記錄，而三個帳號都透過 `default` 使用同樣的 MCP 伺服器。`ccam group set work --add memory` 會以一份計畫更新群組中的每個帳號；[儀表板](/features/dashboard)用核取方塊做同樣的事，並把帳號已經從另一個群組取得的項目顯示為灰色。

- MCP 伺服器在啟動時透過 `--mcp-config` 傳入，所以預設情況下，啟動時不會動到 `.claude.json`，也就是 Claude Code 執行時會改寫的那個檔案。ccam 只在明確說明的步驟中寫入該檔案，而且會在 Claude Code 的鎖定之下進行並先備份：詢問後為 Remote Control 寫入、帳號退出或重新加入 MCP 共用時，以及執行 `ccam mcp import` 時。`ccam mcp add` 一次就能為整個群組新增一個伺服器。
- 其他項目是指向群組資料夾的連結，所以在一個帳號中做的變更會立即出現在其他帳號中。啟動時不複製任何東西。
- 登入資訊、`.claude.json` 和執行狀態絕不共用。

### 每項變更都先顯示並備份

`ccam group` 和 `ccam migrate` 會先列出計畫，並在移動任何東西之前詢問你。當某個檔案在帳號和群組之間不同時，由你選擇：保留群組的版本、使用你的版本，或合併 JSON；資料夾內則保留群組的副本，並由計畫列出內容不同的檔案。原始檔案放到 `.ccam-backup/<time>/`，被移除的帳號會先登出，其資料夾再放進垃圾桶資料夾，檔案以一次重新命名完成替換，所以寫入中斷時舊檔案仍會保留。`--dry-run` 會顯示啟動、共用、遷移、快捷指令、更新和修復將做什麼。`ccam doctor --fix` 會修復損壞的連結和權限。

### 在下一個帳號中接續同一個對話

`ccam next --auto --continue` 會同時查詢所有帳號的用量，開啟剩餘額度最多的帳號並接續對話。共用對話記錄時，接續的是同一個對話；沒有共用時，Claude Code 會把對話複製到新的帳號（`--resume <file> --fork-session`），ccam 也會告訴你第一次回覆會重新讀取多少 token。

當用量上限讓工作階段停止時，ccam 外掛的 hook 會顯示通知，`/ccam:switch` 則會把對話交接出去。預設情況下，何時切換由你決定。開啟 `auto_switch` 時（Linux、macOS），ccam 會自行移動工作階段，但只會在 Claude Code 因上限而停止它之後，而且只會移到和原帳號同類的帳號：同一個組織的另一個席位，或另一個個人方案（[詳細說明](/features/next#auto-switch)）。

沒有任何東西在背景監看你的用量，而且預設情況下，自動選擇絕不會把工作移出 Team 或 Enterprise 組織，也不會移入其中，所以公司的對話絕不會意外落到私人帳號上。帳號達到 90% 時，狀態列會建議這個指令。

指令碼無法等你，所以 `ccam --fallback work -p "…"` 會替它們切換：帳號達到上限時，同一個指令會以剩餘額度最多的帳號再次執行，而且只會輸出完成的那次執行的結果。不是用量上限造成的失敗則原樣傳回。對於平行的工作程序，`ccam pick` 會給每個工作程序一個不同且還有額度的帳號；沒有任何帳號還有額度時，兩者都會以代碼 75 結束。

### 每個資料夾使用正確的帳號

`ccam folder work ~/code/company` 會讓在該資料夾中未指定帳號名稱的每次啟動都使用 `work`，其他專案則繼續使用預設帳號。這項設定存放在每次啟動都會讀取的同一個 `config.toml` 中，所以不會多花任何時間。

### 單一執行檔，經過驗證的下載

ccam 是一個約 10 MB 的靜態執行檔，不需要安裝任何執行環境，支援 Linux、macOS 和 Windows。安裝指令碼（`install.sh`、`install.ps1`）和 `ccam update` 會用發行版本的 SHA-256 校驗碼檢查每個下載，ccam 也絕不會未經詢問就更新。它沒有遙測：它發出的請求只有向 Anthropic 查詢你帳號的用量（平行處理，3 秒逾時，60 秒快取）、儀表板開啟時向 GitHub 檢查是否有較新的發行版本，以及你要求的下載。設定 `usage_source = "claude"` 時，連這些查詢也會透過 Claude Code 進行，ccam 不會傳送任何 token。

## 其他工具做得更多的地方

- claude-swap、clauth、cux、CAAM 和 CCSwitcher 可以在帳號超過用量門檻時，自行切換執行中的互動式工作階段。ccam 以 `--fallback` 切換 `-p` 執行；只有開啟 `auto_switch`（Linux、macOS）時，ccam 才會自行移動互動式工作階段，而且只在 Claude Code 因上限而停止它之後；否則由你自己移動，ccam 會在達到上限時顯示通知，並提供 `/ccam:switch`。
- 在未開啟開發人員模式的 Windows 上，ccam 可以共用資料夾，但無法共用 `settings.json` 這類單一檔案（此時第一次 `ccam add` 只會共用 MCP 伺服器和擴充功能）；claude-swap 和 CCS 則會改為複製這些檔案。
- CCSwitcher、Claude Switcher 和 claude-swap 會在 macOS 選單列中顯示用量。
- CCS、clausona、CAAM、clauth、ccswap 和 Claude Switcher 也能管理 Codex 和其他程式設計代理。
- clausona 還能依帳號開啟或關閉擴充功能（ccam 的 `ccam ext` 只會顯示它們）、自行追蹤本機的權杖費用（ccam 把這件事交給 [ccusage](https://ccusage.com)：`ccam ccusage`），並提供在 tmux 和其他終端機中執行大量工作程序的 Claude Code 外掛。
- CAAM、CCSwitcher 和 Claude Switcher 可以用 Homebrew 安裝，clauth 和 CAAM 會為發行的校驗碼簽章。ccam 以安裝指令碼、`.deb` 或 `go install` 安裝，會檢查校驗碼但沒有簽章。

## 哪一種適合你？

如果你一次只用一個帳號，或希望執行中的工作階段不必重新啟動就能切換帳號（甚至在達到上限之前），請選擇替換登入資訊的工具。如果你想讓多個帳號同時執行、共用一套設定、在下一個帳號接續同一個對話，而且預設不希望有任何東西在 Claude Code 旁邊執行，請選擇 ccam，在 Linux、macOS 或 Windows 上都可以。如果什麼都不需要共用，[一般的 Shell 別名](/compare/aliases)就夠了。

能否使用多個帳號，取決於適用於各帳號的條款。開始之前請先閱讀這些條款。

**Team 和 Enterprise 席位。** 自動選擇帳號（`ccam next`、`--fallback`、`ccam pick`、儀表板的 `n`，以及你開啟後的 [`auto_switch`](/features/next#auto-switch)）也會在同一個 Claude Team 或 Enterprise 組織的席位之間切換，這可能會違反 Anthropic 的條款。在依賴這項功能之前，請先閱讀[常見問題](/docs/faq#am-i-allowed-to-use-several-claude-code-accounts)，並向你所屬的組織和 Anthropic 確認。
