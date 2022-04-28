import Vue from 'vue'
import VueRouter from 'vue-router'
import Index from '../components/Index.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    redirect: 'index'
  },
  {
    path: '/index',
    name: 'index',
    component: Index
  },

]

const router = new VueRouter({
  routes
})

export default router
