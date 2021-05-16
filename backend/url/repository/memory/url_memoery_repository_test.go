package memory_test

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/eyalch/shrtr/backend/domain"
	"github.com/eyalch/shrtr/backend/url/repository/memory"
)

func TestGet(t *testing.T) {
	r := memory.NewURLMemoryRepository()

	// Create the URL
	err := r.Create("abc123", "http://example.com")
	if err != nil {
		t.Fatal("could not create URL:", err)
	}

	url, err := r.Get("abc123")

	if err != nil {
		t.Fatal("could not get URL: ", err)
	} else if url != "http://example.com" {
		t.Fatalf("got wrong URL: %s, expected: %s", url, "http://example.com")
	}
}

func TestGet_NotFound(t *testing.T) {
	r := memory.NewURLMemoryRepository()

	_, err := r.Get("abc123")

	if errors.Cause(err) != domain.ErrKeyNotExists {
		t.Fatal("wrong error was returned for a non-existing key:", err)
	}
}

func TestCreate_Duplicate(t *testing.T) {
	r := memory.NewURLMemoryRepository()

	err := r.Create("abc123", "http://example.com")
	if err != nil {
		t.Fatal("could not create URL:", err)
	}

	err = r.Create("abc123", "http://another-example.com")

	if errors.Cause(err) != domain.ErrKeyAlreadyExists {
		t.Fatal("wrong error was returned for a duplicate key:", err)
	}
}
