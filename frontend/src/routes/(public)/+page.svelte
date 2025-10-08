<script lang="ts">
	import Link from '$lib/components/button/Link.svelte';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { PrivatUnmeiLogo } from '$lib/utils/constants';
	import { Book, CalendarFold, Pencil } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { View } from './view.svelte';

	onMount(() => {
		View.isDesktop = window.innerWidth >= 768;
		function SetIsDesktop() {
			View.isDesktop = window.innerWidth >= 768;
		}
		window.addEventListener('resize', SetIsDesktop);
		return () => {
			window.removeEventListener('resize', SetIsDesktop);
		};
	});
</script>

<div class="flex h-full flex-col gap-8 p-4">
	<div class="flex flex-col gap-4 py-4">
		<CldImage
			width={View.imageWidth}
			height={View.imageHeight}
			src={PrivatUnmeiLogo}
			className="mx-auto"
		/>
		<p class="text-center">
			Lorem ipsum dolor, sit amet consectetur adipisicing elit. Amet odio nesciunt sequi voluptatum
			rerum quae. Est ex voluptatem nobis odit, quis eos obcaecati cum consectetur aspernatur sit
			nisi impedit at!
		</p>

		<div class="mx-auto flex justify-center rounded-2xl bg-[var(--tertiary-color)] p-4">
			<Link href="/courses" theme="light">Start Exploring</Link>
		</div>
	</div>
	<div class="flex flex-col gap-4">
		<h3 class="text-center text-2xl font-medium text-[var(--tertiary-color)]">Lorem ipsum</h3>
		<div class="flex flex-col gap-8 md:flex-row">
			<div class="flex flex-col gap-4">
				<Book />
				<h4 class="text-xl text-[var(--tertiary-color)]">lorem ipsum</h4>
				<p class="text-justify">
					Lorem, ipsum dolor sit amet consectetur adipisicing elit. Dolores deleniti, obcaecati
					maxime cupiditate, facere iure et nam, quidem itaque magni in rerum cumque earum provident
					mollitia animi ducimus expedita possimus.
				</p>
			</div>
			<div class="flex flex-col gap-4">
				<Pencil />
				<h4 class="text-xl text-[var(--tertiary-color)]">lorem ipsum</h4>
				<p class="text-justify">
					Lorem, ipsum dolor sit amet consectetur adipisicing elit. Dolores deleniti, obcaecati
					maxime cupiditate, facere iure et nam, quidem itaque magni in rerum cumque earum provident
					mollitia animi ducimus expedita possimus.
				</p>
			</div>
			<div class="flex flex-col gap-4">
				<CalendarFold />
				<h4 class="text-xl text-[var(--tertiary-color)]">lorem ipsum</h4>
				<p class="text-justify">
					Lorem, ipsum dolor sit amet consectetur adipisicing elit. Dolores deleniti, obcaecati
					maxime cupiditate, facere iure et nam, quidem itaque magni in rerum cumque earum provident
					mollitia animi ducimus expedita possimus.
				</p>
			</div>
		</div>
	</div>
</div>
