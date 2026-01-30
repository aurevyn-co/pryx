"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = SearchableCommandPalette;
var solid_js_1 = require("solid-js");
var solid_1 = require("@opentui/solid");
var theme_1 = require("../theme");
function SearchableCommandPalette(props) {
    var _a = (0, solid_js_1.createSignal)(""), searchQuery = _a[0], setSearchQuery = _a[1];
    var _b = (0, solid_js_1.createSignal)(0), selectedIndex = _b[0], setSelectedIndex = _b[1];
    var _c = (0, solid_js_1.createSignal)("keyboard"), inputMode = _c[0], setInputMode = _c[1];
    var filterCommands = function (query) {
        if (!query.trim()) {
            return props.commands;
        }
        var lowerQuery = query.toLowerCase();
        return props.commands.filter(function (cmd) {
            var _a;
            var nameMatch = cmd.name.toLowerCase().includes(lowerQuery);
            var descMatch = cmd.description.toLowerCase().includes(lowerQuery);
            var categoryMatch = cmd.category.toLowerCase().includes(lowerQuery);
            var keywordMatch = (_a = cmd.keywords) === null || _a === void 0 ? void 0 : _a.some(function (k) { return k.toLowerCase().includes(lowerQuery); });
            return nameMatch || descMatch || categoryMatch || keywordMatch;
        });
    };
    var filteredCommands = (0, solid_js_1.createMemo)(function () { return filterCommands(searchQuery()); });
    var groupedByCategory = (0, solid_js_1.createMemo)(function () {
        var groups = {};
        filteredCommands().forEach(function (cmd) {
            if (!groups[cmd.category]) {
                groups[cmd.category] = [];
            }
            groups[cmd.category].push(cmd);
        });
        return groups;
    });
    var commandsWithIndices = (0, solid_js_1.createMemo)(function () {
        var globalIdx = 0;
        var result = [];
        Object.entries(groupedByCategory()).forEach(function (_a) {
            var category = _a[0], commands = _a[1];
            commands.forEach(function (cmd) {
                result.push({ cmd: cmd, globalIdx: globalIdx++, category: category });
            });
        });
        return result;
    });
    var getCategoryColor = function (category) {
        var colors = {
            Navigation: theme_1.palette.accent,
            Chat: theme_1.palette.success,
            Skills: theme_1.palette.accentSoft,
            MCP: theme_1.palette.info,
            Settings: theme_1.palette.dim,
            System: theme_1.palette.error,
            Help: theme_1.palette.dim,
        };
        return colors[category] || theme_1.palette.text;
    };
    var executeCommand = function (cmd) {
        cmd.action();
        props.onClose();
    };
    (0, solid_1.useKeyboard)(function (evt) {
        var _a;
        setInputMode("keyboard");
        var preventDefaultKeys = [
            "up",
            "down",
            "return",
            "enter",
            "escape",
            "tab",
            "backspace",
            "delete",
            "space",
        ];
        if (preventDefaultKeys.includes(evt.name)) {
            (_a = evt.preventDefault) === null || _a === void 0 ? void 0 : _a.call(evt);
        }
        switch (evt.name) {
            case "up":
            case "arrowup":
                setSelectedIndex(function (i) { return Math.max(0, i - 1); });
                return;
            case "down":
            case "arrowdown":
                setSelectedIndex(function (i) { return Math.min(filteredCommands().length - 1, i + 1); });
                return;
            case "return":
            case "enter": {
                var commands = filteredCommands();
                if (commands.length > 0) {
                    var idx = selectedIndex();
                    if (idx >= 0 && idx < commands.length) {
                        executeCommand(commands[idx]);
                    }
                }
                return;
            }
            case "escape":
                props.onClose();
                return;
            case "backspace":
            case "delete":
                setSearchQuery(function (q) { return q.slice(0, -1); });
                return;
            case "space":
                setSearchQuery(function (q) { return q + " "; });
                return;
            case "tab":
                return;
        }
        if (evt.name.length === 1) {
            setSearchQuery(function (q) { return q + evt.name; });
        }
    });
    var handleMouseMove = function () {
        setInputMode("mouse");
    };
    var handleMouseOver = function (globalIdx) {
        if (inputMode() !== "mouse")
            return;
        setSelectedIndex(globalIdx);
    };
    var handleMouseUp = function (cmd) {
        executeCommand(cmd);
    };
    var handleMouseDown = function (globalIdx) {
        setSelectedIndex(globalIdx);
    };
    (0, solid_js_1.createEffect)(function () {
        var idx = selectedIndex();
        var total = filteredCommands().length;
        if (idx >= total && total > 0) {
            setSelectedIndex(0);
        }
    });
    var totalCommands = (0, solid_js_1.createMemo)(function () { return props.commands.length; });
    return (<box position="absolute" top={3} left="10%" width="80%" height="80%" borderStyle="double" borderColor={theme_1.palette.border} backgroundColor={theme_1.palette.bgPrimary} flexDirection="column" padding={1}>
      <box flexDirection="row" marginBottom={1} gap={1}>
        <text fg={theme_1.palette.accent}>/</text>
        <box flexGrow={1} borderStyle="single" borderColor={searchQuery() ? theme_1.palette.accent : theme_1.palette.dim} padding={0.5}>
          {searchQuery() ? (<text fg={theme_1.palette.text}>{searchQuery()}</text>) : (<text fg={theme_1.palette.dim}>{props.placeholder || "Type to search..."}</text>)}
        </box>
        <box flexGrow={1}/>
        <text fg={theme_1.palette.dim}>
          {filteredCommands().length} / {totalCommands()}
        </text>
      </box>

      <box flexDirection="column" flexGrow={1} overflow="scroll">
        {commandsWithIndices().map(function (_a, idx) {
            var _b;
            var cmd = _a.cmd, globalIdx = _a.globalIdx, category = _a.category;
            var isFirstInCategory = idx === 0 || ((_b = commandsWithIndices()[idx - 1]) === null || _b === void 0 ? void 0 : _b.category) !== category;
            var isSelected = globalIdx === selectedIndex();
            return (<box flexDirection="column">
              {isFirstInCategory && (<box marginTop={idx > 0 ? 1 : 0} marginBottom={0.5} paddingLeft={0.5}>
                  <text fg={getCategoryColor(category)}>{category}</text>
                </box>)}

              <box flexDirection="row" padding={0.5} backgroundColor={isSelected ? theme_1.palette.bgSelected : undefined} onMouseMove={handleMouseMove} onMouseOver={function () { return handleMouseOver(globalIdx); }} onMouseUp={function () { return handleMouseUp(cmd); }} onMouseDown={function () { return handleMouseDown(globalIdx); }}>
                <box width={25}>
                  <text fg={isSelected ? theme_1.palette.accent : theme_1.palette.accentSoft}>{cmd.name}</text>
                </box>
                <box flexGrow={1}>
                  <text fg={isSelected ? theme_1.palette.text : theme_1.palette.dim}>{cmd.description}</text>
                </box>
                {cmd.shortcut && (<box width={10}>
                    <text fg={isSelected ? theme_1.palette.accentSoft : theme_1.palette.dim}>{cmd.shortcut}</text>
                  </box>)}
              </box>
            </box>);
        })}

        {filteredCommands().length === 0 && (<box flexDirection="column" alignItems="center" marginTop={2}>
            <text fg={theme_1.palette.dim}>No commands found matching "{searchQuery()}"</text>
          </box>)}
      </box>

      <box flexDirection="row" marginTop={1} gap={2}>
        <text fg={theme_1.palette.dim}>↑↓ Navigate | Enter Select | Esc Close | Type to filter</text>
      </box>
    </box>);
}
