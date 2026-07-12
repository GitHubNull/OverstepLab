declare module 'sm-crypto' {
  export const sm4: {
    encrypt: (data: string, key: string, mode: number) => string
    decrypt: (data: string, key: string, mode: number) => string
  }
  export const sm2: {
    generateKeyPair: () => { privateKey: string; publicKey: string }
    encrypt: (data: string, publicKey: string, cipherMode?: number) => string
    decrypt: (data: string, privateKey: string, cipherMode?: number) => string
  }
}
