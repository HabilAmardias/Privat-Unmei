import type { Fetch, CookiesData, SameSite } from '$lib/types';
import { FetchData } from '$lib/utils';

class LoginCallbackController {
	async login(fetch: Fetch, req: Request) {
		const formData = await req.formData();
		const otpStr = formData.get('otp');
		const otp = parseInt(otpStr as string);

		const reqBody = JSON.stringify({
			otp
		});
		const url = '/api/v1/login-callback';
		const { success, message, status, res } = await FetchData(fetch, url, 'POST', reqBody);
		if (!success) {
			return { success, message, status };
		}
		const cookiesData = this.#getCookies(res!);
		return { success, message, status, cookiesData };
	}
	async resendOTP(fetch: Fetch) {
		const url = '/api/v1/resend-otp';
		const { success, message, status, res } = await FetchData(fetch, url, 'GET');
		if (!success) {
			return { success, message, status };
		}
		const cookiesData = this.#getCookies(res!);
		return { success, message, status, cookiesData };
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
					case 'login_token':
						data.key = key;
						data.value = val;
						break;
					case 'status':
						data.key = key;
						data.value = val;
						break;
					case 'oauthstate':
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

export const controller = new LoginCallbackController();
