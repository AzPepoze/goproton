import { ipcMain as f, app as d, dialog as v, BrowserWindow as M, shell as k } from "electron";
import l from "path";
import a from "fs";
import { fileURLToPath as T } from "url";
import { exec as x, execSync as C, spawn as N } from "child_process";
import h from "os";
import H from "crypto";
import { promisify as O } from "util";
import I from "https";
function V() {
  f.handle("GetInitialLauncherPath", () => process.env.GOPROTON_LAUNCHER_PATH || ""), f.handle("GetInitialGamePath", () => process.env.GOPROTON_GAME_PATH || ""), f.handle("GetShouldEditLsfg", () => process.env.GOPROTON_EDIT_LSFG === "1"), f.handle("CloseWindow", () => d.quit());
}
function U(n) {
  f.handle("PickFile", async () => {
    if (!n)
      return "";
    const e = await v.showOpenDialog(n, {
      title: "Select Game Executable",
      filters: [
        { name: "Executables", extensions: ["exe"] },
        { name: "All Files", extensions: ["*"] }
      ],
      properties: ["openFile"]
    });
    return e.canceled ? "" : e.filePaths[0];
  }), f.handle("PickFolder", async () => {
    if (!n)
      return "";
    const e = await v.showOpenDialog(n, {
      title: "Select Prefix Directory",
      properties: ["openDirectory"]
    });
    return e.canceled ? "" : e.filePaths[0];
  }), f.handle("PickFileCustom", async (e, t) => {
    if (!n)
      return "";
    const r = await v.showOpenDialog(n, {
      title: t || "Select File",
      properties: ["openFile"]
    });
    return r.canceled ? "" : r.filePaths[0];
  });
}
function z() {
  f.handle("GetExeIcon", async (n, e) => {
    if (!e || !a.existsSync(e))
      return console.warn("Icon path does not exist:", e), "";
    try {
      return await B(e) || "";
    } catch (t) {
      return console.error("Failed to extract PE icon:", t), "";
    }
  });
}
function B(n) {
  return new Promise((e) => {
    const t = l.join(h.tmpdir(), `goproton-icon-${Date.now()}`);
    a.mkdirSync(t, { recursive: !0 }), x(`wrestool -x --output="${t}" "${n}"`, (i) => {
      i ? (console.log("wrestool failed, trying icoextract..."), r()) : a.readdir(t, (o, c) => {
        if (!o && c.length > 0) {
          const s = c.find((u) => u.endsWith(".ico"));
          if (s) {
            const u = l.join(t, s);
            a.readFile(u, (p, m) => {
              if (a.rm(t, { recursive: !0 }, () => {
              }), !p && m.length > 0) {
                console.log(`Successfully extracted icon from ${n} using wrestool`), e("data:image/x-icon;base64," + m.toString("base64"));
                return;
              }
              r();
            });
            return;
          }
        }
        r();
      });
    });
    function r() {
      const i = l.join(t, "icon.ico");
      x(`icoextract "${n}" "${i}"`, (o) => {
        o ? (console.warn(`icoextract failed: ${o.message}`), a.rm(t, { recursive: !0 }, () => {
        }), e(null)) : a.readFile(i, (c, s) => {
          a.rm(t, { recursive: !0 }, () => {
          }), !c && s.length > 0 ? (console.log(`Successfully extracted icon from ${n} using icoextract`), e("data:image/x-icon;base64," + s.toString("base64"))) : (console.warn(`icoextract succeeded but failed to read file: ${c}`), e(null));
        });
      });
    }
  });
}
function E() {
  return l.join(h.homedir(), "GoProton");
}
function q() {
  return l.join(E(), "config", "executables");
}
function y() {
  return l.join(E(), "prefixes");
}
function J(n) {
  const e = H.createHash("sha1").update(n).digest("hex").slice(0, 8), r = `${l.basename(n).replace(/\.[^/.]+$/, "")}-${e}`;
  return l.join(q(), r);
}
function L(n) {
  return l.join(J(n), "config.json");
}
function b(n) {
  return l.join(y(), n, "goproton.json");
}
function _(n) {
  try {
    return C(`command -v ${n}`, { stdio: "ignore" }), !0;
  } catch {
    return !1;
  }
}
function K() {
  try {
    const n = C("vulkaninfo", { stdio: ["ignore", "pipe", "ignore"] }).toString(), e = /* @__PURE__ */ new Set(), t = /GPU\s+id\s*=\s*\d+\s*\((.+)\)/g;
    let r;
    for (; (r = t.exec(n)) !== null; )
      r[1] && e.add(r[1].trim());
    return Array.from(e);
  } catch {
    return [];
  }
}
function Q() {
  return {
    hasGamescope: _("gamescope"),
    hasMangoHud: _("mangohud"),
    hasGameMode: _("gamemoderun")
  };
}
function F() {
  const n = "/usr/share/vulkan/implicit_layer.d";
  return a.existsSync(l.join(n, "VkLayer_LSFGVK_frame_generation.json"));
}
function Z() {
  var e;
  let n = "unknown";
  if (F())
    try {
      const t = "/usr/share/vulkan/implicit_layer.d/VkLayer_LSFGVK_frame_generation.json";
      n = ((e = JSON.parse(a.readFileSync(t, "utf8")).layer) == null ? void 0 : e.implementation_version) || "1.0.0";
    } catch (t) {
      console.error("Failed to read LSFG version:", t), n = "1.0.0";
    }
  return {
    isLsfgInstalled: F(),
    lsfgVersion: n
  };
}
function X() {
  const n = h.homedir(), e = [
    { path: l.join(n, ".steam/root/compatibilitytools.d"), isSteam: !1 },
    { path: l.join(n, ".local/share/Steam/compatibilitytools.d"), isSteam: !1 },
    { path: "/usr/share/steam/compatibilitytools.d", isSteam: !1 },
    { path: l.join(n, ".steam/root/steamapps/common"), isSteam: !0 },
    { path: l.join(n, ".local/share/Steam/steamapps/common"), isSteam: !0 },
    { path: l.join(n, "GoProton/protons"), isSteam: !1 }
  ], t = [], r = /* @__PURE__ */ new Set();
  for (const i of e)
    if (a.existsSync(i.path))
      try {
        const o = a.realpathSync(i.path), c = a.readdirSync(o, { withFileTypes: !0 });
        for (const s of c) {
          if (!s.isDirectory())
            continue;
          const u = l.join(o, s.name);
          r.has(u) || a.existsSync(l.join(u, "proton")) && (r.add(u), t.push({
            Name: s.name,
            Path: u,
            IsSteam: i.isSteam,
            DisplayName: i.isSteam ? `(Steam) ${s.name}` : s.name
          }));
        }
      } catch (o) {
        console.error(`Error scanning ${i.path}:`, o);
      }
  return t.sort((i, o) => i.IsSteam && !o.IsSteam ? -1 : !i.IsSteam && o.IsSteam ? 1 : i.Name.localeCompare(o.Name));
}
function Y() {
  f.handle("GetTotalRam", () => {
    const n = h.totalmem();
    return Math.round(n / (1024 * 1024 * 1024));
  }), f.handle("GetListGpus", () => K()), f.handle("DetectLosslessDll", () => {
    const n = h.homedir(), e = [
      l.join(n, ".steam/root/steamapps/common/Lossless Scaling/Lossless.dll"),
      l.join(n, ".local/share/Steam/steamapps/common/Lossless Scaling/Lossless.dll")
    ];
    for (const t of e)
      if (a.existsSync(t))
        return t;
    return "";
  });
}
function ee() {
  f.handle("GetConfig", async (n, e) => {
    if (!e)
      return null;
    const t = L(e), r = e + ".goproton.json", i = a.existsSync(t) ? t : a.existsSync(r) ? r : null;
    if (i)
      try {
        const o = a.readFileSync(i, "utf8");
        return JSON.parse(o);
      } catch (o) {
        return console.error(o), null;
      }
    return null;
  }), f.handle("SaveConfig", async (n, e) => {
    try {
      const t = e.LauncherPath || e.GamePath;
      if (!t)
        return "Missing Game/Launcher Path";
      const r = L(t);
      return a.mkdirSync(l.dirname(r), { recursive: !0 }), a.writeFileSync(r, JSON.stringify(e, null, 2)), null;
    } catch (t) {
      return t.toString();
    }
  });
}
function te() {
  f.handle("GetPrefixBaseDir", () => y()), f.handle("SavePrefixConfig", async (n, e, t) => {
    try {
      const r = b(e);
      return a.mkdirSync(l.dirname(r), { recursive: !0 }), a.writeFileSync(r, JSON.stringify(t, null, 2)), null;
    } catch (r) {
      return r.toString();
    }
  }), f.handle("GetPrefixConfig", async (n, e) => {
    const t = b(e);
    if (a.existsSync(t))
      try {
        return JSON.parse(a.readFileSync(t, "utf8"));
      } catch {
        return null;
      }
    return null;
  }), f.handle("LoadPrefixConfig", async (n, e) => {
    const t = b(e);
    if (a.existsSync(t))
      try {
        return JSON.parse(a.readFileSync(t, "utf8"));
      } catch {
        return null;
      }
    return null;
  }), f.handle("ListPrefixes", () => {
    const n = y();
    a.existsSync(n) || a.mkdirSync(n, { recursive: !0 });
    try {
      const t = a.readdirSync(n, { withFileTypes: !0 }).filter((r) => r.isDirectory()).map((r) => r.name);
      return t.length === 0 ? (a.mkdirSync(l.join(n, "Default"), { recursive: !0 }), ["Default"]) : t;
    } catch {
      return ["Default"];
    }
  }), f.handle("CreatePrefix", (n, e) => {
    try {
      const t = l.join(y(), e);
      return a.mkdirSync(t, { recursive: !0 }), null;
    } catch (t) {
      return t.toString();
    }
  });
}
function ne(n, e) {
  f.handle("RunGame", async (t, r, i) => {
    try {
      const o = typeof r == "string" ? JSON.parse(r) : r;
      let c = "";
      d.isPackaged ? c = l.join(process.resourcesPath, "goproton-instance") : c = l.resolve(e, "../../bin/goproton-instance"), a.existsSync(c) || (console.error(`Instance binary not found at: ${c}`), c = "goproton-instance");
      const s = [
        `--game=${o.GamePath}`,
        `--launcher=${o.LauncherPath || ""}`,
        `--prefix=${o.PrefixPath}`,
        `--proton-pattern=${o.ProtonPattern || ""}`,
        `--proton-path=${o.ProtonPath || ""}`
      ];
      return o.EnableMangoHud && s.push("--mango"), o.EnableGamemode && s.push("--gamemode"), o.EnableGamescope && (s.push("--gamescope"), o.GamescopeW && s.push(`--gs-w=${o.GamescopeW}`), o.GamescopeH && s.push(`--gs-h=${o.GamescopeH}`), o.GamescopeR && s.push(`--gs-r=${o.GamescopeR}`)), o.EnableLsfgVk && (s.push("--lsfg"), s.push(`--lsfg-mult=${o.LsfgMultiplier || "2"}`), o.LsfgPerfMode && s.push("--lsfg-perf"), o.LsfgDllPath && s.push(`--lsfg-dll-path=${o.LsfgDllPath}`)), o.EnableMemoryMin && (s.push("--memory-min"), o.MemoryMinValue && s.push(`--memory-min-value=${o.MemoryMinValue}`)), i === !0 ? s.push("--logs=true") : i === !1 && s.push("--logs=false"), console.log("Spawning instance manager:", c, s.join(" ")), N(c, s, {
        detached: !0,
        stdio: "ignore"
      }).unref(), d.quit(), null;
    } catch (o) {
      return console.error("RunGame Error:", o), o.toString();
    }
  });
}
function oe() {
  f.handle("ScanProtonVersions", () => X());
}
const S = O(x), re = "PancakeTAS/lsfg-vk", ie = [
  {
    ID: "ge-proton",
    Name: "GE-Proton (GloriousEggroll)",
    Description: "The most popular custom Proton build. Includes many game fixes and codec patches.",
    RepoOwner: "GloriousEggroll",
    RepoName: "proton-ge-custom"
  },
  {
    ID: "proton-cachyos",
    Name: "Proton-CachyOS",
    Description: "Optimized for performance with CachyOS patches and schedulers.",
    RepoOwner: "CachyOS",
    RepoName: "proton-cachyos"
  },
  {
    ID: "kron4ek",
    Name: "Proton-Kron4ek",
    Description: "Vanilla builds and TKG builds. Often smaller and faster updates.",
    RepoOwner: "Kron4ek",
    RepoName: "Proton-Builds"
  },
  {
    ID: "luxtorpeda",
    Name: "Luxtorpeda (Native Tools)",
    Description: "Runs Windows games using native Linux engines (e.g. GZDoom, ScummVM).",
    RepoOwner: "luxtorpeda-dev",
    RepoName: "luxtorpeda"
  }
];
async function ae(n) {
  return new Promise((e, t) => {
    I.get(n, { headers: { "User-Agent": "GoProton-App" } }, (r) => {
      let i = "";
      r.on("data", (o) => i += o), r.on("end", () => {
        r.statusCode !== 200 ? t(new Error(`Failed to fetch: ${r.statusCode}`)) : e(JSON.parse(i));
      });
    }).on("error", t);
  });
}
function $(n, e, t) {
  return new Promise((r, i) => {
    I.get(n, { headers: { "User-Agent": "GoProton-App" } }, (o) => {
      if (o.statusCode === 302 || o.statusCode === 301) {
        $(o.headers.location, e, t).then(r).catch(i);
        return;
      }
      if (o.statusCode !== 200) {
        i(new Error(`Server returned ${o.statusCode}`));
        return;
      }
      const c = parseInt(o.headers["content-length"] || "0", 10);
      let s = 0;
      const u = a.createWriteStream(e);
      o.on("data", (p) => {
        s += p.length, u.write(p), t(s, c);
      }), o.on("end", () => {
        u.end(), r();
      }), o.on("error", (p) => {
        u.close(), a.unlinkSync(e), i(p);
      });
    }).on("error", i);
  });
}
async function se(n) {
  n(0, "Fetching release info from GitHub...");
  const e = await ae(`https://api.github.com/repos/${re}/releases`);
  let t = "", r = "";
  for (const s of e) {
    for (const u of s.assets) {
      const p = u.name.toLowerCase();
      if (p.includes("x86_64") && p.endsWith(".tar.zst") || p.includes("linux") && p.endsWith(".tar.xz")) {
        t = u.browser_download_url, r = u.name;
        break;
      }
    }
    if (t)
      break;
  }
  if (!t)
    throw new Error("lsfg-vk suitable linux asset not found");
  n(5, `Downloading ${r}...`);
  const i = r.endsWith(".tar.zst") ? ".tar.zst" : ".tar.xz", o = l.join(h.tmpdir(), `lsfg-vk-dl${i}`);
  await $(t, o, (s, u) => {
    const p = Math.round(s / u * 80);
    n(5 + p, "Downloading...");
  }), n(85, "Extracting files...");
  const c = a.mkdtempSync(l.join(h.tmpdir(), "lsfg-extract-"));
  try {
    const s = i === ".tar.zst" ? ["--use-compress-program=unzstd", "-xf", o, "-C", c] : ["-xf", o, "-C", c];
    try {
      await S(`tar ${s.join(" ")}`);
    } catch {
      await S(`tar -xf ${o} -C ${c}`);
    }
    n(88, "Installing to system directories (requires sudo)..."), await S(`pkexec sh -c "cp -r ${c}/* /usr"`), n(100, "Installation complete!");
  } finally {
    a.rmSync(c, { recursive: !0, force: !0 }), a.existsSync(o) && a.unlinkSync(o);
  }
}
async function le(n, e) {
  const t = h.homedir();
  let r = l.join(t, ".steam/root/compatibilitytools.d");
  a.existsSync(l.join(t, ".steam/root")) || (r = l.join(t, ".local/share/Steam/compatibilitytools.d")), a.existsSync(r) || a.mkdirSync(r, { recursive: !0 }), e(0, "Downloading...");
  const i = l.join(h.tmpdir(), `proton-dl-${Date.now()}.tar.xz`);
  try {
    await $(n, i, (o, c) => {
      const s = Math.round(o / c * 50);
      e(s, `Downloading... ${Math.round(o / c * 100)}%`);
    }), e(50, "Extracting (using system tools)..."), await S(`tar -xf ${i} -C "${r}"`), e(100, "Installation Complete!");
  } finally {
    a.existsSync(i) && a.unlinkSync(i);
  }
}
const D = O(x);
function ce(n) {
  f.handle("CleanupProcesses", async () => {
    try {
      return await Promise.all(
        ["umu-run", "reaper", "gamescopereaper"].map((e) => D(`pkill -f ${e}`).catch(() => {
        }))
      ), null;
    } catch (e) {
      return e.toString();
    }
  }), f.handle("InstallLsfg", async () => {
    try {
      return await se((e, t) => {
        n == null || n.webContents.send("lsfg-install-progress", { percent: e, message: t });
      }), null;
    } catch (e) {
      throw new Error(e.toString());
    }
  }), f.handle("InstallProton", async (e, t) => {
    try {
      return await le(t, (r, i) => {
        n == null || n.webContents.send("lsfg-install-progress", { percent: r, message: i });
      }), null;
    } catch (r) {
      throw new Error(r.toString());
    }
  }), f.handle("GetUtilsStatus", () => Z()), f.handle("GetSystemToolsStatus", () => Q()), f.handle("GetProtonVariants", () => ie), f.handle("UninstallLsfg", async () => {
    try {
      const e = [
        "rm -f /usr/lib/liblsfg-vk-layer.so",
        "rm -f /usr/share/vulkan/implicit_layer.d/VkLayer_LSFGVK_frame_generation.json",
        "rm -f /usr/share/icons/hicolor/256x256/apps/gay.pancake.lsfg-vk-ui.png",
        "rm -f /usr/share/applications/gay.pancake.lsfg-vk-ui.desktop",
        "rm -f /usr/bin/lsfg-vk-cli",
        "rm -f /usr/bin/lsfg-vk-ui"
      ].join(" && ");
      return await D(`pkexec sh -c "${e}"`), null;
    } catch (e) {
      return e.toString();
    }
  });
}
function R() {
  return l.join(h.homedir(), ".config", "lsfg-vk", "conf.toml");
}
function W(n) {
  const e = n.split(`
`), t = { version: 2, global: {}, profile: [] };
  let r = "", i = null;
  for (let o of e) {
    if (o = o.trim(), !o || o.startsWith("#"))
      continue;
    if (o.startsWith("[") && o.endsWith("]")) {
      r = o.slice(1, -1), r === "profile" && (i && t.profile.push(i), i = {});
      continue;
    }
    const [c, ...s] = o.split("=");
    if (!c)
      continue;
    const u = s.join("=").trim(), p = c.trim();
    let m = u;
    u === "true" ? m = !0 : u === "false" ? m = !1 : u.startsWith('"') && u.endsWith('"') ? m = u.slice(1, -1) : isNaN(Number(u)) ? u.startsWith("[") && u.endsWith("]") && (m = u.slice(1, -1).split(",").map((w) => w.trim().replace(/^"|"$/g, "")).filter((w) => w)) : m = Number(u), r === "global" ? t.global[p] = m : r === "profile" ? i[p] = m : t[p] = m;
  }
  return i && t.profile.push(i), t;
}
function ue(n) {
  let e = `version = ${n.version}

`;
  e += `[global]
`, e += `version = ${n.global.version}
`, e += `allow_fp16 = ${n.global.allow_fp16}
`, e += `dll = "${n.global.dll}"

`;
  for (const t of n.profile)
    e += `[[profile]]
`, e += `name = "${t.name}"
`, Array.isArray(t.active_in) ? e += `active_in = [${t.active_in.map((r) => `"${r}"`).join(", ")}]
` : e += `active_in = "${t.active_in}"
`, e += `multiplier = ${t.multiplier}
`, e += `performance_mode = ${t.performance_mode}
`, e += `gpu = "${t.gpu}"
`, e += `flow_scale = ${t.flow_scale.toFixed(1)}
`, e += `pacing = "${t.pacing}"

`;
  return e;
}
function A(n) {
  const e = R();
  if (!a.existsSync(e))
    return { profile: null, index: -1 };
  try {
    const t = a.readFileSync(e, "utf8"), r = W(t), i = l.basename(n).toLowerCase();
    for (let o = 0; o < r.profile.length; o++) {
      const c = r.profile[o];
      if ((Array.isArray(c.active_in) ? c.active_in : [c.active_in]).some((u) => u.toLowerCase() === i))
        return { profile: c, index: o };
    }
  } catch (t) {
    console.error("Failed to parse LSFG config:", t);
  }
  return { profile: null, index: -1 };
}
function fe(n, e) {
  const t = R(), r = l.dirname(t);
  a.existsSync(r) || a.mkdirSync(r, { recursive: !0 });
  let i = {
    version: 2,
    global: { version: 2, allow_fp16: e.allowFp16 ?? !0, dll: e.dllPath ?? "" },
    profile: []
  };
  if (a.existsSync(t))
    try {
      i = W(a.readFileSync(t, "utf8"));
    } catch {
      console.error("Failed to read existing LSFG config, creating new one");
    }
  e.allowFp16 !== void 0 && (i.global.allow_fp16 = e.allowFp16), e.dllPath !== void 0 && (i.global.dll = e.dllPath);
  const o = l.basename(n), { index: c } = A(n), s = {
    name: o.replace(/\.[^/.]+$/, ""),
    active_in: o,
    multiplier: e.multiplier ?? 2,
    performance_mode: e.performance_mode ?? !1,
    gpu: e.gpu ?? "auto",
    flow_scale: e.flow_scale ?? 1,
    pacing: "none"
  };
  c !== -1 ? i.profile[c] = { ...i.profile[c], ...s } : i.profile.push(s), a.writeFileSync(t, ue(i));
}
function pe() {
  f.handle("GetLsfgProfileForGame", (n, e) => {
    const { profile: t } = A(e);
    return t;
  }), f.handle("SaveLsfgProfile", (n, e, t) => {
    try {
      return fe(e, t), null;
    } catch (r) {
      return r.toString();
    }
  }), f.handle("RemoveProfile", (n, e) => {
    try {
      const t = l.join(h.homedir(), ".config", "lsfg-vk", "conf.toml");
      return a.existsSync(t) && console.log("RemoveProfile called for:", e), null;
    } catch (t) {
      return t.toString();
    }
  });
}
function de(n, e) {
  V(), U(n), z(), Y(), ee(), te(), ne(n, e), oe(), ce(n), pe();
}
const me = T(import.meta.url), P = l.dirname(me);
d.commandLine.appendSwitch("--js-flags", "--max-old-space-size=512");
d.commandLine.appendSwitch("--no-zygote");
d.commandLine.appendSwitch("--no-sandbox");
const G = !d.isPackaged, he = l.join(P, "preload.js"), j = G ? "http://localhost:5173" : l.join(P, "../frontend/index.html");
let g;
const ge = () => {
  g = new M({
    width: 1024,
    height: 768,
    backgroundColor: "#18181b",
    webPreferences: {
      preload: he,
      contextIsolation: !0,
      nodeIntegration: !1,
      backgroundThrottling: !0,
      spellcheck: !1,
      offscreen: !1,
      enableWebSQL: !1
    },
    title: "GoProton",
    icon: l.join(P, G ? "../../src/public" : "../frontend", "favicon.ico")
  }), de(g, P), g.setAutoHideMenuBar(!0), g.setMenuBarVisibility(!1), g.webContents.setWindowOpenHandler(({ url: n }) => ((n.startsWith("http://") || n.startsWith("https://")) && k.openExternal(n), { action: "deny" })), g.webContents.on("new-window", (n, e) => {
    (e.startsWith("http://") || e.startsWith("https://")) && (n.preventDefault(), k.openExternal(e));
  }), G ? g.loadURL(j) : g.loadFile(j);
};
d.on("window-all-closed", () => {
  d.quit();
});
d.whenReady().then(() => {
  const n = process.argv;
  let e = !1;
  if (d.isPackaged ? n.length > 1 && n[1] === "--debug" && (e = !0) : n.length > 2 && n[2] === "--debug" && (e = !0), process.env.RUN_FROM_GOPROTON === "true") {
    ge();
    return;
  }
  const t = l.join(d.getAppPath(), "../goproton");
  if (a.existsSync(t)) {
    console.log("Starting goproton process:", t);
    const r = N(t, [], {
      stdio: "inherit",
      detached: !1,
      env: {
        ...process.env
      }
    });
    r.on("close", (i) => {
      console.log("goproton process exited with code:", i), e ? console.log("Debug mode: keeping app alive, press Ctrl+C to exit") : d.quit();
    }), r.on("error", (i) => {
      console.error("Failed to run goproton:", i), e || d.quit();
    });
  } else
    console.error("goproton binary not found at:", t), d.quit();
});
