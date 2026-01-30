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
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = OnboardingWizard;
var core_1 = require("@opentui/core");
var solid_js_1 = require("solid-js");
var effect_1 = require("effect");
var hooks_1 = require("../lib/hooks");
var ws_1 = require("../services/ws");
function OnboardingWizard(props) {
    var ws = (0, hooks_1.useEffectService)(ws_1.WebSocketService);
    var _a = (0, solid_js_1.createSignal)(1), step = _a[0], setStep = _a[1];
    var _b = (0, solid_js_1.createSignal)({ name: "", path: "" }), workspace = _b[0], setWorkspace = _b[1];
    var _c = (0, solid_js_1.createSignal)({ provider: "", apiKey: "" }), provider = _c[0], setProvider = _c[1];
    var _d = (0, solid_js_1.createSignal)({ botToken: "" }), integration = _d[0], setIntegration = _d[1];
    var _e = (0, solid_js_1.createSignal)(""), input = _e[0], setInput = _e[1];
    var _f = (0, solid_js_1.createSignal)("name"), field = _f[0], setField = _f[1];
    var handleSubmit = function (value) {
        var currentStep = step();
        var currentField = field();
        var service = ws();
        if (currentStep === 1) {
            if (currentField === "name") {
                setWorkspace(function (w) { return (__assign(__assign({}, w), { name: value })); });
                setField("path");
            }
            else {
                setWorkspace(function (w) { return (__assign(__assign({}, w), { path: value })); });
                setStep(2);
                setField("provider");
            }
        }
        else if (currentStep === 2) {
            if (currentField === "provider") {
                setProvider(function (p) { return (__assign(__assign({}, p), { provider: value })); });
                setField("apiKey");
            }
            else {
                setProvider(function (p) { return (__assign(__assign({}, p), { apiKey: value })); });
                setStep(3);
                setField("botToken");
            }
        }
        else if (currentStep === 3) {
            setIntegration({ botToken: value });
            // Save configuration
            if (service) {
                effect_1.Effect.runFork(service.send({
                    event: "config.save",
                    payload: {
                        workspace: workspace(),
                        provider: provider(),
                        integration: { type: "telegram", botToken: value },
                    },
                }));
            }
            setStep("done");
            setTimeout(function () { return props.onComplete(); }, 1500);
        }
        setInput("");
    };
    var getPlaceholder = function () {
        var f = field();
        switch (f) {
            case "name":
                return "Workspace name (e.g., my-project)";
            case "path":
                return "Workspace path (e.g., ~/code/my-project)";
            case "provider":
                return "Model provider (openai, anthropic, google)";
            case "apiKey":
                return "API key (sk-...)";
            case "botToken":
                return "Telegram bot token (from @BotFather)";
        }
    };
    return (<core_1.Box flexDirection="column" flexGrow={1}>
      <core_1.Box marginBottom={1}>
        <core_1.Text bold color="cyan">
          Onboarding Wizard
        </core_1.Text>
        <core_1.Text color="gray"> - Step {step() === "done" ? "✓" : step()} of 3</core_1.Text>
      </core_1.Box>

      <core_1.Box flexDirection="row" marginBottom={1}>
        <core_1.Text color={step() === 1 ? "cyan" : step() === "done" || step() > 1 ? "green" : "gray"}>
          ● Workspace
        </core_1.Text>
        <core_1.Text color="gray"> → </core_1.Text>
        <core_1.Text color={step() === 2 ? "cyan" : step() === "done" || step() > 2 ? "green" : "gray"}>
          ● Provider
        </core_1.Text>
        <core_1.Text color="gray"> → </core_1.Text>
        <core_1.Text color={step() === 3 ? "cyan" : step() === "done" ? "green" : "gray"}>
          ● Integration
        </core_1.Text>
      </core_1.Box>

      <core_1.Box flexDirection="column" flexGrow={1} borderStyle="round" padding={1}>
        <solid_js_1.Switch>
          <solid_js_1.Match when={step() === 1}>
            <core_1.Text bold>Workspace Setup</core_1.Text>
            <core_1.Text color="gray">
              {field() === "name"
            ? "Enter a name for your workspace"
            : "Enter the path to your workspace"}
            </core_1.Text>
          </solid_js_1.Match>
          <solid_js_1.Match when={step() === 2}>
            <core_1.Text bold>AI Provider Setup</core_1.Text>
            <core_1.Text color="gray">
              {field() === "provider" ? "Choose your AI provider" : "Enter your API key"}
            </core_1.Text>
          </solid_js_1.Match>
          <solid_js_1.Match when={step() === 3}>
            <core_1.Text bold>Integration Setup</core_1.Text>
            <core_1.Text color="gray">Enter your Telegram bot token</core_1.Text>
          </solid_js_1.Match>
          <solid_js_1.Match when={step() === "done"}>
            <core_1.Text bold color="green">
              ✓ Setup Complete!
            </core_1.Text>
            <core_1.Text color="gray">Redirecting to main interface...</core_1.Text>
          </solid_js_1.Match>
        </solid_js_1.Switch>

        <solid_js_1.Show when={step() !== "done"}>
          <core_1.Box marginTop={1}>
            <core_1.Input placeholder={getPlaceholder()} value={input()} onChange={setInput} onSubmit={handleSubmit}/>
          </core_1.Box>
        </solid_js_1.Show>
      </core_1.Box>
    </core_1.Box>);
}
