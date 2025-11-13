import type { Fetch, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { MentorPaymentInfo, MentorProfile, MentorScheduleInfo } from './model';

class MentorPageController {
	async getMentorProfile(fetch: Fetch) {
		const url = `http://localhost:8080/api/v1/mentors/me`;
		const { success, res, status, message } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, status, message };
		}
		const resBody: ServerResponse<MentorProfile> = await res?.json();
		return { success, status, message, resBody };
	}
	async getMentorSchedules(fetch: Fetch) {
		const url = `http://localhost:8080/api/v1/mentors/me/availability`;
		const { success, res, status, message } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, status, message };
		}
		const resBody: ServerResponse<MentorScheduleInfo[]> = await res?.json();
		return { success, status, message, resBody };
	}
	async getMentorPayments(fetch: Fetch) {
		const url = `http://localhost:8080/api/v1/mentors/me/payment-methods`;
		const { success, res, status, message } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, status, message };
		}
		const resBody: ServerResponse<MentorPaymentInfo[]> = await res?.json();
		return { success, status, message, resBody };
	}
}

export const controller = new MentorPageController();
