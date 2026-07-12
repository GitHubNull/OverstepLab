/**
 * 前端编码工具集 - 用于 OverstepLab 编码加密挑战
 *
 * 提供与后端 crypto 包对应的常见编码/加密功能。
 * 用户可以在浏览器控制台中使用这些函数进行编码解码操作。
 */

import CryptoJS from 'crypto-js'
import baseX from 'base-x'
import { sm4 } from 'sm-crypto'
import { JSEncrypt } from 'jsencrypt'

// ==================== Base64 ====================

export function base64Encode(data: string): string {
  return CryptoJS.enc.Base64.stringify(CryptoJS.enc.Utf8.parse(data))
}

export function base64Decode(data: string): string {
  return CryptoJS.enc.Base64.parse(data).toString(CryptoJS.enc.Utf8)
}

// ==================== Base32 (RFC 4648) ====================

const BASE32_ALPHABET = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ234567'

export function base32Encode(data: string): string {
  const bytes = new TextEncoder().encode(data)
  let bits = 0
  let value = 0
  let result = ''
  for (let i = 0; i < bytes.length; i++) {
    value = (value << 8) | bytes[i]
    bits += 8
    while (bits >= 5) {
      result += BASE32_ALPHABET[(value >>> (bits - 5)) & 31]
      bits -= 5
    }
  }
  if (bits > 0) {
    result += BASE32_ALPHABET[(value << (5 - bits)) & 31]
  }
  return result
}

export function base32Decode(data: string): string {
  data = data.toUpperCase().replace(/=+$/, '')
  let bits = 0
  let value = 0
  const bytes: number[] = []
  for (let i = 0; i < data.length; i++) {
    const idx = BASE32_ALPHABET.indexOf(data[i])
    if (idx === -1) continue
    value = (value << 5) | idx
    bits += 5
    if (bits >= 8) {
      bytes.push((value >>> (bits - 8)) & 0xff)
      bits -= 8
    }
  }
  return new TextDecoder().decode(new Uint8Array(bytes))
}

// ==================== Hex ====================

export function hexEncode(data: string): string {
  return CryptoJS.enc.Hex.stringify(CryptoJS.enc.Utf8.parse(data))
}

export function hexDecode(data: string): string {
  return CryptoJS.enc.Hex.parse(data).toString(CryptoJS.enc.Utf8)
}

// ==================== Base58 (Bitcoin-style) ====================

const BASE58_ALPHABET = '123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz'
const base58Codec = baseX(BASE58_ALPHABET)

export function base58Encode(data: string): string {
  return base58Codec.encode(new TextEncoder().encode(data))
}

export function base58Decode(data: string): string {
  return new TextDecoder().decode(base58Codec.decode(data))
}

// ==================== Base85 (Z85 variant) ====================

const BASE85_ALPHABET = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!#$%&()*+-;<=>?@^_`{|}~'
const base85Codec = baseX(BASE85_ALPHABET)

export function base85Encode(data: string): string {
  return base85Codec.encode(new TextEncoder().encode(data))
}

export function base85Decode(data: string): string {
  return new TextDecoder().decode(base85Codec.decode(data))
}

// ==================== Custom Base64 (swapped charset: A<->Z, a<->z, 0<->9, +<->/) ====================

const CUSTOM_BASE64_CHARSET = 'ZYXWVUTSRQPONMLKJIHGFEDCBAzyxwvutsrqponmlkjihgfedcba9876543210/+'
const customBase64Codec = baseX(CUSTOM_BASE64_CHARSET)

export function customBase64Encode(data: string): string {
  return customBase64Codec.encode(new TextEncoder().encode(data))
}

export function customBase64Decode(data: string): string {
  return new TextDecoder().decode(customBase64Codec.decode(data))
}

// ==================== Custom Base32 (reversed RFC 4648 charset) ====================

const CUSTOM_BASE32_CHARSET = 'ZYXWVUTSRQPONMLKJIHGFEDCBA765432'
const customBase32Codec = baseX(CUSTOM_BASE32_CHARSET)

export function customBase32Encode(data: string): string {
  return customBase32Codec.encode(new TextEncoder().encode(data))
}

export function customBase32Decode(data: string): string {
  return new TextDecoder().decode(customBase32Codec.decode(data))
}

// ==================== Caesar Cipher ====================

export function caesarEncode(data: string, shift: number = 3): string {
  let result = ''
  for (const c of data) {
    if (c >= '0' && c <= '9') {
      result += String.fromCharCode(48 + ((c.charCodeAt(0) - 48 + shift) % 10))
    } else if (c >= 'A' && c <= 'Z') {
      result += String.fromCharCode(65 + ((c.charCodeAt(0) - 65 + shift) % 26))
    } else if (c >= 'a' && c <= 'z') {
      result += String.fromCharCode(97 + ((c.charCodeAt(0) - 97 + shift) % 26))
    } else {
      result += c
    }
  }
  return result
}

export function caesarDecode(data: string, shift: number = 3): string {
  const reverseShift = (26 - (shift % 26)) % 26
  return caesarEncode(data, reverseShift)
}

