"use strict";
var __spreadArray = (this && this.__spreadArray) || function (to, from, pack) {
    if (pack || arguments.length === 2) for (var i = 0, l = from.length, ar; i < l; i++) {
        if (ar || !(i in from)) {
            if (!ar) ar = Array.prototype.slice.call(from, 0, i);
            ar[i] = from[i];
        }
    }
    return to.concat(ar || Array.prototype.slice.call(from));
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.BUILTIN_PRESETS = exports.LMSTUDIO_PRESET = exports.OLLAMA_PRESET = exports.GOOGLE_PRESET = exports.ANTHROPIC_PRESET = exports.OPENAI_PRESET = exports.LOCAL_MODELS = exports.GOOGLE_MODELS = exports.ANTHROPIC_MODELS = exports.OPENAI_MODELS = void 0;
exports.getPreset = getPreset;
exports.getAllPresets = getAllPresets;
exports.getPresetIds = getPresetIds;
exports.OPENAI_MODELS = [
    {
        id: 'gpt-4o',
        name: 'GPT-4o',
        maxTokens: 128000,
        supportsStreaming: true,
        supportsVision: true,
        supportsTools: true,
        costPer1KInput: 0.005,
        costPer1KOutput: 0.015,
    },
    {
        id: 'gpt-4o-mini',
        name: 'GPT-4o Mini',
        maxTokens: 128000,
        supportsStreaming: true,
        supportsVision: true,
        supportsTools: true,
        costPer1KInput: 0.00015,
        costPer1KOutput: 0.0006,
    },
    {
        id: 'gpt-4-turbo',
        name: 'GPT-4 Turbo',
        maxTokens: 128000,
        supportsStreaming: true,
        supportsVision: true,
        supportsTools: true,
        costPer1KInput: 0.01,
        costPer1KOutput: 0.03,
    },
    {
        id: 'gpt-3.5-turbo',
        name: 'GPT-3.5 Turbo',
        maxTokens: 16385,
        supportsStreaming: true,
        supportsVision: false,
        supportsTools: true,
        costPer1KInput: 0.0005,
        costPer1KOutput: 0.0015,
    },
];
exports.ANTHROPIC_MODELS = [
    {
        id: 'claude-3-opus-20240229',
        name: 'Claude 3 Opus',
        maxTokens: 200000,
        supportsStreaming: true,
        supportsVision: true,
        supportsTools: true,
        costPer1KInput: 0.015,
        costPer1KOutput: 0.075,
    },
    {
        id: 'claude-3-sonnet-20240229',
        name: 'Claude 3 Sonnet',
        maxTokens: 200000,
        supportsStreaming: true,
        supportsVision: true,
        supportsTools: true,
        costPer1KInput: 0.003,
        costPer1KOutput: 0.015,
    },
    {
        id: 'claude-3-haiku-20240307',
        name: 'Claude 3 Haiku',
        maxTokens: 200000,
        supportsStreaming: true,
        supportsVision: true,
        supportsTools: true,
        costPer1KInput: 0.00025,
        costPer1KOutput: 0.00125,
    },
];
exports.GOOGLE_MODELS = [
    {
        id: 'gemini-1.5-pro',
        name: 'Gemini 1.5 Pro',
        maxTokens: 1048576,
        supportsStreaming: true,
        supportsVision: true,
        supportsTools: true,
        costPer1KInput: 0.0035,
        costPer1KOutput: 0.0105,
    },
    {
        id: 'gemini-1.5-flash',
        name: 'Gemini 1.5 Flash',
        maxTokens: 1048576,
        supportsStreaming: true,
        supportsVision: true,
        supportsTools: true,
        costPer1KInput: 0.00035,
        costPer1KOutput: 0.00105,
    },
];
exports.LOCAL_MODELS = [
    {
        id: 'llama2',
        name: 'Llama 2',
        maxTokens: 4096,
        supportsStreaming: true,
        supportsVision: false,
        supportsTools: false,
    },
    {
        id: 'codellama',
        name: 'CodeLlama',
        maxTokens: 16384,
        supportsStreaming: true,
        supportsVision: false,
        supportsTools: false,
    },
    {
        id: 'mistral',
        name: 'Mistral',
        maxTokens: 8192,
        supportsStreaming: true,
        supportsVision: false,
        supportsTools: false,
    },
];
exports.OPENAI_PRESET = {
    id: 'openai',
    name: 'OpenAI',
    type: 'openai',
    enabled: true,
    defaultModel: 'gpt-4o',
    models: exports.OPENAI_MODELS,
    timeout: 30000,
    retries: 3,
};
exports.ANTHROPIC_PRESET = {
    id: 'anthropic',
    name: 'Anthropic',
    type: 'anthropic',
    enabled: true,
    defaultModel: 'claude-3-sonnet-20240229',
    models: exports.ANTHROPIC_MODELS,
    timeout: 30000,
    retries: 3,
};
exports.GOOGLE_PRESET = {
    id: 'google',
    name: 'Google',
    type: 'google',
    enabled: true,
    defaultModel: 'gemini-1.5-pro',
    models: exports.GOOGLE_MODELS,
    timeout: 30000,
    retries: 3,
};
exports.OLLAMA_PRESET = {
    id: 'ollama',
    name: 'Ollama (Local)',
    type: 'local',
    enabled: false,
    defaultModel: 'llama2',
    baseUrl: 'http://localhost:11434',
    models: exports.LOCAL_MODELS,
    timeout: 60000,
    retries: 1,
};
exports.LMSTUDIO_PRESET = {
    id: 'lmstudio',
    name: 'LM Studio (Local)',
    type: 'local',
    enabled: false,
    baseUrl: 'http://localhost:1234',
    models: exports.LOCAL_MODELS,
    timeout: 60000,
    retries: 1,
};
exports.BUILTIN_PRESETS = [
    exports.OPENAI_PRESET,
    exports.ANTHROPIC_PRESET,
    exports.GOOGLE_PRESET,
    exports.OLLAMA_PRESET,
    exports.LMSTUDIO_PRESET,
];
function getPreset(id) {
    return exports.BUILTIN_PRESETS.find(function (p) { return p.id === id; });
}
function getAllPresets() {
    return __spreadArray([], exports.BUILTIN_PRESETS, true);
}
function getPresetIds() {
    return exports.BUILTIN_PRESETS.map(function (p) { return p.id; });
}
