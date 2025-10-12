export class HomeView {
    isDesktop = $state<boolean>(false)
    imageWidth = $derived.by<number>(()=>{
        return this.isDesktop ? 250 : 160
    })
    imageHeight = $derived.by<number>(()=>{
        return this.isDesktop ? 80 : 50
    })
    iconsSize = $derived.by<number>(()=>{
        return this.isDesktop ? 40 : 32
    })
    setIsDesktop(b: boolean){
        this.isDesktop = b
    }
}