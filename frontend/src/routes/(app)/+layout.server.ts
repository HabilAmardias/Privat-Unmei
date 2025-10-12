import type { LayoutServerLoad } from "./$types";


export const load : LayoutServerLoad = async ({cookies})=>{
    const authToken = cookies.get('auth_token')
    return {isLoggedIn: authToken? true : false}
}