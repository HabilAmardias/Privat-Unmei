<script lang="ts">
	import { DatePicker } from 'bits-ui';
	type datePickerProps = {
		label: string;
	};
	let { label }: datePickerProps = $props();
</script>

<DatePicker.Root>
	<DatePicker.Label class="block text-sm">{label}</DatePicker.Label>
	<DatePicker.Input>
		{#snippet children({ segments })}
			{#each segments as { part, value }}
				<DatePicker.Segment {part}>
					{value}
				</DatePicker.Segment>
			{/each}
		{/snippet}
	</DatePicker.Input>
	<DatePicker.Content>
		<DatePicker.Calendar>
			{#snippet children({ months, weekdays })}
				<DatePicker.Header>
					<DatePicker.PrevButton />
					<DatePicker.Heading />
					<DatePicker.NextButton />
				</DatePicker.Header>
				{#each months as month}
					<DatePicker.Grid>
						<DatePicker.GridHead>
							<DatePicker.GridRow>
								{#each weekdays as day}
									<DatePicker.HeadCell>
										{day}
									</DatePicker.HeadCell>
								{/each}
							</DatePicker.GridRow>
						</DatePicker.GridHead>
						<DatePicker.GridBody>
							{#each month.weeks as weekDates}
								<DatePicker.GridRow>
									{#each weekDates as date}
										<DatePicker.Cell {date} month={month.value}>
											<DatePicker.Day />
										</DatePicker.Cell>
									{/each}
								</DatePicker.GridRow>
							{/each}
						</DatePicker.GridBody>
					</DatePicker.Grid>
				{/each}
			{/snippet}
		</DatePicker.Calendar>
	</DatePicker.Content>
</DatePicker.Root>
