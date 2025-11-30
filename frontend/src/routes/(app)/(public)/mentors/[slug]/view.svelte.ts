import type { EnhancementReturn, EnhancementArgs, PaginatedResponse } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';
import type { MentorCourse } from './model';

export class MentorDetailView {
	isDesktop = $state<boolean>(false);
	alertOpen = $state<boolean>(false);
	isLoading = $state<boolean>(false);
	size = $derived.by<number>(() => {
		if (this.isDesktop) {
			return 150;
		}
		return 100;
	});

	courses = $state<MentorCourse[]>([]);
	page = $state<number>(1);
	limit = $state<number>(15);
	totalRow = $state<number>(15);
	paginationForm = $state<HTMLFormElement>();

	constructor(d: PaginatedResponse<MentorCourse>) {
		this.courses = d.entries;
		this.page = d.page_info.page;
		this.limit = d.page_info.limit;
		this.totalRow = d.page_info.total_row;
	}

	onDeleteMentor = () => {
		const loadID = CreateToast('loading', 'deleting....');
		return async ({ result, update }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'redirect') {
				CreateToast('success', 'Successfully delete mentor');
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			update();
		};
	};
	setIsDesktop(b: boolean) {
		this.isDesktop = b;
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
				this.limit = result.data?.courses.page_info.limit;
				this.page = result.data?.courses.page_info.page;
				this.totalRow = result.data?.courses.page_info.total_row;
				this.courses = result.data?.courses.entries;
			}
		};
	};
}
