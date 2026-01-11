import type { Fetch, HTTPMethod } from '$lib/types';
import { PUBLIC_BASE_URL } from '$env/static/public';

export async function FetchData(
	fetch: Fetch,
	url: string | URL,
	method?: HTTPMethod,
	body?: BodyInit
) {
	const fullURL = `${PUBLIC_BASE_URL}${url}`;
	const res = await fetch(fullURL, {
		method,
		body
	});
	if (!res.ok && res.status !== 307) {
		const resBody = await res.json();
		if ('message' in resBody.data) {
			return { success: false, message: resBody.data?.message as string, status: res.status };
		}
		return { success: false, message: 'invalid input', status: res.status };
	}

	return { success: true, message: 'success', status: res.status, res };
}
