import { useKeyboard } from "@opentui/solid";
import { createSignal } from "solid-js";

export interface KeybindInfo {
  ctrl?: boolean;
  meta?: boolean;
  shift?: boolean;
  name: string;
}

export function useKeybind() {
  const [leader, setLeader] = createSignal(false);

  const match = (
    keybind: KeybindInfo,
    evt: { ctrl: boolean; meta: boolean; shift: boolean; name: string }
  ): boolean => {
    return (
      (keybind.ctrl ?? false) === evt.ctrl &&
      (keybind.meta ?? false) === evt.meta &&
      (keybind.shift ?? false) === evt.shift &&
      keybind.name === evt.name
    );
  };

  const parse = (key: string): KeybindInfo => {
    const parts = key.toLowerCase().split("+");
    const info: KeybindInfo = { name: "" };

    for (const part of parts) {
      switch (part) {
        case "ctrl":
          info.ctrl = true;
          break;
        case "alt":
        case "meta":
          info.meta = true;
          break;
        case "shift":
          info.shift = true;
          break;
        case "esc":
          info.name = "escape";
          break;
        case "return":
        case "enter":
          info.name = "return";
          break;
        case "up":
          info.name = "up";
          break;
        case "down":
          info.name = "down";
          break;
        default:
          info.name = part;
      }
    }

    return info;
  };

  return {
    match,
    parse,
    leader,
    setLeader,
  };
}

export function useKeyboardHandler(handlers: Map<string, () => void>) {
  const { match, parse } = useKeybind();

  useKeyboard(evt => {
    for (const [key, handler] of handlers) {
      const keybind = parse(key);
      if (match(keybind, evt)) {
        handler();
        return;
      }
    }
  });
}
