package model

type Project struct {
	Version string `json:"version"`
}

type ProjectV1 struct {
	Version     string   `json:"version"`
	TeamName    string   `json:"teamName"`
	ProjectName string   `json:"projectName"`
	Endpoint    string   `json:"endpoint"`
	Routes      []string `json:"routes"`
	Services    []string `json:"services"`
}
