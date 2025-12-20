import type { Fetch, PaginatedResponse, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { MentorScheduleInfo, paymentMethod, MentorPaymentInfo, MentorProfile } from './model';

class UpdateMentorController {
	async updateProfile(fetch: Fetch, req: Request) {
		const url = '/api/v1/mentors/me';
		const formData = await req.formData();
		const { success, message, status } = await FetchData(fetch, url, 'PATCH', formData);
		return { success, message, status };
	}
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
	async getPaymentMethods(fetch: Fetch, req?: Request) {
		let url = '/api/v1/payment-methods?limit=5';
		if (req) {
			const formData = await req.formData();
			const search = formData.get('search');
			if (search) {
				url += `&search=${search}`;
			}
		}
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<PaginatedResponse<paymentMethod>> = await res?.json();
		return { success, message, status, resBody };
	}
}

export const controller = new UpdateMentorController();
