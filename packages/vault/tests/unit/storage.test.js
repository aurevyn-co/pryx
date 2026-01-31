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
var storage_types_js_1 = require("../src/storage-types.js");
var fs_1 = require("fs");
var path_1 = require("path");
var os_1 = require("os");
(0, vitest_1.describe)('VaultStorage', function () {
    var tempDir;
    var vaultPath;
    var backupDir;
    var storage;
    var password = 'test-password-123';
    (0, vitest_1.beforeEach)(function () {
        tempDir = fs_1.default.mkdtempSync(path_1.default.join(os_1.default.tmpdir(), 'vault-storage-test-'));
        vaultPath = path_1.default.join(tempDir, 'vault.dat');
        backupDir = path_1.default.join(tempDir, 'backups');
        storage = (0, storage_js_1.createVaultStorage)(backupDir);
    });
    (0, vitest_1.afterEach)(function () {
        fs_1.default.rmSync(tempDir, { recursive: true, force: true });
    });
    (0, vitest_1.describe)('createEmptyVault', function () {
        (0, vitest_1.it)('should create empty vault with correct structure', function () {
            var vault = storage.createEmptyVault();
            (0, vitest_1.expect)(vault.version).toBe(1);
            (0, vitest_1.expect)(vault.entries).toEqual([]);
            (0, vitest_1.expect)(vault.metadata).toBeDefined();
            (0, vitest_1.expect)(vault.metadata.algorithm).toBe('argon2id+aes-256-gcm');
            (0, vitest_1.expect)(vault.metadata.salt).toBeDefined();
            (0, vitest_1.expect)(vault.createdAt).toBeDefined();
            (0, vitest_1.expect)(vault.updatedAt).toBeDefined();
        });
    });
    (0, vitest_1.describe)('save and load', function () {
        (0, vitest_1.it)('should save and load vault successfully', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, loaded;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.save(vaultPath, vault, password)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, storage.load(vaultPath, password)];
                    case 2:
                        loaded = _a.sent();
                        (0, vitest_1.expect)(loaded.version).toBe(vault.version);
                        (0, vitest_1.expect)(loaded.metadata.salt).toBe(vault.metadata.salt);
                        (0, vitest_1.expect)(loaded.entries).toEqual([]);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should throw FileNotFoundError for non-existent file', function () { return __awaiter(void 0, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, (0, vitest_1.expect)(storage.load('/nonexistent/vault.dat', password)).rejects.toThrow(storage_types_js_1.FileNotFoundError)];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should throw CorruptedVaultError for invalid JSON', function () { return __awaiter(void 0, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        fs_1.default.writeFileSync(vaultPath, 'invalid json {', { mode: 384 });
                        return [4 /*yield*/, (0, vitest_1.expect)(storage.load(vaultPath, password)).rejects.toThrow(storage_types_js_1.CorruptedVaultError)];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should create parent directories when saving', function () { return __awaiter(void 0, void 0, void 0, function () {
            var nestedPath, vault;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        nestedPath = path_1.default.join(tempDir, 'nested', 'deep', 'vault.dat');
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.save(nestedPath, vault, password)];
                    case 1:
                        _a.sent();
                        (0, vitest_1.expect)(fs_1.default.existsSync(nestedPath)).toBe(true);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should set file permissions to 0o600', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, stats, mode;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.save(vaultPath, vault, password)];
                    case 1:
                        _a.sent();
                        stats = fs_1.default.statSync(vaultPath);
                        mode = stats.mode & 511;
                        (0, vitest_1.expect)(mode).toBe(384);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('addEntry', function () {
        (0, vitest_1.it)('should add entry to vault', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, entryData, entry;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        entryData = {
                            type: 'credential',
                            name: 'Test Entry',
                            data: { username: 'test', password: 'secret' },
                        };
                        return [4 /*yield*/, storage.addEntry(vault, entryData, password)];
                    case 1:
                        entry = _a.sent();
                        (0, vitest_1.expect)(entry.id).toBeDefined();
                        (0, vitest_1.expect)(entry.name).toBe('Test Entry');
                        (0, vitest_1.expect)(entry.type).toBe('credential');
                        (0, vitest_1.expect)(entry.encryptedData).toBeDefined();
                        (0, vitest_1.expect)(entry.iv).toBeDefined();
                        (0, vitest_1.expect)(entry.tag).toBeDefined();
                        (0, vitest_1.expect)(vault.entries).toHaveLength(1);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should use provided id if given', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, entryData, entry;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        entryData = {
                            id: 'custom-id',
                            type: 'credential',
                            name: 'Test Entry',
                            data: { username: 'test' },
                        };
                        return [4 /*yield*/, storage.addEntry(vault, entryData, password)];
                    case 1:
                        entry = _a.sent();
                        (0, vitest_1.expect)(entry.id).toBe('custom-id');
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should throw DuplicateEntryError for duplicate id', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, entryData;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        entryData = {
                            id: 'duplicate-id',
                            type: 'credential',
                            name: 'Test Entry',
                            data: { username: 'test' },
                        };
                        return [4 /*yield*/, storage.addEntry(vault, entryData, password)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, (0, vitest_1.expect)(storage.addEntry(vault, entryData, password)).rejects.toThrow(storage_types_js_1.DuplicateEntryError)];
                    case 2:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('getEntry', function () {
        (0, vitest_1.it)('should retrieve and decrypt entry', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, entryData, entry, retrieved;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        entryData = {
                            type: 'credential',
                            name: 'Test Entry',
                            data: { username: 'test', password: 'secret' },
                        };
                        return [4 /*yield*/, storage.addEntry(vault, entryData, password)];
                    case 1:
                        entry = _a.sent();
                        return [4 /*yield*/, storage.getEntry(vault, entry.id, password)];
                    case 2:
                        retrieved = _a.sent();
                        (0, vitest_1.expect)(retrieved.id).toBe(entry.id);
                        (0, vitest_1.expect)(retrieved.type).toBe('credential');
                        (0, vitest_1.expect)(retrieved.name).toBe('Test Entry');
                        (0, vitest_1.expect)(retrieved.data).toEqual({ username: 'test', password: 'secret' });
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should update access count and last accessed', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, entryData, entry;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        entryData = {
                            type: 'credential',
                            name: 'Test Entry',
                            data: { username: 'test' },
                        };
                        return [4 /*yield*/, storage.addEntry(vault, entryData, password)];
                    case 1:
                        entry = _a.sent();
                        (0, vitest_1.expect)(entry.accessCount).toBe(0);
                        (0, vitest_1.expect)(entry.lastAccessedAt).toBeUndefined();
                        return [4 /*yield*/, storage.getEntry(vault, entry.id, password)];
                    case 2:
                        _a.sent();
                        (0, vitest_1.expect)(vault.entries[0].accessCount).toBe(1);
                        (0, vitest_1.expect)(vault.entries[0].lastAccessedAt).toBeDefined();
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should throw EntryNotFoundError for non-existent entry', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, (0, vitest_1.expect)(storage.getEntry(vault, 'nonexistent', password)).rejects.toThrow(storage_types_js_1.EntryNotFoundError)];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('updateEntry', function () {
        (0, vitest_1.it)('should update entry name', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, entryData, entry, updated;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        entryData = {
                            type: 'credential',
                            name: 'Original Name',
                            data: { username: 'test' },
                        };
                        return [4 /*yield*/, storage.addEntry(vault, entryData, password)];
                    case 1:
                        entry = _a.sent();
                        return [4 /*yield*/, storage.updateEntry(vault, entry.id, { name: 'Updated Name' }, password)];
                    case 2:
                        updated = _a.sent();
                        (0, vitest_1.expect)(updated.name).toBe('Updated Name');
                        (0, vitest_1.expect)(vault.entries[0].name).toBe('Updated Name');
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should update entry data', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, entryData, entry, retrieved;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        entryData = {
                            type: 'credential',
                            name: 'Test Entry',
                            data: { username: 'original' },
                        };
                        return [4 /*yield*/, storage.addEntry(vault, entryData, password)];
                    case 1:
                        entry = _a.sent();
                        return [4 /*yield*/, storage.updateEntry(vault, entry.id, { data: { username: 'updated' } }, password)];
                    case 2:
                        _a.sent();
                        return [4 /*yield*/, storage.getEntry(vault, entry.id, password)];
                    case 3:
                        retrieved = _a.sent();
                        (0, vitest_1.expect)(retrieved.data).toEqual({ username: 'updated' });
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should throw EntryNotFoundError for non-existent entry', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, (0, vitest_1.expect)(storage.updateEntry(vault, 'nonexistent', { name: 'Test' }, password)).rejects.toThrow(storage_types_js_1.EntryNotFoundError)];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('deleteEntry', function () {
        (0, vitest_1.it)('should delete entry from vault', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, entryData, entry;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        entryData = {
                            type: 'credential',
                            name: 'Test Entry',
                            data: { username: 'test' },
                        };
                        return [4 /*yield*/, storage.addEntry(vault, entryData, password)];
                    case 1:
                        entry = _a.sent();
                        (0, vitest_1.expect)(vault.entries).toHaveLength(1);
                        return [4 /*yield*/, storage.deleteEntry(vault, entry.id)];
                    case 2:
                        _a.sent();
                        (0, vitest_1.expect)(vault.entries).toHaveLength(0);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should throw EntryNotFoundError for non-existent entry', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, (0, vitest_1.expect)(storage.deleteEntry(vault, 'nonexistent')).rejects.toThrow(storage_types_js_1.EntryNotFoundError)];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('listEntries', function () {
        (0, vitest_1.it)('should return metadata for all entries', function () {
            var vault = storage.createEmptyVault();
            vault.entries = [
                {
                    id: 'entry-1',
                    type: 'credential',
                    name: 'Entry 1',
                    encryptedData: 'data',
                    iv: 'iv',
                    tag: 'tag',
                    createdAt: '2024-01-01T00:00:00Z',
                    updatedAt: '2024-01-01T00:00:00Z',
                    accessCount: 5,
                    lastAccessedAt: '2024-01-02T00:00:00Z',
                },
            ];
            var entries = storage.listEntries(vault);
            (0, vitest_1.expect)(entries).toHaveLength(1);
            (0, vitest_1.expect)(entries[0].id).toBe('entry-1');
            (0, vitest_1.expect)(entries[0].name).toBe('Entry 1');
            (0, vitest_1.expect)(entries[0].type).toBe('credential');
            (0, vitest_1.expect)(entries[0].accessCount).toBe(5);
            (0, vitest_1.expect)(entries[0].encryptedData).toBeUndefined();
        });
    });
    (0, vitest_1.describe)('verifyIntegrity', function () {
        (0, vitest_1.it)('should return valid for correct vault', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, report;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        return [4 /*yield*/, storage.verifyIntegrity(vault)];
                    case 1:
                        report = _a.sent();
                        (0, vitest_1.expect)(report.valid).toBe(true);
                        (0, vitest_1.expect)(report.errors).toHaveLength(0);
                        (0, vitest_1.expect)(report.entryCount).toBe(0);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should detect invalid version', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, report;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        vault.version = 999;
                        return [4 /*yield*/, storage.verifyIntegrity(vault)];
                    case 1:
                        report = _a.sent();
                        (0, vitest_1.expect)(report.valid).toBe(false);
                        (0, vitest_1.expect)(report.errors).toContain('Unsupported vault version: 999. Expected: 1');
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should detect missing metadata', function () { return __awaiter(void 0, void 0, void 0, function () {
            var vault, report;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        vault = storage.createEmptyVault();
                        vault.metadata = null;
                        return [4 /*yield*/, storage.verifyIntegrity(vault)];
                    case 1:
                        report = _a.sent();
                        (0, vitest_1.expect)(report.valid).toBe(false);
                        (0, vitest_1.expect)(report.errors).toContain('Missing or invalid vault metadata');
                        return [2 /*return*/];
                }
            });
        }); });
    });
});
(0, vitest_1.describe)('createVaultStorage', function () {
    (0, vitest_1.it)('should create VaultStorage instance', function () {
        var storage = (0, storage_js_1.createVaultStorage)();
        (0, vitest_1.expect)(storage).toBeInstanceOf(storage_js_1.VaultStorage);
    });
    (0, vitest_1.it)('should use custom backup directory', function () { return __awaiter(void 0, void 0, void 0, function () {
        var tempDir, backupDir, storage, vault, vaultPath;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    tempDir = fs_1.default.mkdtempSync(path_1.default.join(os_1.default.tmpdir(), 'vault-test-'));
                    backupDir = path_1.default.join(tempDir, 'custom-backups');
                    storage = (0, storage_js_1.createVaultStorage)(backupDir);
                    vault = storage.createEmptyVault();
                    vaultPath = path_1.default.join(tempDir, 'vault.dat');
                    return [4 /*yield*/, storage.save(vaultPath, vault, 'password')];
                case 1:
                    _a.sent();
                    return [4 /*yield*/, storage.createBackup(vaultPath)];
                case 2:
                    _a.sent();
                    (0, vitest_1.expect)(fs_1.default.existsSync(backupDir)).toBe(true);
                    fs_1.default.rmSync(tempDir, { recursive: true, force: true });
                    return [2 /*return*/];
            }
        });
    }); });
});
