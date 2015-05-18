package drschollz

// DefaultTmpl is the default email template
var DefaultTmpl = `

{{.AppName}}

-------------------------------
Error:
-------------------------------

	{{.Err}}

-------------------------------
Info:
-------------------------------

	- {{.Time}}
{{range $e := .Extras}}
	- {{$e}}
{{end}}

-------------------------------
Backtrace:
-------------------------------

{{.BackTrace}}
`
