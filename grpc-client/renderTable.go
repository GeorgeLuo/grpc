package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/GeorgeLuo/grpc/models"
	"github.com/olekukonko/tablewriter"
)

// Render writes to out stream a table with contents of a renderable object.
func Render(out io.Writer, renderables []Renderable, tabled bool) {

	if !tabled {
		for i, renderable := range renderables {
			out.Write([]byte(renderable.Raw()))
			if i != len(renderables) {
				out.Write([]byte("\n"))
			}
		}
		return
	}

	for _, renderable := range renderables {
		t := tablewriter.NewWriter(out)
		t.SetHeader(renderable.Headers())

		// format table
		t.SetRowLine(true)
		t.SetRowSeparator("-")

		t.AppendBulk(renderable.Rows())
		if renderable.Title() != "" {
			out.Write([]byte(renderable.Title() + "\n"))
		}
		t.Render()
	}
}

// Renderable is an interface to support operation to display contents in a
// tablewriter table.
type Renderable interface {
	Title() string
	Headers() []string
	Rows() [][]string
	Raw() string
}

// RenderableStatusResponse a renderable wrapper for models.StatusResponse
type RenderableStatusResponse struct {
	statusResponse *models.StatusResponse
	title          string
}

// NewRenderableStatusResponse is used to return a status response renderable
// wrapper with empty title.
func NewRenderableStatusResponse(r *models.StatusResponse) *RenderableStatusResponse {
	return &RenderableStatusResponse{
		statusResponse: r,
	}
}

// Raw returns the raw json of a status response.
func (r *RenderableStatusResponse) Raw() string {
	s, err := json.Marshal(r.statusResponse)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(s)
}

// Title returns the title of the renderable table.
func (r *RenderableStatusResponse) Title() string {
	return r.title
}

// Headers returns the headers to populate a table of status response fields.
func (r *RenderableStatusResponse) Headers() []string {
	return []string{"task_id", "start_time", "end_time",
		"exit_code", "exec_error"}
}

// Rows produces a row of data for the data returned by status response.
func (r *RenderableStatusResponse) Rows() [][]string {
	var endTime string
	if r.statusResponse.EndTime == nil {
		endTime = ""
	} else {
		endTime = r.statusResponse.EndTime.String()
	}

	var exitCode string
	if r.statusResponse.ExitCode == nil {
		exitCode = ""
	} else {
		exitCode = strconv.Itoa(*r.statusResponse.ExitCode)
	}

	return [][]string{
		{r.statusResponse.TaskID,
			r.statusResponse.StartTime.String(),
			endTime, exitCode,
			r.statusResponse.ExecError},
	}
}

// RenderableStartResponse a renderable wrapper for models.StartResponse
type RenderableStartResponse struct {
	startResponse *models.StartResponse
	title         string
}

// NewRenderableStartResponse is used to return a start response renderable
// wrapper with empty title.
func NewRenderableStartResponse(r *models.StartResponse) *RenderableStartResponse {
	return &RenderableStartResponse{
		startResponse: r,
	}
}

// Raw returns the raw json of a start response.
func (r *RenderableStartResponse) Raw() string {
	s, err := json.Marshal(r.startResponse)
	if err != nil {
		return ""
	}
	return string(s)
}

// Headers returns the headers to display start response data.
func (r *RenderableStartResponse) Headers() []string {
	return []string{"task_id"}
}

// Rows produces a row of data for task_id.
func (r *RenderableStartResponse) Rows() [][]string {
	return [][]string{
		{r.startResponse.TaskID},
	}
}

// Title returns the title of the renderable table.
func (r *RenderableStartResponse) Title() string {
	return r.title
}

// RenderableStopResponse a renderable wrapper for models.StopResponse
type RenderableStopResponse struct {
	stopResponse *models.StopResponse
	title        string
}

// NewRenderableStopResponse is used to return a stop response renderable
// wrapper with empty title.
func NewRenderableStopResponse(r *models.StopResponse) *RenderableStopResponse {
	return &RenderableStopResponse{
		stopResponse: r,
	}
}

// Raw returns the raw json of a stop response.
func (r *RenderableStopResponse) Raw() string {
	s, err := json.Marshal(r.stopResponse)
	if err != nil {
		return ""
	}
	return string(s)
}

// Headers returns the headers to populate a table of stop response fields.
func (r *RenderableStopResponse) Headers() []string {
	return []string{"task_id", "exit_code"}
}

// Rows produces a row of data for the data returned by stop response.
func (r *RenderableStopResponse) Rows() [][]string {
	return [][]string{
		{r.stopResponse.TaskID, strconv.Itoa(*r.stopResponse.ExitCode)},
	}
}

// Title returns the title of the renderable table.
func (r *RenderableStopResponse) Title() string {
	return r.title
}
