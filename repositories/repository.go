package repositories

import "github.com/pocketbase/pocketbase"

type Repository struct {
	PB *pocketbase.PocketBase
}

func NewRepository(
	PB *pocketbase.PocketBase,
) Repository {
	return Repository{
		PB: PB,
	}
}
