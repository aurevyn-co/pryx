"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.useKeybind = useKeybind;
exports.useKeyboardHandler = useKeyboardHandler;
var solid_1 = require("@opentui/solid");
var solid_js_1 = require("solid-js");
function useKeybind() {
    var _a = (0, solid_js_1.createSignal)(false), leader = _a[0], setLeader = _a[1];
    var match = function (keybind, evt) {
        var _a, _b, _c;
        return (((_a = keybind.ctrl) !== null && _a !== void 0 ? _a : false) === evt.ctrl &&
            ((_b = keybind.meta) !== null && _b !== void 0 ? _b : false) === evt.meta &&
            ((_c = keybind.shift) !== null && _c !== void 0 ? _c : false) === evt.shift &&
            keybind.name === evt.name);
    };
    var parse = function (key) {
        var parts = key.toLowerCase().split("+");
        var info = { name: "" };
        for (var _i = 0, parts_1 = parts; _i < parts_1.length; _i++) {
            var part = parts_1[_i];
            switch (part) {
                case "ctrl":
                    info.ctrl = true;
                    break;
                case "alt":
                case "meta":
                    info.meta = true;
                    break;
                case "shift":
                    info.shift = true;
                    break;
                case "esc":
                    info.name = "escape";
                    break;
                case "return":
                case "enter":
                    info.name = "return";
                    break;
                case "up":
                    info.name = "up";
                    break;
                case "down":
                    info.name = "down";
                    break;
                default:
                    info.name = part;
            }
        }
        return info;
    };
    return {
        match: match,
        parse: parse,
        leader: leader,
        setLeader: setLeader,
    };
}
function useKeyboardHandler(handlers) {
    var _a = useKeybind(), match = _a.match, parse = _a.parse;
    (0, solid_1.useKeyboard)(function (evt) {
        for (var _i = 0, handlers_1 = handlers; _i < handlers_1.length; _i++) {
            var _a = handlers_1[_i], key = _a[0], handler = _a[1];
            var keybind = parse(key);
            if (match(keybind, evt)) {
                handler();
                return;
            }
        }
    });
}
