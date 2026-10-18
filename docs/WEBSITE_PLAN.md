# ccam.dev 웹사이트 구현 계획

2026-10-18 작성 · 대상: `https://ccam.dev` · 프레임워크: vinext 1.x · 배포: Cloudflare Workers

## 1. 목표와 성공 기준

ccam.dev는 ccam(Claude Code Account Manager)을 소개하고 설치·사용법을 알려주는 다국어 정적 사이트입니다. 세 가지를 동시에 달성하는 것이 목표입니다.

1. **매우 빠른 사이트**: 모든 페이지를 빌드 때 미리 만들어 Cloudflare 엣지에서 정적 파일로 바로 응답합니다. 클라이언트 JavaScript는 꼭 필요한 곳에만 씁니다.
2. **검색·LLM 노출**: "Claude Code 여러 계정 사용"과 관련된 검색에서 ccam이 상위에 나오고, ChatGPT·Claude·Perplexity·Google AI 개요 같은 LLM 답변에도 인용되도록 합니다.
3. **전 세계 언어 지원**: 30개 언어로 세계 인구의 80% 이상이 모국어 또는 생활 언어로 읽을 수 있게 합니다. 실제 커버리지는 계산 스크립트로 검증합니다(5.2절).

| 지표 | 목표 | 측정 |
| --- | --- | --- |
| LCP (모바일 p75) | 1.0 s 이하 | Cloudflare Web Analytics RUM, Lighthouse CI |
| INP (p75) | 100 ms 이하 | RUM |
| CLS (p75) | 0.02 이하 | RUM |
| TTFB (엣지 캐시 적중) | 100 ms 이하 | RUM, WebPageTest |
| 랜딩 페이지 JS (gzip) | 25 KB 이하 | 빌드 예산 검사 |
| 랜딩 페이지 HTML+CSS (gzip) | 40 KB 이하 | 빌드 예산 검사 |
| Lighthouse 4개 항목 | 모두 100 | Lighthouse CI (모바일) |
| 색인된 언어별 페이지 | 출시 4주 안에 대상 페이지의 90% | Google Search Console, Bing Webmaster Tools |
| 핵심 영어 검색어 순위 | 출시 3개월 안에 1페이지 진입 (목표, 보장 불가) | Search Console 실적 보고서 |

검색 순위는 보장할 수 없습니다. 이 계획은 순위에 영향을 주는 요소(콘텐츠 적합도, 속도, 구조화 데이터, 다국어 신호, 외부 링크)를 모두 갖추는 데 집중하고, 출시 후 데이터를 보며 개선합니다.

범위 밖: 회원 기능, 서버 데이터베이스, 결제, 광고.

## 2. 기술 현황 확인 (2026년 10월 기준)

계획은 아래 사실에 맞춰 세웠습니다. 둘 다 빠르게 바뀌고 있어서 Phase 0에서 설치 버전으로 다시 확인합니다.

