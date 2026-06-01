# Glossary

## Todo

A discrete unit of work. A Todo has:

- **Title** — required, free-text, the primary identifier when scanning.
- **Description** — optional, longer-form text for additional context.
- **Completed** — boolean status (not done / done).
- **Created At** — timestamp when the Todo was created.
- **Updated At** — timestamp of the last modification.
- **Completed At** — timestamp when the Todo was marked done (nil while incomplete).
- **Due Date** — optional date by which the Todo should be completed.

A Todo does not have categories, tags, or priorities in V1.

## Task

Not used — we say "Todo" to avoid collision with the Go language concept of a task/goroutine.
