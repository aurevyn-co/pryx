import { createSignal, createEffect, onCleanup } from "solid-js";

export interface MouseState {
  x: number;
  y: number;
  button: number;
  isDown: boolean;
}

export interface MouseEvent {
  x: number;
  y: number;
  button: number;
}

export function useMouse() {
  const [mouseState, setMouseState] = createSignal<MouseState>({
    x: 0,
    y: 0,
    button: 0,
    isDown: false,
  });

  const [hoveredIndex, setHoveredIndex] = createSignal<number | null>(null);

  const handleMouseMove = (x: number, y: number) => {
    setMouseState(prev => ({ ...prev, x, y }));
  };

  const handleMouseDown = (x: number, y: number, button: number) => {
    setMouseState({ x, y, button, isDown: true });
  };

  const handleMouseUp = (x: number, y: number, button: number) => {
    setMouseState(prev => ({ ...prev, x, y, button, isDown: false }));
  };

  const isInside = (boxX: number, boxY: number, boxWidth: number, boxHeight: number) => {
    const state = mouseState();
    return (
      state.x >= boxX &&
      state.x < boxX + boxWidth &&
      state.y >= boxY &&
      state.y < boxY + boxHeight
    );
  };

  const getItemIndex = (startY: number, itemHeight: number, count: number) => {
    const state = mouseState();
    const index = Math.floor((state.y - startY) / itemHeight);
    if (index >= 0 && index < count) {
      return index;
    }
    return null;
  };

  return {
    mouseState,
    hoveredIndex,
    setHoveredIndex,
    handleMouseMove,
    handleMouseDown,
    handleMouseUp,
    isInside,
    getItemIndex,
  };
}

export function useClipboard() {
  const copy = async (text: string): Promise<boolean> => {
    try {
      if (typeof navigator !== "undefined" && navigator.clipboard) {
        await navigator.clipboard.writeText(text);
        return true;
      }
      
      if (typeof process !== "undefined") {
        const proc = Bun.spawn(["pbcopy"], {
          stdin: "pipe",
        });
        proc.stdin.write(text);
        proc.stdin.end();
        await proc.exited;
        return true;
      }
      
      return false;
    } catch {
      return false;
    }
  };

  return { copy };
}
