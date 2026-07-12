// 前端加解密工具库
// Base64, Base32, Base58, AES, SM4, RSA 等编码/加密/签名工具

// True MD5 implementation using SparkMD5-like algorithm
function md5Cycle(x: number[], k: number[]): number[] {
  let a = x[0], b = x[1], c = x[2], d = x[3]

  const ff = (a: number, b: number, c: number, d: number, x: number, s: number, t: number) => {
    const n = a + (b & c | ~b & d) + x + t
    return ((n << s) | (n >>> (32 - s))) + b
  }
  const gg = (a: number, b: number, c: number, d: number, x: number, s: number, t: number) => {
    const n = a + (b & d | c & ~d) + x + t
    return ((n << s) | (n >>> (32 - s))) + b
  }
  const hh = (a: number, b: number, c: number, d: number, x: number, s: number, t: number) => {
    const n = a + (b ^ c ^ d) + x + t
    return ((n << s) | (n >>> (32 - s))) + b
  }
  const ii = (a: number, b: number, c: number, d: number, x: number, s: number, t: number) => {
    const n = a + (c ^ (b | ~d)) + x + t
    return ((n << s) | (n >>> (32 - s))) + b
  }

  a = ff(a, b, c, d, k[0], 7, -680876936)
  d = ff(d, a, b, c, k[1], 12, -389564586)
  c = ff(c, d, a, b, k[2], 17, 606105819)
  b = ff(b, c, d, a, k[3], 22, -1044525330)
  a = ff(a, b, c, d, k[4], 7, -176418897)
  d = ff(d, a, b, c, k[5], 12, 1200080426)
  c = ff(c, d, a, b, k[6], 17, -1473231341)
  b = ff(b, c, d, a, k[7], 22, -45705983)
  a = ff(a, b, c, d, k[8], 7, 1770035416)
  d = ff(d, a, b, c, k[9], 12, -1958414417)
  c = ff(c, d, a, b, k[10], 17, -42063)
  b = ff(b, c, d, a, k[11], 22, -1990404162)
  a = ff(a, b, c, d, k[12], 7, 1804603682)
  d = ff(d, a, b, c, k[13], 12, -40341101)
  c = ff(c, d, a, b, k[14], 17, -1502002290)
  b = ff(b, c, d, a, k[15], 22, 1236535329)

  a = gg(a, b, c, d, k[1], 5, -165796510)
  d = gg(d, a, b, c, k[6], 9, -1069501632)
  c = gg(c, d, a, b, k[11], 14, 643717713)
  b = gg(b, c, d, a, k[0], 20, -373897302)
  a = gg(a, b, c, d, k[5], 5, -701558691)
  d = gg(d, a, b, c, k[10], 9, 38016083)
  c = gg(c, d, a, b, k[15], 14, -660478335)
  b = gg(b, c, d, a, k[4], 20, -405537848)
  a = gg(a, b, c, d, k[9], 5, 568446438)
  d = gg(d, a, b, c, k[14], 9, -1019803690)
  c = gg(c, d, a, b, k[3], 14, -187363961)
  b = gg(b, c, d, a, k[8], 20, 1163531501)
  a = gg(a, b, c, d, k[13], 5, -1444681467)
  d = gg(d, a, b, c, k[2], 9, -51403784)
  c = gg(c, d, a, b, k[7], 14, 1735328473)
  b = gg(b, c, d, a, k[12], 20, -1926607734)

  a = hh(a, b, c, d, k[5], 4, -378558)
  d = hh(d, a, b, c, k[8], 11, -2022574463)
  c = hh(c, d, a, b, k[11], 16, 1839030562)
  b = hh(b, c, d, a, k[14], 23, -35309556)
  a = hh(a, b, c, d, k[1], 4, -1530992060)
  d = hh(d, a, b, c, k[4], 11, 1272893353)
  c = hh(c, d, a, b, k[7], 16, -155497632)
  b = hh(b, c, d, a, k[10], 23, -1094730640)
  a = hh(a, b, c, d, k[13], 4, 681279174)
  d = hh(d, a, b, c, k[0], 11, -358537222)
  c = hh(c, d, a, b, k[3], 16, -722521979)
  b = hh(b, c, d, a, k[6], 23, 76029189)
  a = hh(a, b, c, d, k[9], 4, -640364487)
  d = hh(d, a, b, c, k[12], 11, -421815835)
  c = hh(c, d, a, b, k[15], 16, 530742520)
  b = hh(b, c, d, a, k[2], 23, -995338651)

  a = ii(a, b, c, d, k[0], 6, -198630844)
  d = ii(d, a, b, c, k[7], 10, 1126891415)
  c = ii(c, d, a, b, k[14], 15, -1416354905)
  b = ii(b, c, d, a, k[5], 21, -57434055)
  a = ii(a, b, c, d, k[12], 6, 1700485571)
  d = ii(d, a, b, c, k[3], 10, -1894986606)
  c = ii(c, d, a, b, k[10], 15, -1051523)
  b = ii(b, c, d, a, k[1], 21, -2054922799)
  a = ii(a, b, c, d, k[8], 6, 1873313359)
  d = ii(d, a, b, c, k[15], 10, -30611744)
  c = ii(c, d, a, b, k[6], 15, -1560198380)
  b = ii(b, c, d, a, k[13], 21, 1309151649)
  a = ii(a, b, c, d, k[4], 6, -145523070)
  d = ii(d, a, b, c, k[11], 10, -1120210379)
  c = ii(c, d, a, b, k[2], 15, 718787259)
  b = ii(b, c, d, a, k[9], 21, -343485551)

  x[0] = a + x[0] | 0
  x[1] = b + x[1] | 0
  x[2] = c + x[2] | 0
  x[3] = d + x[3] | 0
  return x
}

