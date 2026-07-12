<template>
  <div class="space-y-5">
    <PageHeader title="编码工具" description="编码/解码工具箱，助你完成编码加密类越权挑战" />

    <!-- ====== Row 1: Encode + Decode panels ====== -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-5">
      <!-- ---- Encode Panel ---- -->
      <div class="bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl overflow-hidden">
        <div class="px-4 py-3 border-b border-[var(--border-default)] flex items-center gap-2">
          <i class="pi pi-lock-open text-[var(--primary)] text-sm"></i>
          <span class="text-sm font-semibold text-[var(--text-primary)]">编码</span>
        </div>
        <div class="p-4 space-y-3">
          <!-- Input -->
          <div>
            <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5">原始值</label>
            <InputText v-model="encodeInput" placeholder="输入要编码的 ID 或文本..." class="w-full" @keyup.enter="doEncode" />
          </div>
          <!-- Encoding type -->
          <div class="flex gap-2">
            <div class="flex-1">
              <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5">编码类型</label>
              <Select v-model="encodeType" :options="encodingOpts" optionLabel="label" optionValue="value" class="w-full" />
            </div>
            <!-- Optional extra params -->
            <div v-if="encodeType === 'caesar' || encodeType === 'vigenere'" class="w-28">
              <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5">
                {{ encodeType === 'caesar' ? '偏移量' : '密钥' }}
              </label>
              <InputText v-model="encodeExtra" :placeholder="encodeType === 'caesar' ? '3' : 'KEY'" class="w-full" />
            </div>
          </div>
          <!-- Action -->
          <Button label="编码" icon="pi pi-arrow-right" size="small" :loading="encodeLoading" @click="doEncode" />
          <!-- Output -->
          <div v-if="encodeResult" class="bg-[var(--bg-base)] border border-[var(--border-subtle)] rounded-lg p-3">
            <div class="flex items-center justify-between mb-1">
              <span class="text-[10px] text-[var(--text-tertiary)]">编码结果</span>
              <div class="flex gap-1">
                <button class="text-[10px] text-[var(--primary)] hover:underline" @click="copyText(encodeResult)">复制</button>
                <button
                  v-if="isEncodedApiEndpoint"
                  class="text-[10px] text-[var(--info)] hover:underline"
                  @click="testEncodedRequest(encodeResult)"
                >测试</button>
              </div>
            </div>
            <code class="text-xs text-[var(--text-primary)] mono break-all leading-relaxed">{{ encodeResult }}</code>
          </div>
          <div v-if="encodeError" class="text-xs text-[var(--danger)] bg-[var(--danger-subtle)] rounded px-2 py-1.5">{{ encodeError }}</div>
        </div>
      </div>

      <!-- ---- Decode Panel ---- -->
      <div class="bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl overflow-hidden">
        <div class="px-4 py-3 border-b border-[var(--border-default)] flex items-center gap-2">
          <i class="pi pi-lock text-[var(--primary)] text-sm"></i>
          <span class="text-sm font-semibold text-[var(--text-primary)]">解码</span>
        </div>
        <div class="p-4 space-y-3">
          <!-- Input -->
          <div>
            <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5">已编码值</label>
            <InputText v-model="decodeInput" placeholder="输入已编码的值..." class="w-full" @keyup.enter="doDecode" />
          </div>
          <!-- Encoding type -->
          <div>
            <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5">编码类型</label>
            <Select v-model="decodeType" :options="encodingOpts" optionLabel="label" optionValue="value" class="w-full" />
          </div>
          <Button label="解码" icon="pi pi-arrow-left" size="small" :loading="decodeLoading" @click="doDecode" />
          <div v-if="decodeResult !== null" class="bg-[var(--bg-base)] border border-[var(--border-subtle)] rounded-lg p-3">
            <div class="flex items-center justify-between mb-1">
              <span class="text-[10px] text-[var(--text-tertiary)]">解码结果</span>
              <button class="text-[10px] text-[var(--primary)] hover:underline" @click="copyText(decodeResult)">复制</button>
            </div>
            <code class="text-xs text-[var(--text-primary)] mono break-all leading-relaxed">{{ decodeResult }}</code>
          </div>
          <div v-if="decodeError" class="text-xs text-[var(--danger)] bg-[var(--danger-subtle)] rounded px-2 py-1.5">{{ decodeError }}</div>
        </div>
      </div>
    </div>

    <!-- ====== Row 2: API Tester + Keys ====== -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-5">
      <!-- ---- Encoded API Tester ---- -->
      <div class="bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl overflow-hidden">
        <div class="px-4 py-3 border-b border-[var(--border-default)] flex items-center gap-2">
          <i class="pi pi-send text-[var(--warning)] text-sm"></i>
          <span class="text-sm font-semibold text-[var(--text-primary)]">编码接口测试</span>
        </div>
        <div class="p-4 space-y-3">
          <div class="grid grid-cols-2 gap-2">
            <div>
              <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5">端点</label>
              <Select v-model="testEndpoint" :options="endpointOpts" optionLabel="label" optionValue="value" class="w-full" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5">方法</label>
              <Select v-model="testMethod" :options="methodOpts" optionLabel="label" optionValue="value" class="w-full" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-2">
            <div>
              <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5">目标资源 ID</label>
              <InputText v-model="testResourceId" placeholder="如 5" class="w-full" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-[var(--text-secondary)] mb-1.5">编码类型</label>
              <Select v-model="testEncType" :options="encodingOpts" optionLabel="label" optionValue="value" class="w-full" />
            </div>
          </div>
          <!-- Preview -->
          <div v-if="previewUrl" class="bg-[var(--bg-base)] border border-[var(--border-subtle)] rounded-lg p-2.5 space-y-1.5">
            <div class="text-[10px] text-[var(--text-tertiary)]">请求预览</div>
            <code class="text-[11px] text-[var(--text-primary)] mono break-all block">{{ testMethod }} {{ previewUrl }}</code>
            <code v-if="testEncType" class="text-[10px] text-[var(--info)] mono block">X-Encoding-Type: {{ testEncType }}</code>
          </div>
          <div class="flex gap-2">
            <Button label="发送请求" icon="pi pi-play" size="small" severity="warn" :loading="testLoading" @click="sendTestRequest" />
            <Button label="复制 curl" icon="pi pi-copy" size="small" text @click="copyCurl" />
          </div>
          <!-- Response -->
          <div v-if="testResponse" class="bg-[var(--bg-base)] border rounded-lg p-3" :class="testResponse.status >= 200 && testResponse.status < 300 ? 'border-[var(--success)]/30' : 'border-[var(--danger)]/30'">
            <div class="flex items-center gap-2 mb-2">
              <span class="text-[10px] px-1.5 py-0.5 rounded font-semibold" :class="testResponse.status >= 200 && testResponse.status < 300 ? 'bg-[var(--success-subtle)] text-[var(--success)]' : 'bg-[var(--danger-subtle)] text-[var(--danger)]'">
                {{ testResponse.status }}
              </span>
              <span class="text-[10px] text-[var(--text-tertiary)]">响应体</span>
              <button class="text-[10px] text-[var(--primary)] hover:underline ml-auto" @click="copyText(JSON.stringify(testResponse.data, null, 2))">复制</button>
            </div>
            <pre class="text-[11px] text-[var(--text-secondary)] mono whitespace-pre-wrap break-all max-h-40 overflow-y-auto">{{ JSON.stringify(testResponse.data, null, 2) }}</pre>
          </div>
          <div v-if="testError" class="text-xs text-[var(--danger)] bg-[var(--danger-subtle)] rounded px-2 py-1.5 max-h-24 overflow-y-auto whitespace-pre-wrap">{{ testError }}</div>
        </div>
      </div>

      <!-- ---- Crypto Keys ---- -->
      <div class="bg-[var(--bg-surface)] border border-[var(--border-default)] rounded-xl overflow-hidden">
        <div class="px-4 py-3 border-b border-[var(--border-default)] flex items-center justify-between">
          <div class="flex items-center gap-2">
            <i class="pi pi-key text-[var(--success)] text-sm"></i>
            <span class="text-sm font-semibold text-[var(--text-primary)]">加密密钥</span>
          </div>
          <Button icon="pi pi-refresh" text rounded size="small" :loading="keysLoading" @click="loadKeys" />
        </div>
        <div class="p-4 space-y-2">
          <div v-if="keysLoading" class="text-center py-8 text-sm text-[var(--text-tertiary)]">
            <i class="pi pi-spinner pi-spin mr-2"></i>加载中...
          </div>
          <template v-else-if="cryptoKeys">
            <div v-for="(val, name) in cryptoKeys" :key="name" class="bg-[var(--bg-base)] rounded-lg p-2.5 border border-[var(--border-subtle)]">
              <div class="flex items-center justify-between mb-1">
                <span class="text-[10px] font-semibold text-[var(--text-tertiary)] uppercase">{{ name }}</span>
                <button class="text-[10px] text-[var(--primary)] hover:underline" @click="copyText(String(val))">复制</button>
              </div>
              <code class="text-[10px] text-[var(--text-secondary)] mono break-all block leading-relaxed max-h-16 overflow-y-auto">{{ val }}</code>
            </div>
          </template>
          <div v-if="keysError" class="text-xs text-[var(--danger)] bg-[var(--danger-subtle)] rounded px-2 py-1.5">{{ keysError }}</div>
        </div>
      </div>
    </div>

    <!-- ====== Toast notifications for copy ====== -->
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import * as api from '@/api'
import PageHeader from '@/components/PageHeader.vue'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import { useToast } from 'primevue/usetoast'

