import type { Fetch } from "$lib/types";
import { FetchData } from "$lib/utils";

class AdminVerifyController{
    async verifyAdmin(fetch: Fetch, req: Request){
        const formData = await req.formData()
        const newEmail = formData.get('email')
        const newPassword = formData.get('password')
        if (!newEmail){
            return {success: false, status: 400, message: 'please insert an email'}
        }
        if(!this.#validateEmail(newEmail as string)){
            return {success: false, status: 400, message: 'please insert a valid email'}
        }
        if (!newPassword){
            return {success: false, status: 400, message: 'please insert an password'}
        }
        if (!this.#validatePassword(newPassword as string)){
            return {success: false, status: 400, message: 'please insert a valid password'}
        }
        const reqBody = JSON.stringify({
            email: newEmail,
            password: newPassword
        })
        const url = 'http://localhost:8080/api/v1/admins/me/verify'
        const {success, message, status} = await FetchData(fetch, url, 'POST', reqBody)
        return {success, message, status}
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

export const controller = new AdminVerifyController()