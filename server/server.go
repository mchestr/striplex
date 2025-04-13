package server

import "striplex/config"

func Init() {
	config := config.GetConfig()
	r := NewRouter()
	r.SetTrustedProxies(config.GetStringSlice("server.trusted_proxies"))
	r.Run(config.GetString("server.address"))
}
