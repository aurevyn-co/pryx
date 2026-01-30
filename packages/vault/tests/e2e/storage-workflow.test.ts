import { describe, it, expect, beforeEach, afterEach } from 'vitest';
import { VaultStorage, createVaultStorage } from '../../src/storage';
import * as fs from 'fs';
import * as path from 'path';
import * as os from 'os';

describe('VaultStorage E2E', () => {
  let tempDir: string;
  let vaultPath: string;
  let backupDir: string;
  let storage: VaultStorage;
  const masterPassword = 'MasterP@ssw0rd!2024';

  beforeEach(() => {
    tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'vault-e2e-test-'));
    vaultPath = path.join(tempDir, '.pryx', 'vault.dat');
    backupDir = path.join(tempDir, '.pryx', 'vault-backups');
    storage = createVaultStorage(backupDir);
  });

  afterEach(() => {
    fs.rmSync(tempDir, { recursive: true, force: true });
  });

  describe('real-world credential storage workflow', () => {
    it('should store and retrieve multiple API credentials', async () => {
      const vault = storage.createEmptyVault();

      const credentials = [
        {
          type: 'api-key' as const,
          name: 'OpenAI Production',
          data: {
            key: 'sk-prod-1234567890abcdef',
            organization: 'org-123',
            model: 'gpt-4',
          },
        },
        {
          type: 'api-key' as const,
          name: 'Anthropic Claude',
          data: {
            key: 'sk-ant-0987654321fedcba',
            version: 'claude-3-opus-20240229',
          },
        },
        {
          type: 'credential' as const,
          name: 'AWS Production',
          data: {
            accessKeyId: 'AKIAIOSFODNN7EXAMPLE',
            secretAccessKey: 'wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY',
            region: 'us-east-1',
          },
        },
        {
          type: 'token' as const,
          name: 'GitHub Personal Access Token',
          data: {
            token: 'ghp_xxxxxxxxxxxxxxxxxxxx',
            scopes: ['repo', 'workflow', 'read:packages'],
            expiresAt: '2024-12-31T23:59:59Z',
          },
        },
      ];

      const entryIds: string[] = [];
      for (const cred of credentials) {
        const entry = await storage.addEntry(vault, cred, masterPassword);
        entryIds.push(entry.id);
      }

      await storage.save(vaultPath, vault, masterPassword);

      const loaded = await storage.load(vaultPath, masterPassword);
      expect(loaded.entries).toHaveLength(4);

      const openai = await storage.getEntry(loaded, entryIds[0], masterPassword);
      expect(openai.name).toBe('OpenAI Production');
      expect(openai.data.key).toBe('sk-prod-1234567890abcdef');
      expect(openai.data.organization).toBe('org-123');

      const aws = await storage.getEntry(loaded, entryIds[2], masterPassword);
      expect(aws.data.region).toBe('us-east-1');
      expect(aws.data.accessKeyId).toBe('AKIAIOSFODNN7EXAMPLE');
    });

    it('should handle credential rotation workflow', async () => {
      const vault = storage.createEmptyVault();

      const entry = await storage.addEntry(vault, {
        type: 'api-key',
        name: 'Stripe API Key',
        data: { key: 'sk_old_123', mode: 'test' },
      }, masterPassword);

      await storage.save(vaultPath, vault, masterPassword);

      await storage.updateEntry(vault, entry.id, {
        data: { key: 'sk_new_456', mode: 'test', rotatedAt: new Date().toISOString() },
      }, masterPassword);

      await storage.save(vaultPath, vault, masterPassword);

      const loaded = await storage.load(vaultPath, masterPassword);
      const updated = await storage.getEntry(loaded, entry.id, masterPassword);

      expect(updated.data.key).toBe('sk_new_456');
      expect(updated.data.rotatedAt).toBeDefined();
    });
  });

  describe('backup and disaster recovery', () => {
    it('should recover from accidental deletion', async () => {
      const vault = storage.createEmptyVault();
      const entry = await storage.addEntry(vault, {
        type: 'credential',
        name: 'Critical Database',
        data: { host: 'prod.db.example.com', password: 'super-secret' },
      }, masterPassword);

      await storage.save(vaultPath, vault, masterPassword);
      const backupPath = await storage.createBackup(vaultPath);

      await storage.deleteEntry(vault, entry.id);
      await storage.save(vaultPath, vault, masterPassword);

      const corrupted = await storage.load(vaultPath, masterPassword);
      expect(corrupted.entries).toHaveLength(0);

      await storage.restoreFromBackup(backupPath, vaultPath);
      const restored = await storage.load(vaultPath, masterPassword);

      expect(restored.entries).toHaveLength(1);
      const recovered = await storage.getEntry(restored, entry.id, masterPassword);
      expect(recovered.data.password).toBe('super-secret');
    });

    it('should maintain backup rotation', async () => {
      const vault = storage.createEmptyVault();
      await storage.save(vaultPath, vault, masterPassword);

      for (let i = 0; i < 7; i++) {
        await new Promise(resolve => setTimeout(resolve, 20));
        await storage.createBackup(vaultPath);
      }

      const backupFiles = fs.readdirSync(backupDir);
      const vaultBackups = backupFiles.filter(f => f.startsWith('vault.dat') && f.endsWith('.backup'));
      expect(vaultBackups.length).toBeLessThanOrEqual(5);
    });
  });

  describe('integrity and corruption scenarios', () => {
    it('should detect and report corrupted vault file', async () => {
      const vault = storage.createEmptyVault();
      await storage.save(vaultPath, vault, masterPassword);

      const data = fs.readFileSync(vaultPath, 'utf-8');
      const corrupted = JSON.parse(data);
      corrupted.metadata = null;
      fs.writeFileSync(vaultPath, JSON.stringify(corrupted), { mode: 0o600 });

      await expect(storage.load(vaultPath, masterPassword)).rejects.toThrow();
    });

    it('should verify vault integrity on demand', async () => {
      const vault = storage.createEmptyVault();

      await storage.addEntry(vault, {
        type: 'credential',
        name: 'Test',
        data: { value: 'test' },
      }, masterPassword);

      const report = await storage.verifyIntegrity(vault, masterPassword);

      expect(report.valid).toBe(true);
      expect(report.entryCount).toBe(1);
      expect(report.corruptedEntries).toHaveLength(0);
    });
  });

  describe('file permissions and security', () => {
    it('should create vault with correct permissions', async () => {
      const vault = storage.createEmptyVault();
      await storage.save(vaultPath, vault, masterPassword);

      const stats = fs.statSync(vaultPath);
      const mode = stats.mode & 0o777;
      expect(mode).toBe(0o600);
    });

    it('should create backups with correct permissions', async () => {
      const vault = storage.createEmptyVault();
      await storage.save(vaultPath, vault, masterPassword);
      await storage.createBackup(vaultPath);

      const backupFiles = fs.readdirSync(backupDir);
      for (const file of backupFiles) {
        const stats = fs.statSync(path.join(backupDir, file));
        const mode = stats.mode & 0o777;
        expect(mode).toBe(0o600);
      }
    });
  });

  describe('edge cases', () => {
    it('should handle empty vault operations', async () => {
      const vault = storage.createEmptyVault();
      await storage.save(vaultPath, vault, masterPassword);

      const loaded = await storage.load(vaultPath, masterPassword);
      expect(loaded.entries).toHaveLength(0);

      const entries = storage.listEntries(loaded);
      expect(entries).toHaveLength(0);
    });

    it('should handle special characters in data', async () => {
      const vault = storage.createEmptyVault();
      const specialData = {
        content: 'Special chars: !@#$%^&*()_+-=[]{}|;:,.<>?',
        unicode: 'Unicode: 你好世界 🌍 émojis',
        multiline: 'Line 1\nLine 2\nLine 3',
      };

      const entry = await storage.addEntry(vault, {
        type: 'note',
        name: 'Special Characters Test',
        data: specialData,
      }, masterPassword);

      await storage.save(vaultPath, vault, masterPassword);

      const loaded = await storage.load(vaultPath, masterPassword);
      const retrieved = await storage.getEntry(loaded, entry.id, masterPassword);

      expect(retrieved.data.content).toBe(specialData.content);
      expect(retrieved.data.unicode).toBe(specialData.unicode);
      expect(retrieved.data.multiline).toBe(specialData.multiline);
    });

    it('should handle concurrent entry modifications', async () => {
      const vault = storage.createEmptyVault();

      const entries = await Promise.all([
        storage.addEntry(vault, { type: 'credential', name: 'Entry 1', data: { v: 1 } }, masterPassword),
        storage.addEntry(vault, { type: 'credential', name: 'Entry 2', data: { v: 2 } }, masterPassword),
        storage.addEntry(vault, { type: 'credential', name: 'Entry 3', data: { v: 3 } }, masterPassword),
      ]);

      expect(vault.entries).toHaveLength(3);
      expect(entries.map(e => e.name)).toContain('Entry 1');
      expect(entries.map(e => e.name)).toContain('Entry 2');
      expect(entries.map(e => e.name)).toContain('Entry 3');
    });
  });
});
