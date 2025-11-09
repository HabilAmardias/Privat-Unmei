import type { AdditionalCost, Discount } from './model';
import type { EnhancementArgs, EnhancementReturn, PaginatedResponse } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';

export class PaymentManagementView {
	menuState = $state<'costs' | 'discounts'>('costs');

	costs = $state<AdditionalCost[]>([]);
	costPageNumber = $state<number>(1);
	costTotalRow = $state<number>(0);
	costLimit = $state<number>(0);
	costToDelete = $state<number>();
	deleteCostDialogOpen = $state<boolean>(false);
	updateCostDialogOpen = $state<boolean>(false);
	createCostDialogOpen = $state<boolean>(false);
	costPaginationForm = $state<HTMLFormElement>();
	costToUpdate = $state<number>();
	createCostForm = $state<HTMLFormElement>();

	discounts = $state<Discount[]>([]);
	discountPageNumber = $state<number>(1);
	discountTotalRow = $state<number>(0);
	discountLimit = $state<number>(0);
	discountToDelete = $state<number>(1);
	deleteDiscountDialogOpen = $state<boolean>(false);
	updateDiscountDialogOpen = $state<boolean>(false);
	createDiscountDialogOpen = $state<boolean>(false);
	discountPaginationForm = $state<HTMLFormElement>();
	discountToUpdate = $state<number>();
	createDiscountForm = $state<HTMLFormElement>();

	isLoading = $state<boolean>(false);

	constructor(c: PaginatedResponse<AdditionalCost>, d: PaginatedResponse<Discount>) {
		this.costLimit = c.page_info.limit;
		this.costPageNumber = c.page_info.page;
		this.costTotalRow = c.page_info.total_row;
		this.costs = c.entries;

		this.discountLimit = d.page_info.limit;
		this.discountPageNumber = d.page_info.page;
		this.discountTotalRow = d.page_info.total_row;
		this.discounts = d.entries;
	}

	setPageChange = (num: number) => {
		if (this.menuState === 'costs') {
			this.costPageNumber = num;
			this.costPaginationForm?.requestSubmit();
		} else {
			this.discountPageNumber = num;
			this.discountPaginationForm?.requestSubmit();
		}
	};
	onCreateCost = (args: EnhancementArgs) => {
		const loadID = CreateToast('loading', 'Creating Cost.....');
		return async ({ result }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'success') {
				const newCost: AdditionalCost = {
					id: result.data?.newCost.id,
					name: args.formData.get('name') as string,
					amount: parseFloat(args.formData.get('amount') as string)
				};
				if (this.costs.length < this.costLimit) {
					this.costs.push(newCost);
				}
				this.costTotalRow += 1;
				CreateToast('success', 'successfully create cost');
			}
			this.createCostDialogOpen = false;
		};
	};
	onDeleteCost = (args: EnhancementArgs) => {
		this.isLoading = true;
		if (this.costToDelete) {
			args.formData.append('id', `${this.costToDelete}`);
		}
		return async ({ result }: EnhancementReturn) => {
			this.isLoading = false;
			if (result.type === 'success') {
				if (this.costToDelete) {
					this.costs = this.costs.filter((c) => c.id !== this.costToDelete);
				}
				this.costTotalRow -= 1;
				this.costToDelete = undefined;
				CreateToast('success', 'successfully delete cost');
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			this.deleteCostDialogOpen = false;
		};
	};
	onUpdateCost = (args: EnhancementArgs) => {
		this.isLoading = true;
		if (this.costToUpdate) {
			args.formData.append('id', `${this.costToUpdate}`);
		}
		return async ({ result }: EnhancementReturn) => {
			this.isLoading = false;
			if (result.type === 'success') {
				if (this.costToUpdate) {
					this.costs = this.costs.map((m) => {
						if (m.id === this.costToUpdate) {
							return {
								id: this.costToUpdate,
								name: args.formData.get('name') as string,
								amount: parseFloat(args.formData.get('amount') as string)
							};
						}
						return m;
					});
				}
				this.costToUpdate = undefined;
				CreateToast('success', 'successfully update cost');
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			this.updateCostDialogOpen = false;
		};
	};
	onCostPageChangeForm = (args: EnhancementArgs) => {
		this.isLoading = true;
		args.formData.append('page', `${this.costPageNumber}`);
		return async ({ result }: EnhancementReturn) => {
			this.isLoading = false;
			if (result.type === 'success') {
				this.costs = result.data?.costs.entries;
				this.costLimit = result.data?.costs.page_info.limit;
				this.costTotalRow = result.data?.costs.page_info.total_row;
				this.costPageNumber = result.data?.costs.page_info.page;
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
