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
                                Date: {{ payload.date }}
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
                                            {{ payload.periods.length }}
                                        </v-row>
                                        <v-row class="metric-text-box" align="center" justify="center">
                                            Total Work Periods
                                        </v-row>
                                    </v-col>
                                </v-row>
                                <v-row align="center" justify="center">
                                    <v-col cols=12 align="center" justify="center">
                                        <v-row class="metric" align="center" justify="center">
                                            {{ start }}
                                        </v-row>
                                        <v-row class="metric-text-box" align="center" justify="center">
                                            Start Time
                                        </v-row>
                                    </v-col>
                                </v-row>
                                <v-row align="center" justify="center">
                                    <v-col cols=12 align="center" justify="center">
                                        <v-row class="metric" align="center" justify="center">
                                            {{ end }}
                                        </v-row>
                                        <v-row class="metric-text-box" align="center" justify="center">
                                            End Time
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
            default: function() {
                return {
                    periods: []
                }
            }
        }
    },
    computed: {
        hoursWorked: function() {
            var value = 0
            this.payload.periods.forEach((period) => {
                const timespan = moment.duration((moment(period.finishedAt).diff(moment(period.createdAt))))
                value += timespan.asHours()
            })
            return Math.round(value * 10) / 10
        },
        breakHours: function() {
            var value = 0
            const periods = this.payload.periods;
            periods.forEach((period) => {
                period.breaks.forEach((breakPeriod) => {
                    const timespan = moment.duration((moment(breakPeriod.finishedAt).diff(moment(breakPeriod.createdAt))))
                    value += timespan.asHours()
                })
            })
            return Math.round(value * 10) / 10
        },
        netWorkHours: function() {
            return this.hoursWorked - this.breakHours
        },
        chartData: function() {
            return [{name: this.payload.date, data: [this.hoursWorked, this.breakHours, this.netWorkHours]}]
        },
        start: function() {
            if (this.payload.periods.length > 0) {
                const dateString = this.payload.periods[0].createdAt
                return moment(dateString).format('HH:mm')
            }
            return null
        },
        end: function() {
            if (this.payload.periods.length > 0) {
                const index = this.payload.periods.length
                const dateString = this.payload.periods[index - 1].finishedAt
                return moment(dateString).format('HH:mm')
            }
            return null
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