const puppeteer = require('puppeteer-core');
const path = require('path');
const fs = require('fs');
const { execSync } = require('child_process');

const CHROME_PATH = 'C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe';
const OUTPUT_DIR = path.join(__dirname, '..', 'docs', 'images');
const FRAMES_DIR = path.join(__dirname, '..', 'tmp_frames');

async function capture() {
  if (!fs.existsSync(OUTPUT_DIR)) {
    fs.mkdirSync(OUTPUT_DIR, { recursive: true });
  }
  if (!fs.existsSync(FRAMES_DIR)) {
    fs.mkdirSync(FRAMES_DIR, { recursive: true });
  }

  console.log('Launching Chrome from:', CHROME_PATH);
  const browser = await puppeteer.launch({
    executablePath: CHROME_PATH,
    headless: true,
    args: ['--no-sandbox', '--disable-setuid-sandbox', '--disable-gpu'],
  });

  const page = await browser.newPage();

  // 1. Desktop Dashboard (1440x900)
  await page.setViewport({ width: 1440, height: 900, deviceScaleFactor: 2 });
  console.log('Navigating to http://localhost:6363...');
  await page.goto('http://localhost:6363', { waitUntil: 'domcontentloaded', timeout: 10000 });
  await page.waitForSelector('h1', { timeout: 10000 });
  await new Promise((r) => setTimeout(r, 2000)); // wait for initial scan and animations to settle

  console.log('Capturing desktop dashboard screenshot...');
  await page.screenshot({
    path: path.join(OUTPUT_DIR, 'governmesh_dashboard.png'),
    fullPage: false,
  });

  // 2. Mobile Responsive View (390x844)
  console.log('Capturing mobile screenshot (390x844)...');
  await page.setViewport({ width: 390, height: 844, deviceScaleFactor: 2 });
  await new Promise((r) => setTimeout(r, 800));
  await page.screenshot({
    path: path.join(OUTPUT_DIR, 'governmesh_mobile.png'),
    fullPage: false,
  });

  // Reset to Desktop Viewport
  await page.setViewport({ width: 1440, height: 900, deviceScaleFactor: 2 });
  await new Promise((r) => setTimeout(r, 600));

  // 3. AIBOM Inventory Tab
  console.log('Switching to AIBOM tab...');
  const tabs = await page.$$('button');
  for (const tab of tabs) {
    const text = await page.evaluate((el) => el.textContent, tab);
    if (text && text.includes('AIBOM Inventory')) {
      await tab.click();
      break;
    }
  }
  await new Promise((r) => setTimeout(r, 1000));

  // Expand first row to show details
  const firstRow = await page.$('tbody tr');
  if (firstRow) {
    await firstRow.click();
    await new Promise((r) => setTimeout(r, 600));
  }

  console.log('Capturing AIBOM explorer screenshot...');
  await page.screenshot({
    path: path.join(OUTPUT_DIR, 'governmesh_aibom.png'),
    fullPage: false,
  });

  // 4. Code Remediations & Patches Tab
  console.log('Switching to Remediation tab...');
  const tabs2 = await page.$$('button');
  for (const tab of tabs2) {
    const text = await page.evaluate((el) => el.textContent, tab);
    if (text && text.includes('Code Remediations')) {
      await tab.click();
      break;
    }
  }
  await new Promise((r) => setTimeout(r, 1200));

  console.log('Capturing Remediation diff screenshot...');
  await page.screenshot({
    path: path.join(OUTPUT_DIR, 'governmesh_remediation.png'),
    fullPage: false,
  });

  // 5. Annex IV Technical Dossier Tab
  console.log('Switching to Annex IV tab...');
  const tabs3 = await page.$$('button');
  for (const tab of tabs3) {
    const text = await page.evaluate((el) => el.textContent, tab);
    if (text && text.includes('Annex IV')) {
      await tab.click();
      break;
    }
  }
  await new Promise((r) => setTimeout(r, 1000));

  console.log('Capturing Annex IV dossier screenshot...');
  await page.screenshot({
    path: path.join(OUTPUT_DIR, 'governmesh_annex_iv.png'),
    fullPage: false,
  });

  // ==================== GENERATING FRAMES FOR DEMO GIF ====================
  console.log('Capturing frames for animated demo GIF...');

  // Frame 1: Overview top
  const tabsOverview = await page.$$('button');
  for (const tab of tabsOverview) {
    const text = await page.evaluate((el) => el.textContent, tab);
    if (text && text.includes('Executive Overview')) {
      await tab.click();
      break;
    }
  }
  await new Promise((r) => setTimeout(r, 800));
  await page.screenshot({ path: path.join(FRAMES_DIR, 'frame_01.png') });

  // Frame 2: Overview scroll down
  await page.evaluate(() => window.scrollBy(0, 350));
  await new Promise((r) => setTimeout(r, 600));
  await page.screenshot({ path: path.join(FRAMES_DIR, 'frame_02.png') });

  // Frame 3: AIBOM Tab
  for (const tab of tabsOverview) {
    const text = await page.evaluate((el) => el.textContent, tab);
    if (text && text.includes('AIBOM Inventory')) {
      await tab.click();
      break;
    }
  }
  await new Promise((r) => setTimeout(r, 800));
  await page.screenshot({ path: path.join(FRAMES_DIR, 'frame_03.png') });

  // Frame 4: PII Tab
  const tabsPII = await page.$$('button');
  for (const tab of tabsPII) {
    const text = await page.evaluate((el) => el.textContent, tab);
    if (text && text.includes('PII & Egress')) {
      await tab.click();
      break;
    }
  }
  await new Promise((r) => setTimeout(r, 800));
  await page.screenshot({ path: path.join(FRAMES_DIR, 'frame_04.png') });

  // Frame 5: Code Remediations Tab
  const tabsRemed = await page.$$('button');
  for (const tab of tabsRemed) {
    const text = await page.evaluate((el) => el.textContent, tab);
    if (text && text.includes('Code Remediations')) {
      await tab.click();
      break;
    }
  }
  await new Promise((r) => setTimeout(r, 800));
  await page.screenshot({ path: path.join(FRAMES_DIR, 'frame_05.png') });

  // Frame 6: Annex IV Tab
  const tabsAnnex = await page.$$('button');
  for (const tab of tabsAnnex) {
    const text = await page.evaluate((el) => el.textContent, tab);
    if (text && text.includes('Annex IV')) {
      await tab.click();
      break;
    }
  }
  await new Promise((r) => setTimeout(r, 800));
  await page.screenshot({ path: path.join(FRAMES_DIR, 'frame_06.png') });

  await browser.close();
  console.log('Browser session closed.');

  // Render GIF with FFmpeg
  console.log('Assembling governmesh_demo.gif using FFmpeg...');
  const gifPath = path.join(OUTPUT_DIR, 'governmesh_demo.gif');
  // Duplicate frames to create a nice 1.5s per frame loop
  execSync(
    `ffmpeg -y -framerate 0.65 -i "${path.join(FRAMES_DIR, 'frame_%02d.png')}" -vf "fps=10,scale=1080:-1:flags=lanczos,split[s0][s1];[s0]palettegen=max_colors=128[p];[s1][p]paletteuse=dither=bayer" "${gifPath}"`,
    { stdio: 'inherit' }
  );

  console.log('Cleaning up temporary frame files...');
  fs.rmSync(FRAMES_DIR, { recursive: true, force: true });

  console.log('All authentic screenshots and governmesh_demo.gif created successfully!');
}

capture().catch((err) => {
  console.error('Error during capture:', err);
  process.exit(1);
});
