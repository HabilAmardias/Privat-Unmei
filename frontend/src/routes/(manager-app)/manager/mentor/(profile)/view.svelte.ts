import type { EnhancementReturn } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';

export class MentorDetailView {
	isDesktop = $state<boolean>(false);
	openChangePassword = $state<boolean>(false);
	password = $state<string>('');
	passwordErr = $derived.by<Error | undefined>(() => {
		if (this.password) {
			if (!this.#validatePassword(this.password)) {
				return new Error('Need to at least contain 8 characters with number and special character');
			}
			return undefined;
		}
		return undefined;
	});
	changePasswordDisabled = $derived.by<boolean>(() => {
		if (!this.password || this.passwordErr) {
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
	onChangePassword = () => {
		const loadID = CreateToast('loading', 'updating....');
		return async ({ result }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'success') {
				CreateToast('success', 'successfully change password');
				this.openChangePassword = false;
			}
		};
	};
	#validatePassword(password: string) {
		const minLength = password.length >= 8;
		const hasSpecialChar =
			password.includes('!') ||
			password.includes('@') ||
			password.includes('#') ||
			password.includes('?');
		return minLength && hasSpecialChar;
	}
}
