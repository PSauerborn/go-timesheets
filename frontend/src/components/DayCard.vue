<template>
    <v-card class="day-card" max-width="700">
        <v-row>
            <v-col cols=6 align="center" justify="center">
                <v-card-title class="day-card-title">
                    Date: {{ payload.date }}
                </v-card-title>
                <v-card-subtitle>
                    Total Work Periods: {{ payload.periods.length }}
                </v-card-subtitle>
            </v-col>
        </v-row>
        <v-divider class="mx-4"></v-divider>
        <v-row align="center" justify="center">
            <v-col cols=6 align="center" justify="center">
                <v-card-text></v-card-text>
            </v-col>
        </v-row>
    </v-card>
</template>

<script>
export default {
    name: "DayCard",
    props: {
        payload: Object
    }
}
</script>

<style scoped>

.day-card {
    margin-bottom: 30px;
}


</style>