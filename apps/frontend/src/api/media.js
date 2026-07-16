import { apiFetch } from './index'

const MEDIA_BASE = '/media'

function mediaFetch(path, options = {}) {
  return apiFetch(`${MEDIA_BASE}${path}`, options)
}

export const getMediaHealth = () => mediaFetch('/health')

export const getCameras = () => mediaFetch('/api/devices')

export const getCamera = (id) => mediaFetch(`/api/devices/${id}`)

export const createCamera = (body) => mediaFetch('/api/devices', {
  method: 'POST',
  body: JSON.stringify(body)
})

export const updateCamera = (id, body) => mediaFetch(`/api/devices/${id}`, {
  method: 'PUT',
  body: JSON.stringify(body)
})

export const deleteCamera = (id) => mediaFetch(`/api/devices/${id}`, {
  method: 'DELETE'
})

export const getCameraProfiles = (id) => mediaFetch(`/api/devices/${id}/profiles`)

export const setStreamProfile = (id, token) => mediaFetch(`/api/devices/${id}/stream-profile`, {
  method: 'PUT',
  body: JSON.stringify({ token })
})

export const getStreams = () => mediaFetch('/api/streams')

export const startStream = (deviceId) => mediaFetch(`/api/streams/${deviceId}/start`, {
  method: 'POST'
})

export const stopStream = (deviceId) => mediaFetch(`/api/streams/${deviceId}/stop`, {
  method: 'POST'
})
