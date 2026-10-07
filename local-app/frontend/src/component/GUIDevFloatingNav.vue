<template>
  <!-- ===================================================================== -->
  <!-- NAVEGADOR FLOTANTE DE DESARROLLO (DEV SWITCHER DOCK CON ROUTER-LINK)  -->
  <!-- ===================================================================== -->
  <aside
    class="fixed bottom-6 right-6 z-50 bg-chocolate text-white p-2 rounded-full shadow-2xl border-2 border-border flex items-center gap-1 backdrop-blur-md text-xs font-satoshi"
  >
    <span
      class="text-[10px] font-black uppercase tracking-widest pl-3 pr-1 text-border select-none"
    >
      Dev
    </span>
    <router-link
      to="/"
      class="px-3 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider transition-all cursor-pointer text-border hover:bg-white/10"
      active-class="!bg-brand !text-white shadow-md"
      exact
    >
      Inicio
    </router-link>
    <router-link
      to="/dictionary"
      class="px-3 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider transition-all cursor-pointer text-border hover:bg-white/10"
      active-class="!bg-brand !text-white shadow-md"
    >
      Diccionario
    </router-link>
    <router-link
      to="/studio"
      class="px-3 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider transition-all cursor-pointer text-border hover:bg-white/10"
      active-class="!bg-brand !text-white shadow-md"
    >
      Studio
    </router-link>
    <router-link
      to="/play"
      class="px-3 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider transition-all cursor-pointer text-border hover:bg-white/10"
      active-class="!bg-brand !text-white shadow-md"
    >
      Arena
    </router-link>
    <router-link
      to="/editor"
      class="px-3 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider transition-all cursor-pointer text-border hover:bg-white/10"
      active-class="!bg-brand !text-white shadow-md"
    >
      Editorial
    </router-link>
    <router-link
      to="/signin"
      class="px-2.5 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider transition-all cursor-pointer text-border hover:bg-white/10"
      active-class="!bg-brand !text-white shadow-md"
      title="Iniciar Sesión"
    >
      Login
    </router-link>
    <router-link
      to="/signup"
      class="px-2.5 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider transition-all cursor-pointer text-border hover:bg-white/10"
      active-class="!bg-brand !text-white shadow-md"
      title="Crear Cuenta"
    >
      Registro
    </router-link>
    <router-link
      to="/design-language"
      class="px-3 py-1.5 rounded-full text-xs font-bold uppercase tracking-wider transition-all cursor-pointer flex items-center gap-1 text-border hover:bg-white/10"
      active-class="!bg-brand !text-white shadow-md"
    >
      <span>Tokens</span>
      <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"></span>
    </router-link>
  </aside>
</template>
