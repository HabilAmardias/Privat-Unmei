import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
import { CreateToast, debounce, DismissToast } from '$lib/utils/helper';
import type {
	mentorPaymentMethods,
	MentorSchedule,
	paymentMethod,
	paymentMethodOpts
} from './model';

export class CreateMentorView {
	degree = $state<string>('');
	paymentMethods = $state<paymentMethodOpts[]>([]);
	selectedPaymentMethod = $state<string>('');
	mentorPaymentMethods = $state<mentorPaymentMethods[]>([]);
	accountNumber = $state<string>('');
	paymentMethodForm = $state<HTMLFormElement>();
	searchValue = $state<string>('');
	generatedPassword = $state<string>('');
	generatePasswordForm = $state<HTMLFormElement>();
	disableAddPaymentMethod = $derived.by<boolean>(() => {
		if (!this.selectedPaymentMethod || !this.accountNumber) {
			return true;
		}
		return false;
	});
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
			account_number: this.accountNumber
		});
	};
	addMentorSchedule = () => {
		this.mentorSchedules.push();
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

	#paymentMethodSubmit = debounce(() => {
		this.paymentMethodForm?.requestSubmit();
	}, 500);

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
}
