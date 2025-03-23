package marina

type ImportCategory struct {
	Capable bool
	Files   []string
}

type Imports struct {
	Configuration ImportCategory
	Mods          ImportCategory
	Saves         ImportCategory
	Randomizer    ImportCategory
}

type Repository struct {
	Id                int
	Name              string
	Owner             string
	Repository        string
	PathVariableName  string
	LatestBuildUrls   DownloadUrls
	AcceptedRomHashes *[]Rom
	Imports           Imports
}
