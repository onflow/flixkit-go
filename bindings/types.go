package bindings

import (
	"net/url"
)

type simpleParameter struct {
	Name        string
	JsType      string
	Description string
	FclType     string
	CadType     string
}

type templateData struct {
	FclVersion      string
	Version         string
	Parameters      []simpleParameter
	Title           string
	Description     string
	Location        string
	IsScript        bool
	IsLocalTemplate bool
}

type FclGenerator struct {
	Templates []string
}

type FlixParameter struct {
	Name string
	Type string
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func GetFlixFclCompatibility(flixVersion string) string {
	compatibleVersions := map[string]string{
		"1.0.0": "1.3.0",
		"1.1.0": "1.9.0",
		// add more versions here
	}
	v, ok := compatibleVersions[flixVersion]
	if !ok {
		// default to latest if flix version not configured
		return "1.9.0"
	}
	return v
}
