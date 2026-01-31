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
exports.default = Skills;
// @ts-nocheck
var solid_js_1 = require("solid-js");
function Skills() {
    var _this = this;
    var _a, _b, _c, _d, _e;
    var _f = (0, solid_js_1.createSignal)([]), skills = _f[0], setSkills = _f[1];
    var selectedIndex = (0, solid_js_1.createSignal)(0)[0];
    var _g = (0, solid_js_1.createSignal)(true), loading = _g[0], setLoading = _g[1];
    var _h = (0, solid_js_1.createSignal)(""), error = _h[0], setError = _h[1];
    var detailView = (0, solid_js_1.createSignal)(false)[0];
    var fetchSkills = function () { return __awaiter(_this, void 0, void 0, function () {
        var apiUrl, res, data, _1;
        return __generator(this, function (_a) {
            switch (_a.label) {
                case 0:
                    _a.trys.push([0, 3, 4, 5]);
                    setLoading(true);
                    setError("");
                    apiUrl = process.env.PRYX_API_URL || "http://localhost:3000";
                    return [4 /*yield*/, fetch("".concat(apiUrl, "/skills"))];
                case 1:
                    res = _a.sent();
                    if (!res.ok) {
                        throw new Error("HTTP ".concat(res.status));
                    }
                    return [4 /*yield*/, res.json()];
                case 2:
                    data = _a.sent();
                    setSkills(data.skills || []);
                    return [3 /*break*/, 5];
                case 3:
                    _1 = _a.sent();
                    setError("Failed to load skills");
                    return [3 /*break*/, 5];
                case 4:
                    setLoading(false);
                    return [7 /*endfinally*/];
                case 5: return [2 /*return*/];
            }
        });
    }); };
    (0, solid_js_1.onMount)(function () {
        fetchSkills().catch(function () { });
    });
    var selectedSkill = function () {
        var index = selectedIndex();
        var skillsList = skills();
        if (skillsList.length === 0)
            return null;
        return skillsList[index] || skillsList[0];
    };
    return (<box flexDirection="column" flexGrow={1}>
      <text fg="magenta">Skills Manager</text>
      <text fg="gray">Extend agent capabilities with skills</text>

      <solid_js_1.Show when={loading()}>
        <box marginTop={1}>
          <text fg="yellow">Loading skills...</text>
        </box>
      </solid_js_1.Show>

      <solid_js_1.Show when={error()}>
        <box marginTop={1}>
          <text fg="red">{error()}</text>
        </box>
      </solid_js_1.Show>

      <solid_js_1.Show when={!loading() && !error()}>
        <solid_js_1.Show when={!detailView()} fallback={<box marginTop={1} flexDirection="column" borderStyle="round" padding={1}>
              <text fg="cyan">{(_a = selectedSkill()) === null || _a === void 0 ? void 0 : _a.name}</text>
              <box marginTop={1}>
                <text fg="gray">ID: </text>
                <text>{(_b = selectedSkill()) === null || _b === void 0 ? void 0 : _b.id}</text>
              </box>
              <box marginTop={1}>
                <text fg="gray">Description:</text>
              </box>
              <box marginTop={0}>
                <text>{((_c = selectedSkill()) === null || _c === void 0 ? void 0 : _c.description) || "No description available"}</text>
              </box>
              <box marginTop={1}>
                <text fg="gray">Status: </text>
                <text fg={((_d = selectedSkill()) === null || _d === void 0 ? void 0 : _d.enabled) ? "green" : "gray"}>
                  {((_e = selectedSkill()) === null || _e === void 0 ? void 0 : _e.enabled) ? "ENABLED" : "DISABLED"}
                </text>
              </box>
              <box marginTop={2}>
                <text fg="gray">Press Esc to go back</text>
              </box>
            </box>}>
          <box marginTop={1} flexDirection="column" borderStyle="round" padding={1}>
            <solid_js_1.Show when={skills().length === 0}>
              <text fg="gray">No skills available</text>
            </solid_js_1.Show>
            <solid_js_1.For each={skills()}>
              {function (skill, index) {
            var isSelected = index() === selectedIndex();
            return (<box flexDirection="row" marginBottom={0}>
                    <text fg={isSelected ? "cyan" : "gray"}>{isSelected ? "❯ " : "  "}</text>
                    <box width={30}>
                      <text>{skill.name}</text>
                    </box>
                    <text fg={skill.enabled ? "green" : "gray"}>{skill.enabled ? "✓" : "○"}</text>
                  </box>);
        }}
            </solid_js_1.For>
          </box>
        </solid_js_1.Show>
      </solid_js_1.Show>

      <box marginTop={1}>
        <text fg="gray">
          {detailView() ? "Esc: Back" : "↑↓ Navigate │ Enter: Details │ R: Refresh"}
        </text>
      </box>
    </box>);
}
