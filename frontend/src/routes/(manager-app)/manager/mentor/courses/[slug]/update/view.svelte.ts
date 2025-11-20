import type { CourseCategory, CourseCategoryOpts, CourseDetail, CourseTopic } from './model';
import { CreateToast, debounce, DismissToast } from '$lib/utils/helper';
import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
import { MAX_COURSE_CATEGORIES_COUNT } from './constants';
import { MAX_BIO_LENGTH } from '$lib/utils/constants';

export class UpdateCourseView {
	title = $state<string>('');
	description = $state<string>('');
	descriptionErr = $derived.by<Error | undefined>(() => {
		if (this.description.length > MAX_BIO_LENGTH) {
			return new Error(`description can only consist of ${MAX_BIO_LENGTH} characters`);
		}
		return undefined;
	});
	domicile = $state<string>('');
	price = $state<number>(1);
	priceErr = $derived.by<Error | undefined>(() => {
		if (this.price <= 0) {
			return new Error('Price must be greater than zero');
		}
		return undefined;
	});
	method = $state<string>('');
	sessionDuration = $state<number>(1);
	sessionDurationErr = $derived.by<Error | undefined>(() => {
		if (this.sessionDuration <= 0) {
			return new Error('Session duration must be greater than zero');
		}
		return undefined;
	});
	maxSession = $state<number>(1);
	maxSessionErr = $derived.by<Error | undefined>(() => {
		if (this.maxSession <= 0) {
			return new Error('Max Session Count must be greater than zero');
		}
		return undefined;
	});
	categories = $state<CourseCategoryOpts[]>([]);
	addedCategories = $state<CourseCategory[]>([]);
	searchCategoryForm = $state<HTMLFormElement>();
	selectedCategory = $state<string>('');
	disableAddTopic = $derived.by<boolean>(() => {
		if (!this.topicTitle || !this.topicDescription) {
			return true;
		}
		return false;
	});

	topicTitle = $state<string>('');
	topicDescription = $state<string>('');
	addedTopic = $state<CourseTopic[]>([]);

	searchCategory = $state<string>('');
	#searchCategorySubmit = debounce(() => {
		this.searchCategoryForm?.requestSubmit();
	}, 500);

	disableAddCategory = $derived.by<boolean>(() => {
		if (!this.selectedCategory || this.addedCategoryErr) {
			return true;
		}
		return false;
	});

	addedCategoryErr = $derived.by<Error | null>(() => {
		if (this.addedCategories.length > MAX_COURSE_CATEGORIES_COUNT) {
			return new Error(`cannot have more than ${MAX_COURSE_CATEGORIES_COUNT} categories`);
		}
		return null;
	});

	disableUpdateCourse = $derived.by<boolean>(() => {
		if (
			!this.title ||
			!this.description ||
			!this.domicile ||
			!this.price ||
			!this.method ||
			!this.sessionDuration ||
			!this.maxSession ||
			this.descriptionErr ||
			this.addedCategoryErr ||
			this.maxSessionErr ||
			this.priceErr ||
			this.sessionDurationErr
		) {
			return true;
		}
		return false;
	});

	constructor(c: CourseCategory[], t: CourseTopic[], ac: CourseCategory[], d: CourseDetail) {
		this.#convertCategory(c);
		this.addedTopic = t;
		this.addedCategories = ac;
		this.title = d.title;
		this.description = d.description;
		this.price = d.price;
		this.domicile = d.domicile;
		this.method = d.method;
		this.sessionDuration = d.session_duration_minutes;
		this.maxSession = d.max_total_session;
	}
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
	addCourseTopic = () => {
		this.addedTopic.push({
			title: this.topicTitle,
			description: this.topicDescription
		});
		this.topicTitle = '';
		this.topicDescription = '';
	};

	addCourseCategory = () => {
		const name = this.categories.filter((v) => v.value === this.selectedCategory)[0].label;
		this.addedCategories.push({
			id: parseInt(this.selectedCategory),
			name
		});
	};

	removeAddedTopic = (idx: number) => {
		this.addedTopic = this.addedTopic.filter((v, i) => idx !== i);
	};

	removeAddedCategories = (idx: number) => {
		this.addedCategories = this.addedCategories.filter((v, i) => idx !== i);
	};
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
	onUpdateCourse = (args: EnhancementArgs) => {
		const catIDs = this.addedCategories.map<number>((item) => {
			return item.id;
		});
		args.formData.append('categories', catIDs.join(','));
		args.formData.append('topics', JSON.stringify(this.addedTopic));
		const loadID = CreateToast('loading', 'updating course....');
		return async ({ result, update }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'redirect') {
				CreateToast('success', 'successfully update course');
				await update();
			}
		};
	};
}
