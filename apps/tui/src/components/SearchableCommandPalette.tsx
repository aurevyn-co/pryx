import { createSignal, For, onMount, onCleanup, createMemo } from "solid-js";
import { palette } from "../theme";

export interface Command {
    id: string;
    name: string;
    description: string;
    shortcut?: string;
    category: string;
    action: () => void;
    keywords?: string[];
}

interface SearchableCommandPaletteProps {
    commands: Command[];
    onClose: () => void;
    placeholder?: string;
}

export default function SearchableCommandPalette(props: SearchableCommandPaletteProps) {
    const [searchQuery, setSearchQuery] = createSignal("");
    const [selectedIndex, setSelectedIndex] = createSignal(0);
    const [filteredCommands, setFilteredCommands] = createSignal<Command[]>([]);

    const allCommands = () => props.commands;

    const filterCommands = (query: string) => {
        if (!query.trim()) {
            return allCommands();
        }
        
        const lowerQuery = query.toLowerCase();
        return allCommands().filter(cmd => {
            const nameMatch = cmd.name.toLowerCase().includes(lowerQuery);
            const descMatch = cmd.description.toLowerCase().includes(lowerQuery);
            const categoryMatch = cmd.category.toLowerCase().includes(lowerQuery);
            const keywordMatch = cmd.keywords?.some(k => k.toLowerCase().includes(lowerQuery));
            return nameMatch || descMatch || categoryMatch || keywordMatch;
        });
    };

    createMemo(() => {
        const filtered = filterCommands(searchQuery());
        setFilteredCommands(filtered);
        if (selectedIndex() >= filtered.length) {
            setSelectedIndex(0);
        }
    });

    const groupedByCategory = createMemo(() => {
        const groups: Record<string, Command[]> = {};
        filteredCommands().forEach(cmd => {
            if (!groups[cmd.category]) {
                groups[cmd.category] = [];
            }
            groups[cmd.category].push(cmd);
        });
        return groups;
    });

    const getCategoryColor = (category: string) => {
        const colors: Record<string, string> = {
            "Navigation": palette.accent,
            "Chat": palette.success,
            "Skills": palette.accentSoft,
            "MCP": palette.info,
            "Settings": palette.dim,
            "System": palette.error,
            "Help": palette.dim
        };
        return colors[category] || palette.text;
    };

    const executeCommand = (cmd: Command) => {
        cmd.action();
        props.onClose();
    };

    onMount(() => {
        const handleKey = (data: Buffer) => {
            const seq = data.toString();
            const commands = filteredCommands();

            switch (seq) {
                case '\u001b':
                    props.onClose();
                    return;

                case '\r':
                case '\n':
                    if (commands.length > 0) {
                        executeCommand(commands[selectedIndex()]);
                    }
                    return;

                case '\u001b[A':
                    setSelectedIndex(i => Math.max(0, i - 1));
                    return;

                case '\u001b[B':
                    setSelectedIndex(i => Math.min(commands.length - 1, i + 1));
                    return;

                case '\u007f':
                case '\b':
                    setSearchQuery(q => q.slice(0, -1));
                    return;

                case '\t':
                    return;
            }

            if (seq.length === 1 && seq.charCodeAt(0) >= 32) {
                setSearchQuery(q => q + seq);
            }
        };

        if (typeof process !== "undefined" && process.stdin.isTTY) {
            process.stdin.on("data", handleKey);
        }

        onCleanup(() => {
            if (typeof process !== "undefined" && process.stdin) {
                process.stdin.off("data", handleKey);
            }
        });
    });

    return (
        <box
            position="absolute"
            top={3}
            left="10%"
            width="80%"
            height="80%"
            borderStyle="double"
            borderColor={palette.border}
            backgroundColor={palette.bgPrimary}
            flexDirection="column"
            padding={1}
        >
            <box flexDirection="row" marginBottom={1} gap={1}>
                <text fg={palette.accent}>/</text>
                <box flexGrow={1} borderStyle="single" borderColor={searchQuery() ? palette.accent : palette.dim} padding={0.5}>
                    {searchQuery() ? (
                        <text fg={palette.text}>{searchQuery()}</text>
                    ) : (
                        <text fg={palette.dim}>Type to search commands...</text>
                    )}
                </box>
                <box flexGrow={1} />
                <text fg={palette.dim}>{filteredCommands().length} / {allCommands().length}</text>
            </box>

            <box flexDirection="column" flexGrow={1} overflow="scroll">
                <For each={Object.entries(groupedByCategory())}>
                    {([category, commands]) => (
                        <box flexDirection="column" marginBottom={1}>
                            <box marginBottom={0.5}>
                                <text fg={getCategoryColor(category)}>{category}</text>
                            </box>
                            
                            <For each={commands}>
                                {(cmd, idx) => {
                                    const globalIdx = filteredCommands().indexOf(cmd);
                                    const isSelected = globalIdx === selectedIndex();
                                    return (
                                        <box 
                                            flexDirection="row" 
                                            padding={0.5}
                                            backgroundColor={isSelected ? palette.bgSelected : undefined}
                                        >
                                            <box width={25}>
                                                <text fg={isSelected ? palette.accent : palette.accentSoft}>
                                                    {cmd.name}
                                                </text>
                                            </box>
                                            <box flexGrow={1}>
                                                <text fg={isSelected ? palette.text : palette.dim}>
                                                    {cmd.description}
                                                </text>
                                            </box>
                                            {cmd.shortcut && (
                                                <box width={10}>
                                                    <text fg={isSelected ? palette.accentSoft : palette.dim}>
                                                        {cmd.shortcut}
                                                    </text>
                                                </box>
                                            )}
                                        </box>
                                    );
                                }}
                            </For>
                        </box>
                    )}
                </For>

                {filteredCommands().length === 0 && (
                    <box flexDirection="column" alignItems="center" marginTop={2}>
                        <text fg={palette.dim}>No commands found matching "{searchQuery()}"</text>
                    </box>
                )}
            </box>

            <box flexDirection="row" marginTop={1} gap={2}>
                <text fg={palette.dim}>↑↓ Navigate | Enter Select | Esc Close | Type to filter</text>
            </box>
        </box>
    );
}
