import { deriveKey, secureClear } from './crypto.js';
import { DEFAULT_VAULT_CONFIG } from './types.js';
export const DEFAULT_CACHE_CONFIG = {
    maxAgeMs: 30 * 60 * 1000, // 30 minutes max age
    maxIdleMs: 5 * 60 * 1000, // 5 minutes idle timeout
    maxEntries: 10,
};
export class KeyCache {
    _cache = new Map();
    _config;
    _cleanupInterval = null;
    constructor(config = {}) {
        this._config = { ...DEFAULT_CACHE_CONFIG, ...config };
        this._startCleanupInterval();
    }
    async getKey(password, salt, vaultConfig = DEFAULT_VAULT_CONFIG) {
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
    invalidate(password, salt) {
        const cacheKey = this._generateCacheKey(password, salt);
        const entry = this._cache.get(cacheKey);
        if (entry) {
            secureClear(entry.key);
            this._cache.delete(cacheKey);
        }
    }
    invalidateAll() {
        for (const entry of this._cache.values()) {
            secureClear(entry.key);
        }
        this._cache.clear();
    }
    getStats() {
        let totalAccessCount = 0;
        let oldestEntry = null;
        let newestEntry = null;
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
    destroy() {
        if (this._cleanupInterval) {
            clearInterval(this._cleanupInterval);
            this._cleanupInterval = null;
        }
        this.invalidateAll();
    }
    _generateCacheKey(password, salt) {
        // Use synchronous crypto for cache key generation
        const { createHash } = require('crypto');
        return createHash('sha256').update(password).update(salt).digest('hex');
    }
    _isValid(entry) {
        const now = Date.now();
        const age = now - entry.createdAt;
        const idle = now - entry.lastAccessedAt;
        return age < this._config.maxAgeMs && idle < this._config.maxIdleMs;
    }
    _store(cacheKey, key, salt) {
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
    _evictLRU() {
        let oldestKey = null;
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
    _startCleanupInterval() {
        this._cleanupInterval = setInterval(() => {
            this._cleanup();
        }, 60000); // Run cleanup every minute
    }
    _cleanup() {
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
export function createKeyCache(config) {
    return new KeyCache(config);
}
//# sourceMappingURL=key-cache.js.map