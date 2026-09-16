import Vue from 'vue'
import App from './App.vue'
import router from './router'
import './assets/styles.css'

Vue.config.productionTip = false

// 图片懒加载指令
Vue.directive('lazy', {
  inserted(el) {
    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          const img = entry.target
          if (img.dataset.src) {
            img.src = img.dataset.src
            img.removeAttribute('data-src')
          }
          observer.unobserve(img)
        }
      })
    })
    observer.observe(el)
  }
})

new Vue({
  router,
  render: h => h(App),
}).$mount('#app')
