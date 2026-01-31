"use strict";
var __spreadArray = (this && this.__spreadArray) || function (to, from, pack) {
    if (pack || arguments.length === 2) for (var i = 0, l = from.length, ar; i < l; i++) {
        if (ar || !(i in from)) {
            if (!ar) ar = Array.prototype.slice.call(from, 0, i);
            ar[i] = from[i];
        }
    }
    return to.concat(ar || Array.prototype.slice.call(from));
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = Dashboard;
var react_1 = require("react");
var DeviceList_1 = require("./dashboard/DeviceList");
function Dashboard() {
    var _a = (0, react_1.useState)([]), events = _a[0], setEvents = _a[1];
    var _b = (0, react_1.useState)(null), stats = _b[0], setStats = _b[1];
    var _c = (0, react_1.useState)(null), selectedEvent = _c[0], setSelectedEvent = _c[1];
    var _d = (0, react_1.useState)(false), connected = _d[0], setConnected = _d[1];
    (0, react_1.useEffect)(function () {
        var ws = new WebSocket('ws://localhost:3000/ws');
        ws.onopen = function () { return setConnected(true); };
        ws.onclose = function () { return setConnected(false); };
        ws.onmessage = function (msg) {
            try {
                var evt_1 = JSON.parse(msg.data);
                if (evt_1.event === 'trace.event') {
                    setEvents(function (prev) { return __spreadArray(__spreadArray([], prev.slice(-99), true), [evt_1.payload], false); });
                }
                else if (evt_1.event === 'session.stats') {
                    setStats(evt_1.payload);
                }
            }
            catch ( /* ignore */_a) { /* ignore */ }
        };
        return function () { return ws.close(); };
    }, []);
    var formatDuration = function (ms) {
        if (!ms)
            return '-';
        if (ms < 1000)
            return "".concat(ms, "ms");
        return "".concat((ms / 1000).toFixed(2), "s");
    };
    var formatCost = function (cost) { return "$".concat(cost.toFixed(4)); };
    var getEventColor = function (event) {
        if (event.status === 'error')
            return '#ef4444';
        if (event.status === 'running')
            return '#f59e0b';
        switch (event.type) {
            case 'tool_call': return '#3b82f6';
            case 'approval': return '#8b5cf6';
            case 'message': return '#10b981';
            default: return '#6b7280';
        }
    };
    var timelineWidth = 600;
    var now = Date.now();
    var timeWindow = 60000; // 60s window
    return (<div style={{ padding: '1rem', fontFamily: 'system-ui', maxWidth: '1200px', margin: '0 auto' }}>
            <header style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '2rem' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '1.5rem' }}>
                    <h1 style={{ margin: 0, fontSize: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                        <span>⚡</span> Pryx Cloud
                    </h1>
                    <nav style={{ display: 'flex', gap: '1rem', fontSize: '0.9rem' }}>
                        <a href="/dashboard" style={{ color: '#fff', textDecoration: 'none', fontWeight: 'bold' }}>Dashboard</a>
                        <a href="/skills" style={{ color: '#9ca3af', textDecoration: 'none' }}>Skills</a>
                    </nav>
                </div>
                <span style={{
            color: connected ? '#10b981' : '#ef4444',
            display: 'flex',
            alignItems: 'center',
            gap: '0.5rem'
        }}>
                    <span style={{
            width: 8,
            height: 8,
            borderRadius: '50%',
            backgroundColor: connected ? '#10b981' : '#ef4444'
        }}/>
                    {connected ? 'Live' : 'Offline'}
                </span>
            </header>

            <DeviceList_1.default />

            {stats && (<div style={{
                display: 'grid',
                gridTemplateColumns: 'repeat(4, 1fr)',
                gap: '1rem',
                marginBottom: '1.5rem'
            }}>
                        <StatCard label="Cost" value={formatCost(stats.cost)} color="#f59e0b"/>
                        <StatCard label="Tokens" value={stats.tokens.toLocaleString()} color="#3b82f6"/>
                        <StatCard label="Duration" value={formatDuration(stats.duration)} color="#10b981"/>
                        <StatCard label="Events" value={stats.eventCount.toString()} color="#8b5cf6"/>
                    </div>)}

            <section>
                <h2 style={{ fontSize: '1rem', marginBottom: '0.5rem' }}>Trace Timeline</h2>
                <div style={{
            border: '1px solid #333',
            borderRadius: 8,
            padding: '1rem',
            backgroundColor: '#111'
        }}>
                    <svg width={timelineWidth} height={Math.max(events.length * 24 + 20, 100)}>
                        {events.map(function (event, i) {
            var start = Math.max(0, (event.startTime - (now - timeWindow)) / timeWindow);
            var end = event.endTime
                ? Math.min(1, (event.endTime - (now - timeWindow)) / timeWindow)
                : 1;
            var x = start * timelineWidth;
            var width = Math.max(4, (end - start) * timelineWidth);
            return (<g key={event.id} onClick={function () { return setSelectedEvent(event); }}>
                                    <rect x={x} y={i * 24 + 4} width={width} height={18} rx={4} fill={getEventColor(event)} style={{ cursor: 'pointer' }}/>
                                    <text x={x + 4} y={i * 24 + 16} fontSize={10} fill="#fff">
                                        {event.name.slice(0, 20)}
                                    </text>
                                </g>);
        })}
                    </svg>
                </div>
            </section>

            {selectedEvent && (<section style={{ marginTop: '1rem' }}>
                        <h2 style={{ fontSize: '1rem', marginBottom: '0.5rem' }}>Event Details</h2>
                        <div style={{
                border: '1px solid #333',
                borderRadius: 8,
                padding: '1rem',
                backgroundColor: '#111',
                fontSize: '0.875rem'
            }}>
                            <p><strong>ID:</strong> {selectedEvent.id}</p>
                            <p><strong>Type:</strong> {selectedEvent.type}</p>
                            <p><strong>Name:</strong> {selectedEvent.name}</p>
                            <p><strong>Status:</strong> {selectedEvent.status}</p>
                            <p><strong>Duration:</strong> {formatDuration(selectedEvent.duration)}</p>
                            {selectedEvent.correlationId && (<p><strong>Correlation ID:</strong> {selectedEvent.correlationId}</p>)}
                            {selectedEvent.error && (<p style={{ color: '#ef4444' }}><strong>Error:</strong> {selectedEvent.error}</p>)}
                        </div>
                    </section>)}

            <section style={{ marginTop: '1rem' }}>
                <h2 style={{ fontSize: '1rem', marginBottom: '0.5rem' }}>Recent Events</h2>
                <div style={{
            border: '1px solid #333',
            borderRadius: 8,
            overflow: 'hidden',
            backgroundColor: '#111'
        }}>
                    {events.slice(-10).reverse().map(function (event) { return (<div key={event.id} style={{
                padding: '0.5rem 1rem',
                borderBottom: '1px solid #222',
                display: 'flex',
                alignItems: 'center',
                gap: '0.5rem',
                fontSize: '0.875rem'
            }}>
                            <span style={{
                width: 8,
                height: 8,
                borderRadius: '50%',
                backgroundColor: getEventColor(event)
            }}/>
                            <span style={{ color: '#9ca3af' }}>{event.type}</span>
                            <span>{event.name}</span>
                            <span style={{ marginLeft: 'auto', color: '#6b7280' }}>
                                {formatDuration(event.duration)}
                            </span>
                        </div>); })}
                </div>
            </section>
        </div>);
}
function StatCard(_a) {
    var label = _a.label, value = _a.value, color = _a.color;
    return (<div style={{
            border: '1px solid #333',
            borderRadius: 8,
            padding: '1rem',
            backgroundColor: '#111'
        }}>
            <div style={{ fontSize: '0.75rem', color: '#9ca3af', marginBottom: '0.25rem' }}>{label}</div>
            <div style={{ fontSize: '1.5rem', fontWeight: 'bold', color: color }}>{value}</div>
        </div>);
}
