# Marina Database Documentation

## Overview

Marina uses SQLite as its persistence layer to track:
- Downloaded game versions and their installation state
- Validated ROM files
- Unstable/development builds
- GitHub API fetch timestamps (for rate limiting)

**Database Location**: `<InstallDir>/manifest.db`

**Default Install Directory**:
- Linux: `~/.local/share/marina/`
- Windows: `%LOCALAPPDATA%/marina/`
- macOS: `~/Library/Application Support/marina/`

## Database Schema

### Table: Releases

Stores information about stable game releases fetched from GitHub.

```sql
CREATE TABLE IF NOT EXISTS Releases (
    Owner TEXT NOT NULL,                        -- GitHub repository owner
    Repository TEXT NOT NULL,                   -- GitHub repository name
    TagName TEXT NOT NULL,                      -- Git tag (e.g., "v1.0.0")
    Name TEXT NOT NULL,                         -- Display name
    WindowsDownloadUrl TEXT NOT NULL,           -- Windows ZIP download URL
    MacDownloadUrl TEXT NOT NULL,               -- macOS download URL
    LinuxPerformanceDownloadUrl TEXT NOT NULL,  -- Linux download URL
    LinuxCompatabilityDownloadUrl TEXT NOT NULL,-- Linux compatibility build URL
    ReleaseDate DATETIME2 NOT NULL,             -- Release timestamp
    Installed BIT NOT NULL,                     -- Installation status (0/1)
    PRIMARY KEY (Owner, Repository, TagName)
);
```

**Primary Key**: Composite key of `(Owner, Repository, TagName)` ensures one entry per release.

**Indexes**: None explicitly defined (relies on primary key).

**Usage**:
- Populated by `services.SyncReleases()` via GitHub API
- Queried by `stores.GetVersions()` for UI display
- Updated by `db.SetInstalled()` when versions are downloaded/deleted

**Note**: The schema has a typo - `LinuxCompatabilityDownloadUrl` should be `LinuxCompatibilityDownloadUrl`.

### Table: UnstableVersions

Stores unstable/development builds (nightly builds from the latest commit).

```sql
CREATE TABLE IF NOT EXISTS UnstableVersions (
    CommitHash TEXT PRIMARY KEY,              -- Git commit SHA
    Owner TEXT NOT NULL,                      -- GitHub repository owner
    Repository TEXT NOT NULL,                 -- GitHub repository name
    ReleaseDate DATETIME2 NOT NULL,           -- Commit timestamp
    Installed BIT NOT NULL                    -- Installation status (0/1)
);
```

**Primary Key**: `CommitHash` (unique identifier for each build).

**Cleanup Strategy**: Old unstable versions are removed by `stores.RemoveOldUnstableVersions()` to prevent database bloat.

**Usage**:
- Populated by `services.fetchLatestCommit()` via GitHub API
- Only the most recent unstable build is typically kept
- Displays at the top of version list with "Unstable - <timestamp>" format

### Table: InstalledRoms

Tracks ROM files that have been validated and imported into Marina.

```sql
CREATE TABLE IF NOT EXISTS InstalledRoms (
    Hash TEXT PRIMARY KEY,                    -- SHA1 hash of ROM file
    Name TEXT NOT NULL,                       -- ROM variant name (e.g., "PAL 1.0")
    Owner TEXT NOT NULL,                      -- Associated game owner
    Repository TEXT NOT NULL                  -- Associated game repository
);
```

**Primary Key**: `Hash` (SHA1) - prevents duplicate ROMs.

**ROM Validation**:
- ROMs are validated against hardcoded SHA1 hashes in `constants/games/`
- Only validated ROMs are stored in this table
- ROM files are stored in `<InstallDir>/roms/` directory

**Usage**:
- ROMs are game-specific (different games require different ROMs)
- When launching a game, ROMs for that repository are hard-linked into the version directory
- One ROM entry can be used by multiple versions of the same game

### Table: LastFetchedDate

