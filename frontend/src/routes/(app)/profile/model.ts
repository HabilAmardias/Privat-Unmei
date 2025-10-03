export type StudentProfile = {
    name: string
    bio: string
    profile_image: string
    email: string
}

export type StudentOrders = {
    id: number
    student_id: string
    course_id: number
    total_price: number
    status: string
    mentor_name: string
    mentor_email: string
    course_name: string
}

export const statusOptions = [
    {value: 'reserved', label:'reserved'},
    {value: 'pending payment', label:'pending'},
    {value: 'scheduled', label:'scheduled'},
    {value: 'completed', label:'completed'},
    {value: 'cancelled', label:'cancelled'}
]