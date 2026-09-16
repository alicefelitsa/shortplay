import Vue from 'vue'
import VueRouter from 'vue-router'

Vue.use(VueRouter)

const routes = [
    {
        path: '/login',
        name: 'login',
        component: () => import('@/views/login'),
        meta: {title: '管理员登录'}
    },
    {
        path: '/',
        component: () => import('@/layout/index'),
        redirect: '/drama',
        children: [
            {
                path: 'drama',
                name: 'drama',
                component: () => import('@/views/drama/index'),
                meta: {title: '短剧管理'}
            },
            {
                path: 'episode',
                name: 'episode',
                component: () => import('@/views/episode/index'),
                meta: {title: '剧集管理'}
            },
            {
                path: 'category',
                name: 'category',
                component: () => import('@/views/category/index'),
                meta: {title: '分类管理'}
            },
            {
                path: 'setting',
                name: 'setting',
                component: () => import('@/views/setting'),
                meta: {title: '系统设置'}
            },
        ]
    },
    {
        path: '*',
        name: 'NotFound',
        component: () => import('@/views/NotFound'),
    }
]

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
})

// 全局前置守卫
router.beforeEach((to, from, next) => {
    if (to.meta.title) {
        document.title = to.meta.title
    }

    if (to.path === '/login') {
        if (localStorage.getItem("token")) {
            next({path: '/'})
        }
        next()
        return
    }
    if (localStorage.getItem("token")) {
        next()
    } else {
        next({path: '/login'})
    }
})

const originalPush = VueRouter.prototype.push
VueRouter.prototype.push = function push(location) {
    return originalPush.call(this, location).catch(err => err)
}

export default router
