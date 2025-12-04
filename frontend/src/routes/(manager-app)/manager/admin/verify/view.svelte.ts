export class VerifyAdminView {
	email = $state<string>('');
	emailError = $derived.by<Error | undefined>(() => {
		if (this.email && !this.#validateEmail(this.email)) {
			return new Error('please insert an valid email');
		}
		return undefined;
	});
	password = $state<string>('');
	passwordError = $derived.by<Error | undefined>(() => {
		if (this.password && !this.#validatePassword(this.password)) {
			return new Error('min 8 characters with !@#?');
		}
		return undefined;
	});
	isLoading = $state<boolean>(false);
	verifyDisabled = $derived.by<boolean>(() => {
		if (!this.email || this.emailError || this.passwordError || this.isLoading) {
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
}
