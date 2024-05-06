<script>
import SvgIcon from "@/components/SvgIcon.vue";
import ErrorMsg from "@/components/ErrorMsg.vue";
import SetTextModal from "@/components/SetTextModal.vue";
import LoadingSpinner from "@/components/LoadingSpinner.vue";
import Post from "@/components/Post.vue";
import Comment from "@/components/Comment.vue";

export default {
    components: {Comment, ErrorMsg, LoadingSpinner, Post, SetTextModal, SvgIcon},
    data() {
        return {
            errorMsg: null,
            loadingStates: {
                postsCard: false,
                commentsCard: false,
            },

            // Profile
            username: this.$route.params.username,
            userId: null,

            // Stream
            postsList: [], // These posts are fetched from the users you follow, shown from newest to oldest

            // Comments
            showCommentModal: false,
            openCommentCardIndex: null,
            shownComments: [],
        }
    },
    methods: {
        // User

        async fetchStream() {
            this.loadingStates.postsCard = true;
            this.errorMsg = null;

            try {
                // Load the user's stream
                const streamResponse = await this.$axios.get(
                    `/stream`,
                    {headers: {'Authorization': `Bearer ${this.userId}`,}}
                );
                const posts = streamResponse.data;

                // Update the posts with the like status, using parallel requests
                const likeStatusRequests = posts.map(post => this.$axios.get(
                    `/users/${post.ownerUsername}/photos/${post.photoId}/likes/list/${sessionStorage.getItem("username")}`,
                    {headers: {'Authorization': `Bearer ${this.userId}`,}}
                ));
                const likeStatusResponses = await Promise.all(likeStatusRequests);

                // Add the like status to the posts
                this.postsList = posts.map((post, i) => {
                    post.currentUserLiked = likeStatusResponses[i].data.hasLiked;
                    return post;
                });
            } catch (e) {
                console.error(e);
                this.errorMsg = e.response.data;
            } finally {
                this.loadingStates.postsCard = false;
            }
        },

        // Likes

        async toggleLike(index) {
            const post = this.postsList[index];
            try {
                if (!post.currentUserLiked) {
                    // Like the post
                    const likeResponse = await this.$axios.post(
                        `/users/${post.ownerUsername}/photos/${post.photoId}/likes`,
                        {userId: this.userId},
                        {headers: {'Authorization': `Bearer ${this.userId}`, 'Content-Type': 'application/json'}}
                    );

                    // Update the state and data
                    post.likeCount++;
                    post.currentUserLiked = true;
                } else {
                    // Unlike the post
                    const unlikeResponse = await this.$axios.delete(
                        `/users/${post.ownerUsername}/photos/${post.photoId}/likes/${this.userId}`,
                        {headers: {'Authorization': `Bearer ${this.userId}`,}}
                    );

                    // Update the state and data
                    post.likeCount--;
                    post.currentUserLiked = false;
                }
            } catch (e) {
                console.error(e);
                this.errorMsg = e.response.data;
            }
        },

        // Comments

        async fetchComments(index) {
            this.loadingStates.commentsCard = true;
            this.errorMsg = null;

            const post = this.postsList[index];
            try {
                const commentsResponse = await this.$axios.get(
                    `/users/${post.ownerUsername}/photos/${post.photoId}/comments`,
                    {headers: {'Authorization': `Bearer ${this.userId}`,}}
                );
                this.shownComments = commentsResponse.data;
            } catch (e) {
                console.error(e);
                this.errorMsg = e.response.data;
            } finally {
                this.loadingStates.commentsCard = false;
            }
        },
        isCommentOwner(index) {
            return this.shownComments[index].ownerUsername === sessionStorage.getItem("username");
        },
        async uploadComment({inputText}, index) {
            this.loadingStates.commentsCard = true;
            this.errorMsg = null;
            this.showCommentModal = false;

            const post = this.postsList[index];
            try {
                const uploadCommentResponse = await this.$axios.post(
                    `/users/${post.ownerUsername}/photos/${post.photoId}/comments`,
                    {content: inputText},
                    {headers: {'Authorization': `Bearer ${this.userId}`, 'Content-Type': 'application/json'}}
                );

                // Update the state and data
                this.shownComments.push({
                    commentId: uploadCommentResponse.data.commentId,
                    ownerId: this.userId,
                    ownerUsername: sessionStorage.getItem("username"),
                    content: inputText
                });
                post.commentsCount++;
            } catch (e) {
                console.error(e);
                this.errorMsg = e.response.data;
            } finally {
                this.loadingStates.commentsCard = false;
            }
        },
        async deleteComment(index) {
            this.loadingStates.commentsCard = true;
            this.errorMsg = null;

            const post = this.postsList[this.openCommentCardIndex];
            const comment = this.shownComments[index];
            try {
                const deleteCommentResponse = await this.$axios.delete(
                    `/users/${post.ownerUsername}/photos/${post.photoId}/comments/${comment.commentId}`,
                    {headers: {'Authorization': `Bearer ${this.userId}`,}}
                );

                // Update the state and data
                this.shownComments.splice(index, 1);
                this.postsList[this.openCommentCardIndex].commentsCount--;
            } catch (e) {
                console.error(e);
                this.errorMsg = e.response.data;
            } finally {
                this.loadingStates.commentsCard = false;
            }
        },
        async toggleComments(index) {
            if (this.openCommentCardIndex === index) {
                this.openCommentCardIndex = null;
            } else {
                this.openCommentCardIndex = index;
                await this.fetchComments(index);
            }
        },
    },
    mounted() {
        this.userId = sessionStorage.getItem("userId");
        this.fetchStream();
    }
}
</script>

