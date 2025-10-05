import type { StudentOrders } from "./model"

class profileView {
    verifyIsLoading = $state<boolean>(false)
    ordersIsLoading = $state<boolean>(false)
    profileIsLoading = $state<boolean>(false)
    isDesktop = $state<boolean>()
    isEdit = $state<boolean>(false)
    name = $state<string>("")
    bio = $state<string>("")
    profileImage = $state<FileList>()
    status = $state<string>("")
    totalRow = $state<number>(1) // temporary
    limit = $state<number>(15)
    lastID = $state<number>(15)
    pageNumber = $state<number>(1)
    paginationForm = $state<HTMLFormElement | undefined>()
    
    orders = $state<StudentOrders[]>()
    
    size = $derived.by<number>(()=>{
        if (this.isDesktop){
            return 200
        }
        return 100
    })
    onPageChange(num: number){
        if (num < this.pageNumber){
            const lastOrder = (this.orders!)[0]
            this.lastID = lastOrder.id
        } else if (num > this.pageNumber){
            const lastIndex = this.orders!.length - 1
            const lastOrder = (this.orders!)[lastIndex]
            this.lastID = lastOrder.id
        }
        this.pageNumber = num
        this.paginationForm?.requestSubmit()
    }
    setTotalRow(row: number){
        this.totalRow = row
    }
    setOrdersIsLoading(b: boolean){
        this.ordersIsLoading = b
    }
    setVerifyIsLoading(b: boolean){
        this.verifyIsLoading = b
    }
    setProfileIsLoading(b: boolean){
        this.profileIsLoading = b
    }
    setBio(newBio: string){
        this.bio = newBio
    }
    setName(newName: string){
        this.name = newName
    }
    setIsDesktop(b: boolean){
        this.isDesktop = b
    }
    setIsEdit(){
        this.isEdit = !this.isEdit
    }
    setOrders(newOrders: StudentOrders[]){
        this.orders = newOrders
    }
    setLastID(newID: number){
        this.lastID = newID
    }
}

export const View = new profileView()