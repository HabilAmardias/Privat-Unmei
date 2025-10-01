import type { Fetch, ServerResponse } from "$lib/types";
import { FetchData } from "$lib/utils";
import type { StudentProfile } from "./model";

class profileController {
    async getProfile(fetch: Fetch){
        const url = 'http://localhost:8080/api/v1/me'
        const {success, message, status, res} = await FetchData(fetch, url, 'GET')
        if (!success){
            return {success, message, status}
        }
        const resBody : ServerResponse<StudentProfile> = await res?.json()
        return {resBody, status, success, message}
    }
}

export const controller = new profileController()