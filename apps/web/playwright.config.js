"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var test_1 = require("@playwright/test");
exports.default = (0, test_1.defineConfig)({
    testDir: './e2e',
    fullyParallel: true,
    webServer: {
        command: 'bun run dev -- --port 4321 --host 127.0.0.1',
        port: 4321,
        reuseExistingServer: !process.env.CI,
        timeout: 120000,
    },
    use: {
        baseURL: 'http://127.0.0.1:4321',
    },
});
