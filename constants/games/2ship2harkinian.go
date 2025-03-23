package games

import marina "marina/types"

var TwoShipDefinition = marina.Repository{
	Id:         1,
	Name:       "2 Ship 2 Harkinian",
	Owner:      "HarbourMasters",
	Repository: "2Ship2Harkinian",
	LatestBuildUrls: marina.DownloadUrls{
		Windows: "https://nightly.link/HarbourMasters/2ship2harkinian/workflows/main/develop/2ship-windows.zip",
		Linux:   "https://nightly.link/HarbourMasters/2ship2harkinian/workflows/main/develop/2ship-linux.zip",
		Mac:     "https://nightly.link/HarbourMasters/2ship2harkinian/workflows/main/develop/2ship-mac.zip",
	},
	AcceptedRomHashes: &[]marina.Rom{
		{
			Name: "NTSC-U 1.0",
			Sha1: "d6133ace5afaa0882cf214cf88daba39e266c078",
		},
		{
			Name: "NTSC-U GC",
			Sha1: "9743aa026e9269b339eb0e3044cd5830a440c1fd",
		},
	},
	Imports: marina.Imports{
		Configuration: marina.ImportCategory{
			Capable: true,
			Files:   []string{"2ship2harkinian.json", "imgui.ini", "presets"},
		},
		Mods: marina.ImportCategory{
			Capable: true,
			Files:   []string{"mods"},
		},
		Saves: marina.ImportCategory{
			Capable: true,
			Files:   []string{"saves"},
		},
		Randomizer: marina.ImportCategory{
			Capable: true,
			Files:   []string{"randomizer"},
		},
	},
}
