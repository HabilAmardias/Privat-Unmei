import { error, fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';
import { controller } from './controller';

export const load: PageServerLoad = async ({ fetch, params }) => {
	const [profileRes, messageRes, chatroomRes] = await Promise.all([
		controller.getProfile(fetch),
		controller.getMessages(fetch, params.slug),
		controller.getInfo(fetch, params.slug)
	]);
	if (!profileRes.success) {
		throw error(profileRes.status, { message: profileRes.message });
	}
	if (!messageRes.success) {
		throw error(messageRes.status, { message: messageRes.message });
	}
	if (!chatroomRes.success) {
		throw error(chatroomRes.status, { message: chatroomRes.message });
	}

	return {
		chatroom: chatroomRes.resBody.data,
		messages: messageRes.resBody.data,
		profile: profileRes.resBody.data
	};
};

export const actions = {
	getMessage: async ({ fetch, request, params }) => {
		if (!params.slug) {
			return fail(404, { message: 'no data found' });
		}
		const { success, message, status, resBody } = await controller.getMessages(
			fetch,
			params.slug,
			request
		);
		if (!success) {
			return fail(status, { message });
		}
		return {
			messages: resBody.data
		};
	},
	updateLastRead: async ({ fetch, params }) => {
		if (!params.slug) {
			return fail(404, { message: 'no data found' });
		}
		const { success, message, status } = await controller.updateLastRead(fetch, params.slug);
		if (!success) {
			throw error(status, { message });
		}
		return {
			success
		};
	},
	sendMessage: async ({ fetch, request, params }) => {
		if (!params.slug) {
			return fail(404, { message: 'no data found' });
		}
		const { success, message, status } = await controller.sendMessage(fetch, params.slug, request);
		if (!success) {
			return fail(status, { message });
		}
		return { message, status };
	}
} satisfies Actions;
