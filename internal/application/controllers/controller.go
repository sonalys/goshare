package controllers

type (
	Controller struct {
		*Ledgers
		*Users
	}
)

func New(dep Dependencies) *Controller {
	subscriber := newSubscriber()

	return &Controller{
		Ledgers: &Ledgers{
			subscriber: subscriber,
			db:         dep.Database,
		},
		Users: &Users{
			identityEncoder: dep.IdentityEncoder,
			subscriber:      subscriber,
			db:              dep.Database,
		},
	}
}
