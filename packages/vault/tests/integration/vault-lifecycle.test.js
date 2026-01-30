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
var vault_js_1 = require("../../src/vault.js");
var types_js_1 = require("../../src/types.js");
(0, vitest_1.describe)('Vault Lifecycle Integration', function () {
    var vault;
    (0, vitest_1.beforeEach)(function () {
        vault = new vault_js_1.Vault();
    });
    (0, vitest_1.it)('should complete full lifecycle: init → encrypt → decrypt → clear', function () { return __awaiter(void 0, void 0, void 0, function () {
        var password, plaintext, encrypted, decrypted;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    password = 'master-password';
                    plaintext = Buffer.from('sensitive credential data');
                    return [4 /*yield*/, vault.initialize(password)];
                case 1:
                    _a.sent();
                    (0, vitest_1.expect)(vault.isInitialized).toBe(true);
                    return [4 /*yield*/, vault.encrypt(plaintext)];
                case 2:
                    encrypted = _a.sent();
                    (0, vitest_1.expect)(encrypted.ciphertext.length).toBeGreaterThan(0);
                    return [4 /*yield*/, vault.decrypt(encrypted)];
                case 3:
                    decrypted = _a.sent();
                    (0, vitest_1.expect)(decrypted.toString()).toBe(plaintext.toString());
                    vault.clearKey();
                    (0, vitest_1.expect)(vault.isInitialized).toBe(false);
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should handle multiple encryption operations', function () { return __awaiter(void 0, void 0, void 0, function () {
        var data1, data2, data3, encrypted1, encrypted2, encrypted3, decrypted1, decrypted2, decrypted3;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, vault.initialize('password')];
                case 1:
                    _a.sent();
                    data1 = Buffer.from('first secret');
                    data2 = Buffer.from('second secret');
                    data3 = Buffer.from('third secret');
                    return [4 /*yield*/, vault.encrypt(data1)];
                case 2:
                    encrypted1 = _a.sent();
                    return [4 /*yield*/, vault.encrypt(data2)];
                case 3:
                    encrypted2 = _a.sent();
                    return [4 /*yield*/, vault.encrypt(data3)];
                case 4:
                    encrypted3 = _a.sent();
                    return [4 /*yield*/, vault.decrypt(encrypted1)];
                case 5:
                    decrypted1 = _a.sent();
                    return [4 /*yield*/, vault.decrypt(encrypted2)];
                case 6:
                    decrypted2 = _a.sent();
                    return [4 /*yield*/, vault.decrypt(encrypted3)];
                case 7:
                    decrypted3 = _a.sent();
                    (0, vitest_1.expect)(decrypted1.toString()).toBe('first secret');
                    (0, vitest_1.expect)(decrypted2.toString()).toBe('second secret');
                    (0, vitest_1.expect)(decrypted3.toString()).toBe('third secret');
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should use unique IV for each encryption', function () { return __awaiter(void 0, void 0, void 0, function () {
        var data, encrypted1, encrypted2;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, vault.initialize('password')];
                case 1:
                    _a.sent();
                    data = Buffer.from('same data');
                    return [4 /*yield*/, vault.encrypt(data)];
                case 2:
                    encrypted1 = _a.sent();
                    return [4 /*yield*/, vault.encrypt(data)];
                case 3:
                    encrypted2 = _a.sent();
                    (0, vitest_1.expect)(encrypted1.iv.toString('hex')).not.toBe(encrypted2.iv.toString('hex'));
                    (0, vitest_1.expect)(encrypted1.ciphertext.toString('hex')).not.toBe(encrypted2.ciphertext.toString('hex'));
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should handle concurrent encryption operations', function () { return __awaiter(void 0, void 0, void 0, function () {
        var promises, results, i;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, vault.initialize('password')];
                case 1:
                    _a.sent();
                    promises = Array.from({ length: 10 }, function (_, i) { return __awaiter(void 0, void 0, void 0, function () {
                        var data, encrypted, decrypted;
                        return __generator(this, function (_a) {
                            switch (_a.label) {
                                case 0:
                                    data = Buffer.from("data-".concat(i));
                                    return [4 /*yield*/, vault.encrypt(data)];
                                case 1:
                                    encrypted = _a.sent();
                                    return [4 /*yield*/, vault.decrypt(encrypted)];
                                case 2:
                                    decrypted = _a.sent();
                                    return [2 /*return*/, decrypted.toString()];
                            }
                        });
                    }); });
                    return [4 /*yield*/, Promise.all(promises)];
                case 2:
                    results = _a.sent();
                    for (i = 0; i < 10; i++) {
                        (0, vitest_1.expect)(results[i]).toBe("data-".concat(i));
                    }
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should handle large data encryption', function () { return __awaiter(void 0, void 0, void 0, function () {
        var largeData, encrypted, decrypted;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, vault.initialize('password')];
                case 1:
                    _a.sent();
                    largeData = Buffer.alloc(1024 * 1024, 0x42);
                    return [4 /*yield*/, vault.encrypt(largeData)];
                case 2:
                    encrypted = _a.sent();
                    return [4 /*yield*/, vault.decrypt(encrypted)];
                case 3:
                    decrypted = _a.sent();
                    (0, vitest_1.expect)(decrypted.toString('hex')).toBe(largeData.toString('hex'));
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should handle empty data', function () { return __awaiter(void 0, void 0, void 0, function () {
        var emptyData, encrypted, decrypted;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, vault.initialize('password')];
                case 1:
                    _a.sent();
                    emptyData = Buffer.alloc(0);
                    return [4 /*yield*/, vault.encrypt(emptyData)];
                case 2:
                    encrypted = _a.sent();
                    return [4 /*yield*/, vault.decrypt(encrypted)];
                case 3:
                    decrypted = _a.sent();
                    (0, vitest_1.expect)(decrypted.length).toBe(0);
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should handle unicode data', function () { return __awaiter(void 0, void 0, void 0, function () {
        var unicodeData, encrypted, decrypted;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, vault.initialize('password')];
                case 1:
                    _a.sent();
                    unicodeData = Buffer.from('Hello 世界 🌍 مرحبا', 'utf8');
                    return [4 /*yield*/, vault.encrypt(unicodeData)];
                case 2:
                    encrypted = _a.sent();
                    return [4 /*yield*/, vault.decrypt(encrypted)];
                case 3:
                    decrypted = _a.sent();
                    (0, vitest_1.expect)(decrypted.toString('utf8')).toBe('Hello 世界 🌍 مرحبا');
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should fail decryption with wrong password', function () { return __awaiter(void 0, void 0, void 0, function () {
        var correctPassword, wrongPassword, plaintext, encrypted, wrongVault;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    correctPassword = 'correct-password';
                    wrongPassword = 'wrong-password';
                    plaintext = Buffer.from('secret');
                    return [4 /*yield*/, vault.initialize(correctPassword)];
                case 1:
                    _a.sent();
                    return [4 /*yield*/, vault.encrypt(plaintext)];
                case 2:
                    encrypted = _a.sent();
                    wrongVault = new vault_js_1.Vault();
                    return [4 /*yield*/, wrongVault.initialize(wrongPassword)];
                case 3:
                    _a.sent();
                    return [4 /*yield*/, (0, vitest_1.expect)(wrongVault.decrypt(encrypted)).rejects.toThrow(types_js_1.DecryptionError)];
                case 4:
                    _a.sent();
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should fail decryption with tampered data', function () { return __awaiter(void 0, void 0, void 0, function () {
        var plaintext, encrypted;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, vault.initialize('password')];
                case 1:
                    _a.sent();
                    plaintext = Buffer.from('secret');
                    return [4 /*yield*/, vault.encrypt(plaintext)];
                case 2:
                    encrypted = _a.sent();
                    encrypted.ciphertext[0] ^= 0xFF;
                    return [4 /*yield*/, (0, vitest_1.expect)(vault.decrypt(encrypted)).rejects.toThrow(types_js_1.DecryptionError)];
                case 3:
                    _a.sent();
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should recover from error state', function () { return __awaiter(void 0, void 0, void 0, function () {
        var plaintext, encrypted, decrypted;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, vault.initialize('password')];
                case 1:
                    _a.sent();
                    plaintext = Buffer.from('test');
                    return [4 /*yield*/, vault.encrypt(plaintext)];
                case 2:
                    encrypted = _a.sent();
                    encrypted.ciphertext[0] ^= 0xFF;
                    return [4 /*yield*/, (0, vitest_1.expect)(vault.decrypt(encrypted)).rejects.toThrow()];
                case 3:
                    _a.sent();
                    encrypted.ciphertext[0] ^= 0xFF;
                    return [4 /*yield*/, vault.decrypt(encrypted)];
                case 4:
                    decrypted = _a.sent();
                    (0, vitest_1.expect)(decrypted.toString()).toBe('test');
                    return [2 /*return*/];
            }
        });
    }); });
});
