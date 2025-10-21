# Marina Development Guide

## Overview

This guide covers extending Marina with new features, adding support for new games, understanding the project structure, and general development practices.

## Project Structure

```
marina/
├── assets/                    # Embedded assets
│   └── vendor/fonts/         # Nerd Font for UI
├── build/                     # Platform-specific build configs
│   ├── assets/               # Application icons
│   ├── linux/                # AppImage configuration
│   └── windows/              # Windows manifest and version info
├── constants/                 # Application constants
│   └── games/                # Game port definitions
├── db/                        # Database layer (SQLite)
│   └── scripts/              # SQL schema definitions
├── files/                     # File system operations
├── launcher/                  # Game launching logic
├── services/                  # External service integration (GitHub)
├── settings/                  # Configuration management
├── stores/                    # Data access layer (store pattern)
├── types/                     # Core data structures
├── ui/                        # User interface
│   ├── dialogs/              # Confirmation & error dialogs
│   ├── fonts/                # Font management
│   ├── importconfigs/        # Config import dialogs
│   ├── mainwindow/           # Main window components
│   └── settings/             # Settings dialog
├── docs/                      # Documentation (this directory)
├── .github/workflows/         # CI/CD pipelines
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
├── marina.go                  # Application entry point
└── README.md                  # User-facing documentation
```

## Adding a New Game

Marina is designed to make adding new Harbour Masters ports straightforward. Most game-specific configuration is data-driven.

### Step 1: Create Game Definition

Create a new file in `constants/games/` for your game (e.g., `mygame.go`):

```go
package games

import marina "marina/types"

var MyGameDefinition = marina.Repository{
    Id:         3,  // Unique integer ID
    Name:       "My Game Port",  // Display name in UI
    Owner:      "HarbourMasters",  // GitHub organization/user
    Repository: "MyGamePort",      // GitHub repository name

    // Nightly build URLs from nightly.link or GitHub Actions
    LatestBuildUrls: marina.DownloadUrls{
        Windows: "https://nightly.link/HarbourMasters/MyGamePort/workflows/build/main/mygame-windows.zip",
        Linux:   "https://nightly.link/HarbourMasters/MyGamePort/workflows/build/main/mygame-linux.zip",
        Mac:     "https://nightly.link/HarbourMasters/MyGamePort/workflows/build/main/mygame-mac.zip",
    },

    // Valid ROM SHA1 hashes (find these from the game port documentation)
    AcceptedRomHashes: &[]marina.Rom{
        {
            Name: "US 1.0",
            Sha1: "1234567890abcdef1234567890abcdef12345678",
        },
        {
            Name: "EU 1.0",
            Sha1: "abcdef1234567890abcdef1234567890abcdef12",
        },
    },

    // Define what can be imported/exported between versions
    Imports: marina.Imports{
        Configuration: marina.ImportCategory{
            Capable: true,
            Files:   []string{"mygame.json", "imgui.ini"},
        },
        Mods: marina.ImportCategory{
            Capable: true,
            Files:   []string{"mods"},
        },
        Saves: marina.ImportCategory{
            Capable: true,
            Files:   []string{"Save"},
        },
        Randomizer: marina.ImportCategory{
            Capable: false,  // Not supported
            Files:   []string{},
        },
    },
}
```

### Step 2: Register Game

Add your game to the global repository list in `constants/games/games.go`:

```go
package games

import "marina/types"

var Repositories = []*marina.Repository{
    &SohDefinition,
    &TwoShipDefinition,
    &StarshipDefinition,
    &MyGameDefinition,  // Add your game here
}
```

### Step 3: Find ROM Hashes

ROM hashes must match what the game port accepts. Find these from:
1. Game port's documentation/source code
2. Hash your own ROM: `sha1sum yourrom.z64`
3. Check the port's ROM validation code

**Important**: Only include officially supported ROM versions to avoid confusion.

### Step 4: Test

Build and run Marina:
```bash
go build && ./marina
```

Your game should appear in the game selector. Test:
- ✅ Version list loads from GitHub
- ✅ Unstable version appears
- ✅ Download and install works
- ✅ ROM import validates correctly
- ✅ Game launches with proper ROM

### Step 5: Update Documentation

Update `README.md` to list your newly supported game.

## Understanding Data Flow

### Application Startup

```go
// marina.go
func main() {
    settings.Init()        // Load config from ~/.config/marina/config.toml
    files.Init()           // Create install directory structure
    assets.LoadAssets()    // Load embedded fonts
    db.Init()              // Open SQLite database, initialize schema

    for {
        ui.Start()         // Start UI event loop
        if !settings.ShouldRestart() {
            break
        }
        db.ResetDbConnection()  // Reconnect if install dir changed
    }
}
```

