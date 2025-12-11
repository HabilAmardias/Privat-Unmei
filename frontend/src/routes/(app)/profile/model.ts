export type StudentProfile = {
	name: string;
	bio: string;
	profile_image: string;
	public_id: string;
};

export type StudentOrders = {
	id: string;
	student_id: string;
	course_id: number;
	total_price: number;
	status: string;
	mentor_name: string;
	mentor_public_id: string;
	course_name: string;
};
