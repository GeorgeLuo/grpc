package main

import (
	"io"

	"github.com/GeorgeLuo/grpc/models"
	"github.com/olekukonko/tablewriter"
)

// TODO: handle batch requests

// Render writes to out stream a table with contents of a renderable object.
func Render(renderable models.Renderable, out io.Writer) {

	t := tablewriter.NewWriter(out)
	t.SetHeader(renderable.Headers())

	for _, row := range renderable.Rows() {
		t.Append(row)
	}
	t.Render()

}
