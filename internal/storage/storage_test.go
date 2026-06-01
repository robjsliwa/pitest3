package storage

import (
	"github.com/robjsliwa/pitest3/internal/model"
	"testing"
	"time"
)

func TestStoreCreateReturnsIncrementingIDs(t *testing.T) {
	store := NewStore("")

	todos := make([]model.Todo, 3)
	for i := 0; i < 3; i++ {
		todo := model.Todo{
			Title:     "Todo " + string(rune('A'+i)),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		var err error
		todos[i], err = store.Create(todo)
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	// IDs should start at 1 and increment monotonically
	expectedIDs := []int{1, 2, 3}
	for i, expected := range expectedIDs {
		if todos[i].ID != expected {
			t.Errorf("todo[%d].ID = %d, want %d", i, todos[i].ID, expected)
		}
	}
}

func TestStoreCreateRejectsInvalidTodo(t *testing.T) {
	store := NewStore("")

	_, err := store.Create(model.Todo{})
	if err == nil {
		t.Fatal("expected error for empty title, got nil")
	}
}

func TestStoreGetByID(t *testing.T) {
	store := NewStore("")

	created, err := store.Create(model.Todo{
		Title:     "Buy milk",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := store.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}

	if found.ID != created.ID || found.Title != created.Title {
		t.Errorf("GetByID returned wrong todo: %+v, want %+v", found, created)
	}
}

func TestStoreGetByIDNotFound(t *testing.T) {
	store := NewStore("")

	_, err := store.GetByID(999)
	if err == nil {
		t.Fatal("expected error for non-existent ID, got nil")
	}
}

func TestStoreList(t *testing.T) {
	store := NewStore("")

	for i := 0; i < 3; i++ {
		_, err := store.Create(model.Todo{
			Title:     "Todo " + string(rune('A'+i)),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	list := store.List()
	if len(list) != 3 {
		t.Errorf("List returned %d todos, want 3", len(list))
	}
}
