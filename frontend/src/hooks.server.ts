import { redirect, type Handle, type HandleFetch } from '@sveltejs/kit';
import { IsTokenExpired } from '$lib/utils/helper';
import { PUBLIC_BASE_URL, PUBLIC_ENVIRONMENT_OPTION } from '$env/static/public';
import { controller } from './controller';
import { Production } from '$lib/utils/constants';

export const handleFetch: HandleFetch = async ({ event, request, fetch }) => {
	if (!request.url.startsWith(PUBLIC_BASE_URL)) {
		return fetch(request);
	}
	const authToken = event.cookies.get('auth_token');
	const refreshToken = event.cookies.get('refresh_token');

	if (!IsTokenExpired(refreshToken) && IsTokenExpired(authToken)) {
		const url = `${PUBLIC_BASE_URL}/api/v1/refresh`;
		const res = await fetch(url, {
			method: 'GET',
			headers: {
				...Object.fromEntries(request.headers),
				Cookie: `auth_token=${authToken}; refresh_token=${refreshToken}`
			},
			credentials: 'include'
		});
		if (res.ok) {
			const cookiesData = controller.GetCookies(res);
			cookiesData.forEach((val) => {
				event.cookies.set(val.key, val.value, {
					path: val.path,
					httpOnly: val.httpOnly,
					domain: val.domain,
					maxAge: val.maxAge,
					sameSite: val.sameSite,
					secure: PUBLIC_ENVIRONMENT_OPTION === Production
				});
			});
		}
	}
	const requestCookies: Array<string> = [];
	event.cookies.getAll().forEach((c) => {
		requestCookies.push(`${c.name}=${c.value}`);
	});
	const newReq = new Request(request, {
		headers: {
			...Object.fromEntries(request.headers),
			Cookie: requestCookies.join(';')
		},
		credentials: 'include',
		redirect: 'manual',
		duplex: 'half'
	});
	return fetch(newReq);
};

export const handle: Handle = async ({ event, resolve }) => {
	const authToken = event.cookies.get('auth_token');
	const refreshToken = event.cookies.get('refresh_token');
	const isPublic = event.route.id?.includes('/(public)');
	if (!isPublic) {
		if (IsTokenExpired(authToken) && IsTokenExpired(refreshToken)) {
			throw redirect(
				303,
				event.route.id?.includes('/(manager-app)') ? '/manager/logout' : '/logout'
			);
		}
	}

	return resolve(event);
};
