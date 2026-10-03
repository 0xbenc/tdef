#!/usr/bin/env python3
"""Package GoReleaser's exact binaries with itch Play launchers and manifests."""
import argparse
import json
import pathlib
import plistlib
import re
import shutil

HERE = pathlib.Path(__file__).resolve().parent
ROOT = HERE.parent.parent
TARGETS = {
    ("linux", "amd64"): "linux-amd64",
    ("linux", "arm64"): "linux-arm64",
    ("darwin", "amd64"): "mac-amd64",
    ("darwin", "arm64"): "mac-arm64",
}


def executable(source: pathlib.Path, destination: pathlib.Path) -> None:
    destination.parent.mkdir(parents=True, exist_ok=True)
    shutil.copyfile(source, destination)
    destination.chmod(0o755)


def stage(dist: pathlib.Path, out: pathlib.Path) -> None:
    metadata = json.loads((dist / "metadata.json").read_text())
    version = metadata["version"].removeprefix("v")
    if not re.fullmatch(r"\d+\.\d+\.\d+(?:-[A-Za-z0-9.-]+)?", version):
        raise ValueError("invalid release version")
    binaries = {}
    for artifact in json.loads((dist / "artifacts.json").read_text()):
        if artifact["type"] != "Binary" or artifact.get("extra", {}).get("ID") != "tdef":
            continue
        target = (artifact["goos"], artifact["goarch"])
        if target not in TARGETS:
            continue
        if target in binaries:
            raise ValueError(f"duplicate binary for {target}")
        path = pathlib.Path(artifact["path"])
        if not path.is_file():
            raise ValueError(f"missing binary: {path}")
        binaries[target] = path
    if set(binaries) != set(TARGETS):
        raise ValueError(f"missing release targets: {set(TARGETS) - set(binaries)}")
    for target, channel in TARGETS.items():
        directory = out / channel
        if directory.exists():
            shutil.rmtree(directory)
        directory.mkdir(parents=True)
        for name in ["LICENSE", "README.md", "LORE.md", "TACTICS.md"]:
            shutil.copyfile(ROOT / name, directory / name)
        shutil.copyfile(HERE / "README-PLAY.txt", directory / "README-PLAY.txt")
        (directory / "VERSION.txt").write_text(version + "\n")
        if target[0] == "linux":
            executable(binaries[target], directory / "tdef")
            executable(HERE / "Play.sh", directory / "Play.sh")
            launch_path, platform = "Play.sh", "linux"
        else:
            app = directory / "TDEF.app" / "Contents"
            executable(binaries[target], app / "Resources" / "tdef")
            executable(HERE / "app-launch.sh", app / "MacOS" / "Play")
            executable(HERE / "app-play.command", app / "Resources" / "Play.command")
            executable(HERE / "Play.command", directory / "Play.command")
            info = {
                "CFBundleExecutable": "Play",
                "CFBundleIdentifier": "io.itch.kairuku-studios.tdef",
                "CFBundleName": "TDEF",
                "CFBundleDisplayName": "TDEF",
                "CFBundlePackageType": "APPL",
                "CFBundleShortVersionString": version.split("-")[0],
                "CFBundleVersion": version.split("-")[0],
                "LSMinimumSystemVersion": "12.0",
                "NSHighResolutionCapable": True,
            }
            (app / "Info.plist").write_bytes(plistlib.dumps(info))
            launch_path, platform = "TDEF.app", "osx"
        (directory / ".itch.toml").write_text(
            f'[[actions]]\nname = "play"\npath = "{launch_path}"\nplatform = "{platform}"\n'
        )
        print(f"staged {channel}: {version}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--dist", type=pathlib.Path, default=pathlib.Path("dist"))
    parser.add_argument("--out", type=pathlib.Path, default=pathlib.Path("dist/itch"))
    args = parser.parse_args()
    stage(args.dist, args.out)
