import type { PaginatedResponse } from '$lib/types';
import type { CourseCategoryOpts, CourseList, CourseCategory } from './model';
import { debounce } from '$lib/utils/helper';
import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
import { CreateToast } from '$lib/utils/helper';

export class CourseListView {
	courses = $state<CourseList[]>([]);
	page = $state<number>(1);
	limit = $state<number>(15);
	totalRow = $state<number>(15);

	search = $state<string>('');
	categories = $state<CourseCategoryOpts[]>([]);
	method = $state<string>('');

	isLoading = $state<boolean>(false);
	selectedCategory = $state<string>('');
	categoryForm = $state<HTMLFormElement>();
	searchCategory = $state<string>('');

	paginationForm = $state<HTMLFormElement>();

	constructor(data: PaginatedResponse<CourseList>) {
		this.courses = data.entries;
		this.page = data.page_info.page;
		this.limit = data.page_info.limit;
		this.totalRow = data.page_info.total_row;
	}

	#searchCategorySubmit = debounce(() => {
		this.categoryForm?.requestSubmit();
	}, 500);
	onSearchCategory = (e: Event & { currentTarget: EventTarget & HTMLInputElement }) => {
		this.searchCategory = e.currentTarget.value;
		this.#searchCategorySubmit();
	};
	onGetCategory = (args: EnhancementArgs) => {
		args.formData.append('search', this.searchCategory);
		return async ({ result }: EnhancementReturn) => {
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'success') {
				this.#convertCategory(result.data?.categories);
			}
		};
	};
	#convertCategory = (c: CourseCategory[]) => {
		const options: CourseCategoryOpts[] = [];
		c.forEach((item) => {
			options.push({
				value: `${item.id}`,
				label: item.name
			});
		});
		this.categories = options;
	};
	onSearchCourse = (args: EnhancementArgs) => {
		this.isLoading = true;
		this.page = 1;
		args.formData.append('course_category', `${this.selectedCategory}`);
		args.formData.append('page', `${this.page}`);
		return async ({ result }: EnhancementReturn) => {
			this.isLoading = false;
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'success') {
				this.limit = result.data?.courses.page_info.limit;
				this.page = result.data?.courses.page_info.page;
				this.totalRow = result.data?.courses.page_info.total_row;
				this.courses = result.data?.courses.entries;
			}
		};
	};
	onPageChange = () => {
		this.paginationForm?.requestSubmit();
	};
	onPageChangeEnhance = (args: EnhancementArgs) => {
		args.formData.append('search', this.search);
		args.formData.append('course_category', `${this.selectedCategory}`);
		args.formData.append('page', `${this.page}`);
		args.formData.append('method', this.method);
		return async ({ result }: EnhancementReturn) => {
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'success') {
				this.limit = result.data?.courses.page_info.limit;
				this.page = result.data?.courses.page_info.page;
				this.totalRow = result.data?.courses.page_info.total_row;
				this.courses = result.data?.courses.entries;
			}
		};
	};
}
