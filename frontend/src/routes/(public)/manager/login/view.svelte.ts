import type { ManagerLoginType } from './model';

export class ManagerAuthView {
	loginMenu = $state<ManagerLoginType>('admin');
	email = $state<string>('');
	emailError = $derived.by<Error | undefined>(() => {
		if (this.email && !this.#validateEmail(this.email)) {
			return new Error('please enter a valid email');
		}
		return undefined;
	});
	password = $state<string>('');
	passwordError = $derived.by<Error | undefined>(() => {
		if (this.password && !this.#validatePassword(this.password)) {
			return new Error('need at least 8 chars with 1 special char');
		}
		return undefined;
	});
	isLoading = $state<boolean>(false);
	isDesktop = $state<boolean>(false);

	loginDisabled = $derived.by<boolean>(() => {
		if (!this.email || !this.password || this.emailError || this.passwordError || this.isLoading) {
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

	switchForm(newState: ManagerLoginType) {
		this.loginMenu = newState;
	}

	setIsDesktop(b: boolean) {
		this.isDesktop = b;
	}
}
