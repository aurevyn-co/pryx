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
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = App;
var solid_js_1 = require("solid-js");
var solid_1 = require("@opentui/solid");
var hooks_1 = require("../lib/hooks");
var ws_1 = require("../services/ws");
var config_1 = require("../services/config");
var AppHeader_1 = require("./AppHeader");
var Chat_1 = require("./Chat");
var SessionExplorer_1 = require("./SessionExplorer");
var Settings_1 = require("./Settings");
var Channels_1 = require("./Channels");
var Skills_1 = require("./Skills");
var SearchableCommandPalette_1 = require("./SearchableCommandPalette");
var KeyboardShortcuts_1 = require("./KeyboardShortcuts");
var SetupRequired_1 = require("./SetupRequired");
function App() {
    var _this = this;
    var renderer = (0, solid_1.useRenderer)();
    renderer.disableStdoutInterception();
    var ws = (0, hooks_1.useEffectService)(ws_1.WebSocketService);
    var _a = (0, solid_js_1.createSignal)("chat"), view = _a[0], setView = _a[1];
    var _b = (0, solid_js_1.createSignal)(false), showCommands = _b[0], setShowCommands = _b[1];
    var _c = (0, solid_js_1.createSignal)(false), showHelp = _c[0], setShowHelp = _c[1];
    var _d = (0, solid_js_1.createSignal)("Connecting..."), connectionStatus = _d[0], setConnectionStatus = _d[1];
    var _e = (0, solid_js_1.createSignal)(false), hasProvider = _e[0], setHasProvider = _e[1];
    var _f = (0, solid_js_1.createSignal)(false), setupRequired = _f[0], setSetupRequired = _f[1];
    (0, solid_js_1.onMount)(function () {
        var config = (0, config_1.loadConfig)();
        var hasValidProvider = config.model_provider &&
            (config.openai_key || config.anthropic_key || config.glm_key || config.ollama_endpoint);
        if (!hasValidProvider) {
            setSetupRequired(true);
        }
    });
    var handleSetupComplete = function () {
        setSetupRequired(false);
        setHasProvider(true);
        setConnectionStatus("Ready");
    };
    (0, solid_js_1.createEffect)(function () {
        var service = ws();
        if (!service) {
            setConnectionStatus("Runtime Error");
            return;
        }
        var checkStatus = function () { return __awaiter(_this, void 0, void 0, function () {
            var apiUrl, res, data, _a;
            var _b;
            return __generator(this, function (_c) {
                switch (_c.label) {
                    case 0:
                        _c.trys.push([0, 5, , 6]);
                        apiUrl = process.env.PRYX_API_URL || "http://localhost:3000";
                        return [4 /*yield*/, fetch("".concat(apiUrl, "/health"), { method: "GET" })];
                    case 1:
                        res = _c.sent();
                        if (!res.ok) return [3 /*break*/, 3];
                        return [4 /*yield*/, res.json()];
                    case 2:
                        data = _c.sent();
                        if (((_b = data.providers) === null || _b === void 0 ? void 0 : _b.length) > 0) {
                            setHasProvider(true);
                            setConnectionStatus("Ready");
                        }
                        else {
                            setHasProvider(false);
                            setConnectionStatus("No Provider");
                        }
                        return [3 /*break*/, 4];
                    case 3:
                        setConnectionStatus("Runtime Error");
                        _c.label = 4;
                    case 4: return [3 /*break*/, 6];
                    case 5:
                        _a = _c.sent();
                        setConnectionStatus("Disconnected");
                        return [3 /*break*/, 6];
                    case 6: return [2 /*return*/];
                }
            });
        }); };
        checkStatus();
        var interval = setInterval(checkStatus, 5000);
        return function () { return clearInterval(interval); };
    });
    var allCommands = [
        {
            id: "chat",
            name: "Chat",
            description: "Open chat interface",
            category: "Navigation",
            shortcut: "1",
            keywords: ["chat", "talk", "message", "conversation"],
            action: function () {
                setView("chat");
                setShowCommands(false);
            },
        },
        {
            id: "sessions",
            name: "Sessions",
            description: "Browse and manage sessions",
            category: "Navigation",
            shortcut: "2",
            keywords: ["sessions", "history", "conversations", "browse"],
            action: function () {
                setView("sessions");
                setShowCommands(false);
            },
        },
        {
            id: "channels",
            name: "Channels",
            description: "Manage channel integrations",
            category: "Navigation",
            shortcut: "3",
            keywords: ["channels", "telegram", "discord", "slack", "webhooks", "integrations"],
            action: function () {
                setView("channels");
                setShowCommands(false);
            },
        },
        {
            id: "skills",
            name: "Skills",
            description: "Browse and manage skills",
            category: "Navigation",
            shortcut: "4",
            keywords: ["skills", "abilities", "tools", "capabilities"],
            action: function () {
                setView("skills");
                setShowCommands(false);
            },
        },
        {
            id: "settings",
            name: "Settings",
            description: "Configure Pryx",
            category: "Navigation",
            shortcut: "5",
            keywords: ["settings", "config", "preferences", "options"],
            action: function () {
                setView("settings");
                setShowCommands(false);
            },
        },
        {
            id: "new-chat",
            name: "New Chat",
            description: "Start a new conversation",
            category: "Chat",
            keywords: ["new", "chat", "conversation", "start", "fresh"],
            action: function () {
                setView("chat");
                setShowCommands(false);
            },
        },
        {
            id: "clear-chat",
            name: "Clear Chat",
            description: "Clear current conversation",
            category: "Chat",
            keywords: ["clear", "reset", "clean", "chat"],
            action: function () {
                setShowCommands(false);
            },
        },
        {
            id: "help",
            name: "Keyboard Shortcuts",
            description: "Show all keyboard shortcuts",
            category: "System",
            shortcut: "?",
            keywords: ["help", "shortcuts", "keys", "commands", "?"],
            action: function () {
                setShowHelp(true);
                setShowCommands(false);
            },
        },
        {
            id: "quit",
            name: "Quit",
            description: "Exit Pryx",
            category: "System",
            shortcut: "q",
            keywords: ["quit", "exit", "close", "stop"],
            action: function () { return process.exit(0); },
        },
        {
            id: "reload",
            name: "Reload",
            description: "Refresh connection",
            category: "System",
            keywords: ["reload", "refresh", "reconnect", "restart"],
            action: function () {
                setShowCommands(false);
            },
        },
    ];
    var views = ["chat", "sessions", "channels", "skills", "settings"];
    (0, solid_1.useKeyboard)(function (evt) {
        if (showHelp() || showCommands()) {
            return;
        }
        switch (evt.name) {
            case "/":
                evt.preventDefault();
                setShowCommands(true);
                break;
            case "?":
                evt.preventDefault();
                setShowHelp(true);
                break;
            case "tab":
                evt.preventDefault();
                setView(function (prev) {
                    var idx = views.indexOf(prev);
                    return views[(idx + 1) % views.length];
                });
                break;
            case "1":
            case "2":
            case "3":
            case "4":
            case "5": {
                evt.preventDefault();
                var idx = parseInt(evt.name) - 1;
                if (idx < views.length) {
                    setView(views[idx]);
                }
                break;
            }
            case "c":
                if (evt.ctrl) {
                    evt.preventDefault();
                    process.exit(0);
                }
                break;
        }
    });
    var getStatusColor = function () {
        if (connectionStatus() === "Ready")
            return "green";
        if (connectionStatus() === "Connecting...")
            return "yellow";
        return "red";
    };
    return (<solid_js_1.Show when={!setupRequired()} fallback={<SetupRequired_1.default onSetupComplete={handleSetupComplete}/>}>
      <box flexDirection="column" backgroundColor="#0a0a0a" flexGrow={1}>
        <AppHeader_1.default />

        <box flexDirection="row" padding={1} gap={1}>
          <text fg="gray">/</text>
          <text fg="gray">commands</text>
          <box flexGrow={1}/>
          <solid_js_1.Show when={!hasProvider()}>
            <text fg="yellow">⚠️ No Provider</text>
          </solid_js_1.Show>
          <text fg={getStatusColor()}>{connectionStatus()}</text>
        </box>

        <box flexGrow={1} padding={1}>
          <solid_js_1.Switch>
            <solid_js_1.Match when={view() === "chat"}>
              <Chat_1.default disabled={showCommands() || showHelp()}/>
            </solid_js_1.Match>
            <solid_js_1.Match when={view() === "sessions"}>
              <SessionExplorer_1.default />
            </solid_js_1.Match>
            <solid_js_1.Match when={view() === "channels"}>
              <Channels_1.default />
            </solid_js_1.Match>
            <solid_js_1.Match when={view() === "settings"}>
              <Settings_1.default />
            </solid_js_1.Match>
            <solid_js_1.Match when={view() === "skills"}>
              <Skills_1.default />
            </solid_js_1.Match>
          </solid_js_1.Switch>
        </box>

        <box flexDirection="row" padding={1}>
          <text fg="gray">/: Commands | Tab: Switch | 1-5: Views | ?: Help | Ctrl+C: Quit</text>
          <box flexGrow={1}/>
          <text fg="blue">v0.1.0-alpha</text>
        </box>

        <solid_js_1.Show when={showCommands()}>
          <SearchableCommandPalette_1.default commands={allCommands} onClose={function () { return setShowCommands(false); }} placeholder="Type to search commands..."/>
        </solid_js_1.Show>

        <solid_js_1.Show when={showHelp()}>
          <KeyboardShortcuts_1.default onClose={function () { return setShowHelp(false); }}/>
        </solid_js_1.Show>
      </box>
    </solid_js_1.Show>);
}
