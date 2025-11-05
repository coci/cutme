package services

import (
	"fmt"
	"log"

	"github.com/coci/cutme/internal/core/domain"
	"github.com/coci/cutme/internal/core/ports"
	"github.com/coci/cutme/internal/infra/config"
	"github.com/speps/go-hashids/v2"
)

type ShortenerService struct {
	ports.Shortener

	linkRepo    ports.LinkRepo
	idGenerator ports.IDGenerator
	config      *config.Config
}

func NewShortenerService(linkRepo ports.LinkRepo, idGenerator ports.IDGenerator, config *config.Config) *ShortenerService {
	return &ShortenerService{
		linkRepo:    linkRepo,
		idGenerator: idGenerator,
		config:      config,
	}
}

func (s ShortenerService) Shorten(url string) (string, error) {
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

func (s ShortenerService) Resolve(code string) (domain.Link, error) {
	result := s.linkRepo.FindByCode(code)

	return result, nil
}

func (s ShortenerService) MakeUniqueCode() string {
	newBaseCode := s.idGenerator.NextID()

	hashId := hashids.HashIDData{
		Salt:      s.config.HashCfg.HashSalt,
		Alphabet:  s.config.HashCfg.HashAlphabet,
		MinLength: s.config.HashCfg.HashMinLength,
	}
	h, err := hashids.NewWithData(&hashId)

	if err != nil {
		log.Fatal("instantiating hashids", err)
	}

	e, _ := h.Encode([]int{newBaseCode})

	fmt.Println(e)

	return e
}

func (s ShortenerService) clock() int64 {
	return 0
}
