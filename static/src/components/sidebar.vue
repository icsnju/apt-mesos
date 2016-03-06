<template>
	<div class="sidebar">
		<div class="sidebar-inner" id="side-menu">
			<div class="brand-top">
				<div class="brand">
                    <a href="">
					<img src="../assets/images/logo.png" alt="">
                    <h3>APT-MESOS</h3></a>
				</div>			
			</div>
			<div class="navbar">
                <a class="navbar-minimalize"><span class="navbar-minimalize-inner"><span class="si-menu"></span></span></a>
				<sh-nav :items="items"></sh-nav>
			</div>                			
		</div>	
	</div>
</template>

<script>
import Nav from '../components/nav.vue'

export default {
	data() {
		let items = [
		    {
	        	path: '/dashboard',
	        	icon: 'cloud',
	        	label: 'dashboard'
	      	}, {
	        	path: '/task',
	        	icon: 'browser',
	        	label: 'task'
	      	}, {
	        	path: '/cluster',
	        	icon: 'server',
	        	label: 'cluster'
	      	}
		]
		return {
			items: items
		}
	},
	components: {
		'shNav': Nav
	}
}
</script>
