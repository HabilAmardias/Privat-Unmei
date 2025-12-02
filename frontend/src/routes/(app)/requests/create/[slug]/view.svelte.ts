import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';
import type { DateValue } from '@internationalized/date';
import type { CourseDetail, MentorPaymentInfo, ScheduleSlot, TimeOnly } from './model';

export class CreateRequestView {
	maxSession = $state<number>(0);
	paymentOpts = $state<{ value: string; label: string }[]>([]);
	schedules = $state<ScheduleSlot[]>([]);
	selectedDate = $state<string>('');
	selectedStartTime = $state<string>('');
	selectedPayment = $state<string>('');
	participant = $state<number>(1);

	constructor(d: CourseDetail, p: MentorPaymentInfo[]) {
		this.maxSession = d.max_total_session;
		p.forEach((e) => {
			this.paymentOpts.push({
				label: e.payment_method_name,
				value: `${e.payment_method_id}`
			});
		});
	}

	onCalendarValueChange = (date?: DateValue) => {
		if (date) {
			this.selectedDate = date.toString();
		}
	};

	onCreateRequest = (args: EnhancementArgs) => {
		const loadID = CreateToast('loading', 'creating....');
		const scheduleInput = JSON.stringify(this.schedules);
		args.formData.append('schedules', scheduleInput);

		return async ({ result, update }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'redirect') {
				CreateToast('success', 'successfully create request');
				await update();
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
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

	addSchedule = () => {
		this.schedules.push({
			date: this.selectedDate,
			start_time: this.#stringToTimeOnly(this.selectedStartTime)
		});
		this.selectedDate = '';
		this.selectedStartTime = '';
	};
	disableAddSchedule = $derived.by<boolean>(() => {
		if (!this.selectedDate || !this.selectedStartTime) {
			return true;
		}
		return false;
	});
	disableSubmit = $derived.by<boolean>(() => {
		if (this.schedules.length === 0 || !this.participant || !this.selectedPayment) {
			return true;
		}
		return false;
	});
	#stringToTimeOnly(s: string): TimeOnly {
		const [hour, minute, second] = s.split(':');
		return {
			hour: parseInt(hour),
			minute: parseInt(minute),
			second: second ? parseInt(second) : 0
		};
	}
}
