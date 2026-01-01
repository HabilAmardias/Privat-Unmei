import { goto } from '$app/navigation';
import { resolve } from '$app/paths';
import type { EnhancementReturn } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';

export class GoogleCallbackView {
	openDialog = $state<boolean>(false);
	agreed = $state<boolean>(false);
	verifyForm = $state<HTMLFormElement>();

	verifiedLogin = () => {
		CreateToast('success', 'login success');
		goto(resolve('/(app)/(public)/home'), { replaceState: true });
	};

	onVerify = () => {
		const loadID = CreateToast('loading', 'loading....');
		return async ({ result }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'success') {
				this.agreed = true;
				this.openDialog = false;
				CreateToast('success', 'login success');
				goto(resolve('/(app)/(public)/home'), { replaceState: true });
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
