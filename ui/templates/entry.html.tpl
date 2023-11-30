{{template "layout" .}}
{{define "content"}}
  <div class="container">
    <div class="resource">
      <div class="title">
        <div class="left">
          <h1>{{.Custom.EntryId}}</h1>
            {{$variant := "success"}}
            {{if or (eq .Custom.State 1) (eq .Custom.State 6)}}
                {{$variant = "primary"}}
            {{else if or (or (eq .Custom.State 5) (eq .Custom.State 7)) (eq .Custom.State 3)}}
                {{$variant = "danger"}}
            {{end}}
          <sl-tag size="small" variant="{{$variant}}" pill>{{.Custom.State}}</sl-tag>
          <sl-tag size="small" variant="primary" pill>{{.Custom.OsRelease}}</sl-tag>
        </div>
        <div class="right">
            {{if and (eq .Custom.State 3) (.Consistent.User)}}
              <form action="/{{.Custom.Name}}/rescue" method="POST">
                  {{.Consistent.Csrf}}
                <sl-button type="submit" size="small" variant="primary">Rescue</sl-button>
              </form>
            {{end}}
            {{if and (eq .Custom.State 2) (.Consistent.User)}}
              <form action="/{{.Custom.Name}}/retract" method="POST">
                  {{.Consistent.Csrf}}
                <sl-button type="submit" size="small" variant="error">Retract</sl-button>
              </form>
            {{end}}
        </div>
      </div>
      <div class="info">
        <b>Imported </b> {{pbNaturalTime .Custom.CreateTime}}
      </div>
      <div class="content">
        <div>
          <h3>src.rpm SHA256 Checksum</h3>
          <p>{{.Custom.Sha256Sum}}</p>
        </div>
        <div>
          <h3>Fetched by worker</h3>
          <p>{{.Custom.WorkerId.Value}}</p>
        </div>
          {{if .Custom.CommitUri}}
            <div>
              <h3>Commit URI</h3>
              <p><a href="{{.Custom.CommitUri}}">{{.Custom.CommitUri}}</a></p>
            </div>
          {{end}}
        <div>
          <h3>Commit Hash</h3>
          <p>{{.Custom.CommitHash}}</p>
        </div>
      </div>
    </div>
      {{if eq .Custom.State 3}}
          <code>{{.Custom.ErrorMessage}}</code>
      {{end}}
  </div>
{{end}}