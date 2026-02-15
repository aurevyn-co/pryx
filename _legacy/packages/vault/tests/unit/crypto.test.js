"use strict";
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
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
var crypto_js_1 = require("../../src/crypto.js");
var types_js_1 = require("../../src/types.js");
(0, vitest_1.describe)('deriveKey', function () {
    (0, vitest_1.it)('should derive 32-byte key from password', function () { return __awaiter(void 0, void 0, void 0, function () {
        var password, salt, key;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    password = 'test-password';
                    salt = (0, crypto_js_1.generateSalt)();
                    return [4 /*yield*/, (0, crypto_js_1.deriveKey)(password, salt)];
                case 1:
                    key = _a.sent();
                    (0, vitest_1.expect)(key).toBeInstanceOf(Buffer);
                    (0, vitest_1.expect)(key.length).toBe(32);
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should derive same key with same password and salt', function () { return __awaiter(void 0, void 0, void 0, function () {
        var password, salt, key1, key2;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    password = 'test-password';
                    salt = (0, crypto_js_1.generateSalt)();
                    return [4 /*yield*/, (0, crypto_js_1.deriveKey)(password, salt)];
                case 1:
                    key1 = _a.sent();
                    return [4 /*yield*/, (0, crypto_js_1.deriveKey)(password, salt)];
                case 2:
                    key2 = _a.sent();
                    (0, vitest_1.expect)(key1.toString('hex')).toBe(key2.toString('hex'));
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should derive different keys with different passwords', function () { return __awaiter(void 0, void 0, void 0, function () {
        var salt, key1, key2;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    salt = (0, crypto_js_1.generateSalt)();
                    return [4 /*yield*/, (0, crypto_js_1.deriveKey)('password1', salt)];
                case 1:
                    key1 = _a.sent();
                    return [4 /*yield*/, (0, crypto_js_1.deriveKey)('password2', salt)];
                case 2:
                    key2 = _a.sent();
                    (0, vitest_1.expect)(key1.toString('hex')).not.toBe(key2.toString('hex'));
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should derive different keys with different salts', function () { return __awaiter(void 0, void 0, void 0, function () {
        var password, salt1, salt2, key1, key2;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    password = 'test-password';
                    salt1 = (0, crypto_js_1.generateSalt)();
                    salt2 = (0, crypto_js_1.generateSalt)();
                    return [4 /*yield*/, (0, crypto_js_1.deriveKey)(password, salt1)];
                case 1:
                    key1 = _a.sent();
                    return [4 /*yield*/, (0, crypto_js_1.deriveKey)(password, salt2)];
                case 2:
                    key2 = _a.sent();
                    (0, vitest_1.expect)(key1.toString('hex')).not.toBe(key2.toString('hex'));
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should use custom config parameters', function () { return __awaiter(void 0, void 0, void 0, function () {
        var password, salt, config, key;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    password = 'test-password';
                    salt = (0, crypto_js_1.generateSalt)();
                    config = __assign(__assign({}, types_js_1.DEFAULT_VAULT_CONFIG), { memoryCost: 32768, timeCost: 2 });
                    return [4 /*yield*/, (0, crypto_js_1.deriveKey)(password, salt, config)];
                case 1:
                    key = _a.sent();
                    (0, vitest_1.expect)(key.length).toBe(32);
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should throw on unsupported algorithm', function () { return __awaiter(void 0, void 0, void 0, function () {
        var password, salt, config;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    password = 'test-password';
                    salt = (0, crypto_js_1.generateSalt)();
                    config = __assign(__assign({}, types_js_1.DEFAULT_VAULT_CONFIG), { algorithm: 'pbkdf2' });
                    return [4 /*yield*/, (0, vitest_1.expect)((0, crypto_js_1.deriveKey)(password, salt, config)).rejects.toThrow('Unsupported algorithm')];
                case 1:
                    _a.sent();
                    return [2 /*return*/];
            }
        });
    }); });
});
(0, vitest_1.describe)('generateSalt', function () {
    (0, vitest_1.it)('should generate salt of default length', function () {
        var salt = (0, crypto_js_1.generateSalt)();
        (0, vitest_1.expect)(salt).toBeInstanceOf(Buffer);
        (0, vitest_1.expect)(salt.length).toBe(32);
    });
    (0, vitest_1.it)('should generate salt of custom length', function () {
        var salt = (0, crypto_js_1.generateSalt)(16);
        (0, vitest_1.expect)(salt.length).toBe(16);
    });
    (0, vitest_1.it)('should generate unique salts', function () {
        var salt1 = (0, crypto_js_1.generateSalt)();
        var salt2 = (0, crypto_js_1.generateSalt)();
        (0, vitest_1.expect)(salt1.toString('hex')).not.toBe(salt2.toString('hex'));
    });
});
(0, vitest_1.describe)('generateIV', function () {
    (0, vitest_1.it)('should generate IV of default length (12 bytes)', function () {
        var iv = (0, crypto_js_1.generateIV)();
        (0, vitest_1.expect)(iv).toBeInstanceOf(Buffer);
        (0, vitest_1.expect)(iv.length).toBe(12);
    });
    (0, vitest_1.it)('should generate unique IVs', function () {
        var iv1 = (0, crypto_js_1.generateIV)();
        var iv2 = (0, crypto_js_1.generateIV)();
        (0, vitest_1.expect)(iv1.toString('hex')).not.toBe(iv2.toString('hex'));
    });
});
(0, vitest_1.describe)('encrypt', function () {
    (0, vitest_1.it)('should encrypt data and return ciphertext and tag', function () {
        var plaintext = Buffer.from('Hello, World!');
        var key = Buffer.alloc(32, 0x42);
        var iv = (0, crypto_js_1.generateIV)();
        var result = (0, crypto_js_1.encrypt)(plaintext, key, iv);
        (0, vitest_1.expect)(result.ciphertext).toBeInstanceOf(Buffer);
        (0, vitest_1.expect)(result.ciphertext.length).toBeGreaterThan(0);
        (0, vitest_1.expect)(result.tag).toBeInstanceOf(Buffer);
        (0, vitest_1.expect)(result.tag.length).toBe(16);
    });
    (0, vitest_1.it)('should throw on invalid key length', function () {
        var plaintext = Buffer.from('test');
        var key = Buffer.alloc(16, 0x42);
        var iv = (0, crypto_js_1.generateIV)();
        (0, vitest_1.expect)(function () { return (0, crypto_js_1.encrypt)(plaintext, key, iv); }).toThrow('Invalid key length');
    });
    (0, vitest_1.it)('should throw on invalid IV length', function () {
        var plaintext = Buffer.from('test');
        var key = Buffer.alloc(32, 0x42);
        var iv = Buffer.alloc(16, 0x42);
        (0, vitest_1.expect)(function () { return (0, crypto_js_1.encrypt)(plaintext, key, iv); }).toThrow('Invalid IV length');
    });
});
(0, vitest_1.describe)('decrypt', function () {
    (0, vitest_1.it)('should decrypt encrypted data correctly', function () {
        var plaintext = Buffer.from('Hello, World!');
        var key = Buffer.alloc(32, 0x42);
        var iv = (0, crypto_js_1.generateIV)();
        var encrypted = (0, crypto_js_1.encrypt)(plaintext, key, iv);
        var decrypted = (0, crypto_js_1.decrypt)(encrypted.ciphertext, key, iv, encrypted.tag);
        (0, vitest_1.expect)(decrypted.toString()).toBe(plaintext.toString());
    });
    (0, vitest_1.it)('should throw on wrong key', function () {
        var plaintext = Buffer.from('secret');
        var key1 = Buffer.alloc(32, 0x42);
        var key2 = Buffer.alloc(32, 0x24);
        var iv = (0, crypto_js_1.generateIV)();
        var encrypted = (0, crypto_js_1.encrypt)(plaintext, key1, iv);
        (0, vitest_1.expect)(function () { return (0, crypto_js_1.decrypt)(encrypted.ciphertext, key2, iv, encrypted.tag); }).toThrow(types_js_1.DecryptionError);
    });
    (0, vitest_1.it)('should throw on tampered ciphertext', function () {
        var plaintext = Buffer.from('secret');
        var key = Buffer.alloc(32, 0x42);
        var iv = (0, crypto_js_1.generateIV)();
        var encrypted = (0, crypto_js_1.encrypt)(plaintext, key, iv);
        encrypted.ciphertext[0] ^= 0xFF;
        (0, vitest_1.expect)(function () { return (0, crypto_js_1.decrypt)(encrypted.ciphertext, key, iv, encrypted.tag); }).toThrow(types_js_1.DecryptionError);
    });
    (0, vitest_1.it)('should throw on tampered tag', function () {
        var plaintext = Buffer.from('secret');
        var key = Buffer.alloc(32, 0x42);
        var iv = (0, crypto_js_1.generateIV)();
        var encrypted = (0, crypto_js_1.encrypt)(plaintext, key, iv);
        encrypted.tag[0] ^= 0xFF;
        (0, vitest_1.expect)(function () { return (0, crypto_js_1.decrypt)(encrypted.ciphertext, key, iv, encrypted.tag); }).toThrow(types_js_1.DecryptionError);
    });
    (0, vitest_1.it)('should throw on invalid key length', function () {
        var ciphertext = Buffer.from('test');
        var key = Buffer.alloc(16, 0x42);
        var iv = (0, crypto_js_1.generateIV)();
        var tag = Buffer.alloc(16, 0x42);
        (0, vitest_1.expect)(function () { return (0, crypto_js_1.decrypt)(ciphertext, key, iv, tag); }).toThrow('Invalid key length');
    });
});
(0, vitest_1.describe)('secureClear', function () {
    (0, vitest_1.it)('should clear buffer contents', function () {
        var buffer = Buffer.from('sensitive data');
        (0, crypto_js_1.secureClear)(buffer);
        (0, vitest_1.expect)(buffer.toString()).toBe('\x00'.repeat(buffer.length));
    });
});
(0, vitest_1.describe)('secureCompare', function () {
    (0, vitest_1.it)('should return true for identical buffers', function () {
        var buf1 = Buffer.from('test');
        var buf2 = Buffer.from('test');
        (0, vitest_1.expect)((0, crypto_js_1.secureCompare)(buf1, buf2)).toBe(true);
    });
    (0, vitest_1.it)('should return false for different buffers', function () {
        var buf1 = Buffer.from('test1');
        var buf2 = Buffer.from('test2');
        (0, vitest_1.expect)((0, crypto_js_1.secureCompare)(buf1, buf2)).toBe(false);
    });
    (0, vitest_1.it)('should return false for different lengths', function () {
        var buf1 = Buffer.from('test');
        var buf2 = Buffer.from('testing');
        (0, vitest_1.expect)((0, crypto_js_1.secureCompare)(buf1, buf2)).toBe(false);
    });
});
(0, vitest_1.describe)('serializeEncryptedData', function () {
    (0, vitest_1.it)('should serialize to base64 strings', function () {
        var data = {
            ciphertext: Buffer.from('ciphertext'),
            iv: Buffer.from('iv'),
            salt: Buffer.from('salt'),
            tag: Buffer.from('tag'),
            version: 1,
        };
        var serialized = (0, crypto_js_1.serializeEncryptedData)(data);
        (0, vitest_1.expect)(typeof serialized.ciphertext).toBe('string');
        (0, vitest_1.expect)(typeof serialized.iv).toBe('string');
        (0, vitest_1.expect)(typeof serialized.salt).toBe('string');
        (0, vitest_1.expect)(typeof serialized.tag).toBe('string');
        (0, vitest_1.expect)(serialized.version).toBe(1);
    });
});
(0, vitest_1.describe)('deserializeEncryptedData', function () {
    (0, vitest_1.it)('should deserialize from base64 strings', function () {
        var original = {
            ciphertext: Buffer.from('ciphertext'),
            iv: Buffer.from('iv'),
            salt: Buffer.from('salt'),
            tag: Buffer.from('tag'),
            version: 1,
        };
        var serialized = (0, crypto_js_1.serializeEncryptedData)(original);
        var deserialized = (0, crypto_js_1.deserializeEncryptedData)(serialized);
        (0, vitest_1.expect)(deserialized.ciphertext.toString()).toBe(original.ciphertext.toString());
        (0, vitest_1.expect)(deserialized.iv.toString()).toBe(original.iv.toString());
        (0, vitest_1.expect)(deserialized.salt.toString()).toBe(original.salt.toString());
        (0, vitest_1.expect)(deserialized.tag.toString()).toBe(original.tag.toString());
        (0, vitest_1.expect)(deserialized.version).toBe(original.version);
    });
});
(0, vitest_1.describe)('Error classes', function () {
    (0, vitest_1.it)('should create VaultError', function () {
        var error = new types_js_1.VaultError('test message');
        (0, vitest_1.expect)(error.message).toBe('test message');
        (0, vitest_1.expect)(error.name).toBe('VaultError');
    });
    (0, vitest_1.it)('should create InvalidPasswordError', function () {
        var error = new types_js_1.InvalidPasswordError();
        (0, vitest_1.expect)(error.message).toBe('Invalid password provided');
        (0, vitest_1.expect)(error.name).toBe('InvalidPasswordError');
    });
    (0, vitest_1.it)('should create CorruptedDataError with default message', function () {
        var error = new types_js_1.CorruptedDataError();
        (0, vitest_1.expect)(error.message).toBe('Data appears to be corrupted');
        (0, vitest_1.expect)(error.name).toBe('CorruptedDataError');
    });
    (0, vitest_1.it)('should create CorruptedDataError with custom message', function () {
        var error = new types_js_1.CorruptedDataError('custom error');
        (0, vitest_1.expect)(error.message).toBe('custom error');
        (0, vitest_1.expect)(error.name).toBe('CorruptedDataError');
    });
    (0, vitest_1.it)('should create DecryptionError with default message', function () {
        var error = new types_js_1.DecryptionError();
        (0, vitest_1.expect)(error.message).toBe('Decryption failed');
        (0, vitest_1.expect)(error.name).toBe('DecryptionError');
    });
    (0, vitest_1.it)('should create DecryptionError with custom message', function () {
        var error = new types_js_1.DecryptionError('custom decryption error');
        (0, vitest_1.expect)(error.message).toBe('custom decryption error');
        (0, vitest_1.expect)(error.name).toBe('DecryptionError');
    });
});
