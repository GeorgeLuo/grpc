package main

import (
	"encoding/json"
	"io"

	"github.com/GeorgeLuo/grpc/models"
	"github.com/olekukonko/tablewriter"
)

// table is used to display the response to the client in a human
// readable form. This is primarily useful for displaying the status
// of multiple processes at once.

// TableDisplay encapsulates an object that can be displayed using
// tablewriter.
type TableDisplay struct {
	out    io.Writer
	header []string
	data   [][]string
}

// NewTableDisplay is used to an empty TableDisplay.
func NewTableDisplay(base []string, output io.Writer) TableDisplay {
	return TableDisplay{
		out:    output,
		header: base,
		data:   [][]string{},
	}
}

// AddRow adds input row from Tableble to data if it complies
// with the current headers.
func (table *TableDisplay) AddRow(row models.Tableble) {

	dataRow := make([]string, len(table.header))
	dataMap := row.Data()

	for i, h := range table.header {
		dataRow[i] = dataMap[h]
	}

	table.data = append(table.data, dataRow)
}

// Render prints the contents of a table using tablewriter
func (table *TableDisplay) Render() {

	t := tablewriter.NewWriter(table.out)
	t.SetHeader(table.header)

	for _, row := range table.data {
		t.Append(row)
	}
	t.Render()
}

// PrintTable takes a mapping of string to a string and returns
// a pretty table of the data. PrintTable is not intended to be ran
// frequently (per change) but instead to build a fully populated dataset.
func PrintTable(arrayKey *string, data *[]byte,
	out io.Writer) (*tablewriter.Table, error) {

	rows := [][]string{}
	header := []string{}

	// this is a batch table
	if arrayKey != nil {
		// err := populateBatchRowsAndHeaders(&rows, &header, *arrayKey, data)
		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
	} else {

		dataMap := make(map[string]string)
		error := json.Unmarshal(*data, &dataMap)
		if error != nil {
			return nil, error
		}

		i := 0
		rowData := make([]string, len(dataMap))

		// populate keySet
		for key, val := range dataMap {
			header = append(header, key)
			rowData[i] = val
			i++
		}

		rows = append(rows, rowData)
	}

	table := tablewriter.NewWriter(out)
	table.SetHeader(header)

	for _, row := range rows {
		table.Append(row)
	}

	table.Render()
	return table, nil
}

// NewTable initializes a tablewriter with appropriate
func NewTable(header []string, out io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(out)
	table.SetHeader(header)
	return table
}

// PrintTableble writes tableble contents in a human readable form
// using tablewriter.
func PrintTableble(tableble models.Tableble,
	out io.Writer) {

	table := tablewriter.NewWriter(out)
	// table.SetHeader(tableble.Header())

	// for _, row := range tableble.Data() {
	// table.Append(row)
	// }

	table.Render()
}
