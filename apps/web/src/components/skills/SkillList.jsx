"use strict";
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
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
exports.default = SkillList;
var react_1 = require("react");
var SkillCard_1 = require("./SkillCard");
function SkillList() {
    var _this = this;
    var _a = (0, react_1.useState)([]), skills = _a[0], setSkills = _a[1];
    var _b = (0, react_1.useState)(true), loading = _b[0], setLoading = _b[1];
    var _c = (0, react_1.useState)(null), error = _c[0], setError = _c[1];
    (0, react_1.useEffect)(function () {
        var fetchSkills = function () { return __awaiter(_this, void 0, void 0, function () {
            var res, data, mapped, err_1;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        _a.trys.push([0, 3, 4, 5]);
                        return [4 /*yield*/, fetch('http://localhost:8080/skills')];
                    case 1:
                        res = _a.sent();
                        if (!res.ok)
                            throw new Error('Failed to fetch skills');
                        return [4 /*yield*/, res.json()];
                    case 2:
                        data = _a.sent();
                        mapped = (data.skills || []).map(function (s) {
                            var _a;
                            return ({
                                id: s.id,
                                name: s.name,
                                description: s.description,
                                emoji: ((_a = s.metadata) === null || _a === void 0 ? void 0 : _a.emoji) || '🧩', // Fallback emoji
                                enabled: s.enabled,
                                source: s.source,
                                path: s.path,
                            });
                        });
                        setSkills(mapped);
                        return [3 /*break*/, 5];
                    case 3:
                        err_1 = _a.sent();
                        console.error(err_1);
                        setError('Could not load skills from Runtime API. Is pryx-core running?');
                        return [3 /*break*/, 5];
                    case 4:
                        setLoading(false);
                        return [7 /*endfinally*/];
                    case 5: return [2 /*return*/];
                }
            });
        }); };
        fetchSkills();
    }, []);
    if (loading)
        return <div style={{ padding: '2rem', color: '#6b7280' }}>Loading skills...</div>;
    if (error)
        return <div style={{ padding: '2rem', color: '#ef4444' }}>Error: {error}</div>;
    var handleToggle = function (id) {
        setSkills(function (prev) { return prev.map(function (s) {
            return s.id === id ? __assign(__assign({}, s), { enabled: !s.enabled }) : s;
        }); });
    };
    return (<div style={{ padding: '1rem', fontFamily: 'system-ui', maxWidth: '1200px', margin: '0 auto' }}>
            <header style={{ marginBottom: '2rem' }}>
                <h1 style={{ margin: 0, fontSize: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                    <span>🛠️</span> Skills Manager
                </h1>
                <p style={{ color: '#9ca3af', marginTop: '0.5rem' }}>
                    Manage the capabilities and tools available to your Pryx agent.
                </p>
            </header>

            <div style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))',
            gap: '1.5rem'
        }}>
                {skills.map(function (skill) { return (<SkillCard_1.default key={skill.id} skill={skill} onToggle={handleToggle}/>); })}
            </div>
        </div>);
}
