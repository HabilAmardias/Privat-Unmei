<script lang="ts">
	import LineChart from "$lib/components/charts/LineChart.svelte";
	import ScrollArea from "$lib/components/scrollarea/ScrollArea.svelte";
	import type { PageProps } from "./$types";
	import { AdminDashboardView } from "./view.svelte";

    let {data} : PageProps = $props()

    // svelte-ignore state_referenced_locally
    const View = new AdminDashboardView(data.historyReport, data.incomeReport, data.mentorReport)
</script>

<div class="h-dvh flex flex-col gap-8 p-4">
    <div class="flex justify-center gap-8">
        <div class="text-center">
            <h3 class="text-[var(--tertiary-color)] font-bold text-xl">This Month Session:</h3>
            <h3 class="text-[var(--tertiary-color)] font-bold text-xl">{View.totalSession}</h3>
        </div>
        <div class="text-center">
            <h3 class="text-[var(--tertiary-color)] font-bold text-xl">This Month Income:</h3>
            <h3 class="text-[var(--tertiary-color)] font-bold text-xl">{new Intl.NumberFormat("id-ID", {style:"currency",currency: "IDR"}).format(View.totalCost!)}</h3>
        </div>
    </div>
    <h2 class="text-[var(--tertiary-color)] font-bold text-2xl">This Month Mentor Report</h2>
    {#if View.mentorReports.length > 0}
        <ScrollArea orientation="vertical" class="flex-1" viewportClasses="max-h-[300px]">
            <table class="w-full table-fixed border-spacing-4 border-2 border-[var(--tertiary-color)] border-collapse">
                <thead>
                    <tr class="border-2 border-[var(--tertiary-color)] border-collapse" >
                        <td class="text-center font-bold text-[var(--tertiary-color)]">Name</td>
                        <td class="text-center font-bold text-[var(--tertiary-color)]">Email</td>
                        <td class="text-center font-bold text-[var(--tertiary-color)]">Total Session</td>
                        <td class="text-center font-bold text-[var(--tertiary-color)]">Total Operational Cost</td>
                    </tr>
                </thead>
                <tbody>
                    {#each View.mentorReports as mrp}
                        <tr class="border-2 border-[var(--tertiary-color)] border-collapse">
                            <td class="text-center">{mrp.name}</td>
                            <td class="text-center">{mrp.email}</td>
                            <td class="text-center">{mrp.total_session}</td>
                            <td class="text-center">{new Intl.NumberFormat("id-ID", {style:"currency", currency:"IDR"}).format(mrp.total_cost)}</td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </ScrollArea>
    {:else}
        <div class="w-full flex justify-center items-center h-75">
            <p class="text-[var(--tertiary-color)] font-bold">No Income</p>
        </div>
    {/if}
    <h2 class="text-[var(--tertiary-color)] font-bold text-2xl">This Year Report</h2>
    <div class="grid grid-rows-2 grid-cols-1 md:grid-cols-2 md:grid-rows-1 gap-4 w-full">
        {#if View.costValueHistoryReport.length > 0}
            <LineChart 
        data={$state.snapshot<number[]>(View.costValueHistoryReport)}
        labels={$state.snapshot<string[]>(View.costLabelHistoryReport)}
        graphLabel="This Year Income"
        color="#365486"
        />
        {:else}
        <div class="h-[300px] flex justify-center items-center">
            <p class="text-[var(--tertiary-color)] font-bold text-center">No Income</p>
        </div>
        {/if}
        {#if View.sessionValueHistoryReport.length > 0}
            <LineChart 
        data={$state.snapshot<number[]>(View.sessionValueHistoryReport)}
        labels={$state.snapshot<string[]>(View.sessionLabelHistoryReport)}
        graphLabel="This Year Session"
        color="#365486"
        />
        {:else}
        <div class="h-[300px] flex justify-center items-center">
            <p class="text-[var(--tertiary-color)] font-bold text-center">No Session</p>
        </div>
        {/if}
    </div>
</div>