Tracks the last time version data was fetched from GitHub to implement rate limiting.

```sql
CREATE TABLE IF NOT EXISTS LastFetchedDate (
    Timestamp DATETIME2 NOT NULL,             -- Last fetch time
    Owner TEXT NOT NULL,                      -- GitHub repository owner
    Repository TEXT NOT NULL                  -- GitHub repository name
);
```

**Primary Key**: None explicitly defined (should probably be `(Owner, Repository)`).

**Rate Limiting Strategy**:
- Marina only fetches new versions once per hour
- Manual refresh bypasses this restriction
- Prevents hitting GitHub API rate limits (60 requests/hour for unauthenticated)

**Usage**:
- Checked by `services.SyncReleases()` before fetching
- Updated after successful fetch via `stores.UpdateLastFetched()`

## Data Models

### Go Struct ↔ Database Mappings

#### Version (types/version.go) ↔ Releases Table

```go
type Version struct {
    Repository   *Repository    // Foreign reference to game definition
    Name         string         // → Releases.Name
    TagName      string         // → Releases.TagName
    DownloadUrls DownloadUrls   // → Windows/Mac/Linux URLs
    ReleaseDate  time.Time      // → Releases.ReleaseDate
    Installed    bool           // → Releases.Installed
}
```

**Download URL Mapping**:
```go
type DownloadUrls struct {
    Windows            string  // → Releases.WindowsDownloadUrl
    Linux              string  // → Releases.LinuxPerformanceDownloadUrl
    LinuxCompatibility string  // → Releases.LinuxCompatabilityDownloadUrl
    Mac                string  // → Releases.MacDownloadUrl
}
```

#### UnstableVersion (types/unstableversion.go) ↔ UnstableVersions Table

```go
type UnstableVersion struct {
    Hash        string        // → UnstableVersions.CommitHash
    ReleaseDate time.Time     // → UnstableVersions.ReleaseDate
    Installed   bool          // → UnstableVersions.Installed
    Repository  *Repository   // Foreign reference to game definition
}
```

**Note**: Unstable versions don't have stored download URLs - they use `Repository.LatestBuildUrls` directly.

#### Rom (types/rom.go) ↔ InstalledRoms Table

```go
type Rom struct {
    Name string  // → InstalledRoms.Name
    Sha1 string  // → InstalledRoms.Hash
}
```

**Association**: ROMs are associated with repositories via `(Owner, Repository)` columns.

#### Repository (types/repository.go)

```go
type Repository struct {
    Id                int
    Name              string              // Display name
    Owner             string              // → DB queries
    Repository        string              // → DB queries
    PathVariableName  string              // Not used in current code
    LatestBuildUrls   DownloadUrls        // Nightly build URLs
    AcceptedRomHashes *[]Rom              // Valid ROM hashes (hardcoded)
    Imports           Imports             // Import capabilities
}
```

**Note**: Repository definitions are hardcoded in `constants/games/` and not stored in the database.

## Database Operations

### Initialization

**File**: `db/db.go`

```go
func Init() {
    openDbConnection()
    // Executes schema from db/scripts/init.go
    db.Exec(scripts.DbSchemaInit)
}
```

**Schema Migration**: Marina uses a simple "CREATE TABLE IF NOT EXISTS" approach. There is no formal migration system.

### CRUD Operations

#### Releases Table

**Create**:
```go
// db/versions.go
func AddRelease(v marina.Version)
```
- Inserts new release
- Called from `stores.AddVersion()` after deduplication check

**Read**:
```go
// db/versions.go
func GetVersionList(repo *marina.Repository) []marina.Version
```
- Returns all versions for a repository
- Ordered by `ReleaseDate DESC`
- Populates `Version` structs with data

**Update**:
```go
// db/versions.go
func SetInstalled(v *marina.Version, isInstalled bool)
```
- Updates `Installed` flag
- Called after download/deletion operations

**Delete**: Not implemented (versions are never deleted from database, only marked as uninstalled).

#### InstalledRoms Table

