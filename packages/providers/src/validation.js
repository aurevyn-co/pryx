"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.validateProviderConfig = validateProviderConfig;
exports.assertValidProviderConfig = assertValidProviderConfig;
exports.isValidProviderId = isValidProviderId;
exports.isValidProviderType = isValidProviderType;
exports.isValidUrl = isValidUrl;
var types_js_1 = require("./types.js");
function validateProviderConfig(config) {
    var result = types_js_1.ProviderConfigSchema.safeParse(config);
    if (result.success) {
        var errors = [];
        var validated_1 = result.data;
        if (validated_1.defaultModel) {
            var modelExists = validated_1.models.some(function (m) { return m.id === validated_1.defaultModel; });
            if (!modelExists) {
                errors.push("Default model \"".concat(validated_1.defaultModel, "\" not found in models list"));
            }
        }
        if (validated_1.type === 'custom' && !validated_1.baseUrl) {
            errors.push('Custom providers require a baseUrl');
        }
        if (validated_1.type === 'local' && validated_1.apiKey) {
            errors.push('Local providers should not have an API key');
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
function assertValidProviderConfig(config) {
    var result = validateProviderConfig(config);
    if (!result.valid) {
        throw new types_js_1.ProviderValidationError(result.errors);
    }
    return config;
}
function isValidProviderId(id) {
    return /^[a-z0-9_-]+$/.test(id) && id.length >= 1 && id.length <= 64;
}
function isValidProviderType(type) {
    return ['openai', 'anthropic', 'google', 'local', 'custom'].includes(type);
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
