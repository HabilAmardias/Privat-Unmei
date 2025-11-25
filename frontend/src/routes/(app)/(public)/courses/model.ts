export type CourseList = {
	id: string;
	title: string;
	domicile: string;
	method: string;
	price: number;
	session_duration_minutes: number;
	max_total_session: number;
	mentor_id: string;
	mentor_name: string;
	mentor_email: string;
};

export type CourseCategoryOpts = {
	value: string;
	label: string;
};

export type CourseCategory = {
	id: number;
	name: string;
};
