export type PaymentMethods = {
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