function md5ToHex(arr: number[]): string {
  const hex = '0123456789abcdef'
  let str = ''
  for (let i = 0; i < arr.length * 4; i++) {
    str += hex.charAt((arr[i >> 2] >> ((i % 4) * 8 + 4)) & 0xF) +
           hex.charAt((arr[i >> 2] >> ((i % 4) * 8)) & 0xF)
  }
  return str
}

function md5FromBytes(bytes: Uint8Array): number[] {
  const len = bytes.length
  const tail = len % 64
  const padLen = tail < 56 ? 56 - tail : 120 - tail
  const totalLen = len + padLen + 8
  const buf = new Uint8Array(totalLen)
  buf.set(bytes)
  buf[len] = 0x80
  const bitLen = BigInt(len) * BigInt(8)
  const view = new DataView(buf.buffer, totalLen - 8)
  view.setUint32(0, Number(bitLen & BigInt(0xFFFFFFFF)), true)
  view.setUint32(4, Number(bitLen >> BigInt(32)), true)

  const state = [0x67452301, -0x10325477, -0x67452302, 0x10325476]
  const block = new Int32Array(16)
  for (let i = 0; i < totalLen; i += 64) {
    const chunk = new DataView(buf.buffer, i, 64)
    for (let j = 0; j < 16; j++) {
      block[j] = chunk.getInt32(j * 4, true)
    }
    md5Cycle(state, Array.from(block))
  }
  return state
}

export function md5HashSync(data: string): string {
  const encoder = new TextEncoder()
  const bytes = encoder.encode(data)
  const state = md5FromBytes(bytes)
  return md5ToHex(state)
}

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

// MD5 (true implementation - synchronous)
export async function md5Hash(data: string): Promise<string> {
  // Use the true MD5 implementation (sync)
  return md5HashSync(data)
}

// md5HashSync is defined at the top of the file

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
  // Derive 32-byte key using SHA-256 (consistent with backend)
  const keyHash = await crypto.subtle.digest('SHA-256', encoder.encode(key))
  const keyBuf = new Uint8Array(keyHash).slice(0, 32)
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
  // Derive 16-byte key using SHA-256 (consistent with backend)
  const keyHash = await crypto.subtle.digest('SHA-256', encoder.encode(key))
  const keyBuf = new Uint8Array(keyHash).slice(0, 16)
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
