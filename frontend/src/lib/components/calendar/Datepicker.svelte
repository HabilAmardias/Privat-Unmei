<script lang="ts">
	import { DatePicker } from 'bits-ui';
	import { Calendar, ChevronLeft, ChevronRight } from '@lucide/svelte';
	import type { DateValue } from '@internationalized/date';
	import { getDayOfWeek } from '@internationalized/date';
	type datePickerProps = {
		label?: string;
		dows: number[];
		onChange?: (date: DateValue | undefined) => void;
	};
	let { label, dows, onChange }: datePickerProps = $props();

	function disabledDates(date: DateValue) {
		return !dows.includes(getDayOfWeek(date, 'id-ID'));
	}
</script>

<DatePicker.Root
	onValueChange={onChange}
	isDateDisabled={disabledDates}
	weekdayFormat="short"
	fixedWeeks={true}
>
	<div class="flex w-full max-w-[232px] flex-col gap-1.5">
		{#if label}
			<DatePicker.Label class="block select-none text-sm font-medium">{label}</DatePicker.Label>
		{/if}
		<DatePicker.Input
			class="h-input rounded-input border-border-input bg-background text-foreground focus-within:border-border-input-hover focus-within:shadow-date-field-focus hover:border-border-input-hover flex w-full max-w-[232px] select-none items-center rounded-lg border bg-[var(--tertiary-color)] px-2 py-3 text-sm tracking-[0.01em] text-[var(--secondary-color)]"
		>
			{#snippet children({ segments })}
				{#each segments as { part, value }, i (part + i)}
					<div class="inline-block select-none">
						{#if part === 'literal'}
							<DatePicker.Segment {part} class="text-muted-foreground p-1">
								{value}
							</DatePicker.Segment>
						{:else}
							<DatePicker.Segment
								{part}
								class="rounded-5px hover:bg-muted focus:bg-muted focus:text-foreground aria-[valuetext=Empty]:text-muted-foreground focus-visible:ring-0! focus-visible:ring-offset-0! px-1 py-1"
							>
								{value}
							</DatePicker.Segment>
						{/if}
					</div>
				{/each}
				<DatePicker.Trigger
					class="text-foreground/60 hover:bg-muted active:bg-dark-10 ml-auto inline-flex size-8 items-center justify-center rounded-[5px] transition-all"
				>
					<Calendar class="cursor-pointer hover:text-[var(--primary-color)]" />
				</DatePicker.Trigger>
			{/snippet}
		</DatePicker.Input>
		<DatePicker.Content sideOffset={6} class="z-50">
			<DatePicker.Calendar
				class="border-dark-10 bg-background-alt shadow-popover rounded-[15px] border bg-[var(--tertiary-color)] p-[22px]"
			>
				{#snippet children({ months, weekdays })}
					<DatePicker.Header class="flex items-center justify-between">
						<DatePicker.PrevButton
							class="rounded-9px bg-background-alt hover:bg-muted inline-flex size-10 items-center justify-center transition-all active:scale-[0.98]"
						>
							<ChevronLeft class="cursor-pointer text-[var(--secondary-color)]" />
						</DatePicker.PrevButton>
						<DatePicker.Heading class="text-[15px] font-medium text-[var(--secondary-color)]" />
						<DatePicker.NextButton
							class="rounded-9px bg-background-alt hover:bg-muted inline-flex size-10 items-center justify-center transition-all active:scale-[0.98]"
						>
							<ChevronRight class="cursor-pointer text-[var(--secondary-color)]" />
						</DatePicker.NextButton>
					</DatePicker.Header>
					<div class="flex flex-col space-y-4 pt-4 sm:flex-row sm:space-x-4 sm:space-y-0">
						{#each months as month (month.value)}
							<DatePicker.Grid class="w-full border-collapse select-none space-y-1">
								<DatePicker.GridHead>
									<DatePicker.GridRow class="mb-1 flex w-full justify-between">
										{#each weekdays as day (day)}
											<DatePicker.HeadCell
												class="text-muted-foreground font-normal! w-10 rounded-md text-xs text-[var(--secondary-color)]"
											>
												<div>{day.slice(0, 2)}</div>
											</DatePicker.HeadCell>
										{/each}
									</DatePicker.GridRow>
								</DatePicker.GridHead>
								<DatePicker.GridBody>
									{#each month.weeks as weekDates (weekDates)}
										<DatePicker.GridRow class="flex w-full">
											{#each weekDates as date (date)}
												<DatePicker.Cell
													{date}
													month={month.value}
													class="p-0! relative size-10 text-center text-sm"
												>
													<DatePicker.Day
														class="rounded-9px text-foreground hover:border-foreground data-selected:bg-foreground data-disabled:text-foreground/30 data-selected:text-background data-unavailable:text-muted-foreground data-disabled:pointer-events-none data-outside-month:pointer-events-none data-selected:font-medium data-unavailable:line-through group relative inline-flex size-10 cursor-pointer items-center justify-center whitespace-nowrap border border-transparent bg-transparent p-0 text-sm font-normal text-[var(--secondary-color)] transition-all hover:text-[var(--primary-color)] data-[disabled]:text-[var(--dark-color)]"
													>
														<div
															class="bg-foreground group-data-selected:bg-background group-data-today:block absolute top-[5px] hidden size-1 rounded-full bg-[var(--primary-color)] transition-all"
														></div>
														{date.day}
													</DatePicker.Day>
												</DatePicker.Cell>
											{/each}
										</DatePicker.GridRow>
									{/each}
								</DatePicker.GridBody>
							</DatePicker.Grid>
						{/each}
					</div>
				{/snippet}
			</DatePicker.Calendar>
		</DatePicker.Content>
	</div>
</DatePicker.Root>
