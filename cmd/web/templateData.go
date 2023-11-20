package main

import (
	"artchernov.ru/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}
