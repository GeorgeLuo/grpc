package main

// BatchRenderable is a renderable for multiple status requests.
type BatchRenderable struct {
	Renderables []Renderable
	rows        [][]string
	header      []string
	title       string
}

// NewBatchRenderable is used to return an empty renderable collection.
func NewBatchRenderable(t string) *BatchRenderable {
	return &BatchRenderable{
		Renderables: []Renderable{},
		title:       t,
	}
}

// Headers returns the headers of the renderable.
func (b *BatchRenderable) Headers() []string {
	if len(b.Renderables) > 0 {
		return b.Renderables[0].Headers()
	}
	return b.header
}

// Rows produces a row of data for the data returned by the batch.
func (b *BatchRenderable) Rows() [][]string {

	for _, renderable := range b.Renderables {
		b.rows = append(b.rows, renderable.Rows()...)
	}

	return b.rows
}

// Title is usually the alias, summary of the table contents.
func (b *BatchRenderable) Title() string {
	return b.title
}

// AddRow adds a renderable to the batch.
func (b *BatchRenderable) AddRow(r Renderable) {
	b.Renderables = append(b.Renderables, r)
}
