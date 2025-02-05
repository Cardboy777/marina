package marina

import "time"

type UnstableVersion struct {
	Hash        string
	ReleaseDate time.Time
	Repository  *Repository
}
