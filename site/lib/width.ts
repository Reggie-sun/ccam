// Search results cut titles and descriptions by width; CJK and other
// full-width characters take about two Latin character widths.
export const width = (s: string) => [...s].reduce((n, c) => n + (/[ᄀ-ᅟ⺀-꓏가-힣豈-﫿︰-﹏＀-｠]/.test(c) ? 2 : 1), 0);
export const TITLE_MAX = 60;
