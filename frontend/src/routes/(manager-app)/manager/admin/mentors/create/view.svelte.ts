import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
import { CreateToast, debounce, DismissToast } from '$lib/utils/helper';
import { FILE_IMAGE_THRESHOLD } from '$lib/utils/constants';
import { dayofWeeks } from './constants';
import {
	type mentorPaymentMethods,
	type MentorSchedule,
	type paymentMethod,
	type paymentMethodOpts,
	type TimeOnly
} from './model';

export class CreateMentorView {
	name = $state<string>('');
	degree = $state<string>('');
	resumeFile = $state<FileList>();
	resumeErr = $derived.by<Error | undefined>(() => {
		if (this.resumeFile && this.resumeFile.length > 0) {
			const file = this.resumeFile[0];
			if (file.type !== 'application/pdf' || file.size > FILE_IMAGE_THRESHOLD) {
				return new Error('invalid file');
			}
		}
		return undefined;
	});
	email = $state<string>('');
	emailErr = $derived.by<Error | undefined>(() => {
		if (this.email.length > 0 && !this.#validateEmail(this.email)) {
			return new Error('must insert a valid email');
		}
		return undefined;
	});
	yearsOfExperience = $state<number>(0);
	yoeErr = $derived.by<Error | undefined>(() => {
		if (this.yearsOfExperience < 0 || this.yearsOfExperience > 40) {
			return new Error('invalid years of experience');
		}
		return undefined;
	});
	campus = $state<string>('');
	major = $state<string>('');

	generatePasswordForm = $state<HTMLFormElement>();
	generatedPassword = $state<string>('');

	paymentMethods = $state<paymentMethodOpts[]>([]);
	selectedPaymentMethod = $state<string>('');
	mentorPaymentMethods = $state<mentorPaymentMethods[]>([]);
	accountNumber = $state<string>('');
	paymentMethodForm = $state<HTMLFormElement>();
	searchValue = $state<string>('');
	selectPaymentMethodErr = $derived.by<Error | null>(() => {
		const filtered = this.mentorPaymentMethods.filter(
			(val) => val.payment_method_id === parseInt(this.selectedPaymentMethod)
		);
		if (filtered.length > 0) {
			return new Error('cannot have same payment method');
		}
		return null;
	});
	disableAddPaymentMethod = $derived.by<boolean>(() => {
		if (!this.selectedPaymentMethod || !this.accountNumber || this.selectPaymentMethodErr) {
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
	selectedStartTimeInSec = $derived.by<number>(() => {
		if (this.selectedStartTime) {
			const parsedStartTime = this.#stringToTimeOnly(this.selectedStartTime);
			return parsedStartTime.hour * 3600 + parsedStartTime.minute * 60 + parsedStartTime.second;
		}
		return 0;
	});
	selectedEndTimeInSec = $derived.by<number>(() => {
		if (this.selectedEndTime) {
			const parsedEndTime = this.#stringToTimeOnly(this.selectedEndTime);
			return parsedEndTime.hour * 3600 + parsedEndTime.minute * 60 + parsedEndTime.second;
		}
		return 0;
	});
	selectMentorScheduleErr = $derived.by<Error | null>(() => {
		const filtered = this.mentorSchedules.filter(
			(val) => val.day_of_week === parseInt(this.selectedDayOfWeek)
		);
		if (filtered.length > 0) {
			return new Error('cannot add same day of week');
		}
		if (this.selectedStartTimeInSec >= this.selectedEndTimeInSec) {
			return new Error('invalid time input');
		}
		return null;
	});
	disableAddMentorSchedule = $derived.by<boolean>(() => {
		if (
			!this.selectedDayOfWeek ||
			!this.selectedStartTime ||
			!this.selectedEndTime ||
			this.selectMentorScheduleErr
		) {
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

	disableCreateMentor = $derived.by<boolean>(() => {
		if (
			!this.name ||
			!this.email ||
			this.emailErr ||
			this.yoeErr ||
			!this.generatedPassword ||
			!this.campus ||
			!this.degree ||
			!this.major ||
			!this.resumeFile ||
			this.resumeFile.length === 0 ||
			this.resumeErr ||
			this.mentorPaymentMethods.length === 0 ||
			this.mentorSchedules.length === 0
		) {
			return true;
		}
		return false;
	});

	#validateEmail(email: string) {
		const pattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		return pattern.test(email);
	}
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
		this.selectedPaymentMethod = '';
		this.selectedPaymentLabel = '';
		this.accountNumber = '';
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
		this.selectedDayOfWeek = '';
		this.selectedStartTime = '';
		this.selectedEndTime = '';
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
	onSearchPaymentMethodChange = (e: Event & { currentTarget: EventTarget & HTMLInputElement }) => {
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
