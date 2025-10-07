<script lang="ts">
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

{#if data.resBody.data.length > 0}
	{#each data.resBody.data as course (course.id)}
		<div>
			<p>{course.title}</p>
			<p>{course.price}</p>
			<p>{course.mentor_name}</p>
			<p>{course.session_duration_minutes}</p>
		</div>
	{/each}
{:else}
	<div>No courses found</div>
{/if}
