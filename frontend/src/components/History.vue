<template>
    <v-container align="center" justify="center" class="application-tab-container">
        <v-row align="center" justify="center" style="margin-bottom: 20px;" dense>
            <v-col cols=2 align="center" justify="center" class="overview-cols">
                <v-menu :close-on-content-click="false" v-model="dateSelectorOpen" offset-y>
                <template v-slot:activator="{ on }">
                    <v-btn v-on="on" color="blue" class="date-button" :outlined=true :large=true>{{ range[0] }} - {{ range[1] }}</v-btn>
                </template>
                    <v-date-picker v-model='range' range/>
                </v-menu>
            </v-col>
        </v-row>
        <v-row v-if="sortedPeriods.length < 1" dense>
            <v-col cols=12 align="center" justify="center">
                No historical data found. Adjust the Date range or Completed some work Periods
            </v-col>
        </v-row>

        <v-row align="center" justify="center" style="margin-bottom: 20px;" v-if="sortedPeriods.length > 0" dense>
            <v-col cols=1 align="center" justify="center" class="overview-cols">
                <v-row align="center" justify="center" class="metric" dense>
                    {{ totalWorkHours.workedHours }}
                </v-row>
                <v-row align="center" justify="center" class="metric-text-box" dense>
                    total hours worked
                </v-row>
            </v-col>
            <v-col cols=1 align="center" justify="center" class="overview-cols">
                <v-row align="center" justify="center" class="metric" dense>
                    {{ totalWorkHours.breakHours }}
                </v-row>
                <v-row align="center" justify="center" class="metric-text-box" dense>
                    total break hours
                </v-row>
            </v-col>
            <v-col cols=1 align="center" justify="center" class="overview-cols">
                <v-row align="center" justify="center" class="metric" dense>
                    {{ totalWorkHours.netWorkHours }}
                </v-row>
                <v-row align="center" justify="center" class="metric-text-box" dense>
                    net work hours
                </v-row>
            </v-col>
        </v-row>
        <v-row align="center" justify="center">
            <v-col cols=12 align="center" justify="center">
                <DayCard v-for="(day, index) in sortedPeriods" :key="index" :payload="day" />
            </v-col>
        </v-row>
    </v-container>
</template>

<script>

import DayCard from './DayCard';
import axios from 'axios';
import shared from '../shared';
import moment from 'moment';

export default {
    name: "History",
    components: {
        DayCard
    },
    computed: {
        sortedPeriods: function() {
            var values = []
            Object.keys(this.periods).forEach((date) => {
                if (this.periods[date].length > 0) {
                    values.push({date: date, periods: this.periods[date]})
                }
            })
            return values.sort(function(a, b) {
                return moment(b.date) - moment(a.date)
            })
        },
        totalWorkHours: function() {
            var worked = 0;
            var breaks = 0;
            const periodList = Object.values(this.periods)

            if (periodList.length < 1) {
                return {workedHours: 0, breakHours: 0, netWorkHours: 0}
            }

            periodList.forEach((periods) => {
                periods.forEach((period) => {
                    const timespan = moment.duration((moment(period.finishedAt).diff(moment(period.createdAt))))
                    worked += timespan.asHours()
                    period.breaks.forEach((breakPeriod) => {
                        const timespan = moment.duration((moment(breakPeriod.finishedAt).diff(moment(breakPeriod.createdAt))))
                        breaks += timespan.asHours()
                    })
                })
            })
            return {
                workedHours: Math.round(worked * 10) / 10,
                breakHours: Math.round(breaks * 10) / 10,
                netWorkHours: Math.round((worked - breaks) * 10) / 10
            }

        }
    },
    methods: {
        getPeriods: function() {
            const url = process.env.VUE_APP_BACKEND_URL + `/data/${this.range[0]}/${this.range[1]}?group=true`
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
                    text: 'successfully retrieved historical data'
                })
                // asign payload to variable
                vm.periods = response.data.data
            }).catch(function (error) {
                console.log("error fetching active work period: API return status code " + error.response.status)
                if (error.response.status === 401) {
                    window.location.replace(process.env.VUE_APP_LOGIN_REDIRECT)
                }  else {
                    vm.$notify({
                        group: 'main',
                        title: 'go-timesheets backend',
                        type: 'error',
                        text: 'unable retrieving historical data'
                    })
                }
            })
        },
        onDateRangeChange() {
            console.log("date range changed")
        }
    },
    data: () => ({
        periods: {},
        range: [
            moment().subtract('days', 7).format('YYYY-MM-DD'),
            moment().add('days', 1).format('YYYY-MM-DD')
        ],
        dateSelectorOpen: false
    }),
    mounted() {
        this.getPeriods()
    },
    watch: {
        dateSelectorOpen: function() {
            if (!this.dateSelectorOpen) {
                this.getPeriods()
            }
        }
    }
}
</script>

<style scoped>

@import url('https://fonts.googleapis.com/css?family=Allura&display=swap');

.metric-text-box {
    font-size: 12px;
    text-transform: uppercase;
    padding-bottom: 0px;
    margin-bottom: 0px;
}

.metric {
    font-size: 55px;
    font-weight: bold;
    font-family: 'Allura', 'Avenir', Helvetica, Arial, sans-serif;
    margin-bottom: 0px;
}

.overview-cols {
    margin-left: 20px;
    margin-right: 20px;
}

.date-button {
    color: white;
}

</style>