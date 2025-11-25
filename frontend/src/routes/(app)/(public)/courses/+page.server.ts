import { error, fail, type Actions } from '@sveltejs/kit';
import { controller } from './controller';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	const [mostBought, courses] = await Promise.all([
		controller.getMostBoughtCourses(fetch),
		controller.getCourses(fetch)
	]);
	if (!mostBought.success) {
		throw error(mostBought.status, { message: mostBought.message });
	}
	if (!courses.success) {
		throw error(courses.status, { message: courses.message });
	}
	return { mostBought: mostBought.resBody.data, courses: courses.resBody.data };
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
