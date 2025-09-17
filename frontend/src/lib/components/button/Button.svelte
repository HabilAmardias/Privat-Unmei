<script lang="ts">
    import { Button } from "bits-ui";
	import type { Snippet } from "svelte";
    export type ButtonProps = {
        children: Snippet
        onClick?: () => void
        formAction? : string
    }

    let {children, onClick, formAction} : ButtonProps = $props()
</script>

<Button.Root formaction={formAction} onclick={onClick} class="rounded-md bg-[var(--tertiary-color)] p-2 font-light text-sm text-[var(--secondary-color)] border border-[--tertiary-color] hover:border-[var(--primary-color)] hover:text-[var(--primary-color)] hover:cursor-pointer">
    {@render children()}
</Button.Root>