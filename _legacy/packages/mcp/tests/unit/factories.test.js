"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var vitest_1 = require("vitest");
var storage_js_1 = require("../../src/storage.js");
var registry_js_1 = require("../../src/registry.js");
(0, vitest_1.describe)('Factory Functions', function () {
    (0, vitest_1.describe)('createStorage', function () {
        (0, vitest_1.it)('should create a new MCPStorage instance', function () {
            var storage = (0, storage_js_1.createStorage)();
            (0, vitest_1.expect)(storage).toBeInstanceOf(storage_js_1.MCPStorage);
        });
    });
    (0, vitest_1.describe)('createRegistry', function () {
        (0, vitest_1.it)('should create a new MCPRegistry instance', function () {
            var registry = (0, registry_js_1.createRegistry)();
            (0, vitest_1.expect)(registry).toBeInstanceOf(registry_js_1.MCPRegistry);
            (0, vitest_1.expect)(registry.size).toBe(0);
        });
    });
});
