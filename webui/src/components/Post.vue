<script setup>
import SvgIcon from "@/components/SvgIcon.vue";

const emit = defineEmits(['toggleComments', 'toggleLike', 'deletePost'])
const props = defineProps({
    'post': Object,
    'index': Number,
    'openCommentCardIndex': Number,
    'isCurrentUser': Boolean,
    'isStream': Boolean,
})

const toggleComments = () => {
    emit('toggleComments', props.index)
}

const toggleLike = () => {
    emit('toggleLike', props.index)
}

const deletePost = () => {
    emit('deletePost', props.index)
}
</script>

<template>
    <div class="list-group-item">
        <div class="d-flex flex-column">
            <img :src="`data:image/${post.mimeType};base64,` + post.image"
                 alt="post-thumbnail"
                 class="img-thumbnail posts-thumbnail"/>
            <span class="post-caption" v-if="post.caption">
                {{ post.caption }}
            </span>
        </div>
        <div class="btn-group-vertical">
            <p v-if="isStream">
                Author:
                <router-link :to="`/users/${post.ownerUsername}/profile`">
                    {{ post.ownerUsername }}
                    <svg-icon icon="link"/>
                </router-link>
            </p>
            <p>Likes: {{ post.likeCount }} </p>
            <p>Comments: {{ post.commentsCount }}</p>
            <button class="btn btn-sm btn-outline-primary"
                    @click="toggleComments">
                {{ openCommentCardIndex === index ? 'Hide' : 'Show' }} Comments
                <svg-icon :icon="openCommentCardIndex === index ? 'eye-off' : 'eye'"/>
            </button>
            <button class="btn btn-sm"
                    :class="post.currentUserLiked ? 'btn-outline-danger' : 'btn-outline-success'"
                    @click="toggleLike">
                {{ !post.currentUserLiked ? 'Like' : 'Unlike' }}
                <svg-icon :icon="!post.currentUserLiked ? 'thumbs-up' : 'thumbs-down'"/>
            </button>
            <button class="btn btn-sm btn-outline-danger"
                    v-if="isCurrentUser"
                    @click="deletePost">
                Delete Post
                <svg-icon icon="trash-2"/>
            </button>
        </div>
    </div>
</template>

<style scoped>
/* Posts card rules */
.posts-thumbnail {
    max-height: 200px;
    max-width: 200px;
}

.posts-thumbnail:not(:last-child) {
    border-radius: 0.375rem 0.375rem 0 0;
}

.post-caption {
    max-width: 200px;
    padding-left: 8px;
    border: 1px solid #ddd;
    border-radius: 0 0 0.375rem 0.375rem;
}

/* Item list rules */
.list-group-item {
    display: flex;
    flex-direction: row;
    align-items: center;
}

.list-group-item > *:not(:last-child) {
    margin-right: 10px;
}
</style>
