<template>
  <v-app>
    <v-main>
      <v-container>
        <v-row align="center" justify="center" class="application-title-container">
          <AppTitle />
        </v-row>
        <v-row align="center" justify="center">
          <v-tabs centered="true">
            <v-tab class="application-tab-title">Current</v-tab>
            <v-tab-item>
              <CurrentPeriod />
            </v-tab-item>
            <v-tab class="application-tab-title">History</v-tab>
            <v-tab-item>
              <History />
            </v-tab-item>
            <v-tab class="application-tab-title">Analytics</v-tab>
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

</style>