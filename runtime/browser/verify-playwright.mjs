import { chromium } from "playwright-core";

const port = Number.parseInt(process.env.AIH_CDP_PORT || "9222", 10);
const endpoint = `http://127.0.0.1:${port}`;

const result = {
  ok: false,
  port,
  endpoint,
  connected: false,
  title: null,
  browserVersion: null,
  userAgent: null,
};

let browser;

try {
  browser = await chromium.connectOverCDP(endpoint);
  result.connected = true;
  result.browserVersion = browser.version();

  const context = browser.contexts()[0] ?? await browser.newContext();
  const page = context.pages()[0] ?? await context.newPage();
  await page.goto("data:text/html,<title>aih-playwright-ok</title><h1>aih</h1>", {
    waitUntil: "load",
  });
  result.title = await page.title();
  result.userAgent = await page.evaluate(() => navigator.userAgent);
  result.ok = result.title === "aih-playwright-ok";

  console.log(JSON.stringify(result, null, 2));
  process.exit(result.ok ? 0 : 2);
} catch (error) {
  result.error = error instanceof Error ? error.message : String(error);
  console.log(JSON.stringify(result, null, 2));
  process.exit(2);
} finally {
  await browser?.close().catch(() => {});
}
