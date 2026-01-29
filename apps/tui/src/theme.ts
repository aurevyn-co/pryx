export const palette = {
  text: "#E8E3D5",
  dim: "#7B7F87",
  accent: "#F6C453",
  accentSoft: "#F2A65A",
  accentBright: "#FFD700",
  bgPrimary: "#1a1a1a",
  bgSecondary: "#252525",
  bgSelected: "#2B2F36",
  bgHover: "#2A2D33",
  border: "#3C414B",
  borderAccent: "#F6C453",
  success: "#7DD3A5",
  error: "#F97066",
  warning: "#F2A65A",
  info: "#8CC8FF",
  userBg: "#2B2F36",
  userText: "#F3EEE0",
  systemText: "#9BA3B2",
};

function hexToRgb(hex: string): [number, number, number] {
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
  return result 
    ? [
        parseInt(result[1], 16),
        parseInt(result[2], 16),
        parseInt(result[3], 16),
      ]
    : [255, 255, 255];
}

export const theme = {
  fg: (color: string) => (text: string) => `\x1b[38;2;${hexToRgb(color).join(";")}m${text}\x1b[0m`,
  bg: (color: string) => (text: string) => `\x1b[48;2;${hexToRgb(color).join(";")}m${text}\x1b[0m`,
  bold: (text: string) => `\x1b[1m${text}\x1b[0m`,
  dim: (text: string) => `\x1b[2m${text}\x1b[0m`,
  italic: (text: string) => `\x1b[3m${text}\x1b[0m`,
  underline: (text: string) => `\x1b[4m${text}\x1b[0m`,
  text: (text: string) => theme.fg(palette.text)(text),
  textDim: (text: string) => theme.fg(palette.dim)(text),
  accent: (text: string) => theme.fg(palette.accent)(text),
  accentBold: (text: string) => theme.bold(theme.fg(palette.accent)(text)),
  border: (text: string) => theme.fg(palette.border)(text),
  success: (text: string) => theme.fg(palette.success)(text),
  error: (text: string) => theme.fg(palette.error)(text),
};

export const commandPaletteTheme = {
  borderColor: palette.border,
  backgroundColor: palette.bgPrimary,
  titleColor: palette.accent,
  itemText: (text: string, selected: boolean) => 
    selected ? theme.accentBold(text) : theme.text(text),
  itemNumber: (num: string, selected: boolean) =>
    selected ? theme.accent(num) : theme.textDim(num),
  itemShortcut: (shortcut: string) => theme.textDim(shortcut),
  itemBackground: (selected: boolean) => selected ? palette.bgSelected : undefined,
  footerText: (text: string) => theme.textDim(text),
};

export const colors = {
  text: palette.text,
  dim: palette.dim,
  accent: palette.accent,
  accentSoft: palette.accentSoft,
  bgPrimary: palette.bgPrimary,
  bgSelected: palette.bgSelected,
  border: palette.border,
};
