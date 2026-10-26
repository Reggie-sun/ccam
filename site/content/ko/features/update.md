---
title: ccam 업데이트
description: ccam update는 SHA-256을 확인한 뒤 ccam을 GitHub의 최신 릴리스로 교체합니다. --check는 확인만 하고, --version은 지정한 릴리스 태그를 설치합니다.
updated: 2026-10-26
order: 10
translation: machine
sourceHash: "598931e4801f"
---

`ccam update`는 GitHub에서 ccam 릴리스를 내려받아 릴리스의 `checksums.txt`로 확인한 뒤, 실행 중인 바이너리를 이름 바꾸기 한 번으로 교체합니다. 설치 스크립트나 압축 파일로 ccam을 설치했을 때 쓰면 됩니다. [대시보드](/features/dashboard)도 열릴 때마다 GitHub에 최신 릴리스를 묻고 `U`로 업데이트를 제안하며, 이 키는 확인을 받은 뒤 같은 단계를 실행합니다. 그 밖의 명령은 GitHub에 접속하지 않으므로 오프라인에서도 똑같이 동작합니다.

## 사용법

```sh
ccam update                    # install the latest release
ccam update --check            # only say whether a newer release exists
ccam update --version v0.3.0   # install that release, even an older one
ccam --dry-run update          # show what would be downloaded and replaced
```

```console
$ ccam update --check
ccam 0.5.0 is available (you have 0.4.0). Run `ccam update`.
$ ccam --dry-run update
Would update ccam 0.4.0 → 0.5.0:
  download  https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
  verify    with https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/checksums.txt
  replace   /home/you/.local/bin/ccam
$ ccam update
→ downloading https://github.com/jhyunwoo/ccam/releases/download/v0.5.0/ccam_linux_amd64.tar.gz
✓ Updated ccam 0.4.0 → 0.5.0 (/home/you/.local/bin/ccam)
```

더 새로운 릴리스가 없으면 ccam은 `ccam 0.4.0 is up to date.`를 출력하고 아무것도 바꾸지 않습니다. `--version`에 이미 설치된 버전을 지정하면 `ccam 0.4.0 is already installed.`를 출력합니다.

## 옵션

| 옵션 | 동작 | 기본값 |
| --- | --- | --- |
| `--check` | 더 새로운 릴리스가 있는지만 알려 줍니다 | 꺼짐 |
| `--version <tag>` | `v0.3.0` 같은 릴리스 태그를 설치합니다. 현재 버전보다 새것이든 이전 것이든 상관없습니다 | 최신 |
| `--dry-run` | 지원합니다. 최신 태그를 조회하고(`--version`으로 지정하지 않은 경우), 바이너리에 같은 검사를 실행한 뒤, 내려받지 않고 다운로드, 체크섬, 교체 경로를 출력합니다 | 꺼짐 |

`--check`와 `--version`은 함께 쓸 수 없습니다. 태그는 `vMAJOR.MINOR.PATCH` 형식이어야 하며, 프리릴리스나 빌드 접미사를 붙일 수 있습니다. 앞의 `v`는 생략해도 됩니다. 그 밖의 형식은 실패합니다(예: `--version takes a release tag such as v0.3.0, not "1.2"`).

환경 변수 `CCAM_RELEASES_URL`은 미러나 테스트를 위해 `https://github.com/jhyunwoo/ccam/releases`를 대신합니다. 구조는 같아야 합니다. `<url>/latest`는 `<url>/tag/<tag>`로 리디렉션되고, 파일은 `<url>/download/<tag>/<asset>`에 있어야 합니다.

## 동작 방식

ccam은 설치된 바이너리를 건드리기 전에 내려받은 파일을 확인하고, 이름 바꾸기 한 번으로 바이너리를 교체합니다. 코드는 `internal/update`에 있으며, 이 패키지는 `ccam update`만 가져다 씁니다. 실행 경로에서는 절대 호출하지 않습니다.

1. ccam은 `os.Executable`로 실행 중인 바이너리를 찾고 `filepath.EvalSymlinks`로 심볼릭 링크를 해석합니다(`cli.updateCmd`). 교체되는 것은 이 실제 경로입니다.
2. `--version`이 없으면 ccam이 리디렉션을 따라가지 않고 `<releases>/latest`를 요청한 뒤(`Updater.Latest`), `Location` 헤더(`…/tag/v0.5.0`)에서 태그를 읽습니다. GitHub API 대신 릴리스 페이지를 쓰므로, 시간당 60회라는 API 요청 한도가 적용되지 않습니다. 조회 타임아웃은 10초입니다.
3. ccam이 엄격한 패턴으로 태그를 검사합니다(`parseVersion`). 태그는 다운로드 URL과 출력되는 셸 명령에 들어가므로, `/`, `%`, 공백, 셸 문자가 포함된 것은 통과하지 못합니다. 그다음 ccam은 이 태그를 실행 중인 버전과 비교하며(`compareVersions`), 프리릴리스는 해당 정식 릴리스보다 앞선 버전으로 취급합니다.
4. ccam이 바이너리를 교체해도 되는지 판단합니다(`managed()`, `writable()`, 아래 참고).
5. ccam이 `ccam_<os>_<arch>.tar.gz`(Windows에서는 `ccam_windows_<arch>.zip`)와 `checksums.txt`를 동시에 내려받으며, 두 파일을 합쳐 제한 시간은 2분입니다(`Updater.Fetch`). 각 다운로드는 최대 256 MB입니다.
6. 다른 작업을 하기 전에 압축 파일의 SHA-256이 `checksums.txt`의 해당 줄과 일치해야 합니다. 줄이 없거나 체크섬이 다르면 `nothing was changed`를 출력하고 업데이트를 멈춥니다.
7. ccam이 tar를 읽어 `ccam`이라는 이름의 첫 번째 일반 파일을 꺼냅니다(`extract()`. `extractZip`은 zip에서 `ccam.exe`를 꺼냅니다). 디렉터리, 링크 등 다른 항목은 무시하며, 바이너리도 최대 256 MB입니다.
8. ccam이 `fsutil.AtomicWrite`로 새 바이너리를 기록합니다(`Updater.Install`). 바이너리 옆에 이전 파일과 같은 권한 비트로 임시 파일 `.ccam.tmp*`를 만들고, 디스크에 기록한 뒤 이전 파일 위로 이름을 바꿉니다.

