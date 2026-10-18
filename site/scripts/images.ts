// Converts raster images in public/img to AVIF and WebP at several widths.
// The site currently uses SVG only, so this usually has nothing to do; sharp
// is loaded only when there is work.
import { existsSync, readdirSync } from "node:fs";
import { join, parse } from "node:path";

const dir = "public/img";
const sources = existsSync(dir) ? readdirSync(dir).filter((f) => /\.(png|jpe?g)$/i.test(f)) : [];
if (!sources.length) {
  console.log("images: no raster images");
} else {
  const { default: sharp } = await import("sharp" as string);
  for (const f of sources) {
    for (const width of [480, 960, 1440]) {
      const base = join("dist/client/img", `${parse(f).name}-${width}`);
      await sharp(join(dir, f)).resize({ width, withoutEnlargement: true }).avif({ quality: 55 }).toFile(`${base}.avif`);
      await sharp(join(dir, f)).resize({ width, withoutEnlargement: true }).webp({ quality: 75 }).toFile(`${base}.webp`);
    }
  }
  console.log(`images: ${sources.length} converted`);
}
