<template>
  <!-- ===================================================================== -->
  <!-- NAVEGADOR FLOTANTE DE DESARROLLO (DEV SWITCHER DOCK CON ROUTER-LINK)  -->
  <!-- ===================================================================== -->
  <aside
    class="fixed bottom-6 right-6 z-50 bg-chocolate text-white p-2 rounded-full shadow-2xl border-2 border-border flex items-center gap-1.5 backdrop-blur-md"
  >
    <span
      class="text-[10px] font-black uppercase tracking-widest pl-3 pr-1 text-border select-none"
    >
      Dev Mode
    </span>
    <router-link
      to="/"
      class="px-4 py-2 rounded-full text-xs font-black uppercase tracking-wider transition-all cursor-pointer text-border hover:bg-white/10"
      active-class="!bg-brand !text-white shadow-md"
      exact
    >
      Home Page
    </router-link>
    <router-link
      to="/dictionary"
      class="px-4 py-2 rounded-full text-xs font-black uppercase tracking-wider transition-all cursor-pointer text-border hover:bg-white/10"
      active-class="!bg-brand !text-white shadow-md"
    >
      Diccionario 📖
    </router-link>
    <router-link
      to="/design-language"
      class="px-4 py-2 rounded-full text-xs font-black uppercase tracking-wider transition-all cursor-pointer flex items-center gap-1.5 text-border hover:bg-white/10"
      active-class="!bg-brand !text-white shadow-md"
    >
      <span>Design System</span>
      <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
    </router-link>
  </aside>
</template>
