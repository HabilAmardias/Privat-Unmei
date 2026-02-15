import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';

export class LoginCallbackView {
	countdown = $state<number>(30);
	otp = $state<string>('');
	loginDisabled = $derived<boolean>(this.otp.length < 6);
	resendOTPDisabled = $derived<boolean>(this.countdown > 0);
	countdownDisplay = $derived.by<string>(() => {
		if (this.countdown <= 0) {
			return '';
		}
		return `${this.countdown}s`;
	});
	interval = $state<NodeJS.Timeout>();

	startTimer = () => {
		this.countdown = 30;
		this.interval = setInterval(() => {
			this.countdown -= 1;
			if (this.countdown <= 0) {
				clearInterval(this.interval);
				this.interval = undefined;
			}
		}, 1000);
	};

	onLoginSubmit = (args: EnhancementArgs) => {
		args.formData.append('otp', this.otp);
		const loadID = CreateToast('loading', 'loading....');
		return async ({ result }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'success') {
				if (args.action.search === '?/resendOTP') {
					CreateToast('success', 'success');
					this.startTimer();
					return;
				}
				CreateToast('success', 'login success');
				await goto(resolve('/(app)/(public)/home'), { replaceState: true });
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
