import { createApp } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import App from './App.vue'
import Dashboard from './views/Dashboard.vue'
import AuthView from './views/auth/AuthView.vue'

const router = createRouter({ history: createWebHashHistory(), routes: [{ path: '/', component: Dashboard }, { path: '/auth', component: AuthView }] })
createApp(App).use(router).mount('#app')
