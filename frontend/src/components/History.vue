<template>
    <v-row align="center" justify="center" class="application-tab-container">
        <v-row align="center" justify="center">
            <v-col align="center" justify="center" v-if="sortedPeriods.length < 1">
                No historical data found. Completed some work periods to see historic data
            </v-col>
            <v-col cols=12 align="center" justify="center">
                <DayCard v-for="(day, index) in sortedPeriods" :key="index" :payload="day" />
            </v-col>
        </v-row>
    </v-row>
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
        startDate: function() {
            const timestamp = moment()
            return timestamp.subtract('days', 7).format('YYYY-MM-DD')
        },
        endDate: function() {
            const timestamp = moment()
            return timestamp.add('days', 1).format('YYYY-MM-DD')
        },
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
        }
    },
    methods: {

        getPeriods: function() {
            const url = process.env.VUE_APP_BACKEND_URL + `/data/${this.startDate}/${this.endDate}?group=true`
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
        }
    },
    data: () => ({
        periods: {}
    }),
    mounted() {
        this.getPeriods()
    }
}
</script>