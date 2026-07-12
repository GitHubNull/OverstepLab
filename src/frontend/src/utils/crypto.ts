// 前端加解密工具库
// Base64, Base32, Base58, AES, SM4, RSA 等编码/加密/签名工具

// True MD5 implementation (based on well-tested algorithm)
function safeAdd(x: number, y: number): number {
  const lsw = (x & 0xFFFF) + (y & 0xFFFF)
  const msw = (x >> 16) + (y >> 16) + (lsw >> 16)
  return (msw << 16) | (lsw & 0xFFFF)
}

function bitRotateLeft(num: number, cnt: number): number {
  return (num << cnt) | (num >>> (32 - cnt))
}

function md5cmn(q: number, a: number, b: number, x: number, s: number, t: number): number {
  return safeAdd(bitRotateLeft(safeAdd(safeAdd(a, q), safeAdd(x, t)), s), b)
}

function md5ff(a: number, b: number, c: number, d: number, x: number, s: number, t: number): number {
  return md5cmn((b & c) | (~b & d), a, b, x, s, t)
}

function md5gg(a: number, b: number, c: number, d: number, x: number, s: number, t: number): number {
  return md5cmn((b & d) | (c & ~d), a, b, x, s, t)
}

function md5hh(a: number, b: number, c: number, d: number, x: number, s: number, t: number): number {
  return md5cmn(b ^ c ^ d, a, b, x, s, t)
}

function md5ii(a: number, b: number, c: number, d: number, x: number, s: number, t: number): number {
  return md5cmn(c ^ (b | ~d), a, b, x, s, t)
}

