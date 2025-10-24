import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
import { CreateToast, debounce, DismissToast } from '$lib/utils/helper';
import { dayofWeeks } from './constants';
import {
	type mentorPaymentMethods,
	type MentorSchedule,
	type paymentMethod,
	type paymentMethodOpts,
	type TimeOnly
} from './model';

export class CreateMentorView {
	degree = $state<string>('');
	resumeFile = $state<FileList>();

	generatePasswordForm = $state<HTMLFormElement>();
	generatedPassword = $state<string>('');

	paymentMethods = $state<paymentMethodOpts[]>([]);
	selectedPaymentMethod = $state<string>('');
	mentorPaymentMethods = $state<mentorPaymentMethods[]>([]);
	accountNumber = $state<string>('');
	paymentMethodForm = $state<HTMLFormElement>();
	searchValue = $state<string>('');
	disableAddPaymentMethod = $derived.by<boolean>(() => {
		if (!this.selectedPaymentMethod || !this.accountNumber) {
			return true;
		}
		return false;
	});
	selectedPaymentLabel = $derived.by<string | undefined>(() => {
		if (this.selectedPaymentMethod) {
			const filtered = this.paymentMethods.filter(
				(val) => val.value === this.selectedPaymentMethod
			);
			return filtered[0].label;
		}
	});
	#paymentMethodSubmit = debounce(() => {
		this.paymentMethodForm?.requestSubmit();
	}, 500);

	mentorSchedules = $state<MentorSchedule[]>([]);
	selectedDayOfWeek = $state<string>('');
	selectedStartTime = $state<string>('');
	selectedEndTime = $state<string>('');
	disableAddMentorSchedule = $derived.by<boolean>(() => {
		if (!this.selectedDayOfWeek || !this.selectedStartTime || !this.selectedEndTime) {
			return true;
		}
		return false;
	});
	selectedDayOfWeekLabel = $derived.by<string | undefined>(() => {
		if (this.selectedDayOfWeek) {
			const filtered = dayofWeeks.filter((val) => val.value === this.selectedDayOfWeek);
			return filtered[0].label;
		}
	});

	generatePassword = () => {
		this.generatePasswordForm?.requestSubmit();
	};
	setPaymentMethods(newPayments: paymentMethod[]) {
		const opts: paymentMethodOpts[] = [];
		newPayments.forEach((val) => {
			opts.push({
				value: `${val.payment_method_id}`,
				label: val.payment_method_name
			});
		});
		this.paymentMethods = opts;
	}
	addMentorPaymentMethod = () => {
		this.mentorPaymentMethods.push({
			payment_method_id: parseInt(this.selectedPaymentMethod),
			payment_method_name: this.selectedPaymentLabel!,
			account_number: this.accountNumber
		});
	};
	removeMentorPaymentMethod = (i: number) => {
		this.mentorPaymentMethods = this.mentorPaymentMethods.filter((v, idx) => idx !== i);
	};
	addMentorSchedule = () => {
		this.mentorSchedules.push({
			day_of_week: parseInt(this.selectedDayOfWeek),
			start_time: this.#stringToTimeOnly(this.selectedStartTime),
			end_time: this.#stringToTimeOnly(this.selectedEndTime),
			day_of_week_label: this.selectedDayOfWeekLabel!
		});
	};
	removeMentorSchedule = (i: number) => {
		this.mentorSchedules = this.mentorSchedules.filter((v, idx) => idx !== i);
	};
	onGetPaymentMethods = (args: EnhancementArgs) => {
		args.formData.append('search', this.searchValue);
		return async ({ result }: EnhancementReturn) => {
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'success') {
				this.setPaymentMethods(result.data?.paymentMethods);
			}
		};
	};
	onGetPassword = () => {
		return async ({ result }: EnhancementReturn) => {
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'success') {
				this.generatedPassword = result.data?.password;
			}
		};
	};

	onKeyWordChange = (e: Event & { currentTarget: EventTarget & HTMLInputElement }) => {
		this.searchValue = e.currentTarget.value;
		this.#paymentMethodSubmit();
	};

	setGeneratedPassword(p: string) {
		this.generatedPassword = p;
	}
	onCreateMentor = (args: EnhancementArgs) => {
		const loadID = CreateToast('loading', 'creating....');
		this.mentorPaymentMethods.forEach((val) => {
			args.formData.append('mentor_payment_info', JSON.stringify(val));
		});
		this.mentorSchedules.forEach((val) => {
			args.formData.append('mentor_availability', JSON.stringify(val));
		});
		args.formData.append('password', this.generatedPassword);
		return async ({ result }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'success') {
				CreateToast('success', result.data?.message);
			}
		};
	};
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
}
