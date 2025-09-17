<script lang="ts">
    import Button from "$lib/components/button/Button.svelte";
    import Anchor from "$lib/components/button/Anchor.svelte";
    import toast, { Toaster } from 'svelte-french-toast';
    import Image from "$lib/components/image/Image.svelte";
    import srcImg from '$lib/images/svelte-welcome.png'

    function openToastSuccess(){
        toast.success("Success!",{
            position: "top-right"
        })
    }
</script>

<div>
    <Button>
        Button Example
    </Button>
    <Button onClick={openToastSuccess}>
        Click for success toast
    </Button>
    <Anchor href="/">
        Home
    </Anchor>
    <Image src={srcImg} alt="test-image" rounded={true}/>
</div>

<Toaster />