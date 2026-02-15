import { createVault } from './vault.js';
import { createKeyCache } from './key-cache.js';
import { VaultError } from './types.js';
export const DEFAULT_PASSWORD_MANAGER_CONFIG = {
    autoLockMs: 5 * 60 * 1000,
};
export class PasswordManager {
    _vault = null;
    _keyCache;
    _config;
    _autoLockTimer = null;
    _isLocked = true;
    constructor(config = {}) {
        this._config = {
            ...DEFAULT_PASSWORD_MANAGER_CONFIG,
            ...config,
        };
        this._keyCache = createKeyCache(config.cacheConfig);
    }
    async unlock(password) {
        if (!this._isLocked) {
            return;
        }
        try {
            this._vault = await createVault(password, this._config.vaultConfig);
            this._isLocked = false;
            this._resetAutoLockTimer();
        }
        catch (error) {
            throw new VaultError('Failed to unlock vault: invalid password');
        }
    }
    lock() {
        if (this._vault) {
            this._vault.clearKey();
            this._vault = null;
        }
        this._keyCache.invalidateAll();
        this._isLocked = true;
        this._clearAutoLockTimer();
    }
    async encrypt(plaintext) {
        this._ensureUnlocked();
        this._resetAutoLockTimer();
        return this._vault.encrypt(plaintext);
    }
    async decrypt(encrypted) {
        this._ensureUnlocked();
        this._resetAutoLockTimer();
        return this._vault.decrypt(encrypted);
    }
    async changePassword(_oldPassword, newPassword) {
        this._ensureUnlocked();
        try {
            const newVault = await createVault(newPassword, this._config.vaultConfig);
            this._vault = newVault;
            this._keyCache.invalidateAll();
            this._resetAutoLockTimer();
        }
        catch (error) {
            throw new VaultError('Failed to change password');
        }
    }
    isUnlocked() {
        return !this._isLocked && this._vault !== null;
    }
    isLocked() {
        return this._isLocked;
    }
    getRemainingLockTime() {
        if (this._isLocked || !this._autoLockTimer) {
            return null;
        }
        return this._config.autoLockMs ?? DEFAULT_PASSWORD_MANAGER_CONFIG.autoLockMs;
    }
    destroy() {
        this.lock();
        this._keyCache.destroy();
    }
    _ensureUnlocked() {
        if (this._isLocked || !this._vault) {
            throw new VaultError('Vault is locked. Call unlock() first.');
        }
    }
    _resetAutoLockTimer() {
        this._clearAutoLockTimer();
        const timeout = this._config.autoLockMs ?? DEFAULT_PASSWORD_MANAGER_CONFIG.autoLockMs;
        if (timeout > 0) {
            this._autoLockTimer = setTimeout(() => {
                this.lock();
            }, timeout);
        }
    }
    _clearAutoLockTimer() {
        if (this._autoLockTimer) {
            clearTimeout(this._autoLockTimer);
            this._autoLockTimer = null;
        }
    }
}
export function createPasswordManager(config) {
    return new PasswordManager(config);
}
//# sourceMappingURL=password-manager.js.map