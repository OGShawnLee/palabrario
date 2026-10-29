import { ref } from 'vue'

const isCredentialsOpen = ref(false)

/** Estado global de superposiciones del shell (lo comparten TopBar, guardias del router y layouts). */
export function useOverlays() {
  return {
    isCredentialsOpen,
    openCredentials: () => {
      isCredentialsOpen.value = true
    },
    closeCredentials: () => {
      isCredentialsOpen.value = false
    },
  }
}
