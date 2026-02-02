import os from 'os';
import path from 'path';
import fs from 'fs';
import crypto from 'crypto';
import { execSync } from 'child_process';
import { ProtonTool, SystemToolsStatus, UtilsStatus } from './types.js';

export function expandPath(p: string): string {
    if (p === '~' || p.startsWith('~/')) {
        const home = os.homedir();
        if (p === '~') return home;
        return path.join(home, p.slice(2));
    }
    return p;
}

export function getBaseDir(): string {
    return path.join(os.homedir(), 'GoProton');
}

export function getConfigDir(): string {
    return path.join(getBaseDir(), 'config', 'executables');
}

export function getPrefixBaseDir(): string {
    return path.join(getBaseDir(), 'prefixes');
}

export function getConfigPath(exePath: string): string {
    const hash = crypto.createHash('sha1').update(exePath).digest('hex').slice(0, 8);
    const baseName = path.basename(exePath).replace(/\.[^/.]+$/, "");
    const folderName = `${baseName}-${hash}`;
    return path.join(getConfigDir(), folderName);
}

export function getGameConfigFile(exePath: string): string {
    return path.join(getConfigPath(exePath), 'config.json');
}

export function getPrefixConfigPath(prefixName: string): string {
    return path.join(getPrefixBaseDir(), prefixName, 'goproton.json');
}

export function isCommandAvailable(name: string): boolean {
    try {
        execSync(`command -v ${name}`, { stdio: 'ignore' });
        return true;
    } catch {
        return false;
    }
}

export function getListGpus(): string[] {
    try {
        const output = execSync('vulkaninfo', { stdio: ['ignore', 'pipe', 'ignore'] }).toString();
        const gpus = new Set<string>();
        const gpuPattern = /GPU\s+id\s*=\s*\d+\s*\((.+)\)/g;
        let match;
        while ((match = gpuPattern.exec(output)) !== null) {
            if (match[1]) gpus.add(match[1].trim());
        }
        return Array.from(gpus);
    } catch {
        return [];
    }
}

export function getSystemToolsStatus(): SystemToolsStatus {
    return {
        hasGamescope: isCommandAvailable('gamescope'),
        hasMangoHud: isCommandAvailable('mangohud'),
        hasGameMode: isCommandAvailable('gamemoderun'),
    };
}

export function isLsfgInstalled(): boolean {
    const vulkanDir = '/usr/share/vulkan/implicit_layer.d';
    return fs.existsSync(path.join(vulkanDir, 'VkLayer_LSFGVK_frame_generation.json'));
}

export function getUtilsStatus(): UtilsStatus {
    let version = 'unknown';
    if (isLsfgInstalled()) {
        try {
            const manifestPath = '/usr/share/vulkan/implicit_layer.d/VkLayer_LSFGVK_frame_generation.json';
            const manifest = JSON.parse(fs.readFileSync(manifestPath, 'utf8'));
            version = manifest.layer?.implementation_version || '1.0.0';
        } catch (e) {
            console.error("Failed to read LSFG version:", e);
            version = '1.0.0';
        }
    }
    return {
        isLsfgInstalled: isLsfgInstalled(),
        lsfgVersion: version,
    };
}

export function scanProtonVersions(): ProtonTool[] {
    const home = os.homedir();
    const searchPaths = [
        { path: path.join(home, '.steam/root/compatibilitytools.d'), isSteam: false },
        { path: path.join(home, '.local/share/Steam/compatibilitytools.d'), isSteam: false },
        { path: '/usr/share/steam/compatibilitytools.d', isSteam: false },
        { path: path.join(home, '.steam/root/steamapps/common'), isSteam: true },
        { path: path.join(home, '.local/share/Steam/steamapps/common'), isSteam: true },
        { path: path.join(home, 'GoProton/protons'), isSteam: false },
    ];

    const tools: ProtonTool[] = [];
    const seen = new Set<string>();

    for (const sp of searchPaths) {
        if (!fs.existsSync(sp.path)) continue;
        try {
            const realBase = fs.realpathSync(sp.path);
            
            const entries = fs.readdirSync(realBase, { withFileTypes: true });
            for (const entry of entries) {
                if (!entry.isDirectory()) continue;
                const fullPath = path.join(realBase, entry.name);
                
                if (seen.has(fullPath)) continue;

                if (fs.existsSync(path.join(fullPath, 'proton'))) {
                    seen.add(fullPath);
                    tools.push({
                        Name: entry.name,
                        Path: fullPath,
                        IsSteam: sp.isSteam,
                        DisplayName: sp.isSteam ? `(Steam) ${entry.name}` : entry.name
                    });
                }
            }
        } catch (e) {
            console.error(`Error scanning ${sp.path}:`, e);
        }
    }

    return tools.sort((a, b) => {
        if (a.IsSteam && !b.IsSteam) return -1;
        if (!a.IsSteam && b.IsSteam) return 1;
        return a.Name.localeCompare(b.Name);
    });
}
