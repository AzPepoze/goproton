const { contextBridge, ipcRenderer } = require("electron");

contextBridge.exposeInMainWorld("ipc", {
	invoke: (channel: string, ...args: any[]) => ipcRenderer.invoke(channel, ...args),
	on: (channel: string, listener: (event: any, ...args: any[]) => void) => {
		ipcRenderer.on(channel, listener);
		return () => ipcRenderer.removeListener(channel, listener);
	},
	off: (channel: string, listener: (...args: any[]) => void) => ipcRenderer.removeListener(channel, listener),
});
