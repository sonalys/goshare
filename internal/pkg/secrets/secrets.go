package secrets

type Secrets struct {
	PostgresConn string
}

func LoadSecrets() Secrets {
	secrets := Secrets{
		PostgresConn: "postgresql://user:password@postgres:5432/goshare?sslmode=disable",
	}

	return secrets
}
