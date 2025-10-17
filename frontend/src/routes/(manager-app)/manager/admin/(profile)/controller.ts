import type { Fetch } from "$lib/types";
import { FetchData } from "$lib/utils";

class AdminPageController {
    async getProfile(fetch: Fetch){
        const url = 'http://localhost:8080/api/v1/admins/me'
        const {success, message, status, res} = await FetchData(fetch, url, 'GET')
        if (!success){
            return {success, message, status}
        }
        return {success, message, status, res}
    }
    async changePassword(fetch: Fetch, req: Request){
        const url = 'http://localhost:8080/api/v1/admins/me/change-password'
        const formData = await req.formData()
        const password = formData.get('password')
        const repeatPassword = formData.get('repeat-password')
        if (!password){
            return {success: false, message: 'please insert an valid password', status: 400}
        }
        if (!this.#validatePassword(password as string)){
            return {success: false, message: 'please insert an valid password', status: 400}
        }
        if (!repeatPassword){
            return {success: false, message: 'password does not match', status: 400}
        }
        if (repeatPassword !== password){
            return {success: false, message: 'password does not match', status: 400}
        }
        const body = JSON.stringify({
            password
        })
        const {success, message, status, res} = await FetchData(fetch, url, 'POST',body)
        return {success, message, status, res}
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

export const controller = new AdminPageController()