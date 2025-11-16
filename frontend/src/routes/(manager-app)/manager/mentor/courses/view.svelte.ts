import type { EnhancementArgs, EnhancementReturn, PaginatedResponse } from '$lib/types';
import type { CourseCategory, CourseCategoryOpts, MentorCourse } from './model';
import { CreateToast, DismissToast } from '$lib/utils/helper';
import { debounce } from '$lib/utils/helper';

export class CourseManagementView {
	search = $state<string>('');
	categories = $state<CourseCategoryOpts[]>([]);
	courses = $state<MentorCourse[]>([]);
	searchCategory = $state<string>('');
	selectedCategory = $state<string>('');
	searchCategoryForm = $state<HTMLFormElement>();
	isLoading = $state<boolean>(false);

	limit = $state<number>(15);
	page = $state<number>(1);
	totalRow = $state<number>(0);

	courseToDelete = $state<number | undefined>(undefined);
	deleteCourseOpen = $derived<boolean[]>(new Array<boolean>(this.courses.length).fill(false));
	paginationForm = $state<HTMLFormElement>();
	#searchCategorySubmit = debounce(() => {
		this.searchCategoryForm?.requestSubmit();
	}, 500);

	constructor(categories: CourseCategory[], courses: PaginatedResponse<MentorCourse>) {
		this.limit = courses.page_info.limit;
		this.page = courses.page_info.page;
		this.totalRow = courses.page_info.total_row;
		this.courses = courses.entries;
		this.#ConvertCategoryOpts(categories);
	}

	onSearchCategoryChange = (e: Event & { currentTarget: EventTarget & HTMLInputElement }) => {
		this.searchCategory = e.currentTarget.value;
		this.#searchCategorySubmit();
	};

	onPageChange = () => {
		this.paginationForm?.requestSubmit();
	};

	onPageChangeForm = (args: EnhancementArgs) => {
		args.formData.append('search', this.search);
		args.formData.append('page', `${this.page}`);
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

	onSearchCategory = (args: EnhancementArgs) => {
		args.formData.append('limit', '5');
		args.formData.append('search', this.searchCategory);
		return async ({ result }: EnhancementReturn) => {
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'success') {
				this.#ConvertCategoryOpts(result.data?.categories);
			}
		};
	};

	onSearchCourse = (args: EnhancementArgs) => {
		this.isLoading = true;
		this.page = 1;
		args.formData.append('category', `${this.selectedCategory}`);
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

	onDeleteCourse = (args: EnhancementArgs) => {
		const loadID = CreateToast('loading', 'deleting course....');
		if (this.courseToDelete) {
			args.formData.append('id', `${this.courseToDelete}`);
		}
		return async ({ result }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'success') {
				if (this.courseToDelete) {
					this.courses = this.courses.filter((v) => v.id !== this.courseToDelete);
				}
				this.totalRow -= 1;
				this.courseToDelete = undefined;
				CreateToast('success', 'successfully delete course');
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};

	#ConvertCategoryOpts = (cats: CourseCategory[]) => {
		const opts: CourseCategoryOpts[] = [];
		cats.forEach((item) => {
			opts.push({
				value: `${item.id}`,
				label: item.name
			});
		});
		this.categories = opts;
	};
}
