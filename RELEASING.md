# Releasing TERMTD

Each stable tag produces one set of release binaries and publishes through two
routes: itch.io with Butler and terminal-opening Play launchers, and the
Homebrew formula path used by `0xbenc/uuid`.

1. Push a stable `vX.Y.Z` tag.
2. GitHub Actions vets and tests on Linux, macOS, and Windows.
3. GoReleaser builds Linux/macOS/Windows binaries for amd64/arm64, archives the binary
   with the documentation, and publishes checksums and SBOMs to GitHub Releases.
4. The workflow builds itch packages from those exact binaries and attests the
   GitHub release archives.
5. Independent jobs publish the packages to `kairuku-studios/termtd` with Butler
   and render/push `Formula/termtd.rb` to `0xbenc/homebrew-tap` from the archive
   checksums. A failure in one publishing job does not prevent the other job
   from running.

Install a published stable release with `brew install 0xbenc/tap/termtd`.
Prerelease tags such as `v1.0.0-rc1` publish GitHub release assets but leave both
the tap and itch's stable channels on their last stable versions.

## itch packages

The itch destination is <https://kairuku-studios.itch.io/termtd>. Each build has an
explicit `.itch.toml` Play action and instructions in `README-PLAY.txt`:

| Channel | Player's launcher |
| --- | --- |
| `linux-amd64` | `Play.sh` → desktop terminal, or current terminal |
| `linux-arm64` | `Play.sh` → desktop terminal, or current terminal |
| `mac-amd64` | `TERMTD.app` → Terminal; `Play.command` also works |
| `mac-arm64` | `TERMTD.app` → Terminal; `Play.command` also works |
| `windows-amd64` | `termtd.exe` with itch's `console = true`; `Play.cmd` for browser downloads |
| `windows-arm64` | `termtd.exe` with itch's `console = true`; `Play.cmd` for browser downloads |

