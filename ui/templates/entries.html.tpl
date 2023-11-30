{{template "layout" .}}
{{define "content"}}
  <div class="page_heading">
    <h1>Entries</h1>
    <sl-input clearable id="filter" size="small" value="{{.Consistent.Pagination.Filter}}">
      <sl-icon name="search" slot="prefix"></sl-icon>
    </sl-input>
  </div>
  <div class="table full_page">
    <table>
      <thead>
      <tr>
        <th>Name</th>
        <th>NVRA</th>
        <th>Created</th>
        <th>State</th>
        <th>OS Release</th>
      </tr>
      </thead>
      <tbody>
      {{if not .Custom.Entries}}
        <tr>
          <td colspan="5">No entries found</td>
        </tr>
      {{end}}
      {{range .Custom.Entries}}
        <tr>
          <td><a href="/{{.Name}}">{{.Name}}</a></td>
          <td>{{.EntryId}}</td>
          <td>{{pbNaturalTime .CreateTime}}</td>
          <td>{{.State}}</td>
          <td>{{.OsRelease}}</td>
        </tr>
      {{end}}
      </tbody>
    </table>
    <div class="pages">
      <div class="container">
        {{template "table_range" .}}
        <div>{{.Consistent.Pagination.Range}}</div>
        <sl-button {{if .Consistent.Pagination.PrevQuery}}href="{{.Consistent.Pagination.PrevQuery}}"{{else}}disabled{{end}} size="small">
          <sl-icon name="chevron-left"></sl-icon>
        </sl-button>
        <sl-button {{if .Consistent.Pagination.NextQuery}}href="{{.Consistent.Pagination.NextQuery}}"{{else}}disabled{{end}} size="small">
          <sl-icon name="chevron-right"></sl-icon>
        </sl-button>
      </div>
    </div>
  </div>
  <script>
    const input = document.querySelector('sl-input#filter');

    input.addEventListener('sl-change', event => {
      const value = event.target.value;

      // Set the "q" query parameter to the input value
      const url = new URL(window.location);
      url.searchParams.set('q', window.isEmptyFilter(value, 'entryId:"{value}"'));
      location.search = url.search;
    });
  </script>
{{end}}