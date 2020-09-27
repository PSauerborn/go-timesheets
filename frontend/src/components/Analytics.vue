<template>
    <v-container align="center" justify="center" class="application-tab-container">
        <v-row align="center" justify="center" dense>
            <v-col cols=4 align="center" justify="center">
                <date-selector v-model="dateRange" @dateChanged="updateData" />
            </v-col>
        </v-row>
        <v-row align="center" justify="center">
            <v-col cols=8 align="center" justify="center" >
                <overview ref="overview" :start="dateRange.start" :end="dateRange.end" />
            </v-col>
        </v-row>
    </v-container>
</template>

<script>

import Overview from './analytics/Overview.vue';
import DateSelector from './shared/DateSelector.vue';
import moment from 'moment';

export default {
    name: "Analytics",
    components: {
        Overview,
        DateSelector
    },
    methods: {
        updateData() {
            this.$refs.overview.getData()
        }
    },
    data: () => ({
        dateRange: {
            start: moment().subtract('days', 7).format('YYYY-MM-DD'),
            end: moment().add('days', 1).format('YYYY-MM-DD')
        },
        dateSelectorOpen: false,
        chartOptions: {

        }
    })
}
</script>