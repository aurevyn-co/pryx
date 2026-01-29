import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import {
  PasswordManager,
  createPasswordManager,
  DEFAULT_PASSWORD_MANAGER_CONFIG,
} from '../../src/password-manager.js';

describe('PasswordManager', () => {
  let manager: PasswordManager;

  beforeEach(() => {
    manager = createPasswordManager();
  });

  afterEach(() => {
    manager.destroy();
  });

  describe('unlock', () => {
    it('should unlock with correct password', async () => {
      expect(manager.isLocked()).toBe(true);

      await manager.unlock('correct-password');

      expect(manager.isLocked()).toBe(false);
      expect(manager.isUnlocked()).toBe(true);
    });

    it('should be idempotent if already unlocked', async () => {
      await manager.unlock('password');
      expect(manager.isUnlocked()).toBe(true);

      await manager.unlock('password');
      expect(manager.isUnlocked()).toBe(true);
    });
  });

  describe('lock', () => {
    it('should lock when unlocked', async () => {
      await manager.unlock('password');
      expect(manager.isUnlocked()).toBe(true);

      manager.lock();

      expect(manager.isLocked()).toBe(true);
      expect(manager.isUnlocked()).toBe(false);
    });

    it('should be safe to call when already locked', () => {
      expect(manager.isLocked()).toBe(true);
      manager.lock();
      expect(manager.isLocked()).toBe(true);
    });
  });

  describe('encrypt', () => {
    it('should encrypt when unlocked', async () => {
      await manager.unlock('password');

      const plaintext = Buffer.from('secret data');
      const encrypted = await manager.encrypt(plaintext);

      expect(encrypted.ciphertext).toBeDefined();
      expect(encrypted.iv).toBeDefined();
      expect(encrypted.salt).toBeDefined();
      expect(encrypted.tag).toBeDefined();
    });

    it('should throw when locked', async () => {
      const plaintext = Buffer.from('secret data');

      await expect(manager.encrypt(plaintext)).rejects.toThrow('Vault is locked');
    });
  });

  describe('decrypt', () => {
    it('should decrypt encrypted data', async () => {
      await manager.unlock('password');

      const plaintext = Buffer.from('secret data');
      const encrypted = await manager.encrypt(plaintext);

      const decrypted = await manager.decrypt(encrypted);

      expect(decrypted.toString()).toBe('secret data');
    });

    it('should throw when locked', async () => {
      await manager.unlock('password');
      const plaintext = Buffer.from('secret data');
      const encrypted = await manager.encrypt(plaintext);
      manager.lock();

      await expect(manager.decrypt(encrypted)).rejects.toThrow('Vault is locked');
    });
  });

  describe('changePassword', () => {
    it('should change password and remain functional', async () => {
      const oldPassword = 'old-password';
      const newPassword = 'new-password';

      await manager.unlock(oldPassword);

      const plaintext = Buffer.from('secret data');
      const encrypted = await manager.encrypt(plaintext);

      await manager.changePassword(oldPassword, newPassword);

      expect(manager.isUnlocked()).toBe(true);

      const newPlaintext = Buffer.from('new secret data');
      const newEncrypted = await manager.encrypt(newPlaintext);
      const newDecrypted = await manager.decrypt(newEncrypted);
      expect(newDecrypted.toString()).toBe('new secret data');
    });

    it('should throw when locked', async () => {
      await expect(
        manager.changePassword('old', 'new')
      ).rejects.toThrow('Vault is locked');
    });
  });

  describe('auto-lock', () => {
    it('should auto-lock after timeout', async () => {
      const shortManager = createPasswordManager({
        autoLockMs: 100, // 100ms for testing
      });

      await shortManager.unlock('password');
      expect(shortManager.isUnlocked()).toBe(true);

      // Wait for auto-lock
      await new Promise(resolve => setTimeout(resolve, 150));

      expect(shortManager.isLocked()).toBe(true);

      shortManager.destroy();
    });

    it('should reset timer on activity', async () => {
      const shortManager = createPasswordManager({
        autoLockMs: 200, // 200ms for testing
      });

      await shortManager.unlock('password');

      // Activity at 100ms
      await new Promise(resolve => setTimeout(resolve, 100));
      await shortManager.encrypt(Buffer.from('data'));

      // Should still be unlocked at 250ms (100ms + 200ms timeout)
      await new Promise(resolve => setTimeout(resolve, 100));
      expect(shortManager.isUnlocked()).toBe(true);

      // Wait for auto-lock after activity
      await new Promise(resolve => setTimeout(resolve, 250));
      expect(shortManager.isLocked()).toBe(true);

      shortManager.destroy();
    });

    it('should return null for remaining lock time when locked', () => {
      expect(manager.getRemainingLockTime()).toBeNull();
    });

    it('should return timeout when unlocked', async () => {
      await manager.unlock('password');
      expect(manager.getRemainingLockTime()).toBe(DEFAULT_PASSWORD_MANAGER_CONFIG.autoLockMs);
    });
  });

  describe('destroy', () => {
    it('should lock and cleanup when destroyed', async () => {
      await manager.unlock('password');
      expect(manager.isUnlocked()).toBe(true);

      manager.destroy();

      expect(manager.isLocked()).toBe(true);
    });
  });
});

describe('createPasswordManager', () => {
  it('should create with default config', () => {
    const manager = createPasswordManager();
    expect(manager).toBeInstanceOf(PasswordManager);
    expect(manager.isLocked()).toBe(true);
    manager.destroy();
  });

  it('should create with custom config', () => {
    const manager = createPasswordManager({
      autoLockMs: 10000,
    });
    expect(manager).toBeInstanceOf(PasswordManager);
    manager.destroy();
  });
});
