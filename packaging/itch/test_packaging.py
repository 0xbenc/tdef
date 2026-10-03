import json
import os
import pathlib
import plistlib
import pty
import shlex
import shutil
import subprocess
import tempfile
import unittest

from publish import commands
from stage import HERE, TARGETS, stage


class ItchPackagingTest(unittest.TestCase):
    def setUp(self):
        self.temporary = tempfile.TemporaryDirectory(prefix="TDEF builds with spaces ")
        self.addCleanup(self.temporary.cleanup)
        self.root = pathlib.Path(self.temporary.name)
        self.dist = self.root / "dist"
        self.dist.mkdir()
        (self.dist / "metadata.json").write_text(json.dumps({"version": "1.0.0"}))
        artifacts = []
        self.sources = {}
        for system, arch in TARGETS:
            binary = self.dist / f"{system}_{arch}" / "tdef"
            binary.parent.mkdir()
            binary.write_text('#!/bin/sh\nprintf "%s\\n" "$@"\n')
            binary.chmod(0o755)
            self.sources[(system, arch)] = binary
            artifacts.append({"type": "Binary", "goos": system, "goarch": arch,
                              "path": str(binary), "extra": {"ID": "tdef"}})
        (self.dist / "artifacts.json").write_text(json.dumps(artifacts))
        self.out = self.dist / "itch"
        stage(self.dist, self.out)

    def test_all_packages_preserve_release_binaries_and_can_launch(self):
        for (system, arch), channel in TARGETS.items():
            package = self.out / channel
            if system == "linux":
                binary = package / "tdef"
                launcher = package / "Play.sh"
            else:
                app = package / "TDEF.app" / "Contents"
                binary = app / "Resources" / "tdef"
                launcher = package / "Play.command"
                info = plistlib.loads((app / "Info.plist").read_bytes())
                self.assertTrue((app / "MacOS" / info["CFBundleExecutable"]).is_file())
                self.assertEqual(info["CFBundleShortVersionString"], "1.0.0")
                result = subprocess.run([str(launcher), "argument with spaces", "literal $value"],
                                        capture_output=True, text=True, check=True)
                self.assertEqual(result.stdout.splitlines(), ["argument with spaces", "literal $value"])
            self.assertEqual(binary.read_bytes(), self.sources[(system, arch)].read_bytes())
            self.assertTrue(os.access(binary, os.X_OK))
            self.assertTrue(os.access(launcher, os.X_OK))
            self.assertTrue((package / ".itch.toml").is_file())

    def mock_path(self, terminal):
        directory = self.root / "mock-bin"
        directory.mkdir(exist_ok=True)
        for name in ["bash", "dirname"]:
            path = directory / name
            if not path.exists():
                path.symlink_to(shutil.which(name))
        command = directory / terminal
        command.write_text('#!/bin/sh\nprintf "%s\\n" "$@"\n')
        command.chmod(0o755)
        return directory

    def test_linux_launcher_opens_terminal_and_preserves_arguments(self):
        for terminal, prefix in [("xdg-terminal-exec", []), ("gnome-terminal", ["--"]),
                                 ("konsole", ["-e"]), ("wezterm", ["start", "--"])]:
            with self.subTest(terminal=terminal):
                directory = self.mock_path(terminal)
                env = dict(os.environ, PATH=str(directory))
                package = self.out / "linux-amd64"
                result = subprocess.run([str(package / "Play.sh"), "two words", "literal $value"],
                                        env=env, capture_output=True, text=True, check=True)
                self.assertEqual(result.stdout.splitlines(), prefix + [str(package / "tdef"), "two words", "literal $value"])
                (directory / terminal).unlink()

    def test_missing_terminal_has_actionable_error(self):
        directory = self.mock_path("temporary-terminal")
        (directory / "temporary-terminal").unlink()
        result = subprocess.run([str(self.out / "linux-amd64" / "Play.sh")],
                                env=dict(os.environ, PATH=str(directory)), capture_output=True, text=True)
        self.assertEqual(result.returncode, 1)
        self.assertIn("./tdef", result.stderr)

    def test_xfce_command_string_preserves_paths_and_arguments(self):
        directory = self.mock_path("xfce4-terminal")
        package = self.out / "linux-amd64"
        result = subprocess.run([str(package / "Play.sh"), "two words", "literal $value"],
                                env=dict(os.environ, PATH=str(directory)), capture_output=True,
                                text=True, check=True)
        arguments = result.stdout.splitlines()
        self.assertEqual(arguments[:2], ["--disable-server", "--command"])
        launched = subprocess.run(shlex.split(arguments[2]), capture_output=True, text=True, check=True)
        self.assertEqual(launched.stdout.splitlines(), ["two words", "literal $value"])

    def test_linux_uses_existing_terminal_without_opening_another(self):
        directory = self.mock_path("xdg-terminal-exec")
        master, slave = pty.openpty()
        try:
            process = subprocess.Popen([str(self.out / "linux-amd64" / "Play.sh"), "from existing terminal"],
                                       env=dict(os.environ, PATH=str(directory)),
                                       stdin=slave, stdout=slave, stderr=slave)
            self.assertEqual(process.wait(timeout=5), 0)
            self.assertEqual(os.read(master, 4096).decode().strip(), "from existing terminal")
        finally:
            os.close(master)
            os.close(slave)

    def test_missing_release_target_prevents_packaging(self):
        artifacts = json.loads((self.dist / "artifacts.json").read_text())
        (self.dist / "artifacts.json").write_text(json.dumps(artifacts[:-1]))
        with self.assertRaisesRegex(ValueError, "missing release targets"):
            stage(self.dist, self.root / "incomplete")
        self.assertFalse((self.root / "incomplete").exists())

    def test_publication_validates_all_channels_first(self):
        plan = commands(self.out, "kairuku-studios/tdef", "v1.0.0", "butler")
        self.assertEqual([command[1] for command in plan], ["validate"] * 4 + ["push"] * 4)
        for command, channel in zip(plan[4:], TARGETS.values()):
            self.assertIn("kairuku-studios/tdef:" + channel, command)
            self.assertEqual(command[-2:], ["--userversion", "1.0.0"])
        with self.assertRaisesRegex(ValueError, "stable"):
            commands(self.out, "kairuku-studios/tdef", "v1.0.0-rc1", "butler")
        (self.out / "mac-arm64" / "VERSION.txt").write_text("2.0.0\n")
        with self.assertRaisesRegex(ValueError, "wrong package version"):
            commands(self.out, "kairuku-studios/tdef", "1.0.0", "butler")


if __name__ == "__main__":
    unittest.main()
