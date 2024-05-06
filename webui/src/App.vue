<script setup>
import {RouterLink, RouterView} from 'vue-router'
import {provide, ref} from "vue";
import SvgIcon from "@/components/SvgIcon.vue";

const username = ref(sessionStorage.getItem('username'))

const logout = () => {
    sessionStorage.clear()
    username.value = null
}

const updateUsername = (newUsername) => {
    username.value = newUsername
}

provide('updateUsername', updateUsername)
</script>

<template>
    <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
        <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6">WASAPhoto</a>
        <button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse"
                data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false"
                aria-label="Toggle navigation">
            <span class="navbar-toggler-icon"></span>
        </button>
    </header>

    <div class="container-fluid">
        <div class="row">
            <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
                <div class="position-sticky pt-3 sidebar-sticky">
                    <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
                        <span>Navigation</span>
                    </h6>
                    <ul class="nav flex-column">
                        <li class="nav-item" v-if="!username">
                            <RouterLink to="/session" class="nav-link">
                                <svg-icon icon="key"/>
                                Login
                            </RouterLink>
                        </li>
                        <li class="nav-item" v-if="username">
                            <RouterLink to="/stream" class="nav-link">
                                <svg-icon icon="list"/>
                                My stream
                            </RouterLink>
                        </li>
                        <li class="nav-item" v-if="username">
                            <RouterLink :to="`/users/${username}/profile`" class="nav-link">
                                <svg-icon icon="home"/>
                                My profile
                            </RouterLink>
                        </li>
                        <li class="nav-item" v-if="username">
                            <RouterLink to="/users" class="nav-link">
                                <svg-icon icon="search"/>
                                Search users
                            </RouterLink>
                        </li>
                    </ul>

                    <template v-if="username">
                        <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
                            <span>Settings</span>
                        </h6>
                        <ul class="nav flex-column">
                            <li class="nav-item">
                                <RouterLink to="/session" class="nav-link" @click="logout">
                                    <svg-icon icon="log-out"/>
                                    Logout
                                </RouterLink>
                            </li>
                        </ul>
                    </template>
                </div>
            </nav>

            <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
                <RouterView/>
            </main>
        </div>
    </div>
</template>

<style>
.nav-link {
    max-width: fit-content;
}
</style>
