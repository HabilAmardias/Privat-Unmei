export type PaymentMethods = {
	payment_method_id: number;
	payment_method_name: string;
};

export type adminProfile = {
	name: string;
	email: string;
	bio: string;
	profile_image: string;
	status: 'verified' | 'unverified';
};

export type NewPaymentMethod = {
	payment_method_id: number;
};
