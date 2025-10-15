export type AdminLoginResponse = {
	status: 'verified' | 'unverified';
};

export type ManagerLoginType = 'admin' | 'mentor'