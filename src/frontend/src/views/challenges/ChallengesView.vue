<template>
  <div class="space-y-5">
    <!-- Header with progress -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
      <PageHeader title="漏洞挑战" description="发现并练习越权漏洞，提升安全测试技能" />
      <div class="flex items-center gap-3">
        <div class="flex items-center gap-2 px-3 py-1.5 bg-[var(--success-subtle)] border border-[var(--success)]/20 rounded-lg">
          <i class="pi pi-check-circle text-[var(--success)] text-xs"></i>
          <span class="text-xs font-semibold text-[var(--success)]">{{ completedCount }}/{{ totalCount }}</span>
        </div>
      </div>
    </div>

    <!-- Progress Bar -->
    <div class="bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl p-4">
      <div class="flex items-center justify-between mb-2">
        <span class="text-xs font-medium text-[var(--text-secondary)]">总体进度</span>
        <span class="text-xs font-bold text-[var(--primary)] mono">{{ Math.round((completedCount / totalCount) * 100) }}%</span>
      </div>
      <div class="h-2 bg-[var(--bg-base)] rounded-full overflow-hidden">
        <div
          class="h-full bg-[var(--primary)] rounded-full transition-all duration-500"
          :style="{ width: `${(completedCount / totalCount) * 100}%` }"
        ></div>
      </div>
    </div>

    <!-- Category Filters -->
    <div class="flex gap-2 flex-wrap">
      <button
        v-for="cat in categories"
        :key="cat.value"
        class="px-3 py-1.5 rounded-lg text-xs font-medium transition-all duration-200 border"
        :class="selectedCategory === cat.value
          ? 'bg-[var(--primary-subtle)] border-[var(--primary)]/30 text-[var(--primary)]'
          : 'bg-[var(--bg-surface)] border-[var(--border-default)] text-[var(--text-secondary)] hover:border-[var(--border-strong)]'"
        @click="selectedCategory = cat.value"
      >
        {{ cat.label }}
      </button>
    </div>

    <!-- Global Encoding Challenge State Banner -->
    <div
      v-if="encodingState.active"
      class="bg-[var(--primary-subtle)] border border-[var(--primary)]/20 rounded-xl p-4 flex items-center justify-between"
    >
      <div class="flex items-center gap-3">
        <i class="pi pi-lock text-[var(--primary)]"></i>
        <div>
          <p class="text-sm font-semibold text-[var(--primary)]">
            编码挑战已激活: {{ encodingState.challenge_name || encodingState.challenge_id }}
          </p>
          <p class="text-xs text-[var(--text-secondary)]">
            编码类型: {{ encodingState.encoding_type }} — 所有业务请求将自动编码参数
          </p>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <Tag :value="encodingState.encoding_type.toUpperCase()" severity="info" class="text-xs" />
        <Button
          label="关闭"
          icon="pi pi-times"
          text
          size="small"
          severity="secondary"
          @click="deactivateEncodingChallenge()"
        />
      </div>
    </div>

    <!-- Challenges Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <div
        v-for="ch in filteredChallenges"
        :key="ch.id"
        class="bg-[var(--bg-surface)] border rounded-xl overflow-hidden transition-all duration-200 hover:border-[var(--border-strong)] relative"
        :class="[
          ch.completed ? 'border-[var(--success)]/30' : 'border-[var(--border-default)]'
        ]"
      >
        <!-- Completed indicator -->
        <div
          v-if="ch.completed"
          class="absolute left-0 top-3 bottom-3 w-[3px] bg-[var(--success)] rounded-r-full"
        ></div>
        <div class="p-4 space-y-3" :class="ch.completed ? 'pl-5' : ''">
          <!-- Header -->
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="mono text-xs font-bold text-[var(--text-tertiary)]">{{ ch.id }}</span>
              <Tag
                :value="ch.category"
                class="text-[10px]"
                :style="getCategoryStyle(ch.category)"
              />
            </div>
            <i v-if="ch.completed" class="pi pi-check-circle text-[var(--success)]"></i>
          </div>

          <!-- Title -->
          <h3 class="font-semibold text-sm text-[var(--text-primary)]">{{ ch.title }}</h3>

          <!-- Description -->
          <p class="text-xs text-[var(--text-secondary)] leading-relaxed">{{ ch.description }}</p>

          <!-- Difficulty -->
          <div class="flex items-center gap-1.5">
            <span class="text-[10px] text-[var(--text-tertiary)]">难度:</span>
            <div class="flex gap-0.5">
              <i
                v-for="i in maxDifficulty"
                :key="i"
                class="pi text-xs"
                :class="i <= ch.difficulty ? 'pi-bolt text-[var(--warning)]' : 'pi-bolt text-[var(--border-default)]'"
              />
            </div>
          </div>

          <!-- Endpoint -->
          <div v-if="ch.endpoint && ch.method" class="bg-[var(--bg-base)] rounded-lg p-2.5 border border-[var(--border-subtle)]">
            <div class="flex items-center gap-2">
              <code class="text-[10px] px-1.5 py-0.5 rounded bg-[var(--primary-subtle)] text-[var(--primary)] mono font-semibold">{{ ch.method }}</code>
              <code class="text-[11px] text-[var(--text-secondary)] mono truncate">{{ ch.endpoint }}</code>
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
            <Button
              v-if="ch.category.includes('编码加密')"
              :label="encodingState.active && encodingState.challenge_id === ch.id ? '关闭' : '激活'"
              :icon="encodingState.active && encodingState.challenge_id === ch.id ? 'pi pi-times' : 'pi pi-play'"
              text
              size="small"
              severity="info"
              class="!text-xs"
              @click="toggleEncodingChallenge(ch)"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Hint Dialog -->
    <Dialog v-model:visible="hintDialogVisible" :header="`提示 - ${selectedChallenge?.title}`" modal :style="{ width: '500px' }">
      <div v-if="selectedChallenge" class="space-y-3">
        <div v-for="(hint, index) in challengeDetail?.hints" :key="index" class="space-y-1">
          <div class="flex items-center gap-2">
            <Tag :value="`Level ${index + 1}`" severity="warn" class="text-[10px]" />
          </div>
          <p class="text-sm text-[var(--text-secondary)] pl-1">{{ hint }}</p>
        </div>
      </div>
    </Dialog>

    <!-- WriteUp Dialog -->
    <Dialog v-model:visible="writeupDialogVisible" :header="`WriteUp - ${selectedChallenge?.title}`" modal :style="{ width: '600px' }">
      <div v-if="selectedChallenge" class="space-y-4">
        <div class="bg-[var(--bg-base)] rounded-xl p-4 border border-[var(--border-default)]">
          <h4 class="font-semibold text-sm text-[var(--text-primary)] mb-2">漏洞分析</h4>
          <p class="text-sm text-[var(--text-secondary)] leading-relaxed">{{ challengeDetail?.writeup }}</p>
        </div>
        <div class="bg-[var(--bg-base)] rounded-xl p-4 border border-[var(--border-default)]">
          <h4 class="font-semibold text-sm text-[var(--text-primary)] mb-2">影响</h4>
          <p class="text-sm text-[var(--text-secondary)]">该漏洞可能导致未授权访问敏感数据或执行未授权操作。</p>
        </div>
        <div class="bg-[var(--success-subtle)] border border-[var(--success)]/20 rounded-xl p-4">
          <h4 class="font-semibold text-sm text-[var(--success)] mb-2">修复方案</h4>
          <p class="text-sm text-[var(--text-secondary)]">在安全模式下，系统会启用完整的权限校验，阻止该漏洞的利用。</p>
        </div>
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import * as api from '@/api'
import type { Challenge, ChallengeDetail } from '@/types'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import { useToast } from 'primevue/usetoast'
import PageHeader from '@/components/PageHeader.vue'

