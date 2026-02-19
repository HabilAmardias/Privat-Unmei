<script lang="ts">
    import { Chart } from "chart.js/auto";
    import { type ChartConfiguration, type ChartData, type ChartDataset } from "chart.js"
	import type { Attachment } from "svelte/attachments";

    type LineChartProps = {
        data: number[],
        labels: string[],
        graphLabel: string,
        color: string
    }

    let {data, labels, graphLabel, color} : LineChartProps = $props()

    const drawGraph : Attachment = (node) => {
        // svelte-ignore state_referenced_locally
        const lineDataset : ChartDataset = {
            data: data,
            borderColor: color,
            backgroundColor: color,
            label: graphLabel
        }
        // svelte-ignore state_referenced_locally
        const lineData : ChartData = {
            labels: labels,
            datasets: [lineDataset],
        }
        const lineChartConfig : ChartConfiguration = {
            type: "line",
            data: lineData,
            options: {
                responsive: true,
            },
            
        }
        let graph = new Chart(node as HTMLCanvasElement, lineChartConfig)

        return () => {
            graph.destroy()
        }
    }
</script>

<div class="w-full">
  <canvas {@attach drawGraph}></canvas>
</div>