package internal

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	v1 "github.com/onflow/flixkit-go/internal/v1"
	"github.com/onflow/flixkit-go/internal/v1_1"
	"github.com/onflow/flixkit-go/types"
)

func GetType(s string, f types.FileReader) FlixQueryTypes {
	switch {
	case isFlixPath(s):
		return FlixPath
	case isFilePath(s, f):
		return FlixFilePath
	case isHex(s):
		return FlixId
	case isUrl(s):
		return FlixUrl
	case isJson(s):
		return FlixJson
	default:
		return FlixName
	}
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func isHex(str string) bool {
	if len(str) != 64 {
		return false
	}
	_, err := hex.DecodeString(str)
	return err == nil
}

func isFlixPath(str string) bool {
	if _, err := processFlixPath(str); err != nil {
		return false
	}
	return true
}

func isFilePath(path string, f types.FileReader) bool {
	if f == nil {
		return false
	}
	_, err := f.ReadFile(path)
	return err == nil
}

func isJson(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func getFlixByFilePath(path string, reader types.FileReader) (types.FlixInterface, string, error) {
	rawFlix, err := reader.ReadFile(path)
	if err != nil {
		return nil, path, fmt.Errorf("could not read flix file %s: %w", path, err)
	}

	return processRawFlix(rawFlix, path)
}

// TODO:
func getFlixByFlixPath(path string) (types.FlixInterface, string, error) {
	_, err := processFlixPath(path)
	if err != nil {
		return nil, path, err
	}

	return nil, path, fmt.Errorf("this method is not implemented yet")
}

// TODO:
func getFlixById(id string) (types.FlixInterface, string, error) {
	return nil, id, fmt.Errorf("this method is not implemented yet")
}

// TODO:
func getFlixByName(name string) (types.FlixInterface, string, error) {
	return nil, name, fmt.Errorf("this method is not implemented yet")
}

// TODO:
func getFlixByUrl(url string) (types.FlixInterface, string, error) {
	return nil, url, fmt.Errorf("this method is not implemented yet")
}

func parseFlixJSON(flixJSON string) (types.FlixInterface, string, error) {
	return processRawFlix([]byte(flixJSON), flixJSON)
}

func processRawFlix(rawFlix []byte, source string) (types.FlixInterface, string, error) {
	flixVer, err := getTemplateVersion(rawFlix)
	if err != nil {
		return nil, source, fmt.Errorf("unable to determine flix schema version for %s: %w", source, err)
	}

	var flix types.FlixInterface
	var parseErr error
	switch flixVer {
	case "1.0.0":
		flix, parseErr = v1.ParseJSON(rawFlix)
	case "1.1.0":
		flix, parseErr = v1_1.ParseJSON(rawFlix)
	}

	if parseErr != nil {
		return nil, source, fmt.Errorf("could not parse flix from file %s: %w", source, err)
	}

	return flix, source, parseErr
}

func processFlixPath(path string) ([]string, error) {
	sections := strings.Split(path, "\\")

	if len(sections) != 4 {
		return nil, fmt.Errorf("invalid flix path")
	}

	return sections, nil
}

func getTemplateVersion(templateBytes []byte) (string, error) {
	var verCheck VerCheck

	err := json.Unmarshal([]byte(templateBytes), &verCheck)
	if err != nil {
		return "", err
	}

	if verCheck.FVersion == "" {
		return "", fmt.Errorf("version not found")
	}

	return verCheck.FVersion, nil
}
