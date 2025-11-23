import type { EnhancementReturn } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';
export class adminProfileView {
	isDesktop = $state<boolean>(false);
	isEdit = $state<boolean>(false);
	password = $state<string>('');
	repeatPassword = $state<string>('');
	passwordError = $state<Error | undefined>(undefined);
	repeatPasswordError = $state<Error | undefined>(undefined);
	isLoading = $state<boolean>(false);

	isDisabled = $derived.by<boolean>(() => {
		if (
			this.isLoading ||
			this.passwordError ||
			this.repeatPasswordError ||
			!this.password ||
			!this.repeatPassword
		) {
			return true;
		}
		return false;
	});

	size = $derived.by<number>(() => {
		if (this.isDesktop) {
			return 150;
		}
		return 100;
	});

	setIsDesktop(b: boolean) {
		this.isDesktop = b;
	}

	switchForm() {
		this.isEdit = !this.isEdit;
	}

	setIsLoading(b: boolean) {
		this.isLoading = b;
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

	passwordOnBlur() {
		if (!this.#validatePassword(this.password)) {
			this.passwordError = new Error('min 8 characters with !@#?');
		} else {
			this.passwordError = undefined;
		}
	}

	repeatPasswordOnBlur() {
		if (this.repeatPassword !== this.password) {
			this.repeatPasswordError = new Error('password does not match');
		} else {
			this.repeatPasswordError = undefined;
		}
	}
	onChangePasswordSubmit = () => {
		this.setIsLoading(true);
		const loadID = CreateToast('loading', 'loading....');
		return async ({ result }: EnhancementReturn) => {
			DismissToast(loadID);
			this.setIsLoading(false);
			if (result.type === 'success') {
				CreateToast('success', 'successfully change password');
				this.switchForm();
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
