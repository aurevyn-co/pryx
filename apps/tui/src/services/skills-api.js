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
exports.SkillsServiceLive = exports.SkillsService = exports.SkillsFetchError = void 0;
var effect_1 = require("effect");
var node_fs_1 = require("node:fs");
var node_path_1 = require("node:path");
var node_os_1 = require("node:os");
function getApiUrl() {
    if (process.env.PRYX_API_URL)
        return process.env.PRYX_API_URL;
    try {
        var port = (0, node_fs_1.readFileSync)((0, node_path_1.join)((0, node_os_1.homedir)(), ".pryx", "runtime.port"), "utf-8").trim();
        return "http://localhost:".concat(port);
    }
    catch (_a) {
        return "http://localhost:3000";
    }
}
var SkillsFetchError = /** @class */ (function () {
    function SkillsFetchError(message, cause) {
        this.message = message;
        this.cause = cause;
        this._tag = "SkillsFetchError";
    }
    return SkillsFetchError;
}());
exports.SkillsFetchError = SkillsFetchError;
exports.SkillsService = effect_1.Context.GenericTag("@pryx/tui/SkillsService");
var makeSkillsService = effect_1.Effect.gen(function () {
    var fetchSkills, toggleSkill, installSkill, uninstallSkill;
    return __generator(this, function (_a) {
        fetchSkills = effect_1.Effect.gen(function () {
            var result;
            var _this = this;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [5 /*yield**/, __values(effect_1.Effect.tryPromise({
                            try: function () { return __awaiter(_this, void 0, void 0, function () {
                                var res, data;
                                return __generator(this, function (_a) {
                                    switch (_a.label) {
                                        case 0: return [4 /*yield*/, fetch("".concat(getApiUrl(), "/skills"))];
                                        case 1:
                                            res = _a.sent();
                                            if (!res.ok)
                                                throw new Error("HTTP ".concat(res.status));
                                            return [4 /*yield*/, res.json()];
                                        case 2:
                                            data = (_a.sent());
                                            return [2 /*return*/, data.skills || []];
                                    }
                                });
                            }); },
                            catch: function (error) { return new SkillsFetchError("Failed to fetch skills", error); },
                        }))];
                    case 1:
                        result = _a.sent();
                        return [2 /*return*/, result];
                }
            });
        });
        toggleSkill = function (skillId, enabled) {
            return effect_1.Effect.gen(function () {
                var _this = this;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [5 /*yield**/, __values(effect_1.Effect.tryPromise({
                                try: function () { return __awaiter(_this, void 0, void 0, function () {
                                    var endpoint, res;
                                    return __generator(this, function (_a) {
                                        switch (_a.label) {
                                            case 0:
                                                endpoint = enabled ? "/skills/enable" : "/skills/disable";
                                                return [4 /*yield*/, fetch("".concat(getApiUrl()).concat(endpoint), {
                                                        method: "POST",
                                                        headers: { "Content-Type": "application/json" },
                                                        body: JSON.stringify({ id: skillId }),
                                                    })];
                                            case 1:
                                                res = _a.sent();
                                                if (!res.ok)
                                                    throw new Error("HTTP ".concat(res.status));
                                                return [2 /*return*/];
                                        }
                                    });
                                }); },
                                catch: function (error) {
                                    return new SkillsFetchError("Failed to ".concat(enabled ? "enable" : "disable", " skill"), error);
                                },
                            }))];
                        case 1:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            });
        };
        installSkill = function (skillId) {
            return effect_1.Effect.gen(function () {
                var _this = this;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [5 /*yield**/, __values(effect_1.Effect.tryPromise({
                                try: function () { return __awaiter(_this, void 0, void 0, function () {
                                    var res;
                                    return __generator(this, function (_a) {
                                        switch (_a.label) {
                                            case 0: return [4 /*yield*/, fetch("".concat(getApiUrl(), "/skills/install"), {
                                                    method: "POST",
                                                    headers: { "Content-Type": "application/json" },
                                                    body: JSON.stringify({ id: skillId }),
                                                })];
                                            case 1:
                                                res = _a.sent();
                                                if (!res.ok)
                                                    throw new Error("HTTP ".concat(res.status));
                                                return [2 /*return*/];
                                        }
                                    });
                                }); },
                                catch: function (error) { return new SkillsFetchError("Failed to install skill", error); },
                            }))];
                        case 1:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            });
        };
        uninstallSkill = function (skillId) {
            return effect_1.Effect.gen(function () {
                var _this = this;
                return __generator(this, function (_a) {
                    switch (_a.label) {
                        case 0: return [5 /*yield**/, __values(effect_1.Effect.tryPromise({
                                try: function () { return __awaiter(_this, void 0, void 0, function () {
                                    var res;
                                    return __generator(this, function (_a) {
                                        switch (_a.label) {
                                            case 0: return [4 /*yield*/, fetch("".concat(getApiUrl(), "/skills/uninstall"), {
                                                    method: "POST",
                                                    headers: { "Content-Type": "application/json" },
                                                    body: JSON.stringify({ id: skillId }),
                                                })];
                                            case 1:
                                                res = _a.sent();
                                                if (!res.ok)
                                                    throw new Error("HTTP ".concat(res.status));
                                                return [2 /*return*/];
                                        }
                                    });
                                }); },
                                catch: function (error) { return new SkillsFetchError("Failed to uninstall skill", error); },
                            }))];
                        case 1:
                            _a.sent();
                            return [2 /*return*/];
                    }
                });
            });
        };
        return [2 /*return*/, {
                fetchSkills: fetchSkills,
                toggleSkill: toggleSkill,
                installSkill: installSkill,
                uninstallSkill: uninstallSkill,
            }];
    });
});
exports.SkillsServiceLive = effect_1.Layer.effect(exports.SkillsService, makeSkillsService);
