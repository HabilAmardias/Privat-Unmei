export class ResetView {
	email = $state<string>('');
	emailError = $state<Error | undefined>();
	password = $state<string>('');
	passwordError = $state<Error | undefined>();
	repeatPassword = $state<string>('');
	repeatPasswordError = $state<Error | undefined>();
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

	passwordOnBlur() {
		if (!this.#validatePassword(this.password)) {
			this.passwordError = new Error('min 8 characters with !@#?');
			return;
		}
		this.passwordError = undefined;
	}
	emailOnBlur() {
		if (!this.#validateEmail(this.email)) {
			this.emailError = new Error('please insert an valid email');
			return;
		}
		this.emailError = undefined;
	}
	repeatPasswordOnBlur() {
		if (this.repeatPassword !== this.password) {
			this.repeatPasswordError = new Error('password does not match');
			return;
		}
		this.repeatPasswordError = undefined;
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