**vinext 1.0** (2026-09-28 정식 출시, [Cloudflare 블로그](https://blog.cloudflare.com/vinext-nextjs-on-vite)): Vite 위에서 Next.js API를 다시 구현한 프레임워크입니다. [README](https://github.com/cloudflare/vinext) 기준 이 사이트에 중요한 지원 현황은 다음과 같습니다.

| 기능 | vinext 상태 | 이 사이트의 대응 |
| --- | --- | --- |
| App Router, RSC, `generateStaticParams` | 지원 | 그대로 사용 |
| 정적 생성, `output: "export"` | 지원 | 모든 페이지 빌드 시 생성 |
| Metadata API (`generateMetadata`, viewport) | 지원 | 그대로 사용 |
| `proxy.ts` / middleware | 지원 | 쓰지 않음 (정적 응답 유지) |
| `next/image` | 부분 지원 (빌드 시 최적화 없음) | 자체 빌드 스크립트로 AVIF/WebP 생성 후 `<img>` |
| `next/font` | 부분 지원 (자체 호스팅·서브셋 없음) | 폰트 파일을 직접 서브셋해서 자체 호스팅 |
| i18n 라우팅 | 부분 지원 (Pages Router 위주) | App Router `[locale]` 세그먼트로 직접 구현 |
| `"use cache"`, Cache Components | 부분 지원 | 쓰지 않음 (전부 빌드 시 정적 생성) |
| PPR | 미지원 | 필요 없음 (동적 영역 없음) |
| `satori`, `resvg`, `sharp` 등 네이티브 모듈 | App Router 개발 중 실패 가능 | 별도 Node 스크립트에서만 실행 |

**Next.js 16.4** (2026-10-06, [릴리스 노트](https://nextjs.org/blog/next-16-4)): Cache Components가 모든 앱에 권장되고, 라우트를 반드시 정적으로 유지하는 `ensureStatic`이 추가됐습니다. React 19.3에서 View Transitions가 안정화됐습니다. vinext는 Cache Components와 `ensureStatic`을 아직 온전히 지원하지 않으므로, 이 사이트는 "모든 페이지를 빌드 시 정적으로 만든다"는 같은 목표를 `generateStaticParams` + `dynamicParams = false` + 빌드 검사로 달성합니다. vinext가 `ensureStatic`을 지원하면 루트 레이아웃에 `ensureStatic = 'navigation'`을 추가합니다.

**Cloudflare AI 크롤러 기본값** ([정리 글](https://suganthan.com/notes/cloudflare-ai-crawler-defaults/)): 2026-09-15부터 새 도메인은 광고가 있는 페이지에서 Agent·Training 봇을 기본 차단합니다. Googlebot·Bingbot처럼 검색과 학습을 겸하는 크롤러는 가장 엄격한 규칙으로 판단되므로, Training 범주를 막으면 Google 검색에서도 빠질 수 있습니다. 관리형 robots.txt는 기본으로 `use=reference` 신호를 붙입니다. 9.4절에서 설정 방법을 정합니다.

## 3. 전체 아키텍처

모든 HTML을 빌드 시 생성해 Workers Static Assets로 배포하고, 아주 작은 Worker는 정적 파일로 처리할 수 없는 몇 가지 경로만 맡습니다.

```
            ┌───────────────────────────── 빌드 (GitHub Actions) ─────────────────────────────┐
 content/   │ vinext build (output: "export")  →  out/<locale>/**/index.html                   │
 (MDX, JSON)│ scripts/og.ts (satori+resvg)     →  out/og/<locale>/<page>.png                   │
 messages/  │ scripts/images.ts (sharp)        →  out/img/*.avif|webp                          │
            │ scripts/llms.ts                  →  out/llms.txt, llms-full.txt, <page>.md        │
            │ scripts/sitemap.ts               →  out/sitemap-index.xml, sitemaps/<locale>.xml │
            │ 예산·hreflang·링크 검사          →  실패 시 배포 중단                            │
            └──────────────────────────────────────────┬──────────────────────────────────────┘
                                                       │ wrangler deploy
            ┌──────────────────────────── Cloudflare (ccam.dev) ──────────────────────────────┐
 사용자 ───▶│ 엣지 캐시 ─▶ Workers Static Assets (HTML, CSS, JS, 이미지, md, xml)             │
 크롤러 ───▶│        └─▶ 작은 Worker: /install.sh, /latest(릴리스 리다이렉트), 헤더 보정        │
            └─────────────────────────────────────────────────────────────────────────────────┘
```

- **정적 우선**: 페이지 응답에 서버 연산이 없습니다. TTFB가 낮고, 비용이 거의 들지 않으며, 장애 지점이 적습니다.
- **Worker의 역할**: `curl -fsSL https://ccam.dev/install.sh | sh`용 스크립트 제공, GitHub 최신 릴리스로 리다이렉트, 언어 추천 배너용 `Vary` 없는 헤더 보정 정도로 제한합니다.
- **배포 방식 결정**: 기본안은 `output: "export"` 결과물을 `wrangler deploy`(Static Assets)로 올리는 것입니다. 대안은 `npx @vinext/cloudflare deploy --warm-cache`로 vinext 런타임과 CDN 캐시를 쓰는 방식입니다. Phase 0에서 둘 다 배포해 보고 TTFB·운영 단순성으로 고릅니다.

## 4. 정보 구조와 콘텐츠

### 4.1 URL 구조

- 언어는 경로 앞에 둡니다: `/en/…`, `/ko/…`, `/zh-hans/…`, `/pt-br/…`
- 루트 `/`는 영어 페이지를 그대로 보여주고 `x-default`로 지정합니다. 언어에 따른 자동 리다이렉트는 하지 않습니다. 자동 리다이렉트는 크롤러가 다른 언어 페이지를 못 보게 만들 수 있고, 정적 응답도 깨기 때문입니다.
- 브라우저 언어가 현재 페이지 언어와 다르면 작은 배너로 "한국어로 보기"를 제안합니다(클라이언트에서 `navigator.languages`만 보고 판단, 선택은 `localStorage`에 기억).
- 슬러그는 모든 언어에서 영어로 통일합니다(`/ko/docs/install`). 언어 간 연결과 유지보수가 쉬워지고, 검색 순위에는 페이지 내용과 hreflang이 더 중요합니다.

### 4.2 페이지 목록

| 경로 | 목적 | 핵심 검색 의도 |
| --- | --- | --- |
| `/` | 랜딩: 한 줄 설명, 설치 명령, 30초 데모, 기능 요약 | "claude code multiple accounts", "claude code account switcher" |
| `/docs/install` | 설치 (install.sh, .deb, 소스 빌드) | "install ccam" |
| `/docs/quick-start` | 계정 2개 추가하고 전환하기 | "switch claude code account" |
| `/docs/commands` | 명령어 전체 레퍼런스 | "ccam usage", "ccam next" |
| `/docs/sharing` | MCP·설정·대화 내역 공유 | "share mcp servers claude code accounts" |
| `/docs/usage-limits` | 사용량 조회와 `next --auto` | "claude code usage limit switch account" |
| `/docs/faq` | 자주 묻는 질문 | 다양한 롱테일 질문 |
| `/guides/multiple-claude-code-accounts` | 원리 설명 (ccam 없이 하는 방법 포함) | "how to use two claude code accounts" |
| `/guides/claude-config-dir` | `CLAUDE_CONFIG_DIR` 설명 | "CLAUDE_CONFIG_DIR" |
| `/guides/usage-limit-workflow` | 한도에 걸렸을 때 이어서 작업하기 | "claude code usage limit reached what to do" |
| `/compare/aliases` | 셸 별칭 방식과 비교 | "claude code alias multiple accounts" |
| `/changelog` | 릴리스 기록 (GitHub Releases에서 빌드 시 생성) | 브랜드 검색 |
| `/blog/*` | 출시 글, 기술 글 | 외부 공유용 |

가이드 페이지는 ccam을 쓰지 않는 방법(셸 별칭)도 정직하게 설명합니다. 검색 의도에 정확히 답하는 페이지가 순위와 LLM 인용에 유리하고, 그 위에 "이 과정을 자동화한 것이 ccam"이라고 연결합니다.

### 4.3 콘텐츠 작성 원칙

- **답을 먼저**: 각 페이지 첫 문단에 질문에 대한 답을 한두 문장으로 씁니다. 검색 결과 스니펫과 LLM 답변이 이 부분을 가져갑니다.
- **사실은 구체적으로**: 명령어, 경로, 버전, 동작 조건을 정확히 씁니다. 모호한 마케팅 문구는 피합니다.
- **비공식 명시**: 모든 페이지 푸터와 소개 문단에 "Anthropic과 관련 없는 비공식 오픈소스 도구"라고 밝히고, Claude 로고나 Anthropic 브랜드 요소는 쓰지 않습니다. "Claude"는 상표이므로 제품명은 `ccam`으로 앞세웁니다.
- **약관 안내**: 여러 계정 사용은 각 계정 약관을 사용자가 확인해야 한다는 문장을 FAQ와 가이드에 넣습니다.
- **README와 일치**: GitHub README, 사이트, `llms.txt`의 핵심 설명 문장을 똑같이 맞춥니다. LLM은 GitHub를 많이 학습하므로 같은 문장이 여러 곳에 있으면 인용 정확도가 올라갑니다.

## 5. 다국어 (i18n)

### 5.1 대상 언어

1차 출시에 18개, 2차에 12개를 더해 총 30개 언어를 지원합니다.

| 단계 | 언어 (URL 코드) |
| --- | --- |
| 1차 (출시) | English (`en`), 简体中文 (`zh-hans`), 繁體中文 (`zh-hant`), हिन्दी (`hi`), Español (`es`), العربية (`ar`), Français (`fr`), বাংলা (`bn`), Português (`pt-br`), Русский (`ru`), Bahasa Indonesia (`id`), اردو (`ur`), Deutsch (`de`), 日本語 (`ja`), 한국어 (`ko`), Tiếng Việt (`vi`), Türkçe (`tr`), Italiano (`it`) |
| 2차 (출시 후 4주 안) | मराठी (`mr`), తెలుగు (`te`), தமிழ் (`ta`), ગુજરાતી (`gu`), ಕನ್ನಡ (`kn`), فارسی (`fa`), Kiswahili (`sw`), Filipino (`fil`), ไทย (`th`), Polski (`pl`), Українська (`uk`), Hausa (`ha`) |

인도 지역 언어(마라티어·텔루구어·타밀어·구자라트어·칸나다어)를 2차에 넣은 이유는 인도 인구의 상당수가 힌디어나 영어를 주 언어로 쓰지 않기 때문입니다. 이 언어들이 빠지면 80% 기준을 넘기 어렵습니다.

### 5.2 커버리지 검증

"인류의 80%"를 감이 아니라 숫자로 확인합니다. `scripts/coverage.ts`가 Unicode CLDR의 `supplementalData`(국가별 인구, 국가별 언어 사용 비율)를 읽어 두 값을 계산합니다.

- **보수적 추정(하한)**: 나라마다 지원 언어 중 가장 많이 쓰이는 언어의 비율만 셉니다. 여러 언어 사용자를 중복으로 세지 않습니다.
- **낙관적 추정(상한)**: 나라마다 지원 언어 비율을 모두 더하고 100%에서 자릅니다.

출시 기준은 **보수적 추정 80% 이상**입니다. 결과는 빌드 산출물(`coverage.json`)과 사이트의 `/languages` 페이지에 공개합니다. 기준에 못 미치면 기여도가 가장 큰 언어부터 추가합니다.

### 5.3 구현 방식

- **라우팅**: `app/[locale]/…` 세그먼트, `generateStaticParams`로 모든 언어를 미리 생성, `dynamicParams = false`로 목록 밖의 언어는 404입니다.
- **문자열**: `messages/<locale>.json`에 ICU MessageFormat으로 둡니다. 서버 컴포넌트에서만 불러오는 작은 자체 모듈(`lib/i18n.ts`)을 씁니다. `getDictionary(locale)`, `t(key, values)`, `Intl.PluralRules`·`Intl.NumberFormat`·`Intl.DateTimeFormat`을 제공합니다. 번역 문자열이 클라이언트 JS에 들어가지 않으므로 언어 수가 늘어도 번들 크기는 그대로입니다. 키는 영어 파일에서 TypeScript 타입을 만들어 오타를 컴파일 단계에서 잡습니다.
- **본문**: `content/<locale>/<section>/<slug>.mdx`. 번역이 없는 페이지는 영어로 대체하지 않고 그 언어에서 아예 만들지 않습니다. 반쯤 영어인 페이지는 품질 신호를 떨어뜨리기 때문입니다. hreflang 목록에도 실제로 있는 언어만 넣습니다.
- **오른쪽에서 왼쪽 언어**: `ar`, `ur`, `fa`는 `<html dir="rtl">`로 렌더링합니다. CSS는 처음부터 논리 속성(`margin-inline-start`, `padding-block`)만 씁니다. 코드 블록과 명령어는 `dir="ltr"`로 고정합니다.
- **글꼴**: 라틴 문자용 가변 글꼴 하나만 서브셋해서 자체 호스팅합니다. 한중일·아랍·인도 문자는 운영체제 기본 글꼴(`system-ui` 스택과 언어별 `font-family` 지정)을 씁니다. 수 MB짜리 한중일 웹폰트를 내려받지 않아도 되어 LCP가 크게 줄어듭니다.
- **언어 전환기**: 각 언어 이름을 그 언어로 표시하고(`Intl.DisplayNames`), 현재 페이지의 같은 경로로 이동합니다. JS 없이 `<details>` 기반으로 만들어 클라이언트 코드를 늘리지 않습니다.

### 5.4 번역 작업 흐름

1. 영어(`en`)가 원본입니다. 원본이 바뀌면 문단 단위 해시가 바뀌고, CI가 "번역이 오래된 언어"를 표시합니다.
2. Claude Code 서브에이전트가 번역 초안을 만듭니다. 이때 다음을 함께 넘깁니다.
   - **용어집**(`i18n/glossary.yaml`): `ccam`, `Claude Code`, `MCP`, 명령어, 경로는 번역하지 않습니다.
   - **언어별 검색어 목록**(`i18n/keywords/<locale>.yaml`)
   - **문체 지침**(존댓말 수준, 용어 선택 등)
3. 자동 검사를 거칩니다.
   - ICU 자리표시자·마크다운 구조·링크·코드 블록이 원본과 일치하는지 확인합니다.
   - 길이가 지나치게 다르면 경고합니다.
   - 일부 문단은 영어로 되번역해 의미 차이를 점검합니다.
4. 1차 언어 중 사용자가 많은 언어(중국어, 스페인어, 일본어, 한국어, 포르투갈어, 프랑스어, 독일어)는 원어민이 검수합니다. GitHub 이슈로 기여자를 모집합니다.
5. 각 페이지 front matter에 `translation: machine | reviewed` 상태를 기록합니다. 검수 전에도 공개하되, 사용자가 번역 오류를 바로 고칠 수 있도록 GitHub 편집 링크를 둡니다.

언어별 검색어 예시 (Phase 0에서 각 언어 검색 결과와 자동완성으로 다듬습니다):

| 언어 | 대표 검색어 |
| --- | --- |
| en | claude code multiple accounts, switch claude code account, CLAUDE_CONFIG_DIR |
| ko | 클로드 코드 여러 계정, 클로드 코드 계정 전환 |
| ja | Claude Code 複数アカウント, Claude Code アカウント 切り替え |
| zh-hans | Claude Code 多账号, Claude Code 切换账号 |
| es | Claude Code varias cuentas, cambiar cuenta Claude Code |
| vi | Claude Code nhiều tài khoản, đổi tài khoản Claude Code |
| fr | Claude Code plusieurs comptes, changer de compte Claude Code |

## 6. 성능 설계

### 6.1 렌더링

- **완전 정적**: 모든 페이지가 `out/` 아래 HTML로 생성됩니다. 빌드 후 검사 스크립트가 서버 렌더링이 필요한 라우트가 하나라도 있으면 실패시킵니다(vinext가 `ensureStatic`을 지원하기 전까지 대신하는 장치).
- **서버 컴포넌트 기본**: 클라이언트 컴포넌트는 세 개만 둡니다.
  - 설치 명령 복사 버튼
  - 언어 추천 배너
  - 테마 전환
  
  셋 다 작은 아일랜드로 만들고, 화면에 보일 때 하이드레이션합니다.
- **JS 없이도 동작**: 메뉴, 언어 전환기, FAQ 접기는 `<details>`/`<summary>`와 CSS로 구현합니다. JS가 꺼져 있어도 모든 내용과 이동이 됩니다.
- **페이지 이동**:
  - Speculation Rules(`<script type="speculationrules">`)로 링크에 마우스를 올리거나 손가락을 대는 순간 다음 페이지를 미리 렌더링합니다(`eagerness: "moderate"`).
  - 같은 사이트 안 이동은 CSS View Transitions(`@view-transition { navigation: auto; }`)로 부드럽게 처리합니다.
  - 둘 다 정적 MPA에서 동작하므로 클라이언트 라우터 JS가 필요 없습니다.

### 6.2 자원 최적화

| 자원 | 방법 |
| --- | --- |
| CSS | Tailwind CSS v4로 사용한 클래스만 생성. 전체 CSS가 15 KB(gzip) 이하이면 `<head>`에 인라인, 넘으면 첫 화면용 CSS만 인라인 |
| 글꼴 | 라틴 가변 글꼴 1개를 `pyftsubset`으로 서브셋한 WOFF2(목표 25 KB 이하), `preload`, `font-display: swap`, 대체 글꼴에 `size-adjust`를 맞춰 CLS 방지 |
| 이미지 | `scripts/images.ts`(sharp)로 AVIF·WebP·여러 폭 생성, `<img srcset sizes width height loading decoding>` 직접 작성, 첫 화면 이미지만 `fetchpriority="high"` |
| 데모 영상 | 터미널 데모는 영상 대신 asciinema 형식 데이터를 작은 SVG 애니메이션으로 변환(수십 KB) |
| 아이콘 | 인라인 SVG 스프라이트 하나 |
| OG 이미지 | `scripts/og.ts`(satori + resvg)로 언어·페이지별 1200×630 PNG를 빌드 시 생성 |

### 6.3 엣지 전송

- **캐시 헤더**: 해시가 붙은 자원은 `Cache-Control: public, max-age=31536000, immutable`, HTML은 `public, max-age=0, s-maxage=86400, stale-while-revalidate=604800`로 `_headers`에 지정합니다. 배포할 때 HTML 캐시를 비웁니다.
- **103 Early Hints**: HTML 응답의 `Link: rel=preload` 헤더(글꼴, CSS)를 Cloudflare Early Hints가 103으로 먼저 보냅니다.
- **압축**: Cloudflare 자동 Brotli·Zstandard를 씁니다.
- **HTTPS**: `.dev`는 HSTS 사전 등록 도메인이라 HTTPS가 강제됩니다. HTTP/3를 켭니다.
- **보안 헤더**: 엄격한 CSP(nonce 대신 해시 기반, 외부 스크립트 없음), `X-Content-Type-Options`, `Referrer-Policy`를 `_headers`로 넣습니다.

### 6.4 성능 예산 검사

`scripts/budget.ts`가 빌드 결과를 페이지별로 검사하고, 기준을 넘으면 CI를 실패시킵니다.

- 랜딩 JS 25 KB, 문서 페이지 JS 15 KB, HTML+CSS 40 KB (모두 gzip)
- 글꼴 요청 1개, 첫 화면 이미지 1개 이하

Lighthouse CI는 대표 페이지 10개 × 언어 5개(en, zh-hans, ar, hi, ko)를 모바일 조건으로 돌립니다. 아랍어는 RTL, 힌디어는 인도 문자, 중국어·한국어는 한중일 글꼴을 대표합니다.

## 7. 검색엔진 최적화 (SEO)

### 7.1 페이지별 메타데이터

`generateMetadata`에서 언어별로 다음을 만듭니다.

- `<title>`: 60자 안팎, 핵심 검색어를 앞에 둡니다. 예: "Use Multiple Claude Code Accounts on One Machine — ccam"
- `<meta name="description">`: 150자 안팎, 해당 언어로 작성
- `<link rel="canonical">`: 자기 자신
- `<link rel="alternate" hreflang="…">`: 실제로 있는 모든 언어와 `x-default`(영어). 양방향으로 서로 가리키게 만듭니다.
- Open Graph·Twitter 카드: 언어별 OG 이미지, `og:locale`, `og:locale:alternate`
- `<html lang>`과 `dir`

### 7.2 구조화 데이터 (JSON-LD)

| 페이지 | 스키마 |
| --- | --- |
| 전체 | `WebSite`, `Organization`(오픈소스 프로젝트 정보, GitHub `sameAs`) |
| 랜딩 | `SoftwareApplication` (`applicationCategory: DeveloperApplication`, `operatingSystem: Ubuntu 22.04, 24.04`, `offers: price 0`, `softwareVersion`, `downloadUrl`) |
| 문서·가이드 | `TechArticle` + `BreadcrumbList` |
| FAQ | `FAQPage` (Google 리치 결과 노출은 제한적이지만, 검색엔진과 LLM이 질문·답 구조를 이해하는 데 도움) |

빌드 검사가 JSON-LD 문법을 검증하고, 출시 전 Google 리치 결과 테스트로 대표 페이지를 확인합니다.

### 7.3 크롤링과 색인

- **사이트맵**: `sitemap-index.xml` 아래 언어별 사이트맵을 둡니다. 각 URL에 `xhtml:link` 대체 언어와 `lastmod`(콘텐츠 해시가 바뀐 날짜)를 넣습니다.
- **robots.txt**: 직접 작성해 배포합니다. 검색·AI 크롤러를 모두 허용하고 사이트맵 위치를 적습니다(9.4절).
- **색인 요청**:
  - Google Search Console과 Bing Webmaster Tools에 도메인을 등록합니다. Bing은 ChatGPT 검색과 Copilot의 바탕이므로 중요합니다.
  - 배포 후 바뀐 URL을 IndexNow로 알립니다. Cloudflare Crawler Hints를 켜거나 배포 스크립트에서 IndexNow API를 호출합니다.
- **내부 링크**: 모든 문서 페이지에 이전·다음 링크와 관련 가이드 링크를 둡니다. 랜딩에서 핵심 가이드 4개로 바로 연결합니다.
- **404와 리다이렉트**: 언어 코드 대소문자 차이(`/ko-KR/`, `/zh-CN/`)는 Worker에서 정식 경로로 301 리다이렉트합니다.

### 7.4 외부 신호

새 도메인은 외부 링크와 언급이 없으면 순위가 잘 오르지 않습니다. 출시 계획에 다음을 넣습니다.

- **GitHub 저장소**: 설명과 topics(`claude-code`, `account-manager`, `cli`, `ubuntu`), README 첫 줄에 ccam.dev 링크를 넣습니다.
- **커뮤니티 목록**: awesome-claude-code 같은 목록에 추가를 요청합니다.
- **출시 글**: Hacker News (Show HN), Reddit, dev.to, 언어별 개발자 커뮤니티(Qiita·Zenn, 掘金, velog 등)에 게시합니다. 홍보 문구 대신 문제 설명과 해결 과정을 씁니다.
- **기술 글**: `CLAUDE_CONFIG_DIR` 동작 원리와 사용량 조회 방식처럼 그 자체로 쓸모 있는 글을 블로그에 올려 자연스러운 링크를 모읍니다.

## 8. LLM 최적화 (GEO)

검색 결과 링크뿐 아니라 LLM 답변 안에서 ccam이 추천되도록, LLM이 읽기 쉬운 형태를 함께 제공합니다.

- **`/llms.txt`**: [llmstxt.org](https://llmstxt.org) 형식으로 프로젝트 한 줄 설명, 핵심 사실 목록, 주요 문서 링크(마크다운판)를 둡니다. 영어 기준이며, 언어별 `/ko/llms.txt` 등도 생성합니다.
- **`/llms-full.txt`**: 영어 문서 전체를 마크다운 하나로 합칩니다.
- **페이지별 마크다운판**: 모든 문서를 `/en/docs/install.md`처럼 순수 마크다운으로도 제공하고, HTML `<head>`에 `<link rel="alternate" type="text/markdown" href="…">`를 넣습니다.
- **핵심 사실 블록**: 랜딩과 FAQ에 다음 사실을 짧은 문장으로 고정해 둡니다. 모든 언어와 README에서 같은 사실을 같은 순서로 씁니다.
  - "ccam은 무엇인가"
  - "무엇을 해결하는가"
  - "어떻게 설치하는가"
  - "공식 도구인가"
- **질문형 제목**: FAQ와 가이드 제목을 실제 질문 형태로 씁니다. 예: "How do I use two Claude Code accounts on one computer?"
- **최신성 표시**: 각 문서에 `최종 수정일`과 대상 ccam·Claude Code 버전을 표시합니다. LLM과 검색엔진 모두 최신 정보를 선호합니다.
- **측정**: 매월 주요 LLM(ChatGPT, Claude, Perplexity, Gemini)에 대표 질문 20개를 언어별로 묻고, ccam이 언급·인용되는지 기록합니다(`docs/geo-log.md`).

## 9. Cloudflare 설정

### 9.1 Worker와 정적 자원

`wrangler.jsonc`:

```jsonc
{
  "name": "ccam-site",
  "compatibility_date": "2026-10-01",
  "compatibility_flags": ["nodejs_compat"],
  "main": "worker/index.ts",
  "assets": {
    "directory": "./out",
    "binding": "ASSETS",
    "html_handling": "auto-trailing-slash",
    "not_found_handling": "404-page",
    "run_worker_first": ["/install.sh", "/latest", "/latest/*"]
  },
  "routes": [{ "pattern": "ccam.dev", "custom_domain": true }],
  "observability": { "enabled": true }
}
```

- `run_worker_first`에 적은 경로만 Worker 코드가 실행되고, 나머지는 정적 자원이 바로 응답합니다. 페이지 요청에 Worker 실행 비용이 없습니다.
- `www.ccam.dev`는 `ccam.dev`로 301 리다이렉트합니다(Redirect Rules).

### 9.2 Worker 기능

- **`/install.sh`**: GitHub 최신 릴리스의 설치 스크립트를 가져와 1시간 캐시하고, `Content-Type: text/x-shellscript`로 응답합니다.
- **`/latest`, `/latest/<asset>`**: GitHub 최신 릴리스 파일로 302 리다이렉트합니다.
- **언어 코드 정규화 리다이렉트**: `/zh-CN` → `/zh-hans` 등을 처리합니다. 이 경로들도 `run_worker_first`에 추가합니다.

### 9.3 캐시와 배포

- **GitHub Actions**: main 브랜치에 푸시하면 빌드 → 검사 → `wrangler deploy` 순서로 실행합니다. PR마다 `wrangler versions upload`로 미리보기 URL을 만듭니다.
- **배포 후 처리**: HTML 경로 캐시를 지우고(Cache Purge API), 바뀐 URL을 IndexNow로 알립니다.
- **분석**: Cloudflare Web Analytics(쿠키 없음, 동의 배너 불필요)로 방문과 Core Web Vitals를 수집합니다. 별도 분석 스크립트는 넣지 않습니다.

### 9.4 크롤러 정책

Cloudflare 대시보드의 AI 크롤러 설정(Configure AI bot policies)에서 다음과 같이 둡니다.

- **Search 봇**: 허용
- **Agent 봇**: 허용 (ChatGPT-User처럼 사용자 질문에 답하려고 실시간으로 페이지를 읽는 봇)
- **Training 봇**: 허용. 이 범주를 막으면 Googlebot·Bingbot까지 막힐 수 있습니다. 또 이 사이트는 학습 데이터에 포함될수록 LLM이 ccam을 알게 되어 유리합니다.
- **관리형 robots.txt**: 끄고, 직접 작성한 robots.txt를 씁니다.
- **Bot Fight Mode**: 정상 크롤러가 막히지 않는지 출시 전 Search Console의 URL 검사와 Bing 도구로 확인합니다.

직접 작성하는 robots.txt:

```
User-agent: *
Allow: /

Sitemap: https://ccam.dev/sitemap-index.xml
```

## 10. 저장소 구조

ccam 본체 저장소 안의 `site/` 폴더에 둡니다. 명령어 레퍼런스와 변경 기록을 본체 코드에서 바로 생성하고, README와 문구를 함께 관리하기 위해서입니다.

```
site/
  app/
    [locale]/
      layout.tsx            # <html lang dir>, 공통 헤더·푸터, JSON-LD
      page.tsx              # 랜딩
      docs/[slug]/page.tsx
      guides/[slug]/page.tsx
      compare/[slug]/page.tsx
      blog/[slug]/page.tsx
      changelog/page.tsx
      languages/page.tsx    # 지원 언어와 커버리지
    page.tsx                # / (영어 랜딩, x-default)
    not-found.tsx
  components/               # 서버 컴포넌트 위주, client/ 아래에 아일랜드 3개
  content/<locale>/…        # MDX 본문
  messages/<locale>.json    # UI 문자열 (ICU)
  i18n/
    locales.ts              # 언어 목록, 방향, 표시명
    glossary.yaml
    keywords/<locale>.yaml
  lib/
    i18n.ts  seo.ts  jsonld.ts  content.ts
  scripts/
    og.ts  images.ts  llms.ts  sitemap.ts  coverage.ts
    budget.ts  check-static.ts  check-hreflang.ts  indexnow.ts
    gen-commands.ts         # ccam --help 출력에서 명령어 레퍼런스 생성
  worker/index.ts           # install.sh, latest, 리다이렉트
  public/  _headers  robots.txt
  vite.config.ts  next.config.ts  wrangler.jsonc
  tests/                    # Playwright, 접근성, 시각 회귀
```

## 11. 테스트와 품질 검사

| 검사 | 도구 | 실패 조건 |
| --- | --- | --- |
| 정적 생성 여부 | `scripts/check-static.ts` | 요청 시 렌더링이 필요한 라우트가 있음 |
| 성능 예산 | `scripts/budget.ts`, Lighthouse CI | 6.4절 기준 초과 |
| hreflang | `scripts/check-hreflang.ts` | 양방향 누락, 없는 페이지 가리킴, 코드 오류 |
| 링크 | lychee | 깨진 내부·외부 링크 |
| HTML·JSON-LD | html-validate, schema 검사 | 문법 오류 |
| 번역 | `scripts/i18n-check.ts` | 자리표시자·구조 불일치, 원본보다 오래된 번역 |
| 접근성 | axe (Playwright) | 심각도 serious 이상 |
| 시각 회귀 | Playwright 스크린샷 (en, ar, ja, hi) | 레이아웃 깨짐, RTL 뒤집힘 오류 |
| 언어 커버리지 | `scripts/coverage.ts` | 보수적 추정 80% 미만 (2차 출시 이후) |
| AI 특유의 저품질 코드·문구 | antislop 스킬 (11.1절) | 지적 사항 중 수정하지 않은 항목이 남아 있음 |

### 11.1 antislop 스킬로 전체 코드 검증

Claude Code가 작성한 코드와 문구에는 AI가 자주 남기는 저품질 패턴이 섞이기 쉽습니다. 예를 들면 의미 없는 주석, 불필요한 추상화, 뻔한 디자인, 과장된 마케팅 문구 같은 것입니다. 이를 걸러내기 위해 antislop 스킬([skills-hub](https://skills-hub.ai/skills/anti-slop-antislop))로 사이트 전체를 검증합니다.

- **설치**: 저장소 루트에서 `npx @skills-hub-ai/cli install anti-slop-antislop --target claude-code`를 실행합니다. 이 스킬은 서명·검증되지 않은 배포본입니다. 설치 후 `SKILL.md` 내용을 직접 읽어 확인하고, 확인한 버전을 저장소에 커밋해 고정합니다.
- **같은 계열 스킬**: 같은 배포처의 관련 스킬(코드 주석, UI·시각 디자인, 카피라이팅, 모바일 레이아웃, 접근성)이 있다면 함께 설치해 antislop과 같이 적용합니다.
- **검증 범위**: `site/` 아래 전체입니다.
  - 코드: `app/`, `components/`, `lib/`, `scripts/`, `worker/`
  - 스타일
  - 영어 원본 콘텐츠(`content/en/`, `messages/en.json`)
  
  번역본은 원본을 고친 뒤 다시 번역하므로 영어 원본을 기준으로 검증합니다.
- **실행 시점**:
  1. 모든 Phase 완료 직전에 그 Phase에서 바뀐 파일 전체를 검증합니다.
  2. 출시 전(Phase 4 시작 전)에 사이트 전체를 한 번 더 검증합니다.
  3. 이후 큰 변경이 있는 PR마다 바뀐 파일을 검증합니다.
- **진행 방법**: Claude Code에서 `/anti-slop-antislop`을 실행하고(등록된 명령 이름은 설치 후 확인), 지적 사항을 `docs/antislop-review.md`에 파일·항목·조치(수정함/유지 이유)로 기록합니다. 모든 항목을 수정하거나 유지 이유를 적어야 Phase 완료로 봅니다.
- **다른 검사와의 순서**: antislop으로 수정한 뒤에는 11절의 자동 검사(예산, 정적 생성, hreflang, 접근성, 시각 회귀)를 다시 돌려 수정이 다른 기준을 깨지 않았는지 확인합니다.

## 12. 단계별 로드맵

모든 Phase의 완료 기준에는 "그 Phase에서 바뀐 파일의 antislop 검증 통과(11.1절)"가 공통으로 들어갑니다.

**Phase 0 — 기술 검증**

- 확인:
  - vinext 최신 버전으로 App Router `[locale]` 정적 생성과 `output: "export"`가 되는지
  - `wrangler deploy`(Static Assets)와 `@vinext/cloudflare deploy` 각각의 TTFB
  - 언어별 검색 결과와 자동완성으로 검색어 목록 정리
  - Cloudflare 크롤러 설정 화면
- 완료 기준: `docs/site-decisions.md`에 배포 방식, 검색어 목록, 확인된 vinext 제약 정리
- 프롬프트: "site/ 폴더에 vinext 최신 버전으로 최소 App Router 앱을 만들고, [locale] 세그먼트 3개 언어를 output export로 정적 생성해줘. 결과를 wrangler deploy(Static Assets)와 @vinext/cloudflare deploy 두 방식으로 미리보기 배포하고 TTFB와 제약 사항을 docs/site-decisions.md에 표로 정리해줘."

**Phase 1 — 뼈대와 영어 랜딩**

- 구현: 레이아웃, 디자인 토큰, 랜딩, 설치·빠른 시작 문서, i18n 모듈(영어만), 메타데이터·JSON-LD, 사이트맵, robots.txt, 정적 생성·예산 검사, CI 배포
- 완료 기준: ccam.dev에 영어 사이트 공개, Lighthouse 4개 항목 100, 예산 검사 통과
- 프롬프트: "WEBSITE_PLAN.md의 Phase 1을 구현해줘. 클라이언트 컴포넌트는 6.1절의 세 개만 허용하고, scripts/budget.ts와 check-static.ts를 먼저 만든 뒤 그 검사를 통과하도록 페이지를 작성해줘."

**Phase 2 — 문서 전체와 LLM 자료**

- 구현: 4.2절의 문서·가이드·비교·FAQ 페이지, 명령어 레퍼런스 자동 생성, `llms.txt`·`llms-full.txt`·페이지별 마크다운, OG 이미지 생성, Worker(`/install.sh`, `/latest`)
- 완료 기준: 영어 페이지 전체 공개, Search Console·Bing 등록과 사이트맵 제출, `curl -fsSL https://ccam.dev/install.sh | sh` 동작

**Phase 3 — 1차 다국어 (18개 언어)**

- 구현: 언어 목록·용어집·검색어 파일, 번역 파이프라인과 검사, RTL 지원, 언어 전환기와 추천 배너, hreflang 검사
- 완료 기준: 18개 언어 전체 페이지 공개, hreflang·번역 검사 통과, ar·hi·ja·ko 시각 회귀 통과
- 프롬프트: "Phase 3을 진행해줘. 영어 원본에서 번역할 때 i18n/glossary.yaml의 용어는 절대 번역하지 말고, 언어마다 서브에이전트를 하나씩 써서 병렬로 번역한 뒤 scripts/i18n-check.ts를 통과시켜줘."

**Phase 4 — 출시와 외부 신호**

- 시작 전: antislop 스킬로 사이트 전체 코드와 영어 콘텐츠를 검증하고 `docs/antislop-review.md`의 모든 항목을 처리
- 구현: 블로그 출시 글, GitHub 저장소 정비, 커뮤니티 게시, IndexNow 자동화, GEO 측정 기록 시작
- 완료 기준: 전체 antislop 검증 완료, 출시 글 게시, 대표 질문 20개의 첫 GEO 측정 기록
- 프롬프트: "출시 전에 antislop 스킬로 site/ 전체를 검증해줘. 코드, 스타일, content/en, messages/en.json을 모두 보고, 지적 사항을 docs/antislop-review.md에 파일·항목·조치로 기록한 뒤 고쳐줘. 고친 뒤 pnpm check와 Lighthouse CI를 다시 돌려 결과를 알려줘."

**Phase 5 — 2차 다국어와 커버리지 검증**

- 구현: 12개 언어 추가, `scripts/coverage.ts`, `/languages` 페이지
- 완료 기준: 보수적 추정 커버리지 80% 이상, 30개 언어 색인 시작

**Phase 6 — 운영 개선 (지속)**

- 매월:
  - Search Console 검색어 분석으로 새 가이드 주제를 정합니다.
  - 순위가 낮은 페이지의 제목·첫 문단을 고칩니다.
  - GEO 측정을 반복합니다.
  - Core Web Vitals RUM 회귀를 점검합니다.
  - 원어민 검수 언어를 늘립니다.
- vinext가 `ensureStatic`·Cache Components를 지원하면 6.1절의 자체 정적 검사를 프레임워크 기능으로 바꿉니다.

## 13. 위험 요소

| 위험 | 대응 |
| --- | --- |
| vinext 기능 부족이나 버그 | 정적 생성만 쓰는 단순한 구조로 의존 범위를 줄임. 막히면 같은 코드로 Next.js 16 `output: "export"` 빌드로 전환 가능하게 Next.js 공식 API만 사용 |
| 기계 번역 품질 | 용어집·자동 검사·되번역 점검, 주요 언어 원어민 검수, 편집 링크로 커뮤니티 수정 |
| 저품질 다국어 페이지로 인한 검색 품질 신호 하락 | 번역이 없는 페이지는 만들지 않음, 얇은 페이지 금지, 각 언어 검색어에 맞춘 제목·설명 |
| Cloudflare 크롤러 차단 | 9.4절 설정, 출시 전 Search Console URL 검사로 확인 |
| 상표 문제 | 제품명은 ccam, 비공식 명시, Anthropic 로고·디자인 미사용 |
| 새 도메인이라 순위가 늦게 오름 | 정직하고 구체적인 가이드 콘텐츠, GitHub·커뮤니티에서 자연스러운 링크 확보, 꾸준한 업데이트 |
| 인기 검색어에서 공식 문서나 큰 사이트에 밀림 | 롱테일 질문(한도 도달 시 이어서 작업, MCP 공유 등)에서 먼저 1위를 확보하고 넓혀감 |

## 14. Claude Code로 시작하기

1. ccam 저장소에 이 파일을 `docs/WEBSITE_PLAN.md`로 저장합니다.
2. `site/CLAUDE.md`를 만듭니다.

   ```markdown
   # ccam.dev 사이트

   전체 계획은 ../docs/WEBSITE_PLAN.md를 먼저 읽는다. 결정 기록은 ../docs/site-decisions.md.

   ## 규칙
   - 프레임워크는 vinext, Next.js 공식 API만 사용 (vinext 전용 API 금지)
   - 모든 페이지는 빌드 시 정적 생성. 요청 시 렌더링이 필요한 코드는 금지
   - 클라이언트 컴포넌트는 components/client/ 아래 세 개(복사 버튼, 언어 배너, 테마)만 허용
   - 번역 문자열은 서버에서만 사용, 클라이언트 번들에 넣지 않는다
   - CSS는 논리 속성만 사용 (RTL 대응)
   - 이미지·OG·사이트맵·llms.txt는 scripts/의 빌드 스크립트로 생성
   - i18n/glossary.yaml의 용어는 번역하지 않는다
   - 모든 페이지 푸터에 비공식 도구 안내를 넣는다

   ## 명령
   - 개발: pnpm vinext dev
   - 빌드: pnpm build (vinext build + scripts/*)
   - 검사: pnpm check (budget, static, hreflang, i18n, links)
   - 배포: pnpm deploy (wrangler deploy)

   ## 작업 방식
   - 한 번에 한 Phase만 진행하고 완료 기준을 직접 확인해 보고한다
   - 검사 스크립트를 먼저 만들고, 그 검사를 통과하도록 구현한다
   - Phase를 마치기 전에 antislop 스킬로 바뀐 파일을 검증하고, 결과를 docs/antislop-review.md에 기록한다
   - antislop 수정 후에는 pnpm check를 다시 실행한다
   ```

   antislop 스킬은 작업 시작 전에 설치합니다.

   ```bash
   npx @skills-hub-ai/cli install anti-slop-antislop --target claude-code
   # 설치된 SKILL.md를 직접 읽어 확인한 뒤 커밋
   ```

3. Claude Code에서 첫 프롬프트를 입력합니다.

   ```
   docs/WEBSITE_PLAN.md와 site/CLAUDE.md를 읽고 Phase 0부터 시작해줘.
   Cloudflare 배포가 필요한 단계에서는 wrangler login이 되어 있는지 확인하고,
   안 되어 있으면 멈추고 나에게 알려줘.
   ```

## 출처

- [Next.js applications, powered by Vite: introducing Vinext 1.0 — Cloudflare Blog](https://blog.cloudflare.com/vinext-nextjs-on-vite)
- [cloudflare/vinext — GitHub](https://github.com/cloudflare/vinext)
- [Next.js 16.4 — Next.js Blog](https://nextjs.org/blog/next-16-4)
- [Cloudflare just changed the AI crawler defaults, and there's a trap in it](https://suganthan.com/notes/cloudflare-ai-crawler-defaults/)
- [llms.txt 형식 제안](https://llmstxt.org)
- [anti-slop (antislop) 스킬 — skills-hub](https://skills-hub.ai/skills/anti-slop-antislop)

