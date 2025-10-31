package repositories

import "github.com/coci/cutme/internal/core/domain"

type LinkRepository struct {
}

func NewLinkRepository() *LinkRepository {
	return &LinkRepository{}
}

func (l LinkRepository) Save(link domain.Link) error {
	//TODO implement me
	panic("implement me")
}

func (l LinkRepository) FindByCode(code string) domain.Link {
	//TODO implement me
	panic("implement me")
}
