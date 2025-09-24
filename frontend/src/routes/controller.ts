import type { LoginResponse } from './model';

class AuthController {
	async register(req: Request): Promise<boolean> {
		const formData = await req.formData();
		const name = formData.get('name');
		const email = formData.get('email');
		const password = formData.get('password');
		const repeatPassword = formData.get('repeat-password');
		const token = formData.get('token');
		if (!email) {
			return false;
		}
		if (!this.#validateEmail(email as string)) {
			return false;
		}
		if (!password) {
			return false;
		}
		if (!this.#validatePassword(password as string)) {
			return false;
		}
		if (repeatPassword !== password) {
			return false;
		}
		if (!name) {
			return false;
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
			return false;
		}
		return true;
	}
	async login(req: Request): Promise<{ data: LoginResponse | undefined; success: boolean }> {
		const formData = await req.formData();
		const email = formData.get('email');
		const password = formData.get('password');
		if (!email) {
			return { data: undefined, success: false };
		}
		if (!this.#validateEmail(email as string)) {
			return { data: undefined, success: false };
		}
		if (!password) {
			return { data: undefined, success: false };
		}
		if (!this.#validatePassword(password as string)) {
			return { data: undefined, success: false };
		}
		const url = '/api/v1/login';
		const body = JSON.stringify({
			email: email,
			password: password
		});
		const res = await fetch(url, {
			method: 'POST',
			body: body
		});
		if (!res.ok) {
			return { data: undefined, success: false };
		}
		const data: LoginResponse = await res.json();
		return { data, success: true };
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
}

export const controller = new AuthController();
