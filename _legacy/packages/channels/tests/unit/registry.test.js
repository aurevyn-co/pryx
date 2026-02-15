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
var registry_js_1 = require("../../src/registry.js");
var types_js_1 = require("../../src/types.js");
(0, vitest_1.describe)('ChannelRegistry', function () {
    var registry;
    (0, vitest_1.beforeEach)(function () {
        registry = new registry_js_1.ChannelRegistry();
    });
    (0, vitest_1.describe)('constructor', function () {
        (0, vitest_1.it)('should create empty registry', function () {
            (0, vitest_1.expect)(registry.size).toBe(0);
        });
    });
    (0, vitest_1.describe)('addChannel', function () {
        (0, vitest_1.it)('should add telegram channel', function () {
            var _a;
            var channel = {
                id: 'telegram-bot',
                name: 'My Telegram Bot',
                type: 'telegram',
                enabled: true,
                config: {
                    botToken: 'test-token',
                    chatId: '123456',
                },
            };
            registry.addChannel(channel);
            (0, vitest_1.expect)(registry.hasChannel('telegram-bot')).toBe(true);
            (0, vitest_1.expect)((_a = registry.getChannel('telegram-bot')) === null || _a === void 0 ? void 0 : _a.name).toBe('My Telegram Bot');
        });
        (0, vitest_1.it)('should add discord channel', function () {
            var channel = {
                id: 'discord-bot',
                name: 'Discord Bot',
                type: 'discord',
                enabled: true,
                config: {
                    botToken: 'discord-token',
                    applicationId: 'app-id',
                    intents: ['GUILDS', 'GUILD_MESSAGES'],
                },
            };
            registry.addChannel(channel);
            (0, vitest_1.expect)(registry.hasChannel('discord-bot')).toBe(true);
        });
        (0, vitest_1.it)('should add webhook channel', function () {
            var channel = {
                id: 'webhook-endpoint',
                name: 'Webhook Endpoint',
                type: 'webhook',
                enabled: true,
                config: {
                    url: 'https://api.example.com/webhook',
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                },
            };
            registry.addChannel(channel);
            (0, vitest_1.expect)(registry.hasChannel('webhook-endpoint')).toBe(true);
        });
        (0, vitest_1.it)('should throw when channel already exists', function () {
            var channel = {
                id: 'test',
                name: 'Test',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token' },
            };
            registry.addChannel(channel);
            (0, vitest_1.expect)(function () { return registry.addChannel(channel); }).toThrow(types_js_1.ChannelAlreadyExistsError);
        });
        (0, vitest_1.it)('should throw on invalid config', function () {
            var channel = {
                id: 'invalid',
                name: '',
                type: 'telegram',
                enabled: true,
                config: {},
            };
            (0, vitest_1.expect)(function () { return registry.addChannel(channel); }).toThrow(types_js_1.ChannelValidationError);
        });
    });
    (0, vitest_1.describe)('updateChannel', function () {
        (0, vitest_1.it)('should update existing channel', function () {
            var _a;
            registry.addChannel({
                id: 'test',
                name: 'Test',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token' },
            });
            var updated = registry.updateChannel('test', { name: 'Updated' });
            (0, vitest_1.expect)(updated.name).toBe('Updated');
            (0, vitest_1.expect)((_a = registry.getChannel('test')) === null || _a === void 0 ? void 0 : _a.name).toBe('Updated');
        });
        (0, vitest_1.it)('should throw when channel not found', function () {
            (0, vitest_1.expect)(function () { return registry.updateChannel('nonexistent', { name: 'Test' }); }).toThrow(types_js_1.ChannelNotFoundError);
        });
    });
    (0, vitest_1.describe)('removeChannel', function () {
        (0, vitest_1.it)('should remove channel', function () {
            registry.addChannel({
                id: 'test',
                name: 'Test',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token' },
            });
            registry.removeChannel('test');
            (0, vitest_1.expect)(registry.hasChannel('test')).toBe(false);
        });
        (0, vitest_1.it)('should throw when channel not found', function () {
            (0, vitest_1.expect)(function () { return registry.removeChannel('nonexistent'); }).toThrow(types_js_1.ChannelNotFoundError);
        });
    });
    (0, vitest_1.describe)('getChannel', function () {
        (0, vitest_1.it)('should return channel by id', function () {
            registry.addChannel({
                id: 'test',
                name: 'Test',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token' },
            });
            var channel = registry.getChannel('test');
            (0, vitest_1.expect)(channel).toBeDefined();
            (0, vitest_1.expect)(channel === null || channel === void 0 ? void 0 : channel.id).toBe('test');
        });
        (0, vitest_1.it)('should return undefined for nonexistent channel', function () {
            var channel = registry.getChannel('nonexistent');
            (0, vitest_1.expect)(channel).toBeUndefined();
        });
    });
    (0, vitest_1.describe)('getAllChannels', function () {
        (0, vitest_1.it)('should return all channels', function () {
            registry.addChannel({
                id: 'channel1',
                name: 'Channel 1',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token1' },
            });
            registry.addChannel({
                id: 'channel2',
                name: 'Channel 2',
                type: 'discord',
                enabled: true,
                config: { botToken: 'token2', applicationId: 'app' },
            });
            var channels = registry.getAllChannels();
            (0, vitest_1.expect)(channels.length).toBe(2);
        });
    });
    (0, vitest_1.describe)('getEnabledChannels', function () {
        (0, vitest_1.it)('should return only enabled channels', function () {
            registry.addChannel({
                id: 'enabled',
                name: 'Enabled',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token' },
            });
            registry.addChannel({
                id: 'disabled',
                name: 'Disabled',
                type: 'telegram',
                enabled: false,
                config: { botToken: 'token' },
            });
            var enabled = registry.getEnabledChannels();
            (0, vitest_1.expect)(enabled.length).toBe(1);
            (0, vitest_1.expect)(enabled[0].id).toBe('enabled');
        });
    });
    (0, vitest_1.describe)('getChannelsByType', function () {
        (0, vitest_1.it)('should return channels by type', function () {
            registry.addChannel({
                id: 'telegram1',
                name: 'Telegram 1',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token' },
            });
            registry.addChannel({
                id: 'discord1',
                name: 'Discord 1',
                type: 'discord',
                enabled: true,
                config: { botToken: 'token', applicationId: 'app' },
            });
            registry.addChannel({
                id: 'telegram2',
                name: 'Telegram 2',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token' },
            });
            var telegramChannels = registry.getChannelsByType('telegram');
            (0, vitest_1.expect)(telegramChannels.length).toBe(2);
            (0, vitest_1.expect)(telegramChannels.every(function (c) { return c.type === 'telegram'; })).toBe(true);
        });
    });
    (0, vitest_1.describe)('enableChannel / disableChannel', function () {
        (0, vitest_1.it)('should disable channel', function () {
            var _a;
            registry.addChannel({
                id: 'test',
                name: 'Test',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token' },
            });
            registry.disableChannel('test');
            (0, vitest_1.expect)((_a = registry.getChannel('test')) === null || _a === void 0 ? void 0 : _a.enabled).toBe(false);
        });
        (0, vitest_1.it)('should enable channel', function () {
            var _a;
            registry.addChannel({
                id: 'test',
                name: 'Test',
                type: 'telegram',
                enabled: false,
                config: { botToken: 'token' },
            });
            registry.enableChannel('test');
            (0, vitest_1.expect)((_a = registry.getChannel('test')) === null || _a === void 0 ? void 0 : _a.enabled).toBe(true);
        });
    });
    (0, vitest_1.describe)('enableAll / disableAll', function () {
        (0, vitest_1.it)('should enable all channels', function () {
            registry.addChannel({
                id: 'test1',
                name: 'Test 1',
                type: 'telegram',
                enabled: false,
                config: { botToken: 'token' },
            });
            registry.addChannel({
                id: 'test2',
                name: 'Test 2',
                type: 'telegram',
                enabled: false,
                config: { botToken: 'token' },
            });
            registry.enableAll();
            (0, vitest_1.expect)(registry.getAllChannels().every(function (c) { return c.enabled; })).toBe(true);
        });
        (0, vitest_1.it)('should disable all channels', function () {
            registry.addChannel({
                id: 'test1',
                name: 'Test 1',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token' },
            });
            registry.addChannel({
                id: 'test2',
                name: 'Test 2',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token' },
            });
            registry.disableAll();
            (0, vitest_1.expect)(registry.getAllChannels().every(function (c) { return !c.enabled; })).toBe(true);
        });
    });
    (0, vitest_1.describe)('enableType / disableType', function () {
        (0, vitest_1.it)('should enable channels by type', function () {
            var _a, _b;
            registry.addChannel({
                id: 'telegram1',
                name: 'Telegram 1',
                type: 'telegram',
                enabled: false,
                config: { botToken: 'token' },
            });
            registry.addChannel({
                id: 'discord1',
                name: 'Discord 1',
                type: 'discord',
                enabled: false,
                config: { botToken: 'token', applicationId: 'app' },
            });
            registry.enableType('telegram');
            (0, vitest_1.expect)((_a = registry.getChannel('telegram1')) === null || _a === void 0 ? void 0 : _a.enabled).toBe(true);
            (0, vitest_1.expect)((_b = registry.getChannel('discord1')) === null || _b === void 0 ? void 0 : _b.enabled).toBe(false);
        });
    });
    (0, vitest_1.describe)('updateChannelStatus', function () {
        (0, vitest_1.it)('should update channel status', function () {
            var _a;
            registry.addChannel({
                id: 'test',
                name: 'Test',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token' },
            });
            registry.updateChannelStatus('test', {
                connected: true,
                messageCount: 10,
            });
            var status = (_a = registry.getChannel('test')) === null || _a === void 0 ? void 0 : _a.status;
            (0, vitest_1.expect)(status === null || status === void 0 ? void 0 : status.connected).toBe(true);
            (0, vitest_1.expect)(status === null || status === void 0 ? void 0 : status.messageCount).toBe(10);
        });
        (0, vitest_1.it)('should throw when channel not found', function () {
            (0, vitest_1.expect)(function () { return registry.updateChannelStatus('nonexistent', { connected: true }); }).toThrow(types_js_1.ChannelNotFoundError);
        });
    });
    (0, vitest_1.describe)('validateChannel', function () {
        (0, vitest_1.it)('should validate existing channel', function () {
            registry.addChannel({
                id: 'test',
                name: 'Test',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token' },
            });
            var result = registry.validateChannel('test');
            (0, vitest_1.expect)(result.valid).toBe(true);
        });
        (0, vitest_1.it)('should return error for nonexistent channel', function () {
            var result = registry.validateChannel('nonexistent');
            (0, vitest_1.expect)(result.valid).toBe(false);
            (0, vitest_1.expect)(result.errors[0]).toContain('not found');
        });
    });
    (0, vitest_1.describe)('testConnection', function () {
        (0, vitest_1.it)('should fail for nonexistent channel', function () { return __awaiter(void 0, void 0, void 0, function () {
            var result;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, registry.testConnection('nonexistent')];
                    case 1:
                        result = _a.sent();
                        (0, vitest_1.expect)(result.success).toBe(false);
                        (0, vitest_1.expect)(result.error).toContain('not found');
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should return success for valid channel', function () { return __awaiter(void 0, void 0, void 0, function () {
            var result;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        registry.addChannel({
                            id: 'test',
                            name: 'Test',
                            type: 'telegram',
                            enabled: true,
                            config: { botToken: 'token' },
                        });
                        return [4 /*yield*/, registry.testConnection('test')];
                    case 1:
                        result = _a.sent();
                        (0, vitest_1.expect)(result.success).toBe(true);
                        (0, vitest_1.expect)(result.latency).toBeDefined();
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('toJSON / fromJSON', function () {
        (0, vitest_1.it)('should serialize to JSON', function () {
            registry.addChannel({
                id: 'test',
                name: 'Test',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token' },
            });
            var json = registry.toJSON();
            (0, vitest_1.expect)(json.version).toBe(1);
            (0, vitest_1.expect)(json.channels.length).toBe(1);
        });
        (0, vitest_1.it)('should deserialize from JSON', function () {
            var json = {
                version: 1,
                channels: [{
                        id: 'test',
                        name: 'Test',
                        type: 'telegram',
                        enabled: true,
                        config: { botToken: 'token' },
                    }],
            };
            registry.fromJSON(json);
            (0, vitest_1.expect)(registry.size).toBe(1);
            (0, vitest_1.expect)(registry.hasChannel('test')).toBe(true);
        });
        (0, vitest_1.it)('should throw on unsupported version', function () {
            var json = { version: 999, channels: [] };
            (0, vitest_1.expect)(function () { return registry.fromJSON(json); }).toThrow(types_js_1.ChannelValidationError);
        });
    });
    (0, vitest_1.describe)('clear', function () {
        (0, vitest_1.it)('should clear all channels', function () {
            registry.addChannel({
                id: 'test',
                name: 'Test',
                type: 'telegram',
                enabled: true,
                config: { botToken: 'token' },
            });
            registry.clear();
            (0, vitest_1.expect)(registry.size).toBe(0);
        });
    });
});
(0, vitest_1.describe)('createRegistry', function () {
    (0, vitest_1.it)('should create new registry', function () {
        var registry = (0, registry_js_1.createRegistry)();
        (0, vitest_1.expect)(registry).toBeInstanceOf(registry_js_1.ChannelRegistry);
        (0, vitest_1.expect)(registry.size).toBe(0);
    });
});