const toast = useToast()
const challenges = ref<Challenge[]>([])
const selectedCategory = ref('all')
const hintDialogVisible = ref(false)
const writeupDialogVisible = ref(false)
const selectedChallenge = ref<Challenge | null>(null)
const challengeDetail = ref<ChallengeDetail | null>(null)

// Global encoding state (managed by backend admin)
const encodingState = ref({
  active: false,
  challenge_id: null as string | null,
  encoding_type: 'none',
  challenge_name: null as string | null,
})

const categories = [
  { label: '全部', value: 'all' },
  { label: '水平越权', value: '水平越权' },
  { label: '垂直越权', value: '垂直越权' },
  { label: '上下文越权', value: '上下文越权' },
  { label: '编码加密', value: '编码加密' },
]

const maxDifficulty = computed(() => {
  if (challenges.value.length === 0) return 3
  return Math.max(...challenges.value.map(ch => ch.difficulty))
})

onMounted(async () => {
  try {
    const [chRes, encRes] = await Promise.all([
      api.getChallenges(),
      api.getEncodingChallengeState().catch(() => null),
    ])
    challenges.value = chRes.data.data!
    if (encRes?.data?.data) {
      encodingState.value = encRes.data.data
    }
  } catch (e: any) {
    toast.add({
      severity: 'error',
      summary: '加载失败',
      detail: e.response?.data?.message || '无法获取挑战列表',
      life: 3000,
    })
  }
})

