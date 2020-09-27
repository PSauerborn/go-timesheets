<template>
    <v-container fluid>
        <v-menu :close-on-content-click="false" v-model="dateSelectorOpen" offset-y>
            <template v-slot:activator="{ on }">
                <v-card v-on="on" :max-width="400" :min-width="400" :max-height="100" :flat=true>
                    <v-row align="center" justify="center" dense>
                        <v-col align="center" justify="center" cols=5>
                            <v-row align="center" justify="center">
                                <v-col cols=3>
                                    <v-icon>mdi-clock-start</v-icon>
                                </v-col>
                                <v-col cols=7>
                                    <v-row class="date-title" dense>start</v-row>
                                    <v-row align="center" justify="center" class="date-value" dense>{{ startDate }}</v-row>
                                </v-col>
                            </v-row>
                        </v-col>
                        <v-divider :vertical=true class="button-divider"></v-divider>
                        <v-col align="center" justify="center" cols=5>
                            <v-row align="center" justify="center">
                                <v-col cols=3>
                                    <v-icon>mdi-clock-end</v-icon>
                                </v-col>
                                <v-col cols=7>
                                    <v-row class="date-title" dense>end</v-row>
                                    <v-row align="center" justify="center" class="date-value" dense>{{ endDate }}</v-row>
                                </v-col>
                            </v-row>
                        </v-col>
                    </v-row>
                </v-card>
            </template>
            <v-date-picker :width="400" v-model='range' @input="updateDateRange()" range/>
        </v-menu>
    </v-container>
</template>

<script>
import moment from 'moment';
export default {
    name: "DateSelector",
    computed: {
        startDate() {
            return moment(this.range[0]).format('YYYY-MM-DD')
        },
        endDate() {
            return moment(this.range[1]).format('YYYY-MM-DD')
        },
        buttonColor() {
            if (this.validDateRange) {
                return "blue"
            }
            return "red"
        }
    },
    methods: {
        updateDateRange() {
            this.$emit('input', {
                start: this.startDate,
                end: this.endDate
            })
        }
    },
    watch: {
        dateSelectorOpen: function() {
            if (!this.dateSelectorOpen) {
                if (this.startDate > this.endDate) {
                    this.validDateRange = false
                } else if (this.startDate != this.oldDates.start || this.endDate != this.oldDates.end) {
                    this.validDateRange = true
                    this.$emit('dateChanged')
                    // set old dates to cross reference on next change
                    this.oldDates.start = this.startDate
                    this.oldDates.end = this.endDate
                }
            }
        }
    },
    data: () => ({
        dateSelectorOpen: false,
        range: [
            moment().subtract('days', 7).format('YYYY-MM-DD'),
            moment().add('days', 1).format('YYYY-MM-DD')
        ],
        validDateRange: true,
        oldDates: {
            start: moment().subtract('days', 7).format('YYYY-MM-DD'),
            end: moment().add('days', 1).format('YYYY-MM-DD')
        }
    })
}
</script>

<style scoped>

.date-button {
    font-size: 20px;
}

.button-divider {
    margin-left: 5px;
    margin-right: 5px;
}

.date-title {
    text-transform: uppercase;
    font-size: 12px;
    font-weight: bold;
}

.date-value {
    font-size: 16px;
}

</style>