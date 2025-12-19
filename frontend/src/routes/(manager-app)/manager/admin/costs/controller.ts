import type { Fetch, PaginatedResponse, ServerResponse } from '$lib/types';
import { FetchData } from '$lib/utils';
import type { AdditionalCost, adminProfile, Discount, newCost, newDiscount } from './model';

class CostManagementController {
	getDiscounts = async (fetch: Fetch, req?: Request) => {
		let url = '/api/v1/discounts?';
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
		const resBody: ServerResponse<PaginatedResponse<Discount>> = await res?.json();
		return { success, message, status, resBody };
	};
	updateDiscountAmount = async (fetch: Fetch, req: Request) => {
		const formData = await req.formData();
		const id = formData.get('id');
		const amount = formData.get('amount');
		if (!id) {
			return { success: false, message: 'no discount selected', status: 400 };
		}
		const reqBody = JSON.stringify({
			amount: amount ? parseFloat(amount as string) : null
		});
		const url = `/api/v1/discounts/${id}`;

		const { success, message, status } = await FetchData(fetch, url, 'PATCH', reqBody);
		return { success, message, status };
	};
	getCosts = async (fetch: Fetch, req?: Request) => {
		let url = '/api/v1/additional-costs?';
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
		const resBody: ServerResponse<PaginatedResponse<AdditionalCost>> = await res?.json();
		return { success, message, status, resBody };
	};
	deleteDiscount = async (fetch: Fetch, req: Request) => {
		const formData = await req.formData();
		const id = formData.get('id');
		if (!id) {
			return { success: false, message: 'no discount selected', status: 400 };
		}
		const url = `/api/v1/discounts/${id}`;
		const { success, message, status } = await FetchData(fetch, url, 'DELETE');
		return { success, message, status };
	};
	createDiscount = async (fetch: Fetch, req: Request) => {
		const formData = await req.formData();
		const amount = formData.get('amount');
		const number_of_participant = formData.get('number_of_participant');
		if (!amount) {
			return { success: false, message: 'provide discount amount', status: 400 };
		}
		if (!number_of_participant) {
			return { success: false, message: 'provide discount number of participant', status: 400 };
		}
		const reqBody = JSON.stringify({
			number_of_participant: parseInt(number_of_participant as string),
			amount: parseFloat(amount as string)
		});
		const url = '/api/v1/discounts';
		const { success, message, status, res } = await FetchData(fetch, url, 'POST', reqBody);
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<newDiscount> = await res?.json();
		return { success, message, status, resBody };
	};
	updateCostAmount = async (fetch: Fetch, req: Request) => {
		const formData = await req.formData();
		const id = formData.get('id');
		const amount = formData.get('amount');
		if (!id) {
			return { success: false, message: 'no costs selected', status: 400 };
		}
		const reqBody = JSON.stringify({
			amount: amount ? parseFloat(amount as string) : null
		});
		const url = `/api/v1/additional-costs/${id}`;

		const { success, message, status } = await FetchData(fetch, url, 'PATCH', reqBody);
		return { success, message, status };
	};
	deleteCost = async (fetch: Fetch, req: Request) => {
		const formData = await req.formData();
		const id = formData.get('id');
		if (!id) {
			return { success: false, message: 'no costs selected', status: 400 };
		}
		const url = `/api/v1/additional-costs/${id}`;
		const { success, message, status } = await FetchData(fetch, url, 'DELETE');
		return { success, message, status };
	};
	createCost = async (fetch: Fetch, req: Request) => {
		const formData = await req.formData();
		const amount = formData.get('amount');
		const name = formData.get('name');
		if (!amount) {
			return { success: false, message: 'provide cost amount', status: 400 };
		}
		if (!name) {
			return { success: false, message: 'provide cost name', status: 400 };
		}
		const reqBody = JSON.stringify({
			name,
			amount: parseFloat(amount as string)
		});
		const url = '/api/v1/additional-costs';
		const { success, message, status, res } = await FetchData(fetch, url, 'POST', reqBody);
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<newCost> = await res?.json();
		return { success, message, status, resBody };
	};
	getProfile = async (fetch: Fetch) => {
		const url = '/api/v1/admins/me';
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const resBody: ServerResponse<adminProfile> = await res?.json();
		return { success, message, status, resBody };
	};
}

export const controller = new CostManagementController();
