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
var index_js_1 = require("../../src/index.js");
var defaultSettings = {
    allowCommands: true,
    autoReply: false,
    filterPatterns: [],
    allowedUsers: [],
    blockedUsers: [],
};
(0, vitest_1.describe)('Channel Workflow E2E', function () {
    var tempDir;
    var configPath;
    (0, vitest_1.beforeEach)(function () { return __awaiter(void 0, void 0, void 0, function () {
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, (0, promises_1.mkdtemp)((0, path_1.join)((0, os_1.tmpdir)(), 'channels-e2e-test-'))];
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
    (0, vitest_1.it)('should configure Telegram bot', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry, storage, channel, _a;
        return __generator(this, function (_b) {
            switch (_b.label) {
                case 0:
                    registry = (0, index_js_1.createRegistry)();
                    storage = (0, index_js_1.createStorage)();
                    registry.addChannel({
                        id: 'telegram-bot',
                        name: 'My Telegram Bot',
                        type: 'telegram',
                        enabled: true,
                        config: {
                            botToken: '123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11',
                            chatId: '123456789',
                            parseMode: 'Markdown',
                        },
                        settings: defaultSettings,
                    });
                    channel = registry.getChannel('telegram-bot');
                    (0, vitest_1.expect)(channel === null || channel === void 0 ? void 0 : channel.type).toBe('telegram');
                    (0, vitest_1.expect)(channel === null || channel === void 0 ? void 0 : channel.config.botToken).toBe('123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11');
                    return [4 /*yield*/, storage.save(configPath, registry)];
                case 1:
                    _b.sent();
                    _a = vitest_1.expect;
                    return [4 /*yield*/, storage.exists(configPath)];
                case 2:
                    _a.apply(void 0, [_b.sent()]).toBe(true);
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should configure Discord bot', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry;
        return __generator(this, function (_a) {
            registry = (0, index_js_1.createRegistry)();
            registry.addChannel({
                id: 'discord-bot',
                name: 'Discord Bot',
                type: 'discord',
                enabled: true,
                config: {
                    botToken: 'discord-token',
                    applicationId: '123456789',
                    guildId: '987654321',
                    intents: ['GUILDS', 'GUILD_MESSAGES'],
                },
                settings: defaultSettings,
            });
            (0, vitest_1.expect)(registry.hasChannel('discord-bot')).toBe(true);
            return [2 /*return*/];
        });
    }); });
    (0, vitest_1.it)('should configure webhook endpoint', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry;
        return __generator(this, function (_a) {
            registry = (0, index_js_1.createRegistry)();
            registry.addChannel({
                id: 'webhook-endpoint',
                name: 'Webhook Endpoint',
                type: 'webhook',
                enabled: true,
                config: {
                    url: 'https://api.example.com/webhook',
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'X-Custom-Header': 'value',
                    },
                },
                settings: defaultSettings,
            });
            (0, vitest_1.expect)(registry.hasChannel('webhook-endpoint')).toBe(true);
            return [2 /*return*/];
        });
    }); });
    (0, vitest_1.it)('should configure email channel', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry;
        return __generator(this, function (_a) {
            registry = (0, index_js_1.createRegistry)();
            registry.addChannel({
                id: 'email-channel',
                name: 'Email Channel',
                type: 'email',
                enabled: true,
                config: {
                    imap: {
                        host: 'imap.gmail.com',
                        port: 993,
                        secure: true,
                        username: 'user@gmail.com',
                        password: 'app-password',
                    },
                    smtp: {
                        host: 'smtp.gmail.com',
                        port: 587,
                        secure: true,
                        username: 'user@gmail.com',
                        password: 'app-password',
                    },
                    checkInterval: 60000,
                    markAsRead: true,
                },
                settings: defaultSettings,
            });
            (0, vitest_1.expect)(registry.hasChannel('email-channel')).toBe(true);
            return [2 /*return*/];
        });
    }); });
    (0, vitest_1.it)('should enable and disable channels', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry, enabledChannels;
        var _a, _b;
        return __generator(this, function (_c) {
            registry = (0, index_js_1.createRegistry)();
            registry.addChannel({
                id: 'test-channel',
                name: 'Test Channel',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token' },
                settings: defaultSettings,
            });
            registry.disableChannel('test-channel');
            (0, vitest_1.expect)((_a = registry.getChannel('test-channel')) === null || _a === void 0 ? void 0 : _a.enabled).toBe(false);
            enabledChannels = registry.getEnabledChannels();
            (0, vitest_1.expect)(enabledChannels.some(function (c) { return c.id === 'test-channel'; })).toBe(false);
            registry.enableChannel('test-channel');
            (0, vitest_1.expect)((_b = registry.getChannel('test-channel')) === null || _b === void 0 ? void 0 : _b.enabled).toBe(true);
            return [2 /*return*/];
        });
    }); });
    (0, vitest_1.it)('should persist configuration across restarts', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry1, storage, registry2;
        var _a;
        return __generator(this, function (_b) {
            switch (_b.label) {
                case 0:
                    registry1 = (0, index_js_1.createRegistry)();
                    storage = (0, index_js_1.createStorage)();
                    registry1.addChannel({
                        id: 'persistent-channel',
                        name: 'Persistent Channel',
                        type: 'telegram',
                        enabled: true,
                        config: { botToken: 'persistent-token' },
                        settings: defaultSettings,
                    });
                    return [4 /*yield*/, storage.save(configPath, registry1)];
                case 1:
                    _b.sent();
                    return [4 /*yield*/, storage.load(configPath)];
                case 2:
                    registry2 = _b.sent();
                    (0, vitest_1.expect)(registry2.hasChannel('persistent-channel')).toBe(true);
                    (0, vitest_1.expect)((_a = registry2.getChannel('persistent-channel')) === null || _a === void 0 ? void 0 : _a.config.botToken).toBe('persistent-token');
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should filter channels by type', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry, telegramChannels;
        return __generator(this, function (_a) {
            registry = (0, index_js_1.createRegistry)();
            registry.addChannel({
                id: 'telegram-1',
                name: 'Telegram 1',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token1' },
                settings: defaultSettings,
            });
            registry.addChannel({
                id: 'discord-1',
                name: 'Discord 1',
                type: 'discord',
                enabled: true,
                config: { botToken: 'token2', applicationId: 'app' },
                settings: defaultSettings,
            });
            registry.addChannel({
                id: 'telegram-2',
                name: 'Telegram 2',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token3' },
                settings: defaultSettings,
            });
            telegramChannels = registry.getChannelsByType('telegram');
            (0, vitest_1.expect)(telegramChannels.length).toBe(2);
            return [2 /*return*/];
        });
    }); });
    (0, vitest_1.it)('should handle bulk operations', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry;
        return __generator(this, function (_a) {
            registry = (0, index_js_1.createRegistry)();
            registry.addChannel({
                id: 'channel1',
                name: 'Channel 1',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token1' },
                settings: defaultSettings,
            });
            registry.addChannel({
                id: 'channel2',
                name: 'Channel 2',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token2' },
                settings: defaultSettings,
            });
            registry.disableAll();
            (0, vitest_1.expect)(registry.getAllChannels().every(function (c) { return !c.enabled; })).toBe(true);
            registry.enableAll();
            (0, vitest_1.expect)(registry.getAllChannels().every(function (c) { return c.enabled; })).toBe(true);
            return [2 /*return*/];
        });
    }); });
    (0, vitest_1.it)('should update channel status', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry, status;
        var _a;
        return __generator(this, function (_b) {
            registry = (0, index_js_1.createRegistry)();
            registry.addChannel({
                id: 'test',
                name: 'Test',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token' },
                settings: defaultSettings,
            });
            registry.updateChannelStatus('test', {
                connected: true,
                messageCount: 42,
            });
            status = (_a = registry.getChannel('test')) === null || _a === void 0 ? void 0 : _a.status;
            (0, vitest_1.expect)(status === null || status === void 0 ? void 0 : status.connected).toBe(true);
            (0, vitest_1.expect)(status === null || status === void 0 ? void 0 : status.messageCount).toBe(42);
            return [2 /*return*/];
        });
    }); });
});
