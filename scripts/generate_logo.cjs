const fs = require('fs');
const path = require('path');
const { Resvg } = require('@resvg/resvg-js');

function buildLogoSvg() {
  return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024" width="1024" height="1024">
  <defs>
    <clipPath id="squircle-clip">
      <rect x="24" y="24" width="976" height="976" rx="220" />
    </clipPath>

    <linearGradient id="cyan-glow" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#00f5ff"/>
      <stop offset="100%" stop-color="#0284c7"/>
    </linearGradient>

    <linearGradient id="cyan-light" x1="0%" y1="0%" x2="0%" y2="100%">
      <stop offset="0%" stop-color="#38bdf8"/>
      <stop offset="100%" stop-color="#0284c7"/>
    </linearGradient>

    <linearGradient id="slate-facet-l" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#1e293b"/>
      <stop offset="100%" stop-color="#0f172a"/>
    </linearGradient>

    <linearGradient id="slate-facet-r" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#334155"/>
      <stop offset="100%" stop-color="#1e293b"/>
    </linearGradient>

    <filter id="subtle-shadow" x="-10%" y="-10%" width="120%" height="120%">
      <feDropShadow dx="0" dy="16" stdDeviation="20" flood-color="#000000" flood-opacity="0.16" />
    </filter>
  </defs>

  <!-- Luxury White Squircle Container -->
  <rect x="24" y="24" width="976" height="976" rx="220" fill="#ffffff" stroke="#e2e8f0" stroke-width="6" />

  <g clip-path="url(#squircle-clip)">
    <g transform="translate(512, 512)" filter="url(#subtle-shadow)">

      <!-- Hexagonal Architectural Gateway Frame -->
      <polygon points="
        0,-390
        338,-195
        338,195
        0,390
        -338,195
        -338,-195
      " fill="none" stroke="#0f172a" stroke-width="36" stroke-linejoin="round" />

      <polygon points="
        0,-355
        307,-177
        307,177
        0,355
        -307,177
        -307,-177
      " fill="none" stroke="#00f5ff" stroke-width="4" opacity="0.4" stroke-dasharray="16, 12" />

      <!-- 100% Watertight Solid Silhouette Underlayer (Zero-Leak Rule) -->
      <path d="
        M 0 -260
        L 130 -305
        L 205 -225
        L 165 -135
        L 255 -30
        L 245 110
        L 190 220
        L 120 200
        L 65 270
        L 0 300
        L -65 270
        L -120 200
        L -190 220
        L -245 110
        L -255 -30
        L -165 -135
        L -205 -225
        L -130 -305
        Z
      " fill="#0b0f19" />

      <!-- ==================== MASCOT FACETS (THE SOVEREIGN AEGIS OWL) ==================== -->

      <!-- Top Crown & Ear Tufts -->
      <!-- Left Ear Tuft -->
      <polygon points="0,-260 -130,-305 -165,-135" fill="#1e293b" />
      <polygon points="-130,-305 -205,-225 -165,-135" fill="#0f172a" />
      <!-- Right Ear Tuft -->
      <polygon points="0,-260 130,-305 165,-135" fill="#334155" />
      <polygon points="130,-305 205,-225 165,-135" fill="#1e293b" />

      <!-- Central Forehead Wedge -->
      <polygon points="0,-260 -55,-140 0,-105" fill="#1e293b" />
      <polygon points="0,-260 55,-140 0,-105" fill="#475569" />

      <!-- Outer Brow / Temples -->
      <polygon points="-165,-135 -255,-30 -140,15 -100,-45" fill="#0f172a" />
      <polygon points="165,-135 255,-30 140,15 100,-45" fill="#334155" />

      <!-- Eye Socket Cavities (Deep Obsidian) -->
      <polygon points="-55,-140 -165,-135 -100,-45 0,-60" fill="#0b0f19" />
      <polygon points="55,-140 165,-135 100,-45 0,-60" fill="#111827" />

      <!-- Focused Predatory Almond Optics (Electric Cyan) -->
      <!-- Left Eye -->
      <polygon points="-145,-55 -80,-52 -95,-28 -155,-32" fill="url(#cyan-glow)" />
      <polygon points="-132,-50 -92,-48 -102,-34 -140,-36" fill="#00f5ff" />
      <polygon points="-116,-48 -110,-48 -112,-36 -118,-36" fill="#061224" />

      <!-- Right Eye -->
      <polygon points="145,-55 80,-52 95,-28 155,-32" fill="url(#cyan-glow)" />
      <polygon points="132,-50 92,-48 102,-34 140,-36" fill="#00f5ff" />
      <polygon points="116,-48 110,-48 112,-36 118,-36" fill="#061224" />

      <!-- Cheeks / Zygomatic Facets -->
      <polygon points="-100,-45 -140,15 -65,45 0,-15" fill="#1e293b" />
      <polygon points="100,-45 140,15 65,45 0,-15" fill="#475569" />

      <!-- Sharp Obsidian Raptor Beak -->
      <polygon points="0,-15 -24,40 0,95" fill="#090d16" />
      <polygon points="0,-15 24,40 0,95" fill="#1e293b" />
      <polygon points="0,-15 0,95 10,40" fill="#334155" opacity="0.7" />

      <!-- ==================== GOVERNANCE MESH CHEST PLATES (LATTICE) ==================== -->
      <!-- Upper Collar -->
      <polygon points="0,95 -65,45 -115,130 0,165" fill="#0f172a" />
      <polygon points="0,95 65,45 115,130 0,165" fill="#334155" />

      <!-- Lateral Shoulder Armor Plates -->
      <polygon points="-140,15 -245,110 -185,190 -115,130" fill="#0b0f19" />
      <polygon points="140,15 245,110 185,190 115,130" fill="#1e293b" />

      <!-- Mid-Chest Governance Shield Facets -->
      <polygon points="0,165 -115,130 -85,220 0,230" fill="#1e293b" />
      <polygon points="0,165 115,130 85,220 0,230" fill="#475569" />

      <!-- Lower Wing Base Plates -->
      <polygon points="-115,130 -185,190 -120,200 -85,220" fill="#0f172a" />
      <polygon points="115,130 185,190 120,200 85,220" fill="#334155" />

      <!-- Lower Crest & Central Anchor -->
      <polygon points="0,230 -85,220 -45,270 0,300" fill="#0f172a" />
      <polygon points="0,230 85,220 45,270 0,300" fill="#1e293b" />

      <!-- Core Telemetry Node (Aegis Heart Cyan Accent) -->
      <polygon points="0,155 -16,180 0,205 16,180" fill="url(#cyan-glow)" />

    </g>
  </g>
</svg>`;
}

async function run() {
  const outputDir = path.join(__dirname, '..', 'docs', 'images');
  const svg = buildLogoSvg();
  const svgPath = path.join(outputDir, 'logo.svg');
  const pngPath = path.join(outputDir, 'logo.png');

  fs.writeFileSync(svgPath, svg, 'utf8');

  const resvg = new Resvg(svg, {
    fitTo: {
      mode: 'width',
      value: 1024,
    },
  });

  const pngData = resvg.render().asPng();
  fs.writeFileSync(pngPath, pngData);

  console.log('Successfully regenerated logo.svg and logo.png in docs/images/');
}

run().catch(console.error);