**Create**:
```go
// db/installedroms.go
func AddInstalledRom(rom marina.Rom, repo *marina.Repository)
```
- Inserts validated ROM
- Primary key prevents duplicates

**Read**:
```go
// db/installedroms.go
func GetInstalledRomList(repo *marina.Repository) []marina.Rom
```
- Returns all ROMs for a specific repository
- Used to check if game can be launched

**Update/Delete**: Not implemented.

#### UnstableVersions Table

**Create**:
```go
// db/unstableversions.go
func AddUnstableVersion(v *marina.UnstableVersion)
```
- Inserts latest commit version
- Primary key prevents duplicates

**Read**:
```go
// db/unstableversions.go
func GetUnstableVersionList(repo *marina.Repository) []marina.UnstableVersion
```
- Returns unstable versions for repository
- Ordered by `ReleaseDate DESC`

**Update**:
```go
// db/unstableversions.go
func SetUnstableVersionInstalled(v *marina.UnstableVersion, installed bool)
```

**Delete**:
```go
// db/unstableversions.go
func RemoveAllUninstalledUnstableVersions()
```
- Removes uninstalled unstable versions
- Prevents database bloat from nightly builds
- Called after unstable version deletion

#### LastFetchedDate Table

**Read**:
```go
// db/lastupdated.go
func GetLastFetched(repo *marina.Repository) *time.Time
```
- Returns last fetch timestamp for repository
- Returns `nil` if never fetched

**Update**:
```go
// db/lastupdated.go
func UpdateLastFetched(repo *marina.Repository, timestamp time.Time)
```
- Uses `INSERT OR REPLACE` pattern
- Called after successful version sync

## Data Flow Examples

### Version Sync Flow

```
services.SyncReleases(repo, force)
    ↓
Check: stores.GetVersions(repo) - returns cached data
    ↓
Check: stores.GetLastFetched(repo)
    ↓
If > 1 hour OR force:
    ↓
    services.fetchReleaseVersions(repo)
        ↓
        GitHub API: GET /repos/:owner/:repo/releases
        ↓
        Parse assets → DownloadUrls
        ↓
        For each release:
            stores.AddVersion(version)
                ↓
                Check if exists: GetVersions()
                ↓
                If new: db.AddRelease(version)
                    ↓
                    INSERT INTO Releases (...)
    ↓
    services.fetchLatestCommit(repo)
        ↓
        GitHub API: GET /repos/:owner/:repo/branches/:default
        ↓
        stores.AddUnstableVersion(unstableVersion)
            ↓
            db.AddUnstableVersion(unstableVersion)
                ↓
                INSERT INTO UnstableVersions (...)
    ↓
    stores.UpdateLastFetched(repo, now)
        ↓
        db.UpdateLastFetched(repo, now)
            ↓
            INSERT OR REPLACE INTO LastFetchedDate (...)
```

### ROM Import Flow

```
User selects ROM file
    ↓
files.IsValidRom(validHashes, filepath)
    ↓
    Read file → SHA1 hash
    ↓
    Compare against constants/games/*.go ROM hashes
    ↓
    If valid:
        files.CopyRomToInstallDir(repo, filepath)
            ↓
            Copy to <InstallDir>/roms/<name>.z64
            ↓
            stores.AddInstalledRom(rom, repo)
                ↓
                db.AddInstalledRom(rom, repo)
                    ↓
                    INSERT INTO InstalledRoms (Hash, Name, Owner, Repository)
```

### Game Launch Flow

```
launcher.LaunchGame(version, callback)
    ↓
files.CopyRomsToVersionInstall(version)
    ↓
    stores.GetInstalledRomsList(version.Repository)
        ↓
        db.GetInstalledRomList(repo)
            ↓
            SELECT Hash, Name FROM InstalledRoms
            WHERE Owner = ? AND Repository = ?
    ↓
    For each ROM: hard link from roms/ to version install dir
    ↓
    Detect executable in install dir
    ↓
    exec.Command() with SHIP_HOME env var
```

