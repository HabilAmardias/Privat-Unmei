export type TimeOnly = {
	hour: number;
	minute: number;
	second: number;
};

export type MentorScheduleInfo = {
	day_of_week: number;
	start_time: string;
	end_time: string;
};

export type MentorPaymentInfo = {
	payment_method_id: number;
	payment_method_name: string;
	account_number: string;
};

export type StudentProfile = {
	id: string;
	name: string;
	bio: string;
	profile_image: string;
	email: string;
};

export type MentorProfile = {
	id: string;
	profile_image: string;
	name: string;
	bio: string;
	years_of_experience: number;
	degree: string;
	major: string;
	campus: string;
	email: string;
};

export type MentorCourse = {
	id: number;
	title: string;
	domicile: string;
	method: string;
	price: number;
	session_duration_minutes: number;
	max_total_session: number;
};

export type ChatroomID = {
	id: number;
};
