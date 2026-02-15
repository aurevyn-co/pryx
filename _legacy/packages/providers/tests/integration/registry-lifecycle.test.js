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
var registry_js_1 = require("../../src/registry.js");
var storage_js_1 = require("../../src/storage.js");
(0, vitest_1.describe)('Provider Registry Lifecycle', function () {
    var registry;
    var storage;
    var tempDir;
    var configPath;
    (0, vitest_1.beforeEach)(function () { return __awaiter(void 0, void 0, void 0, function () {
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    registry = (0, registry_js_1.createRegistry)();
                    storage = (0, storage_js_1.createStorage)();
                    return [4 /*yield*/, (0, promises_1.mkdtemp)((0, path_1.join)((0, os_1.tmpdir)(), 'providers-lifecycle-test-'))];
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
    (0, vitest_1.it)('should complete full lifecycle: init → add → save → load → update → remove', function () { return __awaiter(void 0, void 0, void 0, function () {
        var customProvider, _a, loadedRegistry;
        var _b, _c;
        return __generator(this, function (_d) {
            switch (_d.label) {
                case 0:
                    customProvider = {
                        id: 'custom-ai',
                        name: 'Custom AI Service',
                        type: 'custom',
                        enabled: true,
                        baseUrl: 'https://api.custom-ai.com/v1',
                        apiKey: 'sk-custom-key',
                        defaultModel: 'custom-model',
                        models: [{
                                id: 'custom-model',
                                name: 'Custom Model',
                                maxTokens: 8192,
                                supportsStreaming: true,
                                supportsVision: false,
                                supportsTools: true,
                            }],
                        timeout: 30000,
                        retries: 3,
                    };
                    registry.addProvider(customProvider);
                    (0, vitest_1.expect)(registry.hasProvider('custom-ai')).toBe(true);
                    return [4 /*yield*/, storage.save(configPath, registry)];
                case 1:
                    _d.sent();
                    _a = vitest_1.expect;
                    return [4 /*yield*/, storage.exists(configPath)];
                case 2:
                    _a.apply(void 0, [_d.sent()]).toBe(true);
                    return [4 /*yield*/, storage.load(configPath)];
                case 3:
                    loadedRegistry = _d.sent();
                    (0, vitest_1.expect)(loadedRegistry.hasProvider('custom-ai')).toBe(true);
                    (0, vitest_1.expect)((_b = loadedRegistry.getProvider('custom-ai')) === null || _b === void 0 ? void 0 : _b.apiKey).toBe('sk-custom-key');
                    loadedRegistry.updateProvider('custom-ai', { apiKey: 'sk-updated-key' });
                    (0, vitest_1.expect)((_c = loadedRegistry.getProvider('custom-ai')) === null || _c === void 0 ? void 0 : _c.apiKey).toBe('sk-updated-key');
                    loadedRegistry.removeProvider('custom-ai');
                    (0, vitest_1.expect)(loadedRegistry.hasProvider('custom-ai')).toBe(false);
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should handle multiple providers', function () { return __awaiter(void 0, void 0, void 0, function () {
        var providers, _i, providers_1, provider, loaded;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    providers = [
                        {
                            id: 'provider-1',
                            name: 'Provider 1',
                            type: 'openai',
                            enabled: true,
                            apiKey: 'key-1',
                            models: [{ id: 'model-1', name: 'Model 1', maxTokens: 1000, supportsStreaming: true, supportsVision: false, supportsTools: false }],
                            timeout: 30000,
                            retries: 3,
                        },
                        {
                            id: 'provider-2',
                            name: 'Provider 2',
                            type: 'anthropic',
                            enabled: true,
                            apiKey: 'key-2',
                            models: [{ id: 'model-2', name: 'Model 2', maxTokens: 2000, supportsStreaming: true, supportsVision: false, supportsTools: false }],
                            timeout: 30000,
                            retries: 3,
                        },
                    ];
                    for (_i = 0, providers_1 = providers; _i < providers_1.length; _i++) {
                        provider = providers_1[_i];
                        registry.addProvider(provider);
                    }
                    (0, vitest_1.expect)(registry.size).toBeGreaterThanOrEqual(7);
                    return [4 /*yield*/, storage.save(configPath, registry)];
                case 1:
                    _a.sent();
                    return [4 /*yield*/, storage.load(configPath)];
                case 2:
                    loaded = _a.sent();
                    (0, vitest_1.expect)(loaded.hasProvider('provider-1')).toBe(true);
                    (0, vitest_1.expect)(loaded.hasProvider('provider-2')).toBe(true);
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should persist default provider', function () { return __awaiter(void 0, void 0, void 0, function () {
        var loaded;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    registry.setDefaultProvider('anthropic');
                    (0, vitest_1.expect)(registry.getDefaultProviderId()).toBe('anthropic');
                    return [4 /*yield*/, storage.save(configPath, registry)];
                case 1:
                    _a.sent();
                    return [4 /*yield*/, storage.load(configPath)];
                case 2:
                    loaded = _a.sent();
                    (0, vitest_1.expect)(loaded.getDefaultProviderId()).toBe('anthropic');
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should handle enable/disable providers', function () { return __awaiter(void 0, void 0, void 0, function () {
        var loaded;
        var _a, _b, _c;
        return __generator(this, function (_d) {
            switch (_d.label) {
                case 0:
                    registry.disableProvider('openai');
                    (0, vitest_1.expect)((_a = registry.getProvider('openai')) === null || _a === void 0 ? void 0 : _a.enabled).toBe(false);
                    return [4 /*yield*/, storage.save(configPath, registry)];
                case 1:
                    _d.sent();
                    return [4 /*yield*/, storage.load(configPath)];
                case 2:
                    loaded = _d.sent();
                    (0, vitest_1.expect)((_b = loaded.getProvider('openai')) === null || _b === void 0 ? void 0 : _b.enabled).toBe(false);
                    loaded.enableProvider('openai');
                    (0, vitest_1.expect)((_c = loaded.getProvider('openai')) === null || _c === void 0 ? void 0 : _c.enabled).toBe(true);
                    return [2 /*return*/];
            }
        });
    }); });
    (0, vitest_1.it)('should handle provider updates', function () { return __awaiter(void 0, void 0, void 0, function () {
        var provider;
        return __generator(this, function (_a) {
            registry.updateProvider('openai', {
                apiKey: 'sk-test',
                defaultModel: 'gpt-4o',
            });
            provider = registry.getProvider('openai');
            (0, vitest_1.expect)(provider === null || provider === void 0 ? void 0 : provider.apiKey).toBe('sk-test');
            (0, vitest_1.expect)(provider === null || provider === void 0 ? void 0 : provider.defaultModel).toBe('gpt-4o');
            return [2 /*return*/];
        });
    }); });
    (0, vitest_1.it)('should handle registry clear', function () { return __awaiter(void 0, void 0, void 0, function () {
        return __generator(this, function (_a) {
            registry.clear();
            (0, vitest_1.expect)(registry.size).toBe(0);
            (0, vitest_1.expect)(registry.getDefaultProviderId()).toBeNull();
            return [2 /*return*/];
        });
    }); });
});
