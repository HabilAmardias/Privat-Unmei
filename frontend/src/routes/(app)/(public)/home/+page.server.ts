import { error } from '@sveltejs/kit';
import { controller } from './controller';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	const { success, message, status, resBody } = await controller.getCourses(fetch);
	if (!success) {
		throw error(status, { message });
	}
	return {
		courses: resBody.data
	};
};
