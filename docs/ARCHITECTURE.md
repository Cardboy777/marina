# Marina Architecture Documentation

## Overview

Marina is a cross-platform desktop application written in Go that serves as an unofficial launcher and version manager for Harbour Masters game ports (Ship of Harkinian, 2 Ship 2 Harkinian, and Starship). The application provides a GUI for downloading, managing, and launching different versions of these game ports.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Marina Application                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌────────────────── Presentation Layer ──────────────────────┐ │
│  │                                                              │ │
│  │  ui/                                                         │ │
│  │  ├── mainwindow/  (Main UI components)                      │ │
│  │  ├── settings/    (Settings dialog)                         │ │
│  │  ├── dialogs/     (Confirmation & error dialogs)            │ │
│  │  ├── fonts/       (Font management)                         │ │
│  │  └── importconfigs/ (Config import UI)                      │ │
│  │                                                              │ │
│  │  Technology: Dear ImGui (via AllenDang/giu)                 │ │
│  └──────────────────────────────────────────────────────────────┘ │
│                              ↕                                    │
│  ┌─────────────── Business Logic Layer ────────────────────────┐ │
│  │                                                              │ │
│  │  services/     (GitHub API integration)                     │ │
│  │  ├── Fetch releases from GitHub                             │ │
│  │  ├── Rate limiting (1-hour cache)                           │ │
│  │  └── Parse release assets                                   │ │
│  │                                                              │ │
│  │  stores/       (Data access & state management)             │ │
│  │  ├── Version management                                     │ │
│  │  ├── ROM management                                         │ │
│  │  └── Unstable version tracking                              │ │
│  │                                                              │ │
│  │  launcher/     (Game execution)                             │ │
│  │  ├── ROM linking                                            │ │
│  │  ├── Executable detection                                   │ │
│  │  └── Process spawning                                       │ │
│  │                                                              │ │
│  │  files/        (File system operations)                     │ │
│  │  ├── Download & extraction                                  │ │
│  │  ├── ROM validation (SHA1)                                  │ │
│  │  └── Import/export configs                                  │ │
│  └──────────────────────────────────────────────────────────────┘ │
│                              ↕                                    │
│  ┌───────────────── Data Layer ───────────────────────────────┐ │
│  │                                                              │ │
│  │  db/           (SQLite persistence)                         │ │
│  │  ├── Version tracking                                       │ │
│  │  ├── ROM registry                                           │ │
│  │  ├── Installation state                                     │ │
│  │  └── Fetch timestamps                                       │ │
│  │                                                              │ │
│  │  settings/     (Configuration management)                   │ │
│  │  └── Uses Viper for TOML config                            │ │
│  └──────────────────────────────────────────────────────────────┘ │
│                              ↕                                    │
│  ┌──────────────── External Services ──────────────────────────┐ │
│  │                                                              │ │
│  │  • GitHub REST API (releases & commits)                     │ │
│  │  • Nightly.link (unstable builds)                           │ │
│  │  • File System (game installations, ROMs, configs)          │ │
│  └──────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## Component Descriptions

### 1. Presentation Layer (ui/)

**Purpose**: Provides the graphical user interface using Dear ImGui through the giu wrapper.

**Key Files**:
- `ui/start.go`: Application entry point, initializes master window
- `ui/mainwindow/mainwindow.go`: Main layout with game selector and version list
- `ui/mainwindow/versions.go`: Version list rendering and interactions
- `ui/mainwindow/gameselector.go`: Game selection dropdown
- `ui/mainwindow/roms.go`: ROM status display
- `ui/mainwindow/toolbar.go`: Refresh and settings toolbar

**Technology Stack**:
- Dear ImGui (Immediate Mode GUI)
- AllenDang/giu (Go bindings)
- AllenDang/cimgui-go (C bindings layer)

**UI Structure**:
```
┌─────────────────────────────────────────┐
│  Marina                            [⚙]  │
├───────────┬─────────────────────────────┤
│  Game     │  Toolbar: [Refresh] [⚙]    │
│  Selector │                             │
│  ┌─────┐  │  Version List:              │
│  │ SoH │  │  ┌──────────────────────┐   │
│  │2Ship│  │  │ Unstable - [Install] │   │
│  │Star │  │  ├──────────────────────┤   │
│  └─────┘  │  │ v1.0.0     [⋮] [Play]│   │
│           │  │ 2025-01-15           │   │
│  ROM      │  ├──────────────────────┤   │
│  Status   │  │ v0.9.0     [Install] │   │
│           │  └──────────────────────┘   │
└───────────┴─────────────────────────────┘
```

