import type { EnhancementReturn } from '$lib/types';
import { CreateToast } from '$lib/utils/helper';
import { SvelteDate } from 'svelte/reactivity';

export class RequestDetailView {
	confirmDialogOpen = $state<boolean>(false);
	acceptDialogOpen = $state<boolean>(false);
	rejectDialogOpen = $state<boolean>(false);
	paymentDetailDialogOpen = $state<boolean>(false);
	onReject = () => {
		return async ({ result }: EnhancementReturn) => {
			this.rejectDialogOpen = false;
			if (result.type === 'success') {
				CreateToast('success', 'successfully reject request');
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	onAccept = () => {
		return async ({ result }: EnhancementReturn) => {
			this.acceptDialogOpen = false;
			if (result.type === 'success') {
				CreateToast('success', 'successfully accept request');
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	onConfirm = () => {
		return async ({ result }: EnhancementReturn) => {
			this.confirmDialogOpen = false;
			if (result.type === 'success') {
				CreateToast('success', 'successfully confirm payment');
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	convertToDate = (datetime: string) => {
		const date = new SvelteDate(datetime);
		return `${date.getFullYear()}-${date.getMonth()}-${date.getDate()}`;
	};
}
