'use strict'

import Vue from 'vue'
import VueRouter from 'vue-router'
import VueResource from 'vue-resource'
import validator from 'vue-validator'
import filters from './filters'
import routerMap from './routers'
import './assets/css/style.min.css'
import './assets/css/app.min.css'
import app from './main.vue'

Vue.use(VueResource);
Vue.use(VueRouter);
Vue.use(validator);
Vue.config.debug = true

//实例化VueRouter
let router = new VueRouter({
    hashbang: true,
    history: false,
    saveScrollPosition: true,
    transitionOnLoad: true
});

routerMap(router);

router.start(app, "#app");

toastr.options = {
  "newestOnTop": true,
  "progressBar": true,
  "positionClass": "toast-top-right",
  "timeOut": "3000",
  "showEasing": "swing",
  "hideEasing": "linear",
}