Launchers do not install anything or change PATH. Homebrew provides the global
`termtd` command. The macOS bundle requires macOS 12+, matching
[Go 1.26's platform support](https://go.dev/doc/go1.26#darwin).

Butler is pinned to 15.31.0 with an archive SHA-256 in
`packaging/itch/install-butler.sh`. The builder uses GoReleaser's
`dist/artifacts.json` to locate the binaries and restores executable permissions
when rebuilding packages from downloaded CI artifacts.

Windows packages require Windows 10/11. The executable uses native console
keyboard, mouse, and resize events and enables UTF-8 virtual terminal output.
The itch manifest opens a console directly; Windows Terminal is optional.
GitHub Windows archives and portable itch downloads use ZIP format.

Itch publishing validates all six packages before uploading the first one,
then pushes with `--userversion X.Y.Z`. A retry replaces the same channels;
publishing across all six channels is sequential, not an atomic transaction.

## One-time GitHub setup

The TERMTD repository needs an Actions secret named `TAP_GITHUB_TOKEN` with
Contents write access to `0xbenc/homebrew-tap`. This is the same secret name
used by the other projects. TERMTD's secret has been configured using the existing
local `gh` login, which has tap write access and TERMTD admin access.

For another repository, check the existing login before requesting a token:

```sh
gh auth status
gh api repos/0xbenc/homebrew-tap --jq '.permissions.push'
gh api repos/0xbenc/termtd --jq '.permissions.admin'
```

When that login has the required access, configure the secret directly without
printing or saving its value:

```sh
set -o pipefail
gh auth token --hostname github.com | gh secret set TAP_GITHUB_TOKEN --repo 0xbenc/termtd
```

GitHub does not return existing repository secret values. Reusing the local
authenticated credential avoids needing to retrieve another repository's secret.
Alternatively, supply a dedicated tap token through the private prompt:

```sh
gh secret set TAP_GITHUB_TOKEN --repo 0xbenc/termtd
```

Stable releases fail before publication when this secret is absent.
The regular `GITHUB_TOKEN` is supplied
by Actions for publishing the TERMTD release itself.

The following itch settings are configured on `0xbenc/termtd`:

- Actions variable `ITCH_TARGET`: `kairuku-studios/termtd`.
- Actions secret `BUTLER_API_KEY`: copied directly from the local Butler login.

If the itch credential needs renewal, run `butler login` locally, then transfer
its saved value without printing it:

```sh
gh secret set BUTLER_API_KEY --repo 0xbenc/termtd < ~/.config/itch/butler_creds
```

See [Butler authentication](https://itch.io/docs/butler/login.html) for the
credential location on other OSes. Both publisher credentials and `ITCH_TARGET`
are checked before a stable GitHub release is published.

## 1.1.0: Windows support

This release adds Windows 10/11 downloads for Intel/AMD and ARM, native console
keyboard/mouse/resize support, and itch Play actions that open a console window.
Browser downloads include Play.cmd, which keeps startup errors visible. CI now
tests Windows alongside Linux and macOS; all six release targets are packaged.

## Prepare 1.1.0

Review and commit the intended game changes and release configuration first.
Generated images and release artifacts under `output/` and `dist/` are ignored.

```sh
go vet ./...
go test ./...
python3 -m unittest discover -s packaging/homebrew -p 'test_*.py'
python3 -m unittest discover -s packaging/itch -p 'test_*.py'
goreleaser check
goreleaser release --snapshot --clean --skip=sbom
python3 packaging/itch/stage.py
butler validate --platform linux --arch amd64 dist/itch/linux-amd64
butler validate --platform osx --arch amd64 dist/itch/mac-amd64
```

The snapshot command builds all six targets and archives without publishing.
Remove `--skip=sbom` if `syft` is installed to also exercise SBOM generation.
Release CI installs syft before publishing. Local binaries report `termtd dev`;
GoReleaser embeds the version from the tag.

Before tagging, confirm that Linux, macOS, and Windows CI passed and smoke-test
playing, resizing, mouse input, exiting, and continuing a saved campaign on all
three OSes.
Cross-compilation checks builds, while a real terminal checks platform behavior.
Also test the macOS app and `.command` launcher and the Linux graphical Play
launcher from the itch app. Linux launcher tests cover terminal selection and
paths/arguments with spaces; a graphical macOS launch requires a Mac.
Butler's validator currently accepts only `386` and `amd64`, so ARM package
validation checks the script/app launch target using `--arch amd64`.

Windows ARM packages contain a native ARM executable. The pinned Butler validates
their launch manifest but warns about the ARM PE header and reports amd64, so its
output does not verify their executable architecture. Smoke-test on ARM Windows
before release.
On Windows, also test itch Play and Play.cmd from an extracted folder with spaces,
both in Windows Terminal and the classic console. Verify that normal exit restores
the console settings and startup errors remain visible through Play.cmd.

To preview the stable channel commands using staged stable-version packages:

```sh
python3 packaging/itch/publish.py --target kairuku-studios/termtd --version 1.1.0 --dry-run
```

Remove `--dry-run` only when intentionally publishing those staged packages.

When ready to publish the committed revision:

```sh
git tag -a v1.1.0 -m 'TERMTD 1.1.0: Windows support'
git push origin main
git push origin v1.1.0
```

After the release workflow succeeds:

```sh
brew update
brew install 0xbenc/tap/termtd
termtd --version
brew test 0xbenc/tap/termtd
butler status kairuku-studios/termtd
```

If the GitHub release succeeds but the tap push fails, use its existing
`checksums.txt` with `packaging/homebrew/render.py` to regenerate the formula;
do not rebuild the published archives with different checksums. The publisher
leaves the tap unchanged if that formula is already current.

## Rename rollout

Before the first termtd release, rename the itch.io project slug to `termtd`
and update the repository Actions variable `ITCH_TARGET` to
`kairuku-studios/termtd`. Repository secrets retain their existing names.
The release workflow publishes `Formula/termtd.rb` into `0xbenc/homebrew-tap`;
retire the old `Formula/tdef.rb` there once the new formula is available.
Existing release archives retain their original names. New tags produce
`termtd` archives, executables, and `TERMTD.app`.

Player data now uses `~/.termtd-hiscores.json`, `~/.termtd-journal.json`, and
`~/.termtd-lair.json`. When a new file is absent, the game reads its old
`.tdef-` counterpart. Subsequent saves use the new name and leave the old
file intact. Existing new files always take precedence.