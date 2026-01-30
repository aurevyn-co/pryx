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
var __spreadArray = (this && this.__spreadArray) || function (to, from, pack) {
    if (pack || arguments.length === 2) for (var i = 0, l = from.length, ar; i < l; i++) {
        if (ar || !(i in from)) {
            if (!ar) ar = Array.prototype.slice.call(from, 0, i);
            ar[i] = from[i];
        }
    }
    return to.concat(ar || Array.prototype.slice.call(from));
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = Chat;
var solid_js_1 = require("solid-js");
var effect_1 = require("effect");
var hooks_1 = require("../lib/hooks");
var ws_1 = require("../services/ws");
var Message_1 = require("./Message");
// ANSI escape sequences for special keys
var KEYS = {
    ARROW_UP: "\u001b[A",
    ARROW_DOWN: "\u001b[B",
    ARROW_RIGHT: "\u001b[C",
    ARROW_LEFT: "\u001b[D",
    HOME: "\u001b[H",
    END: "\u001b[F",
    DELETE: "\u001b[3~",
    BACKSPACE: "\u007f",
    RETURN: "\r",
    NEWLINE: "\n",
    TAB: "\t",
    ESCAPE: "\u001b",
    CTRL_A: "\u0001",
    CTRL_E: "\u0005",
    CTRL_K: "\u000b",
    CTRL_U: "\u0015",
    CTRL_W: "\u0017",
    CTRL_C: "\u0003",
};
function Chat(props) {
    var ws = (0, hooks_1.useEffectService)(ws_1.WebSocketService);
    var _a = (0, solid_js_1.createSignal)([]), messages = _a[0], setMessages = _a[1];
    var _b = (0, solid_js_1.createSignal)(""), inputValue = _b[0], setInputValue = _b[1];
    var _c = (0, solid_js_1.createSignal)(0), cursorPosition = _c[0], setCursorPosition = _c[1];
    var sessionId = (0, solid_js_1.createSignal)(crypto.randomUUID())[0];
    var _d = (0, solid_js_1.createSignal)(null), pendingApproval = _d[0], setPendingApproval = _d[1];
    var _e = (0, solid_js_1.createSignal)(false), isStreaming = _e[0], setIsStreaming = _e[1];
    var _f = (0, solid_js_1.createSignal)(""), streamingContent = _f[0], setStreamingContent = _f[1];
    var _g = (0, solid_js_1.createSignal)([]), history = _g[0], setHistory = _g[1];
    var _h = (0, solid_js_1.createSignal)(-1), historyIndex = _h[0], setHistoryIndex = _h[1];
    (0, solid_js_1.createEffect)(function () {
        var service = ws();
        if (!service)
            return;
        var connectFiber = effect_1.Effect.runFork(service.connect);
        var messageFiber = effect_1.Effect.runFork(service.messages.pipe(effect_1.Stream.runForEach(function (evt) { return effect_1.Effect.sync(function () { return handleEvent(evt); }); })));
        (0, solid_js_1.onCleanup)(function () {
            effect_1.Effect.runFork(effect_1.Fiber.interrupt(connectFiber));
            effect_1.Effect.runFork(effect_1.Fiber.interrupt(messageFiber));
            effect_1.Effect.runFork(service.disconnect);
        });
    });
    var handleEvent = function (evt) {
        var _a, _b, _c;
        switch (evt.event) {
            case "message.delta":
                setIsStreaming(true);
                setStreamingContent(function (prev) { var _a, _b; return prev + ((_b = (_a = evt.payload) === null || _a === void 0 ? void 0 : _a.content) !== null && _b !== void 0 ? _b : ""); });
                break;
            case "message.done":
                setIsStreaming(false);
                setMessages(function (prev) { return __spreadArray(__spreadArray([], prev, true), [
                    {
                        type: "assistant",
                        content: streamingContent(),
                        pending: false,
                    },
                ], false); });
                setStreamingContent("");
                break;
            case "tool.start":
                setMessages(function (prev) {
                    var _a;
                    return __spreadArray(__spreadArray([], prev, true), [
                        {
                            type: "tool",
                            content: "Running...",
                            toolName: (_a = evt.payload) === null || _a === void 0 ? void 0 : _a.name,
                            toolStatus: "running",
                        },
                    ], false);
                });
                break;
            case "tool.end":
                setMessages(function (prev) {
                    var _a, _b, _c;
                    var idx = prev.findLastIndex(function (m) { var _a; return m.toolName === ((_a = evt.payload) === null || _a === void 0 ? void 0 : _a.name) && m.toolStatus === "running"; });
                    if (idx >= 0) {
                        var updated = __spreadArray([], prev, true);
                        updated[idx] = __assign(__assign({}, updated[idx]), { content: (_b = (_a = evt.payload) === null || _a === void 0 ? void 0 : _a.result) !== null && _b !== void 0 ? _b : "Done", toolStatus: ((_c = evt.payload) === null || _c === void 0 ? void 0 : _c.error) ? "error" : "done" });
                        return updated;
                    }
                    return prev;
                });
                break;
            case "approval.request":
                setPendingApproval({
                    id: (_a = evt.payload) === null || _a === void 0 ? void 0 : _a.approval_id,
                    description: (_c = (_b = evt.payload) === null || _b === void 0 ? void 0 : _b.description) !== null && _c !== void 0 ? _c : "Action requires approval",
                });
                break;
        }
    };
    var handleSubmit = function () {
        var value = inputValue();
        if (!value.trim())
            return;
        var service = ws();
        if (!service)
            return;
        if (pendingApproval()) {
            var approval = pendingApproval();
            if (value.toLowerCase() === "y" || value.toLowerCase() === "yes") {
                effect_1.Effect.runFork(service.send({
                    type: "approval.response",
                    sessionId: sessionId(),
                    approvalId: approval.id,
                    approved: true,
                }));
                setMessages(function (prev) { return __spreadArray(__spreadArray([], prev, true), [{ type: "system", content: "✅ Approved" }], false); });
                setPendingApproval(null);
                setInputValue("");
                setCursorPosition(0);
                return;
            }
            else if (value.toLowerCase() === "n" || value.toLowerCase() === "no") {
                effect_1.Effect.runFork(service.send({
                    type: "approval.response",
                    sessionId: sessionId(),
                    approvalId: approval.id,
                    approved: false,
                }));
                setMessages(function (prev) { return __spreadArray(__spreadArray([], prev, true), [{ type: "system", content: "❌ Denied" }], false); });
                setPendingApproval(null);
                setInputValue("");
                setCursorPosition(0);
                return;
            }
        }
        // Add to history
        setHistory(function (prev) { return __spreadArray([value], prev, true).slice(0, 100); });
        setHistoryIndex(-1);
        setMessages(function (prev) { return __spreadArray(__spreadArray([], prev, true), [{ type: "user", content: value }], false); });
        effect_1.Effect.runFork(service.send({
            type: "chat.message",
            sessionId: sessionId(),
            content: value,
        }));
        setInputValue("");
        setCursorPosition(0);
        setIsStreaming(true);
    };
    var handleKey = function (data) {
        var _a, _b, _c;
        if (props.disabled)
            return;
        var seq = data.toString();
        var pos = cursorPosition();
        var value = inputValue();
        switch (seq) {
            case KEYS.RETURN:
            case KEYS.NEWLINE:
                handleSubmit();
                break;
            case KEYS.BACKSPACE:
                if (pos > 0) {
                    var newValue = value.slice(0, pos - 1) + value.slice(pos);
                    setInputValue(newValue);
                    setCursorPosition(pos - 1);
                }
                break;
            case KEYS.DELETE:
                if (pos < value.length) {
                    var newValue = value.slice(0, pos) + value.slice(pos + 1);
                    setInputValue(newValue);
                }
                break;
            case KEYS.ARROW_LEFT:
                setCursorPosition(Math.max(0, pos - 1));
                break;
            case KEYS.ARROW_RIGHT:
                setCursorPosition(Math.min(value.length, pos + 1));
                break;
            case KEYS.HOME:
            case KEYS.CTRL_A:
                setCursorPosition(0);
                break;
            case KEYS.END:
            case KEYS.CTRL_E:
                setCursorPosition(value.length);
                break;
            case KEYS.CTRL_K:
                // Clear from cursor to end
                setInputValue(value.slice(0, pos));
                break;
            case KEYS.CTRL_U:
                // Clear from start to cursor
                setInputValue(value.slice(pos));
                setCursorPosition(0);
                break;
            case KEYS.CTRL_W: {
                var beforeCursor = value.slice(0, pos);
                var match = beforeCursor.match(/^(.*\s)?(\S+)$/);
                if (match) {
                    var newValue = (match[1] || "") + value.slice(pos);
                    setInputValue(newValue);
                    setCursorPosition(((_a = match[1]) === null || _a === void 0 ? void 0 : _a.length) || 0);
                }
                break;
            }
            case KEYS.ARROW_UP: {
                var h = history();
                if (h.length > 0) {
                    var newIndex = Math.min(historyIndex() + 1, h.length - 1);
                    setHistoryIndex(newIndex);
                    setInputValue(h[newIndex]);
                    setCursorPosition(((_b = h[newIndex]) === null || _b === void 0 ? void 0 : _b.length) || 0);
                }
                break;
            }
            case KEYS.ARROW_DOWN: {
                var idx = historyIndex();
                if (idx > 0) {
                    var newIndex = idx - 1;
                    setHistoryIndex(newIndex);
                    setInputValue(history()[newIndex]);
                    setCursorPosition(((_c = history()[newIndex]) === null || _c === void 0 ? void 0 : _c.length) || 0);
                }
                else if (idx === 0) {
                    setHistoryIndex(-1);
                    setInputValue("");
                    setCursorPosition(0);
                }
                break;
            }
            case KEYS.ESCAPE:
                // Cancel/Clear
                setInputValue("");
                setCursorPosition(0);
                setHistoryIndex(-1);
                break;
            case KEYS.TAB:
                // Ignore tab in chat input
                break;
            default:
                // Handle printable characters (including multi-byte for copy-paste)
                if (seq.length >= 1 && !seq.startsWith("\u001b")) {
                    // Check if it's a printable character or paste
                    var isPrintable = seq.split("").every(function (c) {
                        var code = c.charCodeAt(0);
                        return code >= 32 && code < 127;
                    });
                    if (isPrintable) {
                        var newValue = value.slice(0, pos) + seq + value.slice(pos);
                        setInputValue(newValue);
                        setCursorPosition(pos + seq.length);
                    }
                }
                break;
        }
    };
    (0, solid_js_1.onMount)(function () {
        if (typeof process !== "undefined" && process.stdin.isTTY) {
            process.stdin.on("data", handleKey);
        }
    });
    (0, solid_js_1.onCleanup)(function () {
        if (typeof process !== "undefined" && process.stdin) {
            process.stdin.off("data", handleKey);
        }
    });
    var displayMessages = function () { return __spreadArray([], messages(), true); };
    // Render input with cursor
    var renderInput = function () {
        var value = inputValue();
        var pos = cursorPosition();
        if (!value) {
            return (<box flexDirection="row">
          <text fg="gray">Type a message... (Enter to send, ↑↓ for history)</text>
          <box flexGrow={1}/>
          <text fg="cyan">▌</text>
        </box>);
        }
        return (<box flexDirection="row" flexWrap="wrap">
        <text fg="white">{value.slice(0, pos)}</text>
        <text fg="cyan" bg="cyan">
          {" "}
        </text>
        <text fg="white">{value.slice(pos)}</text>
      </box>);
    };
    return (<box flexDirection="column" flexGrow={1}>
      <box flexDirection="column" flexGrow={1} borderStyle="single" borderColor="cyan" padding={1} gap={1}>
        <solid_js_1.For each={displayMessages()}>{function (msg) { return <Message_1.default {...msg}/>; }}</solid_js_1.For>

        {isStreaming() && streamingContent() && (<Message_1.default type="assistant" content={streamingContent()} pending={true}/>)}
      </box>

      {pendingApproval() && (<box borderStyle="double" borderColor="yellow" padding={1} marginTop={1} flexDirection="row">
          <text fg="yellow">⚠️ {pendingApproval().description}</text>
          <box flexGrow={1}/>
          <text fg="gray">(y/n)</text>
        </box>)}

      <box borderStyle="single" borderColor={inputValue() ? "cyan" : "gray"} marginTop={1} padding={1} flexDirection="row" gap={1}>
        <text fg="cyan">❯</text>
        <box flexGrow={1}>{renderInput()}</box>
      </box>
    </box>);
}
