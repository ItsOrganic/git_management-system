package models

type GithubUser struct {
	Name   string `json:"name"`
	Id     int    `json:"id"`
	Avatar string `json:"avatar_url"`
}

type GitlabUser struct {
	Name   string `json:"name"`
	Id     int    `json:"id"`
	Avatar string `json:"avatar_url"`
}

type GitlabRepo struct {
	Name   string `json:"name"`
	ID     int    `json:"id"`
	Review bool   `json:"review"`
}

type GithubRepo struct {
	FullName string `json:"full_name"`
	Id       int    `json:"id" omitempty`
	Review   bool   `json:"review"`
}
