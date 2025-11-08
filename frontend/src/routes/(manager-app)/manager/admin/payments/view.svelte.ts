import type { PaymentMethods } from './model';
import type { EnhancementArgs, EnhancementReturn, SeekPaginatedResponse } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';
import { debounce } from '$lib/utils/helper';

export class PaymentManagementView {
	last_id = $state<number>(0);
	payments = $state<PaymentMethods[]>([]);
	page = $state<number>(1);
	totalRow = $state<number>(0);
	search = $state<string>('');
	isLoading = $state<boolean>(false);
	selectedDeletePayment = $state<number>(1);
	deleteDialogOpen = $state<boolean>(false);
	createDialogOpen = $state<boolean>(false);
	paymentToDelete = $state<number>();
	limit = $state<number>(0);
	createPaymentForm = $state<HTMLFormElement>();
	searchForm = $state<HTMLFormElement>();

	#SearchSubmit = debounce(() => {
		this.searchForm?.requestSubmit();
	}, 500);

	constructor(p: SeekPaginatedResponse<PaymentMethods>) {
		this.setLastID(p.page_info.last_id);
		this.setPayments(p.entries);
		this.setTotalRow(p.page_info.total_row);
		this.setLimit(p.page_info.limit);
	}
	onSearchInput = () => {
		this.#SearchSubmit();
	};
	setTotalRow = (n: number) => {
		this.totalRow = n;
	};
	setLimit = (n: number) => {
		this.limit = n;
	};
	setLastID = (n: number) => {
		this.last_id = n;
	};
	setPayments = (p: PaymentMethods[]) => {
		this.payments = p;
	};
	setPaymentToDelete = (id: number | undefined) => {
		this.paymentToDelete = id;
	};
	setIsLoading = (b: boolean) => {
		this.isLoading = b;
	};
	filterPayments(id: number) {
		this.payments = this.payments.filter((m) => m.payment_method_id !== id);
	}
	onPageChange(num: number) {
		this.page = num;
	}
	onCreatePayment = () => {
		const loadID = CreateToast('loading', 'Creating Payment Method.....');
		return async ({ result, update }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'success') {
				CreateToast('success', 'successfully create payment method');
			}
			update();
			this.createDialogOpen = false;
		};
	};
	onDeletePayment = (args: EnhancementArgs) => {
		this.setIsLoading(true);
		if (this.paymentToDelete) {
			args.formData.append('id', `${this.paymentToDelete}`);
		}
		return async ({ result }: EnhancementReturn) => {
			this.setIsLoading(false);
			if (result.type === 'success') {
				if (this.paymentToDelete) {
					this.filterPayments(this.paymentToDelete);
				}
				this.setPaymentToDelete(undefined);
				CreateToast('success', 'successfully delete payment');
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			this.deleteDialogOpen = false;
		};
	};
	setPaginationData(page: number, limit: number, total_row: number) {
		this.page = page;
		this.limit = limit;
		this.totalRow = total_row;
	}
	onSearchPayments = (args: EnhancementArgs) => {
		this.setIsLoading(true);
		args.formData.append('page', `${this.page}`);
		return async ({ result }: EnhancementReturn) => {
			this.setIsLoading(false);
			if (result.type === 'success') {
				this.setPayments(result.data?.payments.entries);
				this.setPaginationData(
					result.data?.payments.page_info.page,
					result.data?.payments.page_info.limit,
					result.data?.payments.page_info.total_row
				);
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
