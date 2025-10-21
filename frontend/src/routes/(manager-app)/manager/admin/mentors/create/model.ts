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
	account_number: string;
};
