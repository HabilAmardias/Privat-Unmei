import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
import { CreateToast } from '$lib/utils/helper';
import { SvelteDate } from 'svelte/reactivity';
import { DismissToast } from '$lib/utils/helper';
import type { RequestDetail } from './model';
import { Banknote, BookMarked, CalendarCheck, Check, X } from '@lucide/svelte';
import { goto } from '$app/navigation';
import { resolve } from '$app/paths';

export class RequestDetailView {
	mentorID = $state<string>('');
	detailState = $state<'detail' | 'payment'>('detail');
	status = $state<string>('');
	now = $state<SvelteDate>(new SvelteDate());
	expiredAt = $state<string | null>('');
	feedback = $state<string>('');
	feedbackErr = $derived.by<Error | undefined>(() => {
		if (this.feedback.length < 15) {
			return new Error('Feedback must at least contain 15 characters');
		}
		return undefined;
	});
	reviewDisabled = $derived.by<boolean>(() => {
		return this.feedbackErr ? true : false;
	});
	expiredIn = $derived.by<string | null>(() => {
		if (this.expiredAt) {
			return this.getTimeDifference(this.expiredAt, this.now);
		}
		return null;
	});
	star = $state<number>(1);
	courseID = $state<number>(0);

	statuses = [
		{ id: 'reserved', label: 'Reserved', icon: BookMarked },
		{ id: 'pending payment', label: 'Pending Payment', icon: Banknote },
		{ id: 'scheduled', label: 'Scheduled', icon: CalendarCheck },
		{ id: 'completed', label: 'Completed', icon: Check },
		{ id: 'cancelled', label: 'Cancelled', icon: X }
	];

	constructor(d: RequestDetail) {
		this.mentorID = d.mentor_id;
		this.status = d.status;
		this.expiredAt = d.expired_at;
		this.courseID = d.course_id;
	}
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
	getTimeDifference(dbDate: string, now: SvelteDate) {
		const date = new SvelteDate(dbDate);
		const diffInMs = date.valueOf() - now.valueOf();

		const totalMinutes = Math.floor(diffInMs / (1000 * 60));
		const hours = Math.floor(totalMinutes / 60);
		const minutes = totalMinutes % 60;

		return `${hours}h ${minutes}m`;
	}
	onCreateReview = (args: EnhancementArgs) => {
		const loadID = CreateToast('loading', 'loading....');
		args.formData.append('id', `${this.courseID}`);
		return async ({ result }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'success') {
				CreateToast('success', 'successfully create review');
				goto(resolve('/(app)/profile'));
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
