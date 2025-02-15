package contract

import "go_link_reducer/types"

type URLRepository interface {
	Create(types.CreateURLPayload) (types.URL, error)
	GetAll() ([]types.URL, error)
	GetOne(URL string) (types.URL, error)
	Delete(URL string) error
}
