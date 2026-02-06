import { error, redirect, type Handle, type HandleFetch } from '@sveltejs/kit';
import { IsTokenExpired } from '$lib/utils/helper';
import { PUBLIC_BASE_URL, PUBLIC_ENVIRONMENT_OPTION } from '$env/static/public';
import { controller } from './controller';
import { Production, SESSION_EXPIRED } from '$lib/utils/constants';
import type { MessageResponse, ServerResponse } from '$lib/types';

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
			if (res.status !== 200 && !event.route.id?.includes('/(public)')) {
				const resBody: ServerResponse<MessageResponse> = await res.json();
				if (resBody.data.message === SESSION_EXPIRED) {
					throw redirect(
						303,
						event.route.id?.includes('/(manager-app)') ? '/manager/logout' : '/logout'
					);
				}
				throw error(res.status, { message: 'unauthorized' });
			}
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
		if (IsTokenExpired(refreshToken) && IsTokenExpired(authToken)) {
			if (!event.route.id?.includes('/(public)')) {
				throw redirect(
					303,
					event.route.id?.includes('/(manager-app)') ? '/manager/logout' : '/logout'
				);
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
		const res = await fetch(url, {
			headers: {
				Cookie: `refresh_token=${refreshToken}; auth_token=${authToken}`
			},
			credentials: 'include'
		});
		if (res.status !== 200 && !event.route.id?.includes('/(public)')) {
			const resBody: ServerResponse<MessageResponse> = await res.json();
			if (resBody.data.message === SESSION_EXPIRED) {
				throw redirect(
					303,
					event.route.id?.includes('/(manager-app)') ? '/manager/logout' : '/logout'
				);
			}
			throw error(res.status, { message: 'unauthorized' });
		}
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
	if (IsTokenExpired(refreshToken) && IsTokenExpired(authToken)) {
		if (!event.route.id?.includes('/(public)')) {
			throw redirect(
				303,
				event.route.id?.includes('/(manager-app)') ? '/manager/logout' : '/logout'
			);
		}
	}
	return resolve(event);
};
