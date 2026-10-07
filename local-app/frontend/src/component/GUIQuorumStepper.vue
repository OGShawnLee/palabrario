<script setup lang="ts">
import { Check, Clock, UserCheck } from "lucide-vue-next";
import type { QuorumVoteMock } from "@/mocks/editorial.mock";

const props = defineProps<{
  votes: QuorumVoteMock[];
  requiredVotes?: number;
}>();

const maxVotes = props.requiredVotes ?? 3;
</script>

<template>
  <div class="bg-canvas border border-border rounded-2xl p-4 text-chocolate">
    <div class="flex items-center justify-between mb-3">
      <span class="text-xs font-black uppercase tracking-wider text-chocolate/80 flex items-center gap-1.5 font-satoshi">
        <UserCheck class="w-4 h-4 text-brand" />
        <span>Estado del Quórum Editorial ({{ votes.length }}/{{ maxVotes }})</span>
      </span>
      <span
        :class="[
          'text-[10px] font-black uppercase tracking-wider px-2 py-0.5 rounded-full font-supreme',
          votes.length >= maxVotes
            ? 'bg-pastel-mint text-forest'
            : 'bg-pastel-peach text-chocolate'
        ]"
      >
        {{ votes.length >= maxVotes ? 'Aprobación Alcanzada' : 'En Evaluación' }}
      </span>
    </div>

    <!-- Pasos Progresivos del Quórum (RN06) -->
    <div class="space-y-2.5">
      <!-- Votos emitidos -->
      <div
        v-for="(vote, idx) in votes"
        :key="idx"
        class="flex items-center justify-between p-2.5 rounded-xl bg-surface border border-border shadow-2xs text-xs"
      >
        <div class="flex items-center gap-2">
          <div class="w-6 h-6 rounded-full bg-forest text-white flex items-center justify-center font-bold">
            <Check class="w-3.5 h-3.5" />
          </div>
          <div>
            <span class="font-bold text-chocolate">
              {{ vote.flag }} {{ vote.editor_name }}
            </span>
            <span class="text-[10px] text-chocolate/60 block font-supreme">
              {{ vote.editor_role }} ({{ vote.country }}) • {{ vote.voted_at }}
            </span>
          </div>
        </div>

        <span class="text-[10px] font-bold text-forest bg-pastel-mint/50 px-2 py-0.5 rounded-md font-supreme">
          Aprobado ✅
        </span>
      </div>

      <!-- Votos pendientes para completar los 3 requeridos -->
      <div
        v-for="missingIdx in (maxVotes - votes.length)"
        :key="`pending-${missingIdx}`"
        class="flex items-center justify-between p-2.5 rounded-xl bg-canvas border border-dashed border-border/80 text-xs text-chocolate/50"
      >
        <div class="flex items-center gap-2">
          <div class="w-6 h-6 rounded-full bg-border text-chocolate/50 flex items-center justify-center font-bold">
            <Clock class="w-3.5 h-3.5" />
          </div>
          <span class="font-medium italic">Falta {{ missingIdx }} voto para publicación automática...</span>
        </div>

        <span class="text-[10px] text-chocolate/40 font-supreme">Pendiente ⏳</span>
      </div>
    </div>
  </div>
</template>
