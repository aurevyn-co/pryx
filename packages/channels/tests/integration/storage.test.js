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
var promises_1 = require("fs/promises");
var os_1 = require("os");
var path_1 = require("path");
var storage_js_1 = require("../../src/storage.js");
var registry_js_1 = require("../../src/registry.js");
var defaultSettings = {
    allowCommands: true,
    autoReply: false,
    filterPatterns: [],
    allowedUsers: [],
    blockedUsers: [],
};
(0, vitest_1.describe)('ChannelStorage', function () {
    var tempDir;
    var configPath;
    var storage;
    (0, vitest_1.beforeEach)(function () { return __awaiter(void 0, void 0, void 0, function () {
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, (0, promises_1.mkdtemp)((0, path_1.join)((0, os_1.tmpdir)(), 'channels-test-'))];
                case 1:
                    tempDir = _a.sent();
                    configPath = (0, path_1.join)(tempDir, 'channels.json');
                    storage = new storage_js_1.ChannelStorage();
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.afterEach)(function () { return __awaiter(void 0, void 0, void 0, function () {
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, (0, promises_1.rm)(tempDir, { recursive: true, force: true })];
                case 1:
                    _a.sent();
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.describe)('save', function () {
        (0, vitest_1.it)('should save registry to file', function () { return __awaiter(void 0, void 0, void 0, function () {
            var registry, exists;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        registry = (0, registry_js_1.createRegistry)();
                        registry.addChannel({
                            id: 'telegram-bot',
                            name: 'Telegram Bot',
                            type: 'telegram',
                            enabled: true,
                            config: { botToken: 'test-token' },
                            settings: defaultSettings,
                        });
                        return [4 /*yield*/, storage.save(configPath, registry)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, storage.exists(configPath)];
                    case 2:
                        exists = _a.sent();
                        (0, vitest_1.expect)(exists).toBe(true);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should create directory if not exists', function () { return __awaiter(void 0, void 0, void 0, function () {
            var nestedPath, registry, exists;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        nestedPath = (0, path_1.join)(tempDir, 'nested', 'deep', 'channels.json');
                        registry = (0, registry_js_1.createRegistry)();
                        return [4 /*yield*/, storage.save(nestedPath, registry)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, storage.exists(nestedPath)];
                    case 2:
                        exists = _a.sent();
                        (0, vitest_1.expect)(exists).toBe(true);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('load', function () {
        (0, vitest_1.it)('should load registry from file', function () { return __awaiter(void 0, void 0, void 0, function () {
            var registry, loaded;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        registry = (0, registry_js_1.createRegistry)();
                        registry.addChannel({
                            id: 'telegram-bot',
                            name: 'Telegram Bot',
                            type: 'telegram',
                            enabled: true,
                            config: { botToken: 'test-token' },
                            settings: defaultSettings,
                        });
                        return [4 /*yield*/, storage.save(configPath, registry)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, storage.load(configPath)];
                    case 2:
                        loaded = _a.sent();
                        (0, vitest_1.expect)(loaded.hasChannel('telegram-bot')).toBe(true);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should return new registry when file not found', function () { return __awaiter(void 0, void 0, void 0, function () {
            var nonExistentPath, registry;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        nonExistentPath = (0, path_1.join)(tempDir, 'nonexistent.json');
                        return [4 /*yield*/, storage.load(nonExistentPath)];
                    case 1:
                        registry = _a.sent();
                        (0, vitest_1.expect)(registry).toBeInstanceOf(registry_js_1.ChannelRegistry);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('exists', function () {
        (0, vitest_1.it)('should return true when file exists', function () { return __awaiter(void 0, void 0, void 0, function () {
            var exists;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, storage.save(configPath, (0, registry_js_1.createRegistry)())];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, storage.exists(configPath)];
                    case 2:
                        exists = _a.sent();
                        (0, vitest_1.expect)(exists).toBe(true);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should return false when file does not exist', function () { return __awaiter(void 0, void 0, void 0, function () {
            var exists;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, storage.exists((0, path_1.join)(tempDir, 'nonexistent.json'))];
                    case 1:
                        exists = _a.sent();
                        (0, vitest_1.expect)(exists).toBe(false);
                        return [2 /*return*/];
                }
            });
        }); });
    });
});
(0, vitest_1.describe)('createStorage', function () {
    (0, vitest_1.it)('should create new storage instance', function () {
        var storage = (0, storage_js_1.createStorage)();
        (0, vitest_1.expect)(storage).toBeInstanceOf(storage_js_1.ChannelStorage);
    });
});
