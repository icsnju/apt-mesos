<template>
		<top-header></top-header>
		<sidebar></sidebar>
		<offsidebar></offsidebar>
		<section>
			<div class="content-wrapper">
				<div class="container-fluid">
					<div class="row">
						<div class="col-md-4">
							<monitor title="CPU Allocation" :total="metrics.total_cpus" :used="metrics.used_cpus" unit="" :refresh="refresh_random"></monitor>
						</div>	
						<div class="col-md-4">
							<monitor title="Memory Allocation" :total="metrics.total_mem" :used="metrics.used_mem" unit="MB" :refresh="refresh_random"></monitor>
						</div>	
						<div class="col-md-4">
							<monitor title="Disk Allocation" :total="metrics.total_disk" :used="metrics.used_disk" unit="MB" :refresh="refresh_random"></monitor>
						</div>													
					</div>
				</div>
			</div>
		</section>
		<footer></footer>
</template>
<script>
    export default {
    	data() {
    		return {
    			metrics: {
	    			total_cpus: 0,
	    			total_disk: 0,
	    			total_mem: 0,
	    			used_cpus: 0,
	    			used_disk: 0,
	    			used_mem: 0    				
    			},
    			refresh_random: 0
    		}
    	},
    	ready() {
    		let self = this
    		setInterval(function(){
    			self.getMetrics()
    		},1000)
    	},
    	methods: {
	    	getMetrics() {
	    		this.$http({url: '/api/system/metrics', method: 'GET'}).then(function (response) {
	    			this.metrics = response.data.result
	    			this.refresh_random = Math.random()
			    }, function (response) {
			    	console.log(status)
			    })
	    	},
    	},
        components:{
            "top-header":require('../components/header.vue'),
            "sidebar": require('../components/sidebar.vue'),
            "offsidebar": require('../components/offsidebar.vue'),
            "monitor": require('../components/monitor.vue'),
            "bottom-footer": require('../components/footer.vue')
        }
    }
</script>