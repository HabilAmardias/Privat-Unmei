<script lang="ts">
    import Button from "$lib/components/button/Button.svelte";
    import Link from "$lib/components/button/Link.svelte";
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
    <span class="bg[var(--secondary-color)] p-2 rounded-md">
        <Link href="/">
            Home
        </Link>
    </span>
    <span class="bg-[var(--tertiary-color)] p-2 rounded-md">
        <Link href="/" inverted={true}>
            Home
        </Link>
    </span>
    <Image src={srcImg} alt="test-image" rounded={true}/>
</div>

<Toaster />