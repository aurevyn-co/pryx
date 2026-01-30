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
exports.ChannelRegistry = void 0;
exports.createRegistry = createRegistry;
var types_js_1 = require("./types.js");
var validation_js_1 = require("./validation.js");
var ChannelRegistry = /** @class */ (function () {
    function ChannelRegistry() {
        this._channels = new Map();
        this._version = types_js_1.CURRENT_VERSION;
    }
    ChannelRegistry.prototype.addChannel = function (config) {
        if (this._channels.has(config.id)) {
            throw new types_js_1.ChannelAlreadyExistsError(config.id);
        }
        (0, validation_js_1.assertValidChannelConfig)(config);
        this._channels.set(config.id, __assign({}, config));
    };
    ChannelRegistry.prototype.updateChannel = function (id, updates) {
        var existing = this._channels.get(id);
        if (!existing) {
            throw new types_js_1.ChannelNotFoundError(id);
        }
        var updated = __assign(__assign({}, existing), updates);
        (0, validation_js_1.assertValidChannelConfig)(updated);
        this._channels.set(id, updated);
        return updated;
    };
    ChannelRegistry.prototype.removeChannel = function (id) {
        if (!this._channels.has(id)) {
            throw new types_js_1.ChannelNotFoundError(id);
        }
        this._channels.delete(id);
    };
    ChannelRegistry.prototype.getChannel = function (id) {
        return this._channels.get(id);
    };
    ChannelRegistry.prototype.getAllChannels = function () {
        return Array.from(this._channels.values());
    };
    ChannelRegistry.prototype.getEnabledChannels = function () {
        return this.getAllChannels().filter(function (c) { return c.enabled; });
    };
    ChannelRegistry.prototype.getChannelsByType = function (type) {
        return this.getAllChannels().filter(function (c) { return c.type === type; });
    };
    ChannelRegistry.prototype.hasChannel = function (id) {
        return this._channels.has(id);
    };
    ChannelRegistry.prototype.enableChannel = function (id) {
        this.updateChannel(id, { enabled: true });
    };
    ChannelRegistry.prototype.disableChannel = function (id) {
        this.updateChannel(id, { enabled: false });
    };
    ChannelRegistry.prototype.enableAll = function () {
        for (var _i = 0, _a = this._channels.values(); _i < _a.length; _i++) {
            var channel = _a[_i];
            channel.enabled = true;
        }
    };
    ChannelRegistry.prototype.disableAll = function () {
        for (var _i = 0, _a = this._channels.values(); _i < _a.length; _i++) {
            var channel = _a[_i];
            channel.enabled = false;
        }
    };
    ChannelRegistry.prototype.enableType = function (type) {
        for (var _i = 0, _a = this._channels.values(); _i < _a.length; _i++) {
            var channel = _a[_i];
            if (channel.type === type) {
                channel.enabled = true;
            }
        }
    };
    ChannelRegistry.prototype.disableType = function (type) {
        for (var _i = 0, _a = this._channels.values(); _i < _a.length; _i++) {
            var channel = _a[_i];
            if (channel.type === type) {
                channel.enabled = false;
            }
        }
    };
    ChannelRegistry.prototype.updateChannelStatus = function (id, status) {
        var channel = this._channels.get(id);
        if (!channel) {
            throw new types_js_1.ChannelNotFoundError(id);
        }
        channel.status = __assign(__assign({}, channel.status), status);
    };
    ChannelRegistry.prototype.validateChannel = function (id) {
        var channel = this._channels.get(id);
        if (!channel) {
            return { valid: false, errors: ["Channel not found: ".concat(id)] };
        }
        return (0, validation_js_1.validateChannelConfig)(channel);
    };
    ChannelRegistry.prototype.testConnection = function (id) {
        return __awaiter(this, void 0, void 0, function () {
            var channel, start, latency;
            return __generator(this, function (_a) {
                channel = this._channels.get(id);
                if (!channel) {
                    return [2 /*return*/, {
                            success: false,
                            error: "Channel not found: ".concat(id),
                        }];
                }
                start = performance.now();
                try {
                    latency = performance.now() - start;
                    return [2 /*return*/, {
                            success: true,
                            latency: latency,
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
    ChannelRegistry.prototype.toJSON = function () {
        return {
            version: this._version,
            channels: this.getAllChannels(),
        };
    };
    ChannelRegistry.prototype.fromJSON = function (data) {
        if (data.version !== types_js_1.CURRENT_VERSION) {
            throw new types_js_1.ChannelValidationError(["Unsupported version: ".concat(data.version)]);
        }
        this._channels.clear();
        for (var _i = 0, _a = data.channels; _i < _a.length; _i++) {
            var channel = _a[_i];
            (0, validation_js_1.assertValidChannelConfig)(channel);
            this._channels.set(channel.id, channel);
        }
    };
    ChannelRegistry.prototype.clear = function () {
        this._channels.clear();
    };
    Object.defineProperty(ChannelRegistry.prototype, "size", {
        get: function () {
            return this._channels.size;
        },
        enumerable: false,
        configurable: true
    });
    return ChannelRegistry;
}());
exports.ChannelRegistry = ChannelRegistry;
function createRegistry() {
    return new ChannelRegistry();
}
