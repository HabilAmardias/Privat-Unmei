<script lang="ts">
    import Button from "$lib/components/button/Button.svelte";
    import Link from "$lib/components/button/Link.svelte";
    import toast from 'svelte-french-toast';
    import Image from "$lib/components/image/Image.svelte";
	import Menubar from "$lib/components/menubar/Menubar.svelte";
	import Menu from "$lib/components/menubar/Menu.svelte";
    import imageSrc from "$lib/images/svelte-welcome.png"

    function openToastSuccess(){
        toast.success("Success!",{
            position: "top-right"
        })
    }
</script>

<div>
    <Button className="bg-[var(--tertiary-color)] p-2 text-[var(--secondary-color)] hover:text-[var(--primary-color)] cursor-pointer rounded-md">
        Button Example
    </Button>
    <Button className="bg-[var(--tertiary-color)] p-2 text-[var(--secondary-color)] hover:text-[var(--primary-color)] cursor-pointer rounded-md" onClick={openToastSuccess}>
        Click for success toast
    </Button>
    <span class="bg[var(--secondary-color)] p-2 rounded-md">
        <Link className="text-[var(--tertiary-color)] hover:text-[var(--primary-color)]" href="/">
            Home
        </Link>
    </span>
    <span class="bg-[var(--tertiary-color)] p-2 rounded-md">
        <Link href="/" className="text-[var(--secondary-color)] hover:text-[var(--primary-color)]">
            Home
        </Link>
    </span>
    <Image src={imageSrc} alt="test-image" className="rounded-full" height={32} width={32}/>
    <Menubar className="bg-[var(--tertiary-color)] p-2 flex gap-1">
        <Menu separator={true}>
            <Link href="/home" className="text-[var(--secondary-color)] hover:text-[var(--primary-color)]">
                Home
            </Link>
        </Menu>
        <Menu>
            <Link href="/playground" className="text-[var(--secondary-color)] hover:text-[var(--primary-color)]">
                Playground
            </Link>
        </Menu>
    </Menubar>
</div>