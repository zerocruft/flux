package main

type FluxConfig struct {
	This FluxNode `toml:"this"`
	Cluster []FluxNode `toml:"node"`
	Authentication JwtAuth `toml:"jwt"`
	LogDir string `toml:"logs"`
}


type FluxNode struct {
	Id string `toml:"id"`
	Address string `toml:"address"`
	Port int `toml:"port"`
	Enabled bool `toml:"enabled"`
}

type JwtAuth struct {
	RequiredClaims []JwtClaim `toml:"claim"`
	DecryptionSecret string `toml:"secret"`
}

type JwtClaim struct {
	Key string `toml:"key"`
	Value string `toml:"value"`
}
