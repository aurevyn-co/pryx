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
var storage_js_1 = require("../../src/storage.js");
var validation_js_1 = require("../../src/validation.js");
var fs_1 = require("fs");
var path_1 = require("path");
var os_1 = require("os");
(0, vitest_1.describe)('MCP End-to-End Workflow', function () {
    var tempDir;
    var configPath;
    var storage;
    (0, vitest_1.beforeEach)(function () {
        tempDir = fs_1.default.mkdtempSync(path_1.default.join(os_1.default.tmpdir(), 'mcp-e2e-test-'));
        configPath = path_1.default.join(tempDir, 'mcp-servers.json');
        storage = new storage_js_1.MCPStorage();
    });
    (0, vitest_1.afterEach)(function () {
        fs_1.default.rmSync(tempDir, { recursive: true, force: true });
    });
    (0, vitest_1.describe)('complete workflow', function () {
        (0, vitest_1.it)('should handle full server lifecycle', function () { return __awaiter(void 0, void 0, void 0, function () {
            var registry, server1, server2, fallbacks, loadedRegistry, loadedServer1, loadedServer2, finalRegistry;
            var _a, _b;
            return __generator(this, function (_c) {
                switch (_c.label) {
                    case 0:
                        registry = new registry_js_1.MCPRegistry();
                        server1 = {
                            id: 'filesystem',
                            name: 'Filesystem Server',
                            enabled: true,
                            source: 'curated',
                            transport: {
                                type: 'stdio',
                                command: 'npx',
                                args: ['-y', '@modelcontextprotocol/server-filesystem', '/tmp'],
                                env: {},
                            },
                            capabilities: {
                                tools: [
                                    {
                                        name: 'read_file',
                                        description: 'Read a file',
                                        inputSchema: {
                                            type: 'object',
                                            properties: {
                                                path: { type: 'string' },
                                            },
                                        },
                                    },
                                ],
                                resources: [],
                                prompts: [],
                            },
                            settings: {
                                autoConnect: true,
                                timeout: 30000,
                                reconnect: true,
                                maxReconnectAttempts: 3,
                                fallbackServers: [],
                            },
                        };
                        server2 = {
                            id: 'remote-api',
                            name: 'Remote API Server',
                            enabled: true,
                            source: 'manual',
                            transport: {
                                type: 'sse',
                                url: 'https://api.example.com/mcp/sse',
                                headers: {
                                    Authorization: 'Bearer test-token',
                                },
                            },
                            settings: {
                                autoConnect: true,
                                timeout: 60000,
                                reconnect: true,
                                maxReconnectAttempts: 5,
                                fallbackServers: ['filesystem'],
                            },
                        };
                        registry.addServer(server1);
                        registry.addServer(server2);
                        (0, vitest_1.expect)(registry.size).toBe(2);
                        (0, vitest_1.expect)(registry.getEnabledServers()).toHaveLength(2);
                        fallbacks = registry.getFallbackServers('remote-api');
                        (0, vitest_1.expect)(fallbacks).toHaveLength(1);
                        (0, vitest_1.expect)(fallbacks[0].id).toBe('filesystem');
                        return [4 /*yield*/, storage.save(configPath, registry)];
                    case 1:
                        _c.sent();
                        (0, vitest_1.expect)(fs_1.default.existsSync(configPath)).toBe(true);
                        return [4 /*yield*/, storage.load(configPath)];
                    case 2:
                        loadedRegistry = _c.sent();
                        (0, vitest_1.expect)(loadedRegistry.size).toBe(2);
                        loadedServer1 = loadedRegistry.getServer('filesystem');
                        (0, vitest_1.expect)((_a = loadedServer1 === null || loadedServer1 === void 0 ? void 0 : loadedServer1.capabilities) === null || _a === void 0 ? void 0 : _a.tools).toHaveLength(1);
                        (0, vitest_1.expect)(loadedServer1 === null || loadedServer1 === void 0 ? void 0 : loadedServer1.transport.type).toBe('stdio');
                        loadedServer2 = loadedRegistry.getServer('remote-api');
                        (0, vitest_1.expect)(loadedServer2 === null || loadedServer2 === void 0 ? void 0 : loadedServer2.transport.type).toBe('sse');
                        (0, vitest_1.expect)(loadedServer2 === null || loadedServer2 === void 0 ? void 0 : loadedServer2.settings.timeout).toBe(60000);
                        loadedRegistry.updateServer('filesystem', { enabled: false });
                        (0, vitest_1.expect)(loadedRegistry.getEnabledServers()).toHaveLength(1);
                        loadedRegistry.removeServer('remote-api');
                        (0, vitest_1.expect)(loadedRegistry.size).toBe(1);
                        (0, vitest_1.expect)(loadedRegistry.hasServer('remote-api')).toBe(false);
                        return [4 /*yield*/, storage.save(configPath, loadedRegistry)];
                    case 3:
                        _c.sent();
                        return [4 /*yield*/, storage.load(configPath)];
                    case 4:
                        finalRegistry = _c.sent();
                        (0, vitest_1.expect)(finalRegistry.size).toBe(1);
                        (0, vitest_1.expect)((_b = finalRegistry.getServer('filesystem')) === null || _b === void 0 ? void 0 : _b.enabled).toBe(false);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should validate server configs correctly', function () {
            var validStdio = {
                id: 'valid-stdio',
                name: 'Valid Stdio',
                enabled: true,
                source: 'manual',
                transport: {
                    type: 'stdio',
                    command: 'npx',
                    args: ['-y', '@modelcontextprotocol/server-filesystem'],
                    env: {},
                },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            };
            var result = (0, validation_js_1.validateMCPServerConfig)(validStdio);
            (0, vitest_1.expect)(result.valid).toBe(true);
            (0, vitest_1.expect)(result.errors).toHaveLength(0);
        });
        (0, vitest_1.it)('should reject invalid server configs', function () {
            var invalidServer = {
                id: 'invalid',
                name: '',
                enabled: true,
                source: 'manual',
                transport: {
                    type: 'sse',
                    url: 'not-a-url',
                    headers: {},
                },
                settings: {
                    autoConnect: true,
                    timeout: -1,
                    reconnect: true,
                    maxReconnectAttempts: 3,
                    fallbackServers: [],
                },
            };
            var result = (0, validation_js_1.validateMCPServerConfig)(invalidServer);
            (0, vitest_1.expect)(result.valid).toBe(false);
            (0, vitest_1.expect)(result.errors.length).toBeGreaterThan(0);
        });
        (0, vitest_1.it)('should handle all transport types', function () { return __awaiter(void 0, void 0, void 0, function () {
            var registry, stdioServer, sseServer, wsServer, loaded, loadedStdio, loadedSse, loadedWs;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        registry = new registry_js_1.MCPRegistry();
                        stdioServer = {
                            id: 'stdio-server',
                            name: 'Stdio Server',
                            enabled: true,
                            source: 'manual',
                            transport: {
                                type: 'stdio',
                                command: 'python',
                                args: ['-m', 'mcp.server'],
                                env: { PYTHONPATH: '/app' },
                                cwd: '/home/user',
                            },
                            settings: {
                                autoConnect: true,
                                timeout: 30000,
                                reconnect: true,
                                maxReconnectAttempts: 3,
                                fallbackServers: [],
                            },
                        };
                        sseServer = {
                            id: 'sse-server',
                            name: 'SSE Server',
                            enabled: true,
                            source: 'curated',
                            transport: {
                                type: 'sse',
                                url: 'https://sse.example.com/events',
                                headers: { 'X-API-Key': 'secret' },
                            },
                            settings: {
                                autoConnect: true,
                                timeout: 30000,
                                reconnect: true,
                                maxReconnectAttempts: 3,
                                fallbackServers: [],
                            },
                        };
                        wsServer = {
                            id: 'ws-server',
                            name: 'WebSocket Server',
                            enabled: true,
                            source: 'marketplace',
                            transport: {
                                type: 'websocket',
                                url: 'wss://ws.example.com/mcp',
                                headers: {},
                            },
                            settings: {
                                autoConnect: true,
                                timeout: 30000,
                                reconnect: true,
                                maxReconnectAttempts: 3,
                                fallbackServers: [],
                            },
                        };
                        registry.addServer(stdioServer);
                        registry.addServer(sseServer);
                        registry.addServer(wsServer);
                        (0, vitest_1.expect)(registry.getServersByType('stdio')).toHaveLength(1);
                        (0, vitest_1.expect)(registry.getServersByType('sse')).toHaveLength(1);
                        (0, vitest_1.expect)(registry.getServersByType('websocket')).toHaveLength(1);
                        return [4 /*yield*/, storage.save(configPath, registry)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, storage.load(configPath)];
                    case 2:
                        loaded = _a.sent();
                        loadedStdio = loaded.getServer('stdio-server');
                        (0, vitest_1.expect)(loadedStdio === null || loadedStdio === void 0 ? void 0 : loadedStdio.transport.type).toBe('stdio');
                        if ((loadedStdio === null || loadedStdio === void 0 ? void 0 : loadedStdio.transport.type) === 'stdio') {
                            (0, vitest_1.expect)(loadedStdio.transport.cwd).toBe('/home/user');
                            (0, vitest_1.expect)(loadedStdio.transport.env.PYTHONPATH).toBe('/app');
                        }
                        loadedSse = loaded.getServer('sse-server');
                        (0, vitest_1.expect)(loadedSse === null || loadedSse === void 0 ? void 0 : loadedSse.transport.type).toBe('sse');
                        if ((loadedSse === null || loadedSse === void 0 ? void 0 : loadedSse.transport.type) === 'sse') {
                            (0, vitest_1.expect)(loadedSse.transport.headers['X-API-Key']).toBe('secret');
                        }
                        loadedWs = loaded.getServer('ws-server');
                        (0, vitest_1.expect)(loadedWs === null || loadedWs === void 0 ? void 0 : loadedWs.transport.type).toBe('websocket');
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should handle server status and reconnection', function () {
            var registry = new registry_js_1.MCPRegistry();
            registry.addServer({
                id: 'reconnect-test',
                name: 'Reconnect Test',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'test', args: [], env: {} },
                settings: {
                    autoConnect: true,
                    timeout: 30000,
                    reconnect: true,
                    maxReconnectAttempts: 5,
                    fallbackServers: [],
                },
                status: {
                    connected: false,
                    reconnectAttempts: 3,
                    lastError: 'Connection timeout',
                },
            });
            var delay1 = registry.getReconnectDelay('reconnect-test');
            (0, vitest_1.expect)(delay1).toBe(8000);
            registry.updateServer('reconnect-test', {
                status: {
                    connected: false,
                    reconnectAttempts: 5,
                },
            });
            var delay2 = registry.getReconnectDelay('reconnect-test');
            (0, vitest_1.expect)(delay2).toBe(30000);
        });
        (0, vitest_1.it)('should persist and restore complete server state', function () { return __awaiter(void 0, void 0, void 0, function () {
            var registry, now, loaded, server;
            var _a, _b, _c, _d, _e;
            return __generator(this, function (_f) {
                switch (_f.label) {
                    case 0:
                        registry = new registry_js_1.MCPRegistry();
                        now = new Date().toISOString();
                        registry.addServer({
                            id: 'complete-server',
                            name: 'Complete Server',
                            enabled: false,
                            source: 'curated',
                            transport: {
                                type: 'stdio',
                                command: 'node',
                                args: ['server.js'],
                                env: { NODE_ENV: 'production' },
                            },
                            capabilities: {
                                tools: [
                                    { name: 'tool1', description: 'Tool 1', inputSchema: {} },
                                    { name: 'tool2', description: 'Tool 2', inputSchema: {} },
                                ],
                                resources: [
                                    { uri: 'file:///data', name: 'Data', mimeType: 'application/json' },
                                ],
                                prompts: [
                                    { name: 'prompt1', description: 'Prompt 1' },
                                ],
                            },
                            settings: {
                                autoConnect: false,
                                timeout: 45000,
                                reconnect: false,
                                maxReconnectAttempts: 10,
                                fallbackServers: ['backup1', 'backup2'],
                            },
                            status: {
                                connected: true,
                                lastConnected: now,
                                reconnectAttempts: 0,
                            },
                        });
                        return [4 /*yield*/, storage.save(configPath, registry)];
                    case 1:
                        _f.sent();
                        return [4 /*yield*/, storage.load(configPath)];
                    case 2:
                        loaded = _f.sent();
                        server = loaded.getServer('complete-server');
                        (0, vitest_1.expect)(server === null || server === void 0 ? void 0 : server.enabled).toBe(false);
                        (0, vitest_1.expect)(server === null || server === void 0 ? void 0 : server.source).toBe('curated');
                        (0, vitest_1.expect)((_a = server === null || server === void 0 ? void 0 : server.capabilities) === null || _a === void 0 ? void 0 : _a.tools).toHaveLength(2);
                        (0, vitest_1.expect)((_b = server === null || server === void 0 ? void 0 : server.capabilities) === null || _b === void 0 ? void 0 : _b.resources).toHaveLength(1);
                        (0, vitest_1.expect)((_c = server === null || server === void 0 ? void 0 : server.capabilities) === null || _c === void 0 ? void 0 : _c.prompts).toHaveLength(1);
                        (0, vitest_1.expect)(server === null || server === void 0 ? void 0 : server.settings.autoConnect).toBe(false);
                        (0, vitest_1.expect)(server === null || server === void 0 ? void 0 : server.settings.timeout).toBe(45000);
                        (0, vitest_1.expect)(server === null || server === void 0 ? void 0 : server.settings.fallbackServers).toEqual(['backup1', 'backup2']);
                        (0, vitest_1.expect)((_d = server === null || server === void 0 ? void 0 : server.status) === null || _d === void 0 ? void 0 : _d.connected).toBe(true);
                        (0, vitest_1.expect)((_e = server === null || server === void 0 ? void 0 : server.status) === null || _e === void 0 ? void 0 : _e.lastConnected).toBe(now);
                        return [2 /*return*/];
                }
            });
        }); });
    });
});
