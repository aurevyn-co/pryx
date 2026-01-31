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
exports.ProviderRegistry = void 0;
exports.createRegistry = createRegistry;
var types_js_1 = require("./types.js");
var validation_js_1 = require("./validation.js");
var presets_js_1 = require("./presets.js");
var ProviderRegistry = /** @class */ (function () {
    function ProviderRegistry() {
        this._providers = new Map();
        this._defaultProvider = null;
        this._version = types_js_1.CURRENT_VERSION;
        for (var _i = 0, BUILTIN_PRESETS_1 = presets_js_1.BUILTIN_PRESETS; _i < BUILTIN_PRESETS_1.length; _i++) {
            var preset = BUILTIN_PRESETS_1[_i];
            this._providers.set(preset.id, __assign({}, preset));
        }
    }
    ProviderRegistry.prototype.addProvider = function (config) {
        if (this._providers.has(config.id)) {
            throw new types_js_1.ProviderAlreadyExistsError(config.id);
        }
        (0, validation_js_1.assertValidProviderConfig)(config);
        this._providers.set(config.id, __assign({}, config));
    };
    ProviderRegistry.prototype.updateProvider = function (id, updates) {
        var existing = this._providers.get(id);
        if (!existing) {
            throw new types_js_1.ProviderNotFoundError(id);
        }
        var updated = __assign(__assign({}, existing), updates);
        (0, validation_js_1.assertValidProviderConfig)(updated);
        this._providers.set(id, updated);
        return updated;
    };
    ProviderRegistry.prototype.removeProvider = function (id) {
        if (!this._providers.has(id)) {
            throw new types_js_1.ProviderNotFoundError(id);
        }
        this._providers.delete(id);
        if (this._defaultProvider === id) {
            this._defaultProvider = null;
        }
    };
    ProviderRegistry.prototype.getProvider = function (id) {
        return this._providers.get(id);
    };
    ProviderRegistry.prototype.getAllProviders = function () {
        return Array.from(this._providers.values());
    };
    ProviderRegistry.prototype.getEnabledProviders = function () {
        return this.getAllProviders().filter(function (p) { return p.enabled; });
    };
    ProviderRegistry.prototype.hasProvider = function (id) {
        return this._providers.has(id);
    };
    ProviderRegistry.prototype.setDefaultProvider = function (id) {
        if (!this._providers.has(id)) {
            throw new types_js_1.ProviderNotFoundError(id);
        }
        this._defaultProvider = id;
    };
    ProviderRegistry.prototype.getDefaultProvider = function () {
        if (this._defaultProvider) {
            return this._providers.get(this._defaultProvider);
        }
        var enabled = this.getEnabledProviders();
        return enabled.length > 0 ? enabled[0] : undefined;
    };
    ProviderRegistry.prototype.getDefaultProviderId = function () {
        return this._defaultProvider;
    };
    ProviderRegistry.prototype.enableProvider = function (id) {
        this.updateProvider(id, { enabled: true });
    };
    ProviderRegistry.prototype.disableProvider = function (id) {
        this.updateProvider(id, { enabled: false });
    };
    ProviderRegistry.prototype.validateProvider = function (id) {
        var provider = this._providers.get(id);
        if (!provider) {
            return { valid: false, errors: ["Provider not found: ".concat(id)] };
        }
        return (0, validation_js_1.validateProviderConfig)(provider);
    };
    ProviderRegistry.prototype.testConnection = function (id) {
        return __awaiter(this, void 0, void 0, function () {
            var provider, start, latency;
            return __generator(this, function (_a) {
                provider = this._providers.get(id);
                if (!provider) {
                    return [2 /*return*/, {
                            success: false,
                            error: "Provider not found: ".concat(id),
                        }];
                }
                start = performance.now();
                try {
                    if (!provider.apiKey && provider.type !== 'local') {
                        return [2 /*return*/, {
                                success: false,
                                error: 'API key not configured',
                            }];
                    }
                    latency = performance.now() - start;
                    return [2 /*return*/, {
                            success: true,
                            latency: latency,
                            modelsAvailable: provider.models.map(function (m) { return m.id; }),
                        }];
                }
                catch (error) {
                    return [2 /*return*/, {
                            success: false,
                            error: error instanceof Error ? error.message : 'Unknown error',
                        }];
                }
                return [2 /*return*/];
            });
        });
    };
    ProviderRegistry.prototype.toJSON = function () {
        return {
            version: this._version,
            defaultProvider: this._defaultProvider || undefined,
            providers: this.getAllProviders(),
        };
    };
    ProviderRegistry.prototype.fromJSON = function (data) {
        if (data.version !== types_js_1.CURRENT_VERSION) {
            throw new types_js_1.ProviderValidationError(["Unsupported version: ".concat(data.version)]);
        }
        this._providers.clear();
        for (var _i = 0, _a = data.providers; _i < _a.length; _i++) {
            var provider = _a[_i];
            (0, validation_js_1.assertValidProviderConfig)(provider);
            this._providers.set(provider.id, provider);
        }
        this._defaultProvider = data.defaultProvider || null;
    };
    ProviderRegistry.prototype.clear = function () {
        this._providers.clear();
        this._defaultProvider = null;
    };
    Object.defineProperty(ProviderRegistry.prototype, "size", {
        get: function () {
            return this._providers.size;
        },
        enumerable: false,
        configurable: true
    });
    return ProviderRegistry;
}());
exports.ProviderRegistry = ProviderRegistry;
function createRegistry() {
    return new ProviderRegistry();
}
