import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
import { CreateToast, debounce } from '$lib/utils/helper';
import { PaymentMethodOptionsLimit } from './constants';
import type { mentorPaymentMethods, paymentMethod, paymentMethodOpts } from './model';

export class CreateMentorView {
	degree = $state<string>('');
	paymentMethods = $state<paymentMethodOpts[]>([]);
	selectedPaymentMethod = $state<string>('');
	mentorPaymentMethods = $state<mentorPaymentMethods[]>([]);
	accountNumber = $state<string>('');
	paymentMethodForm = $state<HTMLFormElement>();
	searchValue = $state<string>('');

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
	onGetPaymentMethods = (args: EnhancementArgs) => {
		args.formData.append('search', this.searchValue);
		args.formData.append('limit', PaymentMethodOptionsLimit);
		return async ({ result }: EnhancementReturn) => {
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'success') {
				this.setPaymentMethods(result.data?.paymentMethods);
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
}
