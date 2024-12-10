package models

type GithubUser struct {
	Name   string `json:"name" omitempty`
	Id     int    `json:"id" omitempty`
	Avatar string `json:"avatar_url" omitempty`
}

type GitlabUser struct {
	Name   string `json:"name" omitempty`
	Id     int    `json:"id" omitempty`
	Avatar string `json:"avatar_url" omitempty`
}

type GitlabRepo struct {
	Name   string `json:"name"`
	ID     int    `json:"id"`
	Review bool   `json:"review"`
}

type GithubRepo struct {
	FullName string `json:"full_name" omitempty`
	Id       int    `json:"id" omitempty`
	Review   bool   `json:"review"`
}
