package token

import "time"

type Config struct {
	PrivatePem string        `koanf:"private_pem"`
	PublicPem  string        `koanf:"public_pem"`
	Expiration time.Duration `koanf:"expiration"`
}
