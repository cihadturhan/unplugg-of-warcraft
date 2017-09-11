<template>
  <div class="chart-container">
    <LineChart
      :chart-data="chartData"
      :options="options"
      :width="900"
      :height="320"
    >
    </LineChart>
  </div>
</template>

<script>
  import LineChart from './LineChart';
  import output from '../model/output.json';
  import moment from 'moment';
  import Vue from 'vue';

  output.forecast.sort((a, b) => {
    return a.timestamp - b.timestamp;
  });

  export default {
    props: ['item'],
    watch: {
      item: function () {
        const newdata = output.forecast.map(d => {
          return {y: d.value - 10 + 40 * Math.random(), x: moment(d.timestamp * 1000).toDate()};
        });
        this.chartData = {
          datasets: [
            {
              label: 'Forecast for #1',
              borderColor: '#FC2525',
              pointBackgroundColor: 'transparent',
              borderWidth: 1,
              pointBorderColor: '#FC2525',
              data: newdata
            }]
        };
      },
    },
    data() {
      return {
        chartData: {
          datasets: [
            {
              label: 'Forecast for #1',
              borderColor: '#FC2525',
              pointBackgroundColor: 'transparent',
              borderWidth: 1,
              pointBorderColor: '#FC2525',
              data: output.forecast.map(d => {
                return {y: d.value, x: moment(d.timestamp * 1000).toDate()};
              }),
            }
          ]
        },
        options: {
          cubicInterpolationMode: 'monotone',
          borderDash: [2, 4],
          responsive: true,
          maintainAspectRatio: true,
          scales: {
            xAxes: [{
              type: "time",
              display: true,
              scaleLabel: {
                display: true,
                labelString: 'Date'
              }
            }],
          },

        }
      }
    },
    components: {
      LineChart
    }
  }
</script>

<style>
  .chart-container {
    width: 100%;
    padding: 100px 0 30px 0;
    background-color: var(--color-dark-blue-gray);
    display: flex;
    justify-content: center;
  }
</style>
