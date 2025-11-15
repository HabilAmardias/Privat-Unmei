import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ cookies }) => {
	const role = cookies.get('role');
	const userStatus = cookies.get('status');
	return { role, userStatus };
};
