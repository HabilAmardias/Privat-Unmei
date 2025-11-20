import { error, fail, redirect, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch, params }) => {
	const [categoriesRes, courseCategoriesRes, topicsRes, courseRes] = await Promise.all([
		controller.getCourseCategories(fetch),
		controller.getCourseDetailCategories(fetch, params.slug),
		controller.getCourseTopics(fetch, params.slug),
		controller.getCourseDetail(fetch, params.slug)
	]);
	if (!categoriesRes.success) {
		throw error(categoriesRes.status, { message: categoriesRes.message });
	}
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
		categories: categoriesRes.resBody.data.entries,
		courseCategories: courseCategoriesRes.resBody.data,
		topics: topicsRes.resBody.data,
		detail: courseRes.resBody.data
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
	updateCourse: async ({ fetch, request, params }) => {
		if (!params.slug) {
			throw error(400, { message: 'no course selected' });
		}
		const { success, message, status } = await controller.updateCourse(fetch, request, params.slug);
		if (!success) {
			return fail(status, { message });
		}
		redirect(303, `/manager/mentor/courses/${params.slug}`);
	}
} satisfies Actions;
