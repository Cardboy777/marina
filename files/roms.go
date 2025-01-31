package files

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"marina/types"
	"os"
)

func IsValidRom(validHashes *[]marina.RomDefinition, filepath string) (*[]byte, marina.RomDefinition, bool) {
	hasher := sha1.New()

	file, err := os.Open(filepath)
	if err != nil {
		panic(fmt.Errorf("Error reading rom file: %w", err))
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf("Error reading rom file: %w", err))
	}

	hasher.Write(bytes)

	hash := fmt.Sprint(hex.EncodeToString(hasher.Sum(nil)))

	for _, h := range *validHashes {
		if hash == h.Sha1 {
			return &bytes, h, true
		}
	}

	return nil, marina.RomDefinition{}, false
}
