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
Object.defineProperty(exports, "__esModule", { value: true });
var vitest_1 = require("vitest");
var registry_js_1 = require("../../src/registry.js");
var types_js_1 = require("../../src/types.js");
var defaultSettings = {
    autoConnect: true,
    timeout: 30000,
    reconnect: true,
    maxReconnectAttempts: 3,
    fallbackServers: [],
};
(0, vitest_1.describe)('MCPRegistry', function () {
    var registry;
    (0, vitest_1.beforeEach)(function () {
        registry = new registry_js_1.MCPRegistry();
    });
    (0, vitest_1.describe)('constructor', function () {
        (0, vitest_1.it)('should create empty registry', function () {
            (0, vitest_1.expect)(registry.size).toBe(0);
        });
    });
    (0, vitest_1.describe)('addServer', function () {
        (0, vitest_1.it)('should add stdio server', function () {
            var _a;
            var server = {
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
                settings: defaultSettings,
            };
            registry.addServer(server);
            (0, vitest_1.expect)(registry.hasServer('filesystem')).toBe(true);
            (0, vitest_1.expect)((_a = registry.getServer('filesystem')) === null || _a === void 0 ? void 0 : _a.name).toBe('Filesystem Server');
        });
        (0, vitest_1.it)('should add sse server', function () {
            var server = {
                id: 'remote-api',
                name: 'Remote API',
                enabled: true,
                source: 'manual',
                transport: {
                    type: 'sse',
                    url: 'https://api.example.com/sse',
                    headers: {},
                },
                settings: defaultSettings,
            };
            registry.addServer(server);
            (0, vitest_1.expect)(registry.hasServer('remote-api')).toBe(true);
        });
        (0, vitest_1.it)('should add websocket server', function () {
            var server = {
                id: 'ws-server',
                name: 'WebSocket Server',
                enabled: true,
                source: 'curated',
                transport: {
                    type: 'websocket',
                    url: 'wss://ws.example.com/mcp',
                    headers: {},
                },
                settings: defaultSettings,
            };
            registry.addServer(server);
            (0, vitest_1.expect)(registry.hasServer('ws-server')).toBe(true);
        });
        (0, vitest_1.it)('should throw when server already exists', function () {
            var server = {
                id: 'test',
                name: 'Test',
                enabled: true,
                source: 'manual',
                transport: {
                    type: 'stdio',
                    command: 'test',
                    args: [],
                    env: {},
                },
                settings: defaultSettings,
            };
            registry.addServer(server);
            (0, vitest_1.expect)(function () { return registry.addServer(server); }).toThrow(types_js_1.MCPServerAlreadyExistsError);
        });
    });
    (0, vitest_1.describe)('updateServer', function () {
        (0, vitest_1.it)('should update existing server', function () {
            var _a;
            registry.addServer({
                id: 'test',
                name: 'Test',
                enabled: true,
                source: 'manual',
                transport: {
                    type: 'stdio',
                    command: 'test',
                    args: [],
                    env: {},
                },
                settings: defaultSettings,
            });
            var updated = registry.updateServer('test', { name: 'Updated' });
            (0, vitest_1.expect)(updated.name).toBe('Updated');
            (0, vitest_1.expect)((_a = registry.getServer('test')) === null || _a === void 0 ? void 0 : _a.name).toBe('Updated');
        });
        (0, vitest_1.it)('should throw when server not found', function () {
            (0, vitest_1.expect)(function () { return registry.updateServer('nonexistent', { name: 'Test' }); }).toThrow(types_js_1.MCPServerNotFoundError);
        });
    });
    (0, vitest_1.describe)('removeServer', function () {
        (0, vitest_1.it)('should remove server', function () {
            registry.addServer({
                id: 'test',
                name: 'Test',
                enabled: true,
                source: 'manual',
                transport: {
                    type: 'stdio',
                    command: 'test',
                    args: [],
                    env: {},
                },
                settings: defaultSettings,
            });
            registry.removeServer('test');
            (0, vitest_1.expect)(registry.hasServer('test')).toBe(false);
        });
        (0, vitest_1.it)('should throw when server not found', function () {
            (0, vitest_1.expect)(function () { return registry.removeServer('nonexistent'); }).toThrow(types_js_1.MCPServerNotFoundError);
        });
    });
    (0, vitest_1.describe)('getAllServers', function () {
        (0, vitest_1.it)('should return all servers', function () {
            registry.addServer({
                id: 'server1',
                name: 'Server 1',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd1', args: [], env: {} },
                settings: defaultSettings,
            });
            registry.addServer({
                id: 'server2',
                name: 'Server 2',
                enabled: true,
                source: 'curated',
                transport: { type: 'stdio', command: 'cmd2', args: [], env: {} },
                settings: defaultSettings,
            });
            var servers = registry.getAllServers();
            (0, vitest_1.expect)(servers.length).toBe(2);
        });
    });
    (0, vitest_1.describe)('getEnabledServers', function () {
        (0, vitest_1.it)('should return only enabled servers', function () {
            registry.addServer({
                id: 'enabled',
                name: 'Enabled',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                settings: defaultSettings,
            });
            registry.addServer({
                id: 'disabled',
                name: 'Disabled',
                enabled: false,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                settings: defaultSettings,
            });
            var enabled = registry.getEnabledServers();
            (0, vitest_1.expect)(enabled.length).toBe(1);
            (0, vitest_1.expect)(enabled[0].id).toBe('enabled');
        });
    });
    (0, vitest_1.describe)('getServersByType', function () {
        (0, vitest_1.it)('should return servers by transport type', function () {
            registry.addServer({
                id: 'stdio1',
                name: 'Stdio 1',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                settings: defaultSettings,
            });
            registry.addServer({
                id: 'sse1',
                name: 'SSE 1',
                enabled: true,
                source: 'curated',
                transport: { type: 'sse', url: 'https://example.com', headers: {} },
                settings: defaultSettings,
            });
            var stdioServers = registry.getServersByType('stdio');
            (0, vitest_1.expect)(stdioServers.length).toBe(1);
            (0, vitest_1.expect)(stdioServers[0].id).toBe('stdio1');
        });
    });
    (0, vitest_1.describe)('fallback servers', function () {
        (0, vitest_1.it)('should return fallback servers', function () {
            registry.addServer({
                id: 'primary',
                name: 'Primary',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                settings: __assign(__assign({}, defaultSettings), { fallbackServers: ['fallback1', 'fallback2'] }),
            });
            registry.addServer({
                id: 'fallback1',
                name: 'Fallback 1',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                settings: defaultSettings,
            });
            registry.addServer({
                id: 'fallback2',
                name: 'Fallback 2',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                settings: defaultSettings,
            });
            var fallbacks = registry.getFallbackServers('primary');
            (0, vitest_1.expect)(fallbacks.length).toBe(2);
        });
    });
    (0, vitest_1.describe)('getReconnectDelay', function () {
        (0, vitest_1.it)('should calculate reconnect delay', function () {
            registry.addServer({
                id: 'test',
                name: 'Test',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                settings: defaultSettings,
                status: {
                    connected: false,
                    reconnectAttempts: 2,
                },
            });
            var delay = registry.getReconnectDelay('test');
            (0, vitest_1.expect)(delay).toBe(4000);
        });
    });
    (0, vitest_1.describe)('toJSON / fromJSON', function () {
        (0, vitest_1.it)('should serialize to JSON', function () {
            registry.addServer({
                id: 'test',
                name: 'Test',
                enabled: true,
                source: 'manual',
                transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                settings: defaultSettings,
            });
            var json = registry.toJSON();
            (0, vitest_1.expect)(json.version).toBe(1);
            (0, vitest_1.expect)(json.servers.length).toBe(1);
        });
        (0, vitest_1.it)('should deserialize from JSON', function () {
            var json = {
                version: 1,
                servers: [{
                        id: 'test',
                        name: 'Test',
                        enabled: true,
                        source: 'manual',
                        transport: { type: 'stdio', command: 'cmd', args: [], env: {} },
                        settings: defaultSettings,
                    }],
            };
            registry.fromJSON(json);
            (0, vitest_1.expect)(registry.size).toBe(1);
            (0, vitest_1.expect)(registry.hasServer('test')).toBe(true);
        });
    });
});
(0, vitest_1.describe)('createRegistry', function () {
    (0, vitest_1.it)('should create new registry', function () {
        var registry = (0, registry_js_1.createRegistry)();
        (0, vitest_1.expect)(registry).toBeInstanceOf(registry_js_1.MCPRegistry);
        (0, vitest_1.expect)(registry.size).toBe(0);
    });
});
