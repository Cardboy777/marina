package marina

type Repository struct {
	Id                int
	Name              string
	Owner             string
	Repository        string
	PathVariableName  string
	LatestBuildUrls   DownloadUrls
	AcceptedRomHashes *[]Rom
}
