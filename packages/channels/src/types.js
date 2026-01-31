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
exports.CURRENT_VERSION = exports.ChannelAlreadyExistsError = exports.ChannelValidationError = exports.ChannelNotFoundError = exports.ChannelError = exports.ConnectionTestResultSchema = exports.ValidationResultSchema = exports.ChannelsConfigSchema = exports.ChannelConfigSchema = exports.ConnectionStatusSchema = exports.WebhookSettingsSchema = exports.ChannelSettingsSchema = exports.WebhookConfigSchema = exports.WebhookRetryPolicySchema = exports.WebhookAuthSchema = exports.WhatsAppConfigSchema = exports.EmailConfigSchema = exports.EmailServerConfigSchema = exports.SlackConfigSchema = exports.DiscordConfigSchema = exports.TelegramConfigSchema = exports.RateLimitConfigSchema = exports.ChannelType = void 0;
var zod_1 = require("zod");
exports.ChannelType = zod_1.z.enum(['telegram', 'discord', 'slack', 'email', 'whatsapp', 'webhook']);
exports.RateLimitConfigSchema = zod_1.z.object({
    requestsPerMinute: zod_1.z.number().int().positive().optional(),
    tokensPerMinute: zod_1.z.number().int().positive().optional(),
    requestsPerDay: zod_1.z.number().int().positive().optional(),
});
exports.TelegramConfigSchema = zod_1.z.object({
    botToken: zod_1.z.string().min(1),
    chatId: zod_1.z.string().optional(),
    parseMode: zod_1.z.enum(['HTML', 'Markdown', 'MarkdownV2']).optional(),
    disableNotification: zod_1.z.boolean().optional(),
});
exports.DiscordConfigSchema = zod_1.z.object({
    botToken: zod_1.z.string().min(1),
    applicationId: zod_1.z.string().min(1),
    guildId: zod_1.z.string().optional(),
    channelId: zod_1.z.string().optional(),
    intents: zod_1.z.array(zod_1.z.string()).default([]),
});
exports.SlackConfigSchema = zod_1.z.object({
    botToken: zod_1.z.string().min(1),
    appToken: zod_1.z.string().optional(),
    signingSecret: zod_1.z.string().optional(),
    channelId: zod_1.z.string().optional(),
    socketMode: zod_1.z.boolean().default(false),
});
exports.EmailServerConfigSchema = zod_1.z.object({
    host: zod_1.z.string().min(1),
    port: zod_1.z.number().int().positive(),
    secure: zod_1.z.boolean(),
    username: zod_1.z.string().min(1),
    password: zod_1.z.string().min(1),
});
exports.EmailConfigSchema = zod_1.z.object({
    imap: exports.EmailServerConfigSchema.optional(),
    smtp: exports.EmailServerConfigSchema.optional(),
    checkInterval: zod_1.z.number().int().positive().default(60000),
    markAsRead: zod_1.z.boolean().default(true),
});
exports.WhatsAppConfigSchema = zod_1.z.object({
    sessionName: zod_1.z.string().min(1),
    phoneNumber: zod_1.z.string().optional(),
    qrTimeout: zod_1.z.number().int().positive().default(60000),
    pairingCode: zod_1.z.boolean().default(false),
});
exports.WebhookAuthSchema = zod_1.z.object({
    type: zod_1.z.enum(['bearer', 'basic', 'api-key']),
    token: zod_1.z.string().optional(),
    username: zod_1.z.string().optional(),
    password: zod_1.z.string().optional(),
});
exports.WebhookRetryPolicySchema = zod_1.z.object({
    maxRetries: zod_1.z.number().int().min(0).default(3),
    backoffMs: zod_1.z.number().int().positive().default(1000),
});
exports.WebhookConfigSchema = zod_1.z.object({
    url: zod_1.z.string().url(),
    method: zod_1.z.enum(['GET', 'POST', 'PUT', 'DELETE']).default('POST'),
    headers: zod_1.z.record(zod_1.z.string()).default({}),
    auth: exports.WebhookAuthSchema.optional(),
    retryPolicy: exports.WebhookRetryPolicySchema.default({}),
});
exports.ChannelSettingsSchema = zod_1.z.object({
    allowCommands: zod_1.z.boolean().default(true),
    autoReply: zod_1.z.boolean().default(false),
    filterPatterns: zod_1.z.array(zod_1.z.string()).default([]),
    allowedUsers: zod_1.z.array(zod_1.z.string()).default([]),
    blockedUsers: zod_1.z.array(zod_1.z.string()).default([]),
    rateLimit: exports.RateLimitConfigSchema.optional(),
});
exports.WebhookSettingsSchema = zod_1.z.object({
    url: zod_1.z.string().url(),
    secret: zod_1.z.string().optional(),
    enabled: zod_1.z.boolean().default(false),
});
exports.ConnectionStatusSchema = zod_1.z.object({
    connected: zod_1.z.boolean(),
    lastConnected: zod_1.z.string().datetime().optional(),
    lastError: zod_1.z.string().optional(),
    errorCount: zod_1.z.number().int().min(0).default(0),
    messageCount: zod_1.z.number().int().min(0).default(0),
});
exports.ChannelConfigSchema = zod_1.z.object({
    id: zod_1.z.string().min(1).max(64).regex(/^[a-z0-9_-]+$/),
    name: zod_1.z.string().min(1).max(128),
    type: exports.ChannelType,
    enabled: zod_1.z.boolean().default(true),
    config: zod_1.z.union([
        exports.TelegramConfigSchema,
        exports.DiscordConfigSchema,
        exports.SlackConfigSchema,
        exports.EmailConfigSchema,
        exports.WhatsAppConfigSchema,
        exports.WebhookConfigSchema,
    ]),
    settings: exports.ChannelSettingsSchema.default({}),
    webhook: exports.WebhookSettingsSchema.optional(),
    status: exports.ConnectionStatusSchema.optional(),
});
exports.ChannelsConfigSchema = zod_1.z.object({
    version: zod_1.z.number().int().default(1),
    channels: zod_1.z.array(exports.ChannelConfigSchema),
});
exports.ValidationResultSchema = zod_1.z.object({
    valid: zod_1.z.boolean(),
    errors: zod_1.z.array(zod_1.z.string()),
});
exports.ConnectionTestResultSchema = zod_1.z.object({
    success: zod_1.z.boolean(),
    latency: zod_1.z.number().optional(),
    error: zod_1.z.string().optional(),
});
var ChannelError = /** @class */ (function (_super) {
    __extends(ChannelError, _super);
    function ChannelError(message) {
        var _this = _super.call(this, message) || this;
        _this.name = 'ChannelError';
        return _this;
    }
    return ChannelError;
}(Error));
exports.ChannelError = ChannelError;
var ChannelNotFoundError = /** @class */ (function (_super) {
    __extends(ChannelNotFoundError, _super);
    function ChannelNotFoundError(id) {
        var _this = _super.call(this, "Channel not found: ".concat(id)) || this;
        _this.name = 'ChannelNotFoundError';
        return _this;
    }
    return ChannelNotFoundError;
}(ChannelError));
exports.ChannelNotFoundError = ChannelNotFoundError;
var ChannelValidationError = /** @class */ (function (_super) {
    __extends(ChannelValidationError, _super);
    function ChannelValidationError(errors) {
        var _this = _super.call(this, "Validation failed: ".concat(errors.join(', '))) || this;
        _this.errors = errors;
        _this.name = 'ChannelValidationError';
        return _this;
    }
    return ChannelValidationError;
}(ChannelError));
exports.ChannelValidationError = ChannelValidationError;
var ChannelAlreadyExistsError = /** @class */ (function (_super) {
    __extends(ChannelAlreadyExistsError, _super);
    function ChannelAlreadyExistsError(id) {
        var _this = _super.call(this, "Channel already exists: ".concat(id)) || this;
        _this.name = 'ChannelAlreadyExistsError';
        return _this;
    }
    return ChannelAlreadyExistsError;
}(ChannelError));
exports.ChannelAlreadyExistsError = ChannelAlreadyExistsError;
exports.CURRENT_VERSION = 1;
