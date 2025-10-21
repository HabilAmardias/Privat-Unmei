import type { mentorPaymentMethods, paymentMethod, paymentMethodOpts } from './model';

export class CreateMentorView {
	degree = $state<string>('');
	paymentMethods = $state<paymentMethodOpts[]>([]);
	selectedPaymentMethod = $state<string>('');
	mentorPaymentMethods = $state<mentorPaymentMethods[]>([]);
	accountNumber = $state<string>('');

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
}
