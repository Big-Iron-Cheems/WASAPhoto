import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from "@/views/LoginView.vue";
import ProfileView from "@/views/ProfileView.vue";
import SearchView from "@/views/SearchView.vue";
import StreamView from "@/views/StreamView.vue";
import NotFoundView from "@/views/NotFoundView.vue";

const router = createRouter({
    history: createWebHashHistory(import.meta.env.BASE_URL),
    routes: [
        {path: '/', redirect: '/session'},
        {path: '/session', component: LoginView},
        {path: '/stream', component: StreamView},
        {path: '/users/:username/profile', component: ProfileView},
        {path: '/users', component: SearchView},
        {path: '/:pathMatch(.*)*', component: NotFoundView}
    ]
})

router.beforeEach((to, from, next) => {
    const isAuthenticated = sessionStorage.getItem("userId");
    if (!isAuthenticated && to.path !== '/session') {
        // If the user is not authenticated, redirect to the login page
        next('/session');
    } else if (isAuthenticated && to.path === '/session') {
        /**
         * If the user is already authenticated and tries to access the login page,
         * redirect to the profile page
         */
        next(`/users/${sessionStorage.getItem("username")}/profile`);
    } else {
        // Otherwise, proceed to the next route
        next();
    }
});

export default router
