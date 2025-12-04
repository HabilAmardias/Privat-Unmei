import type { Fetch, ServerResponse } from '$lib/types';
import type { RequestDetail } from './model';
import { FetchData } from '$lib/utils';

class RequestDetailController {
	async getRequestDetail(fetch: Fetch, id: string) {
		const url = `http://localhost:8080/api/v1/me/course-requests/${id}`;
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<RequestDetail> = await res?.json();
		return { success, message, status, resBody };
	}
}

export const controller = new RequestDetailController();