const toast = useToast()

// ---- Encoding types ----
const encodingOpts = [
  { label: 'Base64', value: 'base64' },
  { label: 'Base32', value: 'base32' },
  { label: 'Hex', value: 'hex' },
  { label: 'Base58', value: 'base58' },
  { label: 'Base85 (服务端)', value: 'base85' },
  { label: 'Custom Base64 (服务端)', value: 'custom_base64' },
  { label: 'Custom Base32 (服务端)', value: 'custom_base32' },
  { label: '凯撒密码 Caesar', value: 'caesar' },
  { label: '维吉尼亚 Vigenère', value: 'vigenere' },
  { label: 'AES-256-GCM (服务端)', value: 'aes' },
  { label: 'RSA-OAEP (服务端)', value: 'rsa' },
  { label: 'SM4-CBC 国密 (服务端)', value: 'sm4' },
  { label: 'SM2 国密 (服务端)', value: 'sm2' },
  { label: 'HMAC签名 (服务端)', value: 'signed' },
]
const serverOnlyTypes = new Set(['base85', 'custom_base64', 'custom_base32', 'aes', 'rsa', 'sm4', 'sm2', 'signed'])

// ---- Encode ----
const encodeInput = ref('')
const encodeType = ref('base64')
const encodeExtra = ref('')
const encodeResult = ref('')
const encodeError = ref('')
const encodeLoading = ref(false)

