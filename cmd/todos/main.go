package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/robjsliwa/pitest3/internal/model"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <command> [flags]\n", os.Args[0])
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		doCreate(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func doCreate(args []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	title := fs.String("title", "", "Todo title (required)")
	description := fs.String("description", "", "Todo description (optional)")
	due := fs.String("due", "", "Due date YYYY-MM-DD (optional)")
	jsonOut := fs.Bool("json", false, "Output as JSON")
	fs.Parse(args)

	if strings.TrimSpace(*title) == "" {
		fmt.Fprintln(os.Stderr, "error: --title is required")
		os.Exit(1)
	}

	url := getenv("TODOS_API_URL", "http://localhost:8080")
	todo, err := createTodo(url, *title, *description, *due)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if *jsonOut {
		fmt.Println(formatJSON(todo))
	} else {
		fmt.Println(formatTodo(todo))
	}
}

func createTodo(url, title, description, due string) (model.Todo, error) {
	body := fmt.Sprintf(`{"title":%q,"description":%q,"due_date":%q}`, title, description, due)

	resp, err := http.Post(url+"/todos", "application/json", strings.NewReader(body))
	if err != nil {
		return model.Todo{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return model.Todo{}, fmt.Errorf("server returned %d", resp.StatusCode)
	}

	var result struct {
		Todo model.Todo `json:"todo"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return model.Todo{}, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Todo, nil
}

func formatTodo(t model.Todo) string {
	lines := []string{
		fmt.Sprintf("ID:          %d", t.ID),
		fmt.Sprintf("Title:       %s", t.Title),
		fmt.Sprintf("Description: %s", t.Description),
		fmt.Sprintf("Completed:   %v", t.Completed),
		fmt.Sprintf("Due:         %s", t.DueDate),
		fmt.Sprintf("Created:     %s", t.CreatedAt.Format("2006-01-02 15:04:05")),
	}
	return strings.Join(lines, "\n")
}

func formatJSON(t model.Todo) string {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Sprintf("error formatting todo: %v", err)
	}
	return string(data)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
