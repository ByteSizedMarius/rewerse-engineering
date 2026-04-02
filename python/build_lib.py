#!/usr/bin/env python3
"""
Build script for compiling the Go shared library.

Builds both Linux (.so) and Windows (.dll) by default. Runs on Windows (requires Go 1.21+ and gcc).
macOS (.dylib) must be built on macOS with --platform darwin.

Usage:
    python build_lib.py                    # build linux + windows
    python build_lib.py --platform darwin  # build macOS (on macOS only)
"""

import argparse
import os
import subprocess
import sys
from pathlib import Path


def get_project_root() -> Path:
    """Get the project root directory."""
    return Path(__file__).parent.parent


def get_output_dir() -> Path:
    """Get the output directory for compiled libraries."""
    out = Path(__file__).parent / "rewerse" / "_lib"
    out.mkdir(parents=True, exist_ok=True)
    return out


def build_linux(project_root: Path, output_dir: Path):
    """Cross-compile shared library for Linux from Windows."""
    output = output_dir / "librewerse.so"
    cmd = [
        "go", "build",
        "-buildmode=c-shared",
        "-o", str(output),
        "./cgo",
    ]
    env = os.environ.copy()
    env["CGO_ENABLED"] = "1"
    env["GOOS"] = "linux"
    env["GOARCH"] = "amd64"

    print(f"Cross-compiling Linux library: {output}")
    subprocess.run(cmd, cwd=project_root, env=env, check=True)
    print(f"Built: {output}")

    # Clean up generated header file
    h_file = output_dir / "librewerse.h"
    if h_file.exists():
        h_file.unlink()
        print(f"Cleaned up: {h_file}")


def build_windows(project_root: Path, output_dir: Path):
    """Build shared library for Windows natively."""
    output = output_dir / "rewerse.dll"
    cmd = [
        "go", "build",
        "-buildmode=c-shared",
        "-o", str(output),
        "./cgo",
    ]
    env = os.environ.copy()
    env["CGO_ENABLED"] = "1"

    print(f"Building Windows library: {output}")
    subprocess.run(cmd, cwd=project_root, env=env, check=True)
    print(f"Built: {output}")

    # Clean up generated header file
    h_file = output_dir / "rewerse.h"
    if h_file.exists():
        h_file.unlink()
        print(f"Cleaned up: {h_file}")


def build_darwin(project_root: Path, output_dir: Path):
    """Build shared library for macOS. Must be run on macOS."""
    if sys.platform != "darwin":
        raise RuntimeError("macOS builds must be run on macOS")

    output = output_dir / "librewerse.dylib"
    cmd = [
        "go", "build",
        "-buildmode=c-shared",
        "-o", str(output),
        "./cgo",
    ]
    env = os.environ.copy()
    env["CGO_ENABLED"] = "1"

    print(f"Building macOS library: {output}")
    subprocess.run(cmd, cwd=project_root, env=env, check=True)
    print(f"Built: {output}")

    # Clean up generated header file
    h_file = output_dir / "librewerse.h"
    if h_file.exists():
        h_file.unlink()
        print(f"Cleaned up: {h_file}")


def main():
    parser = argparse.ArgumentParser(description="Build rewerse shared library")
    parser.add_argument(
        "--platform",
        choices=["darwin"],
        default=None,
        help="Build for a specific platform (only needed for macOS)",
    )
    args = parser.parse_args()

    project_root = get_project_root()
    output_dir = get_output_dir()

    try:
        if args.platform == "darwin":
            build_darwin(project_root, output_dir)
        else:
            build_linux(project_root, output_dir)
            build_windows(project_root, output_dir)
    except subprocess.CalledProcessError as e:
        print(f"Build failed with exit code {e.returncode}", file=sys.stderr)
        sys.exit(1)
    except FileNotFoundError as e:
        print(f"Build failed: {e}", file=sys.stderr)
        print("Make sure Go and gcc are installed.", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
