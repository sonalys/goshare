package identitycontroller

import (
	v1 "github.com/sonalys/goshare/internal/application/v1"
	"github.com/sonalys/goshare/internal/ports"
)

type (
	IdentityEncoder interface {
		Encode(identity *v1.Identity) (string, error)
	}

	Dependencies struct {
		ports.LocalDatabase
		IdentityEncoder
	}
)
