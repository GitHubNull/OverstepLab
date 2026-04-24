<template>
  <div
    class="relative bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl p-4 overflow-hidden transition-all duration-200 hover:border-[var(--border-strong)]"
  >
    <!-- Left color bar -->
    <div
      class="absolute left-0 top-3 bottom-3 w-[3px] rounded-r-full"
      :style="{ backgroundColor: color }"
    ></div>
    <div class="flex items-center gap-3 pl-2">
      <!-- Icon container -->
      <div
        class="w-9 h-9 rounded-lg flex items-center justify-center flex-shrink-0"
        :style="{ backgroundColor: `${color}14` }"
      >
        <i :class="[icon, 'text-sm']" :style="{ color }"></i>
      </div>
      <div class="min-w-0">
        <p class="text-[11px] font-medium text-[var(--text-secondary)] uppercase tracking-wider">{{ label }}</p>
        <p class="text-xl font-bold text-[var(--text-primary)] mono leading-tight mt-0.5">{{ value }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  color: string
  icon: string
  value: string | number
  label: string
}>()
</script>
