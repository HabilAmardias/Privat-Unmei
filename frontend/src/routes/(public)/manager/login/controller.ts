import type { Fetch, ServerResponse } from '$lib/types';
import type { CookiesData, SameSite } from '$lib/types';
import type { AdminLoginResponse } from './model';
import { FetchData } from '$lib/utils';

class ManagerAuthController {
	async loginMentor(
		req: Request,
		fetch: Fetch
	): Promise<{
		cookiesData?: CookiesData[];
		success: boolean;
		message: string;
		status: number;
	}> {
		const formData = await req.formData();
		const email = formData.get('email');
		const password = formData.get('password');
		if (!email) {
			return {
				success: false,
				message: 'please insert an email',
				status: 400
			};
		}
		if (!this.#validateEmail(email as string)) {
			return {
				success: false,
				message: 'please insert an valid email',
				status: 400
			};
		}
		if (!password) {
			return {
				success: false,
				message: 'please insert an password',
				status: 400
			};
		}
		if (!this.#validatePassword(password as string)) {
			return {
				success: false,
				message: 'password need to be at least 8 characters and contain any @#!?',
				status: 400
			};
		}
		const url = 'http://localhost:80/api/v1/mentor/login';
		const body = JSON.stringify({
			email: email,
			password: password
		});
		const { success, res, message, status } = await FetchData(fetch, url, 'POST', body);
		if (!success) {
			return {
				success: false,
				message,
				status
			};
		}
		const cookiesData = this.#getCookies(res!);
		return {
			cookiesData,
			success: true,
			message: 'successfully login',
			status: res!.status
		};
	}
	async loginAdmin(
		req: Request,
		fetch: Fetch
	): Promise<{
		cookiesData?: CookiesData[];
		success: boolean;
		message: string;
		userStatus?: 'verified' | 'unverified';
		status: number;
	}> {
		const formData = await req.formData();
		const email = formData.get('email');
		const password = formData.get('password');
		if (!email) {
			return {
				success: false,
				message: 'please insert an email',
				status: 400
			};
		}
		if (!this.#validateEmail(email as string)) {
			return {
				success: false,
				message: 'please insert an valid email',
				status: 400
			};
		}
		if (!password) {
			return {
				success: false,
				message: 'please insert an password',
				status: 400
			};
		}
		if (!this.#validatePassword(password as string)) {
			return {
				success: false,
				message: 'password need to be at least 8 characters and contain any @#!?',
				status: 400
			};
		}
		const url = 'http://localhost:80/api/v1/admin/login';
		const body = JSON.stringify({
			email: email,
			password: password
		});
		const { success, res, message, status } = await FetchData(fetch, url, 'POST', body);
		if (!success) {
			return {
				success: false,
				message,
				status
			};
		}
		const resBody: ServerResponse<AdminLoginResponse> = await res?.json();
		const cookiesData = this.#getCookies(res!);
		return {
			cookiesData,
			success: true,
			message: 'successfully login',
			status: res!.status,
			userStatus: resBody.data.status
		};
	}
	#validateEmail(email: string) {
		const pattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		return pattern.test(email);
	}
	#validatePassword(password: string) {
		const minLength = password.length >= 8;
		const hasSpecialChar =
			password.includes('!') ||
			password.includes('@') ||
			password.includes('#') ||
			password.includes('?');
		return minLength && hasSpecialChar;
	}
	#getCookies(res: Response) {
		const out: CookiesData[] = [];
		const setCookies = res.headers.getSetCookie();
		setCookies.forEach((val) => {
			const data: CookiesData = {
				key: '',
				value: '',
				path: ''
			};

			const keyValPairs = val.split(';');
			keyValPairs.forEach((e) => {
				const [key, val] = e.split('=');
				switch (key.toLowerCase().trim()) {
					case 'refresh_token':
						data.key = key;
						data.value = val;
						break;
					case 'auth_token':
						data.key = key;
						data.value = val;
						break;
					case 'status':
						data.key = key;
						data.value = val;
						break;
					case 'role':
						data.key = key;
						data.value = val;
						break;
					case 'domain':
						data.domain = val;
						break;
					case 'max-age':
						data.maxAge = parseInt(val);
						break;
					case 'path':
						data.path = val;
						break;
					case 'samesite':
						data.sameSite = val.toLowerCase() as SameSite;
						break;
					case 'httponly':
						data.httpOnly = true;
						break;
					default:
						break;
				}
			});
			out.push(data);
		});
		return out;
	}
}

export const controller = new ManagerAuthController();