<template>
    <div class="stream-screen">
        <h1 class="border-bottom">Stream</h1>
        <error-msg v-if="errorMsg" :msg="errorMsg"/>
        <div v-else class="d-flex">
            <div class="card posts-card">
                <div class="card-header d-flex align-items-center">
                    <h2>Posts from users you follow</h2>
                </div>
                <loading-spinner :loading="loadingStates.postsCard">
                    <div class="card-body">
                        <div class="d-flex flex-column list-group posts-container">
                            <div v-if="postsList.length === 0" class="list-group-item">
                                No posts to show yet. Follow someone to see their posts here.
                            </div>
                            <post v-else v-for="(post, index) in this.postsList" :key="index"
                                  :post="post"
                                  :index="index"
                                  :openCommentCardIndex="openCommentCardIndex"
                                  :isCurrentUser="false"
                                  :isStream="true"
                                  @toggleComments="toggleComments"
                                  @toggleLike="toggleLike"
                            />
                        </div>
                    </div>
                </loading-spinner>
            </div>
            <div v-if="openCommentCardIndex !== null" class="card comment-card">
                <div class="card-header d-flex align-items-center">
                    <h2 class="user-comments-title">Post comments</h2>
                    <button class="btn btn-sm btn-outline-primary"
                            @click="showCommentModal = true">
                        Add new comment
                        <svg-icon icon="message-square"/>
                    </button>
                </div>
                <loading-spinner :loading="loadingStates.commentsCard">
                    <div class="card-body">
                        <div class="d-flex flex-column list-group comments-container">
                            <div v-if="this.shownComments.length === 0" class="list-group-item">
                                No comments to show
                            </div>
                            <comment v-else v-for="(comment, index) in this.shownComments" :key="index"
                                     :comment="comment"
                                     :index="index"
                                     :isCurrentUser="false"
                                     :isCommentOwner="isCommentOwner"
                                     @deleteComment="deleteComment"
                            />
                        </div>
                    </div>
                </loading-spinner>
            </div>
        </div>
    </div>

    <set-text-modal v-if="showCommentModal"
                    header="Add a comment"
                    placeholder="Comment text"
                    pattern="^[\p{L}\p{N}\p{M}\p{P}\p{S} \n]{1,256}$"
                    :minlength=1
                    :maxlength=256
                    title="1 to 256 characters (UNICODE supported)"
                    :rows=3
                    :allowEnter=true
                    @confirm="uploadComment($event, openCommentCardIndex)"
                    @cancel="this.showCommentModal = false"
    />
</template>

<style scoped>
.stream-screen {
    padding-top: 16px;
    padding-bottom: 16px;
}

.card {
    width: max-content;
    min-width: max-content;
    height: min-content;
    border: 1px solid #ddd;
    border-radius: 4px;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.15);
}

.card:not(:last-child) {
    margin-right: 20px;
}

.card-header {
    padding: 10px;
    background-color: #f5f5f5;
    border-bottom: 1px solid #ddd;
    justify-content: space-between;
}

.card-header button {
    margin-left: 10px;
}

.card-body {
    padding: 10px;
}

/* Posts card rules */
.posts-container {
    max-height: calc(100vh - 250px);
    overflow-y: auto;
}

/* Comment card rules */
.comments-container {
    max-height: calc(100vh - 250px);
    overflow-y: auto;
}
</style>
