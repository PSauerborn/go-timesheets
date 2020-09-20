<template>
    <v-card class="day-card">
        <v-row>
            <v-col cols=6>
                <v-card-title class="day-card-title">This is a days data</v-card-title>
            </v-col>
        </v-row>
        <v-divider class="mx-4"></v-divider>
        <v-row align="center" justify="center">
            <v-col cols=6 align="center" justify="center">
                <v-card-text>day content</v-card-text>
            </v-col>
        </v-row>
    </v-card>
</template>

<script>
export default {
    name: "DayCard",
    props: {

    }
}
</script>

<style scoped>

.day-card {
    margin-bottom: 30px;
}

.day-card-title {
    font-size: 18px;
}

</style>