<script>
import {inject} from "vue";
import ErrorMsg from "@/components/ErrorMsg.vue";
import SvgIcon from "@/components/SvgIcon.vue";

export default {
    components: {SvgIcon, ErrorMsg},
    data() {
        return {
            errorMsg: null,

            username: '',
            userId: 0,
        };
    },
    setup() {
        return {
            // Inject the updateUsername function from the parent component
            updateUsername: inject('updateUsername'),
        };
    },
    methods: {
        async login() {
            try {
                const loginResponse = await this.$axios.post('/session', {username: this.username});
                this.userId = loginResponse.data.userId;
                this.username = loginResponse.data.username;

                // Update the username in the parent component
                this.updateUsername(this.username);

                // Set the username and userId in sessionStorage
                sessionStorage.setItem("username", this.username);
                sessionStorage.setItem("userId", this.userId);

                // Redirect to the user's profile page
                this.$router.push(`/users/${this.username}/profile`);
            } catch (e) {
                console.error(e);
                this.errorMsg = e.response.data;
            }
        },
    },
};
</script>

<template>
    <div class="login-screen">
        <h1 class="border-bottom">Login</h1>
        <error-msg v-if="errorMsg" :msg="errorMsg"/>
        <div class="input-container">
            <form @submit.prevent="login">
                <div class="input-group">
                    <input type="text"
                           class="form-control"
                           v-model.trim="username"
                           placeholder="Username"
                           required
                           pattern="^[A-Za-z0-9_\-]{3,32}$"
                           minlength="3"
                           maxlength="32"
                           title="3 to 32 alphanumeric characters, allowing _ and -">
                    <div class="input-group-append">
                        <button type="submit" class="btn btn-primary"
                                :disabled="username.length < 3 || username.length > 32">
                            Login
                            <svg-icon icon="log-in"/>
                        </button>
                    </div>
                </div>
            </form>
        </div>
    </div>
</template>

<style scoped>

.login-screen {
    padding-top: 16px;
}

.input-container {
    display: flex;
    min-width: max-content;
}

.input-group input {
    width: 32ch;
}
</style>
