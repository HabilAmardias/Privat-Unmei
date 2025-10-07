import { error, fail, type Actions } from '@sveltejs/kit';
import { controller } from './controller';
import type { PageServerLoad } from './$types';

export const load : PageServerLoad = async ({fetch}) =>{
	const [mostBought] = await Promise.all([controller.getMostBoughtCourses(fetch)])
	if (!mostBought.success){
		error(mostBought.status, {message: mostBought.message})
	}
	return mostBought
}

export const actions = {
	refresh: async ({ cookies, fetch }) => {
		if (!cookies.get('refresh_token')) {
			return fail(401, { message: 'unauthorized' });
		}
		const { success, cookiesData, message, status } = await controller.refresh(fetch);
		if (!success) {
			return fail(status, { message });
		}
		cookiesData?.forEach((val) => {
			cookies.set(val.key, val.value, {
				path: val.path,
				domain: val.domain,
				httpOnly: val.httpOnly,
				maxAge: val.maxAge,
				sameSite: val.sameSite
			});
		});
		return { success: true, message };
	}
} satisfies Actions;
