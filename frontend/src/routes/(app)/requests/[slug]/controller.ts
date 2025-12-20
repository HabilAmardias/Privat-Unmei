import type { Fetch, ServerResponse } from '$lib/types';
import type { RequestDetail, ChatroomID } from './model';
import { FetchData } from '$lib/utils';

class RequestDetailController {
	async getRequestDetail(fetch: Fetch, id: string) {
		const url = `http://localhost:80/api/v1/me/course-requests/${id}`;
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<RequestDetail> = await res?.json();
		return { success, message, status, resBody };
	}
	async messageMentor(fetch: Fetch, req: Request) {
		const formData = await req.formData();
		const id = formData.get('id');
		const url = `http://localhost:80/api/v1/chatrooms/users/${id}`;
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<ChatroomID> = await res?.json();
		return { resBody, status, success, message };
	}
}

export const controller = new RequestDetailController();
