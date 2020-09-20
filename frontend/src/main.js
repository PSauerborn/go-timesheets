import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';
import Notifications from 'vue-notification'
import VueApexCharts from 'vue-apexcharts'

Vue.use(Notifications)
Vue.use(VueApexCharts)
Vue.component('apexchart', VueApexCharts)


Vue.config.productionTip = false

new Vue({
  vuetify,
  render: h => h(App)
}).$mount('#app')