// ==================== Multi-layer Encoding (Base64 -> Base32) ====================

export function multiEncode(data: string): string {
  return base32Encode(base64Encode(data))
}

export function multiDecode(data: string): string {
  return base64Decode(base32Decode(data))
}

// ==================== AES (crypto-js) ====================

// Note: These use a different key format than the backend!
// Use the backend /api/v1/crypto/keys endpoint to get the actual key.
export function aesEncrypt(data: string, key: string): string {
  return CryptoJS.AES.encrypt(data, key).toString()
}

export function aesDecrypt(ciphertext: string, key: string): string {
  const bytes = CryptoJS.AES.decrypt(ciphertext, key)
  return bytes.toString(CryptoJS.enc.Utf8)
}

// ==================== SM4 (国密对称加密) ====================

export function sm4Encrypt(data: string, key: string): string {
  // sm-crypto expects hex key and hex data
  const keyHex = CryptoJS.enc.Hex.stringify(CryptoJS.enc.Utf8.parse(key))
  const dataHex = CryptoJS.enc.Hex.stringify(CryptoJS.enc.Utf8.parse(data))
  // sm4.encrypt returns hex string
  return sm4.encrypt(dataHex, keyHex, 1) // 1 = ECB mode with PKCS7 padding
}

export function sm4Decrypt(ciphertext: string, key: string): string {
  const keyHex = CryptoJS.enc.Hex.stringify(CryptoJS.enc.Utf8.parse(key))
  const decryptedHex = sm4.decrypt(ciphertext, keyHex, 1)
  return CryptoJS.enc.Hex.parse(decryptedHex).toString(CryptoJS.enc.Utf8)
}

// ==================== RSA (jsencrypt) ====================

export function rsaEncrypt(data: string, publicKey: string): string {
  const encryptor = new JSEncrypt()
  encryptor.setPublicKey(publicKey)
  const result = encryptor.encrypt(data)
  if (!result) {
    throw new Error('RSA encryption failed')
  }
  return result
}

export function rsaDecrypt(ciphertext: string, privateKey: string): string {
  const decryptor = new JSEncrypt()
  decryptor.setPrivateKey(privateKey)
  const result = decryptor.decrypt(ciphertext)
  if (!result) {
    throw new Error('RSA decryption failed')
  }
  return result
}

// ==================== SHA256 / HMAC ====================

export function sha256(data: string): string {
  return CryptoJS.SHA256(data).toString()
}

export function hmacSha256(data: string, key: string): string {
  return CryptoJS.HmacSHA256(data, key).toString()
}

// ==================== Signed Param ====================

export function encodeSignedParam(value: string, hmacKey: string): string {
  const encoded = base64Encode(value)
  const sig = hmacSha256(value, hmacKey)
  return encoded + '.' + sig
}

// ==================== Vigenere Cipher ====================

export function vigenereEncode(data: string, key: string): string {
  if (!key) return data
  const upperKey = key.toUpperCase()
  let result = ''
  let keyIdx = 0
  for (const c of data) {
    if (c >= 'A' && c <= 'Z') {
      const shift = upperKey.charCodeAt(keyIdx % key.length) - 65
      result += String.fromCharCode(65 + ((c.charCodeAt(0) - 65 + shift) % 26))
      keyIdx++
    } else if (c >= 'a' && c <= 'z') {
      const shift = upperKey.charCodeAt(keyIdx % key.length) - 65
      result += String.fromCharCode(97 + ((c.charCodeAt(0) - 97 + shift) % 26))
      keyIdx++
    } else {
      result += c
    }
  }
  return result
}

// ==================== 编码类型检测辅助 ====================

export function detectEncoding(data: string): string {
  // Try Base64
  if (/^[A-Za-z0-9+/]+=*$/.test(data) && data.length >= 4) {
    try {
      base64Decode(data)
      return 'base64'
    } catch { /* continue */ }
  }
  // Try Base32
  if (/^[A-Z2-7]+=*$/i.test(data) && data.length >= 4) {
    return 'base32'
  }
  // Try Hex
  if (/^[0-9A-Fa-f]+$/.test(data) && data.length % 2 === 0 && data.length >= 2) {
    return 'hex'
  }
  return 'unknown'
}

// ==================== 便捷导出对象 ====================

export const CryptoUtils = {
  base64Encode,
  base64Decode,
  base32Encode,
  base32Decode,
  base58Encode,
  base58Decode,
  base85Encode,
  base85Decode,
  hexEncode,
  hexDecode,
  customBase64Encode,
  customBase64Decode,
  customBase32Encode,
  customBase32Decode,
  caesarEncode,
  caesarDecode,
  multiEncode,
  multiDecode,
  aesEncrypt,
  aesDecrypt,
  sm4Encrypt,
  sm4Decrypt,
  rsaEncrypt,
  rsaDecrypt,
  sha256,
  hmacSha256,
  encodeSignedParam,
  vigenereEncode,
  detectEncoding,
}

// 挂载到 window 对象，方便在浏览器 Console 中直接使用
if (typeof window !== 'undefined') {
  ;(window as any).CryptoUtils = CryptoUtils
}
