"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var bun_test_1 = require("bun:test");
var ws_1 = require("./ws");
(0, bun_test_1.describe)("WebSocket Service", function () {
    (0, bun_test_1.describe)("ConnectionError", function () {
        (0, bun_test_1.test)("should create error with message", function () {
            var error = new ws_1.ConnectionError("Test error message");
            (0, bun_test_1.expect)(error._tag).toBe("ConnectionError");
            (0, bun_test_1.expect)(error.message).toBe("Test error message");
            (0, bun_test_1.expect)(error.originalError).toBeUndefined();
        });
        (0, bun_test_1.test)("should create error with original error", function () {
            var original = new Error("Original error");
            var error = new ws_1.ConnectionError("Wrapped error", original);
            (0, bun_test_1.expect)(error.originalError).toBe(original);
        });
    });
    (0, bun_test_1.describe)("RuntimeEvent", function () {
        (0, bun_test_1.test)("should parse valid runtime events", function () {
            var testEvent = {
                event: "trace",
                type: "test",
                session_id: "test-session",
                payload: { message: "test" },
            };
            (0, bun_test_1.expect)(testEvent.event).toBe("trace");
            (0, bun_test_1.expect)(testEvent.session_id).toBe("test-session");
            (0, bun_test_1.expect)(testEvent.payload).toHaveProperty("message");
        });
        (0, bun_test_1.test)("should handle events without optional fields", function () {
            var minimalEvent = {
                event: "trace",
            };
            (0, bun_test_1.expect)(minimalEvent.event).toBe("trace");
            (0, bun_test_1.expect)(minimalEvent.session_id).toBeUndefined();
            (0, bun_test_1.expect)(minimalEvent.payload).toBeUndefined();
        });
    });
    (0, bun_test_1.describe)("ConnectionStatus", function () {
        (0, bun_test_1.test)("should define disconnected state", function () {
            var state = { _tag: "Disconnected" };
            (0, bun_test_1.expect)(state._tag).toBe("Disconnected");
        });
        (0, bun_test_1.test)("should define connecting state", function () {
            var state = { _tag: "Connecting" };
            (0, bun_test_1.expect)(state._tag).toBe("Connecting");
        });
        (0, bun_test_1.test)("should define connected state", function () {
            var state = { _tag: "Connected" };
            (0, bun_test_1.expect)(state._tag).toBe("Connected");
        });
        (0, bun_test_1.test)("should define error state", function () {
            var error = new ws_1.ConnectionError("Test");
            var state = { _tag: "Error", error: error };
            (0, bun_test_1.expect)(state._tag).toBe("Error");
            (0, bun_test_1.expect)(state.error).toBe(error);
        });
    });
});
