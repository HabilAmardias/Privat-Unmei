import type { AuthClaim, ToastType } from '$lib/types';
import toast, { type ToastOptions } from 'svelte-french-toast';

function DecodeJWT(token: string) {
	const claim: AuthClaim = JSON.parse(Buffer.from(token.split('.')[1], 'base64').toString());
	return claim;
}

export function IsTokenExpired(token: string | undefined) {
	if (!token) {
		return true;
	}
	const claim = DecodeJWT(token);
	const today = Math.floor(Date.now() / 1000);
	return claim.exp - today < 120;
}

export function IsAlphaOnly(str: string) {
	const reg = /^[a-zA-Z\s]+$/;
	return reg.test(str);
}

export function CreateToast(toastType: ToastType, message: string): string {
	let loadID: string;
	const Opts: ToastOptions = { position: 'top-right' };
	switch (toastType) {
		case 'error':
			loadID = toast.error(message, Opts);
			break;
		case 'success':
			loadID = toast.success(message, Opts);
			break;
		case 'loading':
			loadID = toast.loading(message, Opts);
			break;
	}
	return loadID;
}

export function DismissToast(toastID: string) {
	toast.dismiss(toastID);
}

export function debounce<T extends (...args: any[]) => any>(
	func: T,
	delay: number
): (...args: Parameters<T>) => void {
	let timeoutId: ReturnType<typeof setTimeout> | null = null;

	return function (this: ThisParameterType<T>, ...args: Parameters<T>) {
		if (timeoutId) {
			clearTimeout(timeoutId);
		}

		timeoutId = setTimeout(() => {
			func.apply(this, args);
			timeoutId = null; // Clear the timeoutId after execution
		}, delay);
	};
}
