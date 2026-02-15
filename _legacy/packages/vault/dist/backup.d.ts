export interface BackupInfo {
    path: string;
    createdAt: Date;
    size: number;
}
export declare class BackupManager {
    private backupDir;
    constructor(backupDir: string);
    createBackup(filePath: string): Promise<string>;
    listBackups(originalName: string): Promise<BackupInfo[]>;
    restoreBackup(backupPath: string, targetPath: string): Promise<void>;
    cleanupOldBackups(originalName: string): Promise<void>;
    getLatestBackup(originalName: string): Promise<BackupInfo | null>;
}
//# sourceMappingURL=backup.d.ts.map