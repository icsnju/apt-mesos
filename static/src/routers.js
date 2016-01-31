'use strict'

export default function(router){
	router.map({
		'/': {
			name: 'dashboard',
			component: require('./views/dashboard.vue')
		},
        '*': {
            component: require('./views/dashboard.vue')
        },		
	})
}