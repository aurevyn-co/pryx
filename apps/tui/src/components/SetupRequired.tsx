import { createSignal } from "solid-js";
import { saveConfig } from "../services/config";

interface SetupRequiredProps {
    onSetupComplete: () => void;
}

export default function SetupRequired(props: SetupRequiredProps) {
    const [step, setStep] = createSignal(1);
    const [provider, setProvider] = createSignal("");
    const [apiKey, setApiKey] = createSignal("");
    const [modelName, setModelName] = createSignal("");
    const [error, setError] = createSignal("");

    const providers = [
        { id: "glm", name: "GLM (Zhipu AI)", models: ["glm-4.5", "glm-4.5-air", "glm-4.6", "glm-4.7"] },
        { id: "openai", name: "OpenAI", models: ["gpt-4", "gpt-4-turbo", "gpt-3.5-turbo"] },
        { id: "anthropic", name: "Anthropic", models: ["claude-3-opus", "claude-3-sonnet", "claude-3-haiku"] },
        { id: "ollama", name: "Ollama (Local)", models: ["llama3", "llama2", "mistral", "codellama"] },
    ];

    const handleProviderSelect = (providerId: string) => {
        setProvider(providerId);
        const defaultModel = providers.find(p => p.id === providerId)?.models[0] || "";
        setModelName(defaultModel);
        setStep(2);
        setError("");
    };

    const handleSubmit = () => {
        if (!apiKey().trim()) {
            setError("API key is required");
            return;
        }

        const config: any = {
            model_provider: provider(),
            model_name: modelName(),
        };

        if (provider() === "openai") {
            config.openai_key = apiKey();
        } else if (provider() === "anthropic") {
            config.anthropic_key = apiKey();
        } else if (provider() === "glm") {
            config.glm_key = apiKey();
        } else if (provider() === "ollama") {
            config.ollama_endpoint = apiKey();
        }

        try {
            saveConfig(config);
            setStep(3);
            setTimeout(() => {
                props.onSetupComplete();
            }, 1500);
        } catch (e) {
            setError("Failed to save configuration");
        }
    };

    return (
        <box flexDirection="column" flexGrow={1} padding={2}>
            <box marginBottom={2} flexDirection="column">
                <text fg="cyan">Welcome to Pryx!</text>
                <text fg="white">Setup Required</text>
                <text fg="gray">To start chatting, you need to configure an AI provider.</text>
            </box>

            <box flexDirection="column">
                <box flexDirection="row" marginBottom={1}>
                    <text fg={step() >= 1 ? "cyan" : "gray"}>Step 1: Choose Provider</text>
                    {step() > 1 && <text fg="green"> ✓</text>}
                </box>

                {step() === 1 && (
                    <box flexDirection="column" marginLeft={2}>
                        {providers.map(p => (
                            <box
                                borderStyle="single"
                                borderColor="gray"
                                padding={1}
                                flexDirection="column"
                            >
                                <text fg="white">{p.name}</text>
                                <text fg="gray">Models: {p.models.join(", ")}</text>
                            </box>
                        ))}
                    </box>
                )}

                <box flexDirection="row" marginTop={1} marginBottom={1}>
                    <text fg={step() >= 2 ? "cyan" : "gray"}>Step 2: API Configuration</text>
                    {step() > 2 && <text fg="green"> ✓</text>}
                </box>

                {step() === 2 && (
                    <box flexDirection="column" marginLeft={2}>
                        <box>
                            <text fg="gray">Selected: {providers.find(p => p.id === provider())?.name}</text>
                        </box>

                        <box marginTop={1}>
                            <text fg="gray">Model: {modelName()}</text>
                        </box>

                        <box marginTop={1}>
                            <text fg="gray">{provider() === "ollama" ? "Endpoint:" : "API Key:"}</text>
                            <box
                                borderStyle="single"
                                borderColor={error() ? "red" : "gray"}
                                padding={1}
                                flexDirection="row"
                            >
                                <text fg="white">{apiKey() || "Enter value..."}</text>
                                <box flexGrow={1} />
                                <text fg="cyan">▌</text>
                            </box>
                            {error() && <text fg="red">{error()}</text>}
                        </box>

                        <box marginTop={1}>
                            <box borderStyle="single" borderColor="cyan" padding={1}>
                                <text fg="cyan">Save Configuration</text>
                            </box>
                        </box>
                    </box>
                )}

                {step() === 3 && (
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
