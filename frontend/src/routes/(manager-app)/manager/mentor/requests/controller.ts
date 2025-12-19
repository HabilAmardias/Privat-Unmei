import type { Fetch, PaginatedResponse, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { CourseRequest } from './model';

class RequestManagementController {
	async getMyRequests(fetch: Fetch, req?: Request) {
		let url = 'http://habilog.xyz/api/v1/mentors/me/course-requests?';
		if (req) {
			const args: string[] = [];
			const formData = await req.formData();
			const page = formData.get('page');
			if (page) {
				args.push(`page=${page}`);
			}
			const status = formData.get('status');
			if (status) {
				args.push(`status=${status}`);
			}
			if (args.length > 0) {
				url += args.join('&');
			}
		}
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<PaginatedResponse<CourseRequest>> = await res?.json();
		return { success, message, status, resBody };
	}
}

export const controller = new RequestManagementController();
