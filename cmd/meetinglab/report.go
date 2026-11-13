package main

import "html/template"

var page = template.Must(template.New("report").Parse(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Meeting follow-up register</title>
<style>
:root { background: #f1f5f4; color: #162d34; font: 17px/1.55 system-ui, sans-serif; }
main { max-width: 900px; margin: auto; padding: 36px 22px; }
h1 { font-size: clamp(2.2rem, 6vw, 3.7rem); line-height: 1.06; margin: .2em 0; }
h2 { font-size: 1.25rem; }
a { color: #115e67; }
article { background: white; border: 1px solid #b6cdca; border-radius: 8px; padding: 14px 18px; }
section { margin: 1.4rem 0; overflow-wrap: anywhere; }
li { margin-bottom: .55em; }
small { color: #37565e; }
</style>
</head>
<body><main>
<h1>{{.Title}}</h1>
<p>Meeting follow-up register as of <time>{{.AsOf}}</time></p>
<section aria-label="Actions">
<h2>Actions</h2>
{{range .Actions}}<article>
<strong>{{.Summary}}</strong>
<p>{{.Status}} · {{if .Owner}}{{.Owner}}{{else}}Unassigned{{end}}{{if .Due}} · Due {{.Due}}{{end}}<br><small>{{.Source}}</small></p>
{{if .Evidence}}<p>Completion evidence · {{.Evidence}}</p>{{end}}
{{if .Findings}}<ul>{{range .Findings}}<li>{{.}}</li>{{end}}</ul>{{end}}
</article>{{else}}<p>No actions recorded.</p>{{end}}
</section>
<section aria-label="Decisions"><h2>Decisions</h2>{{template "entries" .Decisions}}</section>
<section aria-label="Questions"><h2>Questions</h2>{{template "entries" .Questions}}</section>
<section aria-label="Notes"><h2>Notes</h2>{{template "entries" .Notes}}</section>
<p><a href="https://www.rfc-editor.org/info/rfc21/">Official RFC 21</a> · <a href="https://github.com/Sreenivas-Sadhu-Prabhakara/rfc-0021-meeting-lab">Go source and run instructions</a></p>
</main></body></html>
{{define "entries"}}{{if .}}<ul>{{range .}}<li>{{.Summary}} <small>({{.Source}})</small></li>{{end}}</ul>{{else}}<p>None recorded in this input.</p>{{end}}{{end}}
`))
