package main

type FluxConfig struct {
	Iam      string      `toml:"iam"`
	Url      string      `toml:"url"`
	Port     int         `toml:"port"`
	Logdir   string      `toml:"logdir"`
	Balancer FluxCluster `toml:"cluster"`
	//Jwts     []JwtAuth   `toml:"jwt"`
}

type JwtAuth struct {
	RequiredClaims   []JwtClaim `toml:"claim"`
	DecryptionSecret string     `toml:"secret"`
}

type JwtClaim struct {
	Key   string `toml:"key"`
	Value string `toml:"value"`
}

type FluxCluster struct {
	Name            string `toml:"name"`
	BalancerAddress string `toml:"address"`
	BalancerPort    int    `toml:"port"`
	Scramble        bool   `toml:"scramble"`
}
