package controllers

type (
	Controller struct {
		*Ledgers
		*Users
	}
)

func New(dep Dependencies) *Controller {
	return &Controller{
		Ledgers: &Ledgers{
			db: dep.Database,
		},
		Users: &Users{
			identityEncoder: dep.IdentityEncoder,
			db:              dep.Database,
		},
	}
}
