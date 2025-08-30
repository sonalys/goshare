package identitycontroller

import (
	"github.com/sonalys/goshare/internal/ports"
	v1 "github.com/sonalys/goshare/pkg/v1"
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
