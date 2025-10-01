import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { controller } from "./controller";

export const load: PageServerLoad = async ({fetch})=>{
    const {success, message, status, resBody} = await controller.getProfile(fetch)
    if (!success){
        throw error(status, {message})
    }
    return {resBody, message, status}
}