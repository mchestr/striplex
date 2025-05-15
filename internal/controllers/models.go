package controllers

type InfoResponse struct {
	RequestsURL      string   `json:"requests_url"`
	DiscordServerUrl string   `json:"discord_server_url"`
	ServerName       string   `json:"server_name"`
	Features         []string `json:"features"`
}
