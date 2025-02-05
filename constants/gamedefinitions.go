package constants

import "marina/types"

var (
	SohDefinition = marina.Repository{
		Id:         0,
		Name:       "Ship of Harkinian",
		Owner:      "HarbourMasters",
		Repository: "Shipwright",
		LatestBuildUrls: marina.DownloadUrls{
			Windows: "https://nightly.link/HarbourMasters/Shipwright/workflows/generate-builds/develop/soh-windows.zip",
			Linux:   "https://nightly.link/HarbourMasters/Shipwright/workflows/generate-builds/develop/soh-linux.zip",
			Mac:     "https://nightly.link/HarbourMasters/Shipwright/workflows/generate-builds/develop/soh-mac.zip",
		},
		AcceptedRomHashes: &[]marina.Rom{
			{
				Name: "PAL 1.0",
				Sha1: "328a1f1beba30ce5e178f031662019eb32c5f3b5",
			},
			{
				Name: "PAL 1.1",
				Sha1: "cfbb98d392e4a9d39da8285d10cbef3974c2f012",
			},
			{
				Name: "PAL GC",
				Sha1: "0227d7c0074f2d0ac935631990da8ec5914597b4",
			},
			{
				Name: "PAL MQ",
				Sha1: "f46239439f59a2a594ef83cf68ef65043b1bffe2",
			},
			{
				Name: "PAL GC (Debug)",
				Sha1: "cee6bc3c2a634b41728f2af8da54d9bf8cc14099",
			},
			{
				Name: "PAL MQ (Debug)",
				Sha1: "079b855b943d6ad8bd1eb026c0ed169ecbdac7da",
			},
			{
				Name: "PAL MQ (Debug)",
				Sha1: "50bebedad9e0f10746a52b07239e47fa6c284d03",
			},
			{
				Name: "PAL MQ (Debug)",
				Sha1: "cfecfdc58d650e71a200c81f033de4e6d617a9f6",
			},
		},
	}
	TwoShipDefinition = marina.Repository{
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
	}
	StarshipDefinition = marina.Repository{
		Id:         2,
		Name:       "Starship",
		Owner:      "HarbourMasters",
		Repository: "Starship",
		LatestBuildUrls: marina.DownloadUrls{
			Windows: "https://nightly.link/HarbourMasters/Starship/workflows/main/main/starship-windows.zip",
			Linux:   "https://nightly.link/HarbourMasters/Starship/workflows/main/main/starship-mac-x64.zip",
			Mac:     "https://nightly.link/HarbourMasters/Starship/workflows/main/main/starship-linux.zip",
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
	}
)

var Repositories = []*marina.Repository{&SohDefinition, &TwoShipDefinition, &StarshipDefinition}
