import { error, type Handle, type HandleFetch } from '@sveltejs/kit';
import { IsTokenExpired } from '$lib/utils/helper';
import { PUBLIC_BASE_URL, PUBLIC_ENVIRONMENT_OPTION } from '$env/static/public';
import { controller } from './controller';
import { Production } from '$lib/utils/constants';

export const handleFetch: HandleFetch = async ({ event, request, fetch }) => {
	if (request.url.startsWith(PUBLIC_BASE_URL)) {
		const authToken = event.cookies.get('auth_token');
		const refreshToken = event.cookies.get('refresh_token');
		if (!IsTokenExpired(refreshToken) && IsTokenExpired(authToken)) {
			const url = `${PUBLIC_BASE_URL}/api/v1/refresh`;
			const res = await fetch(url, {
				headers: {
					Cookie: `refresh_token=${refreshToken}; auth_token=${authToken}`
				},
				credentials: 'include'
			});
			if (res.status !== 200) {
				if (!event.route.id?.includes('/(public)')) {
					throw error(res.status, { message: 'unauthorized' });
				}
			} else {
				const cookiesData = controller.GetCookies(res);
				cookiesData.forEach((val) => {
					event.cookies.set(val.key, val.value, {
						path: val.path,
						httpOnly: val.httpOnly,
						maxAge: val.maxAge,
						sameSite: val.sameSite,
						secure: PUBLIC_ENVIRONMENT_OPTION === Production
					});
				});
			}
		}
		if (IsTokenExpired(refreshToken) && IsTokenExpired(authToken)) {
			if (!event.route.id?.includes('/(public)')) {
				throw error(401, { message: 'unauthorized' });
			}
		}
		const newAuthToken = event.cookies.get('auth_token') || authToken;
		const newRefreshToken = event.cookies.get('refresh_token') || refreshToken;
		const newReq = new Request(request.url, {
			method: request.method,
			headers: {
				...Object.fromEntries(request.headers),
				Cookie: `auth_token=${newAuthToken}; refresh_token=${newRefreshToken}`
			},
			body: request.body,
			credentials: 'include',
			redirect: 'manual',
			duplex: 'half'
		});
		return fetch(newReq);
	}
	return fetch(request);
};

export const handle: Handle = async ({ event, resolve }) => {
	const authToken = event.cookies.get('auth_token');
	const refreshToken = event.cookies.get('refresh_token');

	if (!IsTokenExpired(refreshToken) && IsTokenExpired(authToken)) {
		const url = `${PUBLIC_BASE_URL}/api/v1/refresh`;
		const res = await event.fetch(url);
		if (res.status !== 200) {
			if (!event.route.id?.includes('/(public)')) {
				throw error(res.status, { message: 'unauthorized' });
			}
		} else {
			const cookiesData = controller.GetCookies(res);
			cookiesData.forEach((val) => {
				event.cookies.set(val.key, val.value, {
					path: val.path,
					httpOnly: val.httpOnly,
					maxAge: val.maxAge,
					sameSite: val.sameSite,
					secure: PUBLIC_ENVIRONMENT_OPTION === Production
				});
			});
		}
	}
	if (IsTokenExpired(refreshToken) && IsTokenExpired(authToken)) {
		if (!event.route.id?.includes('/(public)')) {
			throw error(401, { message: 'unauthorized' });
		}
	}
	return resolve(event);
};
