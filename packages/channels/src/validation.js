"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.validateChannelConfig = validateChannelConfig;
exports.assertValidChannelConfig = assertValidChannelConfig;
exports.isValidChannelId = isValidChannelId;
exports.isValidChannelType = isValidChannelType;
exports.matchesFilterPatterns = matchesFilterPatterns;
exports.isUserAllowed = isUserAllowed;
var types_js_1 = require("./types.js");
function validateChannelConfig(config) {
    var _a;
    var result = types_js_1.ChannelConfigSchema.safeParse(config);
    if (result.success) {
        var errors = [];
        var validated = result.data;
        if (((_a = validated.webhook) === null || _a === void 0 ? void 0 : _a.enabled) && !validated.webhook.url) {
            errors.push('Webhook settings enabled but URL is missing');
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
function assertValidChannelConfig(config) {
    var result = validateChannelConfig(config);
    if (!result.valid) {
        throw new types_js_1.ChannelValidationError(result.errors);
    }
    return config;
}
function isValidChannelId(id) {
    return /^[a-z0-9_-]+$/.test(id) && id.length >= 1 && id.length <= 64;
}
function isValidChannelType(type) {
    return ['telegram', 'discord', 'slack', 'email', 'whatsapp', 'webhook'].includes(type);
}
function matchesFilterPatterns(message, patterns) {
    if (patterns.length === 0)
        return true;
    return patterns.some(function (pattern) {
        try {
            var regex = new RegExp(pattern, 'i');
            return regex.test(message);
        }
        catch (_a) {
            return message.toLowerCase().includes(pattern.toLowerCase());
        }
    });
}
function isUserAllowed(userId, allowedList, blockedList) {
    if (blockedList.includes(userId)) {
        return false;
    }
    if (allowedList.length > 0 && !allowedList.includes(userId)) {
        return false;
    }
    return true;
}
