<script lang="ts">
	import Button from '$lib/components/button/Button.svelte';
	import ScrollArea from '$lib/components/scrollarea/ScrollArea.svelte';

	type TACProps = {
		onClick: () => void;
	};
	let { onClick }: TACProps = $props();
</script>

<ScrollArea orientation="vertical" viewportClasses="h-[400px] max-h-[400px]">
	<div class="flex flex-col gap-4">
		<section>
			<h2 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">1. Acceptance of Terms</h2>
			<p>
				By registering for and using this application (“Platform”), you agree to be bound by these
				Terms and Conditions (“Terms”). If you do not agree, please do not use or register on the
				Platform.
			</p>
		</section>
		<section>
			<h2 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">2. Platform Description</h2>
			<p>
				The Platform acts as a marketplace that connects students (“Users”) with tutors (“Tutors”)
				for private lessons, either online or offline. The Platform does not provide tutoring
				services directly and is not a party to any agreement between Users and Tutors.
			</p>
		</section>
		<section>
			<h2 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">3. User Eligibility</h2>
			<p>By registering, you confirm that:</p>
			<ul class="list-inside list-disc">
				<li>You are at least 18 years old, or have parental/legal guardian consent</li>
				<li>The information you provide is accurate, complete, and up to date</li>
				<li>You will use the Platform for lawful purposes only</li>
			</ul>
		</section>
		<section>
			<h2 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">
				4. Account Registration and Security
			</h2>
			<ul class="list-inside list-disc">
				<li>You are responsible for maintaining the confidentiality of your account credentials</li>
				<li>You are fully responsible for all activities conducted through your account</li>
				<li>The Platform is not liable for unauthorized access resulting from your negligence</li>
			</ul>
		</section>
		<section>
			<h2 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">
				5. Tutor and Lesson Information
			</h2>
			<ul class="list-inside list-disc">
				<li>Tutor profiles, availability, pricing, and lesson details are provided by Tutors</li>
				<li>The Platform does not guarantee the quality, outcome, or suitability of any Tutor</li>
				<li>Users are responsible for reviewing Tutor information before booking</li>
			</ul>
		</section>
		<section>
			<h2 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">
				6. Booking, Payments, and Cancellations
			</h2>
			<ul class="list-inside list-disc">
				<li>Lesson bookings are subject to Tutor availability</li>
				<li>Payments and fees are governed by the Platform’s payment rules</li>
				<li>
					The Platform reserves the right to withhold or adjust payments in cases of dispute, fraud,
					or policy violations
				</li>
			</ul>
		</section>
		<section>
			<h2 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">
				7. Offline Lessons Disclaimer
			</h2>
			<p>For offline (in-person) lessons:</p>
			<ul class="list-inside list-disc">
				<li>Users and Tutors are responsible for arranging a safe and lawful meeting location</li>
				<li>
					The Platform is not responsible for any incidents, damages, injuries, or losses occurring
					outside the Platform
				</li>
				<li>Users are advised to exercise caution when meeting in person</li>
			</ul>
		</section>
		<section>
			<h2 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">8. User Conduct</h2>
			<p>You agree not to:</p>
			<ul class="list-inside list-disc">
				<li>Provide false or misleading information</li>
				<li>Harass, abuse, or harm other users</li>
				<li>Circumvent the Platform’s payment or booking system</li>
				<li>Use the Platform for illegal or unauthorized activities</li>
			</ul>
		</section>
		<section>
			<h2 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">9. Reviews and Ratings</h2>
			<ul class="list-inside list-disc">
				<li>Users may provide honest reviews based on real experiences</li>
				<li>Reviews must not contain offensive, defamatory, or misleading content</li>
			</ul>
		</section>
		<section>
			<h2 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">
				10. Limitation of Liability
			</h2>
			<p>To the maximum extent permitted by law:</p>
			<ul class="list-inside list-disc">
				<li>
					The Platform is not liable for disputes, damages, or losses arising between Users and
					Tutors
				</li>
				<li>The Platform does not guarantee uninterrupted or error-free service</li>
			</ul>
		</section>
		<section>
			<h2 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">11. Termination</h2>
			<p>
				The Platform may suspend or terminate your access at any time if you violate these Terms or
				applicable laws.
			</p>
		</section>
		<section>
			<h2 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">12. Changes to Terms</h2>
			<p>
				The Platform may update these Terms from time to time. Continued use of the Platform after
				changes constitutes acceptance of the updated Terms.
			</p>
		</section>
		<section>
			<h2 class="mb-4 text-xl font-bold text-[var(--tertiary-color)]">13. Governing Law</h2>
			<p>
				These Terms are governed by and interpreted in accordance with the laws of the applicable
				jurisdiction.
			</p>
		</section>
		<Button full type="submit" {onClick}
			>By registering, I agree to the Terms and Conditions and Privacy Policy of this Platform.</Button
		>
	</div>
</ScrollArea>
