export type TimeOnly = {
	hour: number;
	minute: number;
	second: number;
};
export type MentorScheduleInfo = {
	day_of_week: number;
	start_time: TimeOnly;
	end_time: TimeOnly;
};

export type MentorPaymentInfo = {
	payment_method_id: number;
	payment_method_name: string;
	account_number: string;
};

export type adminProfile = {
	name: string;
	email: string;
	bio: string;
	profile_image: string;
	status: 'verified' | 'unverified';
};

export type MentorProfile = {
	resume_file: string;
	profile_image: string;
	name: string;
	bio: string;
	years_of_experience: number;
	degree: string;
	major: string;
	campus: string;
	mentor_availability: MentorScheduleInfo[];
	mentor_payment_info: MentorPaymentInfo[];
};
