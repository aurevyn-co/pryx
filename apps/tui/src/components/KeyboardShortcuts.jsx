"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = KeyboardShortcuts;
var solid_js_1 = require("solid-js");
var solid_1 = require("@opentui/solid");
var keybindings_1 = require("../lib/keybindings");
function KeyboardShortcuts(props) {
    var _a = (0, solid_js_1.createSignal)(null), selectedCategory = _a[0], setSelectedCategory = _a[1];
    (0, solid_1.useKeyboard)(function (evt) {
        var _a;
        switch (evt.name) {
            case "escape":
            case "q":
                (_a = evt.preventDefault) === null || _a === void 0 ? void 0 : _a.call(evt);
                props.onClose();
                return;
        }
    });
    var categories = [
        { id: "application", label: "Application", color: "cyan" },
        { id: "navigation", label: "Navigation", color: "green" },
        { id: "editing", label: "Editing", color: "yellow" },
        { id: "history", label: "History", color: "magenta" },
        { id: "scroll", label: "Scroll", color: "blue" },
    ];
    var filteredBindings = function () {
        if (!selectedCategory())
            return keybindings_1.KEYBINDINGS;
        return keybindings_1.KEYBINDINGS.filter(function (b) { return b.category === selectedCategory(); });
    };
    var groupedByCategory = function () {
        var groups = {};
        filteredBindings().forEach(function (binding) {
            if (!groups[binding.category]) {
                groups[binding.category] = [];
            }
            groups[binding.category].push(binding);
        });
        return groups;
    };
    var getCategoryColor = function (cat) {
        var c = categories.find(function (c) { return c.id === cat; });
        return (c === null || c === void 0 ? void 0 : c.color) || "white";
    };
    var getCategoryLabel = function (cat) {
        var c = categories.find(function (c) { return c.id === cat; });
        return (c === null || c === void 0 ? void 0 : c.label) || cat;
    };
    return (<box position="absolute" top={2} left="10%" width="80%" height="90%" borderStyle="double" borderColor="cyan" backgroundColor="#0a0a0a" flexDirection="column" padding={1}>
      <box flexDirection="row" marginBottom={1}>
        <text fg="cyan">Keyboard Shortcuts</text>
        <box flexGrow={1}/>
        <text fg="gray">Press Esc to close</text>
      </box>

      <box flexDirection="row" gap={1} marginBottom={1}>
        <box padding={1} backgroundColor={!selectedCategory() ? "cyan" : undefined}>
          <text fg={!selectedCategory() ? "black" : "white"}>All</text>
        </box>
        <solid_js_1.For each={categories}>
          {function (cat) { return (<box padding={1} backgroundColor={selectedCategory() === cat.id ? cat.color : undefined}>
              <text fg={selectedCategory() === cat.id ? "black" : cat.color}>{cat.label}</text>
            </box>); }}
        </solid_js_1.For>
      </box>

      <box flexDirection="column" flexGrow={1}>
        <solid_js_1.For each={Object.entries(groupedByCategory())}>
          {function (_a) {
            var category = _a[0], bindings = _a[1];
            return (<box flexDirection="column" marginBottom={1}>
              <box marginBottom={0.5}>
                <text fg={getCategoryColor(category)}>{getCategoryLabel(category)}</text>
              </box>

              <solid_js_1.For each={bindings}>
                {function (binding) { return (<box flexDirection="row" padding={0.5}>
                    <box width={20}>
                      <text fg="yellow">{binding.key}</text>
                    </box>
                    <box flexGrow={1}>
                      <text fg="white">{binding.description}</text>
                    </box>
                  </box>); }}
              </solid_js_1.For>
            </box>);
        }}
        </solid_js_1.For>
      </box>

      <box flexDirection="row" marginTop={1} gap={2}>
        <text fg="gray">Total: {keybindings_1.KEYBINDINGS.length} shortcuts</text>
        <box flexGrow={1}/>
        <text fg="gray">? anytime for help</text>
      </box>
    </box>);
}
