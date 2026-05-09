import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import Dashboard from './views/Dashboard.vue'
import AuthView from './views/auth/AuthView.vue'

const router = createRouter({ history: createWebHistory(), routes: [{ path: '/', component: Dashboard }, { path: '/auth', component: AuthView }] })
createApp(App).use(router).mount('#app')
