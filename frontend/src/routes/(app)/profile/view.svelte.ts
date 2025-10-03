class profileView {
    isDesktop = $state<boolean>()
    isEdit = $state<boolean>(false)
    name = $state<string>("")
    bio = $state<string>("")
    profileImage = $state<FileList>()
    status = $state<string>("")
    totalRow = $state<number>(100)
    limit = $state<number>(15)
    paginationForm = $state<HTMLFormElement | undefined>()
    
    size = $derived.by<number>(()=>{
        if (this.isDesktop){
            return 200
        }
        return 100
    })
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

}

export const View = new profileView()