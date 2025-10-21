# Marina Build Documentation

## Overview

Marina is a Go application with CGO dependencies for GUI rendering (Dear ImGui) and SQLite database access. This document provides platform-specific build instructions and addresses common dependency issues.

## Prerequisites

### All Platforms

- **Go**: Version specified in `go.mod` (currently 1.24.0)
  - Check version: `go version`
  - Download: https://go.dev/dl/

- **Git**: For cloning and version control
  - Check version: `git --version`

- **CGO Support**: Marina requires CGO for native dependencies
  - CGO is enabled by default in standard Go distributions
  - Verify: `go env CGO_ENABLED` should return `1`

### Dependency Overview

Marina uses several dependencies that require native compilation:

1. **cimgui-go / giu**: Dear ImGui bindings (requires C++ compiler and graphics libraries)
2. **go-sqlite3**: SQLite driver (requires C compiler)
3. **Platform-specific windowing**: X11 (Linux), Win32 (Windows), Cocoa (macOS)

## Linux

### System Requirements

- **OS**: Any modern Linux distribution (Ubuntu 22.04+ recommended)
- **Architecture**: x86_64 (amd64)

### Dependencies

Install development libraries:

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y \
    build-essential \
    libx11-dev \
    libxcursor-dev \
    libxrandr-dev \
    libxinerama-dev \
    libxi-dev \
    libglx-dev \
    libgl1-mesa-dev \
    libxxf86vm-dev \
    libgtk-3-dev

# Fedora/RHEL
sudo dnf install -y \
    gcc \
    gcc-c++ \
    libX11-devel \
    libXcursor-devel \
    libXrandr-devel \
    libXinerama-devel \
    libXi-devel \
    libGL-devel \
    libXxf86vm-devel \
    gtk3-devel

# Arch Linux
sudo pacman -S --needed \
    base-devel \
    libx11 \
    libxcursor \
    libxrandr \
    libxinerama \
    libxi \
    mesa \
    libxxf86vm \
    gtk3
```

**Why these dependencies?**:
- **libx11-dev**: X Window System client library (core windowing)
- **libxcursor-dev, libxrandr-dev, libxinerama-dev, libxi-dev**: X11 extensions for cursor, multi-monitor, and input
- **libglx-dev, libgl1-mesa-dev**: OpenGL rendering (required by Dear ImGui)
- **libxxf86vm-dev**: XFree86 video mode extension
- **libgtk-3-dev**: GTK3 for native file dialogs (via `sqweek/dialog`)

### Build Instructions

```bash
# Clone the repository
git clone https://github.com/Cardboy777/marina.git
cd marina

# Download Go dependencies
go get

# Build the executable
go build -o marina

# Run
./marina
```

**Build Flags**:
```bash
# Release build with version
go build -ldflags="-s -w -X 'marina/constants.AppVersion=1.0.0'" -o marina

# -s: Omit symbol table
# -w: Omit DWARF debugging information
# -X: Set variable value at link time (version string)
```

### AppImage Packaging

Marina uses AppImage for portable Linux distribution.

**Dependencies**:
```bash
# Install AppImage builder
sudo apt-get install -y python3 python3-pip squashfs-tools
pip3 install appimage-builder
```

**Build AppImage**:
```bash
# Prepare directory structure
mkdir -p build/linux/AppDir/usr/bin
mkdir -p build/linux/AppDir/usr/share/icons/hicolor/48x48/apps
mkdir -p build/linux/AppDir/usr/share/icons/hicolor/scalable/apps

# Copy assets
cp build/assets/marina.png build/linux/AppDir/usr/share/icons/hicolor/48x48/apps/
cp build/assets/marina.ico build/linux/AppDir/usr/share/icons/hicolor/scalable/apps/

# Build binary
go build -ldflags="-X 'marina/constants.AppVersion=1.0.0'" -o build/linux/AppDir/usr/bin/marina

