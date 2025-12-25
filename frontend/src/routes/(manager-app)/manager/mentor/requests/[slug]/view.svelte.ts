import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';
import { SvelteDate } from 'svelte/reactivity';
import type { RequestDetail } from './model';
import { Banknote, BookMarked, CalendarCheck, Check, X } from '@lucide/svelte';

export class RequestDetailView {
	confirmDialogOpen = $state<boolean>(false);
	acceptDialogOpen = $state<boolean>(false);
	rejectDialogOpen = $state<boolean>(false);
	paymentDetailDialogOpen = $state<boolean>(false);
	detailState = $state<'detail' | 'payment'>('detail');
	status = $state<string>('');
	now = $state<SvelteDate>(new SvelteDate());
	expiredAt = $state<string | null>('');
	expiredIn = $derived.by<string | null>(() => {
		if (this.expiredAt) {
			return this.getTimeDifference(this.expiredAt, this.now);
		}
		return null;
	});
	studentID = $state<string>('');

	statuses = [
		{ id: 'reserved', label: 'Reserved', icon: BookMarked },
		{ id: 'pending', label: 'Pending Payment', icon: Banknote },
		{ id: 'scheduled', label: 'Scheduled', icon: CalendarCheck },
		{ id: 'completed', label: 'Completed', icon: Check },
		{ id: 'cancelled', label: 'Cancelled', icon: X }
	];
	constructor(d: RequestDetail) {
		this.studentID = d.student_id;
		this.status = d.status;
		this.expiredAt = d.expired_at;
	}
	getTimeDifference(dbDate: string, now: SvelteDate) {
		const date = new SvelteDate(dbDate);
		const diffInMs = date.valueOf() - now.valueOf();

		const totalMinutes = Math.floor(diffInMs / (1000 * 60));
		const hours = Math.floor(totalMinutes / 60);
		const minutes = totalMinutes % 60;

		return `${hours}h ${minutes}m`;
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
		return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`;
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
	onMessageStudent = (args: EnhancementArgs) => {
		const loadID = CreateToast('loading', 'loading....');
		args.formData.append('id', this.studentID);
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
