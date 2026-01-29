import { createSignal, For, createEffect, onCleanup, onMount } from "solid-js";
import { Effect, Stream, Fiber } from "effect";
import { useEffectService } from "../lib/hooks";
import { WebSocketService } from "../services/ws";
import Message, { MessageProps } from "./Message";

// ANSI escape sequences for special keys
const KEYS = {
    ARROW_UP: "\u001b[A",
    ARROW_DOWN: "\u001b[B",
    ARROW_RIGHT: "\u001b[C",
    ARROW_LEFT: "\u001b[D",
    HOME: "\u001b[H",
    END: "\u001b[F",
    DELETE: "\u001b[3~",
    BACKSPACE: "\u007f",
    RETURN: "\r",
    NEWLINE: "\n",
    TAB: "\t",
    ESCAPE: "\u001b",
    CTRL_A: "\u0001",
    CTRL_E: "\u0005",
    CTRL_K: "\u000b",
    CTRL_U: "\u0015",
    CTRL_W: "\u0017",
    CTRL_C: "\u0003",
};

type RuntimeEvent = any;

interface ChatProps {
    disabled?: boolean;
}

export default function Chat(props: ChatProps) {
    const ws = useEffectService(WebSocketService);
    const [messages, setMessages] = createSignal<MessageProps[]>([]);
    const [inputValue, setInputValue] = createSignal("");
    const [cursorPosition, setCursorPosition] = createSignal(0);
    const [sessionId] = createSignal(crypto.randomUUID());
    const [pendingApproval, setPendingApproval] = createSignal<{ id: string, description: string } | null>(null);
    const [isStreaming, setIsStreaming] = createSignal(false);
    const [streamingContent, setStreamingContent] = createSignal("");
    const [history, setHistory] = createSignal<string[]>([]);
    const [historyIndex, setHistoryIndex] = createSignal(-1);

    createEffect(() => {
        const service = ws();
        if (!service) return;

        const connectFiber = Effect.runFork(service.connect);
        
        const messageFiber = Effect.runFork(
            service.messages.pipe(
                Stream.runForEach((evt) => Effect.sync(() => handleEvent(evt as RuntimeEvent)))
            )
        );

        onCleanup(() => {
            Effect.runFork(Fiber.interrupt(connectFiber));
            Effect.runFork(Fiber.interrupt(messageFiber));
            Effect.runFork(service.disconnect);
        });
    });

    const handleEvent = (evt: RuntimeEvent) => {
        switch (evt.event) {
            case "message.delta":
                setIsStreaming(true);
                setStreamingContent(prev => prev + (evt.payload?.content ?? ""));
                break;
            case "message.done":
                setIsStreaming(false);
                setMessages(prev => [...prev, {
                    type: "assistant",
                    content: streamingContent(),
                    pending: false
                }]);
                setStreamingContent("");
                break;
            case "tool.start":
                setMessages((prev) => [...prev, {
                    type: "tool",
                    content: "Running...",
                    toolName: evt.payload?.name,
                    toolStatus: "running"
                }]);
                break;
            case "tool.end":
                setMessages((prev) => {
                    const idx = prev.findLastIndex(m => m.toolName === evt.payload?.name && m.toolStatus === "running");
                    if (idx >= 0) {
                        const updated = [...prev];
                        updated[idx] = {
                            ...updated[idx],
                            content: evt.payload?.result ?? "Done",
                            toolStatus: evt.payload?.error ? "error" : "done"
                        };
                        return updated;
                    }
                    return prev;
                });
                break;
            case "approval.request":
                setPendingApproval({
                    id: evt.payload?.approval_id,
                    description: evt.payload?.description ?? "Action requires approval"
                });
                break;
        }
    };

    const handleSubmit = () => {
        const value = inputValue();
        if (!value.trim()) return;
        const service = ws();
        if (!service) return;

        if (pendingApproval()) {
            const approval = pendingApproval()!;
            if (value.toLowerCase() === "y" || value.toLowerCase() === "yes") {
                Effect.runFork(service.send({
                    type: "approval.response",
                    sessionId: sessionId(),
                    approvalId: approval.id,
                    approved: true
                }));
                setMessages((prev) => [...prev, { type: "system", content: "✅ Approved" }]);
                setPendingApproval(null);
                setInputValue("");
                setCursorPosition(0);
                return;
            } else if (value.toLowerCase() === "n" || value.toLowerCase() === "no") {
                Effect.runFork(service.send({
                    type: "approval.response",
                    sessionId: sessionId(),
                    approvalId: approval.id,
                    approved: false
                }));
                setMessages((prev) => [...prev, { type: "system", content: "❌ Denied" }]);
                setPendingApproval(null);
                setInputValue("");
                setCursorPosition(0);
                return;
            }
        }

        // Add to history
        setHistory(prev => [value, ...prev].slice(0, 100));
        setHistoryIndex(-1);

        setMessages((prev) => [...prev, { type: "user", content: value }]);
        Effect.runFork(service.send({
            type: "chat.message",
            sessionId: sessionId(),
            content: value
        }));
        setInputValue("");
        setCursorPosition(0);
        setIsStreaming(true);
    };

    const handleKey = (data: Buffer) => {
        if (props.disabled) return;
        
        const seq = data.toString();
        const pos = cursorPosition();
        const value = inputValue();

        switch (seq) {
            case KEYS.RETURN:
            case KEYS.NEWLINE:
                handleSubmit();
                break;

            case KEYS.BACKSPACE:
                if (pos > 0) {
                    const newValue = value.slice(0, pos - 1) + value.slice(pos);
                    setInputValue(newValue);
                    setCursorPosition(pos - 1);
                }
                break;

            case KEYS.DELETE:
                if (pos < value.length) {
                    const newValue = value.slice(0, pos) + value.slice(pos + 1);
                    setInputValue(newValue);
                }
                break;

            case KEYS.ARROW_LEFT:
                setCursorPosition(Math.max(0, pos - 1));
                break;

            case KEYS.ARROW_RIGHT:
                setCursorPosition(Math.min(value.length, pos + 1));
                break;

            case KEYS.HOME:
            case KEYS.CTRL_A:
                setCursorPosition(0);
                break;

            case KEYS.END:
            case KEYS.CTRL_E:
                setCursorPosition(value.length);
                break;

            case KEYS.CTRL_K:
                // Clear from cursor to end
                setInputValue(value.slice(0, pos));
                break;

            case KEYS.CTRL_U:
                // Clear from start to cursor
                setInputValue(value.slice(pos));
                setCursorPosition(0);
                break;

            case KEYS.CTRL_W: {
                const beforeCursor = value.slice(0, pos);
                const match = beforeCursor.match(/^(.*\s)?(\S+)$/);
                if (match) {
                    const newValue = (match[1] || "") + value.slice(pos);
                    setInputValue(newValue);
                    setCursorPosition(match[1]?.length || 0);
                }
                break;
            }

            case KEYS.ARROW_UP: {
                const h = history();
                if (h.length > 0) {
                    const newIndex = Math.min(historyIndex() + 1, h.length - 1);
                    setHistoryIndex(newIndex);
                    setInputValue(h[newIndex]);
                    setCursorPosition(h[newIndex]?.length || 0);
                }
                break;
            }

            case KEYS.ARROW_DOWN: {
                const idx = historyIndex();
                if (idx > 0) {
                    const newIndex = idx - 1;
                    setHistoryIndex(newIndex);
                    setInputValue(history()[newIndex]);
                    setCursorPosition(history()[newIndex]?.length || 0);
                } else if (idx === 0) {
                    setHistoryIndex(-1);
                    setInputValue("");
                    setCursorPosition(0);
                }
                break;
            }

            case KEYS.ESCAPE:
                // Cancel/Clear
                setInputValue("");
                setCursorPosition(0);
                setHistoryIndex(-1);
                break;

            case KEYS.TAB:
                // Ignore tab in chat input
                break;

            default:
                // Handle printable characters (including multi-byte for copy-paste)
                if (seq.length >= 1 && !seq.startsWith("\u001b")) {
                    // Check if it's a printable character or paste
                    const isPrintable = seq.split('').every(c => {
                        const code = c.charCodeAt(0);
                        return code >= 32 && code < 127;
                    });
                    
                    if (isPrintable) {
                        const newValue = value.slice(0, pos) + seq + value.slice(pos);
                        setInputValue(newValue);
                        setCursorPosition(pos + seq.length);
                    }
                }
                break;
        }
    };

    onMount(() => {
        if (typeof process !== "undefined" && process.stdin.isTTY) {
            process.stdin.on("data", handleKey);
        }
    });

    onCleanup(() => {
        if (typeof process !== "undefined" && process.stdin) {
            process.stdin.off("data", handleKey);
        }
    });

    const displayMessages = () => [...messages()];

    // Render input with cursor
    const renderInput = () => {
        const value = inputValue();
        const pos = cursorPosition();
        
        if (!value) {
            return (
                <box flexDirection="row">
                    <text fg="gray">Type a message... (Enter to send, ↑↓ for history)</text>
                    <box flexGrow={1} />
                    <text fg="cyan">▌</text>
                </box>
            );
        }

        return (
            <box flexDirection="row" flexWrap="wrap">
                <text fg="white">{value.slice(0, pos)}</text>
                <text fg="cyan" bg="cyan"> </text>
                <text fg="white">{value.slice(pos)}</text>
            </box>
        );
    };

    return (
        <box flexDirection="column" flexGrow={1}>
            <box 
                flexDirection="column" 
                flexGrow={1} 
                borderStyle="single" 
                borderColor="cyan" 
                padding={1}
                gap={1}
            >
                <For each={displayMessages()}>
                    {(msg) => <Message {...msg} />}
                </For>
                
                {isStreaming() && streamingContent() && (
                    <Message 
                        type="assistant" 
                        content={streamingContent()} 
                        pending={true}
                    />
                )}
            </box>

            {pendingApproval() && (
                <box 
                    borderStyle="double" 
                    borderColor="yellow" 
                    padding={1} 
                    marginTop={1}
                    flexDirection="row"
                >
                    <text fg="yellow">⚠️ {pendingApproval()!.description}</text>
                    <box flexGrow={1} />
                    <text fg="gray">(y/n)</text>
                </box>
            )}

            <box 
                borderStyle="single" 
                borderColor={inputValue() ? "cyan" : "gray"} 
                marginTop={1} 
                padding={1}
                flexDirection="row"
                gap={1}
            >
                <text fg="cyan">❯</text>
                <box flexGrow={1}>
                    {renderInput()}
                </box>
            </box>
        </box>
    );
}
