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
(0, vitest_1.describe)('validateProviderConfig', function () {
    var validConfig = {
        id: 'openai',
        name: 'OpenAI',
        type: 'openai',
        enabled: true,
        apiKey: 'sk-test',
        models: [{
                id: 'gpt-4',
                name: 'GPT-4',
                maxTokens: 8192,
                supportsStreaming: true,
                supportsVision: true,
                supportsTools: true,
            }],
    };
    (0, vitest_1.it)('should validate correct config', function () {
        var result = (0, validation_js_1.validateProviderConfig)(validConfig);
        (0, vitest_1.expect)(result.valid).toBe(true);
        (0, vitest_1.expect)(result.errors).toHaveLength(0);
    });
    (0, vitest_1.it)('should reject config with invalid id format', function () {
        var config = __assign(__assign({}, validConfig), { id: 'Invalid ID!' });
        var result = (0, validation_js_1.validateProviderConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
        (0, vitest_1.expect)(result.errors.length).toBeGreaterThan(0);
    });
    (0, vitest_1.it)('should reject config with empty name', function () {
        var config = __assign(__assign({}, validConfig), { name: '' });
        var result = (0, validation_js_1.validateProviderConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
    });
    (0, vitest_1.it)('should reject config with invalid type', function () {
        var config = __assign(__assign({}, validConfig), { type: 'invalid' });
        var result = (0, validation_js_1.validateProviderConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
    });
    (0, vitest_1.it)('should reject config with empty models array', function () {
        var config = __assign(__assign({}, validConfig), { models: [] });
        var result = (0, validation_js_1.validateProviderConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
    });
    (0, vitest_1.it)('should reject custom provider without baseUrl', function () {
        var config = __assign(__assign({}, validConfig), { type: 'custom', baseUrl: undefined });
        var result = (0, validation_js_1.validateProviderConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
        (0, vitest_1.expect)(result.errors.some(function (e) { return e.includes('baseUrl'); })).toBe(true);
    });
    (0, vitest_1.it)('should reject local provider with apiKey', function () {
        var config = __assign(__assign({}, validConfig), { type: 'local', apiKey: 'should-not-have' });
        var result = (0, validation_js_1.validateProviderConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
        (0, vitest_1.expect)(result.errors.some(function (e) { return e.includes('API key'); })).toBe(true);
    });
    (0, vitest_1.it)('should reject config with invalid defaultModel', function () {
        var config = __assign(__assign({}, validConfig), { defaultModel: 'nonexistent-model' });
        var result = (0, validation_js_1.validateProviderConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
        (0, vitest_1.expect)(result.errors.some(function (e) { return e.includes('Default model'); })).toBe(true);
    });
    (0, vitest_1.it)('should accept config with valid defaultModel', function () {
        var config = __assign(__assign({}, validConfig), { defaultModel: 'gpt-4' });
        var result = (0, validation_js_1.validateProviderConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(true);
    });
    (0, vitest_1.it)('should reject config with invalid baseUrl', function () {
        var config = __assign(__assign({}, validConfig), { baseUrl: 'not-a-url' });
        var result = (0, validation_js_1.validateProviderConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(false);
    });
    (0, vitest_1.it)('should accept config with valid baseUrl', function () {
        var config = __assign(__assign({}, validConfig), { baseUrl: 'https://api.example.com' });
        var result = (0, validation_js_1.validateProviderConfig)(config);
        (0, vitest_1.expect)(result.valid).toBe(true);
    });
});
(0, vitest_1.describe)('assertValidProviderConfig', function () {
    var validConfig = {
        id: 'test',
        name: 'Test',
        type: 'openai',
        enabled: true,
        apiKey: 'key',
        models: [{
                id: 'model',
                name: 'Model',
                maxTokens: 1000,
                supportsStreaming: true,
                supportsVision: false,
                supportsTools: false,
            }],
    };
    (0, vitest_1.it)('should return config when valid', function () {
        var result = (0, validation_js_1.assertValidProviderConfig)(validConfig);
        (0, vitest_1.expect)(result.id).toBe('test');
    });
    (0, vitest_1.it)('should throw when invalid', function () {
        var config = __assign(__assign({}, validConfig), { id: '' });
        (0, vitest_1.expect)(function () { return (0, validation_js_1.assertValidProviderConfig)(config); }).toThrow(types_js_1.ProviderValidationError);
    });
});
(0, vitest_1.describe)('isValidProviderId', function () {
    (0, vitest_1.it)('should return true for valid ids', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidProviderId)('openai')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidProviderId)('anthropic-v2')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidProviderId)('custom_provider')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidProviderId)('a')).toBe(true);
    });
    (0, vitest_1.it)('should return false for invalid ids', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidProviderId)('')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidProviderId)('Invalid ID')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidProviderId)('test@provider')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidProviderId)('a'.repeat(65))).toBe(false);
    });
});
(0, vitest_1.describe)('isValidProviderType', function () {
    (0, vitest_1.it)('should return true for valid types', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidProviderType)('openai')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidProviderType)('anthropic')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidProviderType)('google')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidProviderType)('local')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidProviderType)('custom')).toBe(true);
    });
    (0, vitest_1.it)('should return false for invalid types', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidProviderType)('invalid')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidProviderType)('azure')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidProviderType)('')).toBe(false);
    });
});
(0, vitest_1.describe)('isValidUrl', function () {
    (0, vitest_1.it)('should return true for valid URLs', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidUrl)('https://api.openai.com')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidUrl)('http://localhost:3000')).toBe(true);
        (0, vitest_1.expect)((0, validation_js_1.isValidUrl)('https://example.com/path')).toBe(true);
    });
    (0, vitest_1.it)('should return false for invalid URLs', function () {
        (0, vitest_1.expect)((0, validation_js_1.isValidUrl)('not-a-url')).toBe(false);
        (0, vitest_1.expect)((0, validation_js_1.isValidUrl)('')).toBe(false);
    });
});
