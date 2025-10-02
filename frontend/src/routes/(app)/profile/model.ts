export type StudentProfile = {
    name: string
    bio: string
    profile_image: string
    email: string
}

export const statusOptions = [
    {value: 'reserved', label:'reserved'},
    {value: 'pending payment', label:'pending payment'},
    {value: 'scheduled', label:'scheduled'},
    {value: 'completed', label:'completed'},
    {value: 'cancelled', label:'cancelled'}
]