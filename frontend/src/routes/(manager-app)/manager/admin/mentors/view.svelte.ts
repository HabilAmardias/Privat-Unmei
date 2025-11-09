import type { EnhancementArgs, EnhancementReturn, PaginatedResponse } from '$lib/types';
import type { mentorList } from './model';
import { CreateToast } from '$lib/utils/helper';
import { debounce } from '$lib/utils/helper';

export class MentorManagerView {
	mentors = $state<mentorList[]>([]);
	page = $state<number>(1);
	limit = $state<number>(15);
	total_row = $state<number>(15);
	isDesktop = $state<boolean>(false);
	mentorsIsLoading = $state<boolean>(false);
	mentorToDelete = $state<string>();
	sortByYears = $state<boolean | null>(null);
	search = $state<string>('');
	alertOpen = $state<boolean>(false);
	searchForm = $state<HTMLFormElement | null>(null);
	#SearchSubmit = debounce(() => {
		this.searchForm?.requestSubmit();
	}, 500);

	constructor(m: PaginatedResponse<mentorList>) {
		this.setMentors(m.entries);
		this.setPaginationData(m.page_info.page, m.page_info.limit, m.page_info.total_row);
	}

	onSearchInput = (e: Event & { currentTarget: EventTarget & HTMLInputElement }) => {
		this.search = e.currentTarget.value;
		this.#SearchSubmit();
	};
	setMentors(newList: mentorList[]) {
		this.mentors = newList;
	}
	setMentorToDelete(id: string | undefined) {
		this.mentorToDelete = id;
	}
	filterMentors(id: string) {
		this.mentors = this.mentors.filter((m) => m.id !== id);
	}
	setPaginationData(page: number, limit: number, total_row: number) {
		this.page = page;
		this.limit = limit;
		this.total_row = total_row;
	}
	setIsDesktop(b: boolean) {
		this.isDesktop = b;
	}
	onPageChange(num: number) {
		this.page = num;
	}
	onSort() {
		switch (this.sortByYears) {
			case null:
				this.sortByYears = true;
				break;
			case true:
				this.sortByYears = false;
				break;
			case false:
				this.sortByYears = null;
				break;
		}
	}
	resetFilterForm() {
		this.search = '';
		this.sortByYears = null;
	}
	setMentorsIsLoading(b: boolean) {
		this.mentorsIsLoading = b;
	}
	onDeleteMentor = (args: EnhancementArgs) => {
		this.setMentorsIsLoading(true);
		if (this.mentorToDelete) {
			args.formData.append('id', this.mentorToDelete);
		}
		return async ({ result }: EnhancementReturn) => {
			this.setMentorsIsLoading(false);
			if (result.type === 'success') {
				if (this.mentorToDelete) {
					this.filterMentors(this.mentorToDelete);
				}
				this.setMentorToDelete(undefined);
				this.total_row -= 1;
				CreateToast('success', 'successfully delete mentor');
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	onUpdateMentors = (args: EnhancementArgs) => {
		this.setMentorsIsLoading(true);
		if (this.sortByYears !== null) {
			args.formData.append('sort_year_of_experience', `${this.sortByYears}`);
		}
		args.formData.append('page', `${this.page}`);
		return async ({ result }: EnhancementReturn) => {
			this.setMentorsIsLoading(false);
			if (result.type === 'success') {
				this.setMentors(result.data?.mentorsList.entries);
				this.setPaginationData(
					result.data?.mentorsList.page_info.page,
					result.data?.mentorsList.page_info.limit,
					result.data?.mentorsList.page_info.total_row
				);
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
