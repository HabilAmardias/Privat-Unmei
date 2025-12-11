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
	messageMentor: async ({ fetch, request }) => {
		const { success, status, message, resBody } = await controller.messageMentor(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		throw redirect(303, `/chats/${resBody.data.id}`);
	}
} satisfies Actions;
