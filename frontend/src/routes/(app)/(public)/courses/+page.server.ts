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
	return { mostBought: mostBought.resBody, courses: courses.resBody };
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
	}
} satisfies Actions;
