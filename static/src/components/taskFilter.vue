<template>
	<h5>Status</h5>
	<ul class="category-list" style="padding: 0">
		<li v-for="item in options.status">
			<a  v-link="{path:'task',  query: { status: item.id }}"><i class="fa fa-circle text-{{item.color}}"></i>{{item.label}}</a>
		</li>

	</ul>	
</template>
<script>
export default {
	data() {
		let options = {}
		options.status = [
			{
				label: "Running", 
				id: 1,
				color: "info"
			}, {
				label: "Staging",
				id: 2,
				color: "primary"
			}, {
				label: "Killed",
				id: 3,
				color: "warning"
			}, {
				label: "Failed",
				id: 4,
				color: "danger"
			}, {
				label: "Finished",
				id: 5,
				color: "navy"
			}
		]
		return {
			options: options
		}
	}
}
</script>