import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch, params }) => {
	const { success, message, status, resBody } = await controller.getRequestDetail(
		fetch,
		params.slug
	);
	if (!success) {
		throw error(status, { message });
	}
	return {
		detail: resBody.data
	};
};
