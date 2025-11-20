export class VerifyMentorView {
	password = $state<string>('');
	passwordError = $state<Error | undefined>();
	isLoading = $state<boolean>(false);
	verifyDisabled = $derived.by<boolean>(() => {
		if (this.passwordError || this.isLoading) {
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
