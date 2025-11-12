export type NewCategory = {
	id: number;
};

export type CourseCategory = {
	id: number;
	name: string;
};

export type adminProfile = {
	name: string;
	email: string;
	bio: string;
	profile_image: string;
	status: 'verified' | 'unverified';
};
