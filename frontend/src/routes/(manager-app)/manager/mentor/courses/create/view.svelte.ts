import type { CourseCategory, CourseCategoryOpts, CourseTopic } from './model';

export class CreateCourseView {
	title = $state<string>('');
	description = $state<string>('');
	domicile = $state<string>('');
	price = $state<number>(0);
	method = $state<string>('');
	sessionDuration = $state<number>(1);
	maxSession = $state<number>(1);
	categories = $state<CourseCategoryOpts[]>([]);
	addedCategories = $state<CourseCategory[]>([]);
	selectedCategory = $state<string>('');

	topicTitle = $state<string>('');
	topicDescription = $state<string>('');
	addedTopic = $state<CourseTopic[]>([]);

	searchCategory = $state<string>('');
}