### 2. Business Logic Layer

#### services/ - External API Integration
**Purpose**: Manages communication with GitHub API and external services.

**Key Operations**:
- `SyncReleases()`: Fetches releases and unstable versions, respects 1-hour cache
- `fetchReleaseVersions()`: Paginated GitHub release fetching
- `fetchLatestCommit()`: Gets latest commit for unstable builds

**Rate Limiting Strategy**:
- Caches version data in database
- Only refreshes if > 1 hour since last fetch OR manual refresh
- Handles GitHub rate limit errors gracefully

#### stores/ - State Management
**Purpose**: Provides clean interface between business logic and database.

**Key Patterns**:
- `GetVersions()`: Retrieves cached versions for a repository
- `AddVersion()`: Adds new version (deduplicates by tag)
- `SetVersionInstalled()`: Updates installation state
- Similar patterns for ROMs and unstable versions

#### launcher/ - Game Execution
**Purpose**: Handles launching installed games.

**Process Flow**:
1. Copy/link ROMs to version directory
2. Detect executable (platform-specific extensions)
3. Spawn process with `SHIP_HOME` environment variable
4. Run in goroutine with completion callback

**Platform Detection**:
- Linux: `.appimage` files
- macOS: `.dmg` files (note: this may need refinement for actual .app bundles)
- Windows: `.exe` files

#### files/ - File System Management
**Purpose**: Handles all file operations, downloads, and ROM validation.

**Key Operations**:
- **Download Management**: HTTP download → ZIP extraction → installation tracking
- **ROM Validation**: SHA1 hashing against known-good ROM hashes
- **Import/Export**: Copy configuration, saves, mods between versions
- **Directory Structure**:
  ```
  <InstallDir>/
  ├── manifest.db              # SQLite database
  ├── roms/                    # Validated ROM files
  │   └── pal-1_0.z64
  └── versions/                # Game installations
      ├── shipwright/
      │   ├── v1_0_0/
      │   └── unstable-2025-01-15/
      ├── 2ship2harkinian/
      └── starship/
  ```

### 3. Data Layer

#### db/ - Database Management
**Purpose**: SQLite persistence layer for application state.

**Schema** (see DATABASE.md for details):
- `Releases`: Stable version tracking
- `UnstableVersions`: Development build tracking
- `InstalledRoms`: ROM validation cache
- `LastFetchedDate`: GitHub API rate limiting

**Connection Management**:
- Single global connection via `db.Init()`
- `ResetDbConnection()` called on install directory change
- Automatic schema initialization on startup

#### settings/ - Configuration
**Purpose**: Manages application settings via TOML file.

**Configuration Location**:
- Linux: `~/.config/marina/config.toml`
- Windows: `%APPDATA%/marina/config.toml`
- macOS: `~/Library/Application Support/marina/config.toml`

**Settings**:
- `InstallDir`: Where games and ROMs are stored (default: `~/.local/share/marina`)

### 4. Static Data (constants/ & types/)

#### constants/
**Purpose**: Application-wide constants and game definitions.

**Key Files**:
- `constants/constants.go`: App name, version, file permissions
- `constants/games/games.go`: Array of all supported games
- `constants/games/shipofharkinian.go`: SoH definition with ROM hashes, download URLs
- `constants/games/2ship2harkinian.go`: 2Ship definition
- `constants/games/starship.go`: Starship definition

**Game Definition Structure**:
```go
Repository{
    Id:                int
    Name:              string
    Owner:             string  // GitHub owner
    Repository:        string  // GitHub repo name
    LatestBuildUrls:   DownloadUrls  // Nightly build URLs
    AcceptedRomHashes: []Rom   // Valid ROM SHA1 hashes
    Imports:           Imports // What can be imported between versions
}
```

#### types/
**Purpose**: Core data structures used throughout the application.

**Key Types**:
- `Version`: Represents a stable release
- `UnstableVersion`: Represents a nightly/development build
- `Rom`: ROM identification (name + SHA1)
- `Repository`: Game port definition
- `DownloadUrls`: Platform-specific download URLs

## Data Flow

