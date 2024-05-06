<script setup>
import SvgIcon from "@/components/SvgIcon.vue";

const emit = defineEmits(['deleteComment'])
const props = defineProps({
    'comment': Object,
    'index': Number,
    'isCurrentUser': Boolean,
    'isCommentOwner': Function,
})

const deleteComment = () => {
    emit('deleteComment', props.index)
}
</script>

<template>
    <div class="list-group-item">
        <div class="d-flex flex-column border comment">
            <router-link class="comment-owner-link"
                         :to="`/users/${comment.ownerUsername}/profile`">
                <h5>
                    {{ comment.ownerUsername }}
                    <svg-icon icon="link"/>
                </h5>
            </router-link>
            <hr>
            <p class="text-wrap">{{ comment.content }}</p>
        </div>
        <button v-if="isCurrentUser || isCommentOwner(index)"
                class="btn btn-sm btn-outline-danger"
                @click="deleteComment">
            Delete Comment
            <svg-icon icon="trash-2"/>
        </button>
    </div>
</template>

<style scoped>
/* Comment card rules */
.comment {
    border-radius: 0.375rem;
    flex-grow: 1;
}

.comment > * {
    margin-left: 10px;
    margin-right: 10px;
}

.comment > hr {
    margin-top: 0;
    margin-bottom: 10px;
}

.comment-owner-link {
    max-width: fit-content;
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
