<script lang="ts">
    import { Button } from "bits-ui";
	import type { Snippet } from "svelte";
    export type ButtonProps = {
        children: Snippet
        onClick: () => void
    }

    let {children, onClick} : ButtonProps = $props()
</script>

<Button.Root onclick={onClick} class="rounded-md bg-[#365486] p-2 font-light text-sm text-[#DCF2F1] border border-[#365486] hover:border-[#7FC7D9]">
    {@render children()}
</Button.Root>