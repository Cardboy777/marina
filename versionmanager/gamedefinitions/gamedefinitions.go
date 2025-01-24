package gamedefinitions

import "marina/types"

var (
	SohDefinition = marina.RepositoryDefinition{
		Id:                0,
		Name:              "Ship of Harkinian",
		Owner:             "HarbourMasters",
		Repository:        "Shipwright",
		AcceptedRomHashes: []string{},
	}
	TwoShipDefinition = marina.RepositoryDefinition{
		Id:                1,
		Name:              "2 Ship 2 Harkinian",
		Owner:             "HarbourMasters",
		Repository:        "2Ship2Harkinian",
		AcceptedRomHashes: []string{},
	}
	StarshipDefinition = marina.RepositoryDefinition{
		Id:                2,
		Name:              "Starship",
		Owner:             "HarbourMasters",
		Repository:        "Starship",
		AcceptedRomHashes: []string{},
	}
)

var RepositoryDefinitions = []*marina.RepositoryDefinition{&SohDefinition, &TwoShipDefinition, &StarshipDefinition}
