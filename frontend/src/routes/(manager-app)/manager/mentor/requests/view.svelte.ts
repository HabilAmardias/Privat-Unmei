import type { EnhancementArgs, EnhancementReturn, PaginatedResponse } from '$lib/types';
import type { CourseRequest } from './model';
import { CreateToast } from '$lib/utils/helper';

export class RequestManagementView {
	isLoading = $state<boolean>(false);
	requests = $state<CourseRequest[]>([]);
	status = $state<string>('');
	page = $state<number>(1);
	limit = $state<number>(1);
	totalRow = $state<number>(1);
	paginationForm = $state<HTMLFormElement>();
	searchForm = $state<HTMLFormElement>();

	constructor(r: PaginatedResponse<CourseRequest>) {
		this.requests = r.entries;
		this.page = r.page_info.page;
		this.limit = r.page_info.limit;
		this.totalRow = r.page_info.total_row;
	}

	setPageNumber = (num: number) => {
		this.page = num;
		this.paginationForm?.requestSubmit();
	};

	onPageChange = (args: EnhancementArgs) => {
		this.isLoading = true;
		args.formData.append('page', `${this.page}`);
		args.formData.append('status', this.status);
		return async ({ result }: EnhancementReturn) => {
			this.isLoading = false;
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'success') {
				this.requests = result.data?.entries;
				this.page = result.data?.page_info.page;
				this.limit = result.data?.page_info.limit;
				this.totalRow = result.data?.page_info.total_row;
			}
		};
	};

	onSearchRequest = (args: EnhancementArgs) => {
		this.isLoading = true;
		this.page = 1;
		args.formData.append('page', `${this.page}`);
		return async ({ result }: EnhancementReturn) => {
			this.isLoading = false;
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'success') {
				this.requests = result.data?.entries;
				this.page = result.data?.page_info.page;
				this.limit = result.data?.page_info.limit;
				this.totalRow = result.data?.page_info.total_row;
			}
		};
	};
}
