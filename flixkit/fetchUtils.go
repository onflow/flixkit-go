package flixkit

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"

	v1 "github.com/onflow/flixkit-go/internal/v1"
	"github.com/onflow/flixkit-go/internal/v1_1"
)

func getType(s string, f FileReader) FlixQueryTypes {
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
	return false
}

func isFilePath(path string, f FileReader) bool {
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

func getFlixByFilePath(path string, reader FileReader) (FlixInterface, string, error) {
	rawFlix, err := reader.ReadFile(path)
	if err != nil {
		return nil, path, fmt.Errorf("could not read flix file %s: %w", path, err)
	}

	flixVer, err := getTemplateVersion(rawFlix)
	if err != nil {
		return nil, path, fmt.Errorf("unable to determine flix schema version for %s: %w", path, err)
	}

	var flix FlixInterface
	var parseErr error
	switch flixVer {
	case "1.0.0":
		flix, parseErr = v1.ParseJSON(rawFlix)
	case "1.1.0":
		flix, parseErr = v1_1.ParseJSON(rawFlix)
	}

	if parseErr != nil {
		return nil, path, fmt.Errorf("could not parse flix from file %s: %w", path, err)
	}

	return flix, path, parseErr
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
