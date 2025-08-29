package secrets

import "os"

type Secrets struct {
	PostgresConn string
	JWTSignKey   []byte
}

func LoadSecrets() Secrets {
	secrets := Secrets{
		PostgresConn: os.Getenv("POSTGRES_CONN"),
		JWTSignKey:   []byte("my-super-secret-key"),
	}

	return secrets
}
