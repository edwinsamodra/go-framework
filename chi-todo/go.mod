module github.com/edwinsamodra/go-framework/chi-todo

go 1.26.2

require (
	github.com/go-chi/chi/v5 v5.2.5
	github.com/edwinsamodra/go-framework/common v0.0.0
)

replace github.com/edwinsamodra/go-framework/common => ../common

require filippo.io/edwards25519 v1.1.0 // indirect
