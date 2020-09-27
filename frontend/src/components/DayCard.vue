<template>

    <v-card class="day-card" max-width="800">
        <v-expansion-panels v-model="panel" :disabled=disabled>
            <v-expansion-panel>
                <v-expansion-panel-header>
                    <v-row dense>
                        <v-col cols=6 align="left" justify="left">
                            <v-card-title class="day-card-title">
                                Period Summary
                            </v-card-title>
                            <v-card-subtitle class="day-card-title">
                                Date: {{ formattedBucket }}
                            </v-card-subtitle>
                        </v-col>
                        <v-col cols=2 align="right" justify="right">
                            <v-row class="header-metric" align="center" justify="center">
                                {{ start }}
                            </v-row>
                            <v-row class="metric-text-box" align="center" justify="center">
                                Start
                            </v-row>
                        </v-col>
                        <v-col cols=2 align="right" justify="right">
                            <v-row class="header-metric" align="center" justify="center">
                                {{ end }}
                            </v-row>
                            <v-row class="header-metric-text-box" align="center" justify="center">
                                End
                            </v-row>
                        </v-col>
                        <v-col cols=2 align="right" justify="right">
                            <v-row class="header-metric" align="center" justify="center">
                                {{ netWorkHours }}
                            </v-row>
                            <v-row class="header-metric-text-box" align="center" justify="center">
                                Net Work Hours
                            </v-row>
                        </v-col>
                    </v-row>
                </v-expansion-panel-header>
                <v-expansion-panel-content>
                    <v-divider class="mx-4"></v-divider>
                    <v-row dense>
                        <v-col cols=6 align="left" justify="left">
                            <v-card-text>
                                <apexchart type="bar" :options="chartOptions" :series="chartData"/>
                            </v-card-text>
                        </v-col>
                        <v-divider :vertical=true></v-divider>
                        <v-col cols=5 align="center" justify="center" class="metric-container">
                            <v-card-text>
                                <v-row align="center" justify="center">
                                    <v-col cols=12 align="center" justify="center">
                                        <v-row class="metric" align="center" justify="center">
                                            {{ payload.totalPeriods }}
                                        </v-row>
                                        <v-row class="metric-text-box" align="center" justify="center">
                                            Total Work Periods
                                        </v-row>
                                    </v-col>
                                </v-row>
                                <v-row align="center" justify="center">
                                    <v-col cols=12 align="center" justify="center">
                                        <v-row class="metric" align="center" justify="center">
                                            {{ breakCount }}
                                        </v-row>
                                        <v-row class="metric-text-box" align="center" justify="center">
                                            Total Break Periods
                                        </v-row>
                                    </v-col>
                                </v-row>
                                <v-row align="center" justify="center">
                                    <v-col cols=12 align="center" justify="center">
                                        <v-row class="metric" align="center" justify="center">
                                            {{ averageBreakLength }}
                                        </v-row>
                                        <v-row class="metric-text-box" align="center" justify="center">
                                            Average Break Length
                                        </v-row>
                                    </v-col>
                                </v-row>
                            </v-card-text>
                        </v-col>
                    </v-row>
                </v-expansion-panel-content>
            </v-expansion-panel>
        </v-expansion-panels>
    </v-card>

</template>

<script>

import moment from 'moment';

export default {
    name: "DayCard",
    props: {
        payload: {
            type: Object,
        },
        bucket: {
            type: String
        }
    },
    computed: {
        formattedBucket: function() {
            return moment(this.bucket).format("YYYY-MM-DD")
        },
        hoursWorked: function() {
            return Math.round(this.payload.totalWorkHours * 10) / 10
        },
        breakCount: function() {
            return this.payload.totalBreaks
        },
        breakHours: function() {
            return Math.round(this.payload.totalBreakHours * 10) / 10
        },
        averageBreakLength: function() {
            if (this.breakCount > 0) {
                return this.breakHours / this.breakCount
            }
            return 0
        },
        netWorkHours: function() {
            return Math.round(this.payload.netWorkHours * 10) / 10
        },
        chartData: function() {
            return [{name: this.payload.date, data: [this.hoursWorked, this.breakHours, this.netWorkHours]}]
        },
        start: function() {
            return moment(this.payload.startTime).format('HH:mm')
        },
        end: function() {
            return moment(this.payload.endTime).format('HH:mm')
        }
    },
    data: () => ({
        panel: [0,1],
        disabled: false,
        chartOptions: {
            chart: {
                id: 'Period Statistics',
                toolbar: {
                    show: false
                }
            },
            colors: ['#F44336', '#2196F3', '#4CAF50'],
            xaxis: {
                categories: ['Total', 'Break', 'Net']
            },
            title: {
                text: 'hours'
            },
            plotOptions: {
                bar: {
                    horizontal: true,
                    distributed: true
                }
            },
            grid: {
                show: false
            }
        }
    })
}
</script>

<style scoped>

@import url('https://fonts.googleapis.com/css?family=Allura&display=swap');

.day-card {
    margin-bottom: 30px;
}

.metric-text-box {
    font-size: 12px;
    text-transform: uppercase;
}

.metric {
    font-size: 55px;
    font-weight: bold;
    font-family: 'Allura', 'Avenir', Helvetica, Arial, sans-serif;
    margin-bottom: 10px;
}

.header-metric-text-box {
    font-size: 12px;
    text-transform: uppercase;
}

.header-metric {
    margin-top: 15px;
    font-size: 30px;
    font-weight: bold;
    font-family: 'Allura', 'Avenir', Helvetica, Arial, sans-serif;
    margin-bottom: 10px;
}

.metric-container {
    padding-top: 18px;
}

</style>