<template>
    <v-container class="logo-container">
        <v-row align="center" justify="center">
            <v-col cols=4 align="center" justify="center">
                <h1>Welcome to Go-Timesheets!</h1>
            </v-col>
            <v-col cols=2 align="left" justify="left">
                <v-img :contain=true max-height="200" max-width="200" src="../assets/logo.png"/>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
export default {
    name: "AppTitle"
}
</script>

<style scoped>

@import url('https://fonts.googleapis.com/css?family=Allura&display=swap');

.logo-container {
    font-family: 'Allura', 'Avenir', Helvetica, Arial, sans-serif;
}
</style>