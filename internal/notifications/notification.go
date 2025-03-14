package notifications

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/docker/docker/api/types/events"
)

type Event struct {
	Type            string
	Action          string
	Container       string
	Image           string
	Time            int64
	Status          string
	Labels          map[string]string
	Name            string
	Project         string
	Service         string
	ExitCode        string
	ExitCodeDetails string

	Message string
}

const textTpl = `{{if .Message}}{{.Message}}{{- else -}}
{{.Type}} {{ActionName .Action}} {{.Name}}
{{- if and .Project .Service }} {{.Project}}::{{.Service}}{{- end}}
{{- if .ExitCode }} Exit code: {{.ExitCode}}{{if .ExitCodeDetails}} ({{.ExitCodeDetails}}){{end}}{{- end}}{{end -}}
`

const mdTpl = `{{if .Message}}{{.Message}}{{- else -}}
{{.Type}} **{{ActionName .Action}}** {{WrapCode .Name}}
{{- if and .Project .Service }} {{WrapCode .Project}}::{{WrapCode .Service}}{{- end}}
{{if .ExitCode
-}}Exit code: {{WrapCode .ExitCode}}{{if .ExitCodeDetails}} ({{.ExitCodeDetails}}){{end}}{{-
end}}{{end -}}
`

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

const ansiTpl = `{{if .Message}}{{.Message}}{{- else -}}
{{.Type}} {{Yellow}}{{ActionName .Action}}{{Reset}} {{Cyan}}{{.Name}}{{Reset}}
{{- if and .Project .Service }} {{Blue}}{{.Project}}{{Reset}}::{{Magenta}}{{.Service}}{{Reset}}{{- end -}}
{{if .ExitCode
}} Exit code: {{if eq .ExitCode "0"}}{{Green}}{{.ExitCode}}{{Reset}}{{else}}{{Red}}{{.ExitCode}}{{Reset}}{{end
-}}{{if .ExitCodeDetails}} ({{.ExitCodeDetails}}){{end}}{{- end}}{{end -}}
`

var (
	textTemplate *template.Template
	mdTemplate   *template.Template
	ansiTemplate *template.Template
)

func init() {
	// Initialize templates during package initialization
	var err error

	funcMap := template.FuncMap{
		"ShortID": func(s string) string { return s[0:20] },
		"WrapCode": func(s string) string {
			return "`" + strings.ReplaceAll(s, "`", "'") + "`"
		},
		"ActionName": func(s string) string {
			if s == "die" {
				return "stop"
			}
			if s == "health_status: healthy" {
				return "healthy"
			}
			if s == "health_status: unhealthy" {
				return "unhealthy"
			}
			return s
		},
		// ansi colors
		"Red":     func() string { return Red },
		"Green":   func() string { return Green },
		"Yellow":  func() string { return Yellow },
		"Blue":    func() string { return Blue },
		"Magenta": func() string { return Magenta },
		"Cyan":    func() string { return Cyan },
		"Gray":    func() string { return Gray },
		"White":   func() string { return White },
		"Reset":   func() string { return Reset },
	}

	textTemplate, err = template.New("text").Funcs(funcMap).Parse(textTpl)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse text template: %v", err))
	}

	mdTemplate, err = template.New("md").Funcs(funcMap).Parse(mdTpl)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse markdown template: %v", err))
	}

	ansiTemplate, err = template.New("ansi").Funcs(funcMap).Parse(ansiTpl)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse ANSI template: %v", err))
	}
}

type EventActionMap map[string]map[string]bool
type ExitCodeMap map[string]string

const containerNameLabel = "name"
const dockerComposeProjectLabel = "com.docker.compose.project"
const dockerComposeServiceLabel = "com.docker.compose.service"
const execDurationLabel = "execDuration"
const exitCodeLabel = "exitCode"

var SupportedEvents = EventActionMap{
	"container": {
		"start":                    true,
		"die":                      true,
		"health_status: healthy":   true,
		"health_status: unhealthy": true,
	},
	"connection": {
		"message": true,
	},
}

var exitCodeMapping = ExitCodeMap{
	"0": "Success",
	// https://tldp.org/LDP/abs/html/exitcodes.html
	"1": "Application error",
	"2": "Misuse of builtin",
	// https://docs.docker.com/engine/containers/run/#exit-status
	"125": "Container failed to run",
	"126": "Container command cannot be invoked",
	"127": "Container command cannot be found",
	"128": "Invalid argument used on exit",

	// 128 + n Fatal
	"134": "Abnormal termination SIGABRT",
	"137": "Immediate termination SIGKILL",
	"139": "Segmentation Fault SIGSEGV",
	"143": "Graceful termination SIGTERM",

	"255": "Exit status out of range",
}

func NewEventFromDocker(msg events.Message) Event {
	labels := msg.Actor.Attributes

	name, _ := labels[containerNameLabel]
	project, _ := labels[dockerComposeProjectLabel]
	service, _ := labels[dockerComposeServiceLabel]
	exitCode, _ := labels[exitCodeLabel]

	return Event{
		Type:      string(msg.Type),
		Action:    string(msg.Action),
		Container: msg.Actor.ID,
		Image:     msg.Actor.Attributes["image"],
		Time:      msg.Time,
		Status:    msg.Status,
		Labels:    msg.Actor.Attributes,

		Name:            name,
		Project:         project,
		Service:         service,
		ExitCode:        exitCode,
		ExitCodeDetails: getExitCodeDetails(exitCode),
	}
}

func (e Event) ShouldNotify(debug bool) bool {
	supported := false
	if actionMap, ok := SupportedEvents[e.Type]; ok {
		supported, _ = actionMap[e.Action]
	}
	if !supported {
		if debug {
			fmt.Printf("Skipping unsupported event: %s:%s\n", e.Type, e.Action)
		}
		return false
	}

	return true
}

func getExitCodeDetails(exitCode string) string {
	if msg, ok := exitCodeMapping[exitCode]; ok {
		return msg
	}

	return ""
}

func (e *Event) Text() string {
	var buf bytes.Buffer
	err := textTemplate.Execute(&buf, e)
	if err != nil {
		fmt.Printf("Error generating template: %v\n", err)
		// cheap fallback
		return e.Type + " " + e.Action + " " + e.Name
	}

	return buf.String()
}

func (e *Event) Markdown() string {
	var buf bytes.Buffer
	err := mdTemplate.Execute(&buf, e)
	if err != nil {
		fmt.Printf("Error generating template: %v\n", err)
		// cheap fallback
		return e.Type + " " + e.Action + " " + e.Name
	}

	return buf.String()
}

func (e *Event) ANSI() string {
	var buf bytes.Buffer
	err := ansiTemplate.Execute(&buf, e)
	if err != nil {
		fmt.Printf("Error generating template: %v\n", err)
		// cheap fallback
		return e.Type + " " + e.Action + " " + e.Name
	}

	return buf.String()
}
