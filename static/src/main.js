import	Vue from 'vue'
import VueRouter from 'vue-router'
import routerMap from './routers'

Vue.use(VueRouter);

let router = new VueRouter({
    hashbang: true,
    history: false,
    saveScrollPosition: true,
    transitionOnLoad: true
});

let app = Vue.extend({});

routerMap(router);

router.start(app, "#app");