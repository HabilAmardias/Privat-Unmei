import type { EnhancementArgs, EnhancementReturn, PaginatedResponse } from '$lib/types';
import { CreateToast, DismissToast } from '$lib/utils/helper';
import { SvelteDate } from 'svelte/reactivity';
import type { CourseDetail, CourseReview, StudentProfile } from './model';

export class CourseDetailView {
	reviews = $state<CourseReview[]>([]);
	isLoading = $state<boolean>(false);
	page = $state<number>(1);
	limit = $state<number>(15);
	totalRow = $state<number>(15);
	paginationForm = $state<HTMLFormElement>();
	star = $state<number>(1);
	feedback = $state<string>('');
	profile = $state<StudentProfile>();
	mentorID = $state<string>('');
	feedbackErr = $derived.by<Error | undefined>(() => {
		if (this.feedback.length < 15) {
			return new Error('Feedback must at least contain 15 characters');
		}
		return undefined;
	});
	reviewDisabled = $derived.by<boolean>(() => {
		return this.feedbackErr ? true : false;
	});

	detailState = $state<'description' | 'detail'>('description');

	constructor(d: PaginatedResponse<CourseReview>, c: CourseDetail, p?: StudentProfile) {
		this.reviews = d.entries;
		this.page = d.page_info.page;
		this.limit = d.page_info.limit;
		this.totalRow = d.page_info.total_row;
		if (p) {
			this.profile = p;
		}
		this.mentorID = c.mentor_id;
	}
	capitalizeFirstLetter(s: string) {
		if (s.length === 0) {
			return s;
		}
		if (s.length === 1) {
			return s.toUpperCase();
		}
		return s.charAt(0).toUpperCase() + s.slice(1);
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
	onCreateReview = () => {
		const loadID = CreateToast('loading', 'loading....');
		return async ({ result }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'success') {
				CreateToast('success', 'successfully create review');
				this.totalRow += 1;
				const now = new SvelteDate().toDateString();
				if (this.reviews.length < this.limit) {
					this.reviews.push({
						id: result.data?.id,
						course_id: result.data?.course_id,
						student_id: this.profile!.id,
						name: this.profile!.name,
						rating: this.star,
						feedback: this.feedback,
						created_at: now
					});
				}
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
	onMessageMentor = (args: EnhancementArgs) => {
		const loadID = CreateToast('loading', 'loading....');
		args.formData.append('id', this.mentorID);
		return async ({ result, update }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'redirect') {
				await update();
			}
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
		};
	};
}
