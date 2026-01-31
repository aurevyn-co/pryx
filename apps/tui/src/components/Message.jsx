"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = Message;
var typeColors = {
    user: "cyan",
    assistant: "white",
    tool: "yellow",
    approval: "magenta",
    system: "gray",
    thinking: "gray",
};
var typePrefixes = {
    user: "You",
    assistant: "Pryx",
    tool: "⚙️",
    approval: "⚠️ Approval",
    system: "ℹ️",
    thinking: "💭",
};
function Message(props) {
    var color = typeColors[props.type];
    var prefix = typePrefixes[props.type];
    if (props.type === "tool" && props.toolName) {
        var statusIcon = props.toolStatus === "running"
            ? "⏳"
            : props.toolStatus === "done"
                ? "✅"
                : props.toolStatus === "error"
                    ? "❌"
                    : "⚙️";
        return (<box>
        <text fg="yellow">
          {statusIcon} {props.toolName}:{" "}
        </text>
        <text fg="gray">{props.content}</text>
      </box>);
    }
    if (props.type === "thinking") {
        return (<box borderStyle="single" borderColor="gray" padding={1}>
        <text fg="gray">_Thinking:_ {props.content}</text>
      </box>);
    }
    return (<box>
      <text fg={color}>{prefix}: </text>
      <text fg={color}>{props.content}</text>
      {props.pending && <text fg="gray"> ▌</text>}
    </box>);
}
