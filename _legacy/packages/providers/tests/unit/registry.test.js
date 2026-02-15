"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g = Object.create((typeof Iterator === "function" ? Iterator : Object).prototype);
    return g.next = verb(0), g["throw"] = verb(1), g["return"] = verb(2), typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (g && (g = 0, op[0] && (_ = 0)), _) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
Object.defineProperty(exports, "__esModule", { value: true });
var vitest_1 = require("vitest");
var registry_js_1 = require("../../src/registry.js");
var types_js_1 = require("../../src/types.js");
(0, vitest_1.describe)('ProviderRegistry', function () {
    var registry;
    (0, vitest_1.beforeEach)(function () {
        registry = new registry_js_1.ProviderRegistry();
    });
    (0, vitest_1.describe)('constructor', function () {
        (0, vitest_1.it)('should initialize with builtin presets', function () {
            (0, vitest_1.expect)(registry.size).toBeGreaterThan(0);
            (0, vitest_1.expect)(registry.hasProvider('openai')).toBe(true);
            (0, vitest_1.expect)(registry.hasProvider('anthropic')).toBe(true);
        });
    });
    (0, vitest_1.describe)('addProvider', function () {
        (0, vitest_1.it)('should add new provider', function () {
            var _a;
            var newProvider = {
                id: 'custom',
                name: 'Custom Provider',
                type: 'custom',
                enabled: true,
                baseUrl: 'https://api.custom.com',
                apiKey: 'test-key',
                models: [{
                        id: 'model-1',
                        name: 'Model 1',
                        maxTokens: 4096,
                        supportsStreaming: true,
                        supportsVision: false,
                        supportsTools: false,
                    }],
                timeout: 30000,
                retries: 3,
            };
            registry.addProvider(newProvider);
            (0, vitest_1.expect)(registry.hasProvider('custom')).toBe(true);
            (0, vitest_1.expect)((_a = registry.getProvider('custom')) === null || _a === void 0 ? void 0 : _a.name).toBe('Custom Provider');
        });
        (0, vitest_1.it)('should throw when provider already exists', function () {
            var provider = {
                id: 'openai',
                name: 'Duplicate',
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
                timeout: 30000,
                retries: 3,
            };
            (0, vitest_1.expect)(function () { return registry.addProvider(provider); }).toThrow(types_js_1.ProviderAlreadyExistsError);
        });
        (0, vitest_1.it)('should throw on invalid config', function () {
            var invalidProvider = {
                id: 'invalid',
                name: '',
                type: 'openai',
                enabled: true,
                apiKey: 'key',
                models: [],
                timeout: 30000,
                retries: 3,
            };
            (0, vitest_1.expect)(function () { return registry.addProvider(invalidProvider); }).toThrow(types_js_1.ProviderValidationError);
        });
    });
    (0, vitest_1.describe)('updateProvider', function () {
        (0, vitest_1.it)('should update existing provider', function () {
            var _a;
            var updated = registry.updateProvider('openai', { name: 'Updated OpenAI' });
            (0, vitest_1.expect)(updated.name).toBe('Updated OpenAI');
            (0, vitest_1.expect)((_a = registry.getProvider('openai')) === null || _a === void 0 ? void 0 : _a.name).toBe('Updated OpenAI');
        });
        (0, vitest_1.it)('should throw when provider not found', function () {
            (0, vitest_1.expect)(function () { return registry.updateProvider('nonexistent', { name: 'Test' }); }).toThrow(types_js_1.ProviderNotFoundError);
        });
    });
    (0, vitest_1.describe)('removeProvider', function () {
        (0, vitest_1.it)('should remove provider', function () {
            registry.removeProvider('openai');
            (0, vitest_1.expect)(registry.hasProvider('openai')).toBe(false);
        });
        (0, vitest_1.it)('should throw when provider not found', function () {
            (0, vitest_1.expect)(function () { return registry.removeProvider('nonexistent'); }).toThrow(types_js_1.ProviderNotFoundError);
        });
        (0, vitest_1.it)('should clear default provider when removed', function () {
            registry.setDefaultProvider('openai');
            registry.removeProvider('openai');
            (0, vitest_1.expect)(registry.getDefaultProviderId()).toBeNull();
        });
    });
    (0, vitest_1.describe)('getProvider', function () {
        (0, vitest_1.it)('should return provider by id', function () {
            var provider = registry.getProvider('openai');
            (0, vitest_1.expect)(provider).toBeDefined();
            (0, vitest_1.expect)(provider === null || provider === void 0 ? void 0 : provider.id).toBe('openai');
        });
        (0, vitest_1.it)('should return undefined for nonexistent provider', function () {
            var provider = registry.getProvider('nonexistent');
            (0, vitest_1.expect)(provider).toBeUndefined();
        });
    });
    (0, vitest_1.describe)('getAllProviders', function () {
        (0, vitest_1.it)('should return all providers', function () {
            var providers = registry.getAllProviders();
            (0, vitest_1.expect)(providers.length).toBe(registry.size);
            (0, vitest_1.expect)(providers.some(function (p) { return p.id === 'openai'; })).toBe(true);
        });
    });
    (0, vitest_1.describe)('getEnabledProviders', function () {
        (0, vitest_1.it)('should return only enabled providers', function () {
            registry.disableProvider('openai');
            var enabled = registry.getEnabledProviders();
            (0, vitest_1.expect)(enabled.some(function (p) { return p.id === 'openai'; })).toBe(false);
            (0, vitest_1.expect)(enabled.every(function (p) { return p.enabled; })).toBe(true);
        });
    });
    (0, vitest_1.describe)('setDefaultProvider', function () {
        (0, vitest_1.it)('should set default provider', function () {
            registry.setDefaultProvider('anthropic');
            (0, vitest_1.expect)(registry.getDefaultProviderId()).toBe('anthropic');
        });
        (0, vitest_1.it)('should throw when provider not found', function () {
            (0, vitest_1.expect)(function () { return registry.setDefaultProvider('nonexistent'); }).toThrow(types_js_1.ProviderNotFoundError);
        });
    });
    (0, vitest_1.describe)('getDefaultProvider', function () {
        (0, vitest_1.it)('should return explicitly set default', function () {
            registry.setDefaultProvider('anthropic');
            var provider = registry.getDefaultProvider();
            (0, vitest_1.expect)(provider === null || provider === void 0 ? void 0 : provider.id).toBe('anthropic');
        });
        (0, vitest_1.it)('should return first enabled provider when no default set', function () {
            var provider = registry.getDefaultProvider();
            (0, vitest_1.expect)(provider).toBeDefined();
            (0, vitest_1.expect)(provider === null || provider === void 0 ? void 0 : provider.enabled).toBe(true);
        });
        (0, vitest_1.it)('should return undefined when no providers', function () {
            registry.clear();
            var provider = registry.getDefaultProvider();
            (0, vitest_1.expect)(provider).toBeUndefined();
        });
    });
    (0, vitest_1.describe)('enableProvider / disableProvider', function () {
        (0, vitest_1.it)('should disable provider', function () {
            var _a;
            registry.disableProvider('openai');
            (0, vitest_1.expect)((_a = registry.getProvider('openai')) === null || _a === void 0 ? void 0 : _a.enabled).toBe(false);
        });
        (0, vitest_1.it)('should enable provider', function () {
            var _a;
            registry.disableProvider('openai');
            registry.enableProvider('openai');
            (0, vitest_1.expect)((_a = registry.getProvider('openai')) === null || _a === void 0 ? void 0 : _a.enabled).toBe(true);
        });
    });
    (0, vitest_1.describe)('validateProvider', function () {
        (0, vitest_1.it)('should validate existing provider', function () {
            var result = registry.validateProvider('openai');
            (0, vitest_1.expect)(result.valid).toBe(true);
            (0, vitest_1.expect)(result.errors).toHaveLength(0);
        });
        (0, vitest_1.it)('should return error for nonexistent provider', function () {
            var result = registry.validateProvider('nonexistent');
            (0, vitest_1.expect)(result.valid).toBe(false);
            (0, vitest_1.expect)(result.errors[0]).toContain('not found');
        });
    });
    (0, vitest_1.describe)('testConnection', function () {
        (0, vitest_1.it)('should fail for nonexistent provider', function () { return __awaiter(void 0, void 0, void 0, function () {
            var result;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, registry.testConnection('nonexistent')];
                    case 1:
                        result = _a.sent();
                        (0, vitest_1.expect)(result.success).toBe(false);
                        (0, vitest_1.expect)(result.error).toContain('not found');
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should fail when api key not configured', function () { return __awaiter(void 0, void 0, void 0, function () {
            var result;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, registry.testConnection('openai')];
                    case 1:
                        result = _a.sent();
                        (0, vitest_1.expect)(result.success).toBe(false);
                        (0, vitest_1.expect)(result.error).toContain('API key');
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('toJSON / fromJSON', function () {
        (0, vitest_1.it)('should serialize to JSON', function () {
            var json = registry.toJSON();
            (0, vitest_1.expect)(json.version).toBe(1);
            (0, vitest_1.expect)(json.providers.length).toBe(registry.size);
        });
        (0, vitest_1.it)('should deserialize from JSON', function () {
            var json = registry.toJSON();
            var newRegistry = new registry_js_1.ProviderRegistry();
            newRegistry.clear();
            newRegistry.fromJSON(json);
            (0, vitest_1.expect)(newRegistry.size).toBe(registry.size);
            (0, vitest_1.expect)(newRegistry.hasProvider('openai')).toBe(true);
        });
        (0, vitest_1.it)('should throw on unsupported version', function () {
            var json = { version: 999, providers: [] };
            (0, vitest_1.expect)(function () { return registry.fromJSON(json); }).toThrow(types_js_1.ProviderValidationError);
        });
    });
    (0, vitest_1.describe)('clear', function () {
        (0, vitest_1.it)('should clear all providers', function () {
            registry.clear();
            (0, vitest_1.expect)(registry.size).toBe(0);
            (0, vitest_1.expect)(registry.getDefaultProviderId()).toBeNull();
        });
    });
});
(0, vitest_1.describe)('createRegistry', function () {
    (0, vitest_1.it)('should create new registry', function () {
        var registry = (0, registry_js_1.createRegistry)();
        (0, vitest_1.expect)(registry).toBeInstanceOf(registry_js_1.ProviderRegistry);
        (0, vitest_1.expect)(registry.size).toBeGreaterThan(0);
    });
});
