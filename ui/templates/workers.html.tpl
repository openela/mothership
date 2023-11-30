{{template "layout" .}}
{{define "content"}}
  <sl-dialog label="New worker" class="dialog">
    <form action="/workers" method="POST">
      {{.Consistent.Csrf}}
      <sl-input autofocus name="worker_id" label="Worker ID" help-text="This ID will be used to uniquely identify this worker."></sl-input>
      <sl-button type="submit" slot="footer" variant="primary">Create</sl-button>
    </form>
  </sl-dialog>
  <div class="page_heading">
    <h1>Workers</h1>
    <div class="actions">
      <sl-button id="new_worker" size="small" variant="primary">New worker</sl-button>
      <sl-input clearable id="filter" size="small" value="{{.Consistent.Pagination.Filter}}">
        <sl-icon name="search" slot="prefix"></sl-icon>
      </sl-input>
    </div>
  </div>
  <div class="table full_page">
    <table>
      <thead>
      <tr>
        <th>Name</th>
        <th>Worker ID</th>
        <th>Created</th>
      </tr>
      </thead>
      <tbody>
      {{if not .Custom.Workers}}
        <tr>
          <td colspan="5">No workers found</td>
        </tr>
      {{end}}
      {{range .Custom.Workers}}
        <tr>
          <td><a href="/{{.Name}}">{{.Name}}</a></td>
          <td>{{.WorkerId}}</td>
          <td>{{pbNaturalTime .CreateTime}}</td>
        </tr>
      {{end}}
      </tbody>
    </table>
    <div class="pages">
      <div class="container">
          {{template "table_range" .}}
        <div>{{.Consistent.Pagination.Range}}</div>
        <sl-button
                {{if .Consistent.Pagination.PrevQuery}}href="{{.Consistent.Pagination.PrevQuery}}"
                {{else}}disabled{{end}} size="small">
          <sl-icon name="chevron-left"></sl-icon>
        </sl-button>
        <sl-button
                {{if .Consistent.Pagination.NextQuery}}href="{{.Consistent.Pagination.NextQuery}}"
                {{else}}disabled{{end}} size="small">
          <sl-icon name="chevron-right"></sl-icon>
        </sl-button>
      </div>
    </div>
  </div>
  <script>
    const filterInput = document.querySelector('sl-input#filter');

    filterInput.addEventListener('sl-change', event => {
      const value = event.target.value;

      // Set the "q" query parameter to the input value
      const url = new URL(window.location);
      url.searchParams.set('q', window.isEmptyFilter(value, 'workerId:"{value}"'));
      location.search = url.search;
    });

    const dialog = document.querySelector('.dialog');
    const input = dialog.querySelector('sl-input');
    const openButton = document.querySelector("#new_worker");

    openButton.addEventListener('click', () => dialog.show());
  </script>
{{end}}