import { VaultFile, VaultEntry, EntryData, EntryMetadata, IntegrityReport } from './storage-types.js';
export declare class VaultStorage {
    private backupManager;
    constructor(backupDir?: string);
    load(filePath: string, _password: string): Promise<VaultFile>;
    save(filePath: string, vault: VaultFile, _password: string): Promise<void>;
    addEntry(vault: VaultFile, entryData: EntryData, password: string): Promise<VaultEntry>;
    updateEntry(vault: VaultFile, id: string, updates: Partial<EntryData>, password: string): Promise<VaultEntry>;
    deleteEntry(vault: VaultFile, id: string): Promise<void>;
    getEntry(vault: VaultFile, id: string, password: string): Promise<EntryData>;
    listEntries(vault: VaultFile): EntryMetadata[];
    createBackup(filePath: string): Promise<string>;
    restoreFromBackup(backupPath: string, targetPath: string): Promise<VaultFile>;
    verifyIntegrity(vault: VaultFile, password?: string): Promise<IntegrityReport>;
    createEmptyVault(): VaultFile;
    private validateVaultStructure;
}
export declare function createVaultStorage(backupDir?: string): VaultStorage;
//# sourceMappingURL=storage.d.ts.map