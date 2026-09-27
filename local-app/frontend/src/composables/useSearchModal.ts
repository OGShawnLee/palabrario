import { ref } from "vue";
import { onKeyStroke } from "@vueuse/core";

// Estado reactivo global del SearchModal
const isSearchOpen = ref(false);

export function useSearchModal() {
  function openSearch() {
    isSearchOpen.value = true;
  }

  function closeSearch() {
    isSearchOpen.value = false;
  }

  function toggleSearch() {
    isSearchOpen.value = !isSearchOpen.value;
  }

  // Listener global de teclado con @vueuse/core para Ctrl + K y Cmd + K
  onKeyStroke(["k", "K"], (e) => {
    if (e.ctrlKey || e.metaKey) {
      e.preventDefault();
      toggleSearch();
    }
  });

  return {
    isSearchOpen,
    openSearch,
    closeSearch,
    toggleSearch,
  };
}
