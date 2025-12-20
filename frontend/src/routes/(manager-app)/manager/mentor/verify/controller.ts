import type { Fetch } from '$lib/types';
import { FetchData } from '$lib/utils';

class AdminVerifyController {
	async verifyMentor(fetch: Fetch, req: Request) {
		const formData = await req.formData();
		const newPassword = formData.get('password');
		if (!newPassword) {
			return { success: false, status: 400, message: 'please insert an password' };
		}
		if (!this.#validatePassword(newPassword as string)) {
			return { success: false, status: 400, message: 'please insert a valid password' };
		}
		const reqBody = JSON.stringify({
			password: newPassword
		});
		const url = 'http://localhost:80/api/v1/mentors/me/change-password';
		const { success, message, status } = await FetchData(fetch, url, 'POST', reqBody);
		return { success, message, status };
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

export const controller = new AdminVerifyController();
