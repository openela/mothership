{{define "layout"}}
  <!DOCTYPE html>
  <html lang="en">
  <head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">

    <link rel="icon" type="image/png" href="/static/favicon.png"/>

    <style>
      :not(:defined) {
        visibility: hidden;
      }
    </style>

    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:wght@300;400;500;600;700&display=swap">

    <link rel="stylesheet" href="/static/dist/light.css?ver={{ctx "version"}}"/>
    <link rel="stylesheet" href="/static/dist/styles.css?ver={{ctx "version"}}"/>
    <script src="/static/dist/index.js?ver={{ctx "version"}}"></script>

    <title>{{ctx "instanceName"}}</title>
  </head>
  <body>
  <div class="nav">
    <div class="logo">
      <a href="/">
        <!--<img alt="Mship Logo" src="/static/mothership.svg"/>-->
          {{ctx "instanceName"}}
      </a>
    </div>
    <div class="divider"></div>
    <div class="links">
        {{range $href, $text := (ctx "links")}}
            {{$isActive := false}}
            {{if hasPrefix $.Consistent.Request.URL.Path $href}}
                {{$isActive = true}}
            {{end}}
            <a href="{{$href}}" {{if $isActive}}class="active"{{end}}>{{$text}}</a>
        {{end}}
        {{if .Consistent.User}}
            {{range $href, $text := (ctx "authLinks")}}
                {{$isActive := false}}
                {{if hasPrefix $.Consistent.Request.URL.Path $href}}
                    {{$isActive = true}}
                {{end}}
                <a href="{{$href}}" {{if $isActive}}class="active"{{end}}>{{$text}}</a>
            {{end}}
        {{end}}
    </div>
    <div class="right_links">
        {{if .Consistent.User}}
          <a>{{.Consistent.User.Email}}</a>
          <a href="/auth/logout">Logout</a>
        {{else}}
          <a href="/auth/login">Login</a>
        {{end}}
    </div>
  </div>
  <div class="content">
      {{if .Consistent.Alerts}}
          {{range $alert := .Consistent.Alerts}}
            <sl-alert class="server-alert" variant="{{$alert.Variant}}" duration="2000">
                {{if $alert.Icon}}
                  <sl-icon slot="icon" name="{{$alert.Icon}}"></sl-icon>
                {{end}}
                {{if $alert.Subtitle}}
                  <strong>{{$alert.Title}}</strong><br/>
                {{else}}
                    {{$alert.Title}}
                {{end}}
                {{$alert.Subtitle}}
            </sl-alert>
          {{end}}
      {{end}}
      {{template "content" .}}
  </div>
  <div class="footer">
    Copyright &copy; 2023 Mustafa Gezen, and Ctrl IQ, Inc.
    <div class="right">
      Version: {{ctx "version"}}
    </div>
  </div>
  <script>
    // Auto toast server alerts
    const toasts = document.querySelectorAll('sl-alert.server-alert');
    toasts.forEach(toast => {
      toast.toast().then();
    });
  </script>
  </body>
{{end}}