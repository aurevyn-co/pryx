export { DEFAULT_VAULT_CONFIG, AES_GCM_IV_LENGTH, AES_GCM_TAG_LENGTH, CURRENT_VERSION, VaultError, InvalidPasswordError, CorruptedDataError, DecryptionError, } from './types.js';
export { deriveKey, generateSalt, generateIV, encrypt, decrypt, secureClear, secureCompare, serializeEncryptedData, deserializeEncryptedData, } from './crypto.js';
export { Vault, createVault, encryptWithPassword, decryptWithPassword, } from './vault.js';
export { KeyCache, DEFAULT_CACHE_CONFIG, createKeyCache, } from './key-cache.js';
export { PasswordManager, DEFAULT_PASSWORD_MANAGER_CONFIG, createPasswordManager, } from './password-manager.js';
//# sourceMappingURL=index.js.map