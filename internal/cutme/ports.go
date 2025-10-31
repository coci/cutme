package cutme

import "github.com/coci/cutme/internal/cutme/domain"

type IShortener interface {
	Shorten(url string) (string, error)
	Resolve(code string) (string error)
}

type ILinkRepo interface {
	Save(link domain.Link) error
	FindByCode(code string) domain.Link
}

type IdGenerator interface {
	NewCode() int32
}