// ---- Decode ----
const decodeInput = ref('')
const decodeType = ref('base64')
const decodeResult = ref<string | null>(null)
const decodeError = ref('')
const decodeLoading = ref(false)

// ---- API Tester ----
const endpointOpts = [
  { label: 'GET /api/v1/encoded/vps?v=', value: 'vps_detail' },
  { label: 'POST /api/v1/encoded/vps/start', value: 'vps_start' },
  { label: 'POST /api/v1/encoded/vps/stop', value: 'vps_stop' },
  { label: 'POST /api/v1/encoded/vps/reinstall', value: 'vps_reinstall' },
  { label: 'GET /api/v1/encoded/users?v=', value: 'user' },
  { label: 'GET /api/v1/encoded/orders?v=', value: 'order' },
  { label: 'GET /api/v1/encoded/tickets?v=', value: 'ticket' },
  { label: 'DELETE /api/v1/encoded/apikeys', value: 'apikey' },
  { label: 'POST /api/v1/encoded/company/members', value: 'add_member' },
  { label: 'PUT /api/v1/encoded/company/members/role', value: 'change_role' },
]
const methodOpts = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
]

const testEndpoint = ref('vps_detail')
const testMethod = ref('GET')
const testResourceId = ref('')
const testEncType = ref('base64')
const testResponse = ref<{ status: number; data: any } | null>(null)
const testError = ref('')
const testLoading = ref(false)

// ---- Crypto Keys ----
const cryptoKeys = ref<Record<string, string> | null>(null)
const keysLoading = ref(false)
const keysError = ref('')

// ---- Is the encoded value a potential API endpoint input? ----
const isEncodedApiEndpoint = computed(() => {
  return encodeInput.value.trim().length > 0 && encodeResult.value.length > 0
})

onMounted(() => { loadKeys() })

// ---- Encode logic ----
async function doEncode() {
  const val = encodeInput.value.trim()
  if (!val) return
  encodeError.value = ''
  encodeResult.value = ''
  encodeLoading.value = true

  try {
    if (serverOnlyTypes.has(encodeType.value)) {
      const res = await api.cryptoEncode(val, encodeType.value)
      encodeResult.value = res.data.data!.encoded
    } else {
      // Local encoding via CryptoUtils
      const cu = (window as any).CryptoUtils
      switch (encodeType.value) {
        case 'base64':
          encodeResult.value = cu.base64Encode(val); break
        case 'base32':
          encodeResult.value = cu.base32Encode(val); break
        case 'hex':
          encodeResult.value = cu.hexEncode(val); break
        case 'base58':
          encodeResult.value = cu.base58Encode(val); break
        case 'caesar': {
          const shift = parseInt(encodeExtra.value) || 3
          encodeResult.value = cu.caesarEncode(val, shift); break
        }
        case 'vigenere':
          encodeResult.value = cu.vigenereEncode(val, encodeExtra.value || 'KEY'); break
        default:
          encodeError.value = '不支持的编码类型'
      }
    }
  } catch (e: any) {
    encodeError.value = e.response?.data?.message || e.message || '编码失败'
  } finally {
    encodeLoading.value = false
  }
}

