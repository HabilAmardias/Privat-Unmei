import type { Fetch } from '$lib/types';
import { FetchData } from '$lib/utils';

class ResetController {
	async sendEmail(req: Request, fetch: Fetch) {
		const url = 'http://localhost:80/api/v1/reset-password/send';
		const formData = await req.formData();
		const email = formData.get('email');
		if (!email) {
			return { success: false, message: 'please insert an email', status: 400 };
		}
		if (!this.#validateEmail(email as string)) {
			return { success: false, message: 'please insert an valid email', status: 400 };
		}
		const body = JSON.stringify({
			email
		});
		const { success, status, message } = await FetchData(fetch, url, 'POST', body);
		return { success, status, message };
	}
	async resetPassword(req: Request, fetch: Fetch) {
		const url = 'http://localhost:80/api/v1/reset-password/reset';
		const formData = await req.formData();
		const password = formData.get('password');
		const repeatPassword = formData.get('repeat-password');
		if (!password) {
			return { success: false, message: 'please insert a password', status: 400 };
		}
		if (!this.#validatePassword(password as string)) {
			return { success: false, message: 'please insert a valid password', status: 400 };
		}
		if (password !== repeatPassword) {
			return { success: false, message: 'password does not match', status: 400 };
		}
		const body = JSON.stringify({
			new_password: password
		});
		const { success, message, status } = await FetchData(fetch, url, 'POST', body);
		return { success, message, status };
	}
	#validateEmail(email: string) {
		const pattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		return pattern.test(email);
	}
	#validatePassword(password: string) {
		const minLength = password.length >= 8;
		const hasSpecialChar =
			password.includes('!') ||
			password.includes('@') ||
			password.includes('#') ||
			password.includes('?');
		return minLength && hasSpecialChar;
	}
}

export const controller = new ResetController();
