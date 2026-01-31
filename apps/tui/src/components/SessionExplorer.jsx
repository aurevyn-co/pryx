"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = SessionExplorer;
var solid_js_1 = require("solid-js");
var effect_1 = require("effect");
var hooks_1 = require("../lib/hooks");
var ws_1 = require("../services/ws");
function SessionExplorer() {
    var ws = (0, hooks_1.useEffectService)(ws_1.WebSocketService);
    var _a = (0, solid_js_1.createSignal)([]), sessions = _a[0], setSessions = _a[1];
    var _b = (0, solid_js_1.createSignal)(""), searchQuery = _b[0], setSearchQuery = _b[1];
    var _c = (0, solid_js_1.createSignal)(0), selectedIndex = _c[0], setSelectedIndex = _c[1];
    var _d = (0, solid_js_1.createSignal)(true), loading = _d[0], setLoading = _d[1];
    (0, solid_js_1.createEffect)(function () {
        var service = ws();
        if (!service)
            return;
        var fiber = effect_1.Effect.runFork(service.messages.pipe(effect_1.Stream.runForEach(function (evt) {
            return effect_1.Effect.sync(function () {
                var _a, _b;
                if (evt.event === "sessions.list") {
                    setSessions((_b = (_a = evt.payload) === null || _a === void 0 ? void 0 : _a.sessions) !== null && _b !== void 0 ? _b : []);
                    setLoading(false);
                }
            });
        })));
        (0, solid_js_1.onCleanup)(function () {
            effect_1.Effect.runFork(effect_1.Fiber.interrupt(fiber));
        });
        effect_1.Effect.runFork(service.send({ event: "sessions.list", payload: {} }));
    });
    var filteredSessions = function () {
        var query = searchQuery().toLowerCase();
        if (!query)
            return sessions();
        return sessions().filter(function (s) { return s.title.toLowerCase().includes(query) || s.id.toLowerCase().includes(query); });
    };
    var handleSearch = function (value) {
        setSearchQuery(value);
        setSelectedIndex(0);
    };
    var handleSelect = function () {
        var session = filteredSessions()[selectedIndex()];
        var service = ws();
        if (session && service) {
            effect_1.Effect.runFork(service.send({ event: "session.resume", payload: { session_id: session.id } }));
        }
    };
    var formatDate = function (dateStr) {
        var d = new Date(dateStr);
        return d.toLocaleDateString() + " " + d.toLocaleTimeString().slice(0, 5);
    };
    var formatCost = function (cost) {
        if (!cost)
            return "-";
        return "$".concat(cost.toFixed(4));
    };
    var formatTokens = function (tokens) {
        if (!tokens)
            return "-";
        return tokens > 1000 ? "".concat((tokens / 1000).toFixed(1), "k") : tokens.toString();
    };
    return (<box flexDirection="column" flexGrow={1}>
      <box marginBottom={1}>
        <text fg="cyan">Session Explorer</text>
        <text fg="gray"> ({filteredSessions().length} sessions)</text>
      </box>

      <box borderStyle="single" marginBottom={1} padding={1}>
        <text fg="gray">🔍 </text>
        <box flexGrow={1}>
          {searchQuery() ? (<text fg="white">{searchQuery()}</text>) : (<text fg="gray">Search sessions...</text>)}
        </box>
      </box>

      <box flexDirection="column" flexGrow={1} borderStyle="rounded" padding={1}>
        {loading() ? (<text fg="gray">Loading sessions...</text>) : filteredSessions().length === 0 ? (<text fg="gray">No sessions found</text>) : (<solid_js_1.For each={filteredSessions()}>
            {function (session, index) { return (<box flexDirection="row" gap={1}>
                <text fg={index() === selectedIndex() ? "cyan" : "white"}>
                  {index() === selectedIndex() ? "▶" : " "}
                </text>
                <text fg={index() === selectedIndex() ? "cyan" : "white"}>
                  {session.title.slice(0, 40)}
                </text>
                <text fg="gray">|</text>
                <text fg="gray">{formatDate(session.updatedAt)}</text>
                <text fg="gray">|</text>
                <text fg="yellow">{formatCost(session.cost)}</text>
                <text fg="gray">|</text>
                <text fg="green">{formatTokens(session.tokens)} tok</text>
              </box>); }}
          </solid_js_1.For>)}
      </box>

      <box marginTop={1}>
        <text fg="gray">↑↓ Navigate │ Enter Resume</text>
      </box>
    </box>);
}
