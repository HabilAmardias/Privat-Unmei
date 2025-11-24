<script lang="ts">
	import Link from '$lib/components/button/Link.svelte';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { PrivatUnmeiLogo } from '$lib/utils/constants';
	import { Book, CalendarFold, Pencil } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { HomeView } from './view.svelte';
	import Accordion from '$lib/components/accordion/Accordion.svelte';
	import { faqItems } from './model';
	import AnimatedContainer from '$lib/components/container/AnimatedContainer.svelte';

	const View = new HomeView();

	onMount(() => {
		View.setIsDesktop(window.innerWidth >= 768);
		function SetIsDesktop() {
			View.setIsDesktop(window.innerWidth >= 768);
		}
		window.addEventListener('resize', SetIsDesktop);
		return () => {
			window.removeEventListener('resize', SetIsDesktop);
		};
	});
</script>

<svelte:head>
	<title>Home - Privat Unmei</title>
	<meta name="description" content="Home - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex h-full flex-col gap-16 p-8">
	<AnimatedContainer>
		<div class="flex flex-col gap-4 py-4">
			<div class="mx-auto">
				<CldImage
					width={View.imageWidth}
					height={View.imageHeight}
					src={PrivatUnmeiLogo}
					className="mx-auto"
				/>
			</div>
			<p class="text-center">
				Lorem ipsum dolor, sit amet consectetur adipisicing elit. Amet odio nesciunt sequi
				voluptatum rerum quae. Est ex voluptatem nobis odit, quis eos obcaecati cum consectetur
				aspernatur sit nisi impedit at!
			</p>

			<div class="flex justify-center gap-4">
				<div class="flex justify-center rounded-2xl bg-[var(--tertiary-color)] p-4">
					<Link href="/home" theme="light">Start Exploring</Link>
				</div>
				<div class="flex justify-center rounded-2xl bg-[var(--tertiary-color)] p-4">
					<Link href="/login" theme="light">Sign In</Link>
				</div>
			</div>
		</div>
	</AnimatedContainer>
	<div class="flex flex-col gap-4">
		<AnimatedContainer>
			<h3 class="text-center text-2xl font-medium text-[var(--tertiary-color)]">Lorem ipsum</h3>
		</AnimatedContainer>
		<div class="flex flex-col gap-8 md:flex-row">
			<AnimatedContainer>
				<div class="flex flex-col gap-4">
					<Book size={View.iconsSize} class="text-[var(--tertiary-color)]" />
					<h4 class="text-xl text-[var(--tertiary-color)]">lorem ipsum</h4>
					<p class="text-justify">
						Lorem, ipsum dolor sit amet consectetur adipisicing elit. Dolores deleniti, obcaecati
						maxime cupiditate, facere iure et nam, quidem itaque magni in rerum cumque earum
						provident mollitia animi ducimus expedita possimus.
					</p>
				</div>
			</AnimatedContainer>
			<AnimatedContainer>
				<div class="flex flex-col gap-4">
					<Pencil size={View.iconsSize} class="text-[var(--tertiary-color)]" />
					<h4 class="text-xl text-[var(--tertiary-color)]">lorem ipsum</h4>
					<p class="text-justify">
						Lorem, ipsum dolor sit amet consectetur adipisicing elit. Dolores deleniti, obcaecati
						maxime cupiditate, facere iure et nam, quidem itaque magni in rerum cumque earum
						provident mollitia animi ducimus expedita possimus.
					</p>
				</div>
			</AnimatedContainer>
			<AnimatedContainer>
				<div class="flex flex-col gap-4">
					<CalendarFold size={View.iconsSize} class="text-[var(--tertiary-color)]" />
					<h4 class="text-xl text-[var(--tertiary-color)]">lorem ipsum</h4>
					<p class="text-justify">
						Lorem, ipsum dolor sit amet consectetur adipisicing elit. Dolores deleniti, obcaecati
						maxime cupiditate, facere iure et nam, quidem itaque magni in rerum cumque earum
						provident mollitia animi ducimus expedita possimus.
					</p>
				</div>
			</AnimatedContainer>
		</div>
	</div>
	<AnimatedContainer>
		<div class="flex flex-col gap-4">
			<h3 class="text-center text-2xl font-medium text-[var(--tertiary-color)]">
				Frequently Asked Question (FAQ)
			</h3>
			<Accordion type="multiple" items={faqItems} />
		</div>
	</AnimatedContainer>
</div>
