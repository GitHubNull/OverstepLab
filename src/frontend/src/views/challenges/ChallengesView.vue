<template>
  <div class="space-y-5">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <div class="page-header">
        <h2>漏洞挑战</h2>
        <p>发现并练习越权漏洞，提升安全测试技能</p>
      </div>
      <div class="flex items-center gap-2">
        <div class="flex items-center gap-1.5 px-3 py-1.5 bg-emerald-50 dark:bg-emerald-900/10 rounded-lg">
          <i class="pi pi-check-circle text-emerald-500 text-xs"></i>
          <span class="text-xs font-semibold text-emerald-600 dark:text-emerald-400">{{ completedCount }}/13</span>
        </div>
      </div>
    </div>

    <!-- Category Filters -->
    <div class="flex gap-2 flex-wrap">
      <Button
        v-for="cat in categories"
        :key="cat.value"
        :label="cat.label"
        :severity="selectedCategory === cat.value ? 'primary' : 'secondary'"
        :outlined="selectedCategory !== cat.value"
        size="small"
        @click="selectedCategory = cat.value"
      />
    </div>

    <!-- Challenges Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <Card
        v-for="ch in filteredChallenges"
        :key="ch.id"
        class="shadow-none hover:shadow-md"
        :class="{ 'ring-2 ring-emerald-400/50': ch.completed }"
      >
        <template #content>
          <div class="space-y-3">
            <!-- Header -->
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <span class="font-mono text-xs font-bold text-slate-400">{{ ch.id }}</span>
                <Tag :severity="getCategorySeverity(ch.category)" :value="ch.category" class="text-[10px]" />
              </div>
              <i v-if="ch.completed" class="pi pi-check-circle text-emerald-500"></i>
            </div>

            <!-- Title -->
            <h3 class="font-semibold text-sm text-slate-700 dark:text-white">{{ ch.title }}</h3>

            <!-- Description -->
            <p class="text-xs text-slate-500 leading-relaxed">{{ ch.description }}</p>

            <!-- Difficulty -->
            <div class="flex items-center gap-1.5">
              <span class="text-[10px] text-slate-400">难度:</span>
              <div class="flex gap-0.5">
                <i v-for="i in 3" :key="i" class="pi text-xs" :class="i <= ch.difficulty ? 'pi-star-fill text-amber-400' : 'pi-star text-slate-200'"></i>
              </div>
            </div>

            <!-- Endpoint -->
            <div class="bg-slate-50 dark:bg-slate-700/30 rounded-lg p-2">
              <div class="flex items-center gap-2">
                <Tag :value="(ch as any).method" severity="info" class="text-[10px]" />
                <code class="text-[10px] text-slate-500 dark:text-slate-400 truncate">{{ (ch as any).endpoint }}</code>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex gap-2 pt-1">
              <Button label="提示" icon="pi pi-lightbulb" text size="small" class="flex-1 !text-xs" @click="showHints(ch)" />
              <Button label="WriteUp" icon="pi pi-book" text size="small" class="flex-1 !text-xs" @click="showWriteup(ch)" />
              <Button
                v-if="!ch.completed"
                label="完成"
                icon="pi pi-check"
                text
                size="small"
                severity="success"
                class="!text-xs"
                @click="markComplete(ch)"
              />
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Hint Dialog -->
    <Dialog v-model:visible="hintDialogVisible" :header="`提示 - ${selectedChallenge?.title}`" modal :style="{ width: '500px' }">
      <div v-if="selectedChallenge" class="space-y-3">
        <div v-for="(hint, index) in (selectedChallenge as any).hints" :key="index" class="space-y-1">
          <div class="flex items-center gap-2">
            <Tag :value="`Level ${index + 1}`" severity="warn" class="text-[10px]" />
          </div>
          <p class="text-sm text-slate-600 dark:text-slate-400 pl-1">{{ hint }}</p>
        </div>
      </div>
    </Dialog>

    <!-- WriteUp Dialog -->
    <Dialog v-model:visible="writeupDialogVisible" :header="`WriteUp - ${selectedChallenge?.title}`" modal :style="{ width: '600px' }">
      <div v-if="selectedChallenge" class="space-y-4">
        <div class="bg-slate-50 dark:bg-slate-700/30 rounded-xl p-4">
          <h4 class="font-semibold text-sm text-slate-700 dark:text-white mb-2">漏洞分析</h4>
          <p class="text-sm text-slate-600 dark:text-slate-400 leading-relaxed">{{ (selectedChallenge as any).writeup }}</p>
        </div>
        <div class="bg-slate-50 dark:bg-slate-700/30 rounded-xl p-4">
          <h4 class="font-semibold text-sm text-slate-700 dark:text-white mb-2">影响</h4>
          <p class="text-sm text-slate-600 dark:text-slate-400">该漏洞可能导致未授权访问敏感数据或执行未授权操作。</p>
        </div>
        <div class="bg-emerald-50 dark:bg-emerald-900/10 border border-emerald-200 dark:border-emerald-800/30 rounded-xl p-4">
          <h4 class="font-semibold text-sm text-emerald-700 dark:text-emerald-400 mb-2">修复方案</h4>
          <p class="text-sm text-slate-600 dark:text-slate-400">在安全模式下，系统会启用完整的权限校验，阻止该漏洞的利用。</p>
        </div>
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import * as api from '@/api'
import type { Challenge } from '@/types'
import Card from 'primevue/card'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import { useToast } from 'primevue/usetoast'

const toast = useToast()
const challenges = ref<Challenge[]>([])
const selectedCategory = ref('all')
const hintDialogVisible = ref(false)
const writeupDialogVisible = ref(false)
const selectedChallenge = ref<Challenge | null>(null)

const categories = [
  { label: '全部', value: 'all' },
  { label: '水平越权', value: '水平越权' },
  { label: '垂直越权', value: '垂直越权' },
  { label: '上下文越权', value: '上下文越权' },
]

onMounted(async () => {
  const response = await api.getChallenges()
  challenges.value = response.data.data!
})

const filteredChallenges = computed(() => {
  if (selectedCategory.value === 'all') return challenges.value
  return challenges.value.filter(ch => ch.category.includes(selectedCategory.value))
})

const completedCount = computed(() => challenges.value.filter(ch => ch.completed).length)

function getCategorySeverity(cat: string) {
  if (cat.includes('水平')) return 'info'
  if (cat.includes('垂直')) return 'warn'
  return 'danger'
}

function showHints(ch: Challenge) {
  selectedChallenge.value = ch
  hintDialogVisible.value = true
}

function showWriteup(ch: Challenge) {
  selectedChallenge.value = ch
  writeupDialogVisible.value = true
}

async function markComplete(ch: Challenge) {
  try {
    await api.markChallengeComplete(ch.id)
    ch.completed = true
    toast.add({ severity: 'success', summary: '恭喜', detail: `已完成挑战 ${ch.id}`, life: 2000 })
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: '标记失败', life: 3000 })
  }
}
</script>
