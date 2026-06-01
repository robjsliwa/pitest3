package model

import "testing"

func TestTodoValidateRejectsEmptyTitle(t *testing.T) {
	cases := []struct {
		name string
		todo Todo
	}{
		{
			name: "completely empty title",
			todo: Todo{},
		},
		{
			name: "whitespace-only title",
			todo: Todo{Title: "   "},
		},
		{
			name: "tab and newline only",
			todo: Todo{Title: "\t\n"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.todo.Validate()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestTodoValidateAcceptsValidTitle(t *testing.T) {
	cases := []struct {
		name  string
		title string
	}{
		{"simple title", "Buy milk"},
		{"title with leading space", "  Buy milk"},
		{"title with trailing space", "Buy milk  "},
		{"title with special chars", "Buy #milk & eggs!"},
		{"with valid due_date", "Buy milk"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Todo{Title: tc.title}.Validate()
			if err != nil {
				t.Fatalf("expected nil, got: %v", err)
			}
		})
	}
}

func TestTodoValidateRejectsInvalidDueDate(t *testing.T) {
	cases := []struct {
		name    string
		dueDate string
	}{
		{"free text", "banana"},
		{"wrong format", "01-15-2025"},
		{"partial date", "2025-01"},
		{"slash separated", "2025/01/15"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Todo{Title: "Valid", DueDate: tc.dueDate}.Validate()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestTodoValidateAcceptsValidDueDate(t *testing.T) {
	err := Todo{Title: "Valid", DueDate: "2025-01-20"}.Validate()
	if err != nil {
		t.Fatalf("expected nil, got: %v", err)
	}
}

func TestTodoValidateAcceptsEmptyDueDate(t *testing.T) {
	err := Todo{Title: "Valid", DueDate: ""}.Validate()
	if err != nil {
		t.Fatalf("expected nil for empty due_date, got: %v", err)
	}
}