### UI Event Loop (Dear ImGui Immediate Mode)

```go
// ui/start.go
func Start() {
    wnd := g.NewMasterWindow("Marina", 800, 600, 0)
    fonts.InitializeFonts()
    go mainwindow.RefreshVersions(false)  // Background fetch
    wnd.Run(mainwindow.Loop)              // Run UI loop
}

// ui/mainwindow/mainwindow.go
func Loop() {
    g.SingleWindow().Layout(
        // UI is rebuilt every frame based on current state
        g.SplitLayout(g.DirectionVertical, &mainSplit,
            g.Column(GetGameSelector(), GetRomDisplay()),
            g.Column(GetToolbar(), GetGameList()),
        ),
    )
}
```

**Key Concept**: Immediate Mode GUI means the entire UI is re-rendered each frame (~60 FPS). State is stored externally (database, global variables), not in UI components.

## Configuration System

### Settings Management

Configuration uses Viper library with TOML format.

**Config File Location**:
- Linux: `~/.config/marina/config.toml`
- Windows: `%APPDATA%/marina/config.toml`
- macOS: `~/Library/Application Support/marina/config.toml`

**Current Settings**:
```toml
InstallDir = "/home/user/.local/share/marina"
```

**Adding New Settings**:

1. Define default in `settings/settings.go`:
```go
func setDefaults() {
    config.SetDefault("InstallDir", GetDefaultInstallDir())
    config.SetDefault("MyNewSetting", "default_value")  // Add here
}
```

2. Add getter/setter:
```go
func GetMyNewSetting() string {
    return config.GetString("MyNewSetting")
}

func SetMyNewSetting(value string) {
    config.Set("MyNewSetting", value)
    saveChanges()
}
```

3. Use in your code:
```go
setting := settings.GetMyNewSetting()
```

**Settings that Require Restart**: If a setting change requires restarting (like `InstallDir`), set `restart = true` in the setter.

## Database Operations

### Adding New Tables

1. Update schema in `db/scripts/init.go`:
```sql
const DbSchemaInit = `
    -- Existing tables...

    CREATE TABLE IF NOT EXISTS MyNewTable (
        Id INTEGER PRIMARY KEY AUTOINCREMENT,
        Data TEXT NOT NULL
    );
`
```

2. Create operations file `db/mynew.go`:
```go
package db

import "fmt"

const getMyDataScript = `SELECT Id, Data FROM MyNewTable`

func GetMyData() []MyData {
    rows, err := db.Query(getMyDataScript)
    if err != nil {
        panic(fmt.Errorf("Error querying MyNewTable: %w", err))
    }
    defer rows.Close()

    var results []MyData
    for rows.Next() {
        var item MyData
        if err := rows.Scan(&item.Id, &item.Data); err != nil {
            panic(fmt.Errorf("Error scanning row: %w", err))
        }
        results = append(results, item)
    }
    return results
}
```

3. Use store pattern for higher-level access (see `stores/` package).

**Migration Note**: Marina uses `CREATE TABLE IF NOT EXISTS`, so schema changes require manual migration or database deletion.

## Adding UI Components

### Creating a New Dialog

Example: Adding a "About" dialog.

1. Create `ui/dialogs/about.go`:
```go
package dialogs

import g "github.com/AllenDang/giu"

var showAboutDialog = false

func ShowAboutDialog() {
    showAboutDialog = true
    g.Update()
}

func GetAboutDialog() *g.PopupModalWidget {
    if !showAboutDialog {
        return g.PopupModal("")  // Empty when not shown
    }

    return g.PopupModal("About Marina").Layout(
        g.Label("Marina v1.0.0"),
        g.Label("Unofficial launcher for Harbour Masters ports"),
        g.Separator(),
        g.Button("Close").OnClick(func() {
            showAboutDialog = false
        }),
    )
}
```

2. Add to main layout in `ui/mainwindow/mainwindow.go`:
```go
func Loop() {
    g.SingleWindow().Layout(
        // Existing layout...
        dialogs.GetAboutDialog(),  // Add dialog
    )
}
```

3. Trigger from toolbar or menu:
```go
g.Button("About").OnClick(dialogs.ShowAboutDialog)
```

### Styling UI Components

Marina uses custom colors and fonts defined in `ui/fonts/fonts.go`:

```go
// Using custom color
g.Style().SetColor(g.StyleColorText, fonts.ColorCaption).To(
    g.Label("This is a caption"),
)

// Using custom font size
g.Style().SetFontSize(fonts.CaptionSize).To(
    g.Label("Small text"),
)
```

**Available Style Variables** (see `ui/fonts/fonts.go`):
- `fonts.HeaderSize`
- `fonts.CaptionSize`
- `fonts.ColorCaption`
- `fonts.ColorDestructive`

