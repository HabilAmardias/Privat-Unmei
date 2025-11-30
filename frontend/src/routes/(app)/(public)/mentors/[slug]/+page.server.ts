import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch, params }) => {
	const [mentorProfile, mentorSchedules, mentorCourses] = await Promise.all([
		controller.getMentorProfile(fetch, params.slug),
		controller.getMentorSchedules(fetch, params.slug),
		controller.getMentorCourses(fetch, params.slug)
	]);
	if (!mentorSchedules.success) {
		throw error(mentorSchedules.status, { message: mentorSchedules.message });
	}
	if (!mentorProfile.success) {
		throw error(mentorProfile.status, { message: mentorProfile.message });
	}
	if (!mentorCourses.success) {
		throw error(mentorCourses.status, { message: mentorCourses.message });
	}
	return {
		profile: mentorProfile.resBody.data,
		schedules: mentorSchedules.resBody?.data,
		courses: mentorCourses.resBody.data
	};
};

export const actions = {
	getCourses: async ({ fetch, request, params }) => {
		if (!params.slug) {
			throw error(404, { message: 'no data found' });
		}
		const { success, status, message, resBody } = await controller.getMentorCourses(
			fetch,
			params.slug,
			request
		);
		if (!success) {
			return fail(status, { message });
		}
		return {
			courses: resBody.data
		};
	}
} satisfies Actions;
