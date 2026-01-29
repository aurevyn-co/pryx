import { VaultConfig } from './types.js';
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
export declare const DEFAULT_CACHE_CONFIG: KeyCacheConfig;
export declare class KeyCache {
    private _cache;
    private _config;
    private _cleanupInterval;
    constructor(config?: Partial<KeyCacheConfig>);
    getKey(password: string, salt: Buffer, vaultConfig?: VaultConfig): Promise<Buffer>;
    invalidate(password: string, salt: Buffer): void;
    invalidateAll(): void;
    getStats(): {
        size: number;
        totalAccessCount: number;
        oldestEntry: number | null;
        newestEntry: number | null;
    };
    destroy(): void;
    private _generateCacheKey;
    private _isValid;
    private _store;
    private _evictLRU;
    private _startCleanupInterval;
    private _cleanup;
}
export declare function createKeyCache(config?: Partial<KeyCacheConfig>): KeyCache;
//# sourceMappingURL=key-cache.d.ts.map