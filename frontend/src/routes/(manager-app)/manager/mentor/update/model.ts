export type paymentMethod = {
	payment_method_id: number;
	payment_method_name: string;
};

export type paymentMethodOpts = {
	value: string;
	label: string;
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

export type TimeOnly = {
	hour: number;
	minute: number;
	second: number;
};

export type MentorSchedule = {
	day_of_week: number;
	day_of_week_label: string;
	start_time: TimeOnly;
	end_time: TimeOnly;
};

export type MentorProfile = {
	id: string;
	resume: string;
	profile_image: string;
	name: string;
	bio: string;
	years_of_experience: number;
	degree: string;
	major: string;
	campus: string;
	email: string;
};
