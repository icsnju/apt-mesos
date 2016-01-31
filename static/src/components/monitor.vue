<template>
	<div class="panel panel-default monitor">
		<div class="panel-heading">
			{{title}}
		</div>
		<div class="panel-body">
			<div class="percent"> {{percent| p0}} <small>%</small></div> 
			<div class="detail">{{used | p2}} {{unit}} of {{total | p2}} {{unit}} used</div>
			<div :id="chart_id" >{{data}}</div>
		</div>
	</div>
</template>

<script>
export default {
	props: {
		title: String,
		color: String,
		unit: String,
		total: Number,
		used: Number,
		refresh: Number
	},

	data() {
		let chart_id = 'chart_' + Math.floor(Math.random() * 1000).toString(16),
			data = [],
			totalPoints = 50
		
		while (data.length < totalPoints) {
			data.push(0)
		}	
		return {
			chart_id: chart_id,
			data: data,
			percent: 0,
			chart: null
		}
	},
	ready() {
    	this.chart = $('#' + this.chart_id).peity('line', {
      		width: '100%',
      		height: 100,
	    })		
	},
  	watch: {
    	refresh: function () {
    		this.percent = this.total==0 ? 0 : this.used/this.total*100
	      	this.data.push(this.percent)
	      	this.data.shift()
	      	this.chart.change()
    	}	
  	},	
  	filters: {
    	p2(val) {
    	  	return Math.floor(val*100)/100
    	},
    	p0(val) {
    		return Math.floor(val)
    	}
  	}	
}	
</script>

<style>
.panel.monitor {
	text-align: center;
}
.monitor .percent {
	font-size: 36px;
}
.monitor .percent small {
	font-size: 24px;
}
svg.peity{
	padding: 10px;
}
</style>