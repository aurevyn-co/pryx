"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g = Object.create((typeof Iterator === "function" ? Iterator : Object).prototype);
    return g.next = verb(0), g["throw"] = verb(1), g["return"] = verb(2), typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (g && (g = 0, op[0] && (_ = 0)), _) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
Object.defineProperty(exports, "__esModule", { value: true });
var vitest_1 = require("vitest");
var storage_js_1 = require("../../src/storage.js");
var registry_js_1 = require("../../src/registry.js");
var fs_1 = require("fs");
var path_1 = require("path");
var os_1 = require("os");
(0, vitest_1.describe)('MCPStorage Additional Coverage', function () {
    var tempDir;
    var storage;
    var registry;
    var configPath;
    (0, vitest_1.beforeEach)(function () {
        tempDir = fs_1.default.mkdtempSync(path_1.default.join(os_1.default.tmpdir(), 'mcp-storage-test-'));
        configPath = path_1.default.join(tempDir, 'mcp-servers.json');
        storage = new storage_js_1.MCPStorage();
        registry = new registry_js_1.MCPRegistry();
    });
    (0, vitest_1.afterEach)(function () {
        fs_1.default.rmSync(tempDir, { recursive: true, force: true });
    });
    (0, vitest_1.describe)('exists', function () {
        (0, vitest_1.it)('should return true when file exists', function () { return __awaiter(void 0, void 0, void 0, function () {
            var exists;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        fs_1.default.writeFileSync(configPath, JSON.stringify({ version: 1, servers: [] }));
                        return [4 /*yield*/, storage.exists(configPath)];
                    case 1:
                        exists = _a.sent();
                        (0, vitest_1.expect)(exists).toBe(true);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should return false when file does not exist', function () { return __awaiter(void 0, void 0, void 0, function () {
            var exists;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, storage.exists(configPath)];
                    case 1:
                        exists = _a.sent();
                        (0, vitest_1.expect)(exists).toBe(false);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('save', function () {
        (0, vitest_1.it)('should create parent directories if they do not exist', function () { return __awaiter(void 0, void 0, void 0, function () {
            var nestedPath;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        nestedPath = path_1.default.join(tempDir, 'nested', 'deep', 'mcp-servers.json');
                        registry.addServer({
                            id: 'test',
                            name: 'Test',
                            enabled: true,
                            source: 'manual',
                            transport: { type: 'stdio', command: 'test', args: [], env: {} },
                            settings: {
                                autoConnect: true,
                                timeout: 30000,
                                reconnect: true,
                                maxReconnectAttempts: 3,
                                fallbackServers: [],
                            },
                        });
                        return [4 /*yield*/, storage.save(nestedPath, registry)];
                    case 1:
                        _a.sent();
                        (0, vitest_1.expect)(fs_1.default.existsSync(nestedPath)).toBe(true);
                        return [2 /*return*/];
                }
            });
        }); });
    });
    (0, vitest_1.describe)('load', function () {
        (0, vitest_1.it)('should handle empty servers array', function () { return __awaiter(void 0, void 0, void 0, function () {
            var loaded;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        fs_1.default.writeFileSync(configPath, JSON.stringify({ version: 1, servers: [] }));
                        return [4 /*yield*/, storage.load(configPath)];
                    case 1:
                        loaded = _a.sent();
                        (0, vitest_1.expect)(loaded.size).toBe(0);
                        return [2 /*return*/];
                }
            });
        }); });
        (0, vitest_1.it)('should handle JSON parse error', function () { return __awaiter(void 0, void 0, void 0, function () {
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        fs_1.default.writeFileSync(configPath, '{ invalid json');
                        return [4 /*yield*/, (0, vitest_1.expect)(storage.load(configPath)).rejects.toThrow()];
                    case 1:
                        _a.sent();
                        return [2 /*return*/];
                }
            });
        }); });
    });
});
