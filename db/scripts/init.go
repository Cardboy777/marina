package scripts

const DbSchemaInit = `
	CREATE TABLE IF NOT EXISTS Releases (
		Owner TEXT NOT NULL,
		Repository TEXT NOT NULL,
		TagName TEXT NOT NULL,
		Name TEXT NOT NULL,
		WindowsDownloadUrl TEXT NOT NULL,
		MacDownloadUrl TEXT NOT NULL,
		LinuxPerformanceDownloadUrl TEXT NOT NULL,
		LinuxCompatabilityDownloadUrl TEXT NOT NULL,
		ReleaseDate DATETIME2 NOT NULL,
		Installed BIT NOT NULL,
		PRIMARY KEY (Owner, Repository, TagName)
	);

	CREATE TABLE IF NOT EXISTS UnstableVersions (
		CommitHash TEXTPRIMARY KEY,
		Owner TEXT NOT NULL,
		Repository TEXT NOT NULL,
		ReleaseDate DATETIME2 NOT NULL,
		Installed BIT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS InstalledRoms (
		Hash TEXT PRIMARY KEY,
		Name TEXT NOT NULL,
		Owner TEXT NOT NULL,
		Repository TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS LastFetchedDate (
		Timestamp DATETIME2 NOT NULL,
		Owner TEXT NOT NULL,
		Repository TEXT NOT NULL
	);
`
