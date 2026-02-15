"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g = Object.create((typeof Iterator === "function" ? Iterator : Object).prototype);
    return g.next = verb(0), g["throw"] = verb(1), g["return"] = verb(2), typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (g && (g = 0, op[0] && (_ = 0)), _) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
Object.defineProperty(exports, "__esModule", { value: true });
var vitest_1 = require("vitest");
var storage_js_1 = require("../src/storage.js");
var backup_js_1 = require("../src/backup.js");
var fs_1 = require("fs");
var path_1 = require("path");
var os_1 = require("os");
(0, vitest_1.describe)('VaultStorage Integration', function () {
    var tempDir;
    var vaultPath;
    var backupDir;
    var storage;
    var password = 'integration-test-password';
    (0, vitest_1.beforeEach)(function () {
        tempDir = fs_1.default.mkdtempSync(path_1.default.join(os_1.default.tmpdir(), 'vault-integration-test-'));
        vaultPath = path_1.default.join(tempDir, 'vault.dat');
        backupDir = path_1.default.join(tempDir, 'backups');
        storage = (0, storage_js_1.createVaultStorage)(backupDir);
    });
    (0, vitest_1.afterEach)(function () {
        fs_1.default.rmSync(tempDir, { recursive: true, force: true });
    });
    (0, vitest_1.describe)('full lifecycle', function () {
        (0, vitest_1.it)('should handle complete vault workflow', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, entry1, entry2, loaded, retrieved1, final;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.addEntry(vault, {
                                type: 'credential',
                                name: 'GitHub API Key',
                                data: { key: 'ghp_1234567890', username: 'testuser' },
                            }, password)];
                    case 1:
                        entry1 = _a.sent();
                        return [4 /*yield*/, storage.addEntry(vault, {
                                type: 'api-key',
                                name: 'OpenAI Key',
                                data: { key: 'sk-abc123', organization: 'test-org' },
                            }, password)];
                    case 2:
                        entry2 = _a.sent();
                        (0, vitest_1.expect)(vault.entries).toHaveLength(2);
                        return [4 /*yield*/, storage.save(vaultPath, vault, password)];
                    case 3:
                        _a.sent();
                        return [4 /*yield*/, storage.load(vaultPath, password)];
                    case 4:
                        loaded = _a.sent();
                        (0, vitest_1.expect)(loaded.entries).toHaveLength(2);
                        return [4 /*yield*/, storage.getEntry(loaded, entry1.id, password)];
                    case 5:
                        retrieved1 = _a.sent();
                        (0, vitest_1.expect)(retrieved1.data).toEqual({ key: 'ghp_1234567890', username: 'testuser' });
                        return [4 /*yield*/, storage.updateEntry(loaded, entry2.id, { name: 'OpenAI API Key' }, password)];
                    case 6:
                        _a.sent();
                        return [4 /*yield*/, storage.deleteEntry(loaded, entry1.id)];
                    case 7:
                        _a.sent();
                        (0, vitest_1.expect)(loaded.entries).toHaveLength(1);
                        return [4 /*yield*/, storage.save(vaultPath, loaded, password)];
                    case 8:
                        _a.sent();
                        return [4 /*yield*/, storage.load(vaultPath, password)];
                    case 9:
                        final = _a.sent();
                        (0, vitest_1.expect)(final.entries).toHaveLength(1);
                        (0, vitest_1.expect)(final.entries[0].name).toBe('OpenAI API Key');
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('backup and restore', function () {
        (0, vitest_1.it)('should create and restore backup', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, backupPath, restored, loaded, retrieved;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.addEntry(vault, {
                                type: 'credential',
                                name: 'Test Entry',
                                data: { secret: 'value' },
                            }, password)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, storage.save(vaultPath, vault, password)];
                    case 2:
                        _a.sent();
                        return [4 /*yield*/, storage.createBackup(vaultPath)];
                    case 3:
                        backupPath = _a.sent();
                        (0, vitest_1.expect)(fs_1.default.existsSync(backupPath)).toBe(true);
                        return [4 /*yield*/, storage.deleteEntry(vault, vault.entries[0].id)];
                    case 4:
                        _a.sent();
                        return [4 /*yield*/, storage.save(vaultPath, vault, password)];
                    case 5:
                        _a.sent();
                        return [4 /*yield*/, storage.restoreFromBackup(backupPath, vaultPath)];
                    case 6:
                        restored = _a.sent();
                        (0, vitest_1.expect)(restored.entries).toHaveLength(1);
                        return [4 /*yield*/, storage.load(vaultPath, password)];
                    case 7:
                        loaded = _a.sent();
                        return [4 /*yield*/, storage.getEntry(loaded, restored.entries[0].id, password)];
                    case 8:
                        retrieved = _a.sent();
                        (0, vitest_1.expect)(retrieved.data).toEqual({ secret: 'value' });
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should maintain multiple backups', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, backupPaths, i, path_2, backupManager, backups;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.save(vaultPath, vault, password)];
                    case 1:
                        _a.sent();
                        backupPaths = [];
                        i = 0;
                        _a.label = 2;
                    case 2:
                        if (!(i < 3)) return [3 /*break*/, 6];
                        return [4 /*yield*/, new Promise(function (resolve) { return setTimeout(resolve, 10); })];
                    case 3:
                        _a.sent();
                        return [4 /*yield*/, storage.createBackup(vaultPath)];
                    case 4:
                        path_2 = _a.sent();
                        backupPaths.push(path_2);
                        _a.label = 5;
                    case 5:
                        i++;
                        return [3 /*break*/, 2];
                    case 6:
                        backupManager = new backup_js_1.BackupManager(backupDir);
                        return [4 /*yield*/, backupManager.listBackups('vault.dat')];
                    case 7:
                        backups = _a.sent();
                        (0, vitest_1.expect)(backups.length).toBeGreaterThanOrEqual(3);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('integrity verification', function () {
        (0, vitest_1.it)('should detect corrupted entry', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, data, corrupted, loaded, report;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.addEntry(vault, {
                                type: 'credential',
                                name: 'Test Entry',
                                data: { secret: 'value' },
                            }, password)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, storage.save(vaultPath, vault, password)];
                    case 2:
                        _a.sent();
                        data = fs_1.default.readFileSync(vaultPath, 'utf-8');
                        corrupted = JSON.parse(data);
                        corrupted.entries[0].encryptedData = 'corrupted-data';
                        fs_1.default.writeFileSync(vaultPath, JSON.stringify(corrupted), { mode: 384 });
                        return [4 /*yield*/, storage.load(vaultPath, password)];
                    case 3:
                        loaded = _a.sent();
                        return [4 /*yield*/, storage.verifyIntegrity(loaded, password)];
                    case 4:
                        report = _a.sent();
                        (0, vitest_1.expect)(report.valid).toBe(false);
                        (0, vitest_1.expect)(report.corruptedEntries.length).toBeGreaterThan(0);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should verify integrity without password', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, report;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.addEntry(vault, {
                                type: 'credential',
                                name: 'Test Entry',
                                data: { secret: 'value' },
                            }, password)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, storage.verifyIntegrity(vault)];
                    case 2:
                        report = _a.sent();
                        (0, vitest_1.expect)(report.valid).toBe(true);
                        (0, vitest_1.expect)(report.entryCount).toBe(1);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('concurrent operations', function () {
        (0, vitest_1.it)('should handle multiple entries of different types', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, entries, types;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.addEntry(vault, {
                                type: 'credential',
                                name: 'Database Password',
                                data: { host: 'localhost', port: 5432, password: 'secret' },
                            }, password)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, storage.addEntry(vault, {
                                type: 'api-key',
                                name: 'Stripe API',
                                data: { key: 'sk_test_123', mode: 'test' },
                            }, password)];
                    case 2:
                        _a.sent();
                        return [4 /*yield*/, storage.addEntry(vault, {
                                type: 'token',
                                name: 'JWT Token',
                                data: { token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9', expires: '2024-12-31' },
                            }, password)];
                    case 3:
                        _a.sent();
                        return [4 /*yield*/, storage.addEntry(vault, {
                                type: 'note',
                                name: 'Secure Note',
                                data: { content: 'This is a secret note', tags: ['personal'] },
                            }, password)];
                    case 4:
                        _a.sent();
                        (0, vitest_1.expect)(vault.entries).toHaveLength(4);
                        entries = storage.listEntries(vault);
                        types = entries.map(function (e) { return e.type; });
                        (0, vitest_1.expect)(types).toContain('credential');
                        (0, vitest_1.expect)(types).toContain('api-key');
                        (0, vitest_1.expect)(types).toContain('token');
                        (0, vitest_1.expect)(types).toContain('note');
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('large data handling', function () {
        (0, vitest_1.it)('should handle entries with large data', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, largeData, entry, loaded, retrieved;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        largeData = {
                            content: 'x'.repeat(10000),
                            metadata: { created: new Date().toISOString() },
                        };
                        return [4 /*yield*/, storage.addEntry(vault, {
                                type: 'note',
                                name: 'Large Note',
                                data: largeData,
                            }, password)];
                    case 1:
                        entry = _a.sent();
                        return [4 /*yield*/, storage.save(vaultPath, vault, password)];
                    case 2:
                        _a.sent();
                        return [4 /*yield*/, storage.load(vaultPath, password)];
                    case 3:
                        loaded = _a.sent();
                        return [4 /*yield*/, storage.getEntry(loaded, entry.id, password)];
                    case 4:
                        retrieved = _a.sent();
                        (0, vitest_1.expect)(retrieved.data.content.length).toBe(10000);
                        return [2 /*return*/];
                }
            });
        }); });
    });
});
