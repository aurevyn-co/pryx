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
(0, vitest_1.describe)('User Workflow E2E', function () {
    (0, vitest_1.it)('should complete full user workflow: create vault → store credential → retrieve credential', function () { return __awaiter(void 0, void 0, void 0, function () {
        var masterPassword, credential, vault, encrypted, decrypted, retrieved;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    masterPassword = 'MyStr0ng!Mast3r#P@ssw0rd';
                    credential = JSON.stringify({
                        service: 'openai',
                        apiKey: 'sk-1234567890abcdef',
                        organization: 'org-test123',
                    });
                    vault = new vault_js_1.Vault();
                    return [4 /*yield*/, vault.initialize(masterPassword)];
                case 1:
                    _a.sent();
                    return [4 /*yield*/, vault.encrypt(Buffer.from(credential))];
                case 2:
                    encrypted = _a.sent();
                    return [4 /*yield*/, vault.decrypt(encrypted)];
                case 3:
                    decrypted = _a.sent();
                    retrieved = JSON.parse(decrypted.toString());
                    (0, vitest_1.expect)(retrieved.service).toBe('openai');
                    (0, vitest_1.expect)(retrieved.apiKey).toBe('sk-1234567890abcdef');
                    (0, vitest_1.expect)(retrieved.organization).toBe('org-test123');
                    vault.clearKey();
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should handle wrong password scenario', function () { return __awaiter(void 0, void 0, void 0, function () {
        var correctPassword, wrongPassword, secret, vault, encrypted, wrongVault;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    correctPassword = 'CorrectP@ss123';
                    wrongPassword = 'WrongP@ss456';
                    secret = 'my-secret-data';
                    vault = new vault_js_1.Vault();
                    return [4 /*yield*/, vault.initialize(correctPassword)];
                case 1:
                    _a.sent();
                    return [4 /*yield*/, vault.encrypt(Buffer.from(secret))];
                case 2:
                    encrypted = _a.sent();
                    vault.clearKey();
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
    (0, vitest_1.it)('should handle multiple credentials', function () { return __awaiter(void 0, void 0, void 0, function () {
        var password, credentials, vault, encryptedCredentials, _i, encryptedCredentials_1, _a, service, encrypted, decrypted, parsed;
        return __generator(this, function (_b) {
            switch (_b.label) {
                case 0:
                    password = 'master-password';
                    credentials = [
                        { service: 'openai', key: 'sk-openai-123' },
                        { service: 'anthropic', key: 'sk-anthropic-456' },
                        { service: 'google', key: 'sk-google-789' },
                    ];
                    vault = new vault_js_1.Vault();
                    return [4 /*yield*/, vault.initialize(password)];
                case 1:
                    _b.sent();
                    return [4 /*yield*/, Promise.all(credentials.map(function (cred) { return __awaiter(void 0, void 0, void 0, function () {
                            var _a;
                            return __generator(this, function (_b) {
                                switch (_b.label) {
                                    case 0:
                                        _a = {
                                            service: cred.service
                                        };
                                        return [4 /*yield*/, vault.encrypt(Buffer.from(JSON.stringify(cred)))];
                                    case 1: return [2 /*return*/, (_a.encrypted = _b.sent(),
                                            _a)];
                                }
                            });
                        }); }))];
                case 2:
                    encryptedCredentials = _b.sent();
                    vault.clearKey();
                    _i = 0, encryptedCredentials_1 = encryptedCredentials;
                    _b.label = 3;
                case 3:
                    if (!(_i < encryptedCredentials_1.length)) return [3 /*break*/, 6];
                    _a = encryptedCredentials_1[_i], service = _a.service, encrypted = _a.encrypted;
                    return [4 /*yield*/, (0, vault_js_1.decryptWithPassword)(encrypted, password)];
                case 4:
                    decrypted = _b.sent();
                    parsed = JSON.parse(decrypted.toString());
                    (0, vitest_1.expect)(parsed.service).toBe(service);
                    _b.label = 5;
                case 5:
                    _i++;
                    return [3 /*break*/, 3];
                case 6: return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should encrypt and decrypt using convenience functions', function () { return __awaiter(void 0, void 0, void 0, function () {
        var password, data, encrypted, decrypted;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    password = 'simple-password';
                    data = Buffer.from('test data');
                    return [4 /*yield*/, (0, vault_js_1.encryptWithPassword)(data, password)];
                case 1:
                    encrypted = _a.sent();
                    return [4 /*yield*/, (0, vault_js_1.decryptWithPassword)(encrypted, password)];
                case 2:
                    decrypted = _a.sent();
                    (0, vitest_1.expect)(decrypted.toString()).toBe('test data');
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should fail with wrong password using convenience functions', function () { return __awaiter(void 0, void 0, void 0, function () {
        var correctPassword, wrongPassword, data, encrypted;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    correctPassword = 'correct';
                    wrongPassword = 'wrong';
                    data = Buffer.from('secret');
                    return [4 /*yield*/, (0, vault_js_1.encryptWithPassword)(data, correctPassword)];
                case 1:
                    encrypted = _a.sent();
                    return [4 /*yield*/, (0, vitest_1.expect)((0, vault_js_1.decryptWithPassword)(encrypted, wrongPassword)).rejects.toThrow(types_js_1.DecryptionError)];
                case 2:
                    _a.sent();
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should handle binary data', function () { return __awaiter(void 0, void 0, void 0, function () {
        var password, binaryData, encrypted, decrypted;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    password = 'binary-test';
                    binaryData = Buffer.from([0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD]);
                    return [4 /*yield*/, (0, vault_js_1.encryptWithPassword)(binaryData, password)];
                case 1:
                    encrypted = _a.sent();
                    return [4 /*yield*/, (0, vault_js_1.decryptWithPassword)(encrypted, password)];
                case 2:
                    decrypted = _a.sent();
                    (0, vitest_1.expect)(Buffer.compare(decrypted, binaryData)).toBe(0);
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should handle special characters in password', function () { return __awaiter(void 0, void 0, void 0, function () {
        var specialPassword, data, encrypted, decrypted;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    specialPassword = 'p@$$w0rd!#$%^*()_+-=[]{}|;:,.?<>';
                    data = Buffer.from('test');
                    return [4 /*yield*/, (0, vault_js_1.encryptWithPassword)(data, specialPassword)];
                case 1:
                    encrypted = _a.sent();
                    return [4 /*yield*/, (0, vault_js_1.decryptWithPassword)(encrypted, specialPassword)];
                case 2:
                    decrypted = _a.sent();
                    (0, vitest_1.expect)(decrypted.toString()).toBe('test');
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should handle very long password', function () { return __awaiter(void 0, void 0, void 0, function () {
        var longPassword, data, encrypted, decrypted;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    longPassword = 'a'.repeat(1000);
                    data = Buffer.from('test');
                    return [4 /*yield*/, (0, vault_js_1.encryptWithPassword)(data, longPassword)];
                case 1:
                    encrypted = _a.sent();
                    return [4 /*yield*/, (0, vault_js_1.decryptWithPassword)(encrypted, longPassword)];
                case 2:
                    decrypted = _a.sent();
                    (0, vitest_1.expect)(decrypted.toString()).toBe('test');
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should handle empty password', function () { return __awaiter(void 0, void 0, void 0, function () {
        var emptyPassword, data, encrypted, decrypted;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    emptyPassword = '';
                    data = Buffer.from('test');
                    return [4 /*yield*/, (0, vault_js_1.encryptWithPassword)(data, emptyPassword)];
                case 1:
                    encrypted = _a.sent();
                    return [4 /*yield*/, (0, vault_js_1.decryptWithPassword)(encrypted, emptyPassword)];
                case 2:
                    decrypted = _a.sent();
                    (0, vitest_1.expect)(decrypted.toString()).toBe('test');
                    return [2 /*return*/];
            }
        });
    }); });
});
(0, vitest_1.describe)('Performance Benchmarks', function () {
    (0, vitest_1.it)('should complete key derivation in reasonable time', function () { return __awaiter(void 0, void 0, void 0, function () {
        var password, iterations, start, i, vault, duration, avgTime;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    password = 'benchmark-password';
                    iterations = 5;
                    start = performance.now();
                    i = 0;
                    _a.label = 1;
                case 1:
                    if (!(i < iterations)) return [3 /*break*/, 4];
                    vault = new vault_js_1.Vault();
                    return [4 /*yield*/, vault.initialize(password)];
                case 2:
                    _a.sent();
                    vault.clearKey();
                    _a.label = 3;
                case 3:
                    i++;
                    return [3 /*break*/, 1];
                case 4:
                    duration = performance.now() - start;
                    avgTime = duration / iterations;
                    (0, vitest_1.expect)(avgTime).toBeLessThan(1000);
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should encrypt data efficiently', function () { return __awaiter(void 0, void 0, void 0, function () {
        var vault, data, iterations, start, i, duration, avgTime;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    vault = new vault_js_1.Vault();
                    return [4 /*yield*/, vault.initialize('password')];
                case 1:
                    _a.sent();
                    data = Buffer.from('x'.repeat(10000));
                    iterations = 100;
                    start = performance.now();
                    i = 0;
                    _a.label = 2;
                case 2:
                    if (!(i < iterations)) return [3 /*break*/, 5];
                    return [4 /*yield*/, vault.encrypt(data)];
                case 3:
                    _a.sent();
                    _a.label = 4;
                case 4:
                    i++;
                    return [3 /*break*/, 2];
                case 5:
                    duration = performance.now() - start;
                    avgTime = duration / iterations;
                    (0, vitest_1.expect)(avgTime).toBeLessThan(10);
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should decrypt data efficiently', function () { return __awaiter(void 0, void 0, void 0, function () {
        var vault, data, encrypted, iterations, start, i, duration, avgTime;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    vault = new vault_js_1.Vault();
                    return [4 /*yield*/, vault.initialize('password')];
                case 1:
                    _a.sent();
                    data = Buffer.from('x'.repeat(10000));
                    return [4 /*yield*/, vault.encrypt(data)];
                case 2:
                    encrypted = _a.sent();
                    iterations = 100;
                    start = performance.now();
                    i = 0;
                    _a.label = 3;
                case 3:
                    if (!(i < iterations)) return [3 /*break*/, 6];
                    return [4 /*yield*/, vault.decrypt(encrypted)];
                case 4:
                    _a.sent();
                    _a.label = 5;
                case 5:
                    i++;
                    return [3 /*break*/, 3];
                case 6:
                    duration = performance.now() - start;
                    avgTime = duration / iterations;
                    (0, vitest_1.expect)(avgTime).toBeLessThan(5);
                    return [2 /*return*/];
            }
        });
    }); });
});
