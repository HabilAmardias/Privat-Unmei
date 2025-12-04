import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch }) => {
	const { success, message, status, resBody } = await controller.getMyRequests(fetch);
	if (!success) {
		throw error(status, { message });
	}
	return {
		requests: resBody.data
	};
};

export const actions = {
	getRequests: async ({ fetch, request }) => {
		const { success, message, status, resBody } = await controller.getMyRequests(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { requests: resBody.data };
	}
} satisfies Actions;
