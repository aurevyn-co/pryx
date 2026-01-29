import { Vault, createVault } from './vault.js';
import { KeyCache, createKeyCache } from './key-cache.js';
import { VaultConfig, VaultError } from './types.js';

export interface PasswordManagerConfig {
  vaultConfig?: Partial<VaultConfig>;
  cacheConfig?: Parameters<typeof createKeyCache>[0];
  autoLockMs?: number;
}

export const DEFAULT_PASSWORD_MANAGER_CONFIG = {
  autoLockMs: 5 * 60 * 1000,
};

export class PasswordManager {
  private _vault: Vault | null = null;
  private _keyCache: KeyCache;
  private _config: PasswordManagerConfig;
  private _autoLockTimer: NodeJS.Timeout | null = null;
  private _isLocked = true;

  constructor(config: PasswordManagerConfig = {}) {
    this._config = {
      ...DEFAULT_PASSWORD_MANAGER_CONFIG,
      ...config,
    };
    this._keyCache = createKeyCache(config.cacheConfig);
  }

  async unlock(password: string): Promise<void> {
    if (!this._isLocked) {
      return;
    }

    try {
      this._vault = await createVault(password, this._config.vaultConfig);
      this._isLocked = false;
      this._resetAutoLockTimer();
    } catch (error) {
      throw new VaultError('Failed to unlock vault: invalid password');
    }
  }

  lock(): void {
    if (this._vault) {
      this._vault.clearKey();
      this._vault = null;
    }
    this._keyCache.invalidateAll();
    this._isLocked = true;
    this._clearAutoLockTimer();
  }

  async encrypt(plaintext: Buffer): Promise<import('./types.js').EncryptedData> {
    this._ensureUnlocked();
    this._resetAutoLockTimer();
    return this._vault!.encrypt(plaintext);
  }

  async decrypt(encrypted: import('./types.js').EncryptedData): Promise<Buffer> {
    this._ensureUnlocked();
    this._resetAutoLockTimer();
    return this._vault!.decrypt(encrypted);
  }

  async changePassword(_oldPassword: string, newPassword: string): Promise<void> {
    this._ensureUnlocked();

    try {
      const newVault = await createVault(newPassword, this._config.vaultConfig);

      this._vault = newVault;
      this._keyCache.invalidateAll();
      this._resetAutoLockTimer();
    } catch (error) {
      throw new VaultError('Failed to change password');
    }
  }

  isUnlocked(): boolean {
    return !this._isLocked && this._vault !== null;
  }

  isLocked(): boolean {
    return this._isLocked;
  }

  getRemainingLockTime(): number | null {
    if (this._isLocked || !this._autoLockTimer) {
      return null;
    }
    return this._config.autoLockMs ?? DEFAULT_PASSWORD_MANAGER_CONFIG.autoLockMs;
  }

  destroy(): void {
    this.lock();
    this._keyCache.destroy();
  }

  private _ensureUnlocked(): void {
    if (this._isLocked || !this._vault) {
      throw new VaultError('Vault is locked. Call unlock() first.');
    }
  }

  private _resetAutoLockTimer(): void {
    this._clearAutoLockTimer();
    const timeout = this._config.autoLockMs ?? DEFAULT_PASSWORD_MANAGER_CONFIG.autoLockMs;
    if (timeout > 0) {
      this._autoLockTimer = setTimeout(() => {
        this.lock();
      }, timeout);
    }
  }

  private _clearAutoLockTimer(): void {
    if (this._autoLockTimer) {
      clearTimeout(this._autoLockTimer);
      this._autoLockTimer = null;
    }
  }
}

export function createPasswordManager(config?: PasswordManagerConfig): PasswordManager {
  return new PasswordManager(config);
}
