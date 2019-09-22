package main

import (
	"io"

	"github.com/GeorgeLuo/grpc/models"
	"github.com/olekukonko/tablewriter"
)

// Render writes to out stream a table with contents of a renderable object.
func Render(renderable models.Renderable, out io.Writer) {

	t := tablewriter.NewWriter(out)
	t.SetHeader(renderable.Headers())

	// format table
	t.SetRowLine(true)
	t.SetRowSeparator("-")

	t.AppendBulk(renderable.Rows())

	t.Render()
}
