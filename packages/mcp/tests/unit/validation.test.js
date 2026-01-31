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
var validation_js_1 = require("../../src/validation.js");
var types_js_1 = require("../../src/types.js");
(0, vitest_1.describe)('validateMCPServerConfig', function () {
    var baseConfig = {
        id: 'test-server',
        name: 'Test Server',
        enabled: true,
        transport: {
            type: 'stdio',
            command: 'npx',
            args: ['-y', '@modelcontextprotocol/server-filesystem'],
        },
        settings: {
            autoConnect: true,
            timeout: 30000,
            reconnect: true,
            maxReconnectAttempts: 3,
            fallbackServers: [],
        },
    };
    (0, vitest_1.it)('should validate correct stdio config', function () {
        var result = (0, validation_js_1.validateMCPServerConfig)(baseConfig);
        (0, vitest_1.expect)(result.valid).toBe(true);
        (0, vitest_1.expect)(result.errors).toHaveLength(0);
    });
    (0, vitest_1.it)('should validate correct sse config', function () {
        var config = __assign(__assign({}, baseConfig), { transport: {
                type: 'sse',
                url: 'https://api.example.com/sse',
            } });
        var result = (0, validation_js_1.validateMCPServerConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(true);
    });
    (0, vitest_1.it)('should validate correct websocket config', function () {
        var config = __assign(__assign({}, baseConfig), { transport: {
                type: 'websocket',
                url: 'wss://api.example.com/ws',
            } });
        var result = (0, validation_js_1.validateMCPServerConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(true);
    });
    (0, vitest_1.it)('should reject config with invalid id format', function () {
        var config = __assign(__assign({}, baseConfig), { id: 'Invalid ID!' });
        var result = (0, validation_js_1.validateMCPServerConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
    });
    (0, vitest_1.it)('should reject config with empty name', function () {
        var config = __assign(__assign({}, baseConfig), { name: '' });
        var result = (0, validation_js_1.validateMCPServerConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
    });
    (0, vitest_1.it)('should reject stdio config without command', function () {
        var config = __assign(__assign({}, baseConfig), { transport: {
                type: 'stdio',
                command: '',
            } });
        var result = (0, validation_js_1.validateMCPServerConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
    });
    (0, vitest_1.it)('should reject sse config without url', function () {
        var config = __assign(__assign({}, baseConfig), { transport: {
                type: 'sse',
                url: '',
            } });
        var result = (0, validation_js_1.validateMCPServerConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
    });
    (0, vitest_1.it)('should reject websocket config with invalid url', function () {
        var config = __assign(__assign({}, baseConfig), { transport: {
                type: 'websocket',
                url: 'https://example.com',
            } });
        var result = (0, validation_js_1.validateMCPServerConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
    });
    (0, vitest_1.it)('should reject config with self as fallback', function () {
        var config = __assign(__assign({}, baseConfig), { settings: __assign(__assign({}, baseConfig.settings), { fallbackServers: ['test-server'] }) });
        var result = (0, validation_js_1.validateMCPServerConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
    });
});
(0, vitest_1.describe)('assertValidMCPServerConfig', function () {
    var validConfig = {
        id: 'test',
        name: 'Test',
        enabled: true,
        transport: {
            type: 'stdio',
            command: 'test',
        },
        settings: {
            autoConnect: true,
            timeout: 30000,
            reconnect: true,
            maxReconnectAttempts: 3,
            fallbackServers: [],
        },
    };
    (0, vitest_1.it)('should return config when valid', function () {
        var result = (0, validation_js_1.assertValidMCPServerConfig)(validConfig);
        (0, vitest_1.expect)(result.id).toBe('test');
    });
    (0, vitest_1.it)('should throw when invalid', function () {
        var config = __assign(__assign({}, validConfig), { id: '' });
        (0, vitest_1.expect)(function () { return (0, validation_js_1.assertValidMCPServerConfig)(config); }).toThrow(types_js_1.MCPValidationError);
    });
});
(0, vitest_1.describe)('isValidServerId', function () {
    (0, vitest_1.it)('should return true for valid ids', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidServerId)('filesystem')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidServerId)('github-server')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidServerId)('a')).toBe(true);
    });
    (0, vitest_1.it)('should return false for invalid ids', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidServerId)('')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidServerId)('Invalid ID')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidServerId)('test@server')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidServerId)('a'.repeat(65))).toBe(false);
    });
});
(0, vitest_1.describe)('isValidTransportType', function () {
    (0, vitest_1.it)('should return true for valid types', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidTransportType)('stdio')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidTransportType)('sse')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidTransportType)('websocket')).toBe(true);
    });
    (0, vitest_1.it)('should return false for invalid types', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidTransportType)('http')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidTransportType)('grpc')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidTransportType)('')).toBe(false);
    });
});
(0, vitest_1.describe)('isValidUrl', function () {
    (0, vitest_1.it)('should return true for valid URLs', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidUrl)('https://api.example.com')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidUrl)('http://localhost:3000')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidUrl)('wss://ws.example.com')).toBe(true);
    });
    (0, vitest_1.it)('should return false for invalid URLs', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidUrl)('not-a-url')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidUrl)('')).toBe(false);
    });
});
(0, vitest_1.describe)('isValidWebSocketUrl', function () {
    (0, vitest_1.it)('should return true for valid WebSocket URLs', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidWebSocketUrl)('ws://localhost:3000')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidWebSocketUrl)('wss://api.example.com/ws')).toBe(true);
    });
    (0, vitest_1.it)('should return false for non-WebSocket URLs', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidWebSocketUrl)('https://example.com')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidWebSocketUrl)('http://localhost')).toBe(false);
    });
});
(0, vitest_1.describe)('calculateBackoff', function () {
    (0, vitest_1.it)('should calculate exponential backoff', function () {
        (0, vitest_1.expect)((0, validation_js_1.calculateBackoff)(0, 1000)).toBe(1000);
        (0, vitest_1.expect)((0, validation_js_1.calculateBackoff)(1, 1000)).toBe(2000);
        (0, vitest_1.expect)((0, validation_js_1.calculateBackoff)(2, 1000)).toBe(4000);
        (0, vitest_1.expect)((0, validation_js_1.calculateBackoff)(3, 1000)).toBe(8000);
    });
    (0, vitest_1.it)('should cap at 30 seconds', function () {
        (0, vitest_1.expect)((0, validation_js_1.calculateBackoff)(10, 1000)).toBe(30000);
    });
    (0, vitest_1.it)('should use default base of 1000ms', function () {
        (0, vitest_1.expect)((0, validation_js_1.calculateBackoff)(0)).toBe(1000);
    });
});
