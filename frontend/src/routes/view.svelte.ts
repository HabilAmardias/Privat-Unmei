class AuthView {
	login = $state<boolean>(true);
	email = $state<string>('');
	password = $state<string>('');
	repeatPassword = $state<string>('');
	name = $state<string>('');
	error = $state<Error | undefined>();

	switchForm() {
		this.login = !this.login;
	}
}

export const View = new AuthView();
