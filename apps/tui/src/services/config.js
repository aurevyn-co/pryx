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
var __values = (this && this.__values) || function(o) {
    var s = typeof Symbol === "function" && Symbol.iterator, m = s && o[s], i = 0;
    if (m) return m.call(o);
    if (o && typeof o.length === "number") return {
        next: function () {
            if (o && i >= o.length) o = void 0;
            return { value: o && o[i++], done: !o };
        }
    };
    throw new TypeError(s ? "Object is not iterable." : "Symbol.iterator is not defined.");
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.getConfigValue = exports.updateConfig = exports.saveConfig = exports.loadConfig = exports.ConfigServiceLive = exports.ConfigService = exports.ConfigSaveError = exports.ConfigLoadError = void 0;
var effect_1 = require("effect");
var node_fs_1 = require("node:fs");
var node_path_1 = require("node:path");
var js_yaml_1 = require("js-yaml");
var node_os_1 = require("node:os");
var CONFIG_PATH = node_path_1.default.join(node_os_1.default.homedir(), ".pryx", "config.yaml");
var ConfigLoadError = /** @class */ (function () {
    function ConfigLoadError(message, cause) {
        this.message = message;
        this.cause = cause;
        this._tag = "ConfigLoadError";
    }
    return ConfigLoadError;
}());
exports.ConfigLoadError = ConfigLoadError;
var ConfigSaveError = /** @class */ (function () {
    function ConfigSaveError(message, cause) {
        this.message = message;
        this.cause = cause;
        this._tag = "ConfigSaveError";
    }
    return ConfigSaveError;
}());
exports.ConfigSaveError = ConfigSaveError;
exports.ConfigService = effect_1.Context.GenericTag("@pryx/tui/ConfigService");
var makeConfigService = effect_1.Effect.gen(function () {
    var load, save, update, getValue;
    return __generator(this, function (_a) {
        load = effect_1.Effect.gen(function () {
            var result;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [5 /*yield**/, __values(effect_1.Effect.try({
                            try: function () {
                                if (!node_fs_1.default.existsSync(CONFIG_PATH))
                                    return {};
                                var content = node_fs_1.default.readFileSync(CONFIG_PATH, "utf-8");
                                return js_yaml_1.default.load(content) || {};
                            },
                            catch: function (error) { return new ConfigLoadError("Failed to load config", error); },
                        }))];
                    case 1:
                        result = _a.sent();
                        return [2 /*return*/, result];
                }
            });
        });
        save = function (cfg) {
            return effect_1.Effect.gen(function () {
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [5 /*yield**/, __values(effect_1.Effect.try({
                                try: function () {
                                    var dir = node_path_1.default.dirname(CONFIG_PATH);
                                    if (!node_fs_1.default.existsSync(dir))
                                        node_fs_1.default.mkdirSync(dir, { recursive: true });
                                    node_fs_1.default.writeFileSync(CONFIG_PATH, js_yaml_1.default.dump(cfg), "utf-8");
                                },
                                catch: function (error) { return new ConfigSaveError("Failed to save config", error); },
                            }))];
                        case 1:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            });
        };
        update = function (updates) {
            return effect_1.Effect.gen(function () {
                var current, updated;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [5 /*yield**/, __values(load)];
                        case 1:
                            current = _a.sent();
                            updated = __assign(__assign({}, current), updates);
                            return [5 /*yield**/, __values(save(updated))];
                        case 2:
                            _a.sent();
                            return [2 /*return*/, updated];
                    }
                });
            });
        };
        getValue = function (key, defaultValue) {
            return effect_1.Effect.gen(function () {
                var config;
                var _a;
                return __generator(this, function (_b) {
                    switch (_b.label) {
                        case 0: return [5 /*yield**/, __values(load)];
                        case 1:
                            config = _b.sent();
                            return [2 /*return*/, (_a = config[key]) !== null && _a !== void 0 ? _a : defaultValue];
                    }
                });
            });
        };
        return [2 /*return*/, {
                load: load,
                save: save,
                update: update,
                getValue: getValue,
            }];
    });
});
exports.ConfigServiceLive = effect_1.Layer.effect(exports.ConfigService, makeConfigService);
var loadConfig = function () {
    try {
        if (!node_fs_1.default.existsSync(CONFIG_PATH))
            return {};
        var content = node_fs_1.default.readFileSync(CONFIG_PATH, "utf-8");
        return js_yaml_1.default.load(content) || {};
    }
    catch (e) {
        console.error("[Config] Failed to load config:", e);
        return {};
    }
};
exports.loadConfig = loadConfig;
var saveConfig = function (cfg) {
    try {
        var dir = node_path_1.default.dirname(CONFIG_PATH);
        if (!node_fs_1.default.existsSync(dir))
            node_fs_1.default.mkdirSync(dir, { recursive: true });
        node_fs_1.default.writeFileSync(CONFIG_PATH, js_yaml_1.default.dump(cfg), "utf-8");
        console.log("[Config] Saved successfully");
    }
    catch (e) {
        console.error("[Config] Failed to save config:", e);
        throw e;
    }
};
exports.saveConfig = saveConfig;
var updateConfig = function (updates) {
    var current = (0, exports.loadConfig)();
    var updated = __assign(__assign({}, current), updates);
    (0, exports.saveConfig)(updated);
    return updated;
};
exports.updateConfig = updateConfig;
var getConfigValue = function (key, defaultValue) {
    var _a;
    var config = (0, exports.loadConfig)();
    return (_a = config[key]) !== null && _a !== void 0 ? _a : defaultValue;
};
exports.getConfigValue = getConfigValue;
