import { copyFile, mkdir, readdir, stat, unlink } from 'fs/promises';
import { dirname, basename, join } from 'path';
import { MAX_BACKUPS } from './storage-types.js';
export class BackupManager {
    backupDir;
    constructor(backupDir) {
        this.backupDir = backupDir;
    }
    async createBackup(filePath) {
        const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
        const originalName = basename(filePath);
        const backupPath = join(this.backupDir, `${originalName}.${timestamp}.backup`);
        await mkdir(this.backupDir, { recursive: true });
        await copyFile(filePath, backupPath);
        await this.cleanupOldBackups(originalName);
        return backupPath;
    }
    async listBackups(originalName) {
        try {
            const files = await readdir(this.backupDir);
            const backups = [];
            for (const file of files) {
                if (file.startsWith(originalName) && file.endsWith('.backup')) {
                    const path = join(this.backupDir, file);
                    const stats = await stat(path);
                    backups.push({
                        path,
                        createdAt: stats.mtime,
                        size: stats.size,
                    });
                }
            }
            return backups.sort((a, b) => b.createdAt.getTime() - a.createdAt.getTime());
        }
        catch {
            return [];
        }
    }
    async restoreBackup(backupPath, targetPath) {
        await mkdir(dirname(targetPath), { recursive: true });
        await copyFile(backupPath, targetPath);
    }
    async cleanupOldBackups(originalName) {
        const backups = await this.listBackups(originalName);
        if (backups.length > MAX_BACKUPS) {
            const toDelete = backups.slice(MAX_BACKUPS);
            for (const backup of toDelete) {
                try {
                    await unlink(backup.path);
                }
                catch {
                }
            }
        }
    }
    async getLatestBackup(originalName) {
        const backups = await this.listBackups(originalName);
        return backups[0] || null;
    }
}
//# sourceMappingURL=backup.js.map