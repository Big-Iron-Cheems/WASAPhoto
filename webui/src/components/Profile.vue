<script setup>
import SvgIcon from "@/components/SvgIcon.vue";
import {ref} from "vue";

const emit = defineEmits([
    'toggleFollow',
    'toggleBan',
    'toggleFollowers',
    'toggleFollowing',
    'toggleBanned'
])
const props = defineProps({
    'profile': Object,
    'isCurrentUser': Boolean,
    'followers': Array,
    'following': Array,
    'bannedList': Array,
    'isFollowing': Boolean,
    'hasBannedUser': Boolean,
    'isBannedByProfileUser': Boolean,
    'showFollowing': Boolean,
    'showBanned': Boolean,
})

let showFollowers = ref(false);
let showFollowing = ref(false);
let showBanned = ref(false);

const toggleFollow = () => {
    emit('toggleFollow')
}

const toggleBan = () => {
    emit('toggleBan')
}

const toggleFollowers = () => {
    showFollowers.value = !showFollowers.value
    emit('toggleFollowers')
}

const toggleFollowing = () => {
    showFollowing.value = !showFollowing.value
    emit('toggleFollowing')
}

const toggleBanned = () => {
    showBanned.value = !showBanned.value
    emit('toggleBanned')
}
</script>

<template>
    <div class="profile-info border-bottom" v-if="!isBannedByProfileUser">
        <p>Photos: {{ profile?.photoCount }}</p>
        <p>Followers: {{ profile?.followersCount }}</p>
        <p>Following: {{ profile?.followingCount }}</p>
        <p v-if="isCurrentUser">Banned: {{ profile?.bannedCount }}</p>
    </div>
    <div class="profile-info border-bottom" v-else>
        <p>
            You are banned by this user.
            Profile info, posts and comments are hidden until the ban is lifted.
        </p>
    </div>

    <template v-if="!isCurrentUser">
        <div class="btn-group-vertical d-flex">
            <button class="btn btn-sm" @click="toggleFollow"
                    :class="!isFollowing ? 'btn-outline-success' : 'btn-outline-danger'"
                    v-if="!isBannedByProfileUser">
                {{ isFollowing ? 'Unfollow' : 'Follow' }}
                <svg-icon :icon="!isFollowing ? 'user-plus' : 'user-minus'"/>
            </button>
            <button class="btn btn-sm" @click="toggleBan"
                    :class="!hasBannedUser? 'btn-outline-danger' : 'btn-outline-success'">
                {{ hasBannedUser ? 'Unban' : 'Ban' }}
                <svg-icon :icon="!hasBannedUser ?'slash':'check-circle'"/>
            </button>
        </div>
    </template>
    <template v-else>
        <div class="btn-group-vertical d-flex">
            <button class="btn btn-sm btn-outline-primary"
                    @click="toggleFollowers"
                    :disabled="profile?.followersCount === 0">
                {{ showFollowers ? 'Hide' : 'Show' }} Followers
                <svg-icon :icon="showFollowers ? 'eye-off' : 'eye'"/>
            </button>
            <button class="btn btn-sm btn-outline-primary"
                    @click="toggleFollowing"
                    :disabled="profile?.followingCount === 0">
                {{ showFollowing ? 'Hide' : 'Show' }} Following
                <svg-icon :icon="showFollowing ? 'eye-off' : 'eye'"/>
            </button>
            <button class="btn btn-sm btn-outline-primary"
                    @click="toggleBanned"
                    :disabled="profile?.bannedCount === 0">
                {{ showBanned ? 'Hide' : 'Show' }} Banned
                <svg-icon :icon="showBanned ? 'eye-off' : 'eye'"/>
            </button>
        </div>
    </template>

    <table class="table table-bordered caption-top" v-if="showFollowers">
        <caption>List of followers</caption>
        <thead>
        <tr>
            <th scope="col">Username</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="follower in followers" :key="follower.id">
            <td>
                <router-link :to="`/users/${follower.username}/profile`">
                    {{ follower.username }}
                    <svg-icon icon="link"/>
                </router-link>
            </td>
        </tr>
        </tbody>
    </table>

    <table class="table table-bordered caption-top" v-if="showFollowing">
        <caption>List of followed users</caption>
        <thead>
        <tr>
            <th scope="col">Username</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="followed in following" :key="followed.id">
            <td>
                <router-link :to="`/users/${followed.username}/profile`">
                    {{ followed.username }}
                    <svg-icon icon="link"/>
                </router-link>
            </td>
        </tr>
        </tbody>
    </table>

    <table class="table table-bordered caption-top" v-if="showBanned">
        <caption>List of banned users</caption>
        <thead>
        <tr>
            <th scope="col">Username</th>
        </tr>
        </thead>
        <tbody>
        <tr v-for="banned in bannedList" :key="banned.id">
            <td>
                <router-link :to="`/users/${banned.username}/profile`">
                    {{ banned.username }}
                    <svg-icon icon="link"/>
                </router-link>
            </td>
        </tr>
        </tbody>
    </table>
</template>

<style scoped>
.profile-info {
    margin-bottom: 20px;
}

.profile-info p {
    width: 32ch;
}
</style>
