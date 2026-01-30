"use strict";
// Comprehensive keyboard shortcuts for Pryx TUI
// ANSI escape sequences reference: https://en.wikipedia.org/wiki/ANSI_escape_code
Object.defineProperty(exports, "__esModule", { value: true });
exports.KEYBINDINGS = exports.KEYS = void 0;
exports.parseKeySequence = parseKeySequence;
exports.isPrintable = isPrintable;
// Key sequence constants
exports.KEYS = {
    // Navigation
    ARROW_UP: "\u001b[A",
    ARROW_DOWN: "\u001b[B",
    ARROW_RIGHT: "\u001b[C",
    ARROW_LEFT: "\u001b[D",
    ARROW_UP_ALT: "\u001bOA",
    ARROW_DOWN_ALT: "\u001bOB",
    ARROW_RIGHT_ALT: "\u001bOC",
    ARROW_LEFT_ALT: "\u001bOD",
    HOME: "\u001b[H",
    END: "\u001b[F",
    HOME_ALT: "\u001b[1~",
    END_ALT: "\u001b[4~",
    PAGE_UP: "\u001b[5~",
    PAGE_DOWN: "\u001b[6~",
    // Control combinations
    CTRL_A: "\u0001",
    CTRL_E: "\u0005",
    CTRL_K: "\u000b",
    CTRL_U: "\u0015",
    CTRL_W: "\u0017",
    CTRL_Y: "\u0019",
    CTRL_L: "\u000c",
    CTRL_C: "\u0003",
    CTRL_D: "\u0004",
    CTRL_N: "\u000e",
    CTRL_P: "\u0010",
    // Special keys
    TAB: "\t",
    RETURN: "\r",
    NEWLINE: "\n",
    BACKSPACE: "\u007f",
    BACKSPACE_ALT: "\b",
    DELETE: "\u001b[3~",
    ESCAPE: "\u001b",
    SPACE: " ",
};
// All keybindings
exports.KEYBINDINGS = [
    // Application
    { key: "/", description: "Open command palette", category: "application" },
    { key: "Esc", description: "Close palette/cancel", category: "application" },
    { key: "Ctrl+C", description: "Copy / Cancel operation", category: "application" },
    { key: "Ctrl+L", description: "Clear screen / Refresh", category: "application" },
    { key: "Tab", description: "Switch to next view", category: "application" },
    { key: "Shift+Tab", description: "Switch to previous view", category: "application" },
    // Navigation in input
    { key: "←", description: "Move cursor left", category: "navigation" },
    { key: "→", description: "Move cursor right", category: "navigation" },
    { key: "Home", description: "Go to start of line", category: "navigation" },
    { key: "End", description: "Go to end of line", category: "navigation" },
    { key: "Ctrl+A", description: "Go to start of line", category: "navigation" },
    { key: "Ctrl+E", description: "Go to end of line", category: "navigation" },
    { key: "Ctrl+←", description: "Jump word left", category: "navigation" },
    { key: "Ctrl+→", description: "Jump word right", category: "navigation" },
    // Editing
    { key: "Backspace", description: "Delete character before cursor", category: "editing" },
    { key: "Delete", description: "Delete character after cursor", category: "editing" },
    { key: "Ctrl+K", description: "Clear from cursor to end", category: "editing" },
    { key: "Ctrl+U", description: "Clear from cursor to start", category: "editing" },
    { key: "Ctrl+W", description: "Delete word before cursor", category: "editing" },
    { key: "Ctrl+Y", description: "Paste (yank)", category: "editing" },
    { key: "Ctrl+D", description: "Delete character under cursor", category: "editing" },
    // History
    { key: "↑", description: "Previous message in history", category: "history" },
    { key: "↓", description: "Next message in history", category: "history" },
    { key: "Ctrl+P", description: "Previous message", category: "history" },
    { key: "Ctrl+N", description: "Next message", category: "history" },
    // Scroll
    { key: "Page Up", description: "Scroll messages up", category: "scroll" },
    { key: "Page Down", description: "Scroll messages down", category: "scroll" },
    { key: "Ctrl+↑", description: "Scroll up one line", category: "scroll" },
    { key: "Ctrl+↓", description: "Scroll down one line", category: "scroll" },
];
// Parse key sequence from buffer
function parseKeySequence(data) {
    var seq = data.toString();
    // Handle Ctrl+ combinations
    if (seq.length === 1) {
        var code = seq.charCodeAt(0);
        if (code < 32) {
            return "Ctrl+".concat(String.fromCharCode(code + 64));
        }
    }
    // Map escape sequences
    switch (seq) {
        case exports.KEYS.ARROW_UP:
        case exports.KEYS.ARROW_UP_ALT:
            return "ArrowUp";
        case exports.KEYS.ARROW_DOWN:
        case exports.KEYS.ARROW_DOWN_ALT:
            return "ArrowDown";
        case exports.KEYS.ARROW_LEFT:
        case exports.KEYS.ARROW_LEFT_ALT:
            return "ArrowLeft";
        case exports.KEYS.ARROW_RIGHT:
        case exports.KEYS.ARROW_RIGHT_ALT:
            return "ArrowRight";
        case exports.KEYS.HOME:
        case exports.KEYS.HOME_ALT:
            return "Home";
        case exports.KEYS.END:
        case exports.KEYS.END_ALT:
            return "End";
        case exports.KEYS.PAGE_UP:
            return "PageUp";
        case exports.KEYS.PAGE_DOWN:
            return "PageDown";
        case exports.KEYS.DELETE:
            return "Delete";
        case exports.KEYS.BACKSPACE:
        case exports.KEYS.BACKSPACE_ALT:
            return "Backspace";
        case exports.KEYS.ESCAPE:
            return "Escape";
        case exports.KEYS.TAB:
            return "Tab";
        case exports.KEYS.RETURN:
        case exports.KEYS.NEWLINE:
            return "Enter";
        default:
            return seq;
    }
}
// Check if key is printable
function isPrintable(key) {
    return key.length === 1 && key.charCodeAt(0) >= 32 && key.charCodeAt(0) < 127;
}
