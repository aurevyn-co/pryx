"use strict";
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
exports.WebSocketServiceLive = exports.WebSocketService = exports.ConnectionError = void 0;
var effect_1 = require("effect");
var ws_1 = require("ws");
var node_fs_1 = require("node:fs");
var node_path_1 = require("node:path");
var node_os_1 = require("node:os");
// Define errors
var ConnectionError = /** @class */ (function () {
    function ConnectionError(message, originalError) {
        this.message = message;
        this.originalError = originalError;
        this._tag = "ConnectionError";
    }
    return ConnectionError;
}());
exports.ConnectionError = ConnectionError;
exports.WebSocketService = effect_1.Context.GenericTag("@pryx/tui/WebSocketService");
// Implementation
var make = effect_1.Effect.gen(function (_) {
    var statusHub, messageHub, socketRef, getRuntimeURL, connect, send, disconnect;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0: return [5 /*yield**/, __values(effect_1.PubSub.unbounded())];
            case 1:
                statusHub = _a.sent();
                return [5 /*yield**/, __values(effect_1.PubSub.unbounded())];
            case 2:
                messageHub = _a.sent();
                return [5 /*yield**/, __values(effect_1.Ref.make(null))];
            case 3:
                socketRef = _a.sent();
                // Initial status
                return [5 /*yield**/, __values(effect_1.PubSub.publish(statusHub, { _tag: "Disconnected" }))];
            case 4:
                // Initial status
                _a.sent();
                getRuntimeURL = function () {
                    if (process.env.PRYX_WS_URL)
                        return process.env.PRYX_WS_URL;
                    try {
                        var port = (0, node_fs_1.readFileSync)((0, node_path_1.join)((0, node_os_1.homedir)(), ".pryx", "runtime.port"), "utf-8").trim();
                        return "ws://localhost:".concat(port, "/ws");
                    }
                    catch (_a) {
                        return "ws://localhost:3000/ws";
                    }
                };
                connect = effect_1.Effect.gen(function (_) {
                    var url;
                    return __generator(this, function (_a) {
                        switch (_a.label) {
                            case 0: return [5 /*yield**/, __values(effect_1.PubSub.publish(statusHub, { _tag: "Connecting" }))];
                            case 1:
                                _a.sent();
                                url = getRuntimeURL();
                                return [5 /*yield**/, __values(effect_1.Effect.async(function (resume) {
                                        var ws;
                                        try {
                                            ws = new ws_1.default(url);
                                        }
                                        catch (e) {
                                            var err = new ConnectionError("Failed to create WebSocket", e);
                                            effect_1.Effect.runSync(effect_1.PubSub.publish(statusHub, { _tag: "Error", error: err }));
                                            resume(effect_1.Effect.fail(err));
                                            return;
                                        }
                                        ws.onopen = function () {
                                            effect_1.Effect.runSync(effect_1.Ref.set(socketRef, ws));
                                            effect_1.Effect.runSync(effect_1.PubSub.publish(statusHub, { _tag: "Connected" }));
                                            resume(effect_1.Effect.void);
                                        };
                                        ws.onmessage = function (event) {
                                            try {
                                                var raw = event.data.toString();
                                                var parsed = JSON.parse(raw);
                                                effect_1.Effect.runSync(effect_1.PubSub.publish(messageHub, parsed));
                                            }
                                            catch (e) {
                                                effect_1.Effect.runSync(effect_1.Console.error("Failed to parse message", e));
                                            }
                                        };
                                        ws.onerror = function (err) {
                                            var error = new ConnectionError(err.message, err);
                                            console.error("WebSocket error:", err);
                                            effect_1.Effect.runSync(effect_1.PubSub.publish(statusHub, { _tag: "Error", error: error }));
                                        };
                                        ws.onclose = function () {
                                            effect_1.Effect.runSync(effect_1.Ref.set(socketRef, null));
                                            effect_1.Effect.runSync(effect_1.PubSub.publish(statusHub, { _tag: "Disconnected" }));
                                        };
                                    }))];
                            case 2:
                                _a.sent();
                                return [2 /*return*/];
                        }
                    });
                });
                send = function (msg) {
                    return effect_1.Effect.gen(function (_) {
                        var ws, e_1;
                        return __generator(this, function (_a) {
                            switch (_a.label) {
                                case 0: return [5 /*yield**/, __values(effect_1.Ref.get(socketRef))];
                                case 1:
                                    ws = _a.sent();
                                    if (!(!ws || ws.readyState !== ws_1.default.OPEN)) return [3 /*break*/, 3];
                                    return [5 /*yield**/, __values(effect_1.Effect.fail(new ConnectionError("Not connected")))];
                                case 2: return [2 /*return*/, _a.sent()];
                                case 3:
                                    _a.trys.push([3, 4, , 6]);
                                    ws.send(JSON.stringify(msg));
                                    return [3 /*break*/, 6];
                                case 4:
                                    e_1 = _a.sent();
                                    return [5 /*yield**/, __values(effect_1.Effect.fail(new ConnectionError("Send failed", e_1)))];
                                case 5: return [2 /*return*/, _a.sent()];
                                case 6: return [2 /*return*/];
                            }
                        });
                    });
                };
                disconnect = effect_1.Effect.gen(function (_) {
                    var ws;
                    return __generator(this, function (_a) {
                        switch (_a.label) {
                            case 0: return [5 /*yield**/, __values(effect_1.Ref.get(socketRef))];
                            case 1:
                                ws = _a.sent();
                                if (!ws) return [3 /*break*/, 4];
                                ws.close();
                                return [5 /*yield**/, __values(effect_1.Ref.set(socketRef, null))];
                            case 2:
                                _a.sent();
                                return [5 /*yield**/, __values(effect_1.PubSub.publish(statusHub, { _tag: "Disconnected" }))];
                            case 3:
                                _a.sent();
                                _a.label = 4;
                            case 4: return [2 /*return*/];
                        }
                    });
                });
                return [2 /*return*/, {
                        status: effect_1.Stream.fromPubSub(statusHub),
                        messages: effect_1.Stream.fromPubSub(messageHub),
                        connect: connect,
                        send: send,
                        disconnect: disconnect,
                    }];
        }
    });
});
exports.WebSocketServiceLive = effect_1.Layer.effect(exports.WebSocketService, make);
