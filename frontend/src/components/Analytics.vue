<template>
    <v-container align="center" justify="center" class="application-tab-container">
        <v-row align="center" justify="center" style="margin-bottom: 20px;" dense>
            <v-col cols=4 align="center" justify="center">
                <date-selector v-model="dateRange" />
            </v-col>
        </v-row>
        <v-row align="center" justify="center">
            <v-col cols=5 align="center" justify="center" >
                <overview />
            </v-col>
            <v-col cols=3 align="center" justify="center">
                <trends />
            </v-col>
        </v-row>
    </v-container>
</template>

<script>

import Overview from './analytics/Overview.vue';
import Trends from './analytics/Trends.vue';
import DateSelector from './shared/DateSelector.vue';

export default {
    name: "Analytics",
    components: {
        Overview,
        DateSelector,
        Trends
    },
    data: () => ({
        dateSelectorOpen: false,
        dateRange: {},
        chartOptions: {

        }
    })
}
</script>