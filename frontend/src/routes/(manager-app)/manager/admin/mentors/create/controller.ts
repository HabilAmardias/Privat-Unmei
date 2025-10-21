import type { Fetch, SeekPaginatedResponse, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { paymentMethod } from './model';

class CreateMentorController {
	async getPaymentMethods(fetch: Fetch) {
		const url = 'http://localhost:8080/api/v1/payment-methods';
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<SeekPaginatedResponse<paymentMethod>> = await res?.json();
		return { success, message, status, resBody };
	}
}

export const controller = new CreateMentorController();
