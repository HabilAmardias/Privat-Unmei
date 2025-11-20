import type { CourseCategory } from './model';
import type { EnhancementArgs, EnhancementReturn, PaginatedResponse } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';
import { debounce } from '$lib/utils/helper';

export class CourseCategoryManagementView {
	categories = $state<CourseCategory[]>([]);
	pageNumber = $state<number>(1);
	totalRow = $state<number>(0);
	search = $state<string>('');
	isLoading = $state<boolean>(false);
	deleteDialogOpen = $state<boolean[]>([]);
	updateDialogOpen = $state<boolean[]>([]);
	createDialogOpen = $state<boolean>(false);
	categoryToDelete = $state<number>();
	limit = $state<number>(0);
	createCategoryForm = $state<HTMLFormElement>();
	searchForm = $state<HTMLFormElement>();
	paginationForm = $state<HTMLFormElement>();
	categoryToUpdate = $state<number>();

	#SearchSubmit = debounce(() => {
		this.searchForm?.requestSubmit();
	}, 500);

	constructor(c: PaginatedResponse<CourseCategory>) {
		this.pageNumber = c.page_info.page;
		this.categories = c.entries;
		this.totalRow = c.page_info.total_row;
		this.limit = c.page_info.limit;
		this.deleteDialogOpen = new Array<boolean>(this.categories.length).fill(false);
		this.updateDialogOpen = new Array<boolean>(this.categories.length).fill(false);
	}
	onSearchInput = () => {
		this.#SearchSubmit();
	};
	setPageNumber = () => {
		this.paginationForm?.requestSubmit();
	};
	onCreateCategory = (args: EnhancementArgs) => {
		const loadID = CreateToast('loading', 'Creating Course Category.....');
		return async ({ result }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'success') {
				const newCategory: CourseCategory = {
					id: result.data?.newCategory.id,
					name: args.formData.get('name') as string
				};
				if (this.categories.length < this.limit) {
					this.categories.push(newCategory);
					this.deleteDialogOpen = new Array<boolean>(this.categories.length).fill(false);
					this.updateDialogOpen = new Array<boolean>(this.categories.length).fill(false);
				}
				this.totalRow += 1;
				CreateToast('success', 'successfully create course category');
			}
			this.createDialogOpen = false;
		};
	};
	onDeleteCategory = (args: EnhancementArgs) => {
		this.isLoading = true;
		if (this.categoryToDelete) {
			args.formData.append('id', `${this.categoryToDelete}`);
		}
		return async ({ result }: EnhancementReturn) => {
			this.isLoading = false;
			if (result.type === 'success') {
				if (this.categoryToDelete) {
					this.categories = this.categories.filter((m) => m.id !== this.categoryToDelete);
					this.deleteDialogOpen = new Array<boolean>(this.categories.length).fill(false);
					this.updateDialogOpen = new Array<boolean>(this.categories.length).fill(false);
				}
				this.totalRow -= 1;
				this.categoryToDelete = undefined;
				CreateToast('success', 'successfully delete category');
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	onUpdateCategory = (args: EnhancementArgs) => {
		this.isLoading = true;
		if (this.categoryToUpdate) {
			args.formData.append('id', `${this.categoryToUpdate}`);
		}
		return async ({ result }: EnhancementReturn) => {
			this.isLoading = false;
			if (result.type === 'success') {
				if (this.categoryToUpdate) {
					this.categories = this.categories.map((m) => {
						if (m.id === this.categoryToUpdate) {
							return {
								id: this.categoryToUpdate,
								name: args.formData.get('name') as string
							};
						}
						return m;
					});
				}
				this.categoryToUpdate = undefined;
				CreateToast('success', 'successfully update course category');
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
	onSearchCategory = (args: EnhancementArgs) => {
		this.isLoading = true;
		this.pageNumber = 1;
		args.formData.append('page', `${this.pageNumber}`);
		return async ({ result }: EnhancementReturn) => {
			this.isLoading = false;
			if (result.type === 'success') {
				this.categories = result.data?.categories.entries;
				this.deleteDialogOpen = new Array<boolean>(this.categories.length).fill(false);
				this.updateDialogOpen = new Array<boolean>(this.categories.length).fill(false);
				this.setPaginationData(
					result.data?.categories.page_info.limit,
					result.data?.categories.page_info.total_row,
					result.data?.categories.page_info.page
				);
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	onPageChangeForm = (args: EnhancementArgs) => {
		this.isLoading = true;
		args.formData.append('page', `${this.pageNumber}`);
		args.formData.append('search', this.search);
		return async ({ result }: EnhancementReturn) => {
			this.isLoading = false;
			if (result.type === 'success') {
				this.categories = result.data?.categories.entries;
				this.deleteDialogOpen = new Array<boolean>(this.categories.length).fill(false);
				this.updateDialogOpen = new Array<boolean>(this.categories.length).fill(false);
				this.setPaginationData(
					result.data?.categories.page_info.limit,
					result.data?.categories.page_info.total_row,
					result.data?.categories.page_info.page
				);
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
