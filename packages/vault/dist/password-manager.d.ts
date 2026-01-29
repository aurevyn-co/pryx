import { createKeyCache } from './key-cache.js';
import { VaultConfig } from './types.js';
export interface PasswordManagerConfig {
    vaultConfig?: Partial<VaultConfig>;
    cacheConfig?: Parameters<typeof createKeyCache>[0];
    autoLockMs?: number;
}
export declare const DEFAULT_PASSWORD_MANAGER_CONFIG: {
    autoLockMs: number;
};
export declare class PasswordManager {
    private _vault;
    private _keyCache;
    private _config;
    private _autoLockTimer;
    private _isLocked;
    constructor(config?: PasswordManagerConfig);
    unlock(password: string): Promise<void>;
    lock(): void;
    encrypt(plaintext: Buffer): Promise<import('./types.js').EncryptedData>;
    decrypt(encrypted: import('./types.js').EncryptedData): Promise<Buffer>;
    changePassword(_oldPassword: string, newPassword: string): Promise<void>;
    isUnlocked(): boolean;
    isLocked(): boolean;
    getRemainingLockTime(): number | null;
    destroy(): void;
    private _ensureUnlocked;
    private _resetAutoLockTimer;
    private _clearAutoLockTimer;
}
export declare function createPasswordManager(config?: PasswordManagerConfig): PasswordManager;
//# sourceMappingURL=password-manager.d.ts.map