#!/usr/bin/env python3
"""
Build script for compiling the Go shared library.

Usage:
    python build_lib.py [--platform PLATFORM]

Platforms:
    linux   - Build for Linux (default, requires gcc)
    windows - Cross-compile for Windows (requires mingw-w64)
    darwin  - Cross-compile for macOS (requires osxcross)
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
    """Build shared library for Linux."""
    output = output_dir / "librewerse.so"
    cmd = [
        "go", "build",
        "-buildmode=c-shared",
        "-o", str(output),
        "./cgo",
    ]
    env = os.environ.copy()
    env["CGO_ENABLED"] = "1"

    print(f"Building Linux library: {output}")
    subprocess.run(cmd, cwd=project_root, env=env, check=True)
    print(f"Built: {output}")

    # Clean up generated header file
    h_file = output_dir / "librewerse.h"
    if h_file.exists():
        h_file.unlink()
        print(f"Cleaned up: {h_file}")


def build_windows(project_root: Path, output_dir: Path):
    """Cross-compile shared library for Windows."""
    # Check mingw-w64 is available
    try:
        subprocess.run(
            ["x86_64-w64-mingw32-gcc", "--version"],
            capture_output=True,
            check=True,
        )
    except FileNotFoundError:
        raise FileNotFoundError(
            "mingw-w64 not found. Install with: apt install mingw-w64"
        )

    output = output_dir / "rewerse.dll"
    cmd = [
        "go", "build",
        "-buildmode=c-shared",
        "-o", str(output),
        "./cgo",
    ]
    env = os.environ.copy()
    env["CGO_ENABLED"] = "1"
    env["GOOS"] = "windows"
    env["GOARCH"] = "amd64"
    env["CC"] = "x86_64-w64-mingw32-gcc"

    print(f"Cross-compiling Windows library: {output}")
    subprocess.run(cmd, cwd=project_root, env=env, check=True)
    print(f"Built: {output}")

    # Clean up generated header file
    h_file = output_dir / "rewerse.h"
    if h_file.exists():
        h_file.unlink()
        print(f"Cleaned up: {h_file}")


def build_darwin(project_root: Path, output_dir: Path):
    """Cross-compile shared library for macOS."""
    output = output_dir / "librewerse.dylib"
    cmd = [
        "go", "build",
        "-buildmode=c-shared",
        "-o", str(output),
        "./cgo",
    ]
    env = os.environ.copy()
    env["CGO_ENABLED"] = "1"
    env["GOOS"] = "darwin"
    env["GOARCH"] = "amd64"

    print(f"Cross-compiling macOS library: {output}")
    print("Note: This requires osxcross or building on macOS")
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
        choices=["linux", "windows", "darwin"],
        default="linux",
        help="Target platform (default: linux)",
    )
    args = parser.parse_args()

    project_root = get_project_root()
    output_dir = get_output_dir()

    builders = {
        "linux": build_linux,
        "windows": build_windows,
        "darwin": build_darwin,
    }

    try:
        builders[args.platform](project_root, output_dir)
    except subprocess.CalledProcessError as e:
        print(f"Build failed with exit code {e.returncode}", file=sys.stderr)
        sys.exit(1)
    except FileNotFoundError as e:
        print(f"Build failed: {e}", file=sys.stderr)
        print("Make sure Go and the required C compiler are installed.", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
