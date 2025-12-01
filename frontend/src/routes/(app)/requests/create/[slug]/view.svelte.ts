import type { CourseDetail, MentorPaymentInfo, ScheduleSlot } from './model';

export class CreateRequestView {
	maxSession = $state<number>(0);
	paymentOpts = $state<{ value: string; label: string }[]>([]);
	schedules = $state<ScheduleSlot[]>([]);
	selectedDate = $state<string>('');
	selectedStartTime = $state<string>('');

	constructor(d: CourseDetail, p: MentorPaymentInfo[]) {
		this.maxSession = d.max_total_session;
		p.forEach((e) => {
			this.paymentOpts.push({
				label: e.payment_method_name,
				value: `${e.payment_method_id}`
			});
		});
	}
	capitalizeFirstLetter(s: string) {
		if (s.length === 0) {
			return s;
		}
		if (s.length === 1) {
			return s.toUpperCase();
		}
		return s.charAt(0).toUpperCase() + s.slice(1);
	}
}
