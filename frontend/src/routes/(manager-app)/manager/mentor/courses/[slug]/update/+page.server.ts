import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch }) => {
	const { success, message, status, resBody } = await controller.getCourseCategories(fetch);
	if (!success) {
		throw error(status, { message });
	}
	return {
		categories: resBody.data.entries
	};
};

export const actions = {
	getCategories: async ({ fetch, request }) => {
		const { success, message, status, resBody } = await controller.getCourseCategories(
			fetch,
			request
		);
		if (!success) {
			return fail(status, { message });
		}
		return {
			categories: resBody.data.entries
		};
	},
	updateCourse: async ({ fetch, request }) => {
		const { success, message, status } = await controller.updateCourse(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { success };
	}
} satisfies Actions;
