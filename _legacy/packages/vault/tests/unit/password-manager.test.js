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
var password_manager_js_1 = require("../../src/password-manager.js");
(0, vitest_1.describe)('PasswordManager', function () {
    var manager;
    (0, vitest_1.beforeEach)(function () {
        manager = (0, password_manager_js_1.createPasswordManager)();
    });
    (0, vitest_1.afterEach)(function () {
        manager.destroy();
    });
    (0, vitest_1.describe)('unlock', function () {
        (0, vitest_1.it)('should unlock with correct password', function () { return __awaiter(void 0, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        (0, vitest_1.expect)(manager.isLocked()).toBe(true);
                        return [4 /*yield*/, manager.unlock('correct-password')];
                    case 1:
                        _a.sent();
                        (0, vitest_1.expect)(manager.isLocked()).toBe(false);
                        (0, vitest_1.expect)(manager.isUnlocked()).toBe(true);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should be idempotent if already unlocked', function () { return __awaiter(void 0, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, manager.unlock('password')];
                    case 1:
                        _a.sent();
                        (0, vitest_1.expect)(manager.isUnlocked()).toBe(true);
                        return [4 /*yield*/, manager.unlock('password')];
                    case 2:
                        _a.sent();
                        (0, vitest_1.expect)(manager.isUnlocked()).toBe(true);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('lock', function () {
        (0, vitest_1.it)('should lock when unlocked', function () { return __awaiter(void 0, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, manager.unlock('password')];
                    case 1:
                        _a.sent();
                        (0, vitest_1.expect)(manager.isUnlocked()).toBe(true);
                        manager.lock();
                        (0, vitest_1.expect)(manager.isLocked()).toBe(true);
                        (0, vitest_1.expect)(manager.isUnlocked()).toBe(false);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should be safe to call when already locked', function () {
            (0, vitest_1.expect)(manager.isLocked()).toBe(true);
            manager.lock();
            (0, vitest_1.expect)(manager.isLocked()).toBe(true);
        });
    });
    (0, vitest_1.describe)('encrypt', function () {
        (0, vitest_1.it)('should encrypt when unlocked', function () { return __awaiter(void 0, void 0, void 0, function () {
            var plaintext, encrypted;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, manager.unlock('password')];
                    case 1:
                        _a.sent();
                        plaintext = Buffer.from('secret data');
                        return [4 /*yield*/, manager.encrypt(plaintext)];
                    case 2:
                        encrypted = _a.sent();
                        (0, vitest_1.expect)(encrypted.ciphertext).toBeDefined();
                        (0, vitest_1.expect)(encrypted.iv).toBeDefined();
                        (0, vitest_1.expect)(encrypted.salt).toBeDefined();
                        (0, vitest_1.expect)(encrypted.tag).toBeDefined();
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should throw when locked', function () { return __awaiter(void 0, void 0, void 0, function () {
            var plaintext;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        plaintext = Buffer.from('secret data');
                        return [4 /*yield*/, (0, vitest_1.expect)(manager.encrypt(plaintext)).rejects.toThrow('Vault is locked')];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('decrypt', function () {
        (0, vitest_1.it)('should decrypt encrypted data', function () { return __awaiter(void 0, void 0, void 0, function () {
            var plaintext, encrypted, decrypted;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, manager.unlock('password')];
                    case 1:
                        _a.sent();
                        plaintext = Buffer.from('secret data');
                        return [4 /*yield*/, manager.encrypt(plaintext)];
                    case 2:
                        encrypted = _a.sent();
                        return [4 /*yield*/, manager.decrypt(encrypted)];
                    case 3:
                        decrypted = _a.sent();
                        (0, vitest_1.expect)(decrypted.toString()).toBe('secret data');
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should throw when locked', function () { return __awaiter(void 0, void 0, void 0, function () {
            var plaintext, encrypted;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, manager.unlock('password')];
                    case 1:
                        _a.sent();
                        plaintext = Buffer.from('secret data');
                        return [4 /*yield*/, manager.encrypt(plaintext)];
                    case 2:
                        encrypted = _a.sent();
                        manager.lock();
                        return [4 /*yield*/, (0, vitest_1.expect)(manager.decrypt(encrypted)).rejects.toThrow('Vault is locked')];
                    case 3:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('changePassword', function () {
        (0, vitest_1.it)('should change password and remain functional', function () { return __awaiter(void 0, void 0, void 0, function () {
            var oldPassword, newPassword, plaintext, encrypted, newPlaintext, newEncrypted, newDecrypted;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        oldPassword = 'old-password';
                        newPassword = 'new-password';
                        return [4 /*yield*/, manager.unlock(oldPassword)];
                    case 1:
                        _a.sent();
                        plaintext = Buffer.from('secret data');
                        return [4 /*yield*/, manager.encrypt(plaintext)];
                    case 2:
                        encrypted = _a.sent();
                        return [4 /*yield*/, manager.changePassword(oldPassword, newPassword)];
                    case 3:
                        _a.sent();
                        (0, vitest_1.expect)(manager.isUnlocked()).toBe(true);
                        newPlaintext = Buffer.from('new secret data');
                        return [4 /*yield*/, manager.encrypt(newPlaintext)];
                    case 4:
                        newEncrypted = _a.sent();
                        return [4 /*yield*/, manager.decrypt(newEncrypted)];
                    case 5:
                        newDecrypted = _a.sent();
                        (0, vitest_1.expect)(newDecrypted.toString()).toBe('new secret data');
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should throw when locked', function () { return __awaiter(void 0, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, (0, vitest_1.expect)(manager.changePassword('old', 'new')).rejects.toThrow('Vault is locked')];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('auto-lock', function () {
        (0, vitest_1.it)('should auto-lock after timeout', function () { return __awaiter(void 0, void 0, void 0, function () {
            var shortManager;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        shortManager = (0, password_manager_js_1.createPasswordManager)({
                            autoLockMs: 100, // 100ms for testing
                        });
                        return [4 /*yield*/, shortManager.unlock('password')];
                    case 1:
                        _a.sent();
                        (0, vitest_1.expect)(shortManager.isUnlocked()).toBe(true);
                        // Wait for auto-lock
                        return [4 /*yield*/, new Promise(function (resolve) { return setTimeout(resolve, 150); })];
                    case 2:
                        // Wait for auto-lock
                        _a.sent();
                        (0, vitest_1.expect)(shortManager.isLocked()).toBe(true);
                        shortManager.destroy();
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should reset timer on activity', function () { return __awaiter(void 0, void 0, void 0, function () {
            var shortManager;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        shortManager = (0, password_manager_js_1.createPasswordManager)({
                            autoLockMs: 200, // 200ms for testing
                        });
                        return [4 /*yield*/, shortManager.unlock('password')];
                    case 1:
                        _a.sent();
                        // Activity at 100ms
                        return [4 /*yield*/, new Promise(function (resolve) { return setTimeout(resolve, 100); })];
                    case 2:
                        // Activity at 100ms
                        _a.sent();
                        return [4 /*yield*/, shortManager.encrypt(Buffer.from('data'))];
                    case 3:
                        _a.sent();
                        // Should still be unlocked at 250ms (100ms + 200ms timeout)
                        return [4 /*yield*/, new Promise(function (resolve) { return setTimeout(resolve, 100); })];
                    case 4:
                        // Should still be unlocked at 250ms (100ms + 200ms timeout)
                        _a.sent();
                        (0, vitest_1.expect)(shortManager.isUnlocked()).toBe(true);
                        // Wait for auto-lock after activity
                        return [4 /*yield*/, new Promise(function (resolve) { return setTimeout(resolve, 250); })];
                    case 5:
                        // Wait for auto-lock after activity
                        _a.sent();
                        (0, vitest_1.expect)(shortManager.isLocked()).toBe(true);
                        shortManager.destroy();
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should return null for remaining lock time when locked', function () {
            (0, vitest_1.expect)(manager.getRemainingLockTime()).toBeNull();
        });
        (0, vitest_1.it)('should return timeout when unlocked', function () { return __awaiter(void 0, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, manager.unlock('password')];
                    case 1:
                        _a.sent();
                        (0, vitest_1.expect)(manager.getRemainingLockTime()).toBe(password_manager_js_1.DEFAULT_PASSWORD_MANAGER_CONFIG.autoLockMs);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('destroy', function () {
        (0, vitest_1.it)('should lock and cleanup when destroyed', function () { return __awaiter(void 0, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, manager.unlock('password')];
                    case 1:
                        _a.sent();
                        (0, vitest_1.expect)(manager.isUnlocked()).toBe(true);
                        manager.destroy();
                        (0, vitest_1.expect)(manager.isLocked()).toBe(true);
                        return [2 /*return*/];
                }
            });
        }); });
    });
});
(0, vitest_1.describe)('createPasswordManager', function () {
    (0, vitest_1.it)('should create with default config', function () {
        var manager = (0, password_manager_js_1.createPasswordManager)();
        (0, vitest_1.expect)(manager).toBeInstanceOf(password_manager_js_1.PasswordManager);
        (0, vitest_1.expect)(manager.isLocked()).toBe(true);
        manager.destroy();
    });
    (0, vitest_1.it)('should create with custom config', function () {
        var manager = (0, password_manager_js_1.createPasswordManager)({
            autoLockMs: 10000,
        });
        (0, vitest_1.expect)(manager).toBeInstanceOf(password_manager_js_1.PasswordManager);
        manager.destroy();
    });
});
