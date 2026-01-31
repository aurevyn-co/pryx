import { readFile, writeFile, mkdir, rename, unlink, access, constants } from 'fs/promises';
import { dirname, join } from 'path';
import { randomUUID } from 'crypto';
import { deriveKey, encrypt, decrypt, generateSalt, generateIV } from './crypto.js';
import { DEFAULT_VAULT_CONFIG } from './types.js';
import { BackupManager } from './backup.js';
import {
  VaultFile,
  VaultEntry,
  EntryData,
  EntryMetadata,
  IntegrityReport,
  VAULT_FORMAT_VERSION,
  StorageError,
  FileNotFoundError,
  CorruptedVaultError,
  EntryNotFoundError,
  DuplicateEntryError,
} from './storage-types.js';

export class VaultStorage {
  private backupManager: BackupManager;

  constructor(backupDir?: string) {
    const defaultBackupDir = join(process.env.HOME || process.env.USERPROFILE || '.', '.pryx', 'vault-backups');
    this.backupManager = new BackupManager(backupDir || defaultBackupDir);
  }

  async load(filePath: string, _password: string): Promise<VaultFile> {
    try {
      await access(filePath, constants.R_OK);
    } catch {
      throw new FileNotFoundError(filePath);
    }

    const data = await readFile(filePath, 'utf-8');
    let vault: VaultFile;

    try {
      vault = JSON.parse(data);
    } catch {
      throw new CorruptedVaultError('Invalid JSON format');
    }

    this.validateVaultStructure(vault);
    
    await this.verifyIntegrity(vault, _password);

    return vault;
  }

  async save(filePath: string, vault: VaultFile, _password: string): Promise<void> {
    vault.updatedAt = new Date().toISOString();
    
    await mkdir(dirname(filePath), { recursive: true });
    
    const tempPath = `${filePath}.tmp.${Date.now()}`;
    
    try {
      await writeFile(tempPath, JSON.stringify(vault, null, 2), { mode: 0o600 });
      await rename(tempPath, filePath);
    } catch (error) {
      try {
        await unlink(tempPath);
      } catch {}
      throw new StorageError(`Failed to save vault: ${(error as Error).message}`);
    }
  }

  async addEntry(vault: VaultFile, entryData: EntryData, password: string): Promise<VaultEntry> {
    const existingIndex = vault.entries.findIndex(e => e.id === entryData.id);
    if (existingIndex !== -1) {
      throw new DuplicateEntryError(entryData.id || 'unknown');
    }

    const id = entryData.id || randomUUID();
    const now = new Date().toISOString();
    
    const salt = Buffer.from(vault.metadata.salt, 'base64');
    const key = await deriveKey(password, salt);
    const iv = generateIV();
    
    const plaintext = Buffer.from(JSON.stringify(entryData.data));
    const { ciphertext, tag } = encrypt(plaintext, key, iv);
    
    const entry: VaultEntry = {
      id,
      type: entryData.type,
      name: entryData.name,
      encryptedData: ciphertext.toString('base64'),
      iv: iv.toString('base64'),
      tag: tag.toString('base64'),
      createdAt: now,
      updatedAt: now,
      accessCount: 0,
    };

    vault.entries.push(entry);
    
    return entry;
  }

  async updateEntry(
    vault: VaultFile,
    id: string,
    updates: Partial<EntryData>,
    password: string
  ): Promise<VaultEntry> {
    const index = vault.entries.findIndex(e => e.id === id);
    if (index === -1) {
      throw new EntryNotFoundError(id);
    }

    const entry = vault.entries[index];
    
    if (updates.name !== undefined) {
      entry.name = updates.name;
    }
    
    if (updates.data !== undefined) {
      const salt = Buffer.from(vault.metadata.salt, 'base64');
      const key = await deriveKey(password, salt);
      const iv = generateIV();
      
      const plaintext = Buffer.from(JSON.stringify(updates.data));
      const { ciphertext, tag } = encrypt(plaintext, key, iv);
      
      entry.encryptedData = ciphertext.toString('base64');
      entry.iv = iv.toString('base64');
      entry.tag = tag.toString('base64');
    }

    entry.updatedAt = new Date().toISOString();
    
    return entry;
  }

  async deleteEntry(vault: VaultFile, id: string): Promise<void> {
    const index = vault.entries.findIndex(e => e.id === id);
    if (index === -1) {
      throw new EntryNotFoundError(id);
    }
    
    vault.entries.splice(index, 1);
  }

  async getEntry(vault: VaultFile, id: string, password: string): Promise<EntryData> {
    const entry = vault.entries.find(e => e.id === id);
    if (!entry) {
      throw new EntryNotFoundError(id);
    }

    try {
      const salt = Buffer.from(vault.metadata.salt, 'base64');
      const key = await deriveKey(password, salt);
      const iv = Buffer.from(entry.iv, 'base64');
      const ciphertext = Buffer.from(entry.encryptedData, 'base64');
      const tag = Buffer.from(entry.tag, 'base64');
      
      const plaintext = decrypt(ciphertext, key, iv, tag);
      const data = JSON.parse(plaintext.toString('utf-8'));
      
      entry.accessCount++;
      entry.lastAccessedAt = new Date().toISOString();
      
      return {
        id: entry.id,
        type: entry.type,
        name: entry.name,
        data,
      };
    } catch (error) {
      throw new CorruptedVaultError(`Failed to decrypt entry ${id}: ${(error as Error).message}`);
    }
  }

  listEntries(vault: VaultFile): EntryMetadata[] {
    return vault.entries.map(entry => ({
      id: entry.id,
      type: entry.type,
      name: entry.name,
      createdAt: entry.createdAt,
      updatedAt: entry.updatedAt,
      accessCount: entry.accessCount,
      lastAccessedAt: entry.lastAccessedAt,
    }));
  }

  async createBackup(filePath: string): Promise<string> {
    return this.backupManager.createBackup(filePath);
  }

  async restoreFromBackup(backupPath: string, targetPath: string): Promise<VaultFile> {
    await this.backupManager.restoreBackup(backupPath, targetPath);
    
    const data = await readFile(targetPath, 'utf-8');
    return JSON.parse(data);
  }

  async verifyIntegrity(vault: VaultFile, password?: string): Promise<IntegrityReport> {
    const report: IntegrityReport = {
      valid: true,
      errors: [],
      entryCount: vault.entries.length,
      corruptedEntries: [],
    };

    if (vault.version !== VAULT_FORMAT_VERSION) {
      report.valid = false;
      report.errors.push(`Unsupported vault version: ${vault.version}. Expected: ${VAULT_FORMAT_VERSION}`);
    }

    if (!vault.metadata || !vault.metadata.salt) {
      report.valid = false;
      report.errors.push('Missing or invalid vault metadata');
    }

    if (password) {
      for (const entry of vault.entries) {
        try {
          const salt = Buffer.from(vault.metadata.salt, 'base64');
          const key = await deriveKey(password, salt);
          const iv = Buffer.from(entry.iv, 'base64');
          const ciphertext = Buffer.from(entry.encryptedData, 'base64');
          const tag = Buffer.from(entry.tag, 'base64');
          
          decrypt(ciphertext, key, iv, tag);
        } catch {
          report.valid = false;
          report.corruptedEntries.push(entry.id);
        }
      }
    }

    return report;
  }

  createEmptyVault(): VaultFile {
    const salt = generateSalt();
    const now = new Date().toISOString();
    
    return {
      version: VAULT_FORMAT_VERSION,
      createdAt: now,
      updatedAt: now,
      metadata: {
        salt: salt.toString('base64'),
        algorithm: 'argon2id+aes-256-gcm',
        iterations: DEFAULT_VAULT_CONFIG.timeCost,
        memoryCost: DEFAULT_VAULT_CONFIG.memoryCost,
      },
      entries: [],
    };
  }

  private validateVaultStructure(vault: unknown): asserts vault is VaultFile {
    if (typeof vault !== 'object' || vault === null) {
      throw new CorruptedVaultError('Vault is not an object');
    }

    const v = vault as Record<string, unknown>;

    if (typeof v.version !== 'number') {
      throw new CorruptedVaultError('Missing or invalid version');
    }

    if (typeof v.createdAt !== 'string') {
      throw new CorruptedVaultError('Missing or invalid createdAt');
    }

    if (typeof v.updatedAt !== 'string') {
      throw new CorruptedVaultError('Missing or invalid updatedAt');
    }

    if (typeof v.metadata !== 'object' || v.metadata === null) {
      throw new CorruptedVaultError('Missing or invalid metadata');
    }

    if (!Array.isArray(v.entries)) {
      throw new CorruptedVaultError('Missing or invalid entries array');
    }
  }
}

export function createVaultStorage(backupDir?: string): VaultStorage {
  return new VaultStorage(backupDir);
}
