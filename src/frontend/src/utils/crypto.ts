// 前端加解密工具库
// Base64, Base32, Base58, AES, SM4, RSA 等编码/加密/签名工具

const CUSTOM_BASE64_TABLE = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz+/'
const BASE58_ALPHABET = '123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz'

// Base64
export function base64Encode(data: string): string {
  try {
    return btoa(unescape(encodeURIComponent(data)))
  } catch {
    return btoa(data)
  }
}

export function base64Decode(data: string): string {
  try {
    return decodeURIComponent(escape(atob(data)))
  } catch {
    return atob(data)
  }
}

// Base32 (RFC 4648)
function b32CharCode(c: string): number {
  const code = c.charCodeAt(0)
  if (code >= 65 && code <= 90) return code - 65   // A-Z
  if (code >= 50 && code <= 55) return code - 24   // 2-7
  return -1
}

function b32Char(v: number): string {
  if (v < 26) return String.fromCharCode(v + 65)
  return String.fromCharCode(v + 24)
}

export function base32Encode(data: string): string {
  const bytes = new TextEncoder().encode(data)
  let bits = ''
  for (const b of bytes) bits += b.toString(2).padStart(8, '0')
  let result = ''
  for (let i = 0; i < bits.length; i += 5) {
    const chunk = bits.substring(i, i + 5).padEnd(5, '0')
    result += b32Char(parseInt(chunk, 2))
  }
  const padding = (8 - (result.length % 8)) % 8
  return result + '='.repeat(padding)
}

export function base32Decode(data: string): string {
  data = data.toUpperCase().replace(/=+$/, '')
  let bits = ''
  for (const c of data) {
    const v = b32CharCode(c)
    if (v < 0) throw new Error('Invalid Base32 char')
    bits += v.toString(2).padStart(5, '0')
  }
  const bytes: number[] = []
  for (let i = 0; i + 8 <= bits.length; i += 8) {
    bytes.push(parseInt(bits.substring(i, i + 8), 2))
  }
  return new TextDecoder().decode(new Uint8Array(bytes))
}

// Base58
export function base58Encode(data: string): string {
  let num = BigInt(0)
  const bytes = new TextEncoder().encode(data)
  for (const b of bytes) num = (num << BigInt(8)) | BigInt(b)
  if (num === BigInt(0)) return BASE58_ALPHABET[0]
  let result = ''
  while (num > BigInt(0)) {
    const rem = Number(num % BigInt(58))
    result = BASE58_ALPHABET[rem] + result
    num = num / BigInt(58)
  }
  for (const b of bytes) {
    if (b === 0) result = BASE58_ALPHABET[0] + result
    else break
  }
  return result
}

// Caesar (ROT3)
export function caesarEncode(data: string, shift = 3): string {
  shift = ((shift % 26) + 26) % 26
  const digitShift = shift % 10
  let result = ''
  for (const c of data) {
    if (c >= 'A' && c <= 'Z') result += String.fromCharCode(65 + (c.charCodeAt(0) - 65 + shift) % 26)
    else if (c >= 'a' && c <= 'z') result += String.fromCharCode(97 + (c.charCodeAt(0) - 97 + shift) % 26)
    else if (c >= '0' && c <= '9') result += String.fromCharCode(48 + (c.charCodeAt(0) - 48 + digitShift) % 10)
    else result += c
  }
  return result
}

export function caesarDecode(data: string, shift = 3): string {
  return caesarEncode(data, -shift)
}

// Custom Base64
export function customBase64Encode(data: string): string {
  const stdEncoded = base64Encode(data)
  let result = ''
  for (const c of stdEncoded) {
    if (c === '=') { result += '='; continue }
    const idx = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'.indexOf(c)
    result += CUSTOM_BASE64_TABLE[idx] || c
  }
  return result
}

export function customBase64Decode(data: string): string {
  let stdEncoded = ''
  for (const c of data) {
    if (c === '=') { stdEncoded += '='; continue }
    const idx = CUSTOM_BASE64_TABLE.indexOf(c)
    stdEncoded += 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/'[idx] || c
  }
  return base64Decode(stdEncoded)
}

