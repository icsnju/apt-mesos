import Home from './views/home.vue'

export default function (router){
	router.map({
	    '/home': {
	      component: Home
	    },
	})

    router.redirect({
		'/': '/home'
  	})
}