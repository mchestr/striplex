package server

import "striplex/config"

func Init() error {
	r := NewRouter()
	r.SetTrustedProxies(config.Config.GetStringSlice("server.trusted_proxies"))
	return r.Run(config.Config.GetString("server.address"))
}
