"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var vitest_1 = require("vitest");
var presets_js_1 = require("../../src/presets.js");
(0, vitest_1.describe)('OPENAI_MODELS', function () {
    (0, vitest_1.it)('should contain GPT-4 models', function () {
        (0, vitest_1.expect)(presets_js_1.OPENAI_MODELS.some(function (m) { return m.id === 'gpt-4o'; })).toBe(true);
        (0, vitest_1.expect)(presets_js_1.OPENAI_MODELS.some(function (m) { return m.id === 'gpt-4o-mini'; })).toBe(true);
    });
    (0, vitest_1.it)('should have correct properties', function () {
        var gpt4o = presets_js_1.OPENAI_MODELS.find(function (m) { return m.id === 'gpt-4o'; });
        (0, vitest_1.expect)(gpt4o === null || gpt4o === void 0 ? void 0 : gpt4o.maxTokens).toBe(128000);
        (0, vitest_1.expect)(gpt4o === null || gpt4o === void 0 ? void 0 : gpt4o.supportsStreaming).toBe(true);
        (0, vitest_1.expect)(gpt4o === null || gpt4o === void 0 ? void 0 : gpt4o.supportsVision).toBe(true);
        (0, vitest_1.expect)(gpt4o === null || gpt4o === void 0 ? void 0 : gpt4o.supportsTools).toBe(true);
    });
});
(0, vitest_1.describe)('ANTHROPIC_MODELS', function () {
    (0, vitest_1.it)('should contain Claude 3 models', function () {
        (0, vitest_1.expect)(presets_js_1.ANTHROPIC_MODELS.some(function (m) { return m.id === 'claude-3-opus-20240229'; })).toBe(true);
        (0, vitest_1.expect)(presets_js_1.ANTHROPIC_MODELS.some(function (m) { return m.id === 'claude-3-sonnet-20240229'; })).toBe(true);
    });
});
(0, vitest_1.describe)('GOOGLE_MODELS', function () {
    (0, vitest_1.it)('should contain Gemini models', function () {
        (0, vitest_1.expect)(presets_js_1.GOOGLE_MODELS.some(function (m) { return m.id === 'gemini-1.5-pro'; })).toBe(true);
    });
});
(0, vitest_1.describe)('OPENAI_PRESET', function () {
    (0, vitest_1.it)('should have correct configuration', function () {
        (0, vitest_1.expect)(presets_js_1.OPENAI_PRESET.id).toBe('openai');
        (0, vitest_1.expect)(presets_js_1.OPENAI_PRESET.type).toBe('openai');
        (0, vitest_1.expect)(presets_js_1.OPENAI_PRESET.enabled).toBe(true);
        (0, vitest_1.expect)(presets_js_1.OPENAI_PRESET.defaultModel).toBe('gpt-4o');
    });
});
(0, vitest_1.describe)('ANTHROPIC_PRESET', function () {
    (0, vitest_1.it)('should have correct configuration', function () {
        (0, vitest_1.expect)(presets_js_1.ANTHROPIC_PRESET.id).toBe('anthropic');
        (0, vitest_1.expect)(presets_js_1.ANTHROPIC_PRESET.type).toBe('anthropic');
        (0, vitest_1.expect)(presets_js_1.ANTHROPIC_PRESET.defaultModel).toBe('claude-3-sonnet-20240229');
    });
});
(0, vitest_1.describe)('OLLAMA_PRESET', function () {
    (0, vitest_1.it)('should have correct configuration', function () {
        (0, vitest_1.expect)(presets_js_1.OLLAMA_PRESET.id).toBe('ollama');
        (0, vitest_1.expect)(presets_js_1.OLLAMA_PRESET.type).toBe('local');
        (0, vitest_1.expect)(presets_js_1.OLLAMA_PRESET.enabled).toBe(false);
        (0, vitest_1.expect)(presets_js_1.OLLAMA_PRESET.baseUrl).toBe('http://localhost:11434');
    });
});
(0, vitest_1.describe)('BUILTIN_PRESETS', function () {
    (0, vitest_1.it)('should contain all presets', function () {
        (0, vitest_1.expect)(presets_js_1.BUILTIN_PRESETS.length).toBeGreaterThanOrEqual(5);
        (0, vitest_1.expect)(presets_js_1.BUILTIN_PRESETS.some(function (p) { return p.id === 'openai'; })).toBe(true);
        (0, vitest_1.expect)(presets_js_1.BUILTIN_PRESETS.some(function (p) { return p.id === 'anthropic'; })).toBe(true);
    });
});
(0, vitest_1.describe)('getPreset', function () {
    (0, vitest_1.it)('should return preset by id', function () {
        var preset = (0, presets_js_1.getPreset)('openai');
        (0, vitest_1.expect)(preset).toBeDefined();
        (0, vitest_1.expect)(preset === null || preset === void 0 ? void 0 : preset.id).toBe('openai');
    });
    (0, vitest_1.it)('should return undefined for unknown preset', function () {
        var preset = (0, presets_js_1.getPreset)('unknown');
        (0, vitest_1.expect)(preset).toBeUndefined();
    });
});
(0, vitest_1.describe)('getAllPresets', function () {
    (0, vitest_1.it)('should return all presets', function () {
        var presets = (0, presets_js_1.getAllPresets)();
        (0, vitest_1.expect)(presets.length).toBe(presets_js_1.BUILTIN_PRESETS.length);
    });
});
(0, vitest_1.describe)('getPresetIds', function () {
    (0, vitest_1.it)('should return all preset ids', function () {
        var ids = (0, presets_js_1.getPresetIds)();
        (0, vitest_1.expect)(ids).toContain('openai');
        (0, vitest_1.expect)(ids).toContain('anthropic');
    });
});
