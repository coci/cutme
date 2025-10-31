package app

import (
	"fmt"

	"github.com/coci/cutme/internal/cutme"
	"github.com/coci/cutme/internal/cutme/config"
	"github.com/coci/cutme/internal/cutme/domain"
	"github.com/speps/go-hashids/v2"
)

type Service struct {
	cutme.Shortener

	linkRepo    cutme.LinkRepo
	idGenerator cutme.IDGenerator
	clock       func() int64
	config      config.Config
}

func NewService(repo cutme.LinkRepo, idGenerator cutme.IDGenerator, clock func() int64, config config.Config) *Service {
	return &Service{
		linkRepo:    repo,
		idGenerator: idGenerator,
		clock:       clock,
		config:      config,
	}
}

func (s Service) Shorten(url string) (string, error) {
	code := s.MakeUniqueCode()

	data := domain.Link{
		Code:      code,
		Link:      url,
		CreatedAt: s.clock(),
	}

	err := s.linkRepo.Save(data)

	if err != nil {
		return "", err
	}

	return code, nil
}

func (s Service) Resolve(code string) (string, error) {
	return "google.com", nil
}

func (s Service) MakeUniqueCode() string {
	newBaseCode := s.idGenerator.Next()

	hashId := hashids.HashIDData{
		Salt:      s.config.HashCfg.HashSalt,
		Alphabet:  s.config.HashCfg.HashAlphabet,
		MinLength: s.config.HashCfg.HashMinLength,
	}
	h, _ := hashids.NewWithData(&hashId)

	e, _ := h.Encode([]int{newBaseCode})

	fmt.Println(e)

	return e
}
