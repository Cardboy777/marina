package games

import marina "marina/types"

var StarshipDefinition = marina.Repository{
	Id:         2,
	Name:       "Starship",
	Owner:      "HarbourMasters",
	Repository: "Starship",
	LatestBuildUrls: marina.DownloadUrls{
		Windows: "https://nightly.link/HarbourMasters/Starship/workflows/main/main/starship-windows.zip",
		Linux:   "https://nightly.link/HarbourMasters/Starship/workflows/main/main/Starship-linux.zip",
		Mac:     "https://nightly.link/HarbourMasters/Starship/workflows/main/main/starship-mac-x64.zip",
	},
	AcceptedRomHashes: &[]marina.Rom{
		{
			Name: "USA 1.1 Rev A (Compressed)",
			Sha1: "09f0d105f476b00efa5303a3ebc42e60a7753b7a",
		},
		{
			Name: "USA 1.1 Rev A",
			Sha1: "f7475fb11e7e6830f82883412638e8390791ab87",
		},
	},
	Imports: marina.Imports{
		Configuration: marina.ImportCategory{
			Capable: true,
			Files:   []string{"starship.cfg.json", "imgui.ini", "gamecontrollerdb.txt"},
		},
		Mods: marina.ImportCategory{
			Capable: true,
			Files:   []string{"mods"},
		},
		Saves: marina.ImportCategory{
			Capable: true,
			Files:   []string{"default.sav"},
		},
		Randomizer: marina.ImportCategory{
			Capable: false,
			Files:   []string{},
		},
	},
}
