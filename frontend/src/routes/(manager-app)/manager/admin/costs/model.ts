export type AdditionalCost = {
	id: number;
	name: string;
	amount: number;
};

export type adminProfile = {
	name: string;
	email: string;
	bio: string;
	profile_image: string;
	status: 'verified' | 'unverified';
};
