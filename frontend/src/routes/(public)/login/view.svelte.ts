export class AuthView {
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
	googleScript = $state<HTMLScriptElement>();
	isDesktop = $state<boolean>(false);

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

	removeGoogleScript() {
		if (this.googleScript) {
			this.googleScript.remove();
			this.googleScript = undefined;
		}
	}

	passwordOnBlur() {
		if (!this.#validatePassword(this.password)) {
			this.passwordError = new Error('min 8 characters with !@#?');
		}
		this.passwordError = undefined;
	}
	emailOnBlur() {
		if (!this.#validateEmail(this.email)) {
			this.emailError = new Error('please insert an valid email');
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
	nameOnBlur() {
		if (!this.name) {
			this.nameError = new Error('please insert a name');
			return;
		}
		this.nameError = undefined;
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
}