# Create AppImage
cd build/linux
appimage-builder --recipe AppImageBuilder.yml
```

**AppImage Configuration**: See `build/linux/AppImageBuilder.yml` for bundled dependencies.

### Common Issues

#### Issue: `fatal error: X11/Xlib.h: No such file or directory`
**Cause**: Missing X11 development headers
**Solution**: Install `libx11-dev` (or equivalent for your distro)

#### Issue: `undefined reference to 'glXCreateContext'`
**Cause**: Missing OpenGL libraries
**Solution**: Install `libgl1-mesa-dev` and `libglx-dev`

#### Issue: `cannot find -lGL`
**Cause**: OpenGL library not found by linker
**Solution**: Install Mesa development packages or check `/usr/lib/x86_64-linux-gnu/` for `libGL.so`

#### Issue: AppImage crashes on startup
**Cause**: Missing runtime libraries not bundled in AppImage
**Solution**: Check `AppImageBuilder.yml` apt includes section; may need to add missing libraries

## Windows

### System Requirements

- **OS**: Windows 10/11 (64-bit)
- **Architecture**: x86_64 (AMD64)

### Dependencies

Windows builds require **MSYS2** to provide a Unix-like build environment with MinGW-w64 toolchain.

**Install MSYS2**:
1. Download from https://www.msys2.org/
2. Run installer and follow instructions
3. Open "MSYS2 MSYS" from Start Menu

**Install Build Tools**:
```bash
# Update package database
pacman -Syu

# Install MinGW-w64 toolchain and tools
pacman -S --noconfirm \
    git \
    zip \
    mingw-w64-x86_64-toolchain
```

**Install Go**:
- Download Windows installer from https://go.dev/dl/
- Install normally
- MSYS2 will inherit the PATH from Windows

### Build Instructions

**From MSYS2 Shell**:
```bash
# Clone repository
git clone https://github.com/Cardboy777/marina.git
cd marina

# Download dependencies
go get

# Build
go build -o marina.exe

# Release build with Windows subsystem (hides console)
go build \
    -o marina.exe \
    -ldflags="-s -w -H=windowsgui -extldflags='-static' -X 'marina/constants.AppVersion=1.0.0'"
```

**Build Flags Explained**:
- `-ldflags="-s -w"`: Strip debug info (smaller binary)
- `-H=windowsgui`: Use Windows GUI subsystem (no console window)
- `-extldflags='-static'`: Static linking (reduces DLL dependencies)
- `-X 'marina/constants.AppVersion=...'`: Embed version string

### Version Info Embedding (Optional)

Marina supports embedding Windows version information (appears in file properties).

**Setup**:
```bash
# Install goversioninfo
go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest

# Generate version resource
go generate
```

**Configuration**: Edit `build/windows/versioninfo.json` to customize version info.

**Note**: The `//go:generate` directive in `marina.go` automates this, but the GitHub workflow currently has it commented out.

### Common Issues

#### Issue: `gcc: command not found`
**Cause**: MinGW toolchain not installed or not in PATH
**Solution**: Run `pacman -S mingw-w64-x86_64-toolchain` in MSYS2

#### Issue: `undefined reference to 'WinMain'`
**Cause**: Missing `-H=windowsgui` flag
**Solution**: Add `-ldflags="-H=windowsgui"` to build command

#### Issue: Application crashes on startup on other Windows machines
**Cause**: Missing MSYS2 DLLs
**Solution**: Use `-extldflags='-static'` to statically link C runtime

#### Issue: `Package 'x11' not found`
**Cause**: Trying to build Linux code on Windows
**Solution**: Ensure you're using MSYS2 environment, not WSL or native CMD

#### Issue: Binary size is very large (>100MB)
**Cause**: Debug symbols included
**Solution**: Use `-ldflags="-s -w"` to strip symbols (reduces to ~20-30MB)

## macOS

### System Requirements

- **OS**: macOS 10.15 Catalina or later
- **Architecture**: x86_64 or ARM64 (Apple Silicon)

### Dependencies

**Xcode Command Line Tools**:
```bash
xcode-select --install
```

This provides:
- C/C++ compiler (clang)
- macOS SDK
- Headers for Cocoa framework

**Go Installation**:
```bash
# Using Homebrew
brew install go

# Or download from https://go.dev/dl/
```

### Build Instructions

```bash
# Clone repository
git clone https://github.com/Cardboy777/marina.git
cd marina

# Download dependencies
go get

# Build
go build -o marina

# Release build
go build -ldflags="-s -w -X 'marina/constants.AppVersion=1.0.0'" -o marina
```

### App Bundle Creation (Recommended)

macOS applications should be packaged as `.app` bundles for proper integration.

**Manual Bundle Creation**:
```bash
# Create bundle structure
mkdir -p Marina.app/Contents/MacOS
mkdir -p Marina.app/Contents/Resources

# Copy binary
cp marina Marina.app/Contents/MacOS/

# Create Info.plist
cat > Marina.app/Contents/Info.plist << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>marina</string>
    <key>CFBundleIconFile</key>
    <string>marina.icns</string>
    <key>CFBundleIdentifier</key>
    <string>com.cardboy777.marina</string>
    <key>CFBundleName</key>
    <string>Marina</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>NSHighResolutionCapable</key>
    <true/>
</dict>
</plist>
EOF

# Convert icon (requires iconutil)
# Place marina.icns in Marina.app/Contents/Resources/
```

**Launch Bundle**:
```bash
open Marina.app
```

### DMG Creation

For distribution, package the `.app` bundle in a DMG disk image:

```bash
# Create DMG
hdiutil create -volname "Marina" -srcfolder Marina.app -ov -format UDZO Marina.dmg
```

### Common Issues

#### Issue: `ld: library not found for -lc++`
**Cause**: Xcode Command Line Tools not installed
**Solution**: Run `xcode-select --install`

#### Issue: `xcrun: error: invalid active developer path`
**Cause**: Xcode path not set correctly
**Solution**: Run `sudo xcode-select --reset`

#### Issue: Application not opening (bounces in dock and quits)
**Cause**: Not packaged as proper `.app` bundle
**Solution**: Follow app bundle creation steps above

#### Issue: "Marina" cannot be opened because the developer cannot be verified
**Cause**: macOS Gatekeeper security (unsigned application)
**Solution**: Right-click → Open, or run `xattr -cr Marina.app` to remove quarantine attribute

#### Issue: Build fails with "implicit declaration of function" warnings
**Cause**: Missing macOS SDK headers
**Solution**: Update Xcode Command Line Tools: `softwareupdate --install -a`

## Cross-Platform Build Considerations

### CGO and Cross-Compilation

**Important**: CGO makes cross-compilation difficult because it requires target platform's C compiler and libraries.

**Options**:
1. **Native Builds**: Build on each target platform (recommended)
2. **Docker**: Use platform-specific Docker containers
3. **Cross-Compiler Toolchains**: Set up mingw-w64 (Windows on Linux) or osxcross (macOS on Linux)

### CI/CD (GitHub Actions)

Marina uses GitHub Actions for automated builds (see `.github/workflows/generate-builds.yml`).

**Build Matrix**:
- **Linux**: `ubuntu-22.04` runner, AppImage output
- **Windows**: `windows-latest` runner with MSYS2, `.exe` output
- **macOS**: `macos-14` runner, currently outputs `.dmg` (should be `.app`)

**Artifacts**: Built binaries are uploaded as GitHub Actions artifacts.

## Dependency Deep Dive

### cimgui-go (Dear ImGui Bindings)

**What it is**: CGO wrapper around Dear ImGui, a C++ immediate-mode GUI library.

**Build Requirements**:
- C++ compiler (g++, clang++, MSVC)
- OpenGL libraries
- Platform windowing libraries

**Compilation Process**:
1. Go compiler invokes CGO
2. CGO compiles `cimgui` C++ wrapper
3. Linker combines Go code with C++ objects
4. Links against system OpenGL and windowing libraries

**Platform Specifics**:
- **Linux**: Links against X11, OpenGL
- **Windows**: Links against OpenGL32, GDI32, User32
- **macOS**: Links against Cocoa framework, OpenGL framework

**Build Time**: First build is slow (~2-5 minutes) because CGO compiles the entire Dear ImGui library.

### go-sqlite3

**What it is**: SQLite database driver with CGO bindings.

**Build Requirements**:
- C compiler (gcc, clang, MSVC)

**Compilation Process**:
1. Embeds full SQLite3 C source (~230KB)
2. CGO compiles SQLite3 amalgamation
3. Links into Go binary

**Benefits of CGO Version**:
- No external SQLite dependencies
- Full SQLite3 feature set
- Thread-safe with default flags

