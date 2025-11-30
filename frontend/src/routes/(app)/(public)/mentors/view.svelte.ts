import type { PaginatedResponse } from '$lib/types';
import type { MentorList } from './model';
import { debounce } from '$lib/utils/helper';
import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
import { CreateToast } from '$lib/utils/helper';

export class MentorListView {
	mentors = $state<MentorList[]>([]);
	page = $state<number>(1);
	limit = $state<number>(15);
	totalRow = $state<number>(15);

	search = $state<string>('');

	isLoading = $state<boolean>(false);
	searchForm = $state<HTMLFormElement>();

	paginationForm = $state<HTMLFormElement>();

	constructor(data: PaginatedResponse<MentorList>) {
		this.mentors = data.entries;
		this.page = data.page_info.page;
		this.limit = data.page_info.limit;
		this.totalRow = data.page_info.total_row;
	}

	#searchCategorySubmit = debounce(() => {
		this.searchForm?.requestSubmit();
	}, 500);

	onSearchMentor = () => {
		this.#searchCategorySubmit();
	};

	onSearchMentorEnhance = (args: EnhancementArgs) => {
		this.isLoading = true;
		this.page = 1;
		args.formData.append('page', `${this.page}`);
		return async ({ result }: EnhancementReturn) => {
			this.isLoading = false;
			if (result.type === 'success') {
				this.mentors = result.data?.mentors.entries;
				this.page = result.data?.mentorsList.page_info.page;
				this.limit = result.data?.mentorsList.page_info.limit;
				this.totalRow = result.data?.mentorsList.page_info.total_row;
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};

	onPageChange = () => {
		this.paginationForm?.requestSubmit();
	};

	onPageChangeEnhance = (args: EnhancementArgs) => {
		this.isLoading = true;
		args.formData.append('search', this.search);
		args.formData.append('page', `${this.page}`);
		return async ({ result }: EnhancementReturn) => {
			this.isLoading = false;
			if (result.type === 'success') {
				this.mentors = result.data?.mentors.entries;
				this.page = result.data?.mentorsList.page_info.page;
				this.limit = result.data?.mentorsList.page_info.limit;
				this.totalRow = result.data?.mentorsList.page_info.total_row;
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
