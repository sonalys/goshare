package usercontroller

import "github.com/sonalys/goshare/internal/ports"

type (
	Dependencies struct {
		ports.LocalDatabase
	}
)
