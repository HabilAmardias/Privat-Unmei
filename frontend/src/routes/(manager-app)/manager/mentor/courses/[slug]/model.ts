export type CourseCategory = {
	id: number;
	name: string;
};

export type CourseTopic = {
	title: string;
	description: string;
};

export type CourseDetail = {
	id: number;
	title: string;
	domicile: string;
	method: string;
	price: number;
	session_duration_minutes: number;
	max_total_session: number;
	mentor_id: string;
	mentor_name: string;
	mentor_email: string;
	description: string;
};
