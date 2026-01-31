"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = AppHeader;
var theme_1 = require("../theme");
function AppHeader() {
    return (<box flexDirection="column" alignItems="center" padding={1} backgroundColor={theme_1.palette.bgPrimary}>
      <box flexDirection="row">
        <text fg={theme_1.palette.accent}>{"    ██████╗ ██████╗ ██╗   ██╗██╗  ██╗    "}</text>
      </box>
      <box flexDirection="row">
        <text fg={theme_1.palette.accent}>{"    ██╔══██╗██╔══██╗╚██╗ ██╔╝╚██╗██╔╝    "}</text>
      </box>
      <box flexDirection="row">
        <text fg={theme_1.palette.accent}>{"    ██████╔╝██████╔╝ ╚████╔╝  ╚███╔╝     "}</text>
      </box>
      <box flexDirection="row">
        <text fg={theme_1.palette.accent}>{"    ██╔═══╝ ██╔══██╗  ╚██╔╝   ██╔██╗     "}</text>
      </box>
      <box flexDirection="row">
        <text fg={theme_1.palette.accent}>{"    ██║     ██║  ██║   ██║   ██╔╝ ██╗    "}</text>
      </box>
      <box flexDirection="row">
        <text fg={theme_1.palette.accent}>{"    ╚═╝     ╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝    "}</text>
      </box>
      <box marginTop={1}>
        <text fg={theme_1.palette.dim}>Autonomous AI Agent for Any Task</text>
      </box>
    </box>);
}
