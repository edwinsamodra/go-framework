module todo-comparison/chi-todo

go 1.26.2

require (
	github.com/go-chi/chi/v5 v5.2.5
	example.com/todo-comparison/common v0.0.0
)

replace example.com/todo-comparison/common => ../common

require filippo.io/edwards25519 v1.1.0 // indirect
