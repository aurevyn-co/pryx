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
(0, vitest_1.describe)('validateChannelConfig', function () {
    var baseConfig = {
        id: 'test-channel',
        name: 'Test Channel',
        type: 'telegram',
        enabled: true,
        config: {
            botToken: 'test-token',
        },
    };
    (0, vitest_1.it)('should validate correct telegram config', function () {
        var result = (0, validation_js_1.validateChannelConfig)(baseConfig);
        (0, vitest_1.expect)(result.valid).toBe(true);
        (0, vitest_1.expect)(result.errors).toHaveLength(0);
    });
    (0, vitest_1.it)('should validate correct discord config', function () {
        var config = __assign(__assign({}, baseConfig), { type: 'discord', config: {
                botToken: 'discord-token',
                applicationId: 'app-id',
            } });
        var result = (0, validation_js_1.validateChannelConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(true);
    });
    (0, vitest_1.it)('should validate correct slack config', function () {
        var config = __assign(__assign({}, baseConfig), { type: 'slack', config: {
                botToken: 'slack-token',
            } });
        var result = (0, validation_js_1.validateChannelConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(true);
    });
    (0, vitest_1.it)('should validate correct email config with imap', function () {
        var config = __assign(__assign({}, baseConfig), { type: 'email', config: {
                imap: {
                    host: 'imap.example.com',
                    port: 993,
                    secure: true,
                    username: 'user',
                    password: 'pass',
                },
            } });
        var result = (0, validation_js_1.validateChannelConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(true);
    });
    (0, vitest_1.it)('should validate correct email config with smtp', function () {
        var config = __assign(__assign({}, baseConfig), { type: 'email', config: {
                smtp: {
                    host: 'smtp.example.com',
                    port: 587,
                    secure: true,
                    username: 'user',
                    password: 'pass',
                },
            } });
        var result = (0, validation_js_1.validateChannelConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(true);
    });
    (0, vitest_1.it)('should validate correct whatsapp config', function () {
        var config = __assign(__assign({}, baseConfig), { type: 'whatsapp', config: {
                sessionName: 'my-session',
            } });
        var result = (0, validation_js_1.validateChannelConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(true);
    });
    (0, vitest_1.it)('should validate correct webhook config', function () {
        var config = __assign(__assign({}, baseConfig), { type: 'webhook', config: {
                url: 'https://api.example.com/webhook',
                method: 'POST',
            } });
        var result = (0, validation_js_1.validateChannelConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(true);
    });
    (0, vitest_1.it)('should reject config with invalid id format', function () {
        var config = __assign(__assign({}, baseConfig), { id: 'Invalid ID!' });
        var result = (0, validation_js_1.validateChannelConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
    });
    (0, vitest_1.it)('should reject config with empty name', function () {
        var config = __assign(__assign({}, baseConfig), { name: '' });
        var result = (0, validation_js_1.validateChannelConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
    });
    (0, vitest_1.it)('should reject config with invalid type', function () {
        var config = __assign(__assign({}, baseConfig), { type: 'invalid' });
        var result = (0, validation_js_1.validateChannelConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
    });
    (0, vitest_1.it)('should reject webhook settings with enabled but no url', function () {
        var config = __assign(__assign({}, baseConfig), { webhook: {
                enabled: true,
            } });
        var result = (0, validation_js_1.validateChannelConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
    });
});
(0, vitest_1.describe)('assertValidChannelConfig', function () {
    var validConfig = {
        id: 'test',
        name: 'Test',
        type: 'telegram',
        enabled: true,
        config: { botToken: 'token' },
    };
    (0, vitest_1.it)('should return config when valid', function () {
        var result = (0, validation_js_1.assertValidChannelConfig)(validConfig);
        (0, vitest_1.expect)(result.id).toBe('test');
    });
    (0, vitest_1.it)('should throw when invalid', function () {
        var config = __assign(__assign({}, validConfig), { id: '' });
        (0, vitest_1.expect)(function () { return (0, validation_js_1.assertValidChannelConfig)(config); }).toThrow(types_js_1.ChannelValidationError);
    });
});
(0, vitest_1.describe)('isValidChannelId', function () {
    (0, vitest_1.it)('should return true for valid ids', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidChannelId)('telegram-bot')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidChannelId)('discord-server')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidChannelId)('a')).toBe(true);
    });
    (0, vitest_1.it)('should return false for invalid ids', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidChannelId)('')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidChannelId)('Invalid ID')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidChannelId)('test@channel')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidChannelId)('a'.repeat(65))).toBe(false);
    });
});
(0, vitest_1.describe)('isValidChannelType', function () {
    (0, vitest_1.it)('should return true for valid types', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidChannelType)('telegram')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidChannelType)('discord')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidChannelType)('slack')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidChannelType)('email')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidChannelType)('whatsapp')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidChannelType)('webhook')).toBe(true);
    });
    (0, vitest_1.it)('should return false for invalid types', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidChannelType)('invalid')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidChannelType)('signal')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidChannelType)('')).toBe(false);
    });
});
(0, vitest_1.describe)('matchesFilterPatterns', function () {
    (0, vitest_1.it)('should return true when no patterns', function () {
        (0, vitest_1.expect)((0, validation_js_1.matchesFilterPatterns)('hello', [])).toBe(true);
    });
    (0, vitest_1.it)('should match regex patterns', function () {
        (0, vitest_1.expect)((0, validation_js_1.matchesFilterPatterns)('hello world', ['hello'])).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.matchesFilterPatterns)('hello world', ['^hello'])).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.matchesFilterPatterns)('HELLO world', ['hello'])).toBe(true);
    });
    (0, vitest_1.it)('should match any pattern', function () {
        (0, vitest_1.expect)((0, validation_js_1.matchesFilterPatterns)('test', ['hello', 'test', 'world'])).toBe(true);
    });
    (0, vitest_1.it)('should return false when no match', function () {
        (0, vitest_1.expect)((0, validation_js_1.matchesFilterPatterns)('goodbye', ['hello'])).toBe(false);
    });
    (0, vitest_1.it)('should handle invalid regex gracefully', function () {
        (0, vitest_1.expect)((0, validation_js_1.matchesFilterPatterns)('hello', ['[invalid'])).toBe(false);
    });
});
(0, vitest_1.describe)('isUserAllowed', function () {
    (0, vitest_1.it)('should allow user when no restrictions', function () {
        (0, vitest_1.expect)((0, validation_js_1.isUserAllowed)('user1', [], [])).toBe(true);
    });
    (0, vitest_1.it)('should block user in blocked list', function () {
        (0, vitest_1.expect)((0, validation_js_1.isUserAllowed)('user1', [], ['user1'])).toBe(false);
    });
    (0, vitest_1.it)('should allow user in allowed list', function () {
        (0, vitest_1.expect)((0, validation_js_1.isUserAllowed)('user1', ['user1'], [])).toBe(true);
    });
    (0, vitest_1.it)('should block user not in allowed list', function () {
        (0, vitest_1.expect)((0, validation_js_1.isUserAllowed)('user2', ['user1'], [])).toBe(false);
    });
    (0, vitest_1.it)('should prioritize blocked over allowed', function () {
        (0, vitest_1.expect)((0, validation_js_1.isUserAllowed)('user1', ['user1'], ['user1'])).toBe(false);
    });
});
