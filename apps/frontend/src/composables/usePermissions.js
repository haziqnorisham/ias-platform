import { ref, watch } from 'vue'

const STORAGE_KEY = 'ias_viewer_permissions'

// Sidebar pages that can be toggled for viewer users (mirrors SideNav navItems).
export const permissionPages = [
  { label: 'Home', to: '/' },
  { label: 'IAS AI (Preview)', to: '/ai' },
  { label: 'Dashboards', to: '/dashboards' },
  { label: 'Devices', to: '/devices' },
  { label: 'Device Profiles', to: '/device-profiles' },
  { label: 'Data Browser', to: '/data-browser' },
  { label: 'Ingest Logs', to: '/ingest-logs' },
  { label: 'Settings', to: '/settings' },
  { label: 'Extensions', to: '/extensions' },
  { label: 'Diagnostics', to: '/diagnostics' },
  { label: 'About', to: '/about' },
]

function defaults() {
  return Object.fromEntries(permissionPages.map(p => [p.to, true]))
}

function load() {
  try {
    const saved = JSON.parse(localStorage.getItem(STORAGE_KEY) || 'null')
    return { ...defaults(), ...(saved || {}) }
  } catch {
    return defaults()
  }
}

const viewerPermissions = ref(load())

watch(viewerPermissions, (val) => {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(val))
}, { deep: true })

export function usePermissions() {
  function setViewerPermission(to, allowed) {
    viewerPermissions.value[to] = allowed
  }

  // Admins always see everything; viewers are limited to allowed pages.
  function isNavItemVisible(to, isAdmin) {
    if (isAdmin) return true
    return viewerPermissions.value[to] !== false
  }

  return {
    viewerPermissions,
    setViewerPermission,
    isNavItemVisible,
  }
}
