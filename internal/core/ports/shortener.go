package ports

import "github.com/coci/cutme/internal/core/domain"

type Shortener interface {
	Shorten(url string) (string, error)
	Resolve(code string) (domain.Link, error)
	MakeUniqueCode() string
}
