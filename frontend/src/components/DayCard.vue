<template>
    <v-card class="day-card" max-width="800">
        <v-row>
            <v-col cols=12>
                <v-card-title class="day-card-title">
                    Date: {{ payload.date }}
                </v-card-title>
            </v-col>
        </v-row>
        <v-divider class="mx-4"></v-divider>
        <v-row align="center" justify="center">
            <v-col cols=6 align="center" justify="center">
                <v-card-text>
                    <apexchart type="bar" :options="chartOptions" :series="chartData"/>
                </v-card-text>
            </v-col>
            <v-divider :vertical=true></v-divider>
            <v-col cols=4 align="left" justify="left">
                <v-card-title>
                    Period Summary
                </v-card-title>
                <v-divider></v-divider>
                <v-card-text>
                    Total Work Periods: {{ payload.periods.length }}<br>
                    Start Timestamp: {{ start }}<br>
                    End Timestamp: {{ end }}<br>
                </v-card-text>
            </v-col>
        </v-row>
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
        chartOptions: {
            chart: {
                id: 'Period Statistics',
                toolbar: {
                    show: false
                }
            },
            xaxis: {
                categories: ['Total', 'Break', 'Net']
            },
            title: {
                text: 'hours'
            },
            plotOptions: {
                bar: {
                    horizontal: true
                }
            }
        }
    })
}
</script>

<style scoped>

.day-card {
    margin-bottom: 30px;
}

</style>