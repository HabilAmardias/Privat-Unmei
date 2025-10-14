import type { Actions } from "@sveltejs/kit";
import { fail } from "@sveltejs/kit";
import { controller } from "./controller";
import type { PageServerLoad } from "./$types";

export const load : PageServerLoad = async ({cookies}) =>{
    const authToken = cookies.get('auth_token')
    const role = cookies.get('role')
    return {authToken, role}
}

export const actions = {
    loginMentor: async ({ request, cookies, fetch }) => {
        const { cookiesData, success, message, status } = await controller.loginMentor(
            request,
            fetch
        );
        if (!success) {
            return fail(status, { message });
        }
        cookiesData?.forEach((c) => {
            cookies.set(c.key, c.value, {
                path: c.path,
                domain: c.domain,
                httpOnly: c.httpOnly,
                maxAge: c.maxAge,
                sameSite: c.sameSite
            });
        });
        return { success };
    },
    loginAdmin: async ({ request, cookies, fetch }) => {
        const { cookiesData, success, message, status, userStatus } = await controller.loginAdmin(
            request,
            fetch
        );
        if (!success) {
            return fail(status, { message });
        }
        cookiesData?.forEach((c) => {
            cookies.set(c.key, c.value, {
                path: c.path,
                domain: c.domain,
                httpOnly: c.httpOnly,
                maxAge: c.maxAge,
                sameSite: c.sameSite
            });
        });
        return { success, userStatus };
    }
} satisfies Actions