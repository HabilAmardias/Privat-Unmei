import type { PaymentMethods } from './model';
import type { EnhancementArgs, EnhancementReturn, PaginatedResponse } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';
import { debounce } from '$lib/utils/helper';
import type { paymentMethod } from '../mentors/create/model';

export class PaymentManagementView {
	payments = $state<PaymentMethods[]>([]);
	pageNumber = $state<number>(1);
	totalRow = $state<number>(0);
	search = $state<string>('');
	isLoading = $state<boolean>(false);
	selectedDeletePayment = $state<number>(1);
	deleteDialogOpen = $state<boolean[]>([]);
	updateDialogOpen = $state<boolean[]>([]);
	createDialogOpen = $state<boolean>(false);
	paymentToDelete = $state<number>();
	limit = $state<number>(0);
	createPaymentForm = $state<HTMLFormElement>();
	searchForm = $state<HTMLFormElement>();
	paginationForm = $state<HTMLFormElement>();
	paymentToUpdate = $state<number>();

	#SearchSubmit = debounce(() => {
		this.searchForm?.requestSubmit();
	}, 500);

	constructor(p: PaginatedResponse<PaymentMethods>) {
		this.setPageNumber(p.page_info.page);
		this.setPayments(p.entries);
		this.deleteDialogOpen = new Array<boolean>(this.payments.length).fill(false);
		this.updateDialogOpen = new Array<boolean>(this.payments.length).fill(false);
		this.setTotalRow(p.page_info.total_row);
		this.setLimit(p.page_info.limit);
	}
	setPageNumber(n: number) {
		this.pageNumber = n;
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
	setPayments = (p: PaymentMethods[]) => {
		this.payments = p;
	};
	setPaymentToDelete = (id: number | undefined) => {
		this.paymentToDelete = id;
	};
	setPaymentToUpdate = (id: number | undefined) => {
		this.paymentToUpdate = id;
	};
	setIsLoading = (b: boolean) => {
		this.isLoading = b;
	};
	filterPayments(id: number) {
		this.payments = this.payments.filter((m) => m.payment_method_id !== id);
	}
	onPageChange = () => {
		this.paginationForm?.requestSubmit();
	};
	onCreatePayment = (args: EnhancementArgs) => {
		const loadID = CreateToast('loading', 'Creating Payment Method.....');
		return async ({ result }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'success') {
				const newPayment: paymentMethod = {
					payment_method_id: result.data?.newPayment.payment_method_id,
					payment_method_name: args.formData.get('name') as string
				};
				if (this.payments.length < this.limit) {
					this.payments.push(newPayment);
					this.deleteDialogOpen = new Array<boolean>(this.payments.length).fill(false);
					this.updateDialogOpen = new Array<boolean>(this.payments.length).fill(false);
				}
				this.totalRow += 1;
				CreateToast('success', 'successfully create payment method');
			}
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
					this.deleteDialogOpen = new Array<boolean>(this.payments.length).fill(false);
					this.updateDialogOpen = new Array<boolean>(this.payments.length).fill(false);
				}
				this.totalRow -= 1;
				this.setPaymentToDelete(undefined);
				CreateToast('success', 'successfully delete payment');
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	onUpdatePayment = (args: EnhancementArgs) => {
		this.setIsLoading(true);
		if (this.paymentToUpdate) {
			args.formData.append('id', `${this.paymentToUpdate}`);
		}
		return async ({ result }: EnhancementReturn) => {
			this.setIsLoading(false);
			if (result.type === 'success') {
				if (this.paymentToUpdate) {
					this.payments = this.payments.map((m) => {
						if (m.payment_method_id === this.paymentToUpdate) {
							return {
								payment_method_id: this.paymentToUpdate,
								payment_method_name: args.formData.get('name') as string
							};
						}
						return m;
					});
				}
				this.setPaymentToUpdate(undefined);
				CreateToast('success', 'successfully update payment');
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	setPaginationData(limit: number, total_row: number, page: number) {
		this.limit = limit;
		this.totalRow = total_row;
		this.pageNumber = page;
	}
	onSearchPayments = (args: EnhancementArgs) => {
		this.setIsLoading(true);
		this.pageNumber = 1;
		args.formData.append('page', `${this.pageNumber}`);
		return async ({ result }: EnhancementReturn) => {
			this.setIsLoading(false);
			if (result.type === 'success') {
				this.setPayments(result.data?.payments.entries);
				this.deleteDialogOpen = new Array<boolean>(this.payments.length).fill(false);
				this.updateDialogOpen = new Array<boolean>(this.payments.length).fill(false);
				this.setPaginationData(
					result.data?.payments.page_info.limit,
					result.data?.payments.page_info.total_row,
					result.data?.payments.page_info.page
				);
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	onPageChangeForm = (args: EnhancementArgs) => {
		this.setIsLoading(true);
		args.formData.append('page', `${this.pageNumber}`);
		args.formData.append('search', this.search);
		return async ({ result }: EnhancementReturn) => {
			this.setIsLoading(false);
			if (result.type === 'success') {
				this.setPayments(result.data?.payments.entries);
				this.deleteDialogOpen = new Array<boolean>(this.payments.length).fill(false);
				this.updateDialogOpen = new Array<boolean>(this.payments.length).fill(false);
				this.setPaginationData(
					result.data?.payments.page_info.limit,
					result.data?.payments.page_info.total_row,
					result.data?.payments.page_info.page
				);
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
