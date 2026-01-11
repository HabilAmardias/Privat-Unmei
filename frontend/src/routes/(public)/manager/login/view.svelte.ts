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
	isLoading = $state<boolean>(false);
	isDesktop = $state<boolean>(false);

	loginDisabled = $derived.by<boolean>(() => {
		if (!this.email || !this.password || this.emailError || this.isLoading) {
			return true;
		}
		return false;
	});
	#validateEmail(email: string) {
		const pattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		return pattern.test(email);
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
