const BASE_URL = ''

export async function apiFetch(path, options = {}) {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { 'Content-Type': 'application/json', ...options.headers },
    credentials: 'same-origin',
    ...options
  })
  if (!res.ok) throw new Error(`HTTP error: ${res.status}`)
  return res.json()
}
