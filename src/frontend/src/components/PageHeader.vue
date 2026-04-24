<template>
  <div class="flex items-start justify-between mb-6 gap-4">
    <div>
      <h1 class="text-xl font-bold text-[var(--text-primary)] leading-tight">{{ title }}</h1>
      <p v-if="description" class="text-sm text-[var(--text-secondary)] mt-1">{{ description }}</p>
    </div>
    <div v-if="$slots.actions" class="flex-shrink-0">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  description?: string
}>()
</script>
