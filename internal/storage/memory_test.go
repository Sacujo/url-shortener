package storage_test

import (
	"errors"
	"sync"
	"testing"
	"url-shortener/internal/model"
	"url-shortener/internal/storage"
)

func TestMemoryStorage_SaveAndFindByID(t *testing.T) {
	s := storage.NewMemoryStorage()

	link := model.Link{ID: "abc123", URL: "https://example.com"}
	if err := s.Save(link); err != nil {
		t.Fatalf("unexpected error on Save: %v", err)
	}

	got, err := s.FindByID("abc123")
	if err != nil {
		t.Fatalf("unexpected error on FindByID: %v", err)
	}
	if got.URL != link.URL {
		t.Errorf("got URL %q, want %q", got.URL, link.URL)
	}
}

func TestMemoryStorage_SaveDuplicate(t *testing.T) {
	s := storage.NewMemoryStorage()
	link := model.Link{ID: "dup", URL: "https://example.com"}

	if err := s.Save(link); err != nil {
		t.Fatalf("unexpected error on first Save: %v", err)
	}

	err := s.Save(link)
	if !errors.Is(err, storage.ErrAlreadyExists) {
		t.Errorf("got error %v, want ErrAlreadyExists", err)
	}
}

func TestMemoryStorage_FindByID_NotFound(t *testing.T) {
	s := storage.NewMemoryStorage()

	_, err := s.FindByID("missing")
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("got error %v, want ErrNotFound", err)
	}
}

func TestMemoryStorage_IncrementClicks(t *testing.T) {
	s := storage.NewMemoryStorage()
	link := model.Link{ID: "click-me", URL: "https://example.com"}
	if err := s.Save(link); err != nil {
		t.Fatalf("unexpected error on Save: %v", err)
	}

	for i := 0; i < 3; i++ {
		if err := s.IncrementClicks("click-me"); err != nil {
			t.Fatalf("unexpected error on IncrementClicks: %v", err)
		}
	}

	got, err := s.FindByID("click-me")
	if err != nil {
		t.Fatalf("unexpected error on FindByID: %v", err)
	}
	if got.Clicks != 3 {
		t.Errorf("got Clicks=%d, want 3", got.Clicks)
	}
}

func TestMemoryStorage_IncrementClicks_NotFound(t *testing.T) {
	s := storage.NewMemoryStorage()

	err := s.IncrementClicks("missing")
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("got error %v, want ErrNotFound", err)
	}
}

func TestMemoryStorage_ConcurrentIncrementClicks(t *testing.T) {
	s := storage.NewMemoryStorage()
	link := model.Link{ID: "race-me", URL: "https://example.com"}
	if err := s.Save(link); err != nil {
		t.Fatalf("unexpected error on Save: %v", err)
	}

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if err := s.IncrementClicks("race-me"); err != nil {
				t.Errorf("unexpected error on IncrementClicks: %v", err)
			}
		}()
	}
	wg.Wait()

	got, err := s.FindByID("race-me")
	if err != nil {
		t.Fatalf("unexpected error on FindByID: %v", err)
	}
	if got.Clicks != goroutines {
		t.Errorf("got Clicks=%d, want %d (lost updates under concurrent access)", got.Clicks, goroutines)
	}
}
