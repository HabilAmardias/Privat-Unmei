<script lang="ts">
	import Link from '$lib/components/button/Link.svelte';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { PrivatUnmeiLogo } from '$lib/utils/constants';
	import { Book, CalendarFold, Pencil } from '@lucide/svelte';
	import { onMount } from 'svelte';
	import { HomeView } from './view.svelte';
	import Accordion from '$lib/components/accordion/Accordion.svelte';
	import { faqItems } from './constants';
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

	const organizationSchema = {
		'@context': 'https://schema.org',
		'@type': 'EducationalOrganization',
		name: 'Privat Unmei',
		description: 'Find the right private tutor for your learning goals. Connect with qualified tutors and book courses online, offline, or hybrid with ease.',
		url: 'https://privat-unmei.com',
		logo: PrivatUnmeiLogo
	};

	const webPageSchema = {
		'@context': 'https://schema.org',
		'@type': 'WebPage',
		name: 'Privat Unmei - Find Your Perfect Tutor',
		description: 'Connect with qualified private tutors for online, offline, or hybrid lessons. Easy booking, flexible schedules, and transparent pricing.',
		url: 'https://privat-unmei.com',
		publisher: {
			'@type': 'Organization',
			name: 'Privat Unmei'
		}
	};

	const faqSchema = {
		'@context': 'https://schema.org',
		'@type': 'FAQPage',
		mainEntity: faqItems.map((item) => ({
			'@type': 'Question',
			name: item.header,
			acceptedAnswer: {
				'@type': 'Answer',
				text: item.content
			}
		}))
	};
</script>

<svelte:head>
	<title>Privat Unmei - Find Your Perfect Private Tutor | Online & Offline Lessons</title>
	<meta name="description" content="Find the right private tutor for your learning goals. Browse qualified tutors, book courses online, offline, or hybrid with flexible schedules and transparent pricing on Privat Unmei." />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	<meta name="keywords" content="private tutor, online tutoring, offline tutoring, hybrid learning, private lessons, tuition, education, Privat Unmei" />
	<meta name="author" content="Privat Unmei" />
	<meta property="og:title" content="Privat Unmei - Find Your Perfect Private Tutor" />
	<meta property="og:description" content="Connect with qualified private tutors for personalized learning. Book online, offline, or hybrid lessons with flexible schedules." />
	<meta property="og:type" content="website" />
	<meta property="og:url" content="https://privat-unmei.com" />
	<meta property="og:image" content="{PrivatUnmeiLogo}" />
	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:title" content="Privat Unmei - Find Your Perfect Private Tutor" />
	<meta name="twitter:description" content="Connect with qualified private tutors for personalized learning. Book online, offline, or hybrid lessons." />
	<meta name="twitter:image" content="{PrivatUnmeiLogo}" />
	<script type="application/ld+json">
		{JSON.stringify(organizationSchema)}
	</script>
	<script type="application/ld+json">
		{JSON.stringify(webPageSchema)}
	</script>
	<script type="application/ld+json">
		{JSON.stringify(faqSchema)}
	</script>
</svelte:head>

<div class="flex h-full flex-col gap-16 p-8">
	<AnimatedContainer>
		<header class="flex flex-col gap-4 py-4">
			<div class="mx-auto">
				<CldImage
					width={View.imageWidth}
					height={View.imageHeight}
					src={PrivatUnmeiLogo}
					className="mx-auto"
					alt="Privat Unmei - Find Your Perfect Private Tutor"
				/>
			</div>
			<p class="text-center">
				Find the right private tutor for your learning goals, anytime, anywhere. Browse qualified
				tutors, courses, compare profiles, and book courses online, offline or hybrid with ease.
			</p>

			<div class="flex justify-center gap-4">
				<div class="flex justify-center rounded-2xl bg-(--tertiary-color) p-4">
					<Link href="/home">Start Exploring</Link>
				</div>
				<div class="flex justify-center rounded-2xl bg-(--tertiary-color) p-4">
					<Link href="/login">Sign In</Link>
				</div>
			</div>
		</header>
	</AnimatedContainer>
	<section class="flex flex-col gap-4">
		<AnimatedContainer>
			<h1 class="text-center text-2xl font-medium text-(--tertiary-color)">
				Why Choose Privat Unmei?
			</h1>
		</AnimatedContainer>
		<div class="flex flex-col gap-8 md:flex-row">
			<AnimatedContainer>
				<article class="flex flex-col gap-4">
					<Book size={View.iconsSize} class="text-(--tertiary-color)" />
					<h2 class="text-xl text-(--tertiary-color)">Find Courses and Tutors Easily</h2>
					<p class="text-justify">
						Search courses by subject and lesson type. Discover tutors that match your learning
						needs in just a few clicks.
					</p>
				</article>
			</AnimatedContainer>
			<AnimatedContainer>
				<article class="flex flex-col gap-4">
					<Pencil size={View.iconsSize} class="text-(--tertiary-color)" />
					<h2 class="text-xl text-(--tertiary-color)">Flexible Learning</h2>
					<p class="text-justify">
						Choose between online, offline or hybrid lessons based on your preference. Learn at your
						own pace with schedules that fit your routine.
					</p>
				</article>
			</AnimatedContainer>
			<AnimatedContainer>
				<article class="flex flex-col gap-4">
					<CalendarFold size={View.iconsSize} class="text-(--tertiary-color)" />
					<h2 class="text-xl text-(--tertiary-color)">Simple Booking</h2>
					<p class="text-justify">
						Book lessons directly through the app with clear pricing and schedules. No complicated
						process, just choose, book, and learn.
					</p>
				</article>
			</AnimatedContainer>
		</div>
	</section>
	<section class="flex flex-col gap-4">
		<AnimatedContainer>
			<h2 class="text-center text-2xl font-medium text-(--tertiary-color)">
				Frequently Asked Questions
			</h2>
		</AnimatedContainer>
		<AnimatedContainer>
			<Accordion type="multiple" items={faqItems} />
		</AnimatedContainer>
	</section>
</div>
