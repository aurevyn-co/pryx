"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var solid_1 = require("@opentui/solid");
var App_1 = require("../../src/components/App");
// Force TTY for library checks
Object.defineProperty(process.stdout, "isTTY", { value: true });
Object.defineProperty(process.stdout, "columns", { value: 80 });
Object.defineProperty(process.stdout, "rows", { value: 24 });
console.log("Starting Render Script...");
try {
    (0, solid_1.render)(function () { return <App_1.default />; });
    // Allow time to render frames
    setTimeout(function () {
        console.log("Render timeout reached. Exiting.");
        process.exit(0);
    }, 3000);
}
catch (e) {
    console.error("Render failed:", e);
    process.exit(1);
}
