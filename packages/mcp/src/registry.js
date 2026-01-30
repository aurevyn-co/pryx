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
exports.MCPRegistry = void 0;
exports.createRegistry = createRegistry;
var types_js_1 = require("./types.js");
var validation_js_1 = require("./validation.js");
var MCPRegistry = /** @class */ (function () {
    function MCPRegistry() {
        this._servers = new Map();
        this._version = types_js_1.CURRENT_VERSION;
    }
    MCPRegistry.prototype.addServer = function (config) {
        if (this._servers.has(config.id)) {
            throw new types_js_1.MCPServerAlreadyExistsError(config.id);
        }
        (0, validation_js_1.assertValidMCPServerConfig)(config);
        this._servers.set(config.id, __assign({}, config));
    };
    MCPRegistry.prototype.updateServer = function (id, updates) {
        var existing = this._servers.get(id);
        if (!existing) {
            throw new types_js_1.MCPServerNotFoundError(id);
        }
        var updated = __assign(__assign({}, existing), updates);
        (0, validation_js_1.assertValidMCPServerConfig)(updated);
        this._servers.set(id, updated);
        return updated;
    };
    MCPRegistry.prototype.removeServer = function (id) {
        if (!this._servers.has(id)) {
            throw new types_js_1.MCPServerNotFoundError(id);
        }
        this._servers.delete(id);
    };
    MCPRegistry.prototype.getServer = function (id) {
        return this._servers.get(id);
    };
    MCPRegistry.prototype.getAllServers = function () {
        return Array.from(this._servers.values());
    };
    MCPRegistry.prototype.getEnabledServers = function () {
        return this.getAllServers().filter(function (s) { return s.enabled; });
    };
    MCPRegistry.prototype.getServersByType = function (type) {
        return this.getAllServers().filter(function (s) { return s.transport.type === type; });
    };
    MCPRegistry.prototype.hasServer = function (id) {
        return this._servers.has(id);
    };
    MCPRegistry.prototype.enableServer = function (id) {
        this.updateServer(id, { enabled: true });
    };
    MCPRegistry.prototype.disableServer = function (id) {
        this.updateServer(id, { enabled: false });
    };
    MCPRegistry.prototype.enableAll = function () {
        for (var _i = 0, _a = this._servers.values(); _i < _a.length; _i++) {
            var server = _a[_i];
            server.enabled = true;
        }
    };
    MCPRegistry.prototype.disableAll = function () {
        for (var _i = 0, _a = this._servers.values(); _i < _a.length; _i++) {
            var server = _a[_i];
            server.enabled = false;
        }
    };
    MCPRegistry.prototype.enableType = function (type) {
        for (var _i = 0, _a = this._servers.values(); _i < _a.length; _i++) {
            var server = _a[_i];
            if (server.transport.type === type) {
                server.enabled = true;
            }
        }
    };
    MCPRegistry.prototype.disableType = function (type) {
        for (var _i = 0, _a = this._servers.values(); _i < _a.length; _i++) {
            var server = _a[_i];
            if (server.transport.type === type) {
                server.enabled = false;
            }
        }
    };
    MCPRegistry.prototype.updateServerStatus = function (id, status) {
        var server = this._servers.get(id);
        if (!server) {
            throw new types_js_1.MCPServerNotFoundError(id);
        }
        server.status = __assign(__assign({}, server.status), status);
    };
    MCPRegistry.prototype.getFallbackServers = function (id) {
        var _this = this;
        var server = this._servers.get(id);
        if (!server) {
            return [];
        }
        return server.settings.fallbackServers
            .map(function (fallbackId) { return _this._servers.get(fallbackId); })
            .filter(function (s) { return s !== undefined; });
    };
    MCPRegistry.prototype.validateServer = function (id) {
        var server = this._servers.get(id);
        if (!server) {
            return { valid: false, errors: ["Server not found: ".concat(id)] };
        }
        return (0, validation_js_1.validateMCPServerConfig)(server);
    };
    MCPRegistry.prototype.testConnection = function (id) {
        return __awaiter(this, void 0, void 0, function () {
            var server, start, latency;
            return __generator(this, function (_a) {
                server = this._servers.get(id);
                if (!server) {
                    return [2 /*return*/, {
                            success: false,
                            error: "Server not found: ".concat(id),
                        }];
                }
                start = performance.now();
                try {
                    latency = performance.now() - start;
                    return [2 /*return*/, {
                            success: true,
                            latency: latency,
                            capabilities: server.capabilities,
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
    MCPRegistry.prototype.getReconnectDelay = function (id) {
        var server = this._servers.get(id);
        if (!server || !server.status) {
            return 0;
        }
        return (0, validation_js_1.calculateBackoff)(server.status.reconnectAttempts);
    };
    MCPRegistry.prototype.toJSON = function () {
        return {
            version: this._version,
            servers: this.getAllServers(),
        };
    };
    MCPRegistry.prototype.fromJSON = function (data) {
        if (data.version !== types_js_1.CURRENT_VERSION) {
            throw new types_js_1.MCPValidationError(["Unsupported version: ".concat(data.version)]);
        }
        this._servers.clear();
        for (var _i = 0, _a = data.servers; _i < _a.length; _i++) {
            var server = _a[_i];
            (0, validation_js_1.assertValidMCPServerConfig)(server);
            this._servers.set(server.id, server);
        }
    };
    MCPRegistry.prototype.clear = function () {
        this._servers.clear();
    };
    Object.defineProperty(MCPRegistry.prototype, "size", {
        get: function () {
            return this._servers.size;
        },
        enumerable: false,
        configurable: true
    });
    return MCPRegistry;
}());
exports.MCPRegistry = MCPRegistry;
function createRegistry() {
    return new MCPRegistry();
}
