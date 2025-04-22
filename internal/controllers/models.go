package controllers

type InfoResponse struct {
	RequestsURL      string `json:"requests_url,omitempty"`
	DiscordServerUrl string `json:"discord_server_url,omitempty"`
	ServerName       string `json:"server_name,omitempty"`
}
