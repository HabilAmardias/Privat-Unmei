<script lang="ts">
    import { Button } from "bits-ui";
	import type { Snippet } from "svelte";
    export type ButtonProps = {
        children: Snippet
        onClick: () => void
    }

    let {children, onClick} : ButtonProps = $props()
</script>

<Button.Root onclick={onClick} class="rounded-md bg-[var(--tertiary-color)] p-2 font-light text-sm text-[var(--secondary-color)] border border-[--tertiary-color] hover:border-[var(--primary-color)]">
    {@render children()}
</Button.Root>