import { writable } from "svelte/store";

export interface Notification {
	id: number;
	message: string;
	type: "info" | "error" | "success";
}

const { subscribe, update } = writable<Notification[]>([]);

export const notifications = {
	subscribe,
	add: (message: string, type: "info" | "error" | "success" = "info") => {
		const id = Date.now();
		update((n) => [...n, { id, message, type }]);
		setTimeout(() => {
			update((n) => n.filter((i) => i.id !== id));
		}, 5000);
	},
	remove: (id: number) => {
		update((n) => n.filter((i) => i.id !== id));
	},
	error: (message: string) => {
		const id = Date.now();
		update((n) => [...n, { id, message, type: "error" }]);
		setTimeout(() => {
			update((n) => n.filter((i) => i.id !== id));
		}, 5000);
	},
	success: (message: string) => {
		const id = Date.now();
		update((n) => [...n, { id, message, type: "success" }]);
		setTimeout(() => {
			update((n) => n.filter((i) => i.id !== id));
		}, 5000);
	},
	info: (message: string) => {
		const id = Date.now();
		update((n) => [...n, { id, message, type: "info" }]);
		setTimeout(() => {
			update((n) => n.filter((i) => i.id !== id));
		}, 5000);
	},
};
