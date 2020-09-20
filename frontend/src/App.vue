<template>
  <v-app>
    <v-main>
      <notifications group="main" position="bottom right"/>
      <v-container>
        <v-row align="center" justify="center" class="application-title-container">
          <AppTitle />
        </v-row>
        <v-row align="center" justify="center">
          <v-tabs :centered=true>
            <v-tab class="application-tab-title"><v-icon class="tab-icon">mdi-blur</v-icon>Current</v-tab>
            <v-tab-item>
              <CurrentPeriod />
            </v-tab-item>
            <v-tab class="application-tab-title"><v-icon class="tab-icon">mdi-format-list-text</v-icon>History</v-tab>
            <v-tab-item>
              <History />
            </v-tab-item>
            <v-tab class="application-tab-title"><v-icon class="tab-icon">mdi-chart-donut-variant</v-icon>Analytics</v-tab>
            <v-tab-item>
              <Analytics />
            </v-tab-item>
          </v-tabs>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>

import AppTitle from './components/AppTitle';
import Analytics from './components/Analytics';
import History from './components/History';
import CurrentPeriod from './components/CurrentPeriod';
import shared from './shared';

export default {
  name: 'App',

  components: {
    AppTitle,
    Analytics,
    History,
    CurrentPeriod
  },

  data: () => ({
    //
  }),
  beforeMount: function() {
      console.log(process.env.VUE_APP_DEVELOPMENT_MODE)
      if (process.env.VUE_APP_DEVELOPMENT_MODE === "true") {
        shared.setDevToken()
      }
      shared.getAccessToken()
    }
};
</script>

<style>

.application-tab-container {
  margin: 30px;
}

.application-tab-title {
  margin-left: 20px;
  margin-right: 20px;
}

.application-title-container {
  margin: 20px;
}

.tab-icon {
  margin-right: 10px;
}
</style>