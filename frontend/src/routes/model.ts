export type LoginResponse = {
	token: string;
	status: 'verified' | 'unverified';
};
