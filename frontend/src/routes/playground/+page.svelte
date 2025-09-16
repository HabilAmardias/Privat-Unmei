<script lang="ts">
    import Button from "$lib/components/button/Button.svelte";
    import toast, { Toaster } from 'svelte-french-toast';
</script>

<div>
    <Button onClick={() => {}}>
        Button Example
    </Button>
    <Button onClick={() => {toast.success("Success!", {
        position: "top-right"
    })}}>
        Click for success toast
    </Button>
</div>

<Toaster />