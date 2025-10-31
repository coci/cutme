package cutme

import "github.com/coci/cutme/internal/cutme/domain"

type Shortener interface {
	Shorten(url string) (string, error)
	Resolve(code string) (string error)
	MakeUniqueCode() string
}

type LinkRepo interface {
	Save(link domain.Link) error
	FindByCode(code string) domain.Link
}

type IDGenerator interface {
	Next() int32
}
