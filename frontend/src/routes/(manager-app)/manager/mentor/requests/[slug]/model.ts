export type RequestSchedule = {
	date: string;
	start_time: string;
	end_time: string;
};

export type RequestDetail = {
	course_request_id: number;
	course_name: string;
	student_name: string;
	student_email: string;
	total_price: number;
	subtotal: number;
	operational_cost: number;
	number_of_sessions: number;
	status: 'reserved' | 'pending payment' | 'scheduled' | 'completed' | 'cancelled';
	expired_at: string | null;
	payment_method: string;
	account_number: string;
	number_of_participant: number;
	schedules: RequestSchedule[];
};
