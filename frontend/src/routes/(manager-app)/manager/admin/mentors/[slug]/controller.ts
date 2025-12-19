import type { Fetch, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { MentorPaymentInfo, MentorProfile, MentorScheduleInfo, adminProfile } from './model';

class MentorProfileController {
	async getMentorProfile(fetch: Fetch, id: string) {
		const url = `http://localhost:8080/api/v1/mentors/${id}`;
		const { success, res, status, message } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, status, message };
		}
		const resBody: ServerResponse<MentorProfile> = await res?.json();
		return { success, status, message, resBody };
	}
	async getMentorSchedules(fetch: Fetch, id: string) {
		const url = `http://localhost:8080/api/v1/mentors/${id}/availability`;
		const { success, res, status, message } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, status, message };
		}
		const resBody: ServerResponse<MentorScheduleInfo[]> = await res?.json();
		return { success, status, message, resBody };
	}
	async getMentorPayments(fetch: Fetch, id: string) {
		const url = `http://localhost:8080/api/v1/mentors/${id}/payment-methods`;
		const { success, res, status, message } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, status, message };
		}
		const resBody: ServerResponse<MentorPaymentInfo[]> = await res?.json();
		return { success, status, message, resBody };
	}
	async getAdminProfile(fetch: Fetch) {
		const url = 'http://localhost:8080/api/v1/admins/me';
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<adminProfile> = await res?.json();
		return { success, message, status, resBody };
	}
	async deleteMentor(fetch: Fetch, id?: string) {
		if (!id) {
			return { success: false, message: 'no mentor selected', status: 400 };
		}
		const url = `http://localhost:8080/api/v1/mentors/${id}`;
		const { success, message, status } = await FetchData(fetch, url, 'DELETE');
		return { success, message, status };
	}
}

export const controller = new MentorProfileController();
