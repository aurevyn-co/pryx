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
var fs_1 = require("fs");
var path_1 = require("path");
var os_1 = require("os");
(0, vitest_1.describe)('VaultStorage E2E', function () {
    var tempDir;
    var vaultPath;
    var backupDir;
    var storage;
    var masterPassword = 'MasterP@ssw0rd!2024';
    (0, vitest_1.beforeEach)(function () {
        tempDir = fs_1.default.mkdtempSync(path_1.default.join(os_1.default.tmpdir(), 'vault-e2e-test-'));
        vaultPath = path_1.default.join(tempDir, '.pryx', 'vault.dat');
        backupDir = path_1.default.join(tempDir, '.pryx', 'vault-backups');
        storage = (0, storage_js_1.createVaultStorage)(backupDir);
    });
    (0, vitest_1.afterEach)(function () {
        fs_1.default.rmSync(tempDir, { recursive: true, force: true });
    });
    (0, vitest_1.describe)('real-world credential storage workflow', function () {
        (0, vitest_1.it)('should store and retrieve multiple API credentials', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, credentials, entryIds, _i, credentials_1, cred, entry, loaded, openai, aws;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        credentials = [
                            {
                                type: 'api-key',
                                name: 'OpenAI Production',
                                data: {
                                    key: 'sk-prod-1234567890abcdef',
                                    organization: 'org-123',
                                    model: 'gpt-4',
                                },
                            },
                            {
                                type: 'api-key',
                                name: 'Anthropic Claude',
                                data: {
                                    key: 'sk-ant-0987654321fedcba',
                                    version: 'claude-3-opus-20240229',
                                },
                            },
                            {
                                type: 'credential',
                                name: 'AWS Production',
                                data: {
                                    accessKeyId: 'AKIAIOSFODNN7EXAMPLE',
                                    secretAccessKey: 'wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY',
                                    region: 'us-east-1',
                                },
                            },
                            {
                                type: 'token',
                                name: 'GitHub Personal Access Token',
                                data: {
                                    token: 'ghp_xxxxxxxxxxxxxxxxxxxx',
                                    scopes: ['repo', 'workflow', 'read:packages'],
                                    expiresAt: '2024-12-31T23:59:59Z',
                                },
                            },
                        ];
                        entryIds = [];
                        _i = 0, credentials_1 = credentials;
                        _a.label = 1;
                    case 1:
                        if (!(_i < credentials_1.length)) return [3 /*break*/, 4];
                        cred = credentials_1[_i];
                        return [4 /*yield*/, storage.addEntry(vault, cred, masterPassword)];
                    case 2:
                        entry = _a.sent();
                        entryIds.push(entry.id);
                        _a.label = 3;
                    case 3:
                        _i++;
                        return [3 /*break*/, 1];
                    case 4: return [4 /*yield*/, storage.save(vaultPath, vault, masterPassword)];
                    case 5:
                        _a.sent();
                        return [4 /*yield*/, storage.load(vaultPath, masterPassword)];
                    case 6:
                        loaded = _a.sent();
                        (0, vitest_1.expect)(loaded.entries).toHaveLength(4);
                        return [4 /*yield*/, storage.getEntry(loaded, entryIds[0], masterPassword)];
                    case 7:
                        openai = _a.sent();
                        (0, vitest_1.expect)(openai.name).toBe('OpenAI Production');
                        (0, vitest_1.expect)(openai.data.key).toBe('sk-prod-1234567890abcdef');
                        (0, vitest_1.expect)(openai.data.organization).toBe('org-123');
                        return [4 /*yield*/, storage.getEntry(loaded, entryIds[2], masterPassword)];
                    case 8:
                        aws = _a.sent();
                        (0, vitest_1.expect)(aws.data.region).toBe('us-east-1');
                        (0, vitest_1.expect)(aws.data.accessKeyId).toBe('AKIAIOSFODNN7EXAMPLE');
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should handle credential rotation workflow', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, entry, loaded, updated;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.addEntry(vault, {
                                type: 'api-key',
                                name: 'Stripe API Key',
                                data: { key: 'sk_old_123', mode: 'test' },
                            }, masterPassword)];
                    case 1:
                        entry = _a.sent();
                        return [4 /*yield*/, storage.save(vaultPath, vault, masterPassword)];
                    case 2:
                        _a.sent();
                        return [4 /*yield*/, storage.updateEntry(vault, entry.id, {
                                data: { key: 'sk_new_456', mode: 'test', rotatedAt: new Date().toISOString() },
                            }, masterPassword)];
                    case 3:
                        _a.sent();
                        return [4 /*yield*/, storage.save(vaultPath, vault, masterPassword)];
                    case 4:
                        _a.sent();
                        return [4 /*yield*/, storage.load(vaultPath, masterPassword)];
                    case 5:
                        loaded = _a.sent();
                        return [4 /*yield*/, storage.getEntry(loaded, entry.id, masterPassword)];
                    case 6:
                        updated = _a.sent();
                        (0, vitest_1.expect)(updated.data.key).toBe('sk_new_456');
                        (0, vitest_1.expect)(updated.data.rotatedAt).toBeDefined();
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('backup and disaster recovery', function () {
        (0, vitest_1.it)('should recover from accidental deletion', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, entry, backupPath, corrupted, restored, recovered;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.addEntry(vault, {
                                type: 'credential',
                                name: 'Critical Database',
                                data: { host: 'prod.db.example.com', password: 'super-secret' },
                            }, masterPassword)];
                    case 1:
                        entry = _a.sent();
                        return [4 /*yield*/, storage.save(vaultPath, vault, masterPassword)];
                    case 2:
                        _a.sent();
                        return [4 /*yield*/, storage.createBackup(vaultPath)];
                    case 3:
                        backupPath = _a.sent();
                        return [4 /*yield*/, storage.deleteEntry(vault, entry.id)];
                    case 4:
                        _a.sent();
                        return [4 /*yield*/, storage.save(vaultPath, vault, masterPassword)];
                    case 5:
                        _a.sent();
                        return [4 /*yield*/, storage.load(vaultPath, masterPassword)];
                    case 6:
                        corrupted = _a.sent();
                        (0, vitest_1.expect)(corrupted.entries).toHaveLength(0);
                        return [4 /*yield*/, storage.restoreFromBackup(backupPath, vaultPath)];
                    case 7:
                        _a.sent();
                        return [4 /*yield*/, storage.load(vaultPath, masterPassword)];
                    case 8:
                        restored = _a.sent();
                        (0, vitest_1.expect)(restored.entries).toHaveLength(1);
                        return [4 /*yield*/, storage.getEntry(restored, entry.id, masterPassword)];
                    case 9:
                        recovered = _a.sent();
                        (0, vitest_1.expect)(recovered.data.password).toBe('super-secret');
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should maintain backup rotation', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, i, backupFiles, vaultBackups;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.save(vaultPath, vault, masterPassword)];
                    case 1:
                        _a.sent();
                        i = 0;
                        _a.label = 2;
                    case 2:
                        if (!(i < 7)) return [3 /*break*/, 6];
                        return [4 /*yield*/, new Promise(function (resolve) { return setTimeout(resolve, 20); })];
                    case 3:
                        _a.sent();
                        return [4 /*yield*/, storage.createBackup(vaultPath)];
                    case 4:
                        _a.sent();
                        _a.label = 5;
                    case 5:
                        i++;
                        return [3 /*break*/, 2];
                    case 6:
                        backupFiles = fs_1.default.readdirSync(backupDir);
                        vaultBackups = backupFiles.filter(function (f) { return f.startsWith('vault.dat') && f.endsWith('.backup'); });
                        (0, vitest_1.expect)(vaultBackups.length).toBeLessThanOrEqual(5);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('integrity and corruption scenarios', function () {
        (0, vitest_1.it)('should detect and report corrupted vault file', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, data, corrupted;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.save(vaultPath, vault, masterPassword)];
                    case 1:
                        _a.sent();
                        data = fs_1.default.readFileSync(vaultPath, 'utf-8');
                        corrupted = JSON.parse(data);
                        corrupted.metadata = null;
                        fs_1.default.writeFileSync(vaultPath, JSON.stringify(corrupted), { mode: 384 });
                        return [4 /*yield*/, (0, vitest_1.expect)(storage.load(vaultPath, masterPassword)).rejects.toThrow()];
                    case 2:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should verify vault integrity on demand', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, report;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.addEntry(vault, {
                                type: 'credential',
                                name: 'Test',
                                data: { value: 'test' },
                            }, masterPassword)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, storage.verifyIntegrity(vault, masterPassword)];
                    case 2:
                        report = _a.sent();
                        (0, vitest_1.expect)(report.valid).toBe(true);
                        (0, vitest_1.expect)(report.entryCount).toBe(1);
                        (0, vitest_1.expect)(report.corruptedEntries).toHaveLength(0);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('file permissions and security', function () {
        (0, vitest_1.it)('should create vault with correct permissions', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, stats, mode;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.save(vaultPath, vault, masterPassword)];
                    case 1:
                        _a.sent();
                        stats = fs_1.default.statSync(vaultPath);
                        mode = stats.mode & 511;
                        (0, vitest_1.expect)(mode).toBe(384);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should create backups with correct permissions', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, backupFiles, _i, backupFiles_1, file, stats, mode;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.save(vaultPath, vault, masterPassword)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, storage.createBackup(vaultPath)];
                    case 2:
                        _a.sent();
                        backupFiles = fs_1.default.readdirSync(backupDir);
                        for (_i = 0, backupFiles_1 = backupFiles; _i < backupFiles_1.length; _i++) {
                            file = backupFiles_1[_i];
                            stats = fs_1.default.statSync(path_1.default.join(backupDir, file));
                            mode = stats.mode & 511;
                            (0, vitest_1.expect)(mode).toBe(384);
                        }
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('edge cases', function () {
        (0, vitest_1.it)('should handle empty vault operations', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, loaded, entries;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.save(vaultPath, vault, masterPassword)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, storage.load(vaultPath, masterPassword)];
                    case 2:
                        loaded = _a.sent();
                        (0, vitest_1.expect)(loaded.entries).toHaveLength(0);
                        entries = storage.listEntries(loaded);
                        (0, vitest_1.expect)(entries).toHaveLength(0);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should handle special characters in data', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, specialData, entry, loaded, retrieved;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        specialData = {
                            content: 'Special chars: !@#$%^&*()_+-=[]{}|;:,.<>?',
                            unicode: 'Unicode: 你好世界 🌍 émojis',
                            multiline: 'Line 1\nLine 2\nLine 3',
                        };
                        return [4 /*yield*/, storage.addEntry(vault, {
                                type: 'note',
                                name: 'Special Characters Test',
                                data: specialData,
                            }, masterPassword)];
                    case 1:
                        entry = _a.sent();
                        return [4 /*yield*/, storage.save(vaultPath, vault, masterPassword)];
                    case 2:
                        _a.sent();
                        return [4 /*yield*/, storage.load(vaultPath, masterPassword)];
                    case 3:
                        loaded = _a.sent();
                        return [4 /*yield*/, storage.getEntry(loaded, entry.id, masterPassword)];
                    case 4:
                        retrieved = _a.sent();
                        (0, vitest_1.expect)(retrieved.data.content).toBe(specialData.content);
                        (0, vitest_1.expect)(retrieved.data.unicode).toBe(specialData.unicode);
                        (0, vitest_1.expect)(retrieved.data.multiline).toBe(specialData.multiline);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should handle concurrent entry modifications', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, entries;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, Promise.all([
                                storage.addEntry(vault, { type: 'credential', name: 'Entry 1', data: { v: 1 } }, masterPassword),
                                storage.addEntry(vault, { type: 'credential', name: 'Entry 2', data: { v: 2 } }, masterPassword),
                                storage.addEntry(vault, { type: 'credential', name: 'Entry 3', data: { v: 3 } }, masterPassword),
                            ])];
                    case 1:
                        entries = _a.sent();
                        (0, vitest_1.expect)(vault.entries).toHaveLength(3);
                        (0, vitest_1.expect)(entries.map(function (e) { return e.name; })).toContain('Entry 1');
                        (0, vitest_1.expect)(entries.map(function (e) { return e.name; })).toContain('Entry 2');
                        (0, vitest_1.expect)(entries.map(function (e) { return e.name; })).toContain('Entry 3');
                        return [2 /*return*/];
                }
            });
        }); });
    });
});
