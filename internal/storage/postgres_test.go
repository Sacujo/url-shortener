package storage_test

import (
	"database/sql"
	"errors"
	"os"
	"sync"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"url-shortener/internal/model"
	"url-shortener/internal/storage"
)

func newTestPostgresStorage(t *testing.T) (*storage.PostgresStorage, string) {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not set, skipping PostgreSQL integration tests")
	}
	s, err := storage.NewPostgresStorage(dsn)
	if err != nil {
		t.Skipf("could not connect to PostgreSQL: %v", err)
	}
	return s, dsn
}

func cleanupLink(t *testing.T, dsn, id string) {
	t.Helper()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("cleanup: failed to open db: %v", err)
	}
	defer db.Close()
	if _, err := db.Exec(`DELETE FROM links WHERE alias = $1`, id); err != nil {
		t.Fatalf("cleanup: failed to delete link %q: %v", id, err)
	}
}

func TestPostgresStorage_SaveAndFindByID(t *testing.T) {
	s, dsn := newTestPostgresStorage(t)
	const id = "pg-test-save"
	t.Cleanup(func() { cleanupLink(t, dsn, id) })

	link := model.Link{ID: id, URL: "https://example.com"}
	if err := s.Save(link); err != nil {
		t.Fatalf("unexpected error on Save: %v", err)
	}

	got, err := s.FindByID(id)
	if err != nil {
		t.Fatalf("unexpected error on FindByID: %v", err)
	}
	if got.URL != link.URL {
		t.Errorf("got URL %q, want %q", got.URL, link.URL)
	}
	if got.CreatedAt.IsZero() {
		t.Errorf("expected CreatedAt to be set by the database default")
	}
}

func TestPostgresStorage_SaveDuplicate(t *testing.T) {
	s, dsn := newTestPostgresStorage(t)
	const id = "pg-test-dup"
	t.Cleanup(func() { cleanupLink(t, dsn, id) })

	link := model.Link{ID: id, URL: "https://example.com"}
	if err := s.Save(link); err != nil {
		t.Fatalf("unexpected error on first Save: %v", err)
	}

	err := s.Save(link)
	if !errors.Is(err, storage.ErrAlreadyExists) {
		t.Errorf("got error %v, want ErrAlreadyExists", err)
	}
}

func TestPostgresStorage_FindByID_NotFound(t *testing.T) {
	s, _ := newTestPostgresStorage(t)

	_, err := s.FindByID("pg-test-missing")
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("got error %v, want ErrNotFound", err)
	}
}

func TestPostgresStorage_IncrementClicks(t *testing.T) {
	s, dsn := newTestPostgresStorage(t)
	const id = "pg-test-increment"
	t.Cleanup(func() { cleanupLink(t, dsn, id) })

	if err := s.Save(model.Link{ID: id, URL: "https://example.com"}); err != nil {
		t.Fatalf("unexpected error on Save: %v", err)
	}

	for i := 0; i < 3; i++ {
		if err := s.IncrementClicks(id); err != nil {
			t.Fatalf("unexpected error on IncrementClicks: %v", err)
		}
	}

	got, err := s.FindByID(id)
	if err != nil {
		t.Fatalf("unexpected error on FindByID: %v", err)
	}
	if got.Clicks != 3 {
		t.Errorf("got Clicks=%d, want 3", got.Clicks)
	}
}

func TestPostgresStorage_IncrementClicks_NotFound(t *testing.T) {
	s, _ := newTestPostgresStorage(t)

	err := s.IncrementClicks("pg-test-missing")
	if !errors.Is(err, storage.ErrNotFound) {
		t.Errorf("got error %v, want ErrNotFound", err)
	}
}

func TestPostgresStorage_ConcurrentIncrementClicks(t *testing.T) {
	s, dsn := newTestPostgresStorage(t)
	const id = "pg-test-race"
	t.Cleanup(func() { cleanupLink(t, dsn, id) })

	if err := s.Save(model.Link{ID: id, URL: "https://example.com"}); err != nil {
		t.Fatalf("unexpected error on Save: %v", err)
	}

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if err := s.IncrementClicks(id); err != nil {
				t.Errorf("unexpected error on IncrementClicks: %v", err)
			}
		}()
	}
	wg.Wait()

	got, err := s.FindByID(id)
	if err != nil {
		t.Fatalf("unexpected error on FindByID: %v", err)
	}
	if got.Clicks != goroutines {
		t.Errorf("got Clicks=%d, want %d", got.Clicks, goroutines)
	}
}
