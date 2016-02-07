<template>
	<div class="panel panel-default monitor">
		<div class="panel-heading">
			{{title}}
		</div>
		<div class="panel-body">
			<div :id="chart_id" >{{success}}/{{total}}</div>	
			<div class="detail-total">{{total}} <small>Total tasks</small></div>	
			<div class="detail-success detail-data">{{success}}<small>Successful tasks</small></div>
			<div class="detail-fail detail-data">{{total-success}}<small>Failed tasks</small></div>
		</div>
	</div>
</template>

<script>
export default {
	props: {
		title: String,
		color: String,
		success: Number,
		total: Number,
		refresh: Number
	},
	data() {
		let chart_id = 'chart_' + Math.floor(Math.random() * 1000).toString(16)
		return {
			chart_id: chart_id,
			chart: null
		}
	},
	ready() {
    	this.chart = $('#' + this.chart_id).peity('donut', {
      		width: '100%',
      		height: 180,
      		innerRadius: 75, 
      		fill: ["green", "red"]
	    })	
	}
}	
</script>

<style>
.panel-body{
	position: relative;
}
.detail-total{
	position: absolute;
	top: 60px;
	text-align: center;
	width: 100%;
	margin-left: -15px;
	font-size: 45px;
}
.detail-total small{
	display: block;
	font-size: 13px;
}
.detail-data{
	font-size: 36px;
}
.detail-data small{
	font-size: 13px;
	line-height: 13px;
	display: block;
}
.detail-success{
	float: left;
	margin-left: 10px;
}
.detail-fail{
	float: right;
	margin-right: 10px;
}
.detail-success small{
	color: green;
}
.detail-fail small {
	color: red;
}
</style>