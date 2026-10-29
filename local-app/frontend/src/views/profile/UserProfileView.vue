<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { MediaDropzone } from '@palabrario/ui'
import type { PublicUserProfile, Region } from '@palabrario/ui'
import { api } from '@/api/client'
import { useAuth } from '@/composables/useAuth'
import { errorMessage } from '@/utils/errors'
import { fileToBase64, resolveMediaUrl } from '@/utils/files'

const route = useRoute()
const { refreshProfile } = useAuth()

const profile = ref<PublicUserProfile | null>(null)
const regions = ref<Region[]>([])
const isLoading = ref(true)
const loadError = ref<string | null>(null)

const editForm = reactive({ name: '', anchor_dialect: '' })
const avatarFile = ref<File | null>(null)
const profileFeedback = ref<{ ok: boolean; text: string } | null>(null)
const isSaving = ref(false)

const passwordForm = reactive({ current_password: '', new_password: '' })
const passwordFeedback = ref<{ ok: boolean; text: string } | null>(null)
const isChangingPassword = ref(false)

async function load(username: string) {
  isLoading.value = true
  loadError.value = null
  try {
    profile.value = await api.GetPublicProfileByUsername(username)
    if (profile.value?.is_owner) {
      editForm.name = profile.value.name
      editForm.anchor_dialect = profile.value.anchor_dialect
      regions.value = (await api.ListRegions()) ?? []
    }
  } catch (err) {
    profile.value = null
    loadError.value = errorMessage(err, 'No pudimos cargar este perfil.')
  } finally {
    isLoading.value = false
  }
}

watch(() => route.params.username, (username) => username && load(String(username)), { immediate: true })

async function saveProfile() {
  if (!profile.value) return
  isSaving.value = true
  profileFeedback.value = null
  try {
    let avatarUrl = profile.value.avatar_url ?? null
    if (avatarFile.value) {
      const upload = await api.UploadImage('avatar', avatarFile.value.name, await fileToBase64(avatarFile.value))
      avatarUrl = upload?.url ?? avatarUrl
    }
    await api.UpdateProfile({ name: editForm.name.trim(), anchor_dialect: editForm.anchor_dialect, avatar_url: avatarUrl })
    avatarFile.value = null
    profileFeedback.value = { ok: true, text: 'Perfil actualizado.' }
    await Promise.all([refreshProfile(), load(profile.value.user_name)])
  } catch (err) {
    profileFeedback.value = { ok: false, text: errorMessage(err) }
  } finally {
    isSaving.value = false
  }
}

async function changePassword() {
  isChangingPassword.value = true
  passwordFeedback.value = null
  try {
    await api.ChangePassword({ ...passwordForm })
    passwordForm.current_password = ''
    passwordForm.new_password = ''
    passwordFeedback.value = { ok: true, text: 'Contraseña actualizada.' }
  } catch (err) {
    passwordFeedback.value = { ok: false, text: errorMessage(err) }
  } finally {
    isChangingPassword.value = false
  }
}
</script>

