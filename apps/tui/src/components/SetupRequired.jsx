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
exports.default = SetupRequired;
var solid_js_1 = require("solid-js");
var config_1 = require("../services/config");
var API_BASE = "http://localhost:3000";
function SetupRequired(props) {
    var _this = this;
    var _a, _b, _c;
    var _d = (0, solid_js_1.createSignal)(1), step = _d[0], setStep = _d[1];
    var _e = (0, solid_js_1.createSignal)(""), provider = _e[0], setProvider = _e[1];
    var _f = (0, solid_js_1.createSignal)(""), apiKey = _f[0], setApiKey = _f[1];
    var _g = (0, solid_js_1.createSignal)(""), modelName = _g[0], setModelName = _g[1];
    var _h = (0, solid_js_1.createSignal)(""), error = _h[0], setError = _h[1];
    var _j = (0, solid_js_1.createSignal)([]), providers = _j[0], setProviders = _j[1];
    var _k = (0, solid_js_1.createSignal)([]), models = _k[0], setModels = _k[1];
    var _l = (0, solid_js_1.createSignal)(false), loading = _l[0], setLoading = _l[1];
    var _m = (0, solid_js_1.createSignal)(""), fetchError = _m[0], setFetchError = _m[1];
    (0, solid_js_1.onMount)(function () { return __awaiter(_this, void 0, void 0, function () {
        var response, data, e_1;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    setLoading(true);
                    _a.label = 1;
                case 1:
                    _a.trys.push([1, 4, 5, 6]);
                    return [4 /*yield*/, fetch("".concat(API_BASE, "/api/v1/providers"))];
                case 2:
                    response = _a.sent();
                    if (!response.ok) {
                        throw new Error("Failed to fetch providers: ".concat(response.status));
                    }
                    return [4 /*yield*/, response.json()];
                case 3:
                    data = _a.sent();
                    setProviders(data.providers || []);
                    return [3 /*break*/, 6];
                case 4:
                    e_1 = _a.sent();
                    setFetchError(e_1 instanceof Error ? e_1.message : "Failed to connect to runtime");
                    setProviders([
                        { id: "openai", name: "OpenAI", requires_api_key: true },
                        { id: "anthropic", name: "Anthropic", requires_api_key: true },
                        { id: "google", name: "Google AI", requires_api_key: true },
                        { id: "ollama", name: "Ollama (Local)", requires_api_key: false },
                    ]);
                    return [3 /*break*/, 6];
                case 5:
                    setLoading(false);
                    return [7 /*endfinally*/];
                case 6: return [2 /*return*/];
            }
        });
    }); });
    var fetchModels = function (providerId) { return __awaiter(_this, void 0, void 0, function () {
        var response, data, e_2;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    _a.trys.push([0, 3, , 4]);
                    return [4 /*yield*/, fetch("".concat(API_BASE, "/api/v1/providers/").concat(providerId, "/models"))];
                case 1:
                    response = _a.sent();
                    if (!response.ok) {
                        throw new Error("Failed to fetch models: ".concat(response.status));
                    }
                    return [4 /*yield*/, response.json()];
                case 2:
                    data = _a.sent();
                    setModels(data.models || []);
                    return [3 /*break*/, 4];
                case 3:
                    e_2 = _a.sent();
                    setModels([]);
                    return [3 /*break*/, 4];
                case 4: return [2 /*return*/];
            }
        });
    }); };
    var handleProviderSelect = function (providerId) { return __awaiter(_this, void 0, void 0, function () {
        var selectedProvider, availableModels, defaultModel;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    setProvider(providerId);
                    selectedProvider = providers().find(function (p) { return p.id === providerId; });
                    return [4 /*yield*/, fetchModels(providerId)];
                case 1:
                    _a.sent();
                    availableModels = models();
                    defaultModel = availableModels.length > 0 ? availableModels[0].id : "";
                    setModelName(defaultModel);
                    setStep(2);
                    setError("");
                    return [2 /*return*/];
            }
        });
    }); };
    var handleSubmit = function () {
        var selectedProvider = providers().find(function (p) { return p.id === provider(); });
        if ((selectedProvider === null || selectedProvider === void 0 ? void 0 : selectedProvider.requires_api_key) && !apiKey().trim()) {
            setError("API key is required");
            return;
        }
        var config = {
            model_provider: provider(),
            model_name: modelName(),
        };
        var keyMapping = {
            openai: "openai_key",
            anthropic: "anthropic_key",
            google: "google_key",
        };
        var configKey = keyMapping[provider()];
        if (configKey && apiKey().trim()) {
            config[configKey] = apiKey();
        }
        if (provider() === "ollama") {
            config.ollama_endpoint = apiKey().trim() || "http://localhost:11434";
        }
        try {
            (0, config_1.saveConfig)(config);
            setStep(3);
            setTimeout(function () {
                props.onSetupComplete();
            }, 1500);
        }
        catch (e) {
            setError("Failed to save configuration");
        }
    };
    var selectedProvider = function () { return providers().find(function (p) { return p.id === provider(); }); };
    return (<box flexDirection="column" flexGrow={1} padding={2}>
      <box marginBottom={2} flexDirection="column">
        <text fg="cyan">Welcome to Pryx!</text>
        <text fg="white">Setup Required</text>
        <text fg="gray">To start chatting, you need to configure an AI provider.</text>
      </box>

      {fetchError() && (<box marginBottom={1}>
          <text fg="yellow">⚠ {fetchError()}</text>
        </box>)}

      <box flexDirection="column">
        <box flexDirection="row" marginBottom={1}>
          <text fg={step() >= 1 ? "cyan" : "gray"}>Step 1: Choose Provider</text>
          {step() > 1 && <text fg="green"> ✓</text>}
        </box>

        {step() === 1 && (<box flexDirection="column" marginLeft={2}>
            {loading() ? (<text fg="gray">Loading providers...</text>) : (providers().map(function (p) { return (<box borderStyle="single" borderColor={provider() === p.id ? "cyan" : "gray"} padding={1} flexDirection="column">
                  <text fg="white">{p.name}</text>
                  <text fg="gray">
                    {p.requires_api_key ? "Requires API key" : "No API key required"}
                  </text>
                </box>); }))}
          </box>)}

        <box flexDirection="row" marginTop={1} marginBottom={1}>
          <text fg={step() >= 2 ? "cyan" : "gray"}>Step 2: API Configuration</text>
          {step() > 2 && <text fg="green"> ✓</text>}
        </box>

        {step() === 2 && (<box flexDirection="column" marginLeft={2}>
            <box>
              <text fg="gray">Selected: {(_a = selectedProvider()) === null || _a === void 0 ? void 0 : _a.name}</text>
            </box>

            <box marginTop={1}>
              <text fg="gray">Model:</text>
              <box flexDirection="column">
                {models().length > 0 ? (models().map(function (m) { return (<box borderStyle="single" borderColor={modelName() === m.id ? "cyan" : "gray"} padding={1}>
                      <text fg={modelName() === m.id ? "cyan" : "white"}>{m.name}</text>
                    </box>); })) : (<text fg="gray">No models available</text>)}
              </box>
            </box>

            <box marginTop={1}>
              <text fg="gray">
                {((_b = selectedProvider()) === null || _b === void 0 ? void 0 : _b.requires_api_key) ? "API Key:" : "Endpoint (optional):"}
              </text>
              <box borderStyle="single" borderColor={error() ? "red" : "gray"} padding={1} flexDirection="row">
                <text fg="white">
                  {apiKey() ||
                (((_c = selectedProvider()) === null || _c === void 0 ? void 0 : _c.requires_api_key)
                    ? "Enter API key..."
                    : "http://localhost:11434")}
                </text>
                <box flexGrow={1}/>
                <text fg="cyan">▌</text>
              </box>
              {error() && <text fg="red">{error()}</text>}
            </box>

            <box marginTop={1}>
              <box borderStyle="single" borderColor="cyan" padding={1}>
                <text fg="cyan">Save Configuration</text>
              </box>
            </box>
          </box>)}

        {step() === 3 && (<box flexDirection="column" alignItems="center" marginTop={2}>
            <text fg="green">✓ Configuration Saved!</text>
            <text fg="gray">Starting Pryx...</text>
          </box>)}
      </box>

      <box flexGrow={1}/>

      <box flexDirection="row">
        <text fg="gray">Need help? docs.pryx.dev</text>
      </box>
    </box>);
}
