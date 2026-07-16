const BASE_URL = ''

export async function apiFetch(path, options = {}) {
  const { headers: extraHeaders, ...rest } = options
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { 'Content-Type': 'application/json', ...extraHeaders },
    credentials: 'same-origin',
    ...rest
  })
  if (!res.ok) throw new Error(`HTTP error: ${res.status}`)
  return res.json()
}
