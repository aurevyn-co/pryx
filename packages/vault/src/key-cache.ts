import { deriveKey, secureClear } from './crypto.js';
import { VaultConfig, DEFAULT_VAULT_CONFIG } from './types.js';

export interface KeyCacheEntry {
  key: Buffer;
  salt: Buffer;
  createdAt: number;
  lastAccessedAt: number;
  accessCount: number;
}

export interface KeyCacheConfig {
  maxAgeMs: number;
  maxIdleMs: number;
  maxEntries: number;
}

export const DEFAULT_CACHE_CONFIG: KeyCacheConfig = {
  maxAgeMs: 30 * 60 * 1000, // 30 minutes max age
  maxIdleMs: 5 * 60 * 1000,  // 5 minutes idle timeout
  maxEntries: 10,
};

export class KeyCache {
  private _cache: Map<string, KeyCacheEntry> = new Map();
  private _config: KeyCacheConfig;
  private _cleanupInterval: NodeJS.Timeout | null = null;

  constructor(config: Partial<KeyCacheConfig> = {}) {
    this._config = { ...DEFAULT_CACHE_CONFIG, ...config };
    this._startCleanupInterval();
  }

  async getKey(
    password: string,
    salt: Buffer,
    vaultConfig: VaultConfig = DEFAULT_VAULT_CONFIG
  ): Promise<Buffer> {
    const cacheKey = this._generateCacheKey(password, salt);
    const entry = this._cache.get(cacheKey);

    if (entry && this._isValid(entry)) {
      entry.lastAccessedAt = Date.now();
      entry.accessCount++;
      return Buffer.from(entry.key);
    }

    const key = await deriveKey(password, salt, vaultConfig);
    this._store(cacheKey, key, salt);
    return key;
  }

  invalidate(password: string, salt: Buffer): void {
    const cacheKey = this._generateCacheKey(password, salt);
    const entry = this._cache.get(cacheKey);
    if (entry) {
      secureClear(entry.key);
      this._cache.delete(cacheKey);
    }
  }

  invalidateAll(): void {
    for (const entry of this._cache.values()) {
      secureClear(entry.key);
    }
    this._cache.clear();
  }

  getStats(): {
    size: number;
    totalAccessCount: number;
    oldestEntry: number | null;
    newestEntry: number | null;
  } {
    let totalAccessCount = 0;
    let oldestEntry: number | null = null;
    let newestEntry: number | null = null;

    for (const entry of this._cache.values()) {
      totalAccessCount += entry.accessCount;
      if (!oldestEntry || entry.createdAt < oldestEntry) {
        oldestEntry = entry.createdAt;
      }
      if (!newestEntry || entry.createdAt > newestEntry) {
        newestEntry = entry.createdAt;
      }
    }

    return {
      size: this._cache.size,
      totalAccessCount,
      oldestEntry,
      newestEntry,
    };
  }

  destroy(): void {
    if (this._cleanupInterval) {
      clearInterval(this._cleanupInterval);
      this._cleanupInterval = null;
    }
    this.invalidateAll();
  }

  private _generateCacheKey(password: string, salt: Buffer): string {
    // Use synchronous crypto for cache key generation
    const { createHash } = require('crypto');
    return createHash('sha256').update(password).update(salt).digest('hex');
  }

  private _isValid(entry: KeyCacheEntry): boolean {
    const now = Date.now();
    const age = now - entry.createdAt;
    const idle = now - entry.lastAccessedAt;

    return age < this._config.maxAgeMs && idle < this._config.maxIdleMs;
  }

  private _store(cacheKey: string, key: Buffer, salt: Buffer): void {
    if (this._cache.size >= this._config.maxEntries) {
      this._evictLRU();
    }

    this._cache.set(cacheKey, {
      key: Buffer.from(key),
      salt: Buffer.from(salt),
      createdAt: Date.now(),
      lastAccessedAt: Date.now(),
      accessCount: 1,
    });
  }

  private _evictLRU(): void {
    let oldestKey: string | null = null;
    let oldestAccess = Infinity;

    for (const [key, entry] of this._cache.entries()) {
      if (entry.lastAccessedAt < oldestAccess) {
        oldestAccess = entry.lastAccessedAt;
        oldestKey = key;
      }
    }

    if (oldestKey) {
      const entry = this._cache.get(oldestKey);
      if (entry) {
        secureClear(entry.key);
      }
      this._cache.delete(oldestKey);
    }
  }

  private _startCleanupInterval(): void {
    this._cleanupInterval = setInterval(() => {
      this._cleanup();
    }, 60000); // Run cleanup every minute
  }

  private _cleanup(): void {
    const now = Date.now();
    for (const [key, entry] of this._cache.entries()) {
      const age = now - entry.createdAt;
      const idle = now - entry.lastAccessedAt;

      if (age >= this._config.maxAgeMs || idle >= this._config.maxIdleMs) {
        secureClear(entry.key);
        this._cache.delete(key);
      }
    }
  }
}

export function createKeyCache(config?: Partial<KeyCacheConfig>): KeyCache {
  return new KeyCache(config);
}