## GitHub API Integration

### Rate Limiting

GitHub API has strict rate limits:
- **Unauthenticated**: 60 requests/hour per IP
- **Authenticated**: 5000 requests/hour (not currently implemented)

**Current Strategy**:
- Cache version data in database
- Only refresh if > 1 hour since last fetch
- Manual refresh bypasses cache

**Checking Last Fetch**:
```go
lastFetch := stores.GetLastFetched(repository)
if lastFetch != nil && time.Since(*lastFetch) < time.Hour {
    // Use cached data
} else {
    // Fetch from GitHub
}
```

### Adding GitHub Authentication (Future Enhancement)

To increase rate limits, add GitHub token support:

1. Add setting for GitHub token:
```go
// settings/settings.go
func GetGitHubToken() string {
    return config.GetString("GitHubToken")
}
```

2. Update GitHub client in `services/github.go`:
```go
import "golang.org/x/oauth2"

func createClient() *github.Client {
    token := settings.GetGitHubToken()
    if token != "" {
        ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
        tc := oauth2.NewClient(context.Background(), ts)
        return github.NewClient(tc)
    }
    return github.NewClient(nil)  // Unauthenticated
}
```

3. Add UI for token input in settings dialog.

## File System Operations

### Directory Structure

Marina creates this structure in the install directory:

```
<InstallDir>/
├── manifest.db                # SQLite database
├── roms/                      # Validated ROMs
│   ├── pal-1_0.z64
│   └── ntsc-u-1_0.z64
└── versions/                  # Installed versions
    ├── shipwright/
    │   ├── v1_0_0/           # Stable release
    │   │   ├── soh.appimage
    │   │   └── pal-1_0.z64   # Hard-linked ROM
    │   └── unstable-2025-01-15/  # Development build
    └── 2ship2harkinian/
        └── v1_0_0/
```

### Adding New File Operations

Example: Adding backup functionality.

1. Create function in `files/` package:
```go
// files/backup.go
package files

import (
    "marina/types"
    "path/filepath"
    cp "github.com/otiai10/copy"
)

func BackupVersion(version *marina.Version) error {
    srcDir := GetVersionInstallDirPath(version)
    backupDir := filepath.Join(
        settings.GetInstallDirName(),
        "backups",
        DirName(version.Repository.Repository),
        DirName(version.TagName),
    )

    return cp.Copy(srcDir, backupDir)
}
```

2. Add UI button in `ui/mainwindow/versions.go`:
```go
func (i *VersionListItem) getButtons() *g.RowWidget {
    // Existing buttons...
    g.Button("💾 Backup").OnClick(i.backup),
}

func (i *VersionListItem) backup() {
    err := files.BackupVersion(i.StableVersion)
    dialogs.ShowDialogIfError(err)
}
```

## Testing Strategy

### Current State

As noted in the README: "If I was a good developer there would be tests."

### Recommendations for Adding Tests

1. **Unit Tests for Business Logic**:
```go
// files/files_test.go
func TestDirName(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {"Ship of Harkinian", "ship-of-harkinian"},
        {"2 Ship 2 Harkinian", "2-ship-2-harkinian"},
    }

    for _, tt := range tests {
        result := files.DirName(tt.input)
        if result != tt.expected {
            t.Errorf("DirName(%q) = %q, want %q", tt.input, result, tt.expected)
        }
    }
}
```

2. **Database Tests with In-Memory SQLite**:
```go
// db/db_test.go
func TestDatabaseInit(t *testing.T) {
    // Use :memory: database for testing
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()

    // Test schema initialization
    _, err = db.Exec(scripts.DbSchemaInit)
    if err != nil {
        t.Fatalf("Failed to initialize schema: %v", err)
    }

    // Verify tables exist
    tables := []string{"Releases", "UnstableVersions", "InstalledRoms", "LastFetchedDate"}
    for _, table := range tables {
        var name string
        err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
        if err != nil {
            t.Errorf("Table %s not found", table)
        }
    }
}
```

3. **Mock GitHub API for Service Tests**:
```go
// services/github_test.go
// Use httptest to mock GitHub responses
```

4. **UI Tests** (Challenging with ImGui):
- Focus on testing business logic separately from UI
- Consider using headless testing if feasible
- Manual QA checklist for UI changes

## Common Development Tasks

### Running in Development Mode

```bash
# Build with race detector
go build -race -o marina-dev

# Run with verbose logging
./marina-dev -v
```

### Debugging

1. **Add Debug Prints**:
```go
import "fmt"

fmt.Printf("DEBUG: Version: %+v\n", version)
```

2. **Use Delve Debugger**:
```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Run with debugger
dlv debug marina
```

3. **Enable ImGui Demo Window**:

