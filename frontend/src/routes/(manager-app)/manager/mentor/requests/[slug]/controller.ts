import type { Fetch, ServerResponse } from '$lib/types';
import type { RequestDetail } from './model';
import { FetchData } from '$lib/utils';

class CourseDetailController {
	async getRequestDetail(fetch: Fetch, id: string) {
		const url = `http://localhost:8080/api/v1/mentors/me/course-requests/${id}`;
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<RequestDetail> = await res?.json();
		return { success, message, status, resBody };
	}
	async acceptRequest(fetch: Fetch, id: string) {
		const url = `http://localhost:8080/api/v1/course-requests/${id}/approve`;
		const { success, message, status } = await FetchData(fetch, url, 'GET');
		return { success, message, status };
	}
	async rejectRequest(fetch: Fetch, id: string) {
		const url = `http://localhost:8080/api/v1/course-requests/${id}/reject`;
		const { success, message, status } = await FetchData(fetch, url, 'GET');
		return { success, message, status };
	}
	async confirmPayment(fetch: Fetch, id: string) {
		const url = `http://localhost:8080/api/v1/course-requests/${id}/confirm-payment`;
		const { success, message, status } = await FetchData(fetch, url, 'GET');
		return { success, message, status };
	}
}

export const controller = new CourseDetailController();
