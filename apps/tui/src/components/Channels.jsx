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
exports.default = Channels;
var solid_js_1 = require("solid-js");
var config_1 = require("../services/config");
function Channels() {
    var _this = this;
    var _a = (0, solid_js_1.createSignal)({}), config = _a[0], setConfig = _a[1];
    var fields = (0, solid_js_1.createSignal)([
        { id: "h1", type: "header", label: "TELEGRAM BOT" },
        { id: "tg_status", type: "toggle", key: "telegram_enabled", label: "Status" },
        {
            id: "tg_token",
            type: "input",
            key: "telegram_token",
            label: "Bot Token",
            placeholder: "123456:ABC-...",
        },
        { id: "sep1", type: "header", label: " " },
        { id: "h2", type: "header", label: "GENERIC WEBHOOK" },
        { id: "wh_status", type: "toggle", key: "webhook_enabled", label: "Status" },
    ])[0];
    var _b = (0, solid_js_1.createSignal)(1), selectedIndex = _b[0], setSelectedIndex = _b[1];
    var _c = (0, solid_js_1.createSignal)(false), isEditing = _c[0], setIsEditing = _c[1];
    var _d = (0, solid_js_1.createSignal)(""), status = _d[0], setStatus = _d[1];
    var moveSelection = function (dir) {
        var next = selectedIndex();
        var len = fields().length;
        for (var i = 0; i < len; i++) {
            next = (next + dir + len) % len;
            if (fields()[next].type !== "header")
                break;
        }
        setSelectedIndex(next);
    };
    var handleInput = function (data) {
        if (isEditing())
            return;
        var key = data.toString();
        if (key === "\u001B\u005B\u0041") {
            moveSelection(-1);
        }
        else if (key === "\u001B\u005B\u0042") {
            moveSelection(1);
        }
        else if (key === "\r" || key === "\n") {
            var field = fields()[selectedIndex()];
            if (field.type === "input") {
                setIsEditing(true);
            }
            else if (field.type === "toggle") {
                toggleValue(field.key);
            }
        }
        else if (key === " ") {
            var field = fields()[selectedIndex()];
            if (field.type === "toggle") {
                toggleValue(field.key);
            }
        }
    };
    var toggleValue = function (key) { return __awaiter(_this, void 0, void 0, function () {
        var val;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    val = !config()[key];
                    return [4 /*yield*/, handleSave(key, val)];
                case 1:
                    _a.sent();
                    return [2 /*return*/];
            }
        });
    }); };
    var handleSave = function (key, value) { return __awaiter(_this, void 0, void 0, function () {
        var newConfig;
        var _a;
        return __generator(this, function (_b) {
            newConfig = __assign(__assign({}, config()), (_a = {}, _a[key] = value, _a));
            setConfig(newConfig);
            setIsEditing(false);
            (0, config_1.saveConfig)(newConfig);
            setStatus("Saved!");
            setTimeout(function () { return setStatus(""); }, 2000);
            return [2 /*return*/];
        });
    }); };
    (0, solid_js_1.onMount)(function () { return __awaiter(_this, void 0, void 0, function () {
        var loaded;
        return __generator(this, function (_a) {
            loaded = (0, config_1.loadConfig)();
            setConfig(loaded);
            if (typeof process !== "undefined" && process.stdin.isTTY) {
                process.stdin.on("data", handleInput);
            }
            return [2 /*return*/];
        });
    }); });
    (0, solid_js_1.onCleanup)(function () {
        if (typeof process !== "undefined" && process.stdin) {
            process.stdin.off("data", handleInput);
        }
    });
    var renderValue = function (field) {
        var _a;
        var val = config()[field.key];
        if (field.type === "toggle") {
            return val ? <text fg="green">ENABLED</text> : <text fg="gray">DISABLED</text>;
        }
        if (!val)
            return <text fg="gray">empty</text>;
        if ((_a = field.key) === null || _a === void 0 ? void 0 : _a.includes("token")) {
            return (<text>
          {val.substring(0, 4)}...{val.substring(val.length - 4)}
        </text>);
        }
        return <text>{val}</text>;
    };
    return (<box flexDirection="column" flexGrow={1}>
      <text fg="magenta">Channel Setup</text>
      <text fg="gray">Config Path: ~/.pryx/config.yaml</text>

      <box marginTop={1} flexDirection="column" borderStyle="rounded" padding={1}>
        <solid_js_1.For each={fields()}>
          {function (field, index) {
            if (field.type === "header") {
                return (<box marginTop={field.label === " " ? 0 : 1} marginBottom={0}>
                  <text fg="cyan">{field.label}</text>
                </box>);
            }
            var isSelected = index() === selectedIndex();
            return (<box flexDirection="row" marginBottom={0}>
                <text fg={isSelected ? "cyan" : "gray"}>{isSelected ? "❯ " : "  "}</text>
                <box width={15}>
                  <text>{field.label}:</text>
                </box>

                <solid_js_1.Show when={isEditing() && isSelected} fallback={renderValue(field)}>
                  <box>
                    <text fg="cyan">▌{config()[field.key] || ""}</text>
                  </box>
                </solid_js_1.Show>
              </box>);
        }}
        </solid_js_1.For>
      </box>

      <box marginTop={1}>
        <text fg="green">{status()}</text>
      </box>
      <text fg="gray">↑↓ Select │ Enter/Space Toggle │ Enter Edit</text>
    </box>);
}
