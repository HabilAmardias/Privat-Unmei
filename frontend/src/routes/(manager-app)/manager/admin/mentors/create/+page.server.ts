import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch }) => {
	const { success, status, message, resBody } = await controller.getPaymentMethods(fetch);
	if (!success) {
		throw error(status, { message });
	}
	return { paymentMethods: resBody.data.entries };
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
	}
} satisfies Actions;
