import { redirect, type HandleFetch } from '@sveltejs/kit';
import { controller } from './controller';

const publicRoutes = ['/', '/reset', '/home', '/courses', '/playground'];

export const handleFetch: HandleFetch = async ({ event, request, fetch }) => {
	const authToken = event.cookies.get('auth_token');
	const refreshToken = event.cookies.get('refresh_token');

	if (!publicRoutes.includes(event.url.pathname)) {
		if (refreshToken && !authToken) {
			const { success, cookiesData } = await controller.refresh(fetch);
			if (!success) {
				redirect(303, '/');
			}
			cookiesData?.forEach((val) => {
				event.cookies.set(val.key, val.value, {
					path: val.path,
					domain: val.domain,
					httpOnly: val.httpOnly,
					maxAge: val.maxAge
				});
			});
		}
		if (!refreshToken && !authToken) {
			redirect(303, '/');
		}
	}
	return fetch(request);
};