<template>
  <div class="max-w-3xl mx-auto space-y-6">
    <!-- Carga -->
    <div v-if="isLoading" class="space-y-4">
      <div class="h-40 rounded-3xl bg-ink/10 animate-pulse" />
      <div class="h-24 rounded-3xl bg-ink/10 animate-pulse" />
    </div>

    <!-- Error / no encontrado -->
    <div v-else-if="!profile" class="card-editorial text-center space-y-3">
      <p class="font-bold text-ink">{{ loadError ?? 'Este usuario no existe.' }}</p>
      <button type="button" class="btn-secondary" @click="load(String(route.params.username))">Reintentar</button>
    </div>

    <template v-else>
      <section class="card-editorial flex flex-col sm:flex-row items-center sm:items-start gap-5">
        <img
          v-if="profile.avatar_url"
          :src="resolveMediaUrl(profile.avatar_url)"
          :alt="`Avatar de ${profile.name}`"
          class="w-24 h-24 rounded-3xl object-cover border border-ink/10 shadow-soft"
        />
        <div
          v-else
          class="w-24 h-24 rounded-3xl bg-brand/15 text-brand font-black text-3xl flex items-center justify-center"
        >
          {{ profile.name.slice(0, 2).toUpperCase() }}
        </div>
        <div class="space-y-2 text-center sm:text-left">
          <h1 class="font-black text-3xl tracking-tight text-ink">{{ profile.name }}</h1>
          <p class="text-ink/65">@{{ profile.user_name }}</p>
          <div class="flex flex-wrap justify-center sm:justify-start gap-2 tabular-nums">
            <span class="badge-pill bg-ink/6 text-ink/65">{{ profile.role }}</span>
            <span class="badge-pill bg-info/15 text-info border border-info/25">Nivel {{ profile.current_level }}</span>
            <span class="badge-pill bg-info/15 text-info border border-info/25">{{ profile.total_xp }} XP</span>
            <span class="badge-pill bg-brand/12 text-brand border border-brand/25">🔥 {{ profile.current_streak }}</span>
            <span class="badge-pill bg-warning/20 text-ink">🏆 {{ profile.matches_won }} victorias</span>
          </div>
        </div>
      </section>

      <section v-if="profile.unlocked_badges?.length" class="card-editorial space-y-3">
        <span class="text-xs font-bold uppercase tracking-wider text-ink/45 block">Medallas</span>
        <div class="flex flex-wrap gap-2">
          <span v-for="badge in profile.unlocked_badges" :key="badge.id_achievement" class="badge-pill bg-warning/20 text-ink" :title="badge.description">
            {{ badge.name }}
          </span>
        </div>
      </section>

      <!-- Edición (solo propietario) -->
      <template v-if="profile.is_owner">
        <form class="card-editorial space-y-4" @submit.prevent="saveProfile">
          <h2 class="font-bold text-xl tracking-tight">Editar perfil</h2>
          <MediaDropzone
            category="avatar"
            :current-image-url="resolveMediaUrl(profile.avatar_url)"
            @file-selected="avatarFile = $event"
            @error="profileFeedback = { ok: false, text: $event }"
          />
          <div class="space-y-1.5">
            <label for="edit-name" class="text-xs font-bold uppercase tracking-wider text-ink/45">Nombre</label>
            <input id="edit-name" v-model="editForm.name" required minlength="2" maxlength="100" class="input-editorial" />
          </div>
          <div class="space-y-1.5">
            <label for="edit-dialect" class="text-xs font-bold uppercase tracking-wider text-ink/45">Región ancla</label>
            <select id="edit-dialect" v-model="editForm.anchor_dialect" class="input-editorial">
              <option v-for="region in regions" :key="region.code" :value="region.code">{{ region.flag_emoji }} {{ region.name }}</option>
            </select>
          </div>
          <p
            v-if="profileFeedback"
            role="status"
            class="card-subtle text-sm"
            :class="profileFeedback.ok ? '!bg-success/12 !border-success/50' : '!bg-warning/20 !border-warning/40'"
          >
            {{ profileFeedback.text }}
          </p>
          <button type="submit" class="btn-primary" :disabled="isSaving">{{ isSaving ? 'Guardando…' : 'Guardar cambios' }}</button>
        </form>

        <form class="card-editorial space-y-4" @submit.prevent="changePassword">
          <h2 class="font-bold text-xl tracking-tight">Cambiar contraseña</h2>
          <div class="space-y-1.5">
            <label for="current-password" class="text-xs font-bold uppercase tracking-wider text-ink/45">Contraseña actual</label>
            <input id="current-password" v-model="passwordForm.current_password" type="password" required autocomplete="current-password" class="input-editorial" />
          </div>
          <div class="space-y-1.5">
            <label for="new-password" class="text-xs font-bold uppercase tracking-wider text-ink/45">Nueva contraseña</label>
            <input id="new-password" v-model="passwordForm.new_password" type="password" required minlength="8" autocomplete="new-password" class="input-editorial" />
          </div>
          <p
            v-if="passwordFeedback"
            role="status"
            class="card-subtle text-sm"
            :class="passwordFeedback.ok ? '!bg-success/12 !border-success/50' : '!bg-warning/20 !border-warning/40'"
          >
            {{ passwordFeedback.text }}
          </p>
          <button type="submit" class="btn-secondary" :disabled="isChangingPassword">
            {{ isChangingPassword ? 'Actualizando…' : 'Actualizar contraseña' }}
          </button>
        </form>
      </template>
    </template>
  </div>
</template>
