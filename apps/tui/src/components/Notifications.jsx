"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = Notifications;
// @ts-nocheck
var solid_js_1 = require("solid-js");
function Notifications(props) {
    var getColor = function (type) {
        switch (type) {
            case "success":
                return "green";
            case "error":
                return "red";
            case "warning":
                return "yellow";
            default:
                return "blue";
        }
    };
    return (<box flexDirection="column" position="absolute" top={1} right={1} width={40}>
      <solid_js_1.For each={props.items}>
        {function (item) { return (<box borderStyle="single" borderColor={getColor(item.type)} padding={1} marginBottom={1} flexDirection="column">
            <text color={getColor(item.type)} bold>
              {item.type.toUpperCase()}
            </text>
            <text>{item.message}</text>
          </box>); }}
      </solid_js_1.For>
    </box>);
}
