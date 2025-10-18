import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch }) => {
	const [mentorsList, adminProfile] = await Promise.all([
		controller.getMentors(fetch),
		controller.getProfile(fetch)
	]);
	if (!mentorsList.success) {
		throw error(mentorsList.status, { message: mentorsList.message });
	}
	if (!adminProfile.success) {
		throw error(adminProfile.status, { message: adminProfile.message });
	}
	return {
		success: true,
		mentorsList: mentorsList.resBody.data,
		isVerified: adminProfile.resBody.data.status === 'verified'
	};
};

export const actions = {
	getMentors: async ({ fetch, request }) => {
		const { success, message, status, resBody } = await controller.getMentors(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { mentorsList: resBody.data };
	},
	deleteMentor: async ({ fetch, request }) => {
		const { success, message, status } = await controller.deleteMentor(fetch, request);
		if (!success) {
			return fail(status, { message });
		}
		return { success, message, status };
	}
} satisfies Actions;
