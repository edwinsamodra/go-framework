module todo-comparison/native-todo

go 1.26.2

require example.com/todo-comparison/common v0.0.0

replace example.com/todo-comparison/common => ../common

require filippo.io/edwards25519 v1.1.0 // indirect
