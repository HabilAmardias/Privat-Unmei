import { error, type Handle, type HandleFetch } from '@sveltejs/kit';
import { controller } from './controller';
import { IsTokenExpired } from '$lib/utils/helper';

export const handleFetch: HandleFetch = async ({ event, request, fetch }) => {
	const authToken = event.cookies.get('auth_token');
	const refreshToken = event.cookies.get('refresh_token');
	if (!IsTokenExpired(refreshToken) && IsTokenExpired(authToken)) {
		const { success, cookiesData, message, status } = await controller.refresh(fetch);
		if (!success) {
			if (!event.route.id?.includes('/(public)')) {
				throw error(status, { message });
			}
		} else {
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
	}
	if (IsTokenExpired(refreshToken) && IsTokenExpired(authToken)) {
		if (!event.route.id?.includes('/(public)')) {
			throw error(401, { message: 'unauthorized' });
		}
	}

	return fetch(request);
};

export const handle: Handle = async ({ event, resolve }) => {
	const authToken = event.cookies.get('auth_token');
	const refreshToken = event.cookies.get('refresh_token');
	if (!IsTokenExpired(refreshToken) && IsTokenExpired(authToken)) {
		const { success, cookiesData, message, status } = await controller.refresh(event.fetch);
		if (!success) {
			if (!event.route.id?.includes('/(public)')) {
				throw error(status, { message });
			}
		} else {
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
	}
	if (IsTokenExpired(refreshToken) && IsTokenExpired(authToken)) {
		if (!event.route.id?.includes('/(public)')) {
			throw error(401, { message: 'unauthorized' });
		}
	}
	return resolve(event);
};
