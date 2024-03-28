package internal

const (
	DefaultFlixServerUrl = "https://flix.flow.com/v1/templates"
)

type FlixQueryTypes string

const (
	FlixName     FlixQueryTypes = "name"
	FlixFilePath FlixQueryTypes = "filePath"
	FlixPath     FlixQueryTypes = "path"
	FlixId       FlixQueryTypes = "id"
	FlixUrl      FlixQueryTypes = "url"
	FlixJson     FlixQueryTypes = "json"
)

type VerCheck struct {
	FVersion string `json:"f_version"`
}
