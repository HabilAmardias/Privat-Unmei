import { error, fail, redirect, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch, params }) => {
	const [mentorProfile, adminProfile] = await Promise.all([
		controller.getMentorProfile(fetch, params.slug),
		controller.getAdminProfile(fetch)
	]);
	if (!mentorProfile.success) {
		throw error(mentorProfile.status, { message: mentorProfile.message });
	}
	if (!adminProfile.success) {
		throw error(adminProfile.status, { message: adminProfile.message });
	}
	return {
		isVerified: adminProfile.resBody.data.status === 'verified',
		profile: mentorProfile.resBody.data
	};
};

export const actions = {
	deleteMentor: async ({ fetch, params }) => {
		const { success, message, status } = await controller.deleteMentor(fetch, params.slug);
		if (!success) {
			return fail(status, { message });
		}
		throw redirect(303, '/manager/admin/mentors');
	}
} satisfies Actions;
