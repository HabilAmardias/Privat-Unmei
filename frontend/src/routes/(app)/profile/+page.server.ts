import { error, fail, type Actions } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { controller } from "./controller";

export const load: PageServerLoad = async ({fetch, cookies})=>{
    const [profile, orders] = await Promise.all([controller.getProfile(fetch), controller.getOrders(fetch)])
    if (!profile.success){
        throw error(profile.status, {message: profile.message})
    }
    const userStatus = cookies.get('status')
    return {profile: profile.resBody.data, orders: orders.resBody, userStatus}
}

export const actions = {
    updateProfile: async ({request, fetch}) =>{
        const {success, message, status} = await controller.updateProfile(fetch, request)
        if (!success){
            return fail(status, {message})
        }
        return {message}
    },
    sendVerification: async({fetch}) =>{
        const {success, message, status} = await controller.sendVerificationLink(fetch)
        if (!success){
            return fail(status, {message})
        }
        return {success, message, status}
    },
    myOrders: async ({request, fetch})=>{
        const {success, message, status, resBody} = await controller.updateOrdersList(fetch, request)
        if (!success){
            return fail(status, {message})
        }
        return {orders: resBody.entries, lastID: resBody.page_info.last_id, limit: resBody.page_info.limit, totalRow: resBody.page_info.total_row, message}
    }

} satisfies Actions