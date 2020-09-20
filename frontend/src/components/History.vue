<template>
    <v-row align="center" justify="center" class="application-tab-container">
        <v-row align="center" justify="center">
            <v-col cols=4>
                <DayCard v-for="(day, index) in periods" :key="index" />
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
            return moment(this.start).format('YYYY-MM-DD')
        },
        endDate: function() {
            return moment(this.end).format('YYYY-MM-DD')
        }
    },
    methods: {

        getPeriods: function() {
            const url = process.env.VUE_APP_BACKEND_URL + '/data'
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
                vm.periods = response.data.data.workPeriods
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
        start: "2020-09-13",
        end: "2020-09-21",
        periods: []
    }),
    mounted() {
        this.getPeriods()
    }
}
</script>