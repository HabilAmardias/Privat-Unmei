<script lang="ts">
  import { unstable_RatingGroup } from "bits-ui";
  import { Star } from "@lucide/svelte";
  type ratingGroupProps = {
    onChange? : (rate: number) => void
  }
  let {onChange} : ratingGroupProps = $props()
</script>
 
<unstable_RatingGroup.Root onValueChange={onChange} max={5} class="flex gap-1">
  {#snippet children({ items })}
    {#each items as item (item.index)}
      <unstable_RatingGroup.Item
        index={item.index}
        class="text-[var(--primary-color)] hover:text-[var(--primary-color)] data-[state=inactive]:text-[var(--tertiary-color)] group size-10 cursor-pointer transition-colors md:size-8"
      >
        <Star class="size-full fill-current"/>
      </unstable_RatingGroup.Item>
    {/each}
  {/snippet}
</unstable_RatingGroup.Root>