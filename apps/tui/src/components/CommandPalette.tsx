import { createSignal, For, createEffect } from "solid-js";
import { useKeyboard } from "@opentui/solid";
import { palette } from "../theme";

interface Command {
    id: string;
    label: string;
    shortcut: string;
    action: () => void;
}

interface CommandPaletteProps {
    commands: Command[];
    onClose: () => void;
}

export default function CommandPalette(props: CommandPaletteProps) {
    const [selectedIndex, setSelectedIndex] = createSignal(0);
    const [inputMode, setInputMode] = createSignal<"keyboard" | "mouse">("keyboard");

    const executeCommand = (index: number) => {
        const cmd = props.commands[index];
        if (cmd) {
            cmd.action();
        }
    };

    useKeyboard((evt) => {
        setInputMode("keyboard");
        
        switch (evt.name) {
            case "up":
                evt.preventDefault();
                setSelectedIndex(i => Math.max(0, i - 1));
                break;
            case "down":
                evt.preventDefault();
                setSelectedIndex(i => Math.min(props.commands.length - 1, i + 1));
                break;
            case "return":
                evt.preventDefault();
                executeCommand(selectedIndex());
                break;
            case "escape":
                evt.preventDefault();
                props.onClose();
                break;
            default:
                if (evt.name >= "1" && evt.name <= "9") {
                    const idx = parseInt(evt.name) - 1;
                    if (idx < props.commands.length) {
                        setSelectedIndex(idx);
                    }
                }
        }
    });

    const handleMouseOver = (index: number) => {
        if (inputMode() !== "mouse") return;
        setSelectedIndex(index);
    };

    const handleMouseMove = () => {
        setInputMode("mouse");
    };

    const handleMouseUp = (index: number) => {
        setSelectedIndex(index);
        executeCommand(index);
    };

    createEffect(() => {
        const idx = selectedIndex();
        if (idx >= props.commands.length) {
            setSelectedIndex(0);
        }
    });

    return (
        <box
            position="absolute"
            top={3}
            left="20%"
            width="60%"
            borderStyle="single"
            borderColor={palette.border}
            backgroundColor={palette.bgPrimary}
            flexDirection="column"
            padding={1}
        >
            <box marginBottom={1}>
                <text fg={palette.accent}>Commands</text>
            </box>
            
            <box flexDirection="column">
                <For each={props.commands}>
                    {(cmd, index) => {
                        const isSelected = index() === selectedIndex();
                        
                        return (
                            <box 
                                flexDirection="row" 
                                padding={1}
                                backgroundColor={isSelected ? palette.bgSelected : undefined}
                                onMouseMove={handleMouseMove}
                                onMouseOver={() => handleMouseOver(index())}
                                onMouseUp={() => handleMouseUp(index())}
                                onMouseDown={() => setSelectedIndex(index())}
                            >
                                <box width={3}>
                                    <text fg={isSelected ? palette.accent : palette.dim}>
                                        {index() + 1}.
                                    </text>
                                </box>
                                <box flexGrow={1}>
                                    <text fg={isSelected ? palette.accent : palette.text}>
                                        {cmd.label}
                                    </text>
                                </box>
                                <box width={10} alignItems="flex-end">
                                    <text fg={palette.dim}>{cmd.shortcut}</text>
                                </box>
                            </box>
                        );
                    }}
                </For>
            </box>
            
            <box marginTop={1} flexDirection="column" gap={0}>
                <text fg={palette.dim}>↑↓ Navigate | Enter Select | Esc Close | 1-9 Quick Select</text>
            </box>
        </box>
    );
}
