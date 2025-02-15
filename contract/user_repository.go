package contract

import "go_link_reducer/types"

type UserRepository interface {
	Register(user types.RegisterPayload) (string, error)
	GetUserByEmail(email string) (types.User, error)
	GetUserByID(ID int) (types.User, error)
}
