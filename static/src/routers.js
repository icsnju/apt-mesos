'use strict'

export default function(router){
	router.map({
		'/': {
			name: 'dashboard',
			component: require('./views/index.vue')
		},
        '*': {
            component: require('./views/index.vue')
        },		
	})
}