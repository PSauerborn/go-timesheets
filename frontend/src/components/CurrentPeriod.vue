<template>
    <v-row align="center" justify="center" class="application-tab-container">
        <v-card v-if="inActivePeriod" tile min-width="700">
            <v-row align="center" justify="center" dense>
                <v-col cols=6 align="center" justify="center">
                    <v-card-title>
                        Current Work Period
                    </v-card-title>
                    <v-card-subtitle>
                        {{ this.activePeriod.periodId }}
                    </v-card-subtitle>
                </v-col>
                <v-col cols=5 align="center" justify="center">
                    <v-card-text>
                        <v-btn v-if="!inBreak" color="blue" :outlined=true @click="startNewBreakPeriod">Start Break Period</v-btn>
                        <v-btn v-if="inBreak" color="orange" :outlined=true @click="stopCurrentBreakPeriod">Finish Break Period</v-btn>
                    </v-card-text>
                    <v-card-text>
                        <v-btn color="red" :outlined=true @click="endCurrentWorkPeriod" :disabled="inBreak">Stop Active Work Period</v-btn>
                    </v-card-text>
                </v-col>
            </v-row>
            <v-divider class="mx-4"></v-divider>
            <v-row align="center" justify="center" dense>
                <v-col cols=6 align="center" justify="center">
                    <v-card-text class="active-period-text">
                        Active Since: {{ createdTimestamp }}
                    </v-card-text>
                </v-col>
                <v-col cols=6 align="center" justify="center">
                    <v-card-text class="active-period-text">
                        {{ counters.hours }}H {{ counters.minutes }}M {{ counters.seconds }}S
                    </v-card-text>
                </v-col>
            </v-row>
        </v-card>
        <v-card v-if="!inActivePeriod" :flat=true>
            <v-card-text>
                No active period Detected. Start a new Work Period Below
            </v-card-text>
            <v-divider class="mx-4"></v-divider>
            <v-card-text align="center" justify="center">
                <v-btn :outlined=true color="blue" :dark=true @click="startNewWorkPeriod">Start New Period</v-btn>
            </v-card-text>
        </v-card>
    </v-row>
</template>

<script>

import axios from 'axios';
import moment from 'moment';
import shared from '../shared';


