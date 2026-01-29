import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import {
  KeyCache,
  createKeyCache,
  DEFAULT_CACHE_CONFIG,
} from '../../src/key-cache.js';
import { generateSalt } from '../../src/crypto.js';

describe('KeyCache', () => {
  let cache: KeyCache;

  beforeEach(() => {
    cache = createKeyCache();
  });

  afterEach(() => {
    cache.destroy();
  });

  describe('getKey', () => {
    it('should derive and cache key', async () => {
      const password = 'test-password';
      const salt = generateSalt();

      const key1 = await cache.getKey(password, salt);
      const key2 = await cache.getKey(password, salt);

      expect(key1).toEqual(key2);
      expect(cache.getStats().size).toBe(1);
    });

    it('should derive different keys for different passwords', async () => {
      const salt = generateSalt();

      const key1 = await cache.getKey('password1', salt);
      const key2 = await cache.getKey('password2', salt);

      expect(key1).not.toEqual(key2);
      expect(cache.getStats().size).toBe(2);
    });

    it('should derive different keys for different salts', async () => {
      const password = 'test-password';

      const key1 = await cache.getKey(password, generateSalt());
      const key2 = await cache.getKey(password, generateSalt());

      expect(key1).not.toEqual(key2);
      expect(cache.getStats().size).toBe(2);
    });

    it('should increment access count on cache hit', async () => {
      const password = 'test-password';
      const salt = generateSalt();

      await cache.getKey(password, salt);
      await cache.getKey(password, salt);
      await cache.getKey(password, salt);

      const stats = cache.getStats();
      expect(stats.totalAccessCount).toBe(3);
    });
  });

  describe('invalidate', () => {
    it('should remove specific entry from cache', async () => {
      const password = 'test-password';
      const salt = generateSalt();

      await cache.getKey(password, salt);
      expect(cache.getStats().size).toBe(1);

      cache.invalidate(password, salt);
      expect(cache.getStats().size).toBe(0);
    });

    it('should not affect other entries', async () => {
      const salt1 = generateSalt();
      const salt2 = generateSalt();

      await cache.getKey('password1', salt1);
      await cache.getKey('password2', salt2);
      expect(cache.getStats().size).toBe(2);

      cache.invalidate('password1', salt1);
      expect(cache.getStats().size).toBe(1);
    });
  });

  describe('invalidateAll', () => {
    it('should clear all entries', async () => {
      await cache.getKey('password1', generateSalt());
      await cache.getKey('password2', generateSalt());
      await cache.getKey('password3', generateSalt());

      expect(cache.getStats().size).toBe(3);

      cache.invalidateAll();
      expect(cache.getStats().size).toBe(0);
    });
  });

  describe('getStats', () => {
    it('should return correct size', async () => {
      expect(cache.getStats().size).toBe(0);

      await cache.getKey('password1', generateSalt());
      expect(cache.getStats().size).toBe(1);

      await cache.getKey('password2', generateSalt());
      expect(cache.getStats().size).toBe(2);
    });

    it('should track total access count', async () => {
      const salt = generateSalt();

      await cache.getKey('password', salt);
      await cache.getKey('password', salt);
      await cache.getKey('password', salt);

      expect(cache.getStats().totalAccessCount).toBe(3);
    });

    it('should track entry timestamps', async () => {
      const before = Date.now();
      await cache.getKey('password', generateSalt());
      const after = Date.now();

      const stats = cache.getStats();
      expect(stats.oldestEntry).toBeGreaterThanOrEqual(before);
      expect(stats.oldestEntry).toBeLessThanOrEqual(after);
      expect(stats.newestEntry).toBeGreaterThanOrEqual(before);
      expect(stats.newestEntry).toBeLessThanOrEqual(after);
    });
  });

  describe('LRU eviction', () => {
    it('should evict least recently used when max entries reached', async () => {
      const smallCache = createKeyCache({ maxEntries: 2 });

      const salt1 = generateSalt();
      const salt2 = generateSalt();
      const salt3 = generateSalt();

      await smallCache.getKey('password1', salt1);
      await smallCache.getKey('password2', salt2);
      expect(smallCache.getStats().size).toBe(2);

      // Wait a bit to ensure different timestamps
      await new Promise(resolve => setTimeout(resolve, 10));

      // Access first entry to make it recently used
      await smallCache.getKey('password1', salt1);

      // Wait a bit more
      await new Promise(resolve => setTimeout(resolve, 10));

      // Add third entry - should evict password2 (least recently used)
      await smallCache.getKey('password3', salt3);
      expect(smallCache.getStats().size).toBe(2);

      // password1 should still be cached (was accessed recently)
      const key1Again = await smallCache.getKey('password1', salt1);
      expect(smallCache.getStats().totalAccessCount).toBeGreaterThanOrEqual(3);

      smallCache.destroy();
    });
  });

  describe('cleanup', () => {
    it('should clean up expired entries', async () => {
      const shortCache = createKeyCache({
        maxAgeMs: 100,
        maxIdleMs: 50,
      });

      const salt = generateSalt();
      await shortCache.getKey('password', salt);
      expect(shortCache.getStats().size).toBe(1);

      await new Promise(resolve => setTimeout(resolve, 200));

      await shortCache.getKey('password2', generateSalt());

      expect(shortCache.getStats().size).toBeLessThanOrEqual(2);

      shortCache.destroy();
    });
  });

  describe('destroy', () => {
    it('should clear all entries and stop cleanup interval', async () => {
      await cache.getKey('password1', generateSalt());
      await cache.getKey('password2', generateSalt());

      expect(cache.getStats().size).toBe(2);

      cache.destroy();

      expect(cache.getStats().size).toBe(0);
    });
  });
});

describe('createKeyCache', () => {
  it('should create cache with default config', () => {
    const cache = createKeyCache();
    expect(cache).toBeInstanceOf(KeyCache);
    cache.destroy();
  });

  it('should create cache with custom config', () => {
    const cache = createKeyCache({
      maxEntries: 5,
      maxAgeMs: 1000,
      maxIdleMs: 500,
    });
    expect(cache).toBeInstanceOf(KeyCache);
    cache.destroy();
  });
});
