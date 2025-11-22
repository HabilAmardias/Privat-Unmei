import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
import { CreateToast, debounce, DismissToast } from '$lib/utils/helper';
import { FILE_IMAGE_THRESHOLD, MAX_BIO_LENGTH } from '$lib/utils/constants';
import { dayofWeeks } from './constants';
import {
	type MentorPaymentInfo,
	type MentorProfile,
	type MentorSchedule,
	type MentorScheduleInfo,
	type paymentMethod,
	type paymentMethodOpts,
	type TimeOnly
} from './model';

export class UpdateMentorProfileView {
	isDesktop = $state<boolean>(false);
	size = $derived.by<number>(() => {
		if (this.isDesktop) {
			return 150;
		}
		return 100;
	});
	profileImage = $state<FileList>();
	name = $state<string | undefined>();
	bio = $state<string | undefined>();
	bioErr = $derived.by<Error | undefined>(() => {
		if (this.bio && this.bio.length > MAX_BIO_LENGTH) {
			return new Error(`Bio length must less than ${MAX_BIO_LENGTH} characters`);
		}
		return undefined;
	});
	degree = $state<string | undefined>();
	resumeFile = $state<FileList | undefined>();
	resumeErr = $derived.by<Error | undefined>(() => {
		if (this.resumeFile && this.resumeFile.length > 0) {
			const file = this.resumeFile[0];
			if (file.type !== 'application/pdf' || file.size > FILE_IMAGE_THRESHOLD) {
				return new Error('invalid file');
			}
			return undefined;
		}
		return undefined;
	});
	yearsOfExperience = $state<number | undefined>(undefined);
	yoeErr = $derived.by<Error | undefined>(() => {
		if (this.yearsOfExperience && (this.yearsOfExperience < 0 || this.yearsOfExperience > 40)) {
			return new Error('invalid years of experience');
		}
		return undefined;
	});
	campus = $state<string | undefined>();
	major = $state<string | undefined>();

	paymentMethods = $state<paymentMethodOpts[]>([]);
	selectedPaymentMethod = $state<string>('');
	mentorPaymentMethods = $state<MentorPaymentInfo[]>([]);
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

	disableUpdateMentor = $derived.by<boolean>(() => {
		if (
			!this.name ||
			!this.bio ||
			this.yearsOfExperience === undefined ||
			!this.campus ||
			!this.degree ||
			!this.major ||
			this.yoeErr ||
			this.resumeErr
		) {
			return true;
		}
		return false;
	});

	constructor(
		p: paymentMethod[],
		sch: MentorScheduleInfo[],
		pym: MentorPaymentInfo[],
		profile: MentorProfile
	) {
		this.setPaymentMethods(p);
		this.mentorPaymentMethods = pym;
		this.setMentorSchedules(sch);
		this.name = profile.name;
		this.bio = profile.bio;
		this.yearsOfExperience = profile.years_of_experience;
		this.campus = profile.campus;
		this.major = profile.major;
		this.degree = profile.degree;
	}

	setMentorSchedules(newSchedules: MentorScheduleInfo[]) {
		const schedules: MentorSchedule[] = [];
		newSchedules.forEach((v) => {
			const label = dayofWeeks.filter((d) => {
				return parseInt(d.value) === v.day_of_week;
			});
			schedules.push({
				start_time: this.#stringToTimeOnly(v.start_time),
				end_time: this.#stringToTimeOnly(v.end_time),
				day_of_week: v.day_of_week,
				day_of_week_label: label[0].label
			});
		});
		this.mentorSchedules = schedules;
	}
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
	onSearchPaymentMethodChange = (e: Event & { currentTarget: EventTarget & HTMLInputElement }) => {
		this.searchValue = e.currentTarget.value;
		this.#paymentMethodSubmit();
	};
	onUpdateMentor = (args: EnhancementArgs) => {
		const loadID = CreateToast('loading', 'updating....');
		this.mentorPaymentMethods.forEach((val) => {
			args.formData.append('mentor_payment_info', JSON.stringify(val));
		});
		this.mentorSchedules.forEach((val) => {
			args.formData.append('mentor_availability', JSON.stringify(val));
		});
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
