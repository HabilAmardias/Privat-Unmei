import type { AuthClaim } from "$lib/types";

function DecodeJWT(token: string){
    const claim : AuthClaim = JSON.parse(Buffer.from(token.split('.')[1], 'base64').toString())
    return claim
}

export function IsTokenExpired(token: string | undefined){
    if(!token){
        return true
    }
    const claim = DecodeJWT(token)
    const today = Math.floor(Date.now()/1000)
    return claim.exp - today < 120
}