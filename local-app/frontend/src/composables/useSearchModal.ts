import { ref } from "vue";
import { onKeyStroke } from "@vueuse/core";
import router from "@/router";

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
      // Si estamos en la ruta de Diccionario, delegamos al buscador principal de la vista
      if (router.currentRoute.value?.path === "/dictionary") {
        window.dispatchEvent(
          new CustomEvent("palabrario:focus-dictionary-search")
        );
        return;
      }
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
