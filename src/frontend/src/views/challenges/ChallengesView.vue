<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-slate-800 dark:text-white">漏洞挑战</h2>
        <p class="text-slate-500">发现并练习越权漏洞，提升安全测试技能</p>
      </div>
      <div class="flex items-center gap-2">
        <Tag 
          :value="`已完成: ${completedCount}/13`" 
          severity="success"
          class="text-sm"
        />
      </div>
    </div>

    <!-- Category Filters -->
    <div class="flex gap-2">
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
        class="shadow-sm hover:shadow-lg transition-all duration-300"
        :class="{ 'ring-2 ring-green-500': ch.completed }"
      >
        <template #title>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="font-mono text-sm font-bold text-slate-500">{{ ch.id }}</span>
              <Tag 
                :severity="getCategorySeverity(ch.category)" 
                :value="ch.category"
                class="text-xs"
              />
            </div>
            <i 
              v-if="ch.completed"
              class="pi pi-check-circle text-green-500 text-xl"
            ></i>
          </div>
        </template>

        <template #subtitle>
          <div class="flex items-center gap-2 mt-1">
            <span class="font-semibold text-slate-800 dark:text-white">{{ ch.title }}</span>
          </div>
        </template>

        <template #content>
          <div class="space-y-4">
            <p class="text-sm text-slate-600 dark:text-slate-400">{{ ch.description }}</p>

            <!-- Difficulty -->
            <div class="flex items-center gap-2">
              <span class="text-xs text-slate-500">难度:</span>
              <div class="flex gap-1">
                <i 
                  v-for="i in 3" 
                  :key="i"
                  class="pi text-sm"
                  :class="i <= ch.difficulty ? 'pi-star-fill text-yellow-500' : 'pi-star text-slate-300'"
                ></i>
              </div>
            </div>

            <!-- Endpoint -->
            <div class="bg-slate-50 dark:bg-slate-700/50 rounded-lg p-2">
              <div class="flex items-center gap-2">
                <Tag 
                  :value="(ch as any).method" 
                  severity="info"
                  class="text-xs"
                />
                <code class="text-xs text-slate-600 dark:text-slate-400">{{ (ch as any).endpoint }}</code>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex gap-2 pt-2">
              <Button 
                label="查看提示" 
                icon="pi pi-lightbulb" 
                text 
                size="small"
                class="flex-1"
                @click="showHints(ch)"
              />
              <Button 
                label="WriteUp" 
                icon="pi pi-book" 
                text 
                size="small"
                class="flex-1"
                @click="showWriteup(ch)"
              />
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Hint Dialog -->
    <Dialog 
      v-model:visible="hintDialogVisible" 
      :header="`提示 - ${selectedChallenge?.title}`"
      modal
      class="max-w-lg"
    >
      <div v-if="selectedChallenge" class="space-y-4">
        <div v-for="(hint, index) in (selectedChallenge as any).hints" :key="index" class="space-y-2">
          <div class="flex items-center gap-2">
            <Tag :value="`Level ${index + 1}`" severity="warning" class="text-xs" />
          </div>
          <p class="text-sm text-slate-600 dark:text-slate-400">{{ hint }}</p>
        </div>
      </div>
    </Dialog>

    <!-- WriteUp Dialog -->
    <Dialog 
      v-model:visible="writeupDialogVisible" 
      :header="`WriteUp - ${selectedChallenge?.title}`"
      modal
      class="max-w-2xl"
    >
      <div v-if="selectedChallenge" class="space-y-4">
        <div class="bg-slate-50 dark:bg-slate-700/50 rounded-xl p-4">
          <h4 class="font-semibold text-slate-800 dark:text-white mb-2">漏洞分析</h4>
          <p class="text-sm text-slate-600 dark:text-slate-400">{{ (selectedChallenge as any).writeup }}</p>
        </div>

        <div class="bg-slate-50 dark:bg-slate-700/50 rounded-xl p-4">
          <h4 class="font-semibold text-slate-800 dark:text-white mb-2">影响</h4>
          <p class="text-sm text-slate-600 dark:text-slate-400">该漏洞可能导致未授权访问敏感数据或执行未授权操作。</p>
        </div>

        <div class="bg-slate-50 dark:bg-slate-700/50 rounded-xl p-4">
          <h4 class="font-semibold text-slate-800 dark:text-white mb-2">修复方案</h4>
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
  if (cat.includes('IDOR')) return 'info'
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
</script>
