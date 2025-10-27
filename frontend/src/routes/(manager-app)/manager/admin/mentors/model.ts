export type mentorList = {
    id: string
    name: string
    email: string
    years_of_experience: number
}

export type adminProfile = {
    name: string
    email: string
    bio: string
    profile_image: string
    status: 'verified' | 'unverified'
}