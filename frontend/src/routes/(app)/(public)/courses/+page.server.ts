import { error, fail, type Actions } from '@sveltejs/kit';
import { controller } from './controller';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	const { success, message, status, resBody } = await controller.getCourses(fetch);
	if (!success) {
		throw error(status, { message });
	}
	return { courses: resBody.data };
};

export const actions = {
	getCourses: async ({ fetch, request }) => {
		const { success, message, status, resBody } = await controller.getCourses(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return {
			courses: resBody.data
		};
	},
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
	}
} satisfies Actions;
