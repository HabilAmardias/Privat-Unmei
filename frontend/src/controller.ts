import type { CookiesData, SameSite } from '$lib/types';

class CookieController {
	GetCookies(res: Response) {
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

export const controller = new CookieController();
