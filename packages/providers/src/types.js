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
exports.CURRENT_VERSION = exports.ProviderAlreadyExistsError = exports.ProviderValidationError = exports.ProviderNotFoundError = exports.ProviderError = exports.ConnectionTestResultSchema = exports.ValidationResultSchema = exports.ProvidersConfigSchema = exports.ProviderConfigSchema = exports.RateLimitConfigSchema = exports.ModelConfigSchema = exports.ProviderType = void 0;
var zod_1 = require("zod");
exports.ProviderType = zod_1.z.enum(['openai', 'anthropic', 'google', 'local', 'custom']);
exports.ModelConfigSchema = zod_1.z.object({
    id: zod_1.z.string().min(1),
    name: zod_1.z.string().min(1),
    maxTokens: zod_1.z.number().int().positive(),
    supportsStreaming: zod_1.z.boolean().default(true),
    supportsVision: zod_1.z.boolean().default(false),
    supportsTools: zod_1.z.boolean().default(false),
    costPer1KInput: zod_1.z.number().positive().optional(),
    costPer1KOutput: zod_1.z.number().positive().optional(),
});
exports.RateLimitConfigSchema = zod_1.z.object({
    requestsPerMinute: zod_1.z.number().int().positive().optional(),
    tokensPerMinute: zod_1.z.number().int().positive().optional(),
    requestsPerDay: zod_1.z.number().int().positive().optional(),
});
exports.ProviderConfigSchema = zod_1.z.object({
    id: zod_1.z.string().min(1).max(64).regex(/^[a-z0-9_-]+$/),
    name: zod_1.z.string().min(1).max(128),
    type: exports.ProviderType,
    enabled: zod_1.z.boolean().default(true),
    defaultModel: zod_1.z.string().optional(),
    apiKey: zod_1.z.string().optional(),
    baseUrl: zod_1.z.string().url().optional(),
    organization: zod_1.z.string().optional(),
    models: zod_1.z.array(exports.ModelConfigSchema).min(1),
    rateLimits: exports.RateLimitConfigSchema.optional(),
    timeout: zod_1.z.number().int().positive().default(30000),
    retries: zod_1.z.number().int().min(0).max(10).default(3),
});
exports.ProvidersConfigSchema = zod_1.z.object({
    version: zod_1.z.number().int().default(1),
    defaultProvider: zod_1.z.string().optional(),
    providers: zod_1.z.array(exports.ProviderConfigSchema),
});
exports.ValidationResultSchema = zod_1.z.object({
    valid: zod_1.z.boolean(),
    errors: zod_1.z.array(zod_1.z.string()),
});
exports.ConnectionTestResultSchema = zod_1.z.object({
    success: zod_1.z.boolean(),
    latency: zod_1.z.number().optional(),
    error: zod_1.z.string().optional(),
    modelsAvailable: zod_1.z.array(zod_1.z.string()).optional(),
});
var ProviderError = /** @class */ (function (_super) {
    __extends(ProviderError, _super);
    function ProviderError(message) {
        var _this = _super.call(this, message) || this;
        _this.name = 'ProviderError';
        return _this;
    }
    return ProviderError;
}(Error));
exports.ProviderError = ProviderError;
var ProviderNotFoundError = /** @class */ (function (_super) {
    __extends(ProviderNotFoundError, _super);
    function ProviderNotFoundError(id) {
        var _this = _super.call(this, "Provider not found: ".concat(id)) || this;
        _this.name = 'ProviderNotFoundError';
        return _this;
    }
    return ProviderNotFoundError;
}(ProviderError));
exports.ProviderNotFoundError = ProviderNotFoundError;
var ProviderValidationError = /** @class */ (function (_super) {
    __extends(ProviderValidationError, _super);
    function ProviderValidationError(errors) {
        var _this = _super.call(this, "Validation failed: ".concat(errors.join(', '))) || this;
        _this.errors = errors;
        _this.name = 'ProviderValidationError';
        return _this;
    }
    return ProviderValidationError;
}(ProviderError));
exports.ProviderValidationError = ProviderValidationError;
var ProviderAlreadyExistsError = /** @class */ (function (_super) {
    __extends(ProviderAlreadyExistsError, _super);
    function ProviderAlreadyExistsError(id) {
        var _this = _super.call(this, "Provider already exists: ".concat(id)) || this;
        _this.name = 'ProviderAlreadyExistsError';
        return _this;
    }
    return ProviderAlreadyExistsError;
}(ProviderError));
exports.ProviderAlreadyExistsError = ProviderAlreadyExistsError;
exports.CURRENT_VERSION = 1;
