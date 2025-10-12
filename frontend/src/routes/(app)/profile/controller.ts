import type { Fetch, SeekPaginatedResponse, ServerResponse } from "$lib/types";
import { FetchData } from "$lib/utils";
import { FILE_IMAGE_THRESHOLD, MAX_BIO_LENGTH } from "$lib/utils/constants";
import { IsAlphaOnly } from "$lib/utils/helper";
import type { StudentOrders, StudentProfile } from "./model";

class profileController {
    async updateProfile(fetch: Fetch, req: Request){
        const url = 'http://localhost:8080/api/v1/students/me'
        const body = await req.formData()
        const profileImage = body.get('file') as File | null
        const name = body.get('name') as string | null
        const bio = body.get('bio') as string | null
        if (profileImage && profileImage.size > 0){
            if (profileImage.type !== 'image/png'){
                return {success: false, message: 'wrong image format', status: 400}
            }
            if (profileImage.size > FILE_IMAGE_THRESHOLD){
                return {success: false, message: 'file size is too large', status: 400}
            }
        }
        if (name && !IsAlphaOnly(name)){
            return {success: false, message: 'name must contain alphabet only', status: 400}
        }
        if (bio && bio.length > MAX_BIO_LENGTH){
            return {success: false, message: `bio length must be shorter than ${MAX_BIO_LENGTH} characters`, status: 400}
        }
        const {success, message, status} = await FetchData(fetch, url, 'PATCH', body)
        return {success, message, status}
    }
    async sendVerificationLink(fetch : Fetch){
        const url = 'http://localhost:8080/api/v1/verify/send'
        const {success, message, status} = await FetchData(fetch, url)
        return {success, message, status}
    }
    async getProfile(fetch: Fetch){
        const url = 'http://localhost:8080/api/v1/me'
        const {success, message, status, res} = await FetchData(fetch, url, 'GET')
        if (!success){
            return {success, message, status}
        }
        const resBody : ServerResponse<StudentProfile> = await res?.json()
        return {resBody, status, success, message}
    }
    async getOrders(fetch: Fetch){
        const url = `http://localhost:8080/api/v1/me/course-requests`
        const {success, message, status, res} = await FetchData(fetch, url, 'GET')
        if (!success){
            return {success, message, status}
        }
        const resBody : ServerResponse<SeekPaginatedResponse<StudentOrders>> = await res?.json()
        return {success, message, status, resBody: resBody.data}
    }
    async updateOrdersList(fetch: Fetch, req: Request){
        const queries : string[] = []
        const formData = await req.formData()
        const search = formData.get('search')
        const orderStatus = formData.get('status')
        const lastID = formData.get('last_id')
        let url = `http://localhost:8080/api/v1/me/course-requests?`
        if (search){
            queries.push(`search=${search}`)
        }
        if (orderStatus){
            queries.push(`status=${orderStatus}`)
        }
        if (lastID){
            queries.push(`last_id=${lastID}`)
        }
        if (queries.length > 0){
            url += queries.join('&')
        }
        const {success, message, status, res} = await FetchData(fetch, url, 'GET')
        if (!success){
            return {success, message, status}
        }
        const resBody : ServerResponse<SeekPaginatedResponse<StudentOrders>> = await res?.json()
        return {success, message, status, resBody: resBody.data}
    }
}

export const controller = new profileController()