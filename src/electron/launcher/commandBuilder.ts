import { LaunchOptions } from './types.js';
import { expandPath, isCommandAvailable } from './utils.js';

export function buildCommand(opts: LaunchOptions): { cmdArgs: string[], env: Record<string, string> } {
    const env: Record<string, string> = { ...process.env } as Record<string, string>;
    
    // Base Environment
    env['WINEPREFIX'] = expandPath(opts.PrefixPath);
    env['UMU_PROTON_PATTERN'] = opts.ProtonPattern;
    if (opts.ProtonPath) {
        env['PROTONPATH'] = expandPath(opts.ProtonPath);
    }

    const cmdArgs: string[] = [];

    // MangoHud and LSFG (Non-gamescope path)
    if (!opts.EnableGamescope) {
        if (opts.EnableMangoHud) {
            env['MANGOHUD'] = '1';
        }
        // LSFG-VK doesn't need extra env vars currently as per Go code
    }

    // Gamemode
    if (opts.EnableGamemode && isCommandAvailable('gamemoderun')) {
        cmdArgs.push('gamemoderun');
    }

    // Gamescope
    if (opts.EnableGamescope && isCommandAvailable('gamescope')) {
        cmdArgs.push('gamescope');
        if (opts.GamescopeW) cmdArgs.push('-W', opts.GamescopeW);
        if (opts.GamescopeH) cmdArgs.push('-H', opts.GamescopeH);
        if (opts.GamescopeR) cmdArgs.push('-r', opts.GamescopeR);
        
        cmdArgs.push('--', 'env');

        if (opts.EnableMangoHud) {
            env['MANGOHUD'] = '1';
        }
    }

    // umu-run
    cmdArgs.push('umu-run');

    // Executable
    const exePath = opts.LauncherPath || opts.GamePath;
    cmdArgs.push(exePath);

    // Custom Args
    if (opts.CustomArgs) {
        const parts = opts.CustomArgs.split(/\s+/).filter(p => p);
        cmdArgs.push(...parts);
    }

    // Systemd-run (Memory Limit) - Wrap the whole thing
    if (opts.EnableMemoryMin && opts.MemoryMinValue && isCommandAvailable('systemd-run')) {
        cmdArgs.unshift('systemd-run', '--user', '--scope', `-pMemoryMin=${opts.MemoryMinValue}`, '--');
    }

    return { cmdArgs, env };
}

export function formatCommandForDisplay(cmdArgs: string[], opts: LaunchOptions): string {
    let prefix = '';
    if (opts.EnableMemoryMin && opts.MemoryMinValue) {
        prefix += `[MemMin:${opts.MemoryMinValue}] `;
    }
    const envVars = `WINEPREFIX=${opts.PrefixPath} UMU_PROTON_PATTERN=${opts.ProtonPattern} ${opts.EnableMangoHud ? 'MANGOHUD=1 ' : ''}`;
    return `${prefix}${envVars}${cmdArgs.join(' ')}`;
}
