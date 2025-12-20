import type { Fetch, PaginatedResponse, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { NewPaymentMethod, PaymentMethods, adminProfile } from './model';

class PaymentController {
	getPayments = async (fetch: Fetch, req?: Request) => {
		let url = 'http://localhost:80/api/v1/payment-methods?';
		if (req) {
			const args: string[] = [];
			const formData = await req.formData();
			const search = formData.get('search');
			const page = formData.get('page');
			if (search) {
				args.push(`search=${search}`);
			}
			if (page) {
				args.push(`page=${page}`);
			}
			const joinedArgs = args.join('&');
			url += joinedArgs;
		}
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<PaginatedResponse<PaymentMethods>> = await res?.json();
		return { success, message, status, resBody };
	};
	updatePayment = async (fetch: Fetch, req: Request) => {
		const formData = await req.formData();
		const id = formData.get('id');
		const url = `http://localhost:80/api/v1/payment-methods/${id}`;
		const name = formData.get('name');
		const reqBody = JSON.stringify({
			payment_method_name: name
		});
		const { success, message, status } = await FetchData(fetch, url, 'PATCH', reqBody);
		return { success, message, status };
	};
	deletePayment = async (fetch: Fetch, req: Request) => {
		const formData = await req.formData();
		const id = formData.get('id');
		const url = `http://localhost:80/api/v1/payment-methods/${id}`;
		const { success, message, status } = await FetchData(fetch, url, 'DELETE');
		return { success, message, status };
	};
	createPayment = async (fetch: Fetch, req: Request) => {
		const formData = await req.formData();
		const url = `http://localhost:80/api/v1/payment-methods`;
		const name = formData.get('name');
		const reqBody = JSON.stringify({
			payment_method_name: name
		});
		const { success, message, status, res } = await FetchData(fetch, url, 'POST', reqBody);
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<NewPaymentMethod> = await res?.json();
		return { success, message, status, resBody };
	};
	async getProfile(fetch: Fetch) {
		const url = 'http://localhost:80/api/v1/admins/me';
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<adminProfile> = await res?.json();
		return { success, message, status, resBody };
	}
}

export const controller = new PaymentController();
