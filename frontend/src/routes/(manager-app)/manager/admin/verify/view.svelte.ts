export class VerifyAdminView {
	email = $state<string>('');
	emailError = $state<Error | undefined>();
	password = $state<string>('');
	passwordError = $state<Error | undefined>();
	isLoading = $state<boolean>(false);
	verifyDisabled = $derived.by<boolean>(() => {
		if (!this.email || this.emailError || this.passwordError || this.isLoading) {
			return true;
		}
		return false;
	});

	passwordOnBlur() {
		if (!this.#validatePassword(this.password)) {
			this.passwordError = new Error('min 8 characters with !@#?');
		} else {
            this.passwordError = undefined;
        }
	}
	emailOnBlur() {
		if (!this.#validateEmail(this.email)) {
			this.emailError = new Error('please insert an valid email');
		} else {
            this.emailError = undefined;
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
}
