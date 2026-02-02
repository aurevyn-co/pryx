import { createSignal, onCleanup, onMount } from "solid-js";
import type { Accessor } from "solid-js";
import { Effect, Fiber, Runtime, Stream, Context, ManagedRuntime, Layer } from "effect";
import { WebSocketServiceLive } from "../services/ws";
import { HealthCheckServiceLive } from "../services/health-check";
import { ProviderServiceLive } from "../services/provider-service";
import { SkillsServiceLive } from "../services/skills-api";

// Create a managed runtime that includes our Live services
export const AppRuntime = ManagedRuntime.make(
  Layer.mergeAll(
    WebSocketServiceLive,
    HealthCheckServiceLive,
    ProviderServiceLive,
    SkillsServiceLive
  )
);

/**
 * Run an Effect and expose result as SolidJS signal
 */
export function useEffectSignal<A, E = never>(
  effect: Effect.Effect<A, E>
): Accessor<A | undefined> {
  const [value, setValue] = createSignal<A | undefined>();
  const [error, setError] = createSignal<E | undefined>();

  onMount(() => {
    // Run with our managed runtime
    AppRuntime.runFork(
      effect.pipe(
        Effect.tap(a => Effect.sync(() => setValue(() => a))),
        Effect.tapError(e => Effect.sync(() => setError(() => e)))
      )
    );
  });

  return value;
}

/**
 * Subscribe to an Effect Stream as SolidJS signal
 */
export function useEffectStream<A, E = never>(stream: Stream.Stream<A, E>): Accessor<A[]> {
  const [items, setItems] = createSignal<A[]>([]);

  onMount(() => {
    const fiber = AppRuntime.runFork(
      stream.pipe(Stream.runForEach(item => Effect.sync(() => setItems(prev => [...prev, item]))))
    );

    onCleanup(() => {
      Effect.runFork(Fiber.interrupt(fiber));
    });
  });

  return items;
}

/**
 * Access the WebSocketService
 */
export function useEffectService<I, S>(tag: Context.Tag<I, S>): Accessor<S | undefined> {
  const [service, setService] = createSignal<S | undefined>();

  onMount(() => {
    try {
      const svc = AppRuntime.runSync(tag as any);
      setService(() => svc as S);
    } catch (err) {
      console.error("Failed to get service:", err);
    }
  });

  return service;
}

// Global runtime for ad-hoc usage
export const TUIRuntime = Runtime.defaultRuntime;
