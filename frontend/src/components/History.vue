<template>
    <v-container align="center" justify="center" class="application-tab-container">
        <v-row align="center" justify="center" style="margin-bottom: 20px;" dense>
            <v-col cols=9 align="center" justify="center">
                <date-selector v-model="dateRange" @dateChanged="getBucketAnalysis" />
            </v-col>
        </v-row>
        <v-row v-if="Object.keys(buckets).length < 1" dense>
            <v-col cols=12 align="center" justify="center">
                No historical data found. Adjust the Date range or Complete some work Periods
            </v-col>
        </v-row>

        <v-row align="center" justify="center" style="margin-bottom: 20px;" v-if="Object.keys(buckets).length > 0" dense>
            <v-col cols=1 align="center" justify="center" class="overview-cols">
                <v-row align="center" justify="center" class="metric" dense>
                    {{ overview.workedHours }}
                </v-row>
                <v-row align="center" justify="center" class="metric-text-box" dense>
                    total hours worked
                </v-row>
            </v-col>
            <v-col cols=1 align="center" justify="center" class="overview-cols">
                <v-row align="center" justify="center" class="metric" dense>
                    {{ overview.breakHours }}
                </v-row>
                <v-row align="center" justify="center" class="metric-text-box" dense>
                    total break hours
                </v-row>
            </v-col>
            <v-col cols=1 align="center" justify="center" class="overview-cols">
                <v-row align="center" justify="center" class="metric" dense>
                    {{ overview.netWorkHours }}
                </v-row>
                <v-row align="center" justify="center" class="metric-text-box" dense>
                    net work hours
                </v-row>
            </v-col>
        </v-row>
        <v-row align="center" justify="center">
            <v-col cols=12 align="center" justify="center">
                <DayCard v-for="(payload, timestamp) in buckets" :key="timestamp" :payload="payload" :bucket="timestamp"/>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>

import DayCard from './DayCard';
import axios from 'axios';
import shared from '../shared';
import moment from 'moment';
import DateSelector from './shared/DateSelector'

export default {
    name: "History",
    components: {
        DayCard,
        DateSelector
    },
    computed: {
        /**
         * Computed property used to evaluate the total
         * number of worked hours, break hours and the
         * net work hours given the results returned from
         * the bucket analysis
         */
        overview: function() {
            var worked = 0;
            var breaks = 0;
            const buckets = Object.values(this.buckets)

            // return values of 0 is no buckets are found
            if (buckets.length < 1) {
                return {workedHours: 0, breakHours: 0, netWorkHours: 0}
            }
            // iterate over buckets and increment work and break hours
            buckets.forEach((bucket) => {
                worked += bucket.totalWorkHours
                breaks += bucket.totalBreakHours
            })
            // round all values to 1 decimal place and return
            return {
                workedHours: Math.round(worked * 10) / 10,
                breakHours: Math.round(breaks * 10) / 10,
                netWorkHours: Math.round((worked - breaks) * 10) / 10
            }
        },
        /**
         * Computed propert used to evaluate start time of bucket
         * analysis
         */
        startTimestamp: function() {
            return moment(this.dateRange.start).format('YYYY-MM-DDTHH:mm')
        },
        /**
         * Computed propert used to evaluate end time of bucket
         * analysis
         */
        endTimestamp: function() {
            return moment(this.dateRange.end).format('YYYY-MM-DDTHH:mm')
        }
    },
    methods: {
        /**
         * Function used to retrieve bucked analysis from backend database.
         * The bucketed analysis returns an object with {date: values} format,
         * where values contain the total work hours, net work hours and break
         */
        getBucketAnalysis() {
            const url = process.env.VUE_APP_BACKEND_URL + `/bucket_analysis/${this.startTimestamp}/${this.endTimestamp}`
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
                vm.buckets = response.data.payload.buckets
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
        }
    },
    data: () => ({
        periods: {},
        buckets: {},
        dateRange: {
            start: moment().startOf('isoWeek'),
            end: moment().add('days', 1).format('YYYY-MM-DD')
        },
        dateSelectorOpen: false
    }),
    mounted() {
        this.getBucketAnalysis();
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

</style>