<script lang="ts">
	import Link from '$lib/components/button/Link.svelte';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';

	import type { PageProps } from './$types';

	let { data }: PageProps = $props();

	const coursesSchema = {
		'@context': 'https://schema.org',
		'@type': 'CollectionPage',
		name: 'Privat Unmei - Online Courses',
		description: 'Discover top-rated online courses from expert mentors on Privat Unmei',
		mainEntity: {
			'@type': 'ItemList',
			itemListElement: data.courses.entries.map((course, index) => ({
				'@type': 'Course',
				position: index + 1,
				name: course.title,
				description: course.title,
				provider: {
					'@type': 'Organization',
					name: course.mentor_name
				},
				aggregateRating: {
					'@type': 'AggregateRating',
					ratingValue: '4.5',
					ratingCount: '100'
				}
			}))
		}
	};

	const mentorsSchema = {
		'@context': 'https://schema.org',
		'@type': 'ItemList',
		itemListElement: data.mentors.entries.map((mentor, index) => ({
			'@type': 'Person',
			position: index + 1,
			name: mentor.name,
			jobTitle: 'Online Mentor',
			image: mentor.profile_image
		}))
	};
</script>

<svelte:head>
	<title>Home - Privat Unmei | Learn from Expert Mentors</title>
	<meta name="description" content="Discover and learn from expert mentors on Privat Unmei. Browse top-rated courses and connect with professionals." />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<meta property="og:title" content="Privat Unmei - Learn from Expert Mentors" />
	<meta property="og:description" content="Discover and learn from expert mentors on Privat Unmei. Browse top-rated courses and connect with professionals." />
	<meta property="og:type" content="website" />
	<meta name="keywords" content="online courses, mentors, learning, education, Privat Unmei" />
	<meta name="author" content="Privat Unmei" />
	<script type="application/ld+json">
		{JSON.stringify(coursesSchema)}
	</script>
	<script type="application/ld+json">
		{JSON.stringify(mentorsSchema)}
	</script>
</svelte:head>

<main class="space-y-8 px-4 py-6 sm:px-6 md:py-8">
	<!-- Most Bought Section -->
	<section class="space-y-4">
		<h2 class="text-xl font-bold text-(--tertiary-color) sm:text-2xl">Most Bought Courses</h2>
		<ScrollArea orientation="vertical" viewportClasses="h-48 w-full sm:h-56">
			<ul class="flex flex-col gap-3 sm:gap-4 md:grid md:grid-cols-2 lg:grid-cols-3 pr-4">
				{#each data.mostBought as c (c.id)}
					<li>
						<Link href={`/courses/${c.id}`}>
							<article
								class="group flex h-24 flex-col justify-between rounded-lg bg-(--tertiary-color) p-3 transition-all duration-200 hover:-translate-y-1 hover:shadow-md sm:p-4"
							>
								<div class="min-w-0 flex-1">
									<h3 class="truncate font-bold text-(--primary-color) text-sm sm:text-base">{c.title}</h3>
									<p class="truncate text-xs text-(--secondary-color) sm:text-sm">{c.mentor_name}</p>
								</div>
								<p class="text-xs font-semibold text-(--secondary-color) sm:text-sm">
									{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(
										c.price
									)}
								</p>
							</article>
						</Link>
					</li>
				{/each}
			</ul>
		</ScrollArea>
	</section>

	<!-- Courses Section -->
	<section class="space-y-4">
		<div class="flex flex-col items-start justify-between gap-3 sm:flex-row sm:items-center">
			<h2 class="text-xl font-bold text-(--tertiary-color) sm:text-2xl">Browse Courses</h2>
			<Link href="/courses">
				<span class="inline-flex items-center justify-center rounded-lg bg-(--tertiary-color) px-3 py-2 text-xs font-semibold transition-transform hover:scale-105 sm:px-4 sm:py-2 sm:text-sm">
					View All
				</span>
			</Link>
		</div>
		<ScrollArea orientation="vertical" viewportClasses="h-48 w-full sm:h-56">
			<ul class="flex flex-col gap-3 sm:gap-4 md:grid md:grid-cols-2 lg:grid-cols-3 pr-4">
				{#each data.courses.entries as c (c.id)}
					<li>
						<Link href={`/courses/${c.id}`}>
							<article
								class="group flex h-24 flex-col justify-between rounded-lg bg-(--tertiary-color) p-3 transition-all duration-200 hover:-translate-y-1 hover:shadow-md sm:p-4"
							>
								<div class="min-w-0 flex-1">
									<h3 class="truncate font-bold text-(--primary-color) text-sm sm:text-base">{c.title}</h3>
									<p class="truncate text-xs text-(--secondary-color) sm:text-sm">{c.mentor_name}</p>
								</div>
								<p class="text-xs font-semibold text-(--secondary-color) sm:text-sm">
									{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(
										c.price
									)}
								</p>
							</article>
						</Link>
					</li>
				{/each}
			</ul>
		</ScrollArea>
	</section>

	<!-- Mentors Section -->
	<section class="space-y-4">
		<div class="flex flex-col items-start justify-between gap-3 sm:flex-row sm:items-center">
			<h2 class="text-xl font-bold text-(--tertiary-color) sm:text-2xl">Expert Mentors</h2>
			<Link href="/mentors">
				<span class="inline-flex items-center justify-center rounded-lg bg-(--tertiary-color) px-3 py-2 text-xs font-semibold transition-transform hover:scale-105 sm:px-4 sm:py-2 sm:text-sm">
					View All
				</span>
			</Link>
		</div>
		<ScrollArea orientation="vertical" viewportClasses="h-48 w-full sm:h-56">
			<ul class="flex flex-col gap-3 sm:gap-4 md:grid md:grid-cols-2 lg:grid-cols-3 pr-4">
				{#each data.mentors.entries as m (m.id)}
					<li>
						<Link href={`/mentors/${m.id}`}>
							<article
								class="flex h-24 gap-3 rounded-lg bg-(--tertiary-color) p-3 transition-all duration-200 hover:-translate-y-1 hover:shadow-md sm:gap-4 sm:p-4"
							>
								<div class="shrink-0">
									<CldImage
										src={m.profile_image}
										width={70}
										height={70}
										className="rounded-full object-cover"
									/>
								</div>
								<div class="min-w-0 flex-1">
									<h3 class="truncate font-bold text-(--primary-color) text-sm sm:text-base">{m.name}</h3>
									<p class="truncate text-xs text-(--secondary-color) sm:text-sm">{m.public_id}</p>
								</div>
							</article>
						</Link>
					</li>
				{/each}
			</ul>
		</ScrollArea>
	</section>
</main>
