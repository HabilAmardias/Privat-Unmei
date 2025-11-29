import type { EnhancementArgs, EnhancementReturn, PaginatedResponse } from '$lib/types';
import { CreateToast } from '$lib/utils/helper';
import type { CourseReview } from './model';

export class CourseDetailView {
	reviews = $state<CourseReview[]>([]);
	isLoading = $state<boolean>(false);
	page = $state<number>(1);
	limit = $state<number>(15);
	totalRow = $state<number>(15);
	paginationForm = $state<HTMLFormElement>();

	constructor(d: PaginatedResponse<CourseReview>) {
		this.reviews = d.entries;
		this.page = d.page_info.page;
		this.limit = d.page_info.limit;
		this.totalRow = d.page_info.total_row;
	}
	onPageChange = () => {
		this.paginationForm?.requestSubmit();
	};
	onPageChangeEnhance = (args: EnhancementArgs) => {
		this.isLoading = true;
		args.formData.append('page', `${this.page}`);
		return async ({ result }: EnhancementReturn) => {
			this.isLoading = false;
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'success') {
				this.limit = result.data?.reviews.page_info.limit;
				this.page = result.data?.reviews.page_info.page;
				this.totalRow = result.data?.reviews.page_info.total_row;
				this.reviews = result.data?.reviews.entries;
			}
		};
	};
}
