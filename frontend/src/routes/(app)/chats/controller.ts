import type { Fetch, PaginatedResponse, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { Chatroom } from './model';

class ChatController {
	getChatrooms = async (fetch: Fetch, req?: Request) => {
		let url = 'http://habilog.xyz/api/v1/chatrooms/me?';
		if (req) {
			const formData = await req.formData();
			const page = formData.get('page');
			if (page) {
				url += `page=${page}`;
			}
		}
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}

		const resBody: ServerResponse<PaginatedResponse<Chatroom>> = await res?.json();
		return { success, message, status, resBody };
	};
}

export const controller = new ChatController();
