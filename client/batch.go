package main

import (
	"github.com/GeorgeLuo/grpc/models"
)

// BatchStatusRenderable is a renderable for multiple status requests.
type BatchStatusRenderable struct {
	StatusResponses []models.StatusResponse
	rows            [][]string
	header          []string
}

// NewBatchStatusRenderable is used to return an empty SyncMap.
func NewBatchStatusRenderable() *BatchStatusRenderable {
	return &BatchStatusRenderable{
		StatusResponses: []models.StatusResponse{},
	}
}

// Headers returns the headers to populate a table of status responses.
func (b *BatchStatusRenderable) Headers() []string {
	if len(b.StatusResponses) > 0 {
		return b.StatusResponses[0].Headers()
	}
	return b.header
}

// Rows produces a row of data for the data returned by status responses.
func (b *BatchStatusRenderable) Rows() [][]string {

	for _, statusResponse := range b.StatusResponses {
		b.rows = append(b.rows, statusResponse.Rows()...)
	}

	return b.rows
}

// AddResponse adds a status response object to the batch.
func (b *BatchStatusRenderable) AddResponse(r models.StatusResponse) {
	b.StatusResponses = append(b.StatusResponses, r)
}
