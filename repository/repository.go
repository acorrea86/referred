package repository

import "blumer-ms-refers/contracts"

type Repository struct {
	session contracts.NeoSession
}

func NewRepository(session contracts.NeoSession) *Repository {
	return &Repository{
		session: session,
	}
}
