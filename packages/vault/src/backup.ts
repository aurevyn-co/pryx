import { copyFile, mkdir, readdir, stat, unlink } from 'fs/promises';
import { dirname, basename, join } from 'path';
import { MAX_BACKUPS } from './storage-types.js';

export interface BackupInfo {
  path: string;
  createdAt: Date;
  size: number;
}

export class BackupManager {
  private backupDir: string;

  constructor(backupDir: string) {
    this.backupDir = backupDir;
  }

  async createBackup(filePath: string): Promise<string> {
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    const originalName = basename(filePath);
    const backupPath = join(this.backupDir, `${originalName}.${timestamp}.backup`);
    
    await mkdir(this.backupDir, { recursive: true });
    await copyFile(filePath, backupPath);
    
    await this.cleanupOldBackups(originalName);
    
    return backupPath;
  }

  async listBackups(originalName: string): Promise<BackupInfo[]> {
    try {
      const files = await readdir(this.backupDir);
      const backups: BackupInfo[] = [];
      
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
    } catch {
      return [];
    }
  }

  async restoreBackup(backupPath: string, targetPath: string): Promise<void> {
    await mkdir(dirname(targetPath), { recursive: true });
    await copyFile(backupPath, targetPath);
  }

  async cleanupOldBackups(originalName: string): Promise<void> {
    const backups = await this.listBackups(originalName);
    
    if (backups.length > MAX_BACKUPS) {
      const toDelete = backups.slice(MAX_BACKUPS);
      for (const backup of toDelete) {
        try {
          await unlink(backup.path);
        } catch {
        }
      }
    }
  }

  async getLatestBackup(originalName: string): Promise<BackupInfo | null> {
    const backups = await this.listBackups(originalName);
    return backups[0] || null;
  }
}
