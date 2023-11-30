{{define "table_range"}}
  <sl-dropdown id="table_range">
    <sl-button size="small" slot="trigger" caret>Rows per page: {{.Consistent.Pagination.PageSize}}</sl-button>
    <sl-menu>
      <sl-menu-item>25</sl-menu-item>
      <sl-menu-item>50</sl-menu-item>
      <sl-menu-item>100</sl-menu-item>
    </sl-menu>
  </sl-dropdown>
  <script>
    const dropdown = document.querySelector('sl-dropdown#table_range');

    dropdown.addEventListener('sl-select', event => {
      const selectedItem = event.detail.item;

      // Replace the "size" query param with the new value
      const url = new URL(window.location.href);
      url.searchParams.set('offset', '0');
      url.searchParams.set('size', selectedItem.textContent);

      // Navigate to the new URL
      location.search = url.search;
    });
  </script>
{{end}}