## Connection Management

### Single Connection Pattern
Marina maintains a single global SQLite connection:

```go
// db/db.go
var db *sql.DB

func Init() {
    openDbConnection()
}
```

**Thread Safety**: SQLite connections are thread-safe with the `mattn/go-sqlite3` driver when using WAL mode, but Marina primarily operates on the main thread.

### Connection Reset
When the install directory changes:

```go
func ResetDbConnection() {
    closeDbConnection()
    openDbConnection()
}
```

This is necessary because the database file location changes with the install directory.

### Busy Timeout
Connection string includes a 5-second busy timeout:

```go
dbFilePath := fmt.Sprintf("%s?_busy_timeout=5000", filepath.Join(installDir, dbFileName))
```

## Known Issues & Design Decisions

### Schema Issues

1. **Typo in Column Name**: `LinuxCompatabilityDownloadUrl` should be `LinuxCompatibilityDownloadUrl`
2. **No Primary Key on LastFetchedDate**: Could lead to duplicate entries for same repository
3. **No Foreign Keys**: Tables reference repositories but don't enforce referential integrity

### Migration Strategy

The application uses `CREATE TABLE IF NOT EXISTS` which means:
- ✅ Safe for first-time initialization
- ✅ Non-destructive on upgrades
- ❌ No way to modify existing schema
- ❌ Column renames/changes require manual SQL or database deletion

**Recommendation**: Consider using a migration tool like `golang-migrate/migrate` for future schema changes.

### Data Retention

- **Stable Releases**: Never deleted from database (only marked uninstalled)
- **Unstable Versions**: Cleaned up when uninstalled
- **ROMs**: Never deleted (no UI for removal)
- **Fetch Timestamps**: Never cleaned up

**Potential Issue**: Database will grow over time with old release data. Not a practical problem (text data is tiny), but could be optimized.

### Time Format

Time values are stored as `DATETIME2` but formatted as Go's `time.Layout` format:
```go
time.Format(time.Layout) // "Mon Jan 2 15:04:05 MST 2006"
```

This is Go-specific and not a standard SQLite datetime format. Consider using ISO 8601 or Unix timestamps for better portability.

## Querying Tips

### Useful Queries for Development/Debugging

**List all installed versions**:
```sql
SELECT Owner, Repository, TagName, Name, Installed
FROM Releases
WHERE Installed = 1;
```

**Check installed ROMs**:
```sql
SELECT r.Name, r.Hash, r.Owner, r.Repository
FROM InstalledRoms r;
```

**View fetch history**:
```sql
SELECT Owner, Repository, Timestamp
FROM LastFetchedDate
ORDER BY Timestamp DESC;
```

**Find unstable versions**:
```sql
SELECT Owner, Repository, CommitHash, ReleaseDate, Installed
FROM UnstableVersions
ORDER BY ReleaseDate DESC;
```

**Database statistics**:
```sql
SELECT
    (SELECT COUNT(*) FROM Releases) as stable_versions,
    (SELECT COUNT(*) FROM UnstableVersions) as unstable_versions,
    (SELECT COUNT(*) FROM InstalledRoms) as installed_roms,
    (SELECT COUNT(*) FROM LastFetchedDate) as fetch_records;
```

## Testing Recommendations

Since the README notes "If I was a good developer there would be tests", here are recommendations for database testing:

1. **Use in-memory SQLite** for tests: `:memory:` connection string
2. **Test schema initialization**: Verify all tables exist after `Init()`
3. **Test CRUD operations**: Ensure inserts/updates/deletes work correctly
4. **Test deduplication**: Verify primary keys prevent duplicates
5. **Test time parsing**: Ensure dates round-trip correctly
6. **Test concurrent access**: If planning to add threading

Example test structure:
```go
func TestDatabaseInit(t *testing.T) {
    // Use temporary directory for test database
    tmpDir := t.TempDir()
    // Set install dir to temp
    // Call db.Init()
    // Verify tables exist
}
```
