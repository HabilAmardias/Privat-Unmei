export type MentorCourse = {
	id: number;
	title: string;
	domicile: string;
	method: string;
	price: number;
	session_duration_minutes: number;
	max_total_session: number;
};

export type CourseCategory = {
	id: number;
	name: string;
};

export type CourseCategoryOpts = {
	value: string;
	label: string;
};
