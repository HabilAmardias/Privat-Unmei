import type { CookiesData, LoginResponse, SameSiteType } from './model';

class AuthController {
	async register(req: Request): Promise<{ success: boolean; message: string; status: number }> {
		const formData = await req.formData();
		const name = formData.get('name');
		const email = formData.get('email');
		const password = formData.get('password');
		const repeatPassword = formData.get('repeat-password');
		const token = formData.get('token');
		if (!email) {
			return { success: false, message: 'please insert an email', status: 400 };
		}
		if (!this.#validateEmail(email as string)) {
			return { success: false, message: 'please insert an valid email', status: 400 };
		}
		if (!password) {
			return { success: false, message: 'please insert an password', status: 400 };
		}
		if (!this.#validatePassword(password as string)) {
			return {
				success: false,
				message: 'password need to be at least 8 characters and contain any @#!?',
				status: 400
			};
		}
		if (repeatPassword !== password) {
			return { success: false, message: 'both password input are not same', status: 400 };
		}
		if (!name) {
			return { success: false, message: 'please insert a name', status: 400 };
		}
		const url = 'http://localhost:8080/api/v1/register';
		const body = JSON.stringify({
			name: name,
			email: email,
			password: password,
			captcha_token: token
		});
		const res = await fetch(url, {
			method: 'POST',
			body: body
		});
		if (!res.ok) {
			const resData = await res.json();
			if ('message' in resData.data) {
				return { success: false, message: resData.data.message as string, status: res.status };
			}
			return { success: false, message: 'invalid input', status: res.status };
		}
		return { success: true, message: 'successfully registered', status: res.status };
	}
	async login(req: Request): Promise<{
		responseBody: LoginResponse | undefined;
		cookiesData: CookiesData[] | undefined;
		success: boolean;
		message: string;
		status: number;
	}> {
		const formData = await req.formData();
		const email = formData.get('email');
		const password = formData.get('password');
		if (!email) {
			return {
				responseBody: undefined,
				cookiesData: undefined,
				success: false,
				message: 'please insert an email',
				status: 400
			};
		}
		if (!this.#validateEmail(email as string)) {
			return {
				responseBody: undefined,
				cookiesData: undefined,
				success: false,
				message: 'please insert an valid email',
				status: 400
			};
		}
		if (!password) {
			return {
				responseBody: undefined,
				cookiesData: undefined,
				success: false,
				message: 'please insert an password',
				status: 400
			};
		}
		if (!this.#validatePassword(password as string)) {
			return {
				responseBody: undefined,
				success: false,
				cookiesData: undefined,
				message: 'password need to be at least 8 characters and contain any @#!?',
				status: 400
			};
		}
		const url = 'http://localhost:8080/api/v1/login';
		const body = JSON.stringify({
			email: email,
			password: password
		});
		const res = await fetch(url, {
			method: 'POST',
			body: body,
			credentials: 'include'
		});

		if (!res.ok) {
			const resData = await res.json();
			if ('message' in resData.data) {
				return {
					responseBody: undefined,
					cookiesData: undefined,
					success: false,
					message: resData.data.message as string,
					status: res.status
				};
			}
			return {
				responseBody: undefined,
				cookiesData: undefined,
				success: false,
				message: 'invalid input',
				status: res.status
			};
		}
		const cookiesData = this.#getCookies(res);
		const responseBody: LoginResponse = await res.json();
		return {
			responseBody,
			cookiesData,
			success: true,
			message: 'successfully login',
			status: res.status
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
						data.sameSite = val.toLowerCase() as SameSiteType;
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

export const controller = new AuthController();
