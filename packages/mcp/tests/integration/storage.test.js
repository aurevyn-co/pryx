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
var storage_js_1 = require("../../src/storage.js");
var registry_js_1 = require("../../src/registry.js");
var fs_1 = require("fs");
var path_1 = require("path");
var os_1 = require("os");
(0, vitest_1.describe)('MCPRegistryStorage Integration', function () {
    var tempDir;
    var storage;
    var registry;
    var configPath;
    (0, vitest_1.beforeEach)(function () {
        tempDir = fs_1.default.mkdtempSync(path_1.default.join(os_1.default.tmpdir(), 'mcp-storage-test-'));
        configPath = path_1.default.join(tempDir, 'mcp-servers.json');
        storage = new storage_js_1.MCPStorage();
        registry = new registry_js_1.MCPRegistry();
    });
    (0, vitest_1.afterEach)(function () {
        fs_1.default.rmSync(tempDir, { recursive: true, force: true });
    });
    (0, vitest_1.describe)('save and load', function () {
        (0, vitest_1.it)('should save registry to file', function () { return __awaiter(void 0, void 0, void 0, function () {
            var server, content, parsed;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        server = {
                            id: 'filesystem',
                            name: 'Filesystem Server',
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
                        registry.addServer(server);
                        return [4 /*yield*/, storage.save(configPath, registry)];
                    case 1:
                        _a.sent();
                        (0, vitest_1.expect)(fs_1.default.existsSync(configPath)).toBe(true);
                        content = fs_1.default.readFileSync(configPath, 'utf-8');
                        parsed = JSON.parse(content);
                        (0, vitest_1.expect)(parsed.version).toBe(1);
                        (0, vitest_1.expect)(parsed.servers).toHaveLength(1);
                        (0, vitest_1.expect)(parsed.servers[0].id).toBe('filesystem');
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should load registry from file', function () { return __awaiter(void 0, void 0, void 0, function () {
            var config, loaded;
            var _a;
            return __generator(this, function (_b) {
                switch (_b.label) {
                    case 0:
                        config = {
                            version: 1,
                            servers: [{
                                    id: 'test-server',
                                    name: 'Test Server',
                                    enabled: true,
                                    source: 'curated',
                                    transport: {
                                        type: 'sse',
                                        url: 'https://api.example.com/sse',
                                        headers: {},
                                    },
                                    settings: {
                                        autoConnect: true,
                                        timeout: 30000,
                                        reconnect: true,
                                        maxReconnectAttempts: 3,
                                        fallbackServers: [],
                                    },
                                }],
                        };
                        fs_1.default.writeFileSync(configPath, JSON.stringify(config, null, 2));
                        return [4 /*yield*/, storage.load(configPath)];
                    case 1:
                        loaded = _b.sent();
                        (0, vitest_1.expect)(loaded.size).toBe(1);
                        (0, vitest_1.expect)(loaded.hasServer('test-server')).toBe(true);
                        (0, vitest_1.expect)((_a = loaded.getServer('test-server')) === null || _a === void 0 ? void 0 : _a.name).toBe('Test Server');
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should handle non-existent config file', function () { return __awaiter(void 0, void 0, void 0, function () {
            var loaded;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, storage.load(configPath)];
                    case 1:
                        loaded = _a.sent();
                        (0, vitest_1.expect)(loaded.size).toBe(0);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should handle invalid JSON gracefully', function () { return __awaiter(void 0, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        fs_1.default.writeFileSync(configPath, 'invalid json {');
                        return [4 /*yield*/, (0, vitest_1.expect)(storage.load(configPath)).rejects.toThrow()];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should preserve all server properties on save/load', function () { return __awaiter(void 0, void 0, void 0, function () {
            var server, loaded, loadedServer;
            var _a, _b, _c;
            return __generator(this, function (_d) {
                switch (_d.label) {
                    case 0:
                        server = {
                            id: 'complex-server',
                            name: 'Complex Server',
                            enabled: false,
                            source: 'marketplace',
                            transport: {
                                type: 'websocket',
                                url: 'wss://ws.example.com/mcp',
                                headers: { Authorization: 'Bearer token123' },
                            },
                            capabilities: {
                                tools: [
                                    {
                                        name: 'test-tool',
                                        description: 'A test tool',
                                        inputSchema: { type: 'object' },
                                    },
                                ],
                                resources: [
                                    {
                                        uri: 'file:///test',
                                        name: 'Test Resource',
                                        mimeType: 'text/plain',
                                    },
                                ],
                                prompts: [
                                    {
                                        name: 'test-prompt',
                                        description: 'A test prompt',
                                    },
                                ],
                            },
                            settings: {
                                autoConnect: false,
                                timeout: 60000,
                                reconnect: false,
                                maxReconnectAttempts: 5,
                                fallbackServers: ['fallback1', 'fallback2'],
                            },
                            status: {
                                connected: true,
                                lastConnected: new Date().toISOString(),
                                lastError: undefined,
                                reconnectAttempts: 0,
                            },
                        };
                        registry.addServer(server);
                        return [4 /*yield*/, storage.save(configPath, registry)];
                    case 1:
                        _d.sent();
                        return [4 /*yield*/, storage.load(configPath)];
                    case 2:
                        loaded = _d.sent();
                        loadedServer = loaded.getServer('complex-server');
                        (0, vitest_1.expect)(loadedServer).toBeDefined();
                        (0, vitest_1.expect)(loadedServer === null || loadedServer === void 0 ? void 0 : loadedServer.enabled).toBe(false);
                        (0, vitest_1.expect)(loadedServer === null || loadedServer === void 0 ? void 0 : loadedServer.source).toBe('marketplace');
                        (0, vitest_1.expect)((_a = loadedServer === null || loadedServer === void 0 ? void 0 : loadedServer.capabilities) === null || _a === void 0 ? void 0 : _a.tools).toHaveLength(1);
                        (0, vitest_1.expect)((_b = loadedServer === null || loadedServer === void 0 ? void 0 : loadedServer.capabilities) === null || _b === void 0 ? void 0 : _b.resources).toHaveLength(1);
                        (0, vitest_1.expect)(loadedServer === null || loadedServer === void 0 ? void 0 : loadedServer.settings.timeout).toBe(60000);
                        (0, vitest_1.expect)((_c = loadedServer === null || loadedServer === void 0 ? void 0 : loadedServer.status) === null || _c === void 0 ? void 0 : _c.connected).toBe(true);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('multiple servers', function () {
        (0, vitest_1.it)('should handle multiple servers', function () { return __awaiter(void 0, void 0, void 0, function () {
            var servers, loaded;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        servers = [
                            {
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
                            },
                            {
                                id: 'server2',
                                name: 'Server 2',
                                enabled: true,
                                source: 'curated',
                                transport: { type: 'sse', url: 'https://example.com', headers: {} },
                                settings: {
                                    autoConnect: true,
                                    timeout: 30000,
                                    reconnect: true,
                                    maxReconnectAttempts: 3,
                                    fallbackServers: [],
                                },
                            },
                            {
                                id: 'server3',
                                name: 'Server 3',
                                enabled: false,
                                source: 'manual',
                                transport: { type: 'websocket', url: 'wss://example.com', headers: {} },
                                settings: {
                                    autoConnect: true,
                                    timeout: 30000,
                                    reconnect: true,
                                    maxReconnectAttempts: 3,
                                    fallbackServers: [],
                                },
                            },
                        ];
                        servers.forEach(function (s) { registry.addServer(s); });
                        return [4 /*yield*/, storage.save(configPath, registry)];
                    case 1:
                        _a.sent();
                        return [4 /*yield*/, storage.load(configPath)];
                    case 2:
                        loaded = _a.sent();
                        (0, vitest_1.expect)(loaded.size).toBe(3);
                        (0, vitest_1.expect)(loaded.getEnabledServers()).toHaveLength(2);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('file permissions', function () {
        (0, vitest_1.it)('should create file with secure permissions', function () { return __awaiter(void 0, void 0, void 0, function () {
            var server, stats, mode;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        server = {
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
                        };
                        registry.addServer(server);
                        return [4 /*yield*/, storage.save(configPath, registry)];
                    case 1:
                        _a.sent();
                        stats = fs_1.default.statSync(configPath);
                        mode = stats.mode & 511;
                        (0, vitest_1.expect)(mode).toBe(384);
                        return [2 /*return*/];
                }
            });
        }); });
    });
});
