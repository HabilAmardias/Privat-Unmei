import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch }) => {
	const { success, message, status, resBody } = await controller.getChatrooms(fetch);
	if (!success) {
		throw error(status, { message });
	}
	return {
		chatrooms: resBody.data
	};
};

export const actions = {
	getChats: async ({ fetch, request }) => {
		const { success, message, status, resBody } = await controller.getChatrooms(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return {
			chatrooms: resBody.data
		};
	}
} satisfies Actions;
