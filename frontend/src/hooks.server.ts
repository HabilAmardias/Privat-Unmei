import { type Handle, type HandleFetch } from '@sveltejs/kit';
import { controller } from './controller';

export const handleFetch: HandleFetch = async ({ event, request, fetch }) => {
	const authToken = event.cookies.get('auth_token');
	const refreshToken = event.cookies.get('refresh_token');

	if (!event.route.id?.includes('/(public)')) {
		if (refreshToken && !authToken) {
			const { success, cookiesData, message, status } = await controller.refresh(fetch);
			if (!success) {
				const body = JSON.stringify({
					success,
					data: {
						message
					}
				});
				return new Response(body, { status });
			}
			cookiesData?.forEach((val) => {
				event.cookies.set(val.key, val.value, {
					path: val.path,
					domain: val.domain,
					httpOnly: val.httpOnly,
					maxAge: val.maxAge,
					sameSite: val.sameSite
				});
			});
		}
		if (!refreshToken && !authToken) {
			const body = JSON.stringify({
				success: false,
				data: {
					message: 'unauthorized'
				}
			});
			return new Response(body, { status: 401 });
		}
	}
	return fetch(request);
};

export const handle: Handle = async({event, resolve}) =>{
	const authToken = event.cookies.get('auth_token');
	const refreshToken = event.cookies.get('refresh_token');
	if (!event.route.id?.includes('/(public)')) {
		if (refreshToken && !authToken) {
			const { success, cookiesData, message, status } = await controller.refresh(fetch);
			if (!success) {
				const body = JSON.stringify({
					success,
					data: {
						message
					}
				});
				return new Response(body, { status });
			}
			cookiesData?.forEach((val) => {
				event.cookies.set(val.key, val.value, {
					path: val.path,
					domain: val.domain,
					httpOnly: val.httpOnly,
					maxAge: val.maxAge,
					sameSite: val.sameSite
				});
			});
		}
		if (!refreshToken && !authToken) {
			const body = JSON.stringify({
				success: false,
				data: {
					message: 'unauthorized'
				}
			});
			return new Response(body, { status: 401 });
		}
	}
	return resolve(event)
}
