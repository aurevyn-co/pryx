import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { VaultStorage, createVaultStorage } from '../../src/storage';
import { BackupManager } from '../../src/backup';
import * as fs from 'fs';
import * as path from 'path';
import * as os from 'os';

describe('VaultStorage Integration', () => {
  let tempDir: string;
  let vaultPath: string;
  let backupDir: string;
  let storage: VaultStorage;
  const password = 'integration-test-password';

  beforeEach(() => {
    tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'vault-integration-test-'));
    vaultPath = path.join(tempDir, 'vault.dat');
    backupDir = path.join(tempDir, 'backups');
    storage = createVaultStorage(backupDir);
  });

  afterEach(() => {
    fs.rmSync(tempDir, { recursive: true, force: true });
  });

  describe('full lifecycle', () => {
    it('should handle complete vault workflow', async () => {
      const vault = storage.createEmptyVault();

      const entry1 = await storage.addEntry(vault, {
        type: 'credential',
        name: 'GitHub API Key',
        data: { key: 'ghp_1234567890', username: 'testuser' },
      }, password);

      const entry2 = await storage.addEntry(vault, {
        type: 'api-key',
        name: 'OpenAI Key',
        data: { key: 'sk-abc123', organization: 'test-org' },
      }, password);

      expect(vault.entries).toHaveLength(2);

      await storage.save(vaultPath, vault, password);

      const loaded = await storage.load(vaultPath, password);
      expect(loaded.entries).toHaveLength(2);

      const retrieved1 = await storage.getEntry(loaded, entry1.id, password);
      expect(retrieved1.data).toEqual({ key: 'ghp_1234567890', username: 'testuser' });

      await storage.updateEntry(loaded, entry2.id, { name: 'OpenAI API Key' }, password);

      await storage.deleteEntry(loaded, entry1.id);
      expect(loaded.entries).toHaveLength(1);

      await storage.save(vaultPath, loaded, password);

      const final = await storage.load(vaultPath, password);
      expect(final.entries).toHaveLength(1);
      expect(final.entries[0].name).toBe('OpenAI API Key');
    });
  });

  describe('backup and restore', () => {
    it('should create and restore backup', async () => {
      const vault = storage.createEmptyVault();
      await storage.addEntry(vault, {
        type: 'credential',
        name: 'Test Entry',
        data: { secret: 'value' },
      }, password);

      await storage.save(vaultPath, vault, password);

      const backupPath = await storage.createBackup(vaultPath);
      expect(fs.existsSync(backupPath)).toBe(true);

      await storage.deleteEntry(vault, vault.entries[0].id);
      await storage.save(vaultPath, vault, password);

      const restored = await storage.restoreFromBackup(backupPath, vaultPath);
      expect(restored.entries).toHaveLength(1);

      const loaded = await storage.load(vaultPath, password);
      const retrieved = await storage.getEntry(loaded, restored.entries[0].id, password);
      expect(retrieved.data).toEqual({ secret: 'value' });
    });

    it('should maintain multiple backups', async () => {
      const vault = storage.createEmptyVault();
      await storage.save(vaultPath, vault, password);

      const backupPaths: string[] = [];
      for (let i = 0; i < 3; i++) {
        await new Promise(resolve => setTimeout(resolve, 10));
        const path = await storage.createBackup(vaultPath);
        backupPaths.push(path);
      }

      const backupManager = new BackupManager(backupDir);
      const backups = await backupManager.listBackups('vault.dat');
      expect(backups.length).toBeGreaterThanOrEqual(3);
    });
  });

  describe('integrity verification', () => {
    it('should detect corrupted entry', async () => {
      const vault = storage.createEmptyVault();
      await storage.addEntry(vault, {
        type: 'credential',
        name: 'Test Entry',
        data: { secret: 'value' },
      }, password);

      await storage.save(vaultPath, vault, password);

      const data = fs.readFileSync(vaultPath, 'utf-8');
      const corrupted = JSON.parse(data);
      corrupted.entries[0].encryptedData = 'corrupted-data';
      fs.writeFileSync(vaultPath, JSON.stringify(corrupted), { mode: 0o600 });

      const loaded = await storage.load(vaultPath, password);
      const report = await storage.verifyIntegrity(loaded, password);

      expect(report.valid).toBe(false);
      expect(report.corruptedEntries.length).toBeGreaterThan(0);
    });

    it('should verify integrity without password', async () => {
      const vault = storage.createEmptyVault();
      await storage.addEntry(vault, {
        type: 'credential',
        name: 'Test Entry',
        data: { secret: 'value' },
      }, password);

      const report = await storage.verifyIntegrity(vault);

      expect(report.valid).toBe(true);
      expect(report.entryCount).toBe(1);
    });
  });

  describe('concurrent operations', () => {
    it('should handle multiple entries of different types', async () => {
      const vault = storage.createEmptyVault();

      await storage.addEntry(vault, {
        type: 'credential',
        name: 'Database Password',
        data: { host: 'localhost', port: 5432, password: 'secret' },
      }, password);

      await storage.addEntry(vault, {
        type: 'api-key',
        name: 'Stripe API',
        data: { key: 'sk_test_123', mode: 'test' },
      }, password);

      await storage.addEntry(vault, {
        type: 'token',
        name: 'JWT Token',
        data: { token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9', expires: '2024-12-31' },
      }, password);

      await storage.addEntry(vault, {
        type: 'note',
        name: 'Secure Note',
        data: { content: 'This is a secret note', tags: ['personal'] },
      }, password);

      expect(vault.entries).toHaveLength(4);

      const entries = storage.listEntries(vault);
      const types = entries.map(e => e.type);
      expect(types).toContain('credential');
      expect(types).toContain('api-key');
      expect(types).toContain('token');
      expect(types).toContain('note');
    });
  });

  describe('large data handling', () => {
    it('should handle entries with large data', async () => {
      const vault = storage.createEmptyVault();
      const largeData = {
        content: 'x'.repeat(10000),
        metadata: { created: new Date().toISOString() },
      };

      const entry = await storage.addEntry(vault, {
        type: 'note',
        name: 'Large Note',
        data: largeData,
      }, password);

      await storage.save(vaultPath, vault, password);

      const loaded = await storage.load(vaultPath, password);
      const retrieved = await storage.getEntry(loaded, entry.id, password);

      expect(retrieved.data.content.length).toBe(10000);
    });
  });
});
