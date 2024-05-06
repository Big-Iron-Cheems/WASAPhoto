<script>
import ErrorMsg from "@/components/ErrorMsg.vue";
import SvgIcon from "@/components/SvgIcon.vue";
import LoadingSpinner from "@/components/LoadingSpinner.vue";

export default {
    components: {LoadingSpinner, SvgIcon, ErrorMsg},
    data() {
        return {
            errorMsg: null,
            loading: false,

            username: '',
            allUsers: [],
            currentPage: 1,
            pageSize: 50,
        };
    },
    computed: {
        filteredUsers() {
            return this.allUsers
                ? this.allUsers.filter(user => user.username.toLowerCase().includes(this.username.toLowerCase()))
                : [];
        },
    },
    methods: {
        async fetchUsers(page = this.currentPage, pageSize = this.pageSize) {
            this.loading = true;
            try {
                const response = await this.$axios.get(`/users?page=${page}&pageSize=${pageSize}`);
                this.allUsers = response.data;
                this.currentPage = page;
            } catch (e) {
                console.error(e);
                this.errorMsg = e.response.data;
            } finally {
                this.loading = false;
            }
        },
    },
    mounted() {
        this.fetchUsers();
    },
};
</script>

<template>
    <div class="search-screen">
        <h1 class="border-bottom">Search</h1>
        <error-msg v-if="errorMsg" :msg="errorMsg"/>
        <loading-spinner v-if="this.loading" :loading="this.loading"/>
        <div v-else class="content-container">
            <div class="input-group">
                <input type="text"
                       class="form-control"
                       id="username"
                       placeholder="Username"
                       v-model.trim="username"
                       pattern="^[A-Za-z0-9_\-]{3,32}$"
                       minlength="3"
                       maxlength="32"
                       title="3 to 32 alphanumeric characters, allowing _ and -">
            </div>
            <div class="card-body">
                <button class="btn btn-sm" @click="fetchUsers(currentPage - 1)" :disabled="currentPage === 1">Previous
                </button>
                <span>Page {{ currentPage }}</span>
                <button class="btn btn-sm" @click="fetchUsers(currentPage + 1)"
                        :disabled="allUsers ? allUsers.length < pageSize : true">Next
                </button>
            </div>
            <div>
                <table class="table table-bordered caption-top">
                    <caption>List of users</caption>
                    <thead>
                    <tr>
                        <th scope="col">#</th>
                        <th scope="col">Username</th>
                    </tr>
                    </thead>
                    <tbody>
                    <tr v-for="(user, index) in filteredUsers" :key="user.id">
                        <th scope="row">{{ index + 1 }}</th>
                        <td>
                            <router-link :to="`/users/${user.username}/profile`">
                                {{ user.username }}
                                <svg-icon icon="link"/>
                            </router-link>
                        </td>
                    </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</template>


<style scoped>
.search-screen {
    padding-top: 16px;
}

.content-container {
    display: inline-block;
    width: auto;
}

.card-body span {
    margin: 10px;
}

.input-group {
    padding-bottom: 8px;
}

.input-group input {
    width: 32ch;
}
</style>
