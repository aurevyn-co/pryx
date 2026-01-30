"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.colors = exports.commandPaletteTheme = exports.theme = exports.palette = void 0;
exports.palette = {
    text: "#E8E3D5",
    dim: "#7B7F87",
    accent: "#F6C453",
    accentSoft: "#F2A65A",
    accentBright: "#FFD700",
    bgPrimary: "#1a1a1a",
    bgSecondary: "#252525",
    bgSelected: "#2B2F36",
    bgHover: "#2A2D33",
    border: "#3C414B",
    borderAccent: "#F6C453",
    success: "#7DD3A5",
    error: "#F97066",
    warning: "#F2A65A",
    info: "#8CC8FF",
    userBg: "#2B2F36",
    userText: "#F3EEE0",
    systemText: "#9BA3B2",
};
function hexToRgb(hex) {
    var result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result
        ? [parseInt(result[1], 16), parseInt(result[2], 16), parseInt(result[3], 16)]
        : [255, 255, 255];
}
exports.theme = {
    fg: function (color) { return function (text) { return "\u001B[38;2;".concat(hexToRgb(color).join(";"), "m").concat(text, "\u001B[0m"); }; },
    bg: function (color) { return function (text) { return "\u001B[48;2;".concat(hexToRgb(color).join(";"), "m").concat(text, "\u001B[0m"); }; },
    bold: function (text) { return "\u001B[1m".concat(text, "\u001B[0m"); },
    dim: function (text) { return "\u001B[2m".concat(text, "\u001B[0m"); },
    italic: function (text) { return "\u001B[3m".concat(text, "\u001B[0m"); },
    underline: function (text) { return "\u001B[4m".concat(text, "\u001B[0m"); },
    text: function (text) { return exports.theme.fg(exports.palette.text)(text); },
    textDim: function (text) { return exports.theme.fg(exports.palette.dim)(text); },
    accent: function (text) { return exports.theme.fg(exports.palette.accent)(text); },
    accentBold: function (text) { return exports.theme.bold(exports.theme.fg(exports.palette.accent)(text)); },
    border: function (text) { return exports.theme.fg(exports.palette.border)(text); },
    success: function (text) { return exports.theme.fg(exports.palette.success)(text); },
    error: function (text) { return exports.theme.fg(exports.palette.error)(text); },
};
exports.commandPaletteTheme = {
    borderColor: exports.palette.border,
    backgroundColor: exports.palette.bgPrimary,
    titleColor: exports.palette.accent,
    itemText: function (text, selected) {
        return selected ? exports.theme.accentBold(text) : exports.theme.text(text);
    },
    itemNumber: function (num, selected) {
        return selected ? exports.theme.accent(num) : exports.theme.textDim(num);
    },
    itemShortcut: function (shortcut) { return exports.theme.textDim(shortcut); },
    itemBackground: function (selected) { return (selected ? exports.palette.bgSelected : undefined); },
    footerText: function (text) { return exports.theme.textDim(text); },
};
exports.colors = {
    text: exports.palette.text,
    dim: exports.palette.dim,
    accent: exports.palette.accent,
    accentSoft: exports.palette.accentSoft,
    bgPrimary: exports.palette.bgPrimary,
    bgSelected: exports.palette.bgSelected,
    border: exports.palette.border,
};
