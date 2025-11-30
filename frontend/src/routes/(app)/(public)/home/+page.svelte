<script lang="ts">
	import Link from '$lib/components/button/Link.svelte';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';

	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>Home - Privat Unmei</title>
	<meta name="description" content="Home - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex flex-col gap-4 p-4">
	<h2 class="text-2xl font-bold text-[var(--tertiary-color)]">Most Bought</h2>
	<ScrollArea orientation="vertical" viewportClasses="h-[300px] max-h-[300px]">
		<ul class="flex flex-col gap-4 md:grid md:grid-cols-3">
			{#each data.mostBought as c (c.id)}
				<li class="transition-transform hover:-translate-y-1">
					<Link href={`/courses/${c.id}`}>
						<div
							class="flex h-[100px] flex-col justify-between rounded-lg bg-[var(--tertiary-color)] p-4"
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
	</ScrollArea>
	<div class="flex justify-between">
		<h2 class="text-2xl font-bold text-[var(--tertiary-color)]">Courses</h2>
		<div class="flex items-center justify-center rounded-lg bg-[var(--tertiary-color)] p-2">
			<Link href="/courses">See More</Link>
		</div>
	</div>
	<ScrollArea orientation="vertical" viewportClasses="h-[300px] max-h-[300px]">
		<ul class="flex flex-col gap-4 md:grid md:grid-cols-3">
			{#each data.courses.entries as c (c.id)}
				<li class="transition-transform hover:-translate-y-1">
					<Link href={`/courses/${c.id}`}>
						<div
							class="flex h-[100px] flex-col justify-between rounded-lg bg-[var(--tertiary-color)] p-4"
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
	</ScrollArea>
	<div class="flex justify-between">
		<h2 class="text-2xl font-bold text-[var(--tertiary-color)]">Mentors</h2>
		<div class="flex items-center justify-center rounded-lg bg-[var(--tertiary-color)] p-2">
			<Link href="/mentors">See More</Link>
		</div>
	</div>
	<ScrollArea orientation="vertical" viewportClasses="h-[300px] max-h-[300px]">
		<ul class="flex flex-col gap-4 md:grid md:grid-cols-3">
			{#each data.mentors.entries as m (m.id)}
				<li class=" transition-transform hover:-translate-y-1">
					<Link href={`/mentors/${m.id}`}>
						<div class="flex h-[100px] gap-4 rounded-lg bg-[var(--tertiary-color)] p-4">
							<CldImage src={m.profile_image} width={70} height={70} className="rounded-full" />
							<div>
								<p class="font-bold text-[var(--primary-color)]">{m.name}</p>
								<p class="text-[var(--secondary-color)]">{m.email}</p>
							</div>
						</div>
					</Link>
				</li>
			{/each}
		</ul>
	</ScrollArea>
</div>
