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
exports.TUIRuntime = exports.AppRuntime = void 0;
exports.useEffectSignal = useEffectSignal;
exports.useEffectStream = useEffectStream;
exports.useEffectService = useEffectService;
var solid_js_1 = require("solid-js");
var effect_1 = require("effect");
var ws_1 = require("../services/ws");
var fs_1 = require("fs");
function log(msg) {
    (0, fs_1.appendFileSync)("debug.log", "[hooks] ".concat(msg, "\n"));
    // console.error(`[hooks] ${msg}`); // Fallback
}
log("MODULE LOADED");
// Create a managed runtime that includes our Live services
exports.AppRuntime = effect_1.ManagedRuntime.make(ws_1.WebSocketServiceLive);
/**
 * Run an Effect and expose result as SolidJS signal
 */
function useEffectSignal(effect) {
    var _a = (0, solid_js_1.createSignal)(), value = _a[0], setValue = _a[1];
    var _b = (0, solid_js_1.createSignal)(), error = _b[0], setError = _b[1];
    (0, solid_js_1.onMount)(function () {
        // Run with our managed runtime
        exports.AppRuntime.runFork(effect.pipe(effect_1.Effect.tap(function (a) { return effect_1.Effect.sync(function () { return setValue(function () { return a; }); }); }), effect_1.Effect.tapError(function (e) { return effect_1.Effect.sync(function () { return setError(function () { return e; }); }); })));
    });
    return value;
}
/**
 * Subscribe to an Effect Stream as SolidJS signal
 */
function useEffectStream(stream) {
    var _a = (0, solid_js_1.createSignal)([]), items = _a[0], setItems = _a[1];
    (0, solid_js_1.onMount)(function () {
        var fiber = exports.AppRuntime.runFork(stream.pipe(effect_1.Stream.runForEach(function (item) { return effect_1.Effect.sync(function () { return setItems(function (prev) { return __spreadArray(__spreadArray([], prev, true), [item], false); }); }); })));
        (0, solid_js_1.onCleanup)(function () {
            effect_1.Effect.runFork(effect_1.Fiber.interrupt(fiber));
        });
    });
    return items;
}
/**
 * Access the WebSocketService
 */
function useEffectService(tag) {
    log("useEffectService: start");
    try {
        log("useEffectService: calling createSignal");
        var _a = (0, solid_js_1.createSignal)(), service = _a[0], setService_1 = _a[1];
        log("useEffectService: createSignal done");
        log("useEffectService: scheduling onMount");
        (0, solid_js_1.onMount)(function () {
            log("useEffectService: onMount running");
            // Run an effect to extract the service
            // This is safe because we use ManagedRuntime which keeps services alive
            exports.AppRuntime.runPromise(tag)
                .then(function (svc) {
                log("useEffectService: service resolved");
                setService_1(function () { return svc; });
            })
                .catch(function (err) {
                log("useEffectService: service error ".concat(err));
                console.error("Failed to get service:", err);
            });
        });
        log("useEffectService: onMount scheduled");
        return service;
    }
    catch (e) {
        log("useEffectService: CRASHED ".concat(e));
        throw e;
    }
}
// Global runtime for ad-hoc usage
exports.TUIRuntime = effect_1.Runtime.defaultRuntime;
