import { createSignal, Show, Switch, Match } from "solid-js";
import { Effect } from "effect";
import { useEffectService } from "../lib/hooks";
import { WebSocketService } from "../services/ws";

type Step = 1 | 2 | 3 | "done";

interface WorkspaceConfig {
  name: string;
  path: string;
}

interface ProviderConfig {
  provider: string;
  apiKey: string;
}

export default function OnboardingWizard(props: { onComplete: () => void }) {
  const ws = useEffectService(WebSocketService);
  const [step, setStep] = createSignal<Step>(1);
  const [workspace, setWorkspace] = createSignal<WorkspaceConfig>({ name: "", path: "" });
  const [provider, setProvider] = createSignal<ProviderConfig>({ provider: "", apiKey: "" });
  const [input, setInput] = createSignal("");
  const [field, setField] = createSignal<"name" | "path" | "provider" | "apiKey" | "botToken">(
    "name"
  );

  const numericStep = () => (step() === "done" ? 3 : (step() as number));

  const handleSubmit = (value: string) => {
    const currentStep = step();
    const currentField = field();
    const service = ws();

    if (currentStep === 1) {
      if (currentField === "name") {
        setWorkspace(w => ({ ...w, name: value }));
        setField("path");
      } else {
        setWorkspace(w => ({ ...w, path: value }));
        setStep(2);
        setField("provider");
      }
    } else if (currentStep === 2) {
      if (currentField === "provider") {
        setProvider(p => ({ ...p, provider: value }));
        setField("apiKey");
      } else {
        setProvider(p => ({ ...p, apiKey: value }));
        setStep(3);
        setField("botToken");
      }
    } else if (currentStep === 3) {
      // Save configuration
      if (service) {
        Effect.runFork(
          service.send({
            event: "config.save",
            payload: {
              workspace: workspace(),
              provider: provider(),
              integration: { type: "telegram", botToken: value },
            },
          })
        );
      }
      setStep("done");
      setTimeout(() => props.onComplete(), 1500);
    }
    setInput("");
  };

  const getPlaceholder = () => {
    const f = field();
    switch (f) {
      case "name":
        return "Workspace name (e.g., my-project)";
      case "path":
        return "Workspace path (e.g., ~/code/my-project)";
      case "provider":
        return "Model provider (openai, anthropic, google)";
      case "apiKey":
        return "API key (sk-...)";
      case "botToken":
        return "Telegram bot token (from @BotFather)";
    }
  };

  return (
    <box flexDirection="column" flexGrow={1}>
      <box marginBottom={1}>
        <text bold fg="cyan">
          Onboarding Wizard
        </text>
        <text fg="gray"> - Step {step() === "done" ? "✓" : step()} of 3</text>
      </box>

      <box flexDirection="row" marginBottom={1}>
        <text
          fg={step() === 1 ? "cyan" : step() === "done" || numericStep() > 1 ? "green" : "gray"}
        >
          ● Workspace
        </text>
        <text fg="gray"> → </text>
        <text
          fg={step() === 2 ? "cyan" : step() === "done" || numericStep() > 2 ? "green" : "gray"}
        >
          ● Provider
        </text>
        <text fg="gray"> → </text>
        <text fg={step() === 3 ? "cyan" : step() === "done" ? "green" : "gray"}>● Integration</text>
      </box>

      <box flexDirection="column" flexGrow={1} borderStyle="single" padding={1}>
        <Switch>
          <Match when={step() === 1}>
            <text bold>Workspace Setup</text>
            <text fg="gray">
              {field() === "name"
                ? "Enter a name for your workspace"
                : "Enter the path to your workspace"}
            </text>
          </Match>
          <Match when={step() === 2}>
            <text bold>AI Provider Setup</text>
            <text fg="gray">
              {field() === "provider" ? "Choose your AI provider" : "Enter your API key"}
            </text>
          </Match>
          <Match when={step() === 3}>
            <text bold>Integration Setup</text>
            <text fg="gray">Enter your Telegram bot token</text>
          </Match>
          <Match when={step() === "done"}>
            <text bold fg="green">
              ✓ Setup Complete!
            </text>
            <text fg="gray">Redirecting to main interface...</text>
          </Match>
        </Switch>

        <Show when={step() !== "done"}>
          <box marginTop={1}>
            <input
              placeholder={getPlaceholder()}
              value={input()}
              onChange={setInput}
              onSubmit={handleSubmit}
            />
          </box>
        </Show>
      </box>
    </box>
  );
}
