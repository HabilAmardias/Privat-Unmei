import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
import { CreateToast } from '$lib/utils/helper';
import { SvelteDate } from 'svelte/reactivity';
import { DismissToast } from '$lib/utils/helper';
import type { RequestDetail } from './model';

export class RequestDetailView {
	confirmDialogOpen = $state<boolean>(false);
	acceptDialogOpen = $state<boolean>(false);
	rejectDialogOpen = $state<boolean>(false);
	paymentDetailDialogOpen = $state<boolean>(false);
	mentorID = $state<string>('');
	constructor(d: RequestDetail) {
		this.mentorID = d.mentor_id;
	}
	onReject = () => {
		return async ({ result, update }: EnhancementReturn) => {
			this.rejectDialogOpen = false;
			if (result.type === 'redirect') {
				CreateToast('success', 'successfully reject request');
				await update();
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	onAccept = () => {
		return async ({ result, update }: EnhancementReturn) => {
			this.acceptDialogOpen = false;
			if (result.type === 'redirect') {
				CreateToast('success', 'successfully accept request');
				await update();
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	onConfirm = () => {
		return async ({ result, update }: EnhancementReturn) => {
			this.confirmDialogOpen = false;
			if (result.type === 'redirect') {
				CreateToast('success', 'successfully confirm payment');
				await update();
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
	convertToDatetime = (input: string) => {
		const [date, tz] = input.split('T');
		const time = tz.split('.')[0];
		return `${date} ${time}`;
	};
	capitalizeFirstLetter(s: string) {
		if (s.length === 0) {
			return s;
		}
		if (s.length === 1) {
			return s.toUpperCase();
		}
		return s.charAt(0).toUpperCase() + s.slice(1);
	}
	onMessageMentor = (args: EnhancementArgs) => {
		const loadID = CreateToast('loading', 'loading....');
		args.formData.append('id', this.mentorID);
		return async ({ result, update }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'redirect') {
				await update();
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
