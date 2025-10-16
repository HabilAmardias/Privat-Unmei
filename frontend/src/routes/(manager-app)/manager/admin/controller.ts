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
}

export const controller = new AdminPageController()