<script lang="ts">
	import Link from '$lib/components/button/Link.svelte';

	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>Home - Privat Unmei</title>
	<meta name="description" content="Home - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex flex-col gap-4 p-4">
	<h2 class="text-2xl font-bold text-[var(--tertiary-color)]">Courses</h2>
	<ul class="grid grid-cols-2 gap-4 md:grid-cols-3">
		{#each data.courses.entries as c (c.id)}
			<li>
				<Link href={`/courses/${c.id}`}>
					<div
						class="flex h-[200px] flex-col justify-between rounded-lg bg-[var(--tertiary-color)] p-4"
					>
						<div>
							<p class="font-bold text-[var(--primary-color)]">{c.title}</p>
							<p class="text-[var(--secondary-color)]">{c.mentor_name}</p>
						</div>
						<div>
							<p class="text-[var(--secondary-color)]">
								{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(
									c.price
								)}
							</p>
						</div>
					</div>
				</Link>
			</li>
		{/each}
	</ul>
	<h2 class="text-2xl font-bold text-[var(--tertiary-color)]">Mentors</h2>
</div>
