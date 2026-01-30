"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var vitest_1 = require("vitest");
var react_1 = require("@testing-library/react");
var Dashboard_1 = require("./Dashboard");
// Mock WebSocket
vitest_1.vi.stubGlobal('WebSocket', vitest_1.vi.fn(function () { return ({
    onopen: null,
    onclose: null,
    onmessage: null,
    close: vitest_1.vi.fn(),
}); }));
(0, vitest_1.describe)('Dashboard', function () {
    (0, vitest_1.it)('renders without crashing', function () {
        var container = (0, react_1.render)(<Dashboard_1.default />).container;
        (0, vitest_1.expect)(container).toBeDefined();
    });
    (0, vitest_1.it)('shows disconnected state initially', function () {
        var getByText = (0, react_1.render)(<Dashboard_1.default />).getByText;
        (0, vitest_1.expect)(getByText('Disconnected')).toBeDefined();
    });
    (0, vitest_1.it)('displays header', function () {
        var getByText = (0, react_1.render)(<Dashboard_1.default />).getByText;
        (0, vitest_1.expect)(getByText('Observability Dashboard')).toBeDefined();
    });
});
