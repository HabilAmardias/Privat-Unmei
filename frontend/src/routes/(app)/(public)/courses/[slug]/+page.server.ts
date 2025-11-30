import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch, params }) => {
	const [courseCategoriesRes, topicsRes, courseRes, reviewRes] = await Promise.all([
		controller.getCourseDetailCategories(fetch, params.slug),
		controller.getCourseTopics(fetch, params.slug),
		controller.getCourseDetail(fetch, params.slug),
		controller.getCourseReviews(fetch, params.slug)
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
	if (!reviewRes.success) {
		throw error(reviewRes.status, { message: reviewRes.message });
	}
	return {
		courseCategories: courseCategoriesRes.resBody.data,
		topics: topicsRes.resBody.data,
		detail: courseRes.resBody.data,
		reviews: reviewRes.resBody.data
	};
};

export const actions = {
	getReviews: async ({ fetch, request, params }) => {
		if (!params.slug) {
			return fail(404, { message: 'course not found' });
		}
		const { success, message, status, resBody } = await controller.getCourseReviews(
			fetch,
			params.slug,
			request
		);
		if (!success) {
			return fail(status, { message });
		}
		return {
			reviews: resBody.data
		};
	}
} satisfies Actions;
