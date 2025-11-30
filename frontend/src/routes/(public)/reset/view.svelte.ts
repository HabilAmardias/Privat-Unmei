import type { EnhancementReturn } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';

export class ResetView {
	email = $state<string>('');
	emailError = $derived.by<Error | undefined>(() => {
		if (this.email && !this.#validateEmail(this.email)) {
			return new Error('please insert a valid email');
		}
		return undefined;
	});
	password = $state<string>('');
	passwordError = $derived.by<Error | undefined>(() => {
		if (this.password && !this.#validatePassword(this.password)) {
			return new Error('password need at least 8 characters with one special character');
		}
		return undefined;
	});
	repeatPassword = $state<string>('');
	repeatPasswordError = $derived.by<Error | undefined>(() => {
		if (this.password !== this.repeatPassword) {
			return new Error('must be same as password');
		}
		return undefined;
	});
	isLoading = $state<boolean>(false);
	sendDisabled = $derived.by<boolean>(() => {
		if (!this.email || this.emailError || this.isLoading) {
			return true;
		}
		return false;
	});
	resetDisabled = $derived.by<boolean>(() => {
		if (!this.password || !this.repeatPassword || this.passwordError || this.repeatPasswordError) {
			return true;
		}
		return false;
	});
	#validateEmail(email: string) {
		const pattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		return pattern.test(email);
	}
	#validatePassword(password: string) {
		const minLength = password.length >= 8;
		const hasSpecialChar =
			password.includes('!') ||
			password.includes('@') ||
			password.includes('#') ||
			password.includes('?');
		return minLength && hasSpecialChar;
	}
	setIsLoading(b: boolean) {
		this.isLoading = b;
	}
	onSendSubmit = () => {
		this.setIsLoading(true);
		const loadID = CreateToast('loading', 'loading....');
		return async ({ result }: EnhancementReturn) => {
			this.setIsLoading(false);
			DismissToast(loadID);
			if (result.type === 'success') {
				CreateToast('success', result.data?.message);
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
