import { palette } from "../theme";

export default function AppHeader() {
  return (
    <box flexDirection="column" alignItems="center" padding={1} backgroundColor={palette.bgPrimary}>
      <box flexDirection="row">
        <text fg={palette.accent}>{"    ██████╗ ██████╗ ██╗   ██╗██╗  ██╗    "}</text>
      </box>
      <box flexDirection="row">
        <text fg={palette.accent}>{"    ██╔══██╗██╔══██╗╚██╗ ██╔╝╚██╗██╔╝    "}</text>
      </box>
      <box flexDirection="row">
        <text fg={palette.accent}>{"    ██████╔╝██████╔╝ ╚████╔╝  ╚███╔╝     "}</text>
      </box>
      <box flexDirection="row">
        <text fg={palette.accent}>{"    ██╔═══╝ ██╔══██╗  ╚██╔╝   ██╔██╗     "}</text>
      </box>
      <box flexDirection="row">
        <text fg={palette.accent}>{"    ██║     ██║  ██║   ██║   ██╔╝ ██╗    "}</text>
      </box>
      <box flexDirection="row">
        <text fg={palette.accent}>{"    ╚═╝     ╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝    "}</text>
      </box>
      <box marginTop={1}>
        <text fg={palette.dim}>Autonomous AI Agent for Any Task</text>
      </box>
    </box>
  );
}
