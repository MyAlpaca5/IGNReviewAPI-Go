package models

type SystemInfo struct {
	Env     string `json:"env"`
	Version string `json:"version"`
}

type Healthcheck struct {
	Status     string     `json:"status"`
	SystemInfo SystemInfo `json:"system_info"`
}
