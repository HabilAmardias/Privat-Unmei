import { error, fail, redirect, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch, params }) => {
	const { success, message, status, resBody } = await controller.getRequestDetail(
		fetch,
		params.slug
	);
	if (!success) {
		throw error(status, { message });
	}
	return {
		detail: resBody.data
	};
};

export const actions = {
	acceptRequest: async ({ fetch, params }) => {
		if (!params.slug) {
			throw error(500, 'something went wrong');
		}
		const { success, message, status } = await controller.acceptRequest(fetch, params.slug);
		if (!success) {
			return fail(status, { message });
		}
		redirect(303, '/manager/mentor/requests');
	},
	rejectRequest: async ({ fetch, params }) => {
		if (!params.slug) {
			throw error(500, 'something went wrong');
		}
		const { success, message, status } = await controller.rejectRequest(fetch, params.slug);
		if (!success) {
			return fail(status, { message });
		}
		redirect(303, '/manager/mentor/requests');
	},
	confirmPayment: async ({ fetch, params }) => {
		if (!params.slug) {
			throw error(500, 'something went wrong');
		}
		const { success, message, status } = await controller.confirmPayment(fetch, params.slug);
		if (!success) {
			return fail(status, { message });
		}
		redirect(303, '/manager/mentor/requests');
	}
} satisfies Actions;
