import type { Fetch, HTTPMethod } from '$lib/types';

export async function FetchData(
	fetch: Fetch,
	url: string | URL,
	method?: HTTPMethod,
	body?: BodyInit,
	headers?: Headers
) {
	const res = await fetch(url, {
		method,
		body,
		credentials: 'include',
		redirect: 'manual',
		headers
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
