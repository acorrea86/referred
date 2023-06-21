package graph

import "blumer-ms-refers/repository"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repository *repository.Repository
}

func NewResolver(repository *repository.Repository) *Resolver {
	return &Resolver{
		Repository: repository,
	}
}
