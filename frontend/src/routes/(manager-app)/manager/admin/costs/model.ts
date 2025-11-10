export type AdditionalCost = {
	id: number;
	name: string;
	amount: number;
};

export type Discount = {
	id: number;
	number_of_participant: number;
	amount: number;
};

export type adminProfile = {
	name: string;
	email: string;
	bio: string;
	profile_image: string;
	status: 'verified' | 'unverified';
};

export type newCost = {
	id: number;
};

export type newDiscount = {
	id: number;
};
