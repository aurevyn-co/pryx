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
exports.default = Settings;
var solid_js_1 = require("solid-js");
var config_1 = require("../services/config");
function Settings() {
    var _this = this;
    var _a = (0, solid_js_1.createSignal)({}), config = _a[0], setConfig = _a[1];
    var fields = (0, solid_js_1.createSignal)([
        { key: "model_provider", label: "Model Provider", placeholder: "ollama, openai, anthropic" },
        { key: "model_name", label: "Model Name", placeholder: "llama3, gpt-4" },
        { key: "openai_key", label: "OpenAI Key", placeholder: "sk-..." },
        { key: "anthropic_key", label: "Anthropic Key", placeholder: "sk-ant-..." },
        { key: "ollama_endpoint", label: "Ollama URL", placeholder: "http://localhost:11434" },
    ])[0];
    var _b = (0, solid_js_1.createSignal)(0), selectedIndex = _b[0], setSelectedIndex = _b[1];
    var _c = (0, solid_js_1.createSignal)(false), isEditing = _c[0], setIsEditing = _c[1];
    var _d = (0, solid_js_1.createSignal)(""), status = _d[0], setStatus = _d[1];
    (0, solid_js_1.onMount)(function () {
        var loaded = (0, config_1.loadConfig)();
        setConfig(loaded);
        if (typeof process !== "undefined" && process.stdin.isTTY) {
            process.stdin.on("data", handleInput);
        }
    });
    (0, solid_js_1.onCleanup)(function () {
        if (typeof process !== "undefined" && process.stdin) {
            process.stdin.off("data", handleInput);
        }
    });
    var handleInput = function (data) {
        if (isEditing())
            return;
        var key = data.toString();
        if (key === "\u001B\u005B\u0041") {
            setSelectedIndex(function (prev) { return (prev - 1 + fields().length) % fields().length; });
        }
        else if (key === "\u001B\u005B\u0042") {
            setSelectedIndex(function (prev) { return (prev + 1) % fields().length; });
        }
        else if (key === "\r" || key === "\n") {
            setIsEditing(true);
        }
    };
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
    return (<box flexDirection="column" flexGrow={1}>
      <text fg="cyan">Configuration</text>
      <text fg="gray">Config Path: ~/.pryx/config.yaml</text>
      <box marginTop={1} flexDirection="column" borderStyle="rounded" padding={1}>
        <solid_js_1.For each={fields()}>
          {function (field, index) { return (<box flexDirection="row" marginBottom={0}>
              <text fg={index() === selectedIndex() ? "cyan" : "gray"}>
                {index() === selectedIndex() ? "❯ " : "  "}
              </text>
              <box width={20}>
                <text>{field.label}:</text>
              </box>

              <solid_js_1.Show when={isEditing() && index() === selectedIndex()} fallback={<box>
                    {config()[field.key] ? (<text fg="white">{config()[field.key]}</text>) : (<text fg="gray">empty</text>)}
                  </box>}>
                <box>
                  <text fg="cyan">▌{config()[field.key] || ""}</text>
                </box>
              </solid_js_1.Show>
            </box>); }}
        </solid_js_1.For>
      </box>
      <box marginTop={1}>
        <text fg="green">{status()}</text>
      </box>
      <text fg="gray">↑↓ Select │ Enter Edit/Save</text>
    </box>);
}