In `ui/mainwindow/mainwindow.go`, change:
```go
func Loop() {
    if true {  // Change from false to true
        imgui.ShowDemoWindow()
        return
    }
    // ...
}
```

### Changing Database Schema

Since Marina doesn't have migrations:

**Option 1: Development** - Delete database and rebuild:
```bash
rm ~/.local/share/marina/manifest.db
./marina
```

**Option 2: Production** - Write manual migration:
```sql
-- migrations/001_add_column.sql
ALTER TABLE Releases ADD COLUMN NewColumn TEXT DEFAULT '';
```

Then execute in code or provide to users.

### Hot Reload (Not Supported)

Dear ImGui and CGO make hot reload challenging. Recommended workflow:
1. Make changes
2. Rebuild: `go build`
3. Restart application
4. Test changes

**Tip**: Keep build times fast by avoiding unnecessary rebuilds of CGO dependencies.

## Code Style & Conventions

### Naming Conventions

- **Packages**: lowercase, single word (e.g., `files`, `services`)
- **Files**: lowercase, descriptive (e.g., `versions.go`, `github.go`)
- **Types**: PascalCase (e.g., `Version`, `Repository`)
- **Functions**: PascalCase for exported, camelCase for internal
- **Variables**: camelCase

### Error Handling

Marina uses two error patterns:

1. **Panic for Unrecoverable Errors**:
```go
if err != nil {
    panic(fmt.Errorf("Error opening database: %w", err))
}
```
Used for: Database errors, file system errors, configuration errors

2. **Return Error for Recoverable Errors**:
```go
func DownloadVersion(version *marina.Version) error {
    // ...
    if err != nil {
        return fmt.Errorf("Download failed: %w", err)
    }
    return nil
}
```
Used for: Network errors, validation errors, user-correctable issues

### Comments

- **Package Comments**: Brief description at top of main file
- **Function Comments**: For exported functions, especially complex ones
- **Inline Comments**: For non-obvious logic

```go
// GetVersionInstallDirPath returns the filesystem path where a version is installed.
// Example: ~/.local/share/marina/versions/shipwright/v1_0_0
func GetVersionInstallDirPath(version *marina.Version) string {
    return filepath.Join(
        settings.GetInstallDirName(),
        "versions",
        DirName(version.Repository.Repository),
        DirName(version.TagName),
    )
}
```

## Contributing Guidelines

### Before Submitting a PR

1. **Test Your Changes**: Build on all platforms if possible (use CI)
2. **Update Documentation**: Update README.md or docs/ if adding features
3. **Follow Code Style**: Match existing patterns
4. **Keep Commits Clean**: Descriptive commit messages

### PR Process

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-new-feature`
3. Make your changes
4. Test locally: `go build && ./marina`
5. Commit: `git commit -m "Add support for XYZ"`
6. Push: `git push origin feature/my-new-feature`
7. Open Pull Request on GitHub

### CI/CD

GitHub Actions automatically builds for all platforms. Check:
- ✅ Linux build succeeds
- ✅ Windows build succeeds
- ✅ macOS build succeeds

## Roadmap Ideas

Potential features for future development:

### High Priority
- [ ] GitHub token authentication (increase rate limits)
- [ ] Better macOS packaging (proper .app bundles)
- [ ] ROM management UI (delete ROMs)
- [ ] Config import/export between games
- [ ] Update notifications

### Medium Priority
- [ ] Download progress bars
- [ ] Multi-language support (i18n)
- [ ] Theme customization
- [ ] Backup/restore functionality
- [ ] Search/filter version list

### Low Priority
- [ ] Plugin system for custom games
- [ ] Cloud sync for saves
- [ ] Mod manager integration
- [ ] Performance analytics

## Resources

### Go Development
- **Go Documentation**: https://go.dev/doc/
- **Effective Go**: https://go.dev/doc/effective_go
- **Go Modules**: https://go.dev/ref/mod

### Libraries
- **giu (Dear ImGui)**: https://github.com/AllenDang/giu
- **Viper (Config)**: https://github.com/spf13/viper
- **go-sqlite3**: https://github.com/mattn/go-sqlite3
- **go-github**: https://github.com/google/go-github

### Harbour Masters
- **Ship of Harkinian**: https://github.com/HarbourMasters/Shipwright
- **2 Ship 2 Harkinian**: https://github.com/HarbourMasters/2ship2harkinian
- **Starship**: https://github.com/HarbourMasters/Starship

## Getting Help

- **Issues**: https://github.com/Cardboy777/marina/issues
- **Discussions**: Use GitHub Discussions for questions
- **Harbour Masters Discord**: For game-specific questions

## License

Marina is licensed under the MIT License. See LICENSE file for details.

---

**Happy hacking!** If you improve Marina, consider contributing your changes back to the project.
