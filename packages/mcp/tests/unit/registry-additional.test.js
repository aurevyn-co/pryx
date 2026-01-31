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
(0, vitest_1.describe)('MCPRegistry Additional Coverage', function () {
    var registry;
    (0, vitest_1.beforeEach)(function () {
        registry = new registry_js_1.MCPRegistry();
    });
    (0, vitest_1.describe)('enable/disable operations', function () {
        (0, vitest_1.it)('should enable server', function () {
            var _a;
            registry.addServer({
                id: 'test',
                name: 'Test',
                enabled: false,
                source: 'manual',
                transport: { type: 'stdio', command: 'test', args: [], env: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            });
            registry.enableServer('test');
            (0, vitest_1.expect)((_a = registry.getServer('test')) === null || _a === void 0 ? void 0 : _a.enabled).toBe(true);
        });
        (0, vitest_1.it)('should disable server', function () {
            var _a;
            registry.addServer({
                id: 'test',
                name: 'Test',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'test', args: [], env: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            });
            registry.disableServer('test');
            (0, vitest_1.expect)((_a = registry.getServer('test')) === null || _a === void 0 ? void 0 : _a.enabled).toBe(false);
        });
        (0, vitest_1.it)('should enable all servers', function () {
            var _a, _b;
            registry.addServer({
                id: 'server1',
                name: 'Server 1',
                enabled: false,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd1', args: [], env: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            });
            registry.addServer({
                id: 'server2',
                name: 'Server 2',
                enabled: false,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd2', args: [], env: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            });
            registry.enableAll();
            (0, vitest_1.expect)((_a = registry.getServer('server1')) === null || _a === void 0 ? void 0 : _a.enabled).toBe(true);
            (0, vitest_1.expect)((_b = registry.getServer('server2')) === null || _b === void 0 ? void 0 : _b.enabled).toBe(true);
        });
        (0, vitest_1.it)('should disable all servers', function () {
            var _a, _b;
            registry.addServer({
                id: 'server1',
                name: 'Server 1',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd1', args: [], env: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            });
            registry.addServer({
                id: 'server2',
                name: 'Server 2',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd2', args: [], env: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            });
            registry.disableAll();
            (0, vitest_1.expect)((_a = registry.getServer('server1')) === null || _a === void 0 ? void 0 : _a.enabled).toBe(false);
            (0, vitest_1.expect)((_b = registry.getServer('server2')) === null || _b === void 0 ? void 0 : _b.enabled).toBe(false);
        });
        (0, vitest_1.it)('should enable servers by type', function () {
            var _a, _b;
            registry.addServer({
                id: 'stdio-server',
                name: 'Stdio Server',
                enabled: false,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            });
            registry.addServer({
                id: 'sse-server',
                name: 'SSE Server',
                enabled: false,
                source: 'manual',
                transport: { type: 'sse', url: 'https://example.com', headers: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            });
            registry.enableType('stdio');
            (0, vitest_1.expect)((_a = registry.getServer('stdio-server')) === null || _a === void 0 ? void 0 : _a.enabled).toBe(true);
            (0, vitest_1.expect)((_b = registry.getServer('sse-server')) === null || _b === void 0 ? void 0 : _b.enabled).toBe(false);
        });
        (0, vitest_1.it)('should disable servers by type', function () {
            var _a, _b;
            registry.addServer({
                id: 'stdio-server',
                name: 'Stdio Server',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            });
            registry.addServer({
                id: 'sse-server',
                name: 'SSE Server',
                enabled: true,
                source: 'manual',
                transport: { type: 'sse', url: 'https://example.com', headers: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            });
            registry.disableType('stdio');
            (0, vitest_1.expect)((_a = registry.getServer('stdio-server')) === null || _a === void 0 ? void 0 : _a.enabled).toBe(false);
            (0, vitest_1.expect)((_b = registry.getServer('sse-server')) === null || _b === void 0 ? void 0 : _b.enabled).toBe(true);
        });
    });
    (0, vitest_1.describe)('updateServerStatus', function () {
        (0, vitest_1.it)('should update server status', function () {
            var _a, _b;
            registry.addServer({
                id: 'test',
                name: 'Test',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            });
            registry.updateServerStatus('test', {
                connected: true,
                lastConnected: new Date().toISOString(),
            });
            var server = registry.getServer('test');
            (0, vitest_1.expect)((_a = server === null || server === void 0 ? void 0 : server.status) === null || _a === void 0 ? void 0 : _a.connected).toBe(true);
            (0, vitest_1.expect)((_b = server === null || server === void 0 ? void 0 : server.status) === null || _b === void 0 ? void 0 : _b.lastConnected).toBeDefined();
        });
        (0, vitest_1.it)('should throw when updating status of non-existent server', function () {
            (0, vitest_1.expect)(function () {
                registry.updateServerStatus('nonexistent', { connected: true });
            }).toThrow(types_js_1.MCPServerNotFoundError);
        });
    });
    (0, vitest_1.describe)('validateServer', function () {
        (0, vitest_1.it)('should return valid for valid server', function () {
            registry.addServer({
                id: 'test',
                name: 'Test',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            });
            var result = registry.validateServer('test');
            (0, vitest_1.expect)(result.valid).toBe(true);
            (0, vitest_1.expect)(result.errors).toHaveLength(0);
        });
        (0, vitest_1.it)('should return invalid for non-existent server', function () {
            var result = registry.validateServer('nonexistent');
            (0, vitest_1.expect)(result.valid).toBe(false);
            (0, vitest_1.expect)(result.errors).toContain('Server not found: nonexistent');
        });
    });
    (0, vitest_1.describe)('testConnection', function () {
        (0, vitest_1.it)('should return success for existing server', function () { return __awaiter(void 0, void 0, void 0, function () {
            var result;
            var _a;
            return __generator(this, function (_b) {
                switch (_b.label) {
                    case 0:
                        registry.addServer({
                            id: 'test',
                            name: 'Test',
                            enabled: true,
                            source: 'manual',
                            transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                            settings: {
                                autoConnect: true,
                                timeout: 30000,
                                reconnect: true,
                                maxReconnectAttempts: 3,
                                fallbackServers: [],
                            },
                            capabilities: {
                                tools: [{ name: 'tool1', description: 'Tool 1', inputSchema: {} }],
                                resources: [],
                                prompts: [],
                            },
                        });
                        return [4 /*yield*/, registry.testConnection('test')];
                    case 1:
                        result = _b.sent();
                        (0, vitest_1.expect)(result.success).toBe(true);
                        (0, vitest_1.expect)(result.latency).toBeDefined();
                        (0, vitest_1.expect)((_a = result.capabilities) === null || _a === void 0 ? void 0 : _a.tools).toHaveLength(1);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should return error for non-existent server', function () { return __awaiter(void 0, void 0, void 0, function () {
            var result;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, registry.testConnection('nonexistent')];
                    case 1:
                        result = _a.sent();
                        (0, vitest_1.expect)(result.success).toBe(false);
                        (0, vitest_1.expect)(result.error).toBe('Server not found: nonexistent');
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('getReconnectDelay', function () {
        (0, vitest_1.it)('should return 0 for non-existent server', function () {
            var delay = registry.getReconnectDelay('nonexistent');
            (0, vitest_1.expect)(delay).toBe(0);
        });
        (0, vitest_1.it)('should return 0 for server without status', function () {
            registry.addServer({
                id: 'test',
                name: 'Test',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            });
            var delay = registry.getReconnectDelay('test');
            (0, vitest_1.expect)(delay).toBe(0);
        });
    });
    (0, vitest_1.describe)('fromJSON', function () {
        (0, vitest_1.it)('should throw for unsupported version', function () {
            (0, vitest_1.expect)(function () {
                registry.fromJSON({
                    version: 999,
                    servers: [],
                });
            }).toThrow(types_js_1.MCPValidationError);
        });
        (0, vitest_1.it)('should clear existing servers when loading', function () {
            registry.addServer({
                id: 'old',
                name: 'Old',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            });
            registry.fromJSON({
                version: 1,
                servers: [{
                        id: 'new',
                        name: 'New',
                        enabled: true,
                        source: 'manual',
                        transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                        settings: {
                            autoConnect: true,
                            timeout: 30000,
                            reconnect: true,
                            maxReconnectAttempts: 3,
                            fallbackServers: [],
                        },
                    }],
            });
            (0, vitest_1.expect)(registry.hasServer('old')).toBe(false);
            (0, vitest_1.expect)(registry.hasServer('new')).toBe(true);
        });
    });
    (0, vitest_1.describe)('clear', function () {
        (0, vitest_1.it)('should clear all servers', function () {
            registry.addServer({
                id: 'test',
                name: 'Test',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            });
            registry.clear();
            (0, vitest_1.expect)(registry.size).toBe(0);
            (0, vitest_1.expect)(registry.hasServer('test')).toBe(false);
        });
    });
    (0, vitest_1.describe)('getFallbackServers', function () {
        (0, vitest_1.it)('should return empty array for non-existent server', function () {
            var fallbacks = registry.getFallbackServers('nonexistent');
            (0, vitest_1.expect)(fallbacks).toEqual([]);
        });
        (0, vitest_1.it)('should filter out non-existent fallback servers', function () {
            registry.addServer({
                id: 'primary',
                name: 'Primary',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: ['exists', 'missing'],
                },
            });
            registry.addServer({
                id: 'exists',
                name: 'Exists',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            });
            var fallbacks = registry.getFallbackServers('primary');
            (0, vitest_1.expect)(fallbacks).toHaveLength(1);
            (0, vitest_1.expect)(fallbacks[0].id).toBe('exists');
        });
    });
});
