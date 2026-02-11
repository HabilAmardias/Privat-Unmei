import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';

export class LoginCallbackView {
	otp = $state<string>('');
	loginDisabled = $derived<boolean>(this.otp.length < 6);

	onLoginSubmit = (args: EnhancementArgs) => {
		args.formData.append('otp', this.otp);
		const loadID = CreateToast('loading', 'logging in....');
		return async ({ result }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'success') {
				await goto(resolve('/(app)/(public)/home'), { replaceState: true });
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
