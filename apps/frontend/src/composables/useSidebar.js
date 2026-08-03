import { ref } from 'vue'

const sidebarVisible = ref(true)

export function useSidebar() {
  function toggleSidebar() {
    sidebarVisible.value = !sidebarVisible.value
  }

  return {
    sidebarVisible,
    toggleSidebar,
  }
}