const filteredChallenges = computed(() => {
  if (selectedCategory.value === 'all') return challenges.value
  return challenges.value.filter(ch => ch.category.includes(selectedCategory.value))
})

const completedCount = computed(() => challenges.value.filter(ch => ch.completed).length)
const totalCount = computed(() => challenges.value.length)

function getCategoryStyle(cat: string) {
  if (cat.includes('水平')) {
    return {
      background: 'var(--info-subtle)',
      color: 'var(--info)',
    }
  }
  if (cat.includes('垂直')) {
    return {
      background: 'var(--warning-subtle)',
      color: 'var(--warning)',
    }
  }
  return {
    background: 'var(--danger-subtle)',
    color: 'var(--danger)',
  }
}

async function loadChallengeDetail(ch: Challenge) {
  selectedChallenge.value = ch
  try {
    const response = await api.getChallengeDetail(ch.id)
    challengeDetail.value = response.data.data!
  } catch (e: any) {
    challengeDetail.value = null
    toast.add({
      severity: 'error',
      summary: '加载失败',
      detail: e.response?.data?.message || '无法获取挑战详情',
      life: 3000,
    })
  }
}

async function showHints(ch: Challenge) {
  await loadChallengeDetail(ch)
  hintDialogVisible.value = true
}

async function showWriteup(ch: Challenge) {
  await loadChallengeDetail(ch)
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

async function toggleEncodingChallenge(ch: Challenge) {
  const isActive = encodingState.value.active && encodingState.value.challenge_id === ch.id
  try {
    await api.setEncodingChallengeState({
      challenge_id: ch.id,
      encoding_type: ch.encoding_type || 'base64',
      challenge_name: ch.title,
      active: !isActive,
    })
    encodingState.value = {
      active: !isActive,
      challenge_id: !isActive ? ch.id : null,
      encoding_type: !isActive ? (ch.encoding_type || 'base64') : 'none',
      challenge_name: !isActive ? ch.title : null,
    }
    toast.add({
      severity: 'success',
      summary: !isActive ? '已激活' : '已关闭',
      detail: `编码挑战 ${ch.id}: ${ch.title} ${!isActive ? '已激活' : '已关闭'}`,
      life: 2000,
    })
  } catch (e: any) {
    toast.add({
      severity: 'error',
      summary: '操作失败',
      detail: e.response?.data?.message || '无法更新编码挑战状态',
      life: 3000,
    })
  }
}

async function deactivateEncodingChallenge() {
  try {
    await api.setEncodingChallengeState({
      challenge_id: encodingState.value.challenge_id || '',
      encoding_type: encodingState.value.encoding_type,
      challenge_name: encodingState.value.challenge_name || '',
      active: false,
    })
    encodingState.value = { active: false, challenge_id: null, encoding_type: 'none', challenge_name: null }
    toast.add({ severity: 'success', summary: '已关闭', detail: '编码挑战已关闭', life: 2000 })
  } catch (e: any) {
    toast.add({
      severity: 'error',
      summary: '关闭失败',
      detail: e.response?.data?.message || '无法关闭编码挑战',
      life: 3000,
    })
  }
}
</script>