// ---- Decode logic ----
async function doDecode() {
  const val = decodeInput.value.trim()
  if (!val) return
  decodeError.value = ''
  decodeResult.value = null
  decodeLoading.value = true

  try {
    if (serverOnlyTypes.has(decodeType.value)) {
      const res = await api.cryptoDecode(val, decodeType.value)
      decodeResult.value = res.data.data!.decoded
    } else {
      const cu = (window as any).CryptoUtils
      switch (decodeType.value) {
        case 'base64':
          decodeResult.value = cu.base64Decode(val); break
        case 'base32':
          decodeResult.value = cu.base32Decode(val); break
        case 'hex':
          decodeResult.value = cu.hexDecode(val); break
        case 'caesar': {
          const shift = parseInt(encodeExtra.value) || 3
          decodeResult.value = cu.caesarDecode(val, shift); break
        }
        default:
          // Try server-side fallback
          const res = await api.cryptoDecode(val, decodeType.value)
          decodeResult.value = res.data.data!.decoded
      }
    }
  } catch (e: any) {
    decodeError.value = e.response?.data?.message || e.message || '解码失败'
  } finally {
    decodeLoading.value = false
  }
}

// ---- API Tester ----
const endpointUrlMap: Record<string, (id: string) => string> = {
  vps_detail: (id: string) => `/api/v1/encoded/vps?v=${id}`,
  vps_start: (id: string) => `/api/v1/encoded/vps/start`,
  vps_stop: (id: string) => `/api/v1/encoded/vps/stop`,
  vps_reinstall: (id: string) => `/api/v1/encoded/vps/reinstall`,
  user: (id: string) => `/api/v1/encoded/users?v=${id}`,
  order: (id: string) => `/api/v1/encoded/orders?v=${id}`,
  ticket: (id: string) => `/api/v1/encoded/tickets?v=${id}`,
  apikey: (id: string) => `/api/v1/encoded/apikeys`,
  add_member: () => `/api/v1/encoded/company/members`,
  change_role: (id: string) => `/api/v1/encoded/company/members/role`,
}

const previewUrl = computed(() => {
  const id = testResourceId.value.trim()
  if (!id) return ''
  const fn = endpointUrlMap[testEndpoint.value]
  return fn ? fn(id) : ''
})

function testEncodedRequest(encodedVal: string) {
  testResourceId.value = encodedVal
  testEncType.value = encodeType.value
}

async function sendTestRequest() {
  const id = testResourceId.value.trim()
  if (!id) {
    testError.value = '请输入目标资源 ID'
    return
  }
  testError.value = ''
  testResponse.value = null
  testLoading.value = true

  const url = previewUrl.value
  if (!url) {
    testError.value = '无法构建请求 URL'
    testLoading.value = false
    return
  }

  try {
    const token = localStorage.getItem('token')
    const headers: Record<string, string> = {}
    if (token) headers['Authorization'] = `Bearer ${token}`
    if (testEncType.value) headers['X-Encoding-Type'] = testEncType.value

    let resp
    switch (testMethod.value) {
      case 'POST':
        resp = await api.rawClient.post(url, {}, { headers }); break
      case 'PUT':
        resp = await api.rawClient.put(url, {}, { headers }); break
      case 'DELETE':
        resp = await api.rawClient.delete(url, { headers }); break
      default:
        resp = await api.rawClient.get(url, { headers })
    }
    testResponse.value = { status: resp.status, data: resp.data }
  } catch (e: any) {
    if (e.response) {
      testResponse.value = { status: e.response.status, data: e.response.data }
    } else {
      testError.value = e.message || '请求失败'
    }
  } finally {
    testLoading.value = false
  }
}

function copyCurl() {
  const token = localStorage.getItem('token')
  const url = previewUrl.value
  if (!url) return
  let curl = `curl -X ${testMethod.value}`
  if (token) curl += ` -H "Authorization: Bearer ${token}"`
  if (testEncType.value) curl += ` -H "X-Encoding-Type: ${testEncType.value}"`
  curl += ` "${window.location.origin}${url}"`
  copyText(curl)
}

// ---- Keys ----
async function loadKeys() {
  keysLoading.value = true
  keysError.value = ''
  try {
    const res = await api.getCryptoKeys()
    cryptoKeys.value = res.data.data!
  } catch (e: any) {
    keysError.value = e.response?.data?.message || '加载密钥失败'
  } finally {
    keysLoading.value = false
  }
}

// ---- Util ----
async function copyText(text: string) {
  try {
    await navigator.clipboard.writeText(text)
    toast.add({ severity: 'success', summary: '已复制', detail: '内容已复制到剪贴板', life: 1500 })
  } catch {
    toast.add({ severity: 'error', summary: '复制失败', detail: '请手动复制', life: 1500 })
  }
}
</script>
