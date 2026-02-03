import { notifications } from "../notificationStore";

export type NotificationType = "info" | "success" | "error";

/**
 * Shows a success notification
 */
export function showSuccess(message: string): void {
	notifications.success(message);
}

/**
 * Shows an error notification
 */
export function showError(message: string): void {
	notifications.error(message);
}

/**
 * Shows a warning notification (maps to error type)
 */
export function showWarning(message: string): void {
	notifications.error(message);
}

/**
 * Shows an info notification
 */
export function showInfo(message: string): void {
	notifications.info(message);
}

/**
 * Wrapper for async operations with notifications
 */
export async function withNotification<T>(
	promise: Promise<T>,
	options: {
		pending?: string;
		success?: string;
		error?: string;
	},
): Promise<T> {
	if (options.pending) showInfo(options.pending);

	try {
		const result = await promise;
		if (options.success) showSuccess(options.success);
		return result;
	} catch (err) {
		const errorMsg = options.error || String(err);
		showError(errorMsg);
		throw err;
	}
}
