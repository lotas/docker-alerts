package notifications

import (
	"bytes"
	"fmt"
	"text/template"
)

type Notification struct {
	Message         string
	Type            string
	Action          string
	Container       string
	Image           string
	Name            string
	Project         string
	Service         string
	ExitCode        string
	ExitCodeDetails string
}

const textTpl = `
{{.Type}} {{.Action}}
{{.Name}} [{{.Container}}]
{{.Project}} {{.Service}}
{{if .ExitCode}}
Exit with: {{.ExitCode}}{{if .ExitCodeDetails}} ({{.ExitCodeDetails}}){{end}}
{{end}}
`

const backtick = "`"

const mdTpl = `
{{.Type}} **{{.Action}}**
` + backtick + `{{.Name}}` + backtick + ` [{{.Container}}]
` + backtick + `{{.Project}}` + backtick + ` ` + backtick + `{{.Service}}` + backtick + `
{{if .ExitCode}}
Exit with: ` + backtick + `{{.ExitCode}}` + backtick + `{{if .ExitCodeDetails}} ({{.ExitCodeDetails}}){{end}}
{{end}}
`

var (
	textTemplate *template.Template
	mdTemplate   *template.Template
)

func init() {
	// Initialize templates during package initialization
	var err error

	textTemplate, err = template.New("text").Parse(textTpl)
	if err != nil {
		fmt.Errorf("Failed to parse text template: %v", err)
	}

	mdTemplate, err = template.New("md").Parse(mdTpl)
	if err != nil {
		fmt.Errorf("Failed to parse markdown template: %v", err)
	}
}

func (n *Notification) Text() string {
	if n.Message != "" {
		return n.Message
	}

	var buf bytes.Buffer
	err := textTemplate.Execute(&buf, n)
	if err != nil {
		//debug
		fmt.Errorf("Error generating template %v", err.Error())
		// cheap fallback
		return n.Type + " " + n.Action + " " + n.Name
	}

	return buf.String()
}

func (n *Notification) Markdown() string {
	if n.Message != "" {
		return n.Message
	}

	var buf bytes.Buffer
	err := mdTemplate.Execute(&buf, n)
	if err != nil {
		//debug
		fmt.Errorf("Error generating template %v", err.Error())
		// cheap fallback
		return n.Type + " " + n.Action + " " + n.Name
	}

	return buf.String()
}
