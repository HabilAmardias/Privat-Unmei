class ResetView {
	email = $state<string>('');
	password = $state<string>('');
	repeatPassword = $state<string>('');
	isLoading = $state<boolean>(false);

	setIsLoading(b: boolean) {
		this.isLoading = b;
	}
}

export const View = new ResetView();
