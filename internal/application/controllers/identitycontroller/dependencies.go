package identitycontroller

import (
	"github.com/sonalys/goshare/internal/application"
	"github.com/sonalys/goshare/internal/ports"
)

type (
	IdentityEncoder interface {
		Encode(identity *application.Identity) (string, error)
	}

	Dependencies struct {
		ports.LocalDatabase
		IdentityEncoder
	}
)
