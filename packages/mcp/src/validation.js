"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.validateMCPServerConfig = validateMCPServerConfig;
exports.assertValidMCPServerConfig = assertValidMCPServerConfig;
exports.isValidServerId = isValidServerId;
exports.isValidTransportType = isValidTransportType;
exports.isValidUrl = isValidUrl;
exports.isValidWebSocketUrl = isValidWebSocketUrl;
exports.calculateBackoff = calculateBackoff;
var types_js_1 = require("./types.js");
function validateMCPServerConfig(config) {
    var result = types_js_1.MCPServerConfigSchema.safeParse(config);
    if (result.success) {
        var errors = [];
        var validated = result.data;
        switch (validated.transport.type) {
            case 'stdio': {
                if (!validated.transport.command) {
                    errors.push('stdio transport requires command');
                }
                break;
            }
            case 'sse':
            case 'websocket': {
                if (!validated.transport.url) {
                    errors.push("".concat(validated.transport.type, " transport requires url"));
                }
                if (validated.transport.type === 'websocket' &&
                    !validated.transport.url.match(/^wss?:\/\//)) {
                    errors.push('websocket URL must start with ws:// or wss://');
                }
                break;
            }
        }
        if (validated.settings.fallbackServers.includes(validated.id)) {
            errors.push('Server cannot be its own fallback');
        }
        return {
            valid: errors.length === 0,
            errors: errors,
        };
    }
    return {
        valid: false,
        errors: result.error.errors.map(function (e) { return "".concat(e.path.join('.'), ": ").concat(e.message); }),
    };
}
function assertValidMCPServerConfig(config) {
    var result = validateMCPServerConfig(config);
    if (!result.valid) {
        throw new types_js_1.MCPValidationError(result.errors);
    }
    return config;
}
function isValidServerId(id) {
    return /^[a-z0-9_-]+$/.test(id) && id.length >= 1 && id.length <= 64;
}
function isValidTransportType(type) {
    return ['stdio', 'sse', 'websocket'].includes(type);
}
function isValidUrl(url) {
    try {
        new URL(url);
        return true;
    }
    catch (_a) {
        return false;
    }
}
function isValidWebSocketUrl(url) {
    return isValidUrl(url) && /^wss?:\/\//.test(url);
}
function calculateBackoff(attempt, baseMs) {
    if (baseMs === void 0) { baseMs = 1000; }
    return Math.min(baseMs * Math.pow(2, attempt), 30000);
}
