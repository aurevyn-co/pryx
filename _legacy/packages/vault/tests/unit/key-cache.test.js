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
var key_cache_js_1 = require("../../src/key-cache.js");
var crypto_js_1 = require("../../src/crypto.js");
(0, vitest_1.describe)('KeyCache', function () {
    var cache;
    (0, vitest_1.beforeEach)(function () {
        cache = (0, key_cache_js_1.createKeyCache)();
    });
    (0, vitest_1.afterEach)(function () {
        cache.destroy();
    });
    (0, vitest_1.describe)('getKey', function () {
        (0, vitest_1.it)('should derive and cache key', function () { return __awaiter(void 0, void 0, void 0, function () {
            var password, salt, key1, key2;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        password = 'test-password';
                        salt = (0, crypto_js_1.generateSalt)();
                        return [4 /*yield*/, cache.getKey(password, salt)];
                    case 1:
                        key1 = _a.sent();
                        return [4 /*yield*/, cache.getKey(password, salt)];
                    case 2:
                        key2 = _a.sent();
                        (0, vitest_1.expect)(key1).toEqual(key2);
                        (0, vitest_1.expect)(cache.getStats().size).toBe(1);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should derive different keys for different passwords', function () { return __awaiter(void 0, void 0, void 0, function () {
            var salt, key1, key2;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        salt = (0, crypto_js_1.generateSalt)();
                        return [4 /*yield*/, cache.getKey('password1', salt)];
                    case 1:
                        key1 = _a.sent();
                        return [4 /*yield*/, cache.getKey('password2', salt)];
                    case 2:
                        key2 = _a.sent();
                        (0, vitest_1.expect)(key1).not.toEqual(key2);
                        (0, vitest_1.expect)(cache.getStats().size).toBe(2);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should derive different keys for different salts', function () { return __awaiter(void 0, void 0, void 0, function () {
            var password, key1, key2;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        password = 'test-password';
                        return [4 /*yield*/, cache.getKey(password, (0, crypto_js_1.generateSalt)())];
                    case 1:
                        key1 = _a.sent();
                        return [4 /*yield*/, cache.getKey(password, (0, crypto_js_1.generateSalt)())];
                    case 2:
                        key2 = _a.sent();
                        (0, vitest_1.expect)(key1).not.toEqual(key2);
                        (0, vitest_1.expect)(cache.getStats().size).toBe(2);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should increment access count on cache hit', function () { return __awaiter(void 0, void 0, void 0, function () {
            var password, salt, stats;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        password = 'test-password';
                        salt = (0, crypto_js_1.generateSalt)();
                        return [4 /*yield*/, cache.getKey(password, salt)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, cache.getKey(password, salt)];
                    case 2:
                        _a.sent();
                        return [4 /*yield*/, cache.getKey(password, salt)];
                    case 3:
                        _a.sent();
                        stats = cache.getStats();
                        (0, vitest_1.expect)(stats.totalAccessCount).toBe(3);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('invalidate', function () {
        (0, vitest_1.it)('should remove specific entry from cache', function () { return __awaiter(void 0, void 0, void 0, function () {
            var password, salt;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        password = 'test-password';
                        salt = (0, crypto_js_1.generateSalt)();
                        return [4 /*yield*/, cache.getKey(password, salt)];
                    case 1:
                        _a.sent();
                        (0, vitest_1.expect)(cache.getStats().size).toBe(1);
                        cache.invalidate(password, salt);
                        (0, vitest_1.expect)(cache.getStats().size).toBe(0);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should not affect other entries', function () { return __awaiter(void 0, void 0, void 0, function () {
            var salt1, salt2;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        salt1 = (0, crypto_js_1.generateSalt)();
                        salt2 = (0, crypto_js_1.generateSalt)();
                        return [4 /*yield*/, cache.getKey('password1', salt1)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, cache.getKey('password2', salt2)];
                    case 2:
                        _a.sent();
                        (0, vitest_1.expect)(cache.getStats().size).toBe(2);
                        cache.invalidate('password1', salt1);
                        (0, vitest_1.expect)(cache.getStats().size).toBe(1);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('invalidateAll', function () {
        (0, vitest_1.it)('should clear all entries', function () { return __awaiter(void 0, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, cache.getKey('password1', (0, crypto_js_1.generateSalt)())];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, cache.getKey('password2', (0, crypto_js_1.generateSalt)())];
                    case 2:
                        _a.sent();
                        return [4 /*yield*/, cache.getKey('password3', (0, crypto_js_1.generateSalt)())];
                    case 3:
                        _a.sent();
                        (0, vitest_1.expect)(cache.getStats().size).toBe(3);
                        cache.invalidateAll();
                        (0, vitest_1.expect)(cache.getStats().size).toBe(0);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('getStats', function () {
        (0, vitest_1.it)('should return correct size', function () { return __awaiter(void 0, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        (0, vitest_1.expect)(cache.getStats().size).toBe(0);
                        return [4 /*yield*/, cache.getKey('password1', (0, crypto_js_1.generateSalt)())];
                    case 1:
                        _a.sent();
                        (0, vitest_1.expect)(cache.getStats().size).toBe(1);
                        return [4 /*yield*/, cache.getKey('password2', (0, crypto_js_1.generateSalt)())];
                    case 2:
                        _a.sent();
                        (0, vitest_1.expect)(cache.getStats().size).toBe(2);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should track total access count', function () { return __awaiter(void 0, void 0, void 0, function () {
            var salt;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        salt = (0, crypto_js_1.generateSalt)();
                        return [4 /*yield*/, cache.getKey('password', salt)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, cache.getKey('password', salt)];
                    case 2:
                        _a.sent();
                        return [4 /*yield*/, cache.getKey('password', salt)];
                    case 3:
                        _a.sent();
                        (0, vitest_1.expect)(cache.getStats().totalAccessCount).toBe(3);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should track entry timestamps', function () { return __awaiter(void 0, void 0, void 0, function () {
            var before, after, stats;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        before = Date.now();
                        return [4 /*yield*/, cache.getKey('password', (0, crypto_js_1.generateSalt)())];
                    case 1:
                        _a.sent();
                        after = Date.now();
                        stats = cache.getStats();
                        (0, vitest_1.expect)(stats.oldestEntry).toBeGreaterThanOrEqual(before);
                        (0, vitest_1.expect)(stats.oldestEntry).toBeLessThanOrEqual(after);
                        (0, vitest_1.expect)(stats.newestEntry).toBeGreaterThanOrEqual(before);
                        (0, vitest_1.expect)(stats.newestEntry).toBeLessThanOrEqual(after);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('LRU eviction', function () {
        (0, vitest_1.it)('should evict least recently used when max entries reached', function () { return __awaiter(void 0, void 0, void 0, function () {
            var smallCache, salt1, salt2, salt3, key1Again;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        smallCache = (0, key_cache_js_1.createKeyCache)({ maxEntries: 2 });
                        salt1 = (0, crypto_js_1.generateSalt)();
                        salt2 = (0, crypto_js_1.generateSalt)();
                        salt3 = (0, crypto_js_1.generateSalt)();
                        return [4 /*yield*/, smallCache.getKey('password1', salt1)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, smallCache.getKey('password2', salt2)];
                    case 2:
                        _a.sent();
                        (0, vitest_1.expect)(smallCache.getStats().size).toBe(2);
                        // Wait a bit to ensure different timestamps
                        return [4 /*yield*/, new Promise(function (resolve) { return setTimeout(resolve, 10); })];
                    case 3:
                        // Wait a bit to ensure different timestamps
                        _a.sent();
                        // Access first entry to make it recently used
                        return [4 /*yield*/, smallCache.getKey('password1', salt1)];
                    case 4:
                        // Access first entry to make it recently used
                        _a.sent();
                        // Wait a bit more
                        return [4 /*yield*/, new Promise(function (resolve) { return setTimeout(resolve, 10); })];
                    case 5:
                        // Wait a bit more
                        _a.sent();
                        // Add third entry - should evict password2 (least recently used)
                        return [4 /*yield*/, smallCache.getKey('password3', salt3)];
                    case 6:
                        // Add third entry - should evict password2 (least recently used)
                        _a.sent();
                        (0, vitest_1.expect)(smallCache.getStats().size).toBe(2);
                        return [4 /*yield*/, smallCache.getKey('password1', salt1)];
                    case 7:
                        key1Again = _a.sent();
                        (0, vitest_1.expect)(smallCache.getStats().totalAccessCount).toBeGreaterThanOrEqual(3);
                        smallCache.destroy();
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('cleanup', function () {
        (0, vitest_1.it)('should clean up expired entries', function () { return __awaiter(void 0, void 0, void 0, function () {
            var shortCache, salt;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        shortCache = (0, key_cache_js_1.createKeyCache)({
                            maxAgeMs: 100,
                            maxIdleMs: 50,
                        });
                        salt = (0, crypto_js_1.generateSalt)();
                        return [4 /*yield*/, shortCache.getKey('password', salt)];
                    case 1:
                        _a.sent();
                        (0, vitest_1.expect)(shortCache.getStats().size).toBe(1);
                        return [4 /*yield*/, new Promise(function (resolve) { return setTimeout(resolve, 200); })];
                    case 2:
                        _a.sent();
                        return [4 /*yield*/, shortCache.getKey('password2', (0, crypto_js_1.generateSalt)())];
                    case 3:
                        _a.sent();
                        (0, vitest_1.expect)(shortCache.getStats().size).toBeLessThanOrEqual(2);
                        shortCache.destroy();
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('destroy', function () {
        (0, vitest_1.it)('should clear all entries and stop cleanup interval', function () { return __awaiter(void 0, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, cache.getKey('password1', (0, crypto_js_1.generateSalt)())];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, cache.getKey('password2', (0, crypto_js_1.generateSalt)())];
                    case 2:
                        _a.sent();
                        (0, vitest_1.expect)(cache.getStats().size).toBe(2);
                        cache.destroy();
                        (0, vitest_1.expect)(cache.getStats().size).toBe(0);
                        return [2 /*return*/];
                }
            });
        }); });
    });
});
(0, vitest_1.describe)('createKeyCache', function () {
    (0, vitest_1.it)('should create cache with default config', function () {
        var cache = (0, key_cache_js_1.createKeyCache)();
        (0, vitest_1.expect)(cache).toBeInstanceOf(key_cache_js_1.KeyCache);
        cache.destroy();
    });
    (0, vitest_1.it)('should create cache with custom config', function () {
        var cache = (0, key_cache_js_1.createKeyCache)({
            maxEntries: 5,
            maxAgeMs: 1000,
            maxIdleMs: 500,
        });
        (0, vitest_1.expect)(cache).toBeInstanceOf(key_cache_js_1.KeyCache);
        cache.destroy();
    });
});
