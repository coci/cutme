package ports

import "github.com/coci/cutme/internal/core/domain"

type LinkRepo interface {
	Save(link domain.Link) error
	FindByCode(code string) domain.Link
}
