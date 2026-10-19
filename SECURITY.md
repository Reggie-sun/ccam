# Security policy

## Reporting a vulnerability

Please report security problems privately through GitHub: open the [Security tab](https://github.com/jhyunwoo/ccam/security) of this repository and choose **Report a vulnerability**. Do not open a public issue.

Include the ccam version (`ccam --version`), your operating system, and the steps that show the problem.

## Supported versions

Fixes go into the latest release. Run `ccam update` to get it.

## What ccam handles

ccam never copies or shares a login between accounts; each account's `.credentials.json` stays in its own folder. It reads an account's access token only to ask Anthropic's usage endpoint for that account's usage. The install script and `ccam update` check every download against the SHA-256 checksums published with the release.
