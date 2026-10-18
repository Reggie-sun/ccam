import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";
import vinext from "vinext";

export default defineConfig({
  plugins: [tailwindcss(), vinext({ appDir: import.meta.dirname })],
});
