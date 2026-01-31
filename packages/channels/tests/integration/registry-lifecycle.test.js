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
var registry_js_1 = require("../../src/registry.js");
var storage_js_1 = require("../../src/storage.js");
var defaultSettings = {
    allowCommands: true,
    autoReply: false,
    filterPatterns: [],
    allowedUsers: [],
    blockedUsers: [],
};
(0, vitest_1.describe)('Channel Registry Lifecycle', function () {
    var registry;
    var storage;
    var tempDir;
    var configPath;
    (0, vitest_1.beforeEach)(function () { return __awaiter(void 0, void 0, void 0, function () {
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    registry = (0, registry_js_1.createRegistry)();
                    storage = (0, storage_js_1.createStorage)();
                    return [4 /*yield*/, (0, promises_1.mkdtemp)((0, path_1.join)((0, os_1.tmpdir)(), 'channels-lifecycle-test-'))];
                case 1:
                    tempDir = _a.sent();
                    configPath = (0, path_1.join)(tempDir, 'channels.json');
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
    (0, vitest_1.it)('should complete full lifecycle: init → add → save → load → update → remove', function () { return __awaiter(void 0, void 0, void 0, function () {
        var channel, _a, loaded;
        var _b;
        return __generator(this, function (_c) {
            switch (_c.label) {
                case 0:
                    channel = {
                        id: 'telegram-bot',
                        name: 'Telegram Bot',
                        type: 'telegram',
                        enabled: true,
                        config: { botToken: 'test-token' },
                        settings: defaultSettings,
                    };
                    registry.addChannel(channel);
                    (0, vitest_1.expect)(registry.hasChannel('telegram-bot')).toBe(true);
                    return [4 /*yield*/, storage.save(configPath, registry)];
                case 1:
                    _c.sent();
                    _a = vitest_1.expect;
                    return [4 /*yield*/, storage.exists(configPath)];
                case 2:
                    _a.apply(void 0, [_c.sent()]).toBe(true);
                    return [4 /*yield*/, storage.load(configPath)];
                case 3:
                    loaded = _c.sent();
                    (0, vitest_1.expect)(loaded.hasChannel('telegram-bot')).toBe(true);
                    loaded.updateChannel('telegram-bot', { name: 'Updated Bot' });
                    (0, vitest_1.expect)((_b = loaded.getChannel('telegram-bot')) === null || _b === void 0 ? void 0 : _b.name).toBe('Updated Bot');
                    loaded.removeChannel('telegram-bot');
                    (0, vitest_1.expect)(loaded.hasChannel('telegram-bot')).toBe(false);
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should handle multiple channels', function () { return __awaiter(void 0, void 0, void 0, function () {
        var channels, _i, channels_1, channel, loaded;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    channels = [
                        {
                            id: 'telegram-1',
                            name: 'Telegram 1',
                            type: 'telegram',
                            enabled: true,
                            config: { botToken: 'token1' },
                            settings: defaultSettings,
                        },
                        {
                            id: 'discord-1',
                            name: 'Discord 1',
                            type: 'discord',
                            enabled: true,
                            config: { botToken: 'token2', applicationId: 'app-id' },
                            settings: defaultSettings,
                        },
                    ];
                    for (_i = 0, channels_1 = channels; _i < channels_1.length; _i++) {
                        channel = channels_1[_i];
                        registry.addChannel(channel);
                    }
                    (0, vitest_1.expect)(registry.size).toBe(2);
                    return [4 /*yield*/, storage.save(configPath, registry)];
                case 1:
                    _a.sent();
                    return [4 /*yield*/, storage.load(configPath)];
                case 2:
                    loaded = _a.sent();
                    (0, vitest_1.expect)(loaded.hasChannel('telegram-1')).toBe(true);
                    (0, vitest_1.expect)(loaded.hasChannel('discord-1')).toBe(true);
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should handle enable/disable channels', function () { return __awaiter(void 0, void 0, void 0, function () {
        var loaded;
        var _a, _b, _c;
        return __generator(this, function (_d) {
            switch (_d.label) {
                case 0:
                    registry.addChannel({
                        id: 'test',
                        name: 'Test',
                        type: 'telegram',
                        enabled: true,
                        config: { botToken: 'token' },
                        settings: defaultSettings,
                    });
                    registry.disableChannel('test');
                    (0, vitest_1.expect)((_a = registry.getChannel('test')) === null || _a === void 0 ? void 0 : _a.enabled).toBe(false);
                    return [4 /*yield*/, storage.save(configPath, registry)];
                case 1:
                    _d.sent();
                    return [4 /*yield*/, storage.load(configPath)];
                case 2:
                    loaded = _d.sent();
                    (0, vitest_1.expect)((_b = loaded.getChannel('test')) === null || _b === void 0 ? void 0 : _b.enabled).toBe(false);
                    loaded.enableChannel('test');
                    (0, vitest_1.expect)((_c = loaded.getChannel('test')) === null || _c === void 0 ? void 0 : _c.enabled).toBe(true);
                    return [2 /*return*/];
            }
        });
    }); });
});
