<template>
  <div>
    <PageHeader title="公告管理" description="发布、编辑和管理平台公告">
      <template #actions>
        <Button label="新建公告" icon="pi pi-plus" size="small" @click="showCreate" />
      </template>
    </PageHeader>

    <!-- Announcements Table -->
    <Card>
      <template #content>
        <DataTable :value="announcements" stripedRows class="p-datatable-sm" :rows="10" paginator :loading="loading">
          <Column field="id" header="ID" style="width: 60px">
            <template #body="{ data }">
              <span class="mono text-[10px] px-1.5 py-0.5 rounded bg-[var(--bg-base)] text-[var(--text-secondary)]">{{ data.id }}</span>
            </template>
          </Column>
          <Column field="title" header="标题">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <i v-if="data.is_pinned" class="pi pi-thumbtack text-[var(--primary)] text-xs" />
                <span class="font-medium text-sm text-[var(--text-primary)]">{{ data.title }}</span>
              </div>
            </template>
          </Column>
          <Column field="created_at" header="发布时间">
            <template #body="{ data }">
              <span class="text-xs text-[var(--text-secondary)]">{{ formatDate(data.created_at) }}</span>
            </template>
          </Column>
          <Column field="is_pinned" header="置顶" style="width: 80px">
            <template #body="{ data }">
              <Tag :value="data.is_pinned ? '是' : '否'" :severity="data.is_pinned ? 'info' : 'secondary'" class="text-[10px]" />
            </template>
          </Column>
          <Column header="操作" style="width: 120px">
            <template #body="{ data }">
              <div class="flex gap-1">
                <Button icon="pi pi-pencil" text rounded size="small" title="编辑" @click="showEdit(data)" />
                <Button icon="pi pi-trash" text rounded size="small" severity="danger" title="删除" @click="handleDelete(data.id)" />
              </div>
            </template>
          </Column>
        </DataTable>

        <div v-if="announcements.length === 0 && !loading" class="empty-state">
          <div class="empty-state-icon">
            <i class="pi pi-megaphone text-2xl text-[var(--text-tertiary)]"></i>
          </div>
          <p class="text-[var(--text-secondary)] text-sm">暂无公告</p>
        </div>
      </template>
    </Card>

    <!-- Create/Edit Dialog -->
    <Dialog v-model:visible="dialogVisible" :header="editingId ? '编辑公告' : '新建公告'" modal :style="{ width: '500px' }">
      <div class="space-y-4">
        <div>
          <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5">标题</label>
          <InputText v-model="form.title" class="w-full" placeholder="公告标题" />
        </div>
        <div>
          <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5">内容</label>
          <Textarea v-model="form.content" class="w-full" rows="4" placeholder="公告内容" />
        </div>
        <div class="flex items-center gap-2">
          <Checkbox v-model="form.is_pinned" binary inputId="pinned" />
          <label for="pinned" class="text-xs text-[var(--text-secondary)]">置顶公告</label>
        </div>
      </div>
      <template #footer>
        <Button label="取消" text size="small" @click="dialogVisible = false" />
        <Button label="保存" icon="pi pi-check" size="small" :loading="saving" @click="handleSave" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getAnnouncements, adminCreateAnnouncement, adminUpdateAnnouncement, adminDeleteAnnouncement } from '@/api'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import Button from 'primevue/button'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Checkbox from 'primevue/checkbox'
import PageHeader from '@/components/PageHeader.vue'
import { formatDate } from '@/utils/date'
import type { Announcement } from '@/types'

const toast = useToast()
const confirm = useConfirm()
const announcements = ref<Announcement[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const saving = ref(false)
const editingId = ref<number | null>(null)

const form = ref({ title: '', content: '', is_pinned: false })

onMounted(() => { fetchList() })

async function fetchList() {
  loading.value = true
  try {
    const res = await getAnnouncements()
    announcements.value = res.data.data!
  } finally {
    loading.value = false
  }
}

function showCreate() {
  editingId.value = null
  form.value = { title: '', content: '', is_pinned: false }
  dialogVisible.value = true
}

function showEdit(a: Announcement) {
  editingId.value = a.id
  form.value = { title: a.title, content: a.content, is_pinned: a.is_pinned }
  dialogVisible.value = true
}

async function handleSave() {
  if (!form.value.title || !form.value.content) {
    toast.add({ severity: 'warn', summary: '提示', detail: '请填写标题和内容', life: 2000 })
    return
  }
  saving.value = true
  try {
    if (editingId.value) {
      await adminUpdateAnnouncement(editingId.value, form.value)
      toast.add({ severity: 'success', summary: '成功', detail: '公告已更新', life: 2000 })
    } else {
      await adminCreateAnnouncement(form.value)
      toast.add({ severity: 'success', summary: '成功', detail: '公告已发布', life: 2000 })
    }
    dialogVisible.value = false
    fetchList()
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '操作失败', life: 3000 })
  } finally {
    saving.value = false
  }
}

function handleDelete(id: number) {
  confirm.require({
    message: '确定要删除此公告吗？',
    header: '确认删除',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      await adminDeleteAnnouncement(id)
      toast.add({ severity: 'success', summary: '成功', detail: '公告已删除' })
      fetchList()
    },
  })
}
</script>
