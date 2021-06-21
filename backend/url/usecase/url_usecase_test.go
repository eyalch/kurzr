package usecase_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"

	"github.com/eyalch/shrtr/backend/domain"
	"github.com/eyalch/shrtr/backend/url/repository/memory"
	"github.com/eyalch/shrtr/backend/url/usecase"
)

func TestGetURL(t *testing.T) {
	repo := memory.NewURLMemoryRepository()
	uc := usecase.NewURLUsecase(repo, nil)

	_, err := uc.GetURL("abc123")

	if errors.Cause(err) != domain.ErrKeyNotExists {
		t.Fatal("wrong error was returned for a non-existing key:", err)
	}

	err = repo.Create("abc123", "http://example.com")
	if err != nil {
		t.Fatal("could not create URL:", err)
	}

	url, err := uc.GetURL("abc123")

	if err != nil {
		t.Fatal("could not get URL:", err)
	} else if url != "http://example.com" {
		t.Fatalf("got wrong URL: %s, expected: %s", url, "http://example.com")
	}
}

type testKeyGenerator struct {
	counter int
}

// To test the logic which retries the key generation upon getting a duplicate
// key, we implement the key generator so that the key will change once per 2
// generations, which will produce this sequence:
// abc-0, abc-0, abc-2, abc-2, abc-4, ...
func (kg *testKeyGenerator) GenerateKey() string {
	suffix := kg.counter - kg.counter%2
	kg.counter++
	return fmt.Sprintf("abc-%d", suffix)
}

func TestShortenURL(t *testing.T) {
	repo := memory.NewURLMemoryRepository()
	uc := usecase.NewURLUsecase(repo, &testKeyGenerator{0})

	key, err := uc.ShortenURL("http://example.com")

	if err != nil {
		t.Fatal("could not shorten URL:", err)
	} else if key != "abc-0" {
		t.Fatalf("got wrong key: %s, expected: %s", key, "abc-0")
	}

	key, err = uc.ShortenURL("http://example.com")

	if err != nil {
		t.Fatal("could not shorten URL:", err)
	} else if key != "abc-2" {
		t.Fatalf("got wrong key: %s, expected: %s", key, "abc-2")
	}
}
