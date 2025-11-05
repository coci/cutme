package test

import (
	"errors"
	"testing"

	"github.com/coci/cutme/internal/core/domain"
	"github.com/coci/cutme/internal/infra/config"
	"github.com/coci/cutme/internal/services"
)

type fakeRepo struct {
	saved   []domain.Link
	saveErr error
	find    map[string]domain.Link
}

func (f *fakeRepo) Save(link domain.Link) error {
	f.saved = append(f.saved, link)
	return f.saveErr
}

func (f *fakeRepo) FindByCode(code string) domain.Link {
	if f.find == nil {
		return domain.Link{}
	}
	return f.find[code]
}

type seqIDGen struct {
	seq []int
	i   int
}

func (g *seqIDGen) NextID() int {
	if len(g.seq) == 0 {
		return 0
	}
	if g.i >= len(g.seq) {
		// keep returning the last element
		return g.seq[len(g.seq)-1]
	}
	v := g.seq[g.i]
	g.i++
	return v
}

func inAlphabet(s, alphabet string) bool {
	for _, r := range s {
		found := false
		for _, a := range alphabet {
			if r == a {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestShorten_Success_SavesLinkAndReturnsCode(t *testing.T) {
	cfg := config.Default()
	repo := &fakeRepo{}
	idg := &seqIDGen{seq: []int{123}}

	svc := services.NewShortenerService(repo, idg, cfg)

	url := "http://example.com/very/long/path?with=query"
	code, err := svc.Shorten(url)
	if err != nil {
		t.Fatalf("Shorten returned error: %v", err)
	}
	if code == "" {
		t.Fatalf("expected non-empty code")
	}

	// It should save exactly one link with correct fields
	if len(repo.saved) != 1 {
		t.Fatalf("expected 1 saved link, got %d", len(repo.saved))
	}
	saved := repo.saved[0]
	if saved.Link != url {
		t.Errorf("saved link url mismatch: got %q", saved.Link)
	}
	if saved.Code != code {
		t.Errorf("saved link code mismatch: got %q want %q", saved.Code, code)
	}
	if saved.CreatedAt != 0 { // clock() currently returns 0
		t.Errorf("expected CreatedAt to be 0, got %d", saved.CreatedAt)
	}

	// Check properties of generated code
	if len(code) < cfg.HashCfg.HashMinLength {
		t.Errorf("code length %d shorter than min %d", len(code), cfg.HashCfg.HashMinLength)
	}
	if !inAlphabet(code, cfg.HashCfg.HashAlphabet) {
		t.Errorf("code %q contains characters outside configured alphabet", code)
	}
}

func TestShorten_SaveError_Propagates(t *testing.T) {
	cfg := config.Default()
	repo := &fakeRepo{saveErr: errors.New("boom")}
	idg := &seqIDGen{seq: []int{1}}

	svc := services.NewShortenerService(repo, idg, cfg)

	_, err := svc.Shorten("http://example.com")
	if err == nil {
		t.Fatalf("expected error from Shorten when Save fails")
	}
}

func TestResolve_ReturnsRepoResult(t *testing.T) {
	cfg := config.Default()
	link := domain.Link{Code: "abc1234", Link: "http://golang.org"}
	repo := &fakeRepo{find: map[string]domain.Link{"abc1234": link}}
	idg := &seqIDGen{seq: []int{1}}

	svc := services.NewShortenerService(repo, idg, cfg)

	got, err := svc.Resolve("abc1234")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != link {
		t.Fatalf("resolve mismatch: got %+v want %+v", got, link)
	}
}

func TestResolve_NotFound_ReturnsZeroValueNoError(t *testing.T) {
	cfg := config.Default()
	repo := &fakeRepo{find: map[string]domain.Link{}}
	idg := &seqIDGen{seq: []int{1}}

	svc := services.NewShortenerService(repo, idg, cfg)

	got, err := svc.Resolve("does-not-exist")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != (domain.Link{}) {
		t.Fatalf("expected zero value link, got %+v", got)
	}
}

func TestMakeUniqueCode_DifferentIDsProduceDifferentCodes(t *testing.T) {
	cfg := config.Default()
	repo := &fakeRepo{}
	idg := &seqIDGen{seq: []int{1, 2}}

	svc := services.NewShortenerService(repo, idg, cfg)

	c1 := svc.MakeUniqueCode()
	c2 := svc.MakeUniqueCode()

	if c1 == c2 {
		t.Fatalf("expected different codes for different IDs, got %q == %q", c1, c2)
	}
}

func TestMakeUniqueCode_RespectsAlphabetAndMinLength(t *testing.T) {
	cfg := config.Default()
	repo := &fakeRepo{}
	idg := &seqIDGen{seq: []int{0}}

	svc := services.NewShortenerService(repo, idg, cfg)

	code := svc.MakeUniqueCode()
	if len(code) < cfg.HashCfg.HashMinLength {
		t.Errorf("code length %d shorter than min %d", len(code), cfg.HashCfg.HashMinLength)
	}
	if !inAlphabet(code, cfg.HashCfg.HashAlphabet) {
		t.Errorf("code %q contains characters outside configured alphabet", code)
	}
}

func TestMakeUniqueCode_NegativeID_YieldsEmptyCode(t *testing.T) {
	cfg := config.Default()
	repo := &fakeRepo{}
	idg := &seqIDGen{seq: []int{-5}}

	svc := services.NewShortenerService(repo, idg, cfg)

	code := svc.MakeUniqueCode()
	if code != "" {
		t.Fatalf("expected empty code for negative ID (encode error ignored), got %q", code)
	}
}

func TestShorten_AllowsEmptyURL_StillGeneratesCodeAndSaves(t *testing.T) {
	cfg := config.Default()
	repo := &fakeRepo{}
	idg := &seqIDGen{seq: []int{42}}

	svc := services.NewShortenerService(repo, idg, cfg)

	code, err := svc.Shorten("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code == "" {
		t.Fatalf("expected non-empty code even for empty URL input")
	}
	if len(repo.saved) != 1 {
		t.Fatalf("expected 1 saved link, got %d", len(repo.saved))
	}
	if repo.saved[0].Link != "" {
		t.Fatalf("expected saved link to keep empty URL, got %q", repo.saved[0].Link)
	}
}
