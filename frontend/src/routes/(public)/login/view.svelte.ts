import { IsAlphaOnly } from '$lib/utils/helper';
import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';
import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
export class AuthView {
	login = $state<boolean>(true);
	openDialog = $state<boolean>(false);
	agreed = $state<boolean>(false);
	email = $state<string>('');
	emailError = $derived.by<Error | undefined>(() => {
		if (this.email && !this.#validateEmail(this.email)) {
			return new Error('invalid email');
		}
		return undefined;
	});
	password = $state<string>('');
	passwordError = $derived.by<Error | undefined>(() => {
		if (this.password && !this.#validatePassword(this.password)) {
			return new Error('need at least 8 char with 1 special char');
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
	name = $state<string>('');
	nameError = $derived.by<Error | undefined>(() => {
		if (!this.name) {
			return new Error('please insert name');
		}
		if (this.name.replaceAll(' ', '').length === 0) {
			return new Error('please insert a valid name');
		}
		if (!IsAlphaOnly(this.name)) {
			return new Error('name can only consist of alphabet');
		}
		return undefined;
	});
	isLoading = $state<boolean>(false);
	googleScript = $state<HTMLScriptElement>();
	isDesktop = $state<boolean>(false);

	loginDisabled = $derived.by<boolean>(() => {
		if (!this.email || !this.password || this.emailError || this.isLoading) {
			return true;
		}
		return false;
	});

	registerDisabled = $derived.by<boolean>(() => {
		if (
			!this.name ||
			!this.password ||
			!this.repeatPassword ||
			!this.email ||
			this.repeatPasswordError ||
			this.nameError ||
			this.emailError ||
			this.passwordError ||
			this.isLoading ||
			!this.agreed
		) {
			return true;
		}
		return false;
	});

	removeGoogleScript() {
		if (this.googleScript) {
			this.googleScript.remove();
			this.googleScript = undefined;
		}
	}

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

	switchForm() {
		this.login = !this.login;
	}

	setIsDesktop(b: boolean) {
		this.isDesktop = b;
	}
	onLoginSubmit = (args: EnhancementArgs) => {
		if (args.action.search === '?/login') {
			this.setIsLoading(true);
			const loadID = CreateToast('loading', 'logging in....');
			return async ({ result, update }: EnhancementReturn) => {
				if (result.type === 'success') {
					await goto(resolve('/login-callback'), { replaceState: true });
					this.setIsLoading(false);
					DismissToast(loadID);
					CreateToast('success', 'login success');
				}
				if (result.type === 'failure') {
					this.setIsLoading(false);
					DismissToast(loadID);
					CreateToast('error', result.data?.message);
				}
				update();
			};
		}
	};
}
