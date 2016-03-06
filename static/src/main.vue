<template>
	<div>
		<sh-sidebar></sh-sidebar>
		<div class="main-wrapper" id="main-wrapper">
			<sh-header></sh-header>
		    <div class="wrapper wrapper-content">
				<router-view transition="fade" transition-mode="out-in" keep-alive></router-view>
			</div>
		</div>
	</div>
</template>

<script>
import Sidebar from './components/sidebar.vue'
import Header from './components/header.vue'

export default {
	data() {
		
	},
    components:{
        "shSidebar": Sidebar,
        "shHeader": Header
    },
    ready() {
    }
}
</script>