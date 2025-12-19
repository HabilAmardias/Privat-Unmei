import type { Fetch, MessageResponse, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';

class VerifyController {
	async verify(fetch: Fetch) {
		const url = '/api/v1/verify';
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<MessageResponse> = await res?.json();
		return { success, message: resBody.data.message, status };
	}
}

export const controller = new VerifyController();
