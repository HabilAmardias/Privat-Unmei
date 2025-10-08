class HomeView {
    isDesktop = $state<boolean>()
    imageWidth = $derived.by<number>(()=>{
        return this.isDesktop ? 250 : 160
    })
    imageHeight = $derived.by<number>(()=>{
        return this.isDesktop ? 80 : 50
    })
}

export const View = new HomeView()