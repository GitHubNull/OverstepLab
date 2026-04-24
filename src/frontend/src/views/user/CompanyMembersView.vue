<template>
  <div class="space-y-5">
    <PageHeader title="企业成员" description="管理企业内的用户和权限">
      <template #actions>
        <Button label="添加成员" icon="pi pi-plus" size="small" @click="showAddDialog = true" />
      </template>
    </PageHeader>

    <div class="bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl overflow-hidden">
      <div class="p-0">
        <DataTable
          :value="members"
          class="p-datatable-sm"
          :rows="10"
          paginator
          :loading="loading"
        >
          <Column field="username" header="用户名">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <Avatar
                  :label="data.username.charAt(0).toUpperCase()"
                  shape="circle"
                  size="small"
                  class="bg-[var(--primary-subtle)] text-[var(--primary)] text-xs"
                />
                <span class="font-medium text-sm text-[var(--text-primary)]">{{ data.username }}</span>
              </div>
            </template>
          </Column>

          <Column field="role" header="角色">
            <template #body="{ data }">
              <Tag :value="getRoleText(data.role)" :severity="getRoleSeverity(data.role)" class="text-[10px]" />
            </template>
          </Column>

          <Column field="email" header="邮箱">
            <template #body="{ data }">
              <span class="text-sm text-[var(--text-secondary)]">{{ data.email || '-' }}</span>
            </template>
          </Column>

          <Column field="status" header="状态">
            <template #body="{ data }">
              <Tag :value="data.status === 'active' ? '正常' : '禁用'" :severity="data.status === 'active' ? 'success' : 'danger'" class="text-[10px]" />
            </template>
          </Column>

          <Column header="操作" style="width: 100px">
            <template #body="{ data }">
              <div class="flex gap-1">
                <Button icon="pi pi-pencil" text rounded size="small" class="text-[var(--text-tertiary)]" @click="openEditRole(data)" v-tooltip.top="'修改角色'" />
                <Button icon="pi pi-trash" text rounded size="small" severity="danger" class="!text-[var(--danger)]" @click="handleDelete(data)" v-tooltip.top="'删除成员'" />
              </div>
            </template>
          </Column>
        </DataTable>
      </div>
    </div>

    <!-- Add Member Dialog -->
    <Dialog v-model:visible="showAddDialog" header="添加成员" modal :style="{ width: '420px' }">
      <div class="space-y-3">
        <div>
          <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5 uppercase tracking-wider">用户名</label>
          <InputText v-model="addForm.username" class="w-full" placeholder="请输入用户名" />
        </div>
        <div>
          <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5 uppercase tracking-wider">密码</label>
          <Password v-model="addForm.password" class="w-full" placeholder="请输入密码" :feedback="false" />
        </div>
        <div>
          <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5 uppercase tracking-wider">邮箱</label>
          <InputText v-model="addForm.email" type="email" class="w-full" placeholder="请输入邮箱" />
        </div>
        <div>
          <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5 uppercase tracking-wider">角色</label>
          <Select v-model="addForm.role" :options="roleOptions" optionLabel="label" optionValue="value" placeholder="选择角色" class="w-full" />
        </div>
      </div>
      <template #footer>
        <Button label="取消" text size="small" @click="showAddDialog = false" />
        <Button label="确认添加" icon="pi pi-check" size="small" :loading="adding" @click="handleAdd" />
      </template>
    </Dialog>

    <!-- Edit Role Dialog -->
    <Dialog v-model:visible="showRoleDialog" header="修改角色" modal :style="{ width: '380px' }">
      <div class="space-y-3">
        <p class="text-sm text-[var(--text-secondary)]">修改 <strong class="text-[var(--text-primary)]">{{ editingUser?.username }}</strong> 的角色：</p>
        <Select v-model="editRole" :options="roleOptions" optionLabel="label" optionValue="value" placeholder="选择角色" class="w-full" />
      </div>
      <template #footer>
        <Button label="取消" text size="small" @click="showRoleDialog = false" />
        <Button label="确认修改" icon="pi pi-check" size="small" :loading="changingRole" @click="handleChangeRole" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import * as api from '@/api'
import type { User } from '@/types'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Select from 'primevue/select'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import PageHeader from '@/components/PageHeader.vue'

const toast = useToast()
const confirm = useConfirm()
const members = ref<User[]>([])
const loading = ref(false)

const showAddDialog = ref(false)
const adding = ref(false)
const addForm = ref({ username: '', password: '', email: '', role: 'operator' })

const showRoleDialog = ref(false)
const changingRole = ref(false)
const editingUser = ref<User | null>(null)
const editRole = ref('')

const roleOptions = [
  { label: '管理员', value: 'admin' },
  { label: '运维', value: 'operator' },
  { label: '财务', value: 'finance' },
  { label: '只读', value: 'viewer' },
]

onMounted(() => fetchMembers())

async function fetchMembers() {
  loading.value = true
  try {
    const response = await api.getMembers()
    members.value = response.data.data!
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: '获取成员列表失败', life: 3000 })
  } finally {
    loading.value = false
  }
}

async function handleAdd() {
  if (!addForm.value.username || !addForm.value.password) {
    toast.add({ severity: 'warn', summary: '提示', detail: '请填写用户名和密码', life: 2000 })
    return
  }
  adding.value = true
  try {
    await api.addMember(addForm.value)
    toast.add({ severity: 'success', summary: '成功', detail: '成员已添加', life: 2000 })
    showAddDialog.value = false
    addForm.value = { username: '', password: '', email: '', role: 'operator' }
    await fetchMembers()
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '添加失败', life: 3000 })
  } finally {
    adding.value = false
  }
}

function openEditRole(user: User) {
  editingUser.value = user
  editRole.value = user.role
  showRoleDialog.value = true
}

async function handleChangeRole() {
  if (!editingUser.value) return
  changingRole.value = true
  try {
    await api.changeRole(editingUser.value.id, editRole.value)
    toast.add({ severity: 'success', summary: '成功', detail: '角色已修改', life: 2000 })
    showRoleDialog.value = false
    await fetchMembers()
  } catch (e: any) {
    toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '修改失败', life: 3000 })
  } finally {
    changingRole.value = false
  }
}

function handleDelete(user: User) {
  confirm.require({
    message: `确定要删除成员 "${user.username}" 吗？`,
    header: '确认删除',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        await api.deleteMember(user.id)
        toast.add({ severity: 'success', summary: '成功', detail: '成员已删除', life: 2000 })
        await fetchMembers()
      } catch (e: any) {
        toast.add({ severity: 'error', summary: '错误', detail: e.response?.data?.message || '删除失败', life: 3000 })
      }
    },
  })
}

function getRoleText(role: string) {
  const map: Record<string, string> = { admin: '管理员', operator: '运维', finance: '财务', viewer: '只读' }
  return map[role] || role
}

function getRoleSeverity(role: string) {
  const map: Record<string, string> = { admin: 'danger', operator: 'warn', finance: 'info', viewer: 'secondary' }
  return map[role] || 'secondary'
}
</script>
