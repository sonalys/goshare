package identitycontroller

import (
	"github.com/sonalys/goshare/internal/application"
	v1 "github.com/sonalys/goshare/internal/application/pkg/v1"
)

type (
	IdentityEncoder interface {
		Encode(identity *v1.Identity) (string, error)
	}

	Dependencies struct {
		application.Database
		IdentityEncoder
	}
)
