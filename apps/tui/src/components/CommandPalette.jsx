"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = CommandPalette;
var solid_js_1 = require("solid-js");
var solid_1 = require("@opentui/solid");
var theme_1 = require("../theme");
function CommandPalette(props) {
    var _a = (0, solid_js_1.createSignal)(0), selectedIndex = _a[0], setSelectedIndex = _a[1];
    var _b = (0, solid_js_1.createSignal)("keyboard"), inputMode = _b[0], setInputMode = _b[1];
    var executeCommand = function (index) {
        var cmd = props.commands[index];
        if (cmd) {
            cmd.action();
        }
    };
    (0, solid_1.useKeyboard)(function (evt) {
        setInputMode("keyboard");
        switch (evt.name) {
            case "up":
                evt.preventDefault();
                setSelectedIndex(function (i) { return Math.max(0, i - 1); });
                break;
            case "down":
                evt.preventDefault();
                setSelectedIndex(function (i) { return Math.min(props.commands.length - 1, i + 1); });
                break;
            case "return":
                evt.preventDefault();
                executeCommand(selectedIndex());
                break;
            case "escape":
                evt.preventDefault();
                props.onClose();
                break;
            default:
                if (evt.name >= "1" && evt.name <= "9") {
                    var idx = parseInt(evt.name) - 1;
                    if (idx < props.commands.length) {
                        setSelectedIndex(idx);
                    }
                }
        }
    });
    var handleMouseOver = function (index) {
        if (inputMode() !== "mouse")
            return;
        setSelectedIndex(index);
    };
    var handleMouseMove = function () {
        setInputMode("mouse");
    };
    var handleMouseUp = function (index) {
        setSelectedIndex(index);
        executeCommand(index);
    };
    (0, solid_js_1.createEffect)(function () {
        var idx = selectedIndex();
        if (idx >= props.commands.length) {
            setSelectedIndex(0);
        }
    });
    return (<box position="absolute" top={3} left="20%" width="60%" borderStyle="single" borderColor={theme_1.palette.border} backgroundColor={theme_1.palette.bgPrimary} flexDirection="column" padding={1}>
      <box marginBottom={1}>
        <text fg={theme_1.palette.accent}>Commands</text>
      </box>

      <box flexDirection="column">
        <solid_js_1.For each={props.commands}>
          {function (cmd, index) {
            var isSelected = index() === selectedIndex();
            return (<box flexDirection="row" padding={1} backgroundColor={isSelected ? theme_1.palette.bgSelected : undefined} onMouseMove={handleMouseMove} onMouseOver={function () { return handleMouseOver(index()); }} onMouseUp={function () { return handleMouseUp(index()); }} onMouseDown={function () { return setSelectedIndex(index()); }}>
                <box width={3}>
                  <text fg={isSelected ? theme_1.palette.accent : theme_1.palette.dim}>{index() + 1}.</text>
                </box>
                <box flexGrow={1}>
                  <text fg={isSelected ? theme_1.palette.accent : theme_1.palette.text}>{cmd.label}</text>
                </box>
                <box width={10} alignItems="flex-end">
                  <text fg={theme_1.palette.dim}>{cmd.shortcut}</text>
                </box>
              </box>);
        }}
        </solid_js_1.For>
      </box>

      <box marginTop={1} flexDirection="column" gap={0}>
        <text fg={theme_1.palette.dim}>↑↓ Navigate | Enter Select | Esc Close | 1-9 Quick Select</text>
      </box>
    </box>);
}
