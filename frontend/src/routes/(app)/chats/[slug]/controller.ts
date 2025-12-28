import type { Fetch, SeekPaginatedResponse, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { ChatroomInfo, MessageInfo, StudentProfile } from './model';

class ChatroomController {
	async sendMessage(fetch: Fetch, id: string, req: Request) {
		const url = `/api/v1/chatrooms/${id}/messages`;
		const formData = await req.formData();
		const content = formData.get('message');
		if (!content) {
			return { success: false, message: 'message cannot empty', status: 400 };
		}
		if ((content as string).length > 180) {
			return { success: false, message: 'message must less than 180 characters', status: 400 };
		}
		const reqBody = JSON.stringify({
			content
		});
		const { success, message, status } = await FetchData(fetch, url, 'POST', reqBody);
		return { success, message, status };
	}
	async updateLastRead(fetch: Fetch, id: string) {
		const url = `/api/v1/chatrooms/me/${id}/last-read`;
		const { success, message, status } = await FetchData(fetch, url, 'GET');
		return { success, message, status };
	}
	async getInfo(fetch: Fetch, id: string) {
		const url = `/api/v1/chatrooms/${id}`;
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}

		const resBody: ServerResponse<ChatroomInfo> = await res?.json();

		return { success, message, status, resBody };
	}
	async getMessages(fetch: Fetch, id: string, req?: Request) {
		let url = `/api/v1/chatrooms/${id}/messages?`;
		if (req) {
			const formData = await req.formData();
			const lastID = formData.get('last_id');
			if (lastID) {
				url += `last_id=${lastID}`;
			}
		}
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<SeekPaginatedResponse<MessageInfo>> = await res?.json();
		return { success, message, status, resBody };
	}
	async getProfile(fetch: Fetch) {
		const url = '/api/v1/me';
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<StudentProfile> = await res?.json();
		return { resBody, status, success, message };
	}
}

export const controller = new ChatroomController();