태그 조회부터 이름 바꾸기까지의 전체 과정은 다음과 같습니다.

<!-- motion:update-pipeline -->

압축 파일과 `checksums.txt`를 병렬로 내려받고, 체크섬이 일치하기 전에는 아무것도 기록하지 않습니다. 이름 바꾸기는 원자적이므로, 도중에 비정상 종료되거나 디스크가 가득 차도 이전 바이너리나 새 바이너리 중 하나가 온전히 남습니다. 절반만 바뀐 파일은 생기지 않습니다. `~/.local/bin/ccwork` 같은 단축 명령은 ccam 바이너리를 가리키는 심볼릭 링크이므로, 아무것도 바꾸지 않아도 새 버전을 실행합니다.

Windows에서는 실행 중인 프로그램의 파일을 교체하거나 지울 수 없고, 이름만 바꿀 수 있습니다. ccam은 새 바이너리를 `ccam.exe.new`에 쓰고, 실행 중인 `ccam.exe`의 이름을 `ccam.exe.old`로 바꾼 뒤, 새 파일의 이름을 `ccam.exe`로 바꿉니다(`Updater.Install`). 마지막 단계가 실패하면 이전 바이너리의 이름을 원래대로 되돌립니다. 다음 업데이트는 지울 수 있는 `ccam.exe.old` 파일을 지웁니다. 이전의 `ccam.exe.old`로 실행 중인 세션이 아직 있으면, 실행 중인 바이너리는 대신 `ccam.exe.old-<time>`이 됩니다. 심볼릭 링크가 아닌 하드 링크로 된 단축 명령([단축 명령](/features/shortcuts))은 바로 이어서 새로 만들어지므로, 이들도 새 버전을 실행합니다. 방금 교체한 바이너리의 하드 링크뿐 아니라, ccam이 만든 이전 버전 ccam의 하드 링크도 모두 새로 만듭니다.

### 패키지로 설치한 경우

`.deb`로 설치한 ccam 사본은 apt가 관리하므로, apt 모르게 교체하면 패키지 데이터베이스가 실제 파일과 맞지 않게 됩니다. ccam은 내려받기 전에 바이너리를 누가 관리하는지 확인합니다.

<!-- motion:update-managed -->

1. `/usr/bin`, `/bin`, `/usr/sbin`, `/sbin`에 있는 바이너리는 패키지가 관리하는 것으로 봅니다.
2. 그 밖의 위치라도 `dpkg`가 설치되어 있고 `dpkg -S <path>`가 성공하면 패키지가 관리하는 것으로 봅니다.
3. 둘 다 아니면 ccam이 바이너리 옆에 임시 파일 `.ccam-update-*`를 만들었다 지워 보면서 그 폴더에 쓸 수 있는지 확인합니다.

패키지가 관리하는 바이너리라면 ccam은 멈추고 apt로 업데이트하는 방법을 출력합니다.

```console
$ ccam update
ccam: /usr/bin/ccam was installed by a package, so `ccam update` leaves it alone.
To update the .deb, download ccam_0.5.0_linux_amd64.deb from https://github.com/jhyunwoo/ccam/releases/tag/v0.5.0 and run:
  sudo apt install ./ccam_0.5.0_linux_amd64.deb
Or install a copy in ~/.local/bin with the install script:
  curl -fsSL https://ccam.dev/install.sh | sh
```

Linux가 아닌 시스템에서는 설치 스크립트 줄만 출력합니다. Windows에서는 `irm https://ccam.dev/install.ps1 | iex`입니다. 폴더에 쓸 수 없으면 ccam은 설치 스크립트를 다시 실행하라고 안내합니다. 설치 스크립트는 ccam을 `~/.local/bin`에 설치합니다.

## 안전장치와 예외 상황

- 개발 빌드는 버전을 `dev`로 표시하며 비교할 릴리스가 없으므로, `ccam update`와 `ccam update --check`가 거부합니다. `ccam update --version vX.Y.Z`로는 그 위에 릴리스를 설치할 수 있습니다. `go install`로 만든 빌드는 [ccam 설치](/docs/install#other-ways-to-install)를 참고하세요.
- 다운로드가 실패하거나, 사용 중인 플랫폼용 파일이 없거나, 체크섬이 일치하지 않으면 설치된 바이너리는 그대로 남습니다.
- 없는 릴리스를 지정하면 `there is no release <tag>`로 실패합니다. GitHub가 HTTP 403이나 429로 거부하면, ccam은 요청 한도 때문일 가능성이 높다고 알리고 나중에 다시 시도하라고 안내합니다.
- `--check`는 읽기만 합니다. 바이너리가 패키지로 관리되는지, 쓸 수 있는지는 확인하지 않습니다.
- 이전 버전으로 되돌리려면 `--version`에 이전 태그를 주면 됩니다.

## 관련 문서

- [ccam 설치](/docs/install)
- [명령어 레퍼런스](/docs/commands)
- [자주 묻는 질문](/docs/faq)
