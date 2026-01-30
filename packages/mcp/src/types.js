"use strict";
var __extends = (this && this.__extends) || (function () {
    var extendStatics = function (d, b) {
        extendStatics = Object.setPrototypeOf ||
            ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
            function (d, b) { for (var p in b) if (Object.prototype.hasOwnProperty.call(b, p)) d[p] = b[p]; };
        return extendStatics(d, b);
    };
    return function (d, b) {
        if (typeof b !== "function" && b !== null)
            throw new TypeError("Class extends value " + String(b) + " is not a constructor or null");
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
exports.CURRENT_VERSION = exports.MCPServerAlreadyExistsError = exports.MCPValidationError = exports.MCPServerNotFoundError = exports.MCPError = exports.ConnectionTestResultSchema = exports.ValidationResultSchema = exports.MCPServersConfigSchema = exports.MCPServerConfigSchema = exports.ConnectionStatusSchema = exports.ServerSettingsSchema = exports.TransportConfigSchema = exports.WebSocketTransportSchema = exports.SSETransportSchema = exports.StdioTransportSchema = exports.CapabilitiesSchema = exports.PromptDefinitionSchema = exports.ArgumentDefinitionSchema = exports.ResourceDefinitionSchema = exports.ToolDefinitionSchema = exports.ServerSource = exports.TransportType = void 0;
var zod_1 = require("zod");
exports.TransportType = zod_1.z.enum(['stdio', 'sse', 'websocket']);
exports.ServerSource = zod_1.z.enum(['manual', 'curated', 'marketplace']);
exports.ToolDefinitionSchema = zod_1.z.object({
    name: zod_1.z.string().min(1),
    description: zod_1.z.string(),
    inputSchema: zod_1.z.record(zod_1.z.unknown()),
});
exports.ResourceDefinitionSchema = zod_1.z.object({
    uri: zod_1.z.string().url(),
    name: zod_1.z.string().min(1),
    mimeType: zod_1.z.string().optional(),
});
exports.ArgumentDefinitionSchema = zod_1.z.object({
    name: zod_1.z.string().min(1),
    description: zod_1.z.string(),
    required: zod_1.z.boolean().default(false),
});
exports.PromptDefinitionSchema = zod_1.z.object({
    name: zod_1.z.string().min(1),
    description: zod_1.z.string(),
    arguments: zod_1.z.array(exports.ArgumentDefinitionSchema).optional(),
});
exports.CapabilitiesSchema = zod_1.z.object({
    tools: zod_1.z.array(exports.ToolDefinitionSchema).default([]),
    resources: zod_1.z.array(exports.ResourceDefinitionSchema).default([]),
    prompts: zod_1.z.array(exports.PromptDefinitionSchema).default([]),
});
exports.StdioTransportSchema = zod_1.z.object({
    type: zod_1.z.literal('stdio'),
    command: zod_1.z.string().min(1),
    args: zod_1.z.array(zod_1.z.string()).default([]),
    env: zod_1.z.record(zod_1.z.string()).default({}),
    cwd: zod_1.z.string().optional(),
});
exports.SSETransportSchema = zod_1.z.object({
    type: zod_1.z.literal('sse'),
    url: zod_1.z.string().url(),
    headers: zod_1.z.record(zod_1.z.string()).default({}),
});
exports.WebSocketTransportSchema = zod_1.z.object({
    type: zod_1.z.literal('websocket'),
    url: zod_1.z.string().url().regex(/^wss?:\/\//),
    headers: zod_1.z.record(zod_1.z.string()).default({}),
});
exports.TransportConfigSchema = zod_1.z.union([
    exports.StdioTransportSchema,
    exports.SSETransportSchema,
    exports.WebSocketTransportSchema,
]);
exports.ServerSettingsSchema = zod_1.z.object({
    autoConnect: zod_1.z.boolean().default(true),
    timeout: zod_1.z.number().int().positive().default(30000),
    reconnect: zod_1.z.boolean().default(true),
    maxReconnectAttempts: zod_1.z.number().int().min(0).default(3),
    fallbackServers: zod_1.z.array(zod_1.z.string()).default([]),
});
exports.ConnectionStatusSchema = zod_1.z.object({
    connected: zod_1.z.boolean(),
    lastConnected: zod_1.z.string().datetime().optional(),
    lastError: zod_1.z.string().optional(),
    reconnectAttempts: zod_1.z.number().int().min(0).default(0),
});
exports.MCPServerConfigSchema = zod_1.z.object({
    id: zod_1.z.string().min(1).max(64).regex(/^[a-z0-9_-]+$/),
    name: zod_1.z.string().min(1).max(128),
    enabled: zod_1.z.boolean().default(true),
    transport: exports.TransportConfigSchema,
    capabilities: exports.CapabilitiesSchema.optional(),
    source: exports.ServerSource.default('manual'),
    settings: exports.ServerSettingsSchema.default({}),
    status: exports.ConnectionStatusSchema.optional(),
});
exports.MCPServersConfigSchema = zod_1.z.object({
    version: zod_1.z.number().int().default(1),
    servers: zod_1.z.array(exports.MCPServerConfigSchema),
});
exports.ValidationResultSchema = zod_1.z.object({
    valid: zod_1.z.boolean(),
    errors: zod_1.z.array(zod_1.z.string()),
});
exports.ConnectionTestResultSchema = zod_1.z.object({
    success: zod_1.z.boolean(),
    latency: zod_1.z.number().optional(),
    error: zod_1.z.string().optional(),
    capabilities: exports.CapabilitiesSchema.optional(),
});
var MCPError = /** @class */ (function (_super) {
    __extends(MCPError, _super);
    function MCPError(message) {
        var _this = _super.call(this, message) || this;
        _this.name = 'MCPError';
        return _this;
    }
    return MCPError;
}(Error));
exports.MCPError = MCPError;
var MCPServerNotFoundError = /** @class */ (function (_super) {
    __extends(MCPServerNotFoundError, _super);
    function MCPServerNotFoundError(id) {
        var _this = _super.call(this, "MCP server not found: ".concat(id)) || this;
        _this.name = 'MCPServerNotFoundError';
        return _this;
    }
    return MCPServerNotFoundError;
}(MCPError));
exports.MCPServerNotFoundError = MCPServerNotFoundError;
var MCPValidationError = /** @class */ (function (_super) {
    __extends(MCPValidationError, _super);
    function MCPValidationError(errors) {
        var _this = _super.call(this, "Validation failed: ".concat(errors.join(', '))) || this;
        _this.errors = errors;
        _this.name = 'MCPValidationError';
        return _this;
    }
    return MCPValidationError;
}(MCPError));
exports.MCPValidationError = MCPValidationError;
var MCPServerAlreadyExistsError = /** @class */ (function (_super) {
    __extends(MCPServerAlreadyExistsError, _super);
    function MCPServerAlreadyExistsError(id) {
        var _this = _super.call(this, "MCP server already exists: ".concat(id)) || this;
        _this.name = 'MCPServerAlreadyExistsError';
        return _this;
    }
    return MCPServerAlreadyExistsError;
}(MCPError));
exports.MCPServerAlreadyExistsError = MCPServerAlreadyExistsError;
exports.CURRENT_VERSION = 1;
