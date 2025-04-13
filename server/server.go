package server

import "striplex/config"

func Init() {
	r := NewRouter()
	r.SetTrustedProxies(config.Config.GetStringSlice("server.trusted_proxies"))
	r.Run(config.Config.GetString("server.address"))
}
