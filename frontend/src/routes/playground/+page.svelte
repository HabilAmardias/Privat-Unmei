<script lang="ts">
    import Button from "$lib/components/button/Button.svelte";
    import Link from "$lib/components/button/Link.svelte";
    import toast from 'svelte-french-toast';
    import Image from "$lib/components/image/Image.svelte";
	import Menubar from "$lib/components/menubar/Menubar.svelte";
	import Menu from "$lib/components/menubar/Menu.svelte";
    import imageSrc from "$lib/images/svelte-welcome.png"
	import Datepicker from "$lib/components/calendar/Datepicker.svelte";

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
        <Link theme="dark" href="/">
            Home
        </Link>
    </span>
    <span class="bg-[var(--tertiary-color)] p-2 rounded-md">
        <Link href="/">
            Home
        </Link>
    </span>
    <Image src={imageSrc} alt="test-image" height={256} width={256} round="full"/>
    <Menubar>
        <Menu separator={true}>
            <Link href="/">
                Home
            </Link>
        </Menu>
        <Menu>
            <Link href="/playground">
                Playground
            </Link>
        </Menu>
    </Menubar>
    <Datepicker label="test" />
</div>