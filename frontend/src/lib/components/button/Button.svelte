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

<Button.Root formaction={formAction} onclick={onClick} class="bg-[var(--tertiary-color)] p-2 text-[var(--secondary-color)] hover:text-[var(--primary-color)] cursor-pointer rounded-md">
    {@render children()}
</Button.Root>