export type RequestSchedule = {
	date: string;
	start_time: string;
	end_time: string;
};

export type RequestDetail = {
	course_request_id: string;
	course_name: string;
	course_id: number;
	mentor_name: string;
	mentor_id: string;
	mentor_public_id: string;
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

export type ChatroomID = {
	id: number;
};

export type CreateReview = {
	id: number;
};

export type IsReviewed = {
	is_reviewed: boolean;
};
