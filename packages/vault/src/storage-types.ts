

export const VAULT_FORMAT_VERSION = 1;
export const MAX_BACKUPS = 5;

export type EntryType = 'credential' | 'api-key' | 'token' | 'note';

export interface VaultMetadata {
  salt: string;
  algorithm: string;
  iterations: number;
  memoryCost: number;
}

export interface VaultEntry {
  id: string;
  type: EntryType;
  name: string;
  encryptedData: string;
  iv: string;
  tag: string;
  createdAt: string;
  updatedAt: string;
  accessCount: number;
  lastAccessedAt?: string;
}

export interface VaultFile {
  version: number;
  createdAt: string;
  updatedAt: string;
  metadata: VaultMetadata;
  entries: VaultEntry[];
}

export interface EntryData {
  id?: string;
  type: EntryType;
  name: string;
  data: Record<string, unknown>;
}

export interface EntryMetadata {
  id: string;
  type: EntryType;
  name: string;
  createdAt: string;
  updatedAt: string;
  accessCount: number;
  lastAccessedAt?: string;
}

export interface IntegrityReport {
  valid: boolean;
  errors: string[];
  entryCount: number;
  corruptedEntries: string[];
}

export interface Migration {
  fromVersion: number;
  toVersion: number;
  migrate: (data: unknown) => unknown;
}

export class StorageError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'StorageError';
  }
}

export class FileNotFoundError extends StorageError {
  constructor(filePath: string) {
    super(`Vault file not found: ${filePath}`);
    this.name = 'FileNotFoundError';
  }
}

export class CorruptedVaultError extends StorageError {
  constructor(message: string = 'Vault file appears to be corrupted') {
    super(message);
    this.name = 'CorruptedVaultError';
  }
}

export class EntryNotFoundError extends StorageError {
  constructor(entryId: string) {
    super(`Entry not found: ${entryId}`);
    this.name = 'EntryNotFoundError';
  }
}

export class DuplicateEntryError extends StorageError {
  constructor(entryId: string) {
    super(`Entry already exists: ${entryId}`);
    this.name = 'DuplicateEntryError';
  }
}

export class MigrationError extends StorageError {
  constructor(fromVersion: number, toVersion: number) {
    super(`Migration failed from version ${fromVersion} to ${toVersion}`);
    this.name = 'MigrationError';
  }
}
