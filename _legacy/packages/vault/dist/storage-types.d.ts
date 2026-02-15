export declare const VAULT_FORMAT_VERSION = 1;
export declare const MAX_BACKUPS = 5;
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
export declare class StorageError extends Error {
    constructor(message: string);
}
export declare class FileNotFoundError extends StorageError {
    constructor(filePath: string);
}
export declare class CorruptedVaultError extends StorageError {
    constructor(message?: string);
}
export declare class EntryNotFoundError extends StorageError {
    constructor(entryId: string);
}
export declare class DuplicateEntryError extends StorageError {
    constructor(entryId: string);
}
export declare class MigrationError extends StorageError {
    constructor(fromVersion: number, toVersion: number);
}
//# sourceMappingURL=storage-types.d.ts.map