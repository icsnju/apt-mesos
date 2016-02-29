<template>
  <div class="page page-current">
      <tab> </tab>
      <router-view transition="fade" transition-mode="out-in" keep-alive></router-view>
  </div>
</template>

<script>
import Tab from './components/tab.vue'

export default {
  ready () {
    this.options.items = [
      {
        path: '/home',
        icon: 'app',
        label: '首页'
      },
      {
        path: '/article',
        icon: 'star',
        label: '任务'
      },
      {
        path: '/user',
        icon: 'me',
        label: '我自己'
      },
      {
        path: '/more',
        icon: 'gift',
        label: '更多'
      }
    ]
  },  
  data () {
    return {
      isIndex: true,
      title: '首页',
      options: {
        items: [],
        activeClass: 'active'
      }
    }
  },
  components: {
    Tab,
  }
}
</script>

<style>
[v-cloak] {
  display: none;
}
.fade-transition {
  transition: opacity .3s ease;
}
.fade-enter, .fade-leave {
  opacity: 0;
}
</style>
