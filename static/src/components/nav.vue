<template>
	<ul class="nav">
		<li v-for="item in items" v-link="{path:item.path, activeClass: 'active'}">
			<a v-link="{path:item.path}"><span class="si-{{item.icon}}"></span><span class="nav-label">{{item.label}}</span></a>
		</li>
	</ul>	
</template>

<script>
	export default {
		props: {
			items: Array
		},
	}
</script>