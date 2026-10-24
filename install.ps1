# Installs ccam into ~\.local\bin on Windows.
#
#   irm https://raw.githubusercontent.com/jhyunwoo/ccam/main/install.ps1 | iex
#
# Environment:
#   CCAM_VERSION      release tag such as v1.0.0 (default: latest)
#   CCAM_INSTALL_DIR  where to install (default: ~\.local\bin)
#   CCAM_BASE_URL     release download base (default: GitHub releases)
#   CCAM_ZIP          install from a local .zip instead of downloading
#
# It runs inside a script block, so an error never closes the shell that
# piped it into iex.
& {
  $ErrorActionPreference = 'Stop'
  $ProgressPreference = 'SilentlyContinue' # the progress bar slows Invoke-WebRequest down
  [Net.ServicePointManager]::SecurityProtocol = [Net.ServicePointManager]::SecurityProtocol -bor [Net.SecurityProtocolType]::Tls12

  $repo = 'jhyunwoo/ccam'
  $dir = if ($env:CCAM_INSTALL_DIR) { $env:CCAM_INSTALL_DIR } else { Join-Path $HOME '.local\bin' }

  function Fail($msg) { throw "ccam install failed: $msg" }

  switch ($env:PROCESSOR_ARCHITECTURE) {
    'AMD64' { $arch = 'amd64' }
    'ARM64' { $arch = 'arm64' }
    default { Fail "unsupported CPU: $env:PROCESSOR_ARCHITECTURE" }
  }
  $asset = "ccam_windows_$arch.zip"
  $tmp = Join-Path ([IO.Path]::GetTempPath()) ("ccam-install-" + [Guid]::NewGuid())
  New-Item -ItemType Directory -Path $tmp | Out-Null
  try {
    $zip = Join-Path $tmp $asset
    if ($env:CCAM_ZIP) {
      Copy-Item $env:CCAM_ZIP $zip
    } else {
      if ($env:CCAM_VERSION) {
        $tag = $env:CCAM_VERSION
        if (-not $tag.StartsWith('v')) { $tag = "v$tag" } # tags start with v
        $base = if ($env:CCAM_BASE_URL) { $env:CCAM_BASE_URL } else { "https://github.com/$repo/releases/download" }
        $base = "$base/$tag"
      } else {
        $base = if ($env:CCAM_BASE_URL) { $env:CCAM_BASE_URL } else { "https://github.com/$repo/releases/latest/download" }
      }
      Write-Host "-> downloading $base/$asset"
      try { Invoke-WebRequest -UseBasicParsing "$base/$asset" -OutFile $zip } catch { Fail "download failed: $_" }
      $sums = Join-Path $tmp 'checksums.txt'
      try { Invoke-WebRequest -UseBasicParsing "$base/checksums.txt" -OutFile $sums } catch { Fail 'could not download checksums.txt' }
      $line = Get-Content $sums | Where-Object { ($_ -split '\s+')[1] -eq $asset } | Select-Object -First 1
      if (-not $line) { Fail "$asset is not in checksums.txt" }
      $want = ($line -split '\s+')[0].ToLower()
      $got = (Get-FileHash -Algorithm SHA256 $zip).Hash.ToLower()
      if ($want -ne $got) { Fail 'checksum mismatch' }
    }

    $out = Join-Path $tmp 'out'
    try { Expand-Archive -Path $zip -DestinationPath $out } catch { Fail 'could not extract the archive' }
    $new = Join-Path $out 'ccam.exe'
    if (-not (Test-Path $new)) { Fail 'the archive has no ccam.exe' }
    New-Item -ItemType Directory -Force -Path $dir | Out-Null
    $exe = Join-Path $dir 'ccam.exe'
    $upgrade = Test-Path $exe
    if ($upgrade) {
      # A running ccam.exe cannot be replaced or deleted, only renamed, and
      # a session may still run from the ccam.exe.old of the install before.
      Get-ChildItem -Path $dir -Filter 'ccam.exe.old*' | Where-Object { $_.Name -match '^ccam\.exe\.old(-\d+)?$' } | Remove-Item -Force -ErrorAction SilentlyContinue
      $old = "$exe.old"
      if (Test-Path $old) { $old = "$exe.old-$([DateTime]::UtcNow.Ticks)" }
      Move-Item $exe $old
    }
    Copy-Item $new $exe
    Write-Host "OK $(& $exe version) -> $exe"
    if ($upgrade) {
      # Shortcuts made as hard links (without Developer Mode) still run the
      # old binary; ccam links them to the new one.
      $made = & $exe shims install | Out-String
      if ($made -match '(\d+) new shortcuts' -and $Matches[1] -ne '0') { Write-Host "OK $($Matches[1]) shortcuts now run the new ccam" }
    }
  } finally {
    Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
  }

  # The user Path is a REG_EXPAND_SZ value: it is read and written as one,
  # so entries such as %USERPROFILE%\bin keep their variables.
  $envKey = 'HKCU:\Environment'
  $userPath = (Get-Item $envKey).GetValue('Path', '', 'DoNotExpandEnvironmentNames')
  $onPath = ($env:Path -split ';') + ([Environment]::ExpandEnvironmentVariables($userPath) -split ';') | Where-Object { $_ -and ($_.TrimEnd('\') -ieq $dir.TrimEnd('\')) }
  if (-not $onPath) {
    $newPath = if ($userPath) { "$dir;$userPath" } else { $dir }
    Set-ItemProperty -Path $envKey -Name Path -Value $newPath -Type ExpandString
    # Tells Explorer, and so new terminals, that the environment changed.
    [Environment]::SetEnvironmentVariable('CCAM_PATH_CHANGED', $null, 'User')
    $env:Path = "$dir;$env:Path"
    Write-Host ''
    Write-Host "Added $dir to your user PATH. Open a new terminal to use it everywhere."
  }
  if (-not (Get-Command claude -ErrorAction SilentlyContinue)) {
    Write-Host ''
    Write-Host 'Note: Claude Code (claude) is not installed. Install it first: https://docs.claude.com/claude-code'
  }
  Write-Host ''
  Write-Host 'Next: ccam add  (add your first account)'
}
