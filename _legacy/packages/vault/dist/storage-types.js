export const VAULT_FORMAT_VERSION = 1;
export const MAX_BACKUPS = 5;
export class StorageError extends Error {
    constructor(message) {
        super(message);
        this.name = 'StorageError';
    }
}
export class FileNotFoundError extends StorageError {
    constructor(filePath) {
        super(`Vault file not found: ${filePath}`);
        this.name = 'FileNotFoundError';
    }
}
export class CorruptedVaultError extends StorageError {
    constructor(message = 'Vault file appears to be corrupted') {
        super(message);
        this.name = 'CorruptedVaultError';
    }
}
export class EntryNotFoundError extends StorageError {
    constructor(entryId) {
        super(`Entry not found: ${entryId}`);
        this.name = 'EntryNotFoundError';
    }
}
export class DuplicateEntryError extends StorageError {
    constructor(entryId) {
        super(`Entry already exists: ${entryId}`);
        this.name = 'DuplicateEntryError';
    }
}
export class MigrationError extends StorageError {
    constructor(fromVersion, toVersion) {
        super(`Migration failed from version ${fromVersion} to ${toVersion}`);
        this.name = 'MigrationError';
    }
}
//# sourceMappingURL=storage-types.js.map