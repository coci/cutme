package app

import (
	"github.com/coci/cutme/internal/cutme"
)

type IdGenerator struct {
	cutme.HashIDRepo

	idRepo cutme.HashIDRepo
}

func NewIdGenerator(idRepo cutme.HashIDRepo) *IdGenerator {
	return &IdGenerator{
		idRepo: idRepo,
	}
}

func (i IdGenerator) Next() int {
	// TODO : implement me
	return 14_000_000
}
