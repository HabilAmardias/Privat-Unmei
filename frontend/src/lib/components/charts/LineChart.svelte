<script lang="ts">
    import { Chart, type ChartConfiguration, type ChartData, type ChartDataset, type ChartType } from "chart.js";
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
            borderColor: color
        }
        // svelte-ignore state_referenced_locally
        const lineData : ChartData = {
            labels: labels,
            datasets: [lineDataset],
        }
        const lineChartConfig : ChartConfiguration = {
            type: 'line',
            data: lineData,
            options: {
                responsive: true,
                plugins: {
                    title: {
                        text: graphLabel
                    }
                }
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