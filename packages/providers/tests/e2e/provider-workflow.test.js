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
var promises_1 = require("fs/promises");
var os_1 = require("os");
var path_1 = require("path");
var index_js_1 = require("../../src/index.js");
(0, vitest_1.describe)('Provider Workflow E2E', function () {
    var tempDir;
    var configPath;
    (0, vitest_1.beforeEach)(function () { return __awaiter(void 0, void 0, void 0, function () {
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, (0, promises_1.mkdtemp)((0, path_1.join)((0, os_1.tmpdir)(), 'providers-e2e-test-'))];
                case 1:
                    tempDir = _a.sent();
                    configPath = (0, path_1.join)(tempDir, 'providers.json');
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.afterEach)(function () { return __awaiter(void 0, void 0, void 0, function () {
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0: return [4 /*yield*/, (0, promises_1.rm)(tempDir, { recursive: true, force: true })];
                case 1:
                    _a.sent();
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should add and configure OpenAI provider', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry, storage, provider, _a;
        return __generator(this, function (_b) {
            switch (_b.label) {
                case 0:
                    registry = (0, index_js_1.createRegistry)();
                    storage = (0, index_js_1.createStorage)();
                    registry.updateProvider('openai', {
                        apiKey: 'sk-openai-test-key',
                        enabled: true,
                    });
                    provider = registry.getProvider('openai');
                    (0, vitest_1.expect)(provider === null || provider === void 0 ? void 0 : provider.apiKey).toBe('sk-openai-test-key');
                    (0, vitest_1.expect)(provider === null || provider === void 0 ? void 0 : provider.enabled).toBe(true);
                    (0, vitest_1.expect)(provider === null || provider === void 0 ? void 0 : provider.models.some(function (m) { return m.id === 'gpt-4o'; })).toBe(true);
                    return [4 /*yield*/, storage.save(configPath, registry)];
                case 1:
                    _b.sent();
                    _a = vitest_1.expect;
                    return [4 /*yield*/, storage.exists(configPath)];
                case 2:
                    _a.apply(void 0, [_b.sent()]).toBe(true);
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should add and configure Anthropic provider', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry, storage, provider;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    registry = (0, index_js_1.createRegistry)();
                    storage = (0, index_js_1.createStorage)();
                    registry.updateProvider('anthropic', {
                        apiKey: 'sk-ant-test-key',
                        defaultModel: 'claude-3-opus-20240229',
                    });
                    provider = registry.getProvider('anthropic');
                    (0, vitest_1.expect)(provider === null || provider === void 0 ? void 0 : provider.apiKey).toBe('sk-ant-test-key');
                    (0, vitest_1.expect)(provider === null || provider === void 0 ? void 0 : provider.defaultModel).toBe('claude-3-opus-20240229');
                    return [4 /*yield*/, storage.save(configPath, registry)];
                case 1:
                    _a.sent();
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should add custom provider', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry, customProvider;
        var _a;
        return __generator(this, function (_b) {
            registry = (0, index_js_1.createRegistry)();
            customProvider = {
                id: 'my-custom-provider',
                name: 'My Custom AI',
                type: 'custom',
                enabled: true,
                baseUrl: 'https://api.my-ai.com/v1',
                apiKey: 'my-api-key',
                defaultModel: 'custom-v1',
                models: [{
                        id: 'custom-v1',
                        name: 'Custom V1',
                        maxTokens: 4096,
                        supportsStreaming: true,
                        supportsVision: false,
                        supportsTools: true,
                    }],
                timeout: 30000,
                retries: 3,
            };
            registry.addProvider(customProvider);
            (0, vitest_1.expect)(registry.hasProvider('my-custom-provider')).toBe(true);
            (0, vitest_1.expect)((_a = registry.getProvider('my-custom-provider')) === null || _a === void 0 ? void 0 : _a.baseUrl).toBe('https://api.my-ai.com/v1');
            return [2 /*return*/];
        });
    }); });
    (0, vitest_1.it)('should switch between providers', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry;
        var _a, _b;
        return __generator(this, function (_c) {
            registry = (0, index_js_1.createRegistry)();
            registry.setDefaultProvider('openai');
            (0, vitest_1.expect)((_a = registry.getDefaultProvider()) === null || _a === void 0 ? void 0 : _a.id).toBe('openai');
            registry.setDefaultProvider('anthropic');
            (0, vitest_1.expect)((_b = registry.getDefaultProvider()) === null || _b === void 0 ? void 0 : _b.id).toBe('anthropic');
            return [2 /*return*/];
        });
    }); });
    (0, vitest_1.it)('should enable and disable providers', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry, enabledProviders;
        var _a, _b;
        return __generator(this, function (_c) {
            registry = (0, index_js_1.createRegistry)();
            registry.disableProvider('openai');
            (0, vitest_1.expect)((_a = registry.getProvider('openai')) === null || _a === void 0 ? void 0 : _a.enabled).toBe(false);
            enabledProviders = registry.getEnabledProviders();
            (0, vitest_1.expect)(enabledProviders.some(function (p) { return p.id === 'openai'; })).toBe(false);
            registry.enableProvider('openai');
            (0, vitest_1.expect)((_b = registry.getProvider('openai')) === null || _b === void 0 ? void 0 : _b.enabled).toBe(true);
            return [2 /*return*/];
        });
    }); });
    (0, vitest_1.it)('should persist configuration across restarts', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry1, storage, registry2;
        var _a, _b;
        return __generator(this, function (_c) {
            switch (_c.label) {
                case 0:
                    registry1 = (0, index_js_1.createRegistry)();
                    storage = (0, index_js_1.createStorage)();
                    registry1.updateProvider('openai', { apiKey: 'persistent-key' });
                    registry1.setDefaultProvider('anthropic');
                    registry1.disableProvider('google');
                    return [4 /*yield*/, storage.save(configPath, registry1)];
                case 1:
                    _c.sent();
                    return [4 /*yield*/, storage.load(configPath)];
                case 2:
                    registry2 = _c.sent();
                    (0, vitest_1.expect)((_a = registry2.getProvider('openai')) === null || _a === void 0 ? void 0 : _a.apiKey).toBe('persistent-key');
                    (0, vitest_1.expect)(registry2.getDefaultProviderId()).toBe('anthropic');
                    (0, vitest_1.expect)((_b = registry2.getProvider('google')) === null || _b === void 0 ? void 0 : _b.enabled).toBe(false);
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should validate provider configuration', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry, result;
        return __generator(this, function (_a) {
            registry = (0, index_js_1.createRegistry)();
            result = registry.validateProvider('openai');
            (0, vitest_1.expect)(result.valid).toBe(true);
            (0, vitest_1.expect)(result.errors).toHaveLength(0);
            return [2 /*return*/];
        });
    }); });
    (0, vitest_1.it)('should test connection to provider', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry, result;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    registry = (0, index_js_1.createRegistry)();
                    return [4 /*yield*/, registry.testConnection('openai')];
                case 1:
                    result = _a.sent();
                    (0, vitest_1.expect)(result.success).toBe(false);
                    (0, vitest_1.expect)(result.error).toContain('API key');
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should handle provider removal', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry;
        return __generator(this, function (_a) {
            registry = (0, index_js_1.createRegistry)();
            registry.setDefaultProvider('openai');
            registry.removeProvider('openai');
            (0, vitest_1.expect)(registry.hasProvider('openai')).toBe(false);
            (0, vitest_1.expect)(registry.getDefaultProviderId()).toBeNull();
            return [2 /*return*/];
        });
    }); });
    (0, vitest_1.it)('should handle multiple provider operations', function () { return __awaiter(void 0, void 0, void 0, function () {
        var registry, allProviders, enabledProviders;
        var _a;
        return __generator(this, function (_b) {
            registry = (0, index_js_1.createRegistry)();
            registry.updateProvider('openai', { apiKey: 'openai-key' });
            registry.updateProvider('anthropic', { apiKey: 'anthropic-key' });
            registry.setDefaultProvider('anthropic');
            registry.disableProvider('google');
            allProviders = registry.getAllProviders();
            enabledProviders = registry.getEnabledProviders();
            (0, vitest_1.expect)(allProviders.length).toBeGreaterThanOrEqual(5);
            (0, vitest_1.expect)(enabledProviders.some(function (p) { return p.id === 'google'; })).toBe(false);
            (0, vitest_1.expect)((_a = registry.getDefaultProvider()) === null || _a === void 0 ? void 0 : _a.id).toBe('anthropic');
            return [2 /*return*/];
        });
    }); });
});
