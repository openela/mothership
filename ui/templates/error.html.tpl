{{template "layout" .}}
{{define "content"}}
  <div class="centered_page error_page">
    <img src="/static/oh_no.png" alt="Oh noooo gopher" />
    <h1>{{.Custom}}</h1>
  </div>
{{end}}