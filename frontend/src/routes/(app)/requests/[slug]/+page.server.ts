import { error, fail, redirect, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch, params }) => {
	const requestDetailRes = await controller.getRequestDetail(fetch, params.slug);
	if (!requestDetailRes.success) {
		throw error(requestDetailRes.status, { message: requestDetailRes.message });
	}
	const isReviewedRes = await controller.isCourseReviewed(
		fetch,
		requestDetailRes.resBody.data.course_id
	);
	if (!isReviewedRes.success) {
		throw error(isReviewedRes.status, { message: isReviewedRes.message });
	}
	return {
		detail: requestDetailRes.resBody.data,
		isReviewed: isReviewedRes.resBody.data.is_reviewed
	};
};

export const actions = {
	messageMentor: async ({ fetch, request }) => {
		const { success, status, message, resBody } = await controller.messageMentor(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		throw redirect(303, `/chats/${resBody.data.id}`);
	},
	createReview: async ({ fetch, request, params }) => {
		if (!params.slug) {
			return fail(404, { message: 'course not found' });
		}
		const { success, message, status, resBody } = await controller.createReview(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return {
			id: resBody.data.id,
			course_id: parseInt(params.slug)
		};
	}
} satisfies Actions;
