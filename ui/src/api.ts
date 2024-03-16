const API_BASE_URL = '/ui/api';
const ADMIN_API_BASE_URL = '/ui/admin-api';

export function fetchAPI<T>(url: string, options?: RequestInit): Promise<T> {
  return fetch(`${API_BASE_URL}${url}`, options)
    .then((res) => res.json())
    .then((data) => {
      if (data.code) {
        throw data;
      }
      return data;
    });
}

export function fetchAdminAPI<T>(
  url: string,
  options?: RequestInit,
): Promise<T> {
  return fetch(`${ADMIN_API_BASE_URL}${url}`, options).then((res) =>
    res.json(),
  );
}
