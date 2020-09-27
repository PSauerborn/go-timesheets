<template>
    <v-container fluid>
        <v-card>
            <v-row dense>
                <v-col>
                    <v-card-title>Overview</v-card-title>
                    <v-card-subtitle>View general working trends and statistics</v-card-subtitle>
                </v-col>
            </v-row>
            <v-divider class="mx-4"></v-divider>
            <v-card-text>
                <v-row align="center" justify="center" dense>
                    <v-col cols=7 align="center" justify="center">
                        <apexchart type="line" :options="chartOptions" :series="chartData" />
                    </v-col>
                    <v-col cols=4 align="center" justify="center">
                        <v-row align="center" justify="center">
                            <v-col class="value-col">
                                {{ averageWorkHours }}
                            </v-col>
                        </v-row>
                        <v-row align="center" justify="center">
                            <v-col class="description-col">
                                average work hours
                            </v-col>
                        </v-row>
                        <v-row align="center" justify="center">
                            <v-col class="value-col">
                                {{ averageBreakHours }}
                            </v-col>
                        </v-row>
                        <v-row align="center" justify="center">
                            <v-col class="description-col">
                                average break hours
                            </v-col>
                        </v-row>
                        <v-row align="center" justify="center">
                            <v-col class="value-col">
                                {{ totalWorkHours }}
                            </v-col>
                        </v-row>
                        <v-row align="center" justify="center">
                            <v-col class="description-col">
                                total work hours
                            </v-col>
                        </v-row>
                    </v-col>
                </v-row>
            </v-card-text>
        </v-card>
    </v-container>
</template>

<script>

import axios from 'axios';
import shared from '../../shared';
import moment from 'moment';

export default {
    name: "Overview",
    props: {
        start: {
            type: String
        },
        end: {
            type: String
        }
    },
    computed: {
        startDate() {
            return moment(this.start).format('YYYY-MM-DDTHH:mm')
        },
        endDate() {
            return moment(this.end).format('YYYY-MM-DDTHH:mm')
        },
        averageWorkHours() {
            return this.safeDivide(Math.round(this.overview.averageBucketWorkHours * 10), 10)
        },
        averageBreakHours() {
            return this.safeDivide(Math.round(this.overview.averageBucketBreakHours * 10), 10)
        },
        totalWorkHours() {
            return this.safeDivide(Math.round(this.overview.totalWorkHours * 10), 10)
        },
        chartData() {
            let netWorkHours = []
            let breakHours = []
            let totalWorkHours = []
            const buckets = Object.values(this.buckets)
            buckets.forEach((bucket) => {
                netWorkHours.push(this.safeDivide(Math.round(bucket.netWorkHours * 10), 10))
                breakHours.push(this.safeDivide(Math.round(bucket.totalBreakHours * 10), 10))
                totalWorkHours.push(this.safeDivide(Math.round(bucket.totalWorkHours * 10), 10))
            })
            return [
                {
                    name: "Net Work Hours",
                    data: netWorkHours
                },
                {
                    name: "Break Hours",
                    data: breakHours
                },
                {
                    name: "Total Work Hours",
                    data: totalWorkHours
                }
            ]
        },
        chartTimestamps() {
            return Object.keys(this.buckets)
        },
        chartOptions() {
            return {
                chart: {
                    type: 'area',
                    id: 'Overview Stats',
                    stacked: false,
                    toolbar: {show: false},
                },
                fill: {
                    type: 'gradient',
                    gradient: {
                        shadeIntensity: 1,
                        inverseColors: false,
                        opacityFrom: 0.5,
                        opacityTo: 0,
                        stops: [0, 90, 100]
                    }
                },
                stroke: {
                    show: true,
                    curve: 'smooth',
                    lineCap: 'butt',
                    colors: undefined,
                    width: 2,
                    dashArray: 0,
                },
                colors: ['#F44336', '#2196F3', '#4CAF50'],
                xaxis: {
                    type: 'datetime',
                    categories: this.chartTimestamps
                },
                yaxis: {
                    title: {text: 'Hours'}
                },
                title: {
                    text: 'Overview Statistics'
                },
                plotOptions: {

                },
                grid: {show: false},
                legend: {offsetY: 20, position: 'top'}
            }
        }
    },
    methods: {
        getData() {
            const url = process.env.VUE_APP_BACKEND_URL + `/bucket_analysis/${this.startDate}/${this.endDate}?include_empty=true`
            let vm = this

            axios({
                method: 'get',
                url: url,
                headers: {'Authorization': 'Bearer ' + shared.getAccessToken()}
            }).then(function (response) {
                vm.$notify({
                    group: 'main',
                    title: 'go-timesheets backend',
                    type: 'success',
                    text: 'successfully retrieved overview data'
                })
                // asign payload to variable
                vm.buckets = response.data.payload.buckets
                vm.overview = response.data.payload.overview

            }).catch(function (error) {
                console.log("error fetching active work period: API return status code " + error.response.status)
                if (error.response.status === 401) {
                    window.location.replace(process.env.VUE_APP_LOGIN_REDIRECT)
                } else {
                    vm.$notify({
                        group: 'main',
                        title: 'go-timesheets backend',
                        type: 'error',
                        text: 'unable to retrieve overview data'
                    })
                }
            })
        },
        safeDivide(a, b) {
            if (b === 0) {
                return 0
            }
            return a / b
        }
    },
    data: () => ({
        buckets: {},
        overview: {}
    }),
    mounted() {
        this.getData()
    }
}
</script>

<style scoped>

@import url('https://fonts.googleapis.com/css?family=Allura&display=swap');

.value-col {
    font-size: 55px;
    font-weight: bold;
    font-family: 'Allura', 'Avenir', Helvetica, Arial, sans-serif;
}

.description-col {
    font-size: 12px;
    text-transform: uppercase;
    margin-bottom: 10px;
}
</style>