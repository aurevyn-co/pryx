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
var effect_1 = require("effect");
var ws_1 = require("./services/ws");
var program = effect_1.Effect.gen(function () {
    var ws, testMsg;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0: return [5 /*yield**/, __values(ws_1.WebSocketService)];
            case 1:
                ws = _a.sent();
                return [5 /*yield**/, __values(effect_1.Console.log("--- Starting WS Test ---"))];
            case 2:
                _a.sent();
                // 1. Connect
                return [5 /*yield**/, __values(effect_1.Console.log("Connecting..."))];
            case 3:
                // 1. Connect
                _a.sent();
                return [5 /*yield**/, __values(effect_1.Effect.fork(ws.connect))];
            case 4:
                _a.sent();
                // 2. Monitor Status
                return [5 /*yield**/, __values(effect_1.Effect.fork(ws.status.pipe(effect_1.Stream.runForEach(function (status) { return effect_1.Console.log("STATUS CHANGE: ".concat(JSON.stringify(status))); }))))];
            case 5:
                // 2. Monitor Status
                _a.sent();
                // 3. Wait for connection (or timeout)
                return [5 /*yield**/, __values(effect_1.Effect.sleep("2 seconds"))];
            case 6:
                // 3. Wait for connection (or timeout)
                _a.sent();
                // 4. Send a test message
                return [5 /*yield**/, __values(effect_1.Console.log("Sending test message..."))];
            case 7:
                // 4. Send a test message
                _a.sent();
                testMsg = { type: "PING", timestamp: Date.now() };
                return [5 /*yield**/, __values(ws.send(testMsg))];
            case 8:
                _a.sent();
                // 5. Wait for messages
                return [5 /*yield**/, __values(effect_1.Effect.fork(ws.messages.pipe(effect_1.Stream.runForEach(function (msg) { return effect_1.Console.log("RECEIVED: ".concat(JSON.stringify(msg))); }))))];
            case 9:
                // 5. Wait for messages
                _a.sent();
                return [5 /*yield**/, __values(effect_1.Effect.sleep("5 seconds"))];
            case 10:
                _a.sent();
                return [5 /*yield**/, __values(effect_1.Console.log("--- Test Complete ---"))];
            case 11:
                _a.sent();
                return [2 /*return*/];
        }
    });
});
// Run with Live dependencies
var runnable = program.pipe(effect_1.Effect.provide(ws_1.WebSocketServiceLive));
effect_1.Effect.runPromise(runnable).catch(console.error);