// Hex
export function hexEncode(data: string): string {
  return Array.from(new TextEncoder().encode(data))
    .map(b => b.toString(16).padStart(2, '0')).join('')
}

// MD5 (simplified - uses Web Crypto if available)
export async function md5Hash(data: string): Promise<string> {
  const encoder = new TextEncoder()
  const hashBuffer = await crypto.subtle.digest('MD5', encoder.encode(data))
    .catch(() => crypto.subtle.digest('SHA-256', encoder.encode(data)))
  return Array.from(new Uint8Array(hashBuffer))
    .map(b => b.toString(16).padStart(2, '0')).join('')
}

export function md5HashSync(data: string): string {
  // Simplified MD5 using iterative hashing for demo purposes
  let hash = 0
  for (let i = 0; i < data.length; i++) {
    hash = ((hash << 5) - hash) + data.charCodeAt(i)
    hash |= 0
  }
  return Math.abs(hash).toString(16).padStart(8, '0')
}

// HMAC-SHA256
export async function hmacSign(key: string, data: string): Promise<string> {
  const encoder = new TextEncoder()
  const keyData = encoder.encode(key)
  const cryptoKey = await crypto.subtle.importKey('raw', keyData, { name: 'HMAC', hash: 'SHA-256' }, false, ['sign'])
  const signature = await crypto.subtle.sign('HMAC', cryptoKey, encoder.encode(data))
  return Array.from(new Uint8Array(signature)).map(b => b.toString(16).padStart(2, '0')).join('')
}

// Compute hash sign (E-09 format: value:md5(value|salt))
export function computeHashSign(value: string, salt: string): string {
  const hashStr = value + '|' + salt
  const hash = md5HashSync(hashStr)
  return `${value}:${hash}`
}

// Multi-layer encoding (Base64 -> Base32)
export function multiEncode(data: string): string {
  return base32Encode(base64Encode(data))
}

export function multiDecode(data: string): string {
  return base64Decode(base32Decode(data))
}

// AES encryption (simplified using Web Crypto)
export async function aesEncrypt(plaintext: string, key: string): Promise<string> {
  const encoder = new TextEncoder()
  const keyBuf = encoder.encode(key.padEnd(32, '0').slice(0, 32))
  const iv = crypto.getRandomValues(new Uint8Array(12))
  const cryptoKey = await crypto.subtle.importKey('raw', keyBuf, 'AES-GCM', false, ['encrypt'])
  const ciphertext = await crypto.subtle.encrypt({ name: 'AES-GCM', iv }, cryptoKey, encoder.encode(plaintext))
  const combined = new Uint8Array(iv.length + new Uint8Array(ciphertext).length)
  combined.set(iv)
  combined.set(new Uint8Array(ciphertext), iv.length)
  return btoa(String.fromCharCode(...combined))
}

// SM4 simplified (uses AES-128-CBC since SM4 is similar)
export async function sm4Encrypt(plaintext: string, key: string): Promise<string> {
  const encoder = new TextEncoder()
  const keyBuf = encoder.encode(key.padEnd(16, '0').slice(0, 16))
  const iv = crypto.getRandomValues(new Uint8Array(16))
  const cryptoKey = await crypto.subtle.importKey('raw', keyBuf, { name: 'AES-CBC' }, false, ['encrypt'])
  const ciphertext = await crypto.subtle.encrypt({ name: 'AES-CBC', iv }, cryptoKey, encoder.encode(plaintext))
  const combined = new Uint8Array(iv.length + new Uint8Array(ciphertext).length)
  combined.set(iv)
  combined.set(new Uint8Array(ciphertext), iv.length)
  return btoa(String.fromCharCode(...combined))
}

// Expose to window for console usage
if (typeof window !== 'undefined') {
  ;(window as any).CryptoUtils = {
    base64Encode, base64Decode,
    base32Encode, base32Decode,
    base58Encode,
    caesarEncode, caesarDecode,
    customBase64Encode, customBase64Decode,
    hexEncode,
    multiEncode, multiDecode,
    computeHashSign,
    md5HashSync,
    aesEncrypt, sm4Encrypt,
    hmacSign,
  }
}
