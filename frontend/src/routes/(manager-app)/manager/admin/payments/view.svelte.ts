import type { PaymentMethods } from './model';

export class PaymentManagementView {
	last_id = $state<number>(0);
	payments = $state<PaymentMethods[]>([]);
	page = $state<number>(1);
	totalRow = $state<number>(0);
	search = $state<string>('');
	isLoading = $state<boolean>(false);
	selectedDeletePayment = $state<number>(1);

	setLastID = (n: number) => {
		this.last_id = n;
	};
	setPayments = (p: PaymentMethods[]) => {
		this.payments = p;
	};
	setTotalRow = (n: number) => {
		this.totalRow = n;
	};
	setIsLoading = (b: boolean) => {
		this.isLoading = b;
	};
}
