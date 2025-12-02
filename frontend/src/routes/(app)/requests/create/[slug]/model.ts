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
	mentor_profile_image: string;
	description: string;
};

export type TimeOnly = {
	hour: number;
	minute: number;
	second: number;
};

export type ScheduleSlot = {
	date: string;
	start_time: TimeOnly;
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

export type DowRes = {
	day_of_weeks: number[];
};
