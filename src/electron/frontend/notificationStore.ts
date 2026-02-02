import { writable, type Writable } from "svelte/store";

export interface Notification {
	id: number;
	message: string;
	type: "info" | "error" | "success";
}

interface NotificationStore extends Writable<Notification[]> {
	add(message: string, type?: "info" | "error" | "success"): void;
	remove(id: number): void;
	error(message: string): void;
	success(message: string): void;
	info(message: string): void;
}

function createNotifications(): NotificationStore {
	const base = writable<Notification[]>([]);

	function add(message: string, type: "info" | "error" | "success" = "info") {
		const id = Date.now();
		base.update((n) => [...n, { id, message, type }]);
		setTimeout(() => {
			base.update((n) => n.filter((i) => i.id !== id));
		}, 5000);
	}

	function remove(id: number) {
		base.update((n) => n.filter((i) => i.id !== id));
	}

	return {
		subscribe: base.subscribe,
		set: base.set,
		update: base.update,
		add,
		remove,
		error: (message: string) => add(message, "error"),
		success: (message: string) => add(message, "success"),
		info: (message: string) => add(message, "info"),
	};
}

export const notifications: NotificationStore = createNotifications();
