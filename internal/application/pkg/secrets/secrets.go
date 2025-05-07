package secrets

type Secrets struct {
	PostgresConn string
	JWTSignKey   []byte
}

func LoadSecrets() Secrets {
	secrets := Secrets{
		PostgresConn: "postgresql://user:password@localhost:5432/goshare?sslmode=disable",
		JWTSignKey:   []byte("my-super-secret-key"),
	}

	return secrets
}
