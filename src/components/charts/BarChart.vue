<script setup>
    import { ref, onMounted } from 'vue';
    import Chart from 'chart.js/auto';
    import BaseCard from "./BaseCard.vue";
    
    defineProps({
      width: {
        type: String,
        default: "50vw"
      },
      height: {
        type: String,
        default: "auto"
      }
    });

    const canvasRef = ref(null);
    const chartId = `chart-${Math.random().toString(36).substring(2, 8)}`;

    onMounted(() => {
      const ctx = canvasRef.value?.getContext('2d');
      if (!ctx) return;

      new Chart(ctx, {
        type: 'bar',
        data: {
          labels: ['Red', 'Blue', 'Yellow', 'Green', 'Purple', 'Orange'],
          datasets: [{
            label: '# of Votes',
            data: [12, 19, 3, 5, 2, 3],
            backgroundColor: [
              '#48897b', '#48897b', '#48897b',
              '#48897b', '#48897b', '#48897b'
            ],
            borderRadius: 6,
            borderSkipped: false,
          }]
        },
        options: {
          responsive: true,
          maintainAspectRatio: false,
          color: '#a0a0a0',
          plugins: {
            legend: {
              labels: {
                color: '#a0a0a0',
                font: {
                  family: '"Space Grotesk", sans-serif',
                  size: 12,
                },
                usePointStyle: true,
                pointStyleWidth: 8,
              }
            },
            tooltip: {
              backgroundColor: '#202024',
              titleColor: '#ffffff',
              bodyColor: '#c0c0c0',
              borderColor: '#333',
              borderWidth: 1,
              titleFont: {
                family: '"Space Grotesk", sans-serif',
              },
              bodyFont: {
                family: '"Space Grotesk", sans-serif',
              },
            }
          },
          scales: {
            x: {
              ticks: {
                color: '#6d6d6d',
                font: {
                  family: '"Space Grotesk", sans-serif',
                  size: 11,
                }
              },
              grid: {
                color: 'transparent',
              }
            },
            y: {
              beginAtZero: true,
              ticks: {
                color: '#6d6d6d',
                font: {
                  family: '"Space Grotesk", sans-serif',
                  size: 11,
                }
              },
              grid: {
                color: 'rgba(255, 255, 255, 0.06)',
              }
            }
          }
        }
      });
    });
</script>

<template>
    <BaseCard :title="title" :width="width" :height="height" :subtitle="subtitle">
        <div class="chart-container">
            <canvas :ref="el => canvasRef = el" :id="chartId"></canvas>
        </div>
    </BaseCard>
</template>

<style scoped>
.chart-container {
    position: relative;
    width: 100%;
    height: 100%;
    /* min-height: 300px;  Minimum height to ensure chart is visible */
}

canvas {
    width: 100% !important;
    height: 100% !important;
    display: block;
}
</style>