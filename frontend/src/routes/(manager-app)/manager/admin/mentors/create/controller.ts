import type { Fetch, SeekPaginatedResponse, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { paymentMethod } from './model';

class CreateMentorController {
	async getPaymentMethods(fetch: Fetch, req?: Request) {
		let url = 'http://localhost:8080/api/v1/payment-methods?';
		if (req) {
			const formData = await req.formData();
			const search = formData.get('search');
			const limit = formData.get('limit');
			if (search) {
				url += `search=${search}`;
			}
			if (limit) {
				url += `limit=${limit}`;
			}
		}
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<SeekPaginatedResponse<paymentMethod>> = await res?.json();
		return { success, message, status, resBody };
	}
}

export const controller = new CreateMentorController();
