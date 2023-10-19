package generator

import (
	"regexp"

	"github.com/onflow/flixkit-go"
)

func ExtractImports(cadenceCode string) []string {
	// Regex pattern to match Cadence import lines
	pattern := `import [\w\s\"\.]+(?:from 0x[\w]+)?`
	r := regexp.MustCompile(pattern)

	// Find all matches in the given code
	matches := r.FindAllString(cadenceCode, -1)

	return matches
}

func ParseImport(line string) *flixkit.Network {
	// Define regex patterns
	pattern1 := `import "(?P<contract>[^"]+)"`
	pattern2 := `import (?P<contract>\w+) from (?P<address>0x[\w]+)`

	// Use regex to extract relevant information
	if matches, _ := regexpMatch(pattern1, line); matches != nil {
		return &flixkit.Network{
			Contract: matches["contract"],
		}
	} else if matches, _ := regexpMatch(pattern2, line); matches != nil {
		return &flixkit.Network{
			Contract:  matches["contract"],
			Address:   matches["address"],
			FqAddress: matches["address"], // Assuming FqAddress is the same as Address here
		}
	}

	return nil
}

func regexpMatch(pattern, text string) (map[string]string, error) {
	r := regexp.MustCompile(pattern)
	names := r.SubexpNames()
	match := r.FindStringSubmatch(text)
	if match == nil {
		return nil, nil
	}

	m := map[string]string{}
	for i, n := range match {
		m[names[i]] = n
	}

	return m, nil
}