function md5Binl(x: number[], len: number): number[] {
  x[len >> 5] |= 0x80 << (len % 32)
  x[(((len + 64) >>> 9) << 4) + 14] = len

  let i: number
  let olda: number, oldb: number, oldc: number, oldd: number
  let a = 1732584193, b = -271733879, c = -1732584194, d = 271733878

  for (i = 0; i < x.length; i += 16) {
    olda = a; oldb = b; oldc = c; oldd = d

    a = md5ff(a, b, c, d, x[i], 7, -680876936)
    d = md5ff(d, a, b, c, x[i + 1], 12, -389564586)
    c = md5ff(c, d, a, b, x[i + 2], 17, 606105819)
    b = md5ff(b, c, d, a, x[i + 3], 22, -1044525330)
    a = md5ff(a, b, c, d, x[i + 4], 7, -176418897)
    d = md5ff(d, a, b, c, x[i + 5], 12, 1200080426)
    c = md5ff(c, d, a, b, x[i + 6], 17, -1473231341)
    b = md5ff(b, c, d, a, x[i + 7], 22, -45705983)
    a = md5ff(a, b, c, d, x[i + 8], 7, 1770035416)
    d = md5ff(d, a, b, c, x[i + 9], 12, -1958414417)
    c = md5ff(c, d, a, b, x[i + 10], 17, -42063)
    b = md5ff(b, c, d, a, x[i + 11], 22, -1990404162)
    a = md5ff(a, b, c, d, x[i + 12], 7, 1804603682)
    d = md5ff(d, a, b, c, x[i + 13], 12, -40341101)
    c = md5ff(c, d, a, b, x[i + 14], 17, -1502002290)
    b = md5ff(b, c, d, a, x[i + 15], 22, 1236535329)

    a = md5gg(a, b, c, d, x[i + 1], 5, -165796510)
    d = md5gg(d, a, b, c, x[i + 6], 9, -1069501632)
    c = md5gg(c, d, a, b, x[i + 11], 14, 643717713)
    b = md5gg(b, c, d, a, x[i], 20, -373897302)
    a = md5gg(a, b, c, d, x[i + 5], 5, -701558691)
    d = md5gg(d, a, b, c, x[i + 10], 9, 38016083)
    c = md5gg(c, d, a, b, x[i + 15], 14, -660478335)
    b = md5gg(b, c, d, a, x[i + 4], 20, -405537848)
    a = md5gg(a, b, c, d, x[i + 9], 5, 568446438)
    d = md5gg(d, a, b, c, x[i + 14], 9, -1019803690)
    c = md5gg(c, d, a, b, x[i + 3], 14, -187363961)
    b = md5gg(b, c, d, a, x[i + 8], 20, 1163531501)
    a = md5gg(a, b, c, d, x[i + 13], 5, -1444681467)
    d = md5gg(d, a, b, c, x[i + 2], 9, -51403784)
    c = md5gg(c, d, a, b, x[i + 7], 14, 1735328473)
    b = md5gg(b, c, d, a, x[i + 12], 20, -1926607734)

    a = md5hh(a, b, c, d, x[i + 5], 4, -378558)
    d = md5hh(d, a, b, c, x[i + 8], 11, -2022574463)
    c = md5hh(c, d, a, b, x[i + 11], 16, 1839030562)
    b = md5hh(b, c, d, a, x[i + 14], 23, -35309556)
    a = md5hh(a, b, c, d, x[i + 1], 4, -1530992060)
    d = md5hh(d, a, b, c, x[i + 4], 11, 1272893353)
    c = md5hh(c, d, a, b, x[i + 7], 16, -155497632)
    b = md5hh(b, c, d, a, x[i + 10], 23, -1094730640)
    a = md5hh(a, b, c, d, x[i + 13], 4, 681279174)
    d = md5hh(d, a, b, c, x[i], 11, -358537222)
    c = md5hh(c, d, a, b, x[i + 3], 16, -722521979)
    b = md5hh(b, c, d, a, x[i + 6], 23, 76029189)
    a = md5hh(a, b, c, d, x[i + 9], 4, -640364487)
    d = md5hh(d, a, b, c, x[i + 12], 11, -421815835)
    c = md5hh(c, d, a, b, x[i + 15], 16, 530742520)
    b = md5hh(b, c, d, a, x[i + 2], 23, -995338651)

    a = md5ii(a, b, c, d, x[i], 6, -198630844)
    d = md5ii(d, a, b, c, x[i + 7], 10, 1126891415)
    c = md5ii(c, d, a, b, x[i + 14], 15, -1416354905)
    b = md5ii(b, c, d, a, x[i + 5], 21, -57434055)
    a = md5ii(a, b, c, d, x[i + 12], 6, 1700485571)
    d = md5ii(d, a, b, c, x[i + 3], 10, -1894986606)
    c = md5ii(c, d, a, b, x[i + 10], 15, -1051523)
    b = md5ii(b, c, d, a, x[i + 1], 21, -2054922799)
    a = md5ii(a, b, c, d, x[i + 8], 6, 1873313359)
    d = md5ii(d, a, b, c, x[i + 15], 10, -30611744)
    c = md5ii(c, d, a, b, x[i + 6], 15, -1560198380)
    b = md5ii(b, c, d, a, x[i + 13], 21, 1309151649)
    a = md5ii(a, b, c, d, x[i + 4], 6, -145523070)
    d = md5ii(d, a, b, c, x[i + 11], 10, -1120210379)
    c = md5ii(c, d, a, b, x[i + 2], 15, 718787259)
    b = md5ii(b, c, d, a, x[i + 9], 21, -343485551)

    a = safeAdd(a, olda)
    b = safeAdd(b, oldb)
    c = safeAdd(c, oldc)
    d = safeAdd(d, oldd)
  }
  return [a, b, c, d]
}

function md5Binl2Hex(binarray: number[]): string {
  const hexTab = '0123456789abcdef'
  let str = ''
  for (let i = 0; i < binarray.length * 4; i++) {
    str += hexTab.charAt((binarray[i >> 2] >> ((i % 4) * 8 + 4)) & 0xF) +
           hexTab.charAt((binarray[i >> 2] >> ((i % 4) * 8)) & 0xF)
  }
  return str
}

function md5Raw(input: string): string {
  const encoder = new TextEncoder()
  const bytes = encoder.encode(input)
  // Convert bytes to little-endian 32-bit words
  const bin: number[] = []
  for (let i = 0; i < bytes.length * 8; i += 8) {
    bin[i >> 5] |= (bytes[i / 8] & 0xFF) << (i % 32)
  }
  const hash = md5Binl(bin, bytes.length * 8)
  return md5Binl2Hex(hash)
}

export function md5HashSync(data: string): string {
  return md5Raw(data)
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
