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
exports.useMouse = useMouse;
exports.useClipboard = useClipboard;
var solid_js_1 = require("solid-js");
function useMouse() {
    var _a = (0, solid_js_1.createSignal)({
        x: 0,
        y: 0,
        button: 0,
        isDown: false,
    }), mouseState = _a[0], setMouseState = _a[1];
    var _b = (0, solid_js_1.createSignal)(null), hoveredIndex = _b[0], setHoveredIndex = _b[1];
    var handleMouseMove = function (x, y) {
        setMouseState(function (prev) { return (__assign(__assign({}, prev), { x: x, y: y })); });
    };
    var handleMouseDown = function (x, y, button) {
        setMouseState({ x: x, y: y, button: button, isDown: true });
    };
    var handleMouseUp = function (x, y, button) {
        setMouseState(function (prev) { return (__assign(__assign({}, prev), { x: x, y: y, button: button, isDown: false })); });
    };
    var isInside = function (boxX, boxY, boxWidth, boxHeight) {
        var state = mouseState();
        return (state.x >= boxX && state.x < boxX + boxWidth && state.y >= boxY && state.y < boxY + boxHeight);
    };
    var getItemIndex = function (startY, itemHeight, count) {
        var state = mouseState();
        var index = Math.floor((state.y - startY) / itemHeight);
        if (index >= 0 && index < count) {
            return index;
        }
        return null;
    };
    return {
        mouseState: mouseState,
        hoveredIndex: hoveredIndex,
        setHoveredIndex: setHoveredIndex,
        handleMouseMove: handleMouseMove,
        handleMouseDown: handleMouseDown,
        handleMouseUp: handleMouseUp,
        isInside: isInside,
        getItemIndex: getItemIndex,
    };
}
function useClipboard() {
    var _this = this;
    var copy = function (text) { return __awaiter(_this, void 0, void 0, function () {
        var proc, _a;
        return __generator(this, function (_b) {
            switch (_b.label) {
                case 0:
                    _b.trys.push([0, 5, , 6]);
                    if (!(typeof navigator !== "undefined" && navigator.clipboard)) return [3 /*break*/, 2];
                    return [4 /*yield*/, navigator.clipboard.writeText(text)];
                case 1:
                    _b.sent();
                    return [2 /*return*/, true];
                case 2:
                    if (!(typeof process !== "undefined")) return [3 /*break*/, 4];
                    proc = Bun.spawn(["pbcopy"], {
                        stdin: "pipe",
                    });
                    proc.stdin.write(text);
                    proc.stdin.end();
                    return [4 /*yield*/, proc.exited];
                case 3:
                    _b.sent();
                    return [2 /*return*/, true];
                case 4: return [2 /*return*/, false];
                case 5:
                    _a = _b.sent();
                    return [2 /*return*/, false];
                case 6: return [2 /*return*/];
            }
        });
    }); };
    return { copy: copy };
}
