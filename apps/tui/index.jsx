"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var solid_1 = require("@opentui/solid");
var App_1 = require("./src/components/App");
process.on("SIGINT", function () {
    process.exit(0);
});
try {
    (0, solid_1.render)(function () { return <App_1.default />; }, {
        targetFps: 60,
        exitOnCtrlC: false,
    });
}
catch (e) {
    console.error("Failed to start TUI:", e);
    var fs = require("fs");
    fs.writeFileSync("tui-crash.log", String(e) + "\n" + (e instanceof Error ? e.stack : ""));
    process.exit(1);
}
