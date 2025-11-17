import type { CourseCategory, CourseCategoryOpts, CourseTopic } from './model';
import { CreateToast, debounce, DismissToast } from '$lib/utils/helper';
import type { EnhancementArgs, EnhancementReturn } from '$lib/types';

export class CreateCourseView {
	title = $state<string>('');
	description = $state<string>('');
	domicile = $state<string>('');
	price = $state<number>(1);
	method = $state<string>('');
	sessionDuration = $state<number>(1);
	maxSession = $state<number>(1);
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
		if (!this.selectedCategory) {
			return true;
		}
		return false;
	});

	disableCreateCourse = $derived.by<boolean>(() => {
		if (
			!this.title ||
			!this.description ||
			!this.domicile ||
			!this.price ||
			!this.method ||
			!this.sessionDuration ||
			!this.maxSession ||
			this.addedCategories.length === 0 ||
			this.addedTopic.length === 0
		) {
			return true;
		}
		return false;
	});

	constructor(c: CourseCategory[]) {
		this.#convertCategory(c);
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
	onCreateCourse = (args: EnhancementArgs) => {
		const catIDs = this.addedCategories.map<number>((item) => {
			return item.id;
		});
		args.formData.append('categories', catIDs.join(','));
		args.formData.append('topics', JSON.stringify(this.addedTopic));
		const loadID = CreateToast('loading', 'creating course....');
		return async ({ result }: EnhancementReturn) => {
			DismissToast(loadID);
			if (result.type === 'failure') {
				CreateToast('error', result.data?.message);
			}
			if (result.type === 'success') {
				CreateToast('success', 'successfully create course');
			}
		};
	};
}
