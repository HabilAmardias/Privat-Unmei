class AuthView {
	cardState = $state<'Register' | 'Login'>('Login');
	email = $state<string>('');
	password = $state<string>('');
	repeatPassword = $state<string>('');
	name = $state<string>('');
	error = $state<Error | undefined>();
}

export const View = new AuthView();
