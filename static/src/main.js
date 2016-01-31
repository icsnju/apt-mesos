import	Vue from 'vue'
import VueRouter from 'vue-router'
import VueResource from 'vue-resource'
import routerMap from './routers'
import './assets/scss/app.scss'

Vue.use(VueRouter)
Vue.use(VueResource)

let router = new VueRouter({
    hashbang: true,
    history: false,
    saveScrollPosition: true,
    transitionOnLoad: true
});

let app = Vue.extend({});

routerMap(router);

router.start(app, "#app");