### Version Refresh Flow
```
User clicks Refresh
    ↓
ui/mainwindow/toolbar.go → RefreshVersions(force=true)
    ↓
services/services.go → SyncReleases()
    ↓
services/github.go → fetchReleaseVersions()
    ↓                → fetchLatestCommit()
    ↓
stores/versions.go → AddVersion() (deduplicated)
    ↓
db/versions.go → AddRelease() (SQL INSERT)
    ↓
stores/lastfetched.go → UpdateLastFetched()
    ↓
ui updates automatically (giu reactive model)
```

### Download & Install Flow
```
User clicks Install
    ↓
ui/mainwindow/versions.go → install()
    ↓
files/files.go → DownloadVersion()
    ↓
    ├── HTTP GET download URL
    ├── Save to temp ZIP file
    ├── Extract ZIP to version directory
    └── Delete ZIP
    ↓
stores/versions.go → SetVersionInstalled(true)
    ↓
db/versions.go → SetInstalled()
    ↓
ui updates (shows Play button)
```

### Launch Game Flow
```
User clicks Play
    ↓
ui/mainwindow/versions.go → play()
    ↓
launcher/launcher.go → LaunchGame()
    ↓
files/files.go → CopyRomsToVersionInstall()
    ↓                (hard links ROMs into version dir)
launcher/launcher.go → launch()
    ↓                → getGameExecutablePath()
    ↓                → runGame() [in goroutine]
    ↓
exec.Command() with SHIP_HOME env var
```

### ROM Import Flow
```
User selects ROM file
    ↓
ui/mainwindow/roms.go → importRom()
    ↓
files/roms.go → IsValidRom()
    ↓               (SHA1 hash against known ROMs)
    ├── Valid: files/files.go → CopyRomToInstallDir()
    │           ↓
    │           stores/roms.go → AddInstalledRom()
    │           ↓
    │           db/installedroms.go → AddInstalledRom()
    │
    └── Invalid: Show error dialog
```

## Design Patterns

### 1. Store Pattern
The `stores/` package acts as an intermediary between UI/services and database:
- Provides high-level operations
- Handles data transformation
- Encapsulates database operations
- Similar to Repository pattern

### 2. Immediate Mode GUI
UI is rendered every frame based on current state:
- No explicit state management in UI layer
- UI calls `g.Update()` to trigger re-render
- State lives in database and global variables
- Simple but efficient for this use case

### 3. Panic-Based Error Handling
The codebase uses panic for unrecoverable errors:
- Database errors
- File system errors
- Configuration errors
- Recoverable errors (e.g., invalid ROM) return error values

### 4. Global State
Some state is maintained in package-level variables:
- `mainwindow.SelectedGame`: Currently selected game repository
- `db.db`: Database connection
- `settings.config`: Viper configuration

## Threading Model

### Main Thread
- UI rendering (Dear ImGui is not thread-safe)
- Event handling
- Most business logic

### Background Goroutines
- `RefreshVersions()`: Runs in background on startup
- `runGame()`: Launches game in separate goroutine
- HTTP downloads: Blocking but could be made concurrent

## Platform-Specific Considerations

### Windows
- Uses MSYS2 for compilation
- Static linking required for distribution
- `-H=windowsgui` hides console window
- Executable detection: `.exe` extension

### Linux
- Requires X11 development libraries
- GTK3 for file dialogs
- Packaged as AppImage for portability
- Executable detection: `.appimage` extension

### macOS
- Requires Xcode SDK
- Build output: currently `.dmg` (likely needs refinement)
- Executable detection: `.dmg` (should be `.app` bundles)

## Extension Points

### Adding New Games
1. Create definition in `constants/games/`
2. Define `Repository` struct with:
   - GitHub owner/repo
   - Nightly build URLs
   - Accepted ROM SHA1 hashes
   - Import capabilities
3. Add to `games.Repositories` array
4. No code changes needed elsewhere (data-driven design)

### Adding New Platforms
Would require:
1. Update `runtime.GOOS` checks in `types/version.go` and `types/unstableversion.go`
2. Add executable detection logic in `files/files.go`
3. Update build workflow in `.github/workflows/generate-builds.yml`

## Dependencies

### Core Dependencies
- `github.com/AllenDang/giu`: GUI framework
- `github.com/mattn/go-sqlite3`: SQLite driver (requires CGO)
- `github.com/google/go-github/v68`: GitHub API client
- `github.com/spf13/viper`: Configuration management

### Platform Dependencies
- `github.com/adrg/xdg`: XDG Base Directory specification
- `github.com/sqweek/dialog`: Native file dialogs
- `github.com/skratchdot/open-golang`: Cross-platform file opener

See BUILD.md for detailed dependency requirements and compilation notes.
