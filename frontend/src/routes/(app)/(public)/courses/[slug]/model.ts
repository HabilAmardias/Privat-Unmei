export type CourseCategory = {
	id: number;
	name: string;
};

export type CourseTopic = {
	title: string;
	description: string;
};

export type CourseReview = {
	id: number;
	course_id: number;
	student_id: string;
	name: string;
	rating: number;
	feedback: string;
	created_at: string;
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
	mentor_public_id: string;
	mentor_profile_image: string;
	description: string;
};

export type StudentProfile = {
	id: string;
	name: string;
	bio: string;
	profile_image: string;
	public_id: string;
};

export type ChatroomID = {
	id: number;
};
