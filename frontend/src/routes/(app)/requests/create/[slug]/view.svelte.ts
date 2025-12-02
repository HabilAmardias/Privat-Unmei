import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';
import type { DateValue } from '@internationalized/date';
import type { CourseDetail, MentorPaymentInfo, ScheduleSlot, TimeOnly } from './model';
import { debounce } from '$lib/utils/helper';

export class CreateRequestView {
	maxSession = $state<number>(0);
	paymentOpts = $state<{ value: string; label: string }[]>([]);
	schedules = $state<ScheduleSlot[]>([]);
	selectedDate = $state<string>('');
	selectedStartTime = $state<string>('');
	selectedPayment = $state<string>('');
	participant = $state<number>(1);
	participantErr = $derived.by<Error | undefined>(() => {
		return this.participant <= 0 ? new Error('need at least one participant') : undefined;
	});
	price = $state<number>(0);
	cost = $state<number>(0);
	subtotal = $derived<number>(this.price * this.schedules.length);
	operational = $derived<number>(this.cost * this.schedules.length);
	discount = $state<number>(0);
	total = $derived<number>(this.subtotal + this.operational - this.discount);

	getDiscountForm = $state<HTMLFormElement>();
	#getDiscountSubmit = debounce(() => {
		this.getDiscountForm?.requestSubmit();
	}, 500);

	constructor(d: CourseDetail, p: MentorPaymentInfo[], cost: number, discount: number) {
		this.maxSession = d.max_total_session;
		this.price = d.price;
		this.discount = discount;

		p.forEach((e) => {
			this.paymentOpts.push({
				label: e.payment_method_name,
				value: `${e.payment_method_id}`
			});
		});
		this.cost = cost;
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
		this.selectedStartTime = '';
	};
	removeSchedule = (date: string) => {
		this.schedules = this.schedules.filter((e) => {
			return e.date !== date;
		});
	};
	disableAddSchedule = $derived.by<boolean>(() => {
		if (!this.selectedDate || !this.selectedStartTime || this.dateErr) {
			return true;
		}
		return false;
	});
	dateErr = $derived.by<Error | undefined>(() => {
		const numDupes = this.schedules.filter((e) => {
			return e.date === this.selectedDate;
		}).length;
		if (numDupes >= 1) {
			return new Error('cannot add same date');
		}
		return undefined;
	});
	disableSubmit = $derived.by<boolean>(() => {
		if (
			this.schedules.length === 0 ||
			!this.participant ||
			!this.selectedPayment ||
			this.participantErr ||
			this.total < 0
		) {
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
	TimeOnlyToString(t: TimeOnly): string {
		let res: string = '';
		res += t.hour < 10 ? `0${t.hour}:` : `${t.hour}:`;
		res += t.minute < 10 ? `0${t.minute}:` : `${t.minute}:`;
		res += t.second < 10 ? `0${t.second}` : `${t.second}`;
		return res;
	}
	onGetDiscount = (args: EnhancementArgs) => {
		args.formData.append('participant', `${this.participant}`);
		return async ({ result }: EnhancementReturn) => {
			if (result.type === 'success') {
				this.discount = result.data?.discount;
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	onParticipantChange = () => {
		if (this.participant > 0) {
			this.#getDiscountSubmit();
		}
	};
}
