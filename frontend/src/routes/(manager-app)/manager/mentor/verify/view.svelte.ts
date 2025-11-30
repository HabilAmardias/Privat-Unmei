export class VerifyMentorView {
	password = $state<string>('');
	passwordError = $derived.by<Error | undefined>(() => {
		if (this.password && !this.#validatePassword(this.password)) {
			return new Error('password need at least 8 chars with 1 special char');
		}
		return undefined;
	});
	isLoading = $state<boolean>(false);
	verifyDisabled = $derived.by<boolean>(() => {
		if (!this.password || this.passwordError || this.isLoading) {
			return true;
		}
		return false;
	});
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
