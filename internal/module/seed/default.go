package seed

import defaultseed "openresume/internal/module/seed/default"

func LoadDefaultSeed() (SeedData, error) {
	return LoadFromFS(defaultseed.FS, "seed.json")
}

