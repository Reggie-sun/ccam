import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "tests",
  testMatch: "*.spec.ts",
  fullyParallel: true,
  reporter: [["list"]],
  use: { baseURL: "http://localhost:4322", channel: "chromium" },
  webServer: { command: "node scripts/serve.ts 4322", url: "http://localhost:4322/", reuseExistingServer: true },
  projects: [
    { name: "light", use: { ...devices["Desktop Chrome"], channel: "chromium", colorScheme: "light" } },
    { name: "dark", use: { ...devices["Desktop Chrome"], channel: "chromium", colorScheme: "dark" } },
    { name: "mobile", use: { ...devices["Pixel 7"], channel: "chromium" } },
  ],
});
