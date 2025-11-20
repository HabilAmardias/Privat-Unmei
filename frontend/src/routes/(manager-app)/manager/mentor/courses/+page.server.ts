import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch }) => {
	const [categoriesRes, coursesRes] = await Promise.all([
		controller.getCourseCategories(fetch),
		controller.getMyCourses(fetch)
	]);
	if (!categoriesRes.success) {
		throw error(categoriesRes.status, { message: categoriesRes.message });
	}
	if (!coursesRes.success) {
		throw error(coursesRes.status, { message: coursesRes.message });
	}
	return {
		courses: coursesRes.resBody.data,
		categories: categoriesRes.resBody.data.entries
	};
};

export const actions = {
	deleteCourse: async ({ fetch, request }) => {
		const { success, message, status } = await controller.deleteCourse(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { success, message, status };
	},
	getMyCourses: async ({ fetch, request }) => {
		const { success, message, status, resBody } = await controller.getMyCourses(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { courses: resBody.data };
	},
	getCourseCategories: async ({ fetch, request }) => {
		const { success, message, status, resBody } = await controller.getCourseCategories(
			fetch,
			request
		);
		if (!success) {
			return fail(status, { message });
		}
		return { categories: resBody.data.entries };
	}
} satisfies Actions;