**Alternative**: `modernc.org/sqlite` is a pure-Go SQLite implementation (no CGO), but `go-sqlite3` is more mature and faster.

### Platform Dialogs (sqweek/dialog)

**What it is**: Native file picker dialogs for each platform.

**Platform Implementations**:
- **Windows**: Uses Win32 COM interfaces
- **Linux**: Spawns GTK3 dialog via `gtk-3.0` pkg-config
- **macOS**: Uses Cocoa NSOpenPanel

**Why GTK3 on Linux?**:
- `zenity` or `kdialog` could be used, but require external binaries
- GTK3 provides native look and feel
- Most Linux distros have GTK3 installed

## Build Optimization

### Reducing Binary Size

Default binaries are large (~50-100MB) due to:
- Dear ImGui C++ code
- SQLite3 amalgamation
- Go runtime
- Debug symbols

**Optimization Flags**:
```bash
go build -ldflags="-s -w" -o marina
# -s: Strip symbol table (~20% reduction)
# -w: Strip DWARF debug info (~30% reduction)
```

**Further Optimization**:
```bash
# Use UPX compression (not recommended for distribution)
upx --best --lzma marina
```

**Typical Sizes**:
- **With debug info**: ~100MB
- **Stripped (-s -w)**: ~30MB
- **UPX compressed**: ~10MB (may trigger antivirus)

### Build Speed Optimization

**First Build**: Slow due to CGO compilation of Dear ImGui and SQLite3.

**Subsequent Builds**: Go caches CGO objects, so rebuilds are fast.

**Tips**:
- Use `go build -i` to install dependencies (deprecated in Go 1.18+, now default behavior)
- Consider using `ccache` for C/C++ compilation caching
- Use `-p` flag to parallelize compilation: `go build -p 8`

## Version Injection

Marina injects version information at build time using linker flags:

```bash
go build -ldflags="-X 'marina/constants.AppVersion=1.0.0'"
```

**Where it's used**:
- `constants/constants.go`: `var AppVersion = "development"`
- Displayed in UI/logs
- GitHub workflow injects `github.runId` as version

**Best Practice**: Use semantic versioning (e.g., `v1.2.3`) or git commit hash.

## Troubleshooting

### General CGO Issues

#### Issue: `undefined reference` errors during linking
**Diagnosis**: Check what symbols are missing with `nm` or `objdump`
**Solution**: Install missing system libraries or update `CGO_LDFLAGS`

#### Issue: `cgo: C compiler not found`
**Solution**:
- Linux: Install `build-essential` or `gcc`
- Windows: Install MinGW via MSYS2
- macOS: Install Xcode Command Line Tools

### Go Module Issues

#### Issue: `go: finding module for package` errors
**Solution**: Run `go mod tidy` to clean up dependencies

#### Issue: `verifying module` checksum errors
**Solution**: Run `go clean -modcache` and `go get` again

### Platform Detection Issues

Marina uses `runtime.GOOS` for platform detection. Verify with:
```bash
go env GOOS GOARCH
```

Expected values:
- Linux: `GOOS=linux GOARCH=amd64`
- Windows: `GOOS=windows GOARCH=amd64`
- macOS: `GOOS=darwin GOARCH=amd64` or `GOARCH=arm64`

## Development Builds

For development, disable optimization and enable race detection:

```bash
# Development build with race detector (slower, safer)
go build -race -o marina-dev

# Run with verbose logging
./marina-dev -v
```

**Note**: Race detector is not supported with CGO on Windows.

## Recommended Development Workflow

1. **Use Native Development**: Develop on your target platform
2. **Test on All Platforms**: Use VMs or CI/CD for cross-platform testing
3. **Version Control CGO Output**: `.gitignore` should exclude compiled artifacts but keep source
4. **Cache Dependencies**: Go modules cache speeds up builds

## Additional Resources

- **Go CGO Documentation**: https://pkg.go.dev/cmd/cgo
- **Dear ImGui**: https://github.com/ocornut/imgui
- **giu Documentation**: https://github.com/AllenDang/giu
- **SQLite3 Driver**: https://github.com/mattn/go-sqlite3
- **MSYS2**: https://www.msys2.org/docs/what-is-msys2/