export default {
    name: "CurrentPeriod",
    computed: {
        inActivePeriod: function() {
            return this.activePeriod != null
        },
        inBreak: function() {
            return this.activeBreak != null
        },
        createdTimestamp: function() {
            return moment(String(this.activePeriod.createdAt)).format('HH:mm:ss DD/MM/YYYY')
        }
    },
    methods: {
        startNewBreakPeriod: function() {
            const url = process.env.VUE_APP_BACKEND_URL + '/break_period/' + this.activePeriod.periodId
            let vm = this

            axios({
                method: 'post',
                url: url,
                headers: {'Authorization': 'Bearer ' + shared.getAccessToken()}
            }).then(function (response) {
                vm.$notify({
                    group: 'main',
                    title: 'go-timesheets backend',
                    type: 'success',
                    text: 'successfully started new break period'
                })
                // asign payload to variable
                vm.activeBreak = response.data.payload

            }).catch(function (error) {
                console.log("error fetching active work period: API return status code " + error.response.status)
                if (error.response.status === 401) {
                    window.location.replace(process.env.VUE_APP_LOGIN_REDIRECT)
                } else {
                    vm.$notify({
                        group: 'main',
                        title: 'go-timesheets backend',
                        type: 'error',
                        text: 'unable to start new break period'
                    })
                }
            })
        },
        stopCurrentBreakPeriod: function() {
            const url = process.env.VUE_APP_BACKEND_URL + '/break_period/' + this.activeBreak.breakId
            let vm = this

            axios({
                method: 'patch',
                url: url,
                headers: {'Authorization': 'Bearer ' + shared.getAccessToken()}
            }).then(function (response) {
                console.log(response)
                vm.$notify({
                    group: 'main',
                    title: 'go-timesheets backend',
                    type: 'success',
                    text: 'successfully closed current break period'
                })
                // asign payload to variable
                vm.activeBreak = null

            }).catch(function (error) {
                console.log("error fetching active work period: API return status code " + error.response.status)
                if (error.response.status === 401) {
                    window.location.replace(process.env.VUE_APP_LOGIN_REDIRECT)
                } else {
                    vm.$notify({
                        group: 'main',
                        title: 'go-timesheets backend',
                        type: 'error',
                        text: 'unable to stop current break period'
                    })
                }
            })
        },
        increaseCounter: function() {
            let vm = this;
            this.timer = setInterval(() => {
                if (vm.counters.seconds == 59) {
                    if (vm.counters.minutes == 59) {
                        vm.counters.hours += 1
                        vm.counters.minutes = 0
                    } else {
                        vm.counters.minutes += 1
                    }
                    vm.counters.seconds = 0
                } else {
                    vm.counters.seconds += 1
                }
            }, 1000)
        },
        setTimeCounters: function() {
            const timestamp = moment(String(this.activePeriod.createdAt))
            const now = moment()

            const diff = now.diff(timestamp, 'seconds')
            const hours = diff / 3600
            const minutes = (hours - Math.floor(hours)) * 60
            const seconds = (minutes - Math.floor(minutes)) * 60

            this.counters.hours = Math.floor(hours)
            this.counters.minutes = Math.floor(minutes)
            this.counters.seconds = Math.floor(seconds)
        },
        /**
         * Function used to retrieve active work period from backend.
         * Note that the backend returns a 404 if no active work period
         * is detected
         */
        getActivePeriod: function() {
            const url = process.env.VUE_APP_BACKEND_URL + '/active'
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
                    text: 'successfully retrieved active work period'
                })
                // asign payload to variable
                vm.activePeriod = response.data.payload
                vm.activeBreak = response.data.payload.activeBreak
                vm.setTimeCounters()
                if (vm.timer != null) {
                    clearInterval(vm.timer)
                }
                vm.increaseCounter()
            }).catch(function (error) {
                console.log("error fetching active work period: API return status code " + error.response.status)
                if (error.response.status === 401) {
                    window.location.replace(process.env.VUE_APP_LOGIN_REDIRECT)
                } else if (error.response.status === 404) {
                    vm.$notify({
                        group: 'main',
                        title: 'go-timesheets backend',
                        type: 'warn',
                        text: 'no active work period found'
                    })
                    vm.activePeriod = null
                } else {
                    vm.$notify({
                        group: 'main',
                        title: 'go-timesheets backend',
                        type: 'error',
                        text: 'unable retrieving active work period'
                    })
                }
            })
        },
        /**
         * Function used to start new active work period
         */
        startNewWorkPeriod: function() {
            const url = process.env.VUE_APP_BACKEND_URL + '/work_period'
            let vm = this

            axios({
                method: 'post',
                url: url,
                headers: {'Authorization': 'Bearer ' + shared.getAccessToken()}
            }).then(function (response) {
                vm.$notify({
                    group: 'main',
                    title: 'go-timesheets backend',
                    type: 'success',
                    text: 'successfully started new work period'
                })
                // asign payload to variable
                vm.activePeriod = response.data.payload
                vm.setTimeCounters()
                if (vm.timer != null) {
                    clearInterval(vm.timer)
                }
                vm.increaseCounter()
            }).catch(function (error) {
                console.log("error fetching active work period: API return status code " + error.response.status)
                if (error.response.status === 401) {
                    window.location.replace(process.env.VUE_APP_LOGIN_REDIRECT)
                } else {
                    vm.$notify({
                        group: 'main',
                        title: 'go-timesheets backend',
                        type: 'error',
                        text: 'unable to start new work period work period'
                    })
                }
            })
        },
        endCurrentWorkPeriod: function() {
            const url = process.env.VUE_APP_BACKEND_URL + '/work_period/' + this.activePeriod.periodId
            let vm = this

            axios({
                method: 'patch',
                url: url,
                headers: {'Authorization': 'Bearer ' + shared.getAccessToken()}
            }).then(function (response) {
                console.log(response)
                vm.$notify({
                    group: 'main',
                    title: 'go-timesheets backend',
                    type: 'success',
                    text: 'successfully ended work period'
                })
                vm.activePeriod = null
                vm.$emit('endedWorkPeriod')
            }).catch(function (error) {
                console.log("error fetching active work period: API return status code " + error.response.status)
                if (error.response.status === 401) {
                    window.location.replace(process.env.VUE_APP_LOGIN_REDIRECT)
                } else {
                    vm.$notify({
                        group: 'main',
                        title: 'go-timesheets backend',
                        type: 'error',
                        text: 'unable to end current work period'
                    })
                }
            })
        }
    },
    data: () => ({
        activePeriod: null,
        activeBreak: null,
        counters: {seconds: 0, minutes: 0, hours: 0},
        timer: null
    }),
    mounted() {
        this.getActivePeriod()
    }
}
</script>

<style scoped>

.active-period-text {
    font-size: 16px;
}

</style>