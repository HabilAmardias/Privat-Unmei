import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch }) => {
	const res = await controller.getPaymentMethods(fetch);
	if (!res.success) {
		throw error(res.status, { message: res.message });
	}
	return {
		paymentMethods: res.resBody.data.entries
	};
};

export const actions = {
	getPaymentMethods: async ({ fetch, request }) => {
		const { success, status, message, resBody } = await controller.getPaymentMethods(
			fetch,
			request
		);
		if (!success) {
			throw fail(status, { message });
		}
		return { paymentMethods: resBody.data.entries };
	},
	updateProfile: async ({ fetch, request }) => {
		const { success, status, message } = await controller.updateProfile(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { message };
	}
} satisfies Actions;
