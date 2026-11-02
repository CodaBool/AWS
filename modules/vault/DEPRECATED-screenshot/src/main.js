import chromium from "@sparticuz/chromium"; // AWS Lambda-compatible Chromium
import puppeteer from "puppeteer-core"; // Use core version
import { setTimeout } from "node:timers/promises";

export const handler = async (event) => {
  try {
    const queryParams = event.queryStringParameters || {};
    const { map, uuid, z, lng, lat } = queryParams;

    if (!z || !map || !uuid) {
      return {
        statusCode: 400,
        body: "Missing required parameters",
        headers: {
          "Access-Control-Allow-Origin": "*",  // Allow all origins
          "Access-Control-Allow-Methods": "GET, OPTIONS",
        },
      };
    }

    // Construct URL for Stargazer map
    const stargazer = new URL(`https://starlazer.vercel.app/${map}/${uuid}`);
    stargazer.searchParams.set("z", z);
    if (lng) stargazer.searchParams.set("lng", lng);
    if (lat) stargazer.searchParams.set("lat", lat);
    stargazer.searchParams.set("mini", 1);

    // Launch headless Chromium
    const browser = await puppeteer.launch({
      args: chromium.args,
      executablePath: await chromium.executablePath(),
      headless: chromium.headless,
      ignoreHTTPSErrors: true,
    })

    const page = await browser.newPage()
    await page.setViewport({
      width: 3840,  // 4K Width
      height: 2160, // 4K Height (16:9 aspect ratio)
      deviceScaleFactor: 2 // Increase rendering scale
    })
    console.log("Navigating to:", stargazer.toString());
    await page.goto(stargazer.toString(), { waitUntil: "networkidle2" });

    await page.waitForSelector("canvas.maplibregl-canvas", { visible: true })
    await setTimeout(1000)

    const body = await page.screenshot({ type: "webp", optimizeForSpeed: true, encoding: "base64", quality: 60 });
    await browser.close()
    return {
      statusCode: 200,
      headers: {
        "Content-Type": "image/webp",
        'Content-Length': body.length,
        "Access-Control-Allow-Origin": "*",  // Allow all origins
        "Access-Control-Allow-Methods": "GET, OPTIONS",
      },
      body,
      isBase64Encoded: true,
    };
  } catch (error) {
    console.error(error);
    return {
      statusCode: 500,
      body: "Internal Server Error",
      headers: {
        "Access-Control-Allow-Origin": "*",  // Allow all origins
        "Access-Control-Allow-Methods": "GET, OPTIONS",
      },
    };
  }
};
