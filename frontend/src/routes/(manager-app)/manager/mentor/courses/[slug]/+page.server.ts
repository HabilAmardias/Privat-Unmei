import { error, fail, redirect, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch, params }) => {
	const [courseCategoriesRes, topicsRes, courseRes] = await Promise.all([
		controller.getCourseDetailCategories(fetch, params.slug),
		controller.getCourseTopics(fetch, params.slug),
		controller.getCourseDetail(fetch, params.slug)
	]);
	if (!courseCategoriesRes.success) {
		throw error(courseCategoriesRes.status, { message: courseCategoriesRes.message });
	}
	if (!topicsRes.success) {
		throw error(topicsRes.status, { message: topicsRes.message });
	}
	if (!courseRes.success) {
		throw error(courseRes.status, { message: courseRes.message });
	}
	return {
		courseCategories: courseCategoriesRes.resBody.data,
		topics: topicsRes.resBody.data,
		detail: courseRes.resBody.data
	};
};

export const actions = {
	deleteCourse: async ({ fetch, params }) => {
		if (!params.slug) {
			throw error(500, 'something went wrong');
		}
		const { success, message, status } = await controller.deleteCourse(fetch, params.slug);
		if (!success) {
			return fail(status, { message });
		}
		redirect(303, '/manager/mentor/courses');
	}
} satisfies Actions;
