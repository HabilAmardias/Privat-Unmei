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
	mentor_public_id: string;
};

export type MentorList = {
	id: string;
	name: string;
	public_id: string;
	profile_image: string;
	years_of_experience: number;
};
