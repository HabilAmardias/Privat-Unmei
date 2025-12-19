import type { Fetch, PaginatedResponse, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { generatedPassword, paymentMethod, adminProfile } from './model';

class CreateMentorController {
	async CreateNewMentor(fetch: Fetch, req: Request) {
		const url = 'http://habilog.xyz/api/v1/mentors';
		const formData = await req.formData();
		const { success, message, status } = await FetchData(fetch, url, 'POST', formData);
		return { success, message, status };
	}
	async getPaymentMethods(fetch: Fetch, req?: Request) {
		let url = 'http://habilog.xyz/api/v1/payment-methods?limit=5';
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
	async getRandomizedPassword(fetch: Fetch) {
		const url = 'http://habilog.xyz/api/v1/mentors/password';
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<generatedPassword> = await res?.json();
		return { success, message, status, resBody };
	}
	async getAdminProfile(fetch: Fetch) {
		const url = 'http://habilog.xyz/api/v1/admins/me';
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<adminProfile> = await res?.json();
		return { success, message, status, resBody };
	}
}

export const controller = new CreateMentorController();
