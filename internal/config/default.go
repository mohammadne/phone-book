package config

import (
	"github.com/MohammadNE/PhoneBook/pkg/logger"
)

func Default() *Config {
	return &Config{
		Logger: &logger.Config{
			Development: true,
			Level:       "debug",
			Encoding:    "console",
		},
		// Grpc: &grpc.Config{
		// 	ListenPort: 9090,
		// },
		// Token: &token.Config{
		// 	PrivatePem: "-----BEGIN PRIVATE KEY-----\n" +
		// 		"MC4CAQAwBQYDK2VwBCIEINyMNS8h9M9HO73Tg1BPr53p//qlqylO+wPKN8GrlsX7\n" +
		// 		"-----END PRIVATE KEY-----",
		// 	PublicPem: "-----BEGIN PUBLIC KEY-----\n" +
		// 		"MCowBQYDK2VwAyEAqQsZ5iRNP3kdpNn3V/db9o/WkYHY8kkwQqCZGcDvJ+g=\n" +
		// 		"-----END PUBLIC KEY-----",
		// 	Expiration: 30 * time.Minute,
		// },
		// Crypto: &crypto.Config{
		// 	Secret: "7w!z%C*F-JaNdRgU",
		// 	Salt:   "YGp*OfH^za",
		// },
	}
}
