class AuthView {
	login = $state<boolean>(true);
	email = $state<string>('');
	emailError = $state<Error | undefined>();
	password = $state<string>('');
	passwordError = $state<Error | undefined>();
	repeatPassword = $state<string>('');
	repeatPasswordError = $state<Error | undefined>();
	name = $state<string>('');
	nameError = $state<Error | undefined>();
	isLoading = $state<boolean>(false);
	loginDisabled = $derived.by<boolean>(() => {
		if (!this.email || !this.password || this.emailError || this.passwordError || this.isLoading) {
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
			this.isLoading
		) {
			return true;
		}
		return false;
	});
	passwordOnBlur() {
		this.passwordError = undefined;
		if (!this.#validatePassword(this.password)) {
			this.passwordError = new Error('min 8 characters with !@#?');
		}
	}
	emailOnBlur() {
		this.emailError = undefined;
		if (!this.#validateEmail(this.email)) {
			this.emailError = new Error('please insert an valid email');
		}
	}
	repeatPasswordOnBlur() {
		this.repeatPasswordError = undefined;
		if (this.repeatPassword !== this.password) {
			this.repeatPasswordError = new Error('password does not match');
		}
	}
	nameOnBlur() {
		this.nameError = undefined;
		if (!this.name) {
			this.nameError = new Error('please insert a name');
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
}

export const View = new AuthView();
