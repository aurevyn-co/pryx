import { createSignal, createEffect, onMount, Show } from "solid-js";
import { Effect } from "effect";
import { useKeyboard, usePaste } from "@opentui/solid";
import { useEffectService, AppRuntime } from "../lib/hooks";
import {
  ProviderService,
  Provider as ProviderType,
  Model as ModelType,
} from "../services/provider-service";
import { saveConfig, AppConfig } from "../services/config";
import { getRuntimeHttpUrl } from "../services/skills-api";
import { isPrintable } from "../lib/keybindings";

interface SetupRequiredProps {
  onSetupComplete: () => void;
}

export default function SetupRequired(props: SetupRequiredProps) {
  const providerService = useEffectService(ProviderService);
  const [step, setStep] = createSignal(0);
  const [provider, setProvider] = createSignal("");
  const [modelName, setModelName] = createSignal("");
  const [error, setError] = createSignal("");
  const [providers, setProviders] = createSignal<ProviderType[]>([]);
  const [models, setModels] = createSignal<ModelType[]>([]);
  const [loading, setLoading] = createSignal(false);
  const [fetchError, setFetchError] = createSignal("");
  const [selectedProviderIndex, setSelectedProviderIndex] = createSignal(0);
  const [selectedModelIndex, setSelectedModelIndex] = createSignal(0);
  const [inputValue, setInputValue] = createSignal("");
  const [cursorPosition, setCursorPosition] = createSignal(0);
  const [cloudLoggedIn, setCloudLoggedIn] = createSignal(false);
  const [cloudLogin, setCloudLogin] = createSignal<{
    deviceCode: string;
    userCode: string;
    verificationUri: string;
    interval: number;
    expiresIn: number;
  } | null>(null);

  onMount(() => {
    const service = providerService();
    if (!service) return;

    AppRuntime.runFork(
      service.fetchProviders.pipe(
        Effect.tap(providers =>
          Effect.sync(() => {
            setProviders(providers);
          })
        ),
        Effect.catchAll(err =>
          Effect.sync(() => {
            setFetchError(err.message || "Failed to connect to runtime");
            setProviders([
              { id: "openai", name: "OpenAI", requires_api_key: true },
              { id: "anthropic", name: "Anthropic", requires_api_key: true },
              { id: "google", name: "Google AI", requires_api_key: true },
              { id: "ollama", name: "Ollama (Local)", requires_api_key: false },
            ]);
          })
        )
      )
    );

    AppRuntime.runFork(
      Effect.tryPromise({
        try: async () => {
          const res = await fetch(`${getRuntimeHttpUrl()}/api/v1/cloud/status`, { method: "GET" });
          if (!res.ok) {
            return { logged_in: false };
          }
          return (await res.json()) as { logged_in: boolean };
        },
        catch: () => ({ logged_in: false }),
      }).pipe(
        Effect.tap(result =>
          Effect.sync(() => {
            setCloudLoggedIn(!!result.logged_in);
            setStep(result.logged_in ? 1 : 0);
          })
        )
      )
    );
  });

  const handleProviderSelect = (providerId: string) => {
    const service = providerService();
    if (!service) return;

    setProvider(providerId);

    AppRuntime.runFork(
      service.fetchModels(providerId).pipe(
        Effect.tap(availableModels => {
          const defaultModel = availableModels.length > 0 ? availableModels[0].id : "";
          setModels(availableModels);
          setModelName(defaultModel);
          setSelectedModelIndex(0);
          setStep(2);
          setError("");
        }),
        Effect.catchAll(() => Effect.sync(() => setModels([])))
      )
    );
  };

  const startCloudLogin = () => {
    setLoading(true);
    setError("");
    setCloudLogin(null);

    AppRuntime.runFork(
      Effect.tryPromise({
        try: async () => {
          const res = await fetch(`${getRuntimeHttpUrl()}/api/v1/cloud/login/start`, {
            method: "POST",
          });
          if (!res.ok) {
            const text = await res.text();
            throw new Error(text || `HTTP ${res.status}`);
          }
          return (await res.json()) as {
            device_code: string;
            user_code: string;
            verification_uri: string;
            expires_in: number;
            interval: number;
          };
        },
        catch: e => e,
      }).pipe(
        Effect.tap(result =>
          Effect.sync(() => {
            if (result instanceof Error) {
              setError(result.message || "Failed to start cloud login");
              return;
            }
            setCloudLogin({
              deviceCode: result.device_code,
              userCode: result.user_code,
              verificationUri: result.verification_uri,
              interval: result.interval,
              expiresIn: result.expires_in,
            });
          })
        ),
        Effect.tap(() => Effect.sync(() => setLoading(false))),
        Effect.catchAll(() =>
          Effect.sync(() => {
            setLoading(false);
            setError("Failed to start cloud login");
          })
        )
      )
    );
  };

  const pollCloudLogin = () => {
    const login = cloudLogin();
    if (!login) {
      setError("Start login first");
      return;
    }

    setLoading(true);
    setError("");

    AppRuntime.runFork(
      Effect.tryPromise({
        try: async () => {
          const res = await fetch(`${getRuntimeHttpUrl()}/api/v1/cloud/login/poll`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              device_code: login.deviceCode,
              interval: login.interval,
              expires_in: login.expiresIn,
            }),
          });
          if (!res.ok) {
            const text = await res.text();
            throw new Error(text || `HTTP ${res.status}`);
          }
          return (await res.json()) as { ok: boolean };
        },
        catch: e => e,
      }).pipe(
        Effect.tap(result =>
          Effect.sync(() => {
            if (result instanceof Error) {
              setError(result.message || "Login failed");
              return;
            }
            if (result.ok) {
              setCloudLoggedIn(true);
              setStep(1);
              setCloudLogin(null);
            } else {
              setError("Login not complete");
            }
          })
        ),
        Effect.tap(() => Effect.sync(() => setLoading(false))),
        Effect.catchAll(() =>
          Effect.sync(() => {
            setLoading(false);
            setError("Login failed");
          })
        )
      )
    );
  };

  const handleSubmit = () => {
    const selectedProvider = providers().find(p => p.id === provider());

    if (selectedProvider?.requires_api_key && provider() !== "ollama" && !inputValue().trim()) {
      setError("API key is required");
      return;
    }

    const service = providerService();
    if (!service) {
      setError("Runtime not available");
      return;
    }

    setLoading(true);
    setError("");

    AppRuntime.runFork(
      Effect.gen(function* () {
        const providerId = provider();
        const key = inputValue().trim();

        if (selectedProvider?.requires_api_key && providerId !== "ollama") {
          yield* service.setProviderKey(providerId, key);
        }

        const cfg: AppConfig = {
          model_provider: providerId,
          model_name: modelName(),
        };

        if (providerId === "ollama" && key) {
          cfg.ollama_endpoint = key;
        }

        saveConfig(cfg);
      }).pipe(
        Effect.tap(() =>
          Effect.sync(() => {
            setStep(4);
            setTimeout(() => {
              props.onSetupComplete();
            }, 1500);
          })
        ),
        Effect.catchAll(() =>
          Effect.sync(() => {
            setError("Failed to save configuration");
          })
        ),
        Effect.tap(() => Effect.sync(() => setLoading(false)))
      )
    );
  };

  const selectedProvider = () => providers().find(p => p.id === provider());

  useKeyboard(evt => {
    if (loading()) {
      evt.preventDefault();
      return;
    }

    if (evt.ctrl && evt.name === "c") {
      return;
    }

    switch (step()) {
      case 0: {
        switch (evt.name) {
          case "enter":
          case "return":
            evt.preventDefault();
            if (!cloudLogin()) {
              startCloudLogin();
            } else {
              pollCloudLogin();
            }
            break;
          case "escape":
            evt.preventDefault();
            setError("");
            break;
          case "s":
            evt.preventDefault();
            setCloudLoggedIn(false);
            setStep(1);
            break;
        }
        break;
      }
      case 1: {
        const items = providers();
        if (items.length === 0) return;
        switch (evt.name) {
          case "up":
          case "arrowup":
            evt.preventDefault();
            setSelectedProviderIndex(i => (i - 1 + items.length) % items.length);
            break;
          case "down":
          case "arrowdown":
            evt.preventDefault();
            setSelectedProviderIndex(i => (i + 1) % items.length);
            break;
          case "enter":
          case "return": {
            evt.preventDefault();
            const chosen = items[selectedProviderIndex()];
            if (chosen) {
              handleProviderSelect(chosen.id);
              setInputValue("");
              setCursorPosition(0);
            }
            break;
          }
        }
        break;
      }
      case 2: {
        const items = models();
        if (items.length === 0) return;
        switch (evt.name) {
          case "up":
          case "arrowup":
            evt.preventDefault();
            setSelectedModelIndex(i => (i - 1 + items.length) % items.length);
            break;
          case "down":
          case "arrowdown":
            evt.preventDefault();
            setSelectedModelIndex(i => (i + 1) % items.length);
            break;
          case "enter":
          case "return": {
            evt.preventDefault();
            const chosen = items[selectedModelIndex()];
            if (!chosen) return;
            setModelName(chosen.id);
            const p = selectedProvider();
            if (p?.requires_api_key && p.id !== "ollama") {
              setStep(3);
              setInputValue("");
              setCursorPosition(0);
              return;
            }
            if (p?.id === "ollama") {
              setStep(3);
              setInputValue("");
              setCursorPosition(0);
              return;
            }
            handleSubmit();
            break;
          }
          case "escape":
            evt.preventDefault();
            setStep(1);
            break;
        }
        break;
      }
      case 3: {
        const pos = cursorPosition();
        const value = inputValue();
        switch (evt.name) {
          case "enter":
          case "return":
            evt.preventDefault();
            handleSubmit();
            break;
          case "escape":
            evt.preventDefault();
            setStep(2);
            setError("");
            break;
          case "backspace":
            evt.preventDefault();
            if (pos > 0) {
              const newValue = value.slice(0, pos - 1) + value.slice(pos);
              setInputValue(newValue);
              setCursorPosition(pos - 1);
            }
            break;
          case "delete":
            evt.preventDefault();
            if (pos < value.length) {
              const newValue = value.slice(0, pos) + value.slice(pos + 1);
              setInputValue(newValue);
            }
            break;
          case "left":
          case "arrowleft":
            evt.preventDefault();
            setCursorPosition(Math.max(0, pos - 1));
            break;
          case "right":
          case "arrowright":
            evt.preventDefault();
            setCursorPosition(Math.min(value.length, pos + 1));
            break;
          case "home":
            evt.preventDefault();
            setCursorPosition(0);
            break;
          case "end":
            evt.preventDefault();
            setCursorPosition(value.length);
            break;
          default:
            if (isPrintable(evt.name)) {
              evt.preventDefault();
              const newValue = value.slice(0, pos) + evt.name + value.slice(pos);
              setInputValue(newValue);
              setCursorPosition(pos + 1);
            }
            break;
        }
        break;
      }
    }
  });

  usePaste((evt: { text: string }) => {
    if (loading()) return;
    if (step() !== 3) return;
    const pos = cursorPosition();
    const value = inputValue();
    const text = evt.text;
    const newValue = value.slice(0, pos) + text + value.slice(pos);
    setInputValue(newValue);
    setCursorPosition(pos + text.length);
  });

  createEffect(() => {
    if (step() === 1) {
      setSelectedProviderIndex(0);
      setProvider("");
      setModels([]);
      setModelName("");
      setError("");
    }
  });

  createEffect(() => {
    if (step() === 2 && models().length > 0) {
      setSelectedModelIndex(0);
      setModelName(models()[0].id);
    }
  });

  const renderInput = (placeholder: string) => {
    const value = inputValue();
    const pos = cursorPosition();

    if (!value) {
      return (
        <box flexDirection="row">
          <text fg="gray">{placeholder}</text>
          <box flexGrow={1} />
          <text fg="cyan">▌</text>
        </box>
      );
    }

    return (
      <box flexDirection="row" flexWrap="wrap">
        <text fg="white">{value.slice(0, pos)}</text>
        <text fg="cyan" bg="cyan">
          {" "}
        </text>
        <text fg="white">{value.slice(pos)}</text>
      </box>
    );
  };

  return (
    <box flexDirection="column" flexGrow={1} padding={2}>
      <box marginBottom={2} flexDirection="column">
        <text fg="cyan">Welcome to Pryx!</text>
        <text fg="white">Setup Required</text>
        <text fg="gray">
          {step() === 0
            ? "To start, connect your Pryx Cloud account."
            : "To start chatting, configure an AI provider."}
        </text>
      </box>

      {fetchError() && (
        <box marginBottom={1}>
          <text fg="yellow">⚠ {fetchError()}</text>
        </box>
      )}

      <box flexDirection="column">
        <box flexDirection="row" marginBottom={1}>
          <text fg={step() === 0 ? "cyan" : cloudLoggedIn() ? "green" : "gray"}>
            Step 0: Pryx Cloud Login
          </text>
          {cloudLoggedIn() && <text fg="green"> ✓</text>}
        </box>

        {step() === 0 && (
          <box flexDirection="column" marginLeft={2}>
            <Show
              when={cloudLoggedIn()}
              fallback={
                <box flexDirection="column" gap={1}>
                  <Show
                    when={cloudLogin()}
                    fallback={<text fg="gray">Press Enter to start login</text>}
                  >
                    <box flexDirection="column" gap={1}>
                      <text fg="gray">Open this URL in your browser:</text>
                      <text fg="white">{cloudLogin()!.verificationUri}</text>
                      <text fg="gray">Enter this code:</text>
                      <text fg="cyan">{cloudLogin()!.userCode}</text>
                      <text fg="gray">Press Enter after authorizing</text>
                    </box>
                  </Show>
                  <text fg="gray">Press S to skip (offline)</text>
                </box>
              }
            >
              <text fg="green">✓ Logged in</text>
              <text fg="gray">Press Enter to continue</text>
            </Show>
            {error() && <text fg="red">{error()}</text>}
          </box>
        )}

        <box flexDirection="row" marginTop={1} marginBottom={1}>
          <text fg={step() >= 1 ? (step() === 1 ? "cyan" : "green") : "gray"}>
            Step 1: Choose Provider
          </text>
          {step() > 1 && <text fg="green"> ✓</text>}
        </box>

        {step() === 1 && (
          <box flexDirection="column" marginLeft={2}>
            <Show when={!loading()} fallback={<text fg="gray">Loading providers...</text>}>
              <box flexDirection="column" gap={1}>
                {providers().map((p, idx) => (
                  <box flexDirection="row">
                    <text fg={idx === selectedProviderIndex() ? "cyan" : "gray"}>
                      {idx === selectedProviderIndex() ? "❯ " : "  "}
                    </text>
                    <text fg="white">{p.name}</text>
                    <box flexGrow={1} />
                    <text fg="gray">{p.requires_api_key ? "API key" : "Local"}</text>
                  </box>
                ))}
              </box>
            </Show>
            <text fg="gray">↑↓ Select │ Enter Choose</text>
          </box>
        )}

        <box flexDirection="row" marginTop={1} marginBottom={1}>
          <text fg={step() >= 2 ? (step() === 2 ? "cyan" : "green") : "gray"}>
            Step 2: Choose Model
          </text>
          {step() > 2 && <text fg="green"> ✓</text>}
        </box>

        {step() === 2 && (
          <box flexDirection="column" marginLeft={2}>
            <text fg="gray">Provider: {selectedProvider()?.name}</text>
            <Show when={models().length > 0} fallback={<text fg="gray">No models available</text>}>
              <box flexDirection="column" gap={1} marginTop={1}>
                {models().map((m, idx) => (
                  <box flexDirection="row">
                    <text fg={idx === selectedModelIndex() ? "cyan" : "gray"}>
                      {idx === selectedModelIndex() ? "❯ " : "  "}
                    </text>
                    <text fg="white">{m.name}</text>
                    <box flexGrow={1} />
                    <text fg="gray">{m.id}</text>
                  </box>
                ))}
              </box>
            </Show>
            <text fg="gray">↑↓ Select │ Enter Choose │ Esc Back</text>
          </box>
        )}

        <box flexDirection="row" marginTop={1} marginBottom={1}>
          <text fg={step() >= 3 ? (step() === 3 ? "cyan" : "green") : "gray"}>
            Step 3: Configure Credentials
          </text>
          {step() > 3 && <text fg="green"> ✓</text>}
        </box>

        {step() === 3 && (
          <box flexDirection="column" marginLeft={2}>
            <text fg="gray">Provider: {selectedProvider()?.name}</text>
            <text fg="gray">Model: {modelName()}</text>
            <box marginTop={1} flexDirection="column">
              <text fg="gray">
                {selectedProvider()?.id === "ollama" ? "Ollama Endpoint (optional):" : "API Key:"}
              </text>
              <box borderStyle="single" borderColor={error() ? "red" : "gray"} padding={1}>
                {renderInput(
                  selectedProvider()?.id === "ollama"
                    ? "http://localhost:11434"
                    : "Paste your API key..."
                )}
              </box>
              {error() && <text fg="red">{error()}</text>}
            </box>
            <text fg="gray">Enter Save │ Esc Back</text>
          </box>
        )}

        {step() === 4 && (
          <box flexDirection="column" alignItems="center" marginTop={2}>
            <text fg="green">✓ Configuration Saved!</text>
            <text fg="gray">Starting Pryx...</text>
          </box>
        )}
      </box>

      <box flexGrow={1} />

      <box flexDirection="row">
        <text fg="gray">Need help? docs.pryx.dev</text>
      </box>
    </box>
  );
}
