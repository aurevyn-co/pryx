import { createSignal, For, createEffect, onCleanup } from "solid-js";
import { Effect, Stream, Fiber } from "effect";
import { useKeyboard, usePaste } from "@opentui/solid";
import { useEffectService } from "../lib/hooks";
import { WebSocketService } from "../services/ws";
import Message, { MessageProps } from "./Message";
import { isPrintable } from "../lib/keybindings";

type RuntimeEvent = any;

interface ChatProps {
  disabled?: boolean;
  onConnectCommand?: () => void;
}

export default function Chat(props: ChatProps) {
  const ws = useEffectService(WebSocketService);
  const [messages, setMessages] = createSignal<MessageProps[]>([]);
  const [inputValue, setInputValue] = createSignal("");
  const [cursorPosition, setCursorPosition] = createSignal(0);
  const [sessionId] = createSignal(crypto.randomUUID());
  const [pendingApproval, setPendingApproval] = createSignal<{
    id: string;
    description: string;
  } | null>(null);
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
        Stream.runForEach(evt => Effect.sync(() => handleEvent(evt as RuntimeEvent)))
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
        setMessages(prev => [
          ...prev,
          {
            type: "assistant",
            content: streamingContent(),
            pending: false,
          },
        ]);
        setStreamingContent("");
        break;
      case "tool.start":
        setMessages(prev => [
          ...prev,
          {
            type: "tool",
            content: "Running...",
            toolName: evt.payload?.name,
            toolStatus: "running",
          },
        ]);
        break;
      case "tool.end":
        setMessages(prev => {
          const idx = prev.findLastIndex(
            m => m.toolName === evt.payload?.name && m.toolStatus === "running"
          );
          if (idx >= 0) {
            const updated = [...prev];
            updated[idx] = {
              ...updated[idx],
              content: evt.payload?.result ?? "Done",
              toolStatus: evt.payload?.error ? "error" : "done",
            };
            return updated;
          }
          return prev;
        });
        break;
      case "approval.request":
        setPendingApproval({
          id: evt.payload?.approval_id,
          description: evt.payload?.description ?? "Action requires approval",
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
        Effect.runFork(
          service.send({
            type: "approval.response",
            sessionId: sessionId(),
            approvalId: approval.id,
            approved: true,
          })
        );
        setMessages(prev => [...prev, { type: "system", content: "✅ Approved" }]);
        setPendingApproval(null);
        setInputValue("");
        setCursorPosition(0);
        return;
      } else if (value.toLowerCase() === "n" || value.toLowerCase() === "no") {
        Effect.runFork(
          service.send({
            type: "approval.response",
            sessionId: sessionId(),
            approvalId: approval.id,
            approved: false,
          })
        );
        setMessages(prev => [...prev, { type: "system", content: "❌ Denied" }]);
        setPendingApproval(null);
        setInputValue("");
        setCursorPosition(0);
        return;
      }
    }

    if (value.trim() === "/connect") {
      setInputValue("");
      setCursorPosition(0);
      if (props.onConnectCommand) {
        props.onConnectCommand();
      }
      return;
    }

    // Add to history
    setHistory(prev => [value, ...prev].slice(0, 100));
    setHistoryIndex(-1);

    setMessages(prev => [...prev, { type: "user", content: value }]);
    Effect.runFork(
      service.send({
        type: "chat.message",
        sessionId: sessionId(),
        content: value,
      })
    );
    setInputValue("");
    setCursorPosition(0);
    setIsStreaming(true);
  };

  // Handle keyboard input using OpenTUI's useKeyboard hook
  useKeyboard(evt => {
    if (props.disabled) return;

    const pos = cursorPosition();
    const value = inputValue();

    // Handle Ctrl+C for copy when text is selected (if we had selection)
    // For now, let it pass through to allow system copy
    if (evt.ctrl && evt.name === "c") {
      // Allow system copy - don't prevent default
      return;
    }

    // Handle Ctrl+V for paste - OpenTUI's usePaste handles this separately
    // But we need to handle it here to prevent default behavior
    if (evt.ctrl && evt.name === "v") {
      // Paste is handled by usePaste hook
      evt.preventDefault();
      return;
    }

    switch (evt.name) {
      case "return":
      case "enter":
        evt.preventDefault();
        handleSubmit();
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

      case "up":
      case "arrowup": {
        evt.preventDefault();
        const h = history();
        if (h.length > 0) {
          const newIndex = Math.min(historyIndex() + 1, h.length - 1);
          setHistoryIndex(newIndex);
          setInputValue(h[newIndex]);
          setCursorPosition(h[newIndex]?.length || 0);
        }
        break;
      }

      case "down":
      case "arrowdown": {
        evt.preventDefault();
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

      case "escape":
        evt.preventDefault();
        // Cancel/Clear
        setInputValue("");
        setCursorPosition(0);
        setHistoryIndex(-1);
        break;

      case "tab":
        // Ignore tab in chat input - let it propagate for view switching
        break;

      default:
        // Handle Ctrl+A (beginning of line)
        if (evt.ctrl && evt.name === "a") {
          evt.preventDefault();
          setCursorPosition(0);
          return;
        }

        // Handle Ctrl+E (end of line)
        if (evt.ctrl && evt.name === "e") {
          evt.preventDefault();
          setCursorPosition(value.length);
          return;
        }

        // Handle Ctrl+K (clear from cursor to end)
        if (evt.ctrl && evt.name === "k") {
          evt.preventDefault();
          setInputValue(value.slice(0, pos));
          return;
        }

        // Handle Ctrl+U (clear from start to cursor)
        if (evt.ctrl && evt.name === "u") {
          evt.preventDefault();
          setInputValue(value.slice(pos));
          setCursorPosition(0);
          return;
        }

        // Handle Ctrl+W (delete word before cursor)
        if (evt.ctrl && evt.name === "w") {
          evt.preventDefault();
          const beforeCursor = value.slice(0, pos);
          const match = beforeCursor.match(/^(.*\s)?(\S+)$/);
          if (match) {
            const newValue = (match[1] || "") + value.slice(pos);
            setInputValue(newValue);
            setCursorPosition(match[1]?.length || 0);
          }
          return;
        }

        // Handle printable characters
        if (isPrintable(evt.name)) {
          evt.preventDefault();
          const newValue = value.slice(0, pos) + evt.name + value.slice(pos);
          setInputValue(newValue);
          setCursorPosition(pos + 1);
        }
        break;
    }
  });

  usePaste((evt: { text: string }) => {
    if (props.disabled) return;

    const pos = cursorPosition();
    const value = inputValue();
    const text = evt.text;
    const newValue = value.slice(0, pos) + text + value.slice(pos);
    setInputValue(newValue);
    setCursorPosition(pos + text.length);
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
        <text fg="cyan" bg="cyan">
          {" "}
        </text>
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
        <For each={displayMessages()}>{msg => <Message {...msg} />}</For>

        {isStreaming() && streamingContent() && (
          <Message type="assistant" content={streamingContent()} pending={true} />
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
        <box flexGrow={1}>{renderInput()}</box>
      </box>
    </box>
  );
}
