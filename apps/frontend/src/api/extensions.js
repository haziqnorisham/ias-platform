import { apiFetch } from './index'

export const getExtensions = () => apiFetch('/api/extensions', {
  method: 'POST'
})

export async function getExtensionWidgets(name) {
  try {
    const resp = await fetch(`/api/extensions/${name}/health`, {
      credentials: 'same-origin',
    })
    if (!resp.ok) return []
    const data = await resp.json()
    return data.widgets || []
  } catch {
    return []
  }
}

export function loadExtensionComponent(name) {
  return new Promise((resolve, reject) => {
    const existing = document.querySelector(`script[data-extension="${name}"]`)
    if (existing) {
      resolve()
      return
    }
    const script = document.createElement('script')
    script.src = `/api/extensions/${name}/component.js`
    script.setAttribute('data-extension', name)
    script.onload = resolve
    script.onerror = () => reject(new Error(`Failed to load extension component for "${name}"`))
    document.head.appendChild(script)
  })
}
