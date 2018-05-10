package main

import (
	"encoding/csv"
	"os"
)

func init() {
	partnerFile := "./partnerFile.csv"

	f, err := os.Open(partnerFile)
	check(err)
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	check(err)

	for _, line := range lines {
		scope := []string{}
		for i := 2; i < len(line); i++ {
			scope = append(scope, line[i])
		}
		p := partner{line[0], line[1], scope}
		partners = append(partners, p)
	}
}
