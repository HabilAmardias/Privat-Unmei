export type paymentMethod = {
	payment_method_id: number;
	payment_method_name: string;
};

export type paymentMethodOpts = {
	value: string;
	label: string;
};

export type mentorPaymentMethods = {
	payment_method_id: number;
	payment_method_name: string;
	account_number: string;
};

export type generatedPassword = {
	password: string;
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
