<template>
  <div id="app">
    <SiteHeader />
    <main class="main-content">
      <transition name="page" mode="out-in">
        <router-view />
      </transition>
    </main>
    <SiteFooter />
  </div>
</template>

<script>
import SiteHeader from './components/SiteHeader.vue'
import SiteFooter from './components/SiteFooter.vue'

export default {
  name: 'App',
  components: {
    SiteHeader,
    SiteFooter
  }
}
</script>

<style scoped>
#app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  overflow-x: hidden;
  max-width: 100%;
}

.main-content {
  flex: 1;
  overflow-x: hidden;
  max-width: 100%;
  display: flex;
  flex-direction: column;
}
</style>
