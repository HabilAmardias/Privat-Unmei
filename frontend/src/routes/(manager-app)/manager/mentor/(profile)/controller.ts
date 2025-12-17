import type { Fetch, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { MentorPaymentInfo, MentorProfile, MentorScheduleInfo } from './model';

class MentorPageController {
	async getMentorProfile(fetch: Fetch) {
		const url = `/api/v1/mentors/me`;
		const { success, res, status, message } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, status, message };
		}
		const resBody: ServerResponse<MentorProfile> = await res?.json();
		return { success, status, message, resBody };
	}
	async getMentorSchedules(fetch: Fetch) {
		const url = `/api/v1/mentors/me/availability`;
		const { success, res, status, message } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, status, message };
		}
		const resBody: ServerResponse<MentorScheduleInfo[]> = await res?.json();
		return { success, status, message, resBody };
	}
	async getMentorPayments(fetch: Fetch) {
		const url = `/api/v1/mentors/me/payment-methods`;
		const { success, res, status, message } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, status, message };
		}
		const resBody: ServerResponse<MentorPaymentInfo[]> = await res?.json();
		return { success, status, message, resBody };
	}
	async changePassword(fetch: Fetch, req: Request) {
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
		const url = '/api/v1/mentors/me/change-password';
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

export const controller = new MentorPageController();
