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
exports.ALL = void 0;
var hono_1 = require("hono");
/**
 * Unified API route for Pryx Cloud using Hono
 * Ported from vanilla Response logic for better scalability
 */
var app = new hono_1.Hono().basePath('/api');
// --- Constants & Utilities ---
var ALLOWED_TELEMETRY_FIELDS = new Set([
    'correlation_id', 'timestamp', 'level', 'category', 'error_code', 'error_message',
    'duration_ms', 'model_id', 'token_count', 'cost_usd', 'tool_name', 'status',
    'device_id', 'session_id', 'version',
]);
var PII_PATTERNS = [
    /\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b/g,
    /\b(?:sk|pk|api)[-_][a-zA-Z0-9]{20,}\b/g,
    /\b[0-9]{4}[- ]?[0-9]{4}[- ]?[0-9]{4}[- ]?[0-9]{4}\b/g,
    /\b[0-9]{3}[-.]?[0-9]{3}[-.]?[0-9]{4}\b/g,
];
function generateCode(length, charset) {
    var array = new Uint8Array(length);
    crypto.getRandomValues(array);
    return Array.from(array, function (b) { return charset[b % charset.length]; }).join('');
}
function redactPII(value) {
    var result = value;
    for (var _i = 0, PII_PATTERNS_1 = PII_PATTERNS; _i < PII_PATTERNS_1.length; _i++) {
        var pattern = PII_PATTERNS_1[_i];
        result = result.replace(pattern, '[REDACTED]');
    }
    return result;
}
// --- Middleware ---
app.use('*', function (c, next) { return __awaiter(void 0, void 0, void 0, function () {
    var ip, env, success;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0:
                ip = c.req.header('CF-Connecting-IP') || 'unknown';
                env = c.env;
                if (!(env === null || env === void 0 ? void 0 : env.RATE_LIMITER)) return [3 /*break*/, 2];
                return [4 /*yield*/, env.RATE_LIMITER.limit(ip)];
            case 1:
                success = (_a.sent()).success;
                if (!success) {
                    return [2 /*return*/, c.json({ error: 'slow_down' }, 429)];
                }
                _a.label = 2;
            case 2: return [4 /*yield*/, next()];
            case 3:
                _a.sent();
                return [2 /*return*/];
        }
    });
}); });
// --- Auth Routes ---
app.post('/auth/qr/pairing', function (c) { return __awaiter(void 0, void 0, void 0, function () {
    var body, deviceId, pairingCode, pairingToken, entry;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0: return [4 /*yield*/, c.req.json().catch(function () { return ({}); })];
            case 1:
                body = _a.sent();
                deviceId = body.device_id || '';
                pairingCode = generateCode(8, 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789');
                pairingToken = generateCode(40, 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789');
                entry = {
                    deviceId: deviceId,
                    pairingCode: pairingCode,
                    status: 'pending',
                    created_at: Date.now(),
                };
                // Store with 5-minute TTL
                return [4 /*yield*/, c.env.DEVICE_CODES.put("qr:".concat(pairingCode), JSON.stringify(entry), { expirationTtl: 300 })];
            case 2:
                // Store with 5-minute TTL
                _a.sent();
                return [4 /*yield*/, c.env.DEVICE_CODES.put("ptr:".concat(pairingToken), pairingCode, { expirationTtl: 300 })];
            case 3:
                _a.sent();
                return [2 /*return*/, c.json({
                        pairing_code: pairingCode,
                        pairing_token: pairingToken,
                        expires_in: 300,
                    })];
        }
    });
}); });
app.get('/auth/qr/status', function (c) { return __awaiter(void 0, void 0, void 0, function () {
    var token, pairingCode, entryStr, entry;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0:
                token = c.req.query('token');
                if (!token)
                    return [2 /*return*/, c.json({ error: 'missing_token' }, 400)];
                return [4 /*yield*/, c.env.DEVICE_CODES.get("ptr:".concat(token))];
            case 1:
                pairingCode = _a.sent();
                if (!pairingCode)
                    return [2 /*return*/, c.json({ error: 'expired_token' }, 400)];
                return [4 /*yield*/, c.env.DEVICE_CODES.get("qr:".concat(pairingCode))];
            case 2:
                entryStr = _a.sent();
                if (!entryStr)
                    return [2 /*return*/, c.json({ error: 'expired_pairing' }, 400)];
                entry = JSON.parse(entryStr);
                return [2 /*return*/, c.json(entry)];
        }
    });
}); });
app.post('/auth/device/code', function (c) { return __awaiter(void 0, void 0, void 0, function () {
    var body, deviceId, scopeStr, deviceCode, userCode, entry;
    var _a, _b;
    return __generator(this, function (_c) {
        switch (_c.label) {
            case 0: return [4 /*yield*/, c.req.formData().catch(function () { return new FormData(); })];
            case 1:
                body = _c.sent();
                deviceId = ((_a = body.get('device_id')) === null || _a === void 0 ? void 0 : _a.toString()) || '';
                scopeStr = ((_b = body.get('scope')) === null || _b === void 0 ? void 0 : _b.toString()) || 'telemetry.write';
                deviceCode = generateCode(40, 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789');
                userCode = "".concat(generateCode(4, 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'), "-").concat(generateCode(4, 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'));
                entry = {
                    user_code: userCode,
                    device_id: deviceId,
                    scopes: scopeStr.split(' '),
                    created_at: Date.now(),
                    expires_at: Date.now() + 600 * 1000,
                    authorized: false,
                };
                return [4 /*yield*/, c.env.DEVICE_CODES.put(deviceCode, JSON.stringify(entry), { expirationTtl: 660 })];
            case 2:
                _c.sent();
                return [4 /*yield*/, c.env.DEVICE_CODES.put("user:".concat(userCode), deviceCode, { expirationTtl: 660 })];
            case 3:
                _c.sent();
                return [2 /*return*/, c.json({
                        device_code: deviceCode,
                        user_code: userCode,
                        verification_uri: '/link',
                        verification_uri_complete: "/link?code=".concat(userCode),
                        expires_in: 600,
                        interval: 5,
                    })];
        }
    });
}); });
app.post('/auth/device/token', function (c) { return __awaiter(void 0, void 0, void 0, function () {
    var body, deviceCode, entryStr, entry, accessToken, refreshToken;
    var _a;
    return __generator(this, function (_b) {
        switch (_b.label) {
            case 0: return [4 /*yield*/, c.req.formData().catch(function () { return new FormData(); })];
            case 1:
                body = _b.sent();
                deviceCode = ((_a = body.get('device_code')) === null || _a === void 0 ? void 0 : _a.toString()) || '';
                return [4 /*yield*/, c.env.DEVICE_CODES.get(deviceCode)];
            case 2:
                entryStr = _b.sent();
                if (!entryStr)
                    return [2 /*return*/, c.json({ error: 'expired_token' }, 400)];
                entry = JSON.parse(entryStr);
                if (!entry.authorized)
                    return [2 /*return*/, c.json({ error: 'authorization_pending' }, 400)];
                accessToken = "pryx_at_".concat(generateCode(40, 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'));
                refreshToken = "pryx_rt_".concat(generateCode(40, 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'));
                return [4 /*yield*/, c.env.TOKENS.put(accessToken, JSON.stringify(entry), { expirationTtl: 3600 })];
            case 3:
                _b.sent();
                return [4 /*yield*/, c.env.TOKENS.put("refresh:".concat(refreshToken), JSON.stringify(entry), { expirationTtl: 86400 * 30 })];
            case 4:
                _b.sent();
                return [2 /*return*/, c.json({
                        access_token: accessToken,
                        token_type: 'Bearer',
                        expires_in: 3600,
                        refresh_token: refreshToken,
                        scope: entry.scopes.join(' '),
                    })];
        }
    });
}); });
app.post('/auth/token/refresh', function (c) { return __awaiter(void 0, void 0, void 0, function () {
    var body, refreshToken, entryStr, entry, newAccessToken;
    var _a;
    return __generator(this, function (_b) {
        switch (_b.label) {
            case 0: return [4 /*yield*/, c.req.formData().catch(function () { return new FormData(); })];
            case 1:
                body = _b.sent();
                refreshToken = ((_a = body.get('refresh_token')) === null || _a === void 0 ? void 0 : _a.toString()) || '';
                return [4 /*yield*/, c.env.TOKENS.get("refresh:".concat(refreshToken))];
            case 2:
                entryStr = _b.sent();
                if (!entryStr)
                    return [2 /*return*/, c.json({ error: 'invalid_grant' }, 400)];
                entry = JSON.parse(entryStr);
                newAccessToken = "pryx_at_".concat(generateCode(40, 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'));
                return [4 /*yield*/, c.env.TOKENS.put(newAccessToken, JSON.stringify(entry), { expirationTtl: 3600 })];
            case 3:
                _b.sent();
                return [2 /*return*/, c.json({
                        access_token: newAccessToken,
                        token_type: 'Bearer',
                        expires_in: 3600,
                        refresh_token: refreshToken,
                        scope: entry.scopes.join(' '),
                    })];
        }
    });
}); });
// --- Telemetry Routes ---
app.post('/telemetry/ingest', function (c) { return __awaiter(void 0, void 0, void 0, function () {
    var body, events, sanitized, e_1;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0:
                _a.trys.push([0, 2, , 3]);
                return [4 /*yield*/, c.req.json()];
            case 1:
                body = _a.sent();
                events = Array.isArray(body) ? body : [body];
                sanitized = events.map(function (event) {
                    var clean = {};
                    for (var _i = 0, _a = Object.entries(event); _i < _a.length; _i++) {
                        var _b = _a[_i], key = _b[0], value = _b[1];
                        if (ALLOWED_TELEMETRY_FIELDS.has(key)) {
                            clean[key] = typeof value === 'string' ? redactPII(value) : value;
                        }
                    }
                    return clean;
                }).filter(function (e) { return e.correlation_id; });
                return [2 /*return*/, c.json({ accepted: sanitized.length })];
            case 2:
                e_1 = _a.sent();
                return [2 /*return*/, c.json({ error: 'Invalid JSON' }, 400)];
            case 3: return [2 /*return*/];
        }
    });
}); });
// --- Session / Mesh Routes ---
app.post('/sessions/broadcast', function (c) { return __awaiter(void 0, void 0, void 0, function () {
    var body, device_id, session_id, payload, timestamp, key, update;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0: return [4 /*yield*/, c.req.json().catch(function () { return ({}); })];
            case 1:
                body = _a.sent();
                device_id = body.device_id, session_id = body.session_id, payload = body.payload, timestamp = body.timestamp;
                if (!device_id || !session_id)
                    return [2 /*return*/, c.json({ error: 'missing_fields' }, 400)];
                key = "session:".concat(session_id);
                update = {
                    device_id: device_id,
                    payload: payload,
                    timestamp: timestamp || Date.now(),
                };
                // Store session update in KV
                return [4 /*yield*/, c.env.SESSIONS.put(key, JSON.stringify(update), { expirationTtl: 86400 })];
            case 2:
                // Store session update in KV
                _a.sent();
                return [2 /*return*/, c.json({ status: 'broadcasted', key: key })];
        }
    });
}); });
app.get('/sessions/:id', function (c) { return __awaiter(void 0, void 0, void 0, function () {
    var id, entry;
    return __generator(this, function (_a) {
        switch (_a.label) {
            case 0:
                id = c.req.param('id');
                return [4 /*yield*/, c.env.SESSIONS.get("session:".concat(id))];
            case 1:
                entry = _a.sent();
                if (!entry)
                    return [2 /*return*/, c.json({ error: 'not_found' }, 404)];
                return [2 /*return*/, c.json(JSON.parse(entry))];
        }
    });
}); });
// --- Update Routes ---
app.get('/update/manifest', function (c) { return __awaiter(void 0, void 0, void 0, function () {
    var platform, arch, manifest;
    return __generator(this, function (_a) {
        platform = c.req.query('platform') || 'darwin-aarch64';
        arch = c.req.query('arch') || 'aarch64';
        manifest = {
            version: '0.1.1',
            notes: 'Security fixes and Hono API refactor.',
            pub_date: new Date().toISOString(),
            platforms: {
                'darwin-aarch64': {
                    signature: '...', // Placeholder
                    url: "https://github.com/pryx-dev/pryx/releases/download/v0.1.1/Pryx_0.1.1_aarch64.app.tar.gz"
                },
                'darwin-x86_64': {
                    signature: '...',
                    url: "https://github.com/pryx-dev/pryx/releases/download/v0.1.1/Pryx_0.1.1_x64.app.tar.gz"
                },
                'windows-x86_64': {
                    signature: '...',
                    url: "https://github.com/pryx-dev/pryx/releases/download/v0.1.1/Pryx_0.1.1_x64_en-US.msi.zip"
                }
            }
        };
        return [2 /*return*/, c.json(manifest)];
    });
}); });
// --- Health ---
app.get('/', function (c) { return c.json({ name: 'Pryx Cloud API', status: 'operational', engine: 'hono' }); });
// --- Astro Integration ---
var ALL = function (_a) {
    var request = _a.request, locals = _a.locals;
    var runtime = locals.runtime;
    return app.fetch(request, runtime === null || runtime === void 0 ? void 0 : runtime.env);
};
exports.ALL = ALL;
