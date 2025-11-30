import { error, fail, type Actions } from '@sveltejs/kit';
import { controller } from './controller';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	const { success, message, status, resBody } = await controller.getMentors(fetch);
	if (!success) {
		throw error(status, { message });
	}
	return { mentors: resBody.data };
};

export const actions = {
	getMentors: async ({ fetch, request }) => {
		const { success, message, status, resBody } = await controller.getMentors(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return {
			mentors: resBody.data
		};
	}
} satisfies Actions;
