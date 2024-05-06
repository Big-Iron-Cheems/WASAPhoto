<script>
import {inject} from "vue";
import ImageUploadModal from "@/components/ImageUploadModal.vue";
import ErrorMsg from "@/components/ErrorMsg.vue";
import SvgIcon from "@/components/SvgIcon.vue";
import SetTextModal from "@/components/SetTextModal.vue";
import LoadingSpinner from "@/components/LoadingSpinner.vue";
import Post from "@/components/Post.vue";
import Comment from "@/components/Comment.vue";
import Profile from "@/components/Profile.vue";

export default {
    components: {Profile, Comment, ErrorMsg, ImageUploadModal, LoadingSpinner, Post, SetTextModal, SvgIcon},
    data() {
        return {
            errorMsg: null,
            loadingStates: {
                profileCard: false,
                postsCard: false,
                commentsCard: false,
            },

            // Profile
            username: this.$route.params.username,
            userId: null,
            profile: null,
            showUsernameModal: false,
            isFollowing: false,
            hasBannedUser: false, // used by the logged-in user to manage the ban button state
            isBannedByProfileUser: false, // used by the logged-in user to see if they are allowed to see the target's profile

            // Followers, Following, Bans
            followers: [],
            following: [],
            bannedList: [],

            // Posts
            postsList: [], // Each element is post object, with an added currentUserLiked field
            showPostUploadModal: false,

            // Comments
            showCommentModal: false,
            openCommentCardIndex: null,
            shownComments: [],
        }
    },
    setup() {
        return {
            // Inject the updateUsername function from the parent component
            updateUsername: inject('updateUsername'),
        };
    },
    /**
     * Fetch the profile data when the route changes
     */
    beforeRouteUpdate(to, from, next) {
        this.username = to.params.username;
        this.fetchProfile()
        // Reset the state of shown elements
        this.openCommentCardIndex = null;
        next();
    },
    computed: {
        isCurrentUser() {
            return this.username === sessionStorage.getItem("username");
        },
    },
    methods: {
        // User

        async fetchProfile() {
            this.loadingStates.profileCard = true;
            this.errorMsg = null;

            // Reset all boolean flags
            this.isFollowing = false;
            this.isBannedByProfileUser = false;
            this.hasBannedUser = false;

            try {
                // Check these statuses only if the profile is not the current user's
                if (!this.isCurrentUser) {
                    const isBannedByUserResponse = await this.$axios.get(
                        `/users/${this.username}/bans/list/${sessionStorage.getItem("username")}`,
                        {headers: {'Authorization': `Bearer ${this.userId}`,}}
                    )
                    this.isBannedByProfileUser = isBannedByUserResponse.data.isBanned;

                    const hasBannedUserResponse = await this.$axios.get(
                        `/users/${sessionStorage.getItem("username")}/bans/list/${this.username}`,
                        {headers: {'Authorization': `Bearer ${this.userId}`,}}
                    )
                    this.hasBannedUser = hasBannedUserResponse.data.isBanned;

                    // If the session user is banned by the profile user, stop loading further unnecessary data
                    if (this.isBannedByProfileUser) return;

                    // Check the follow status
                    const followStatusResponse = await this.$axios.get(
                        `/users/${sessionStorage.getItem("username")}/followers/list/${this.username}`,
                        {headers: {'Authorization': `Bearer ${this.userId}`,}}
                    )
                    this.isFollowing = followStatusResponse.data.isFollowing;
                }

                const profileResponse = await this.$axios.get(
                    `/users/${this.username}/profile`,
                    {headers: {'Authorization': `Bearer ${this.userId}`,}}
                );
                this.profile = profileResponse.data;

                // Load the user's posts
                this.loadingStates.postsCard = true;
                const postsResponse = await this.$axios.get(
                    `/users/${this.username}/photos`,
                    {headers: {'Authorization': `Bearer ${this.userId}`,}}
                );
                const posts = postsResponse.data;

                // Update the posts with the like status, using parallel requests
                const likeStatusRequests = posts.map(post => this.$axios.get(
                    `/users/${this.username}/photos/${post.photoId}/likes/list/${sessionStorage.getItem("username")}`,
                    {headers: {'Authorization': `Bearer ${this.userId}`,}}
                ));
                const likeStatusResponses = await Promise.all(likeStatusRequests);

                // Add the like status to the posts
                this.postsList = posts.map((post, i) => {
                    post.currentUserLiked = likeStatusResponses[i].data.hasLiked;
                    return post;
                });

                this.loadingStates.postsCard = false;
            } catch (e) {
                console.error(e)
                this.errorMsg = e.response.data;
            } finally {
                this.loadingStates.profileCard = false;
            }
        },
        async setUsername({inputText}) {
            this.loadingStates.profileCard = true;
            this.errorMsg = null;
            this.showUsernameModal = false;

            try {
                const setUsernameResponse = await this.$axios.put(
                    `/users/${sessionStorage.getItem("username")}`,
                    {username: inputText},
                    {headers: {'Authorization': `Bearer ${this.userId}`, 'Content-Type': 'application/json'}}
                );

                // Update local username
                this.username = inputText;

                // Update the username in the parent component
                this.updateUsername(this.username);

                // Update the username in sessionStorage
                sessionStorage.setItem("username", inputText);

                // Redirect to the new username's profile page
                this.$router.push(`/users/${inputText}/profile`);
            } catch (e) {
                console.error(e);
                this.errorMsg = e.response.data;
            } finally {
                this.loadingStates.profileCard = false;
            }
        },
        async toggleFollowers() {
            try {
                const followersListResponse = await this.$axios.get(
                    `/users/${this.username}/followers/list`,
                    {headers: {'Authorization': `Bearer ${this.userId}`},}
                );
                this.followers = followersListResponse.data;
            } catch (e) {
                console.error(e);
                this.errorMsg = e.response.data;
            }
        },
        async toggleFollowing() {
            try {
                const followingListResponse = await this.$axios.get(
                    `/users/${this.username}/following/list`,
                    {headers: {'Authorization': `Bearer ${this.userId}`},}
                );
                this.following = followingListResponse.data;
            } catch (e) {
                console.error(e);
                this.errorMsg = e.response.data;
            }
        },
        async toggleBanned() {
            try {
                const bansListResponse = await this.$axios.get(
                    `/users/${this.username}/bans/list`,
                    {headers: {'Authorization': `Bearer ${this.userId}`},}
                );
                this.bannedList = bansListResponse.data;
            } catch (e) {
                console.error(e);
                this.errorMsg = e.response.data;
            }
        },

        // Posts

        async uploadPost({image, mimeType, caption}) {
            this.loadingStates.postsCard = true;
            this.errorMsg = null;
            this.showPostUploadModal = false;

            // Create form data
            let formData = new FormData();
            formData.append("image", image)
            formData.append("mimeType", mimeType)
            formData.append("caption", caption)

            try {
                const uploadPostResponse = await this.$axios.post(
                    `users/${sessionStorage.getItem("username")}/photos`,
                    formData,
                    {headers: {'Authorization': `Bearer ${this.userId}`, 'Content-Type': 'multipart/form-data'}}
                );

                // Update the state and data, insert the newer post at the beginning
                this.profile.photoCount++;
                this.postsList.unshift(uploadPostResponse.data);
                // Adjust the openCommentCardIndex if necessary
                if (this.openCommentCardIndex !== null && this.openCommentCardIndex >= 0) this.openCommentCardIndex++;
            } catch (e) {
                console.error(e);
                this.errorMsg = e.response.data;
            } finally {
                this.loadingStates.postsCard = false;
            }
        },
        async deletePost(index) {
            this.loadingStates.postsCard = true;
            this.errorMsg = null;

            const post = this.postsList[index];
            try {
                const deletePostResponse = await this.$axios.delete(
                    `/users/${sessionStorage.getItem("username")}/photos/${post.photoId}`,
                    {headers: {'Authorization': `Bearer ${this.userId}`}}
                );

                // Update the state and data
                this.profile.photoCount--;
                this.postsList.splice(index, 1);
                // Adjust the openCommentCardIndex if necessary
                this.openCommentCardIndex = null
            } catch (e) {
                console.error(e);
                this.errorMsg = e.response.data;
            } finally {
                this.loadingStates.postsCard = false;
            }
        },

        // Follow/Ban

        async toggleFollow() {
            try {
                if (!this.isFollowing) {
                    // Follow the user
                    const followResponse = await this.$axios.post(
                        `/users/${this.username}/followers`,
                        {username: this.username},
                        {headers: {'Authorization': `Bearer ${this.userId}`, 'Content-Type': 'application/json'}}
                    );
                } else {
                    // Unfollow the user
                    const unfollowResponse = await this.$axios.delete(
                        `/users/${sessionStorage.getItem("username")}/followers/${this.username}`,
                        {headers: {'Authorization': `Bearer ${this.userId}`,}}
                    );
                }

                // Update the state and data
                this.isFollowing = !this.isFollowing;
                const profileResponse = await this.$axios.get(
                    `/users/${this.username}/profile`,
                    {headers: {'Authorization': `Bearer ${this.userId}`,}}
                );
                this.profile.followersCount = profileResponse.data.followersCount;
            } catch (e) {
                console.error(e);
                this.errorMsg = e.response.data;
            }
        },
        async toggleBan() {
            try {
                if (!this.hasBannedUser) {
                    // Ban the user
                    const banResponse = await this.$axios.post(
                        `/users/${this.username}/bans`,
                        {username: this.username},
                        {headers: {'Authorization': `Bearer ${this.userId}`, 'Content-Type': 'application/json'}}
                    );
                } else {
                    // Unban the user
                    const unbanResponse = await this.$axios.delete(
                        `/users/${sessionStorage.getItem("username")}/bans/${this.username}`,
                        {headers: {'Authorization': `Bearer ${this.userId}`,}}
                    );
                }

                // Update the state and data
                // No need to fetch the profile again, the bannedCount is not displayed on other users' profiles
                this.hasBannedUser = !this.hasBannedUser;
            } catch (e) {
                console.error(e);
                this.errorMsg = e.response.data;
            }
        },

        // Likes

        async toggleLike(index) {
            const post = this.postsList[index];
            try {
                if (!post.currentUserLiked) {
                    // Like the post
                    const likeResponse = await this.$axios.post(
                        `/users/${this.username}/photos/${post.photoId}/likes`,
                        {userId: this.userId},
                        {headers: {'Authorization': `Bearer ${this.userId}`, 'Content-Type': 'application/json'}}
                    );

                    // Update the state and data
                    post.likeCount++;
                    post.currentUserLiked = true;
                } else {
                    // Unlike the post
                    const unlikeResponse = await this.$axios.delete(
                        `/users/${this.username}/photos/${post.photoId}/likes/${this.userId}`,
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
                    `/users/${this.username}/photos/${post.photoId}/comments`,
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
                    `/users/${this.username}/photos/${post.photoId}/comments`,
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
                    `/users/${this.username}/photos/${post.photoId}/comments/${comment.commentId}`,
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
        this.fetchProfile()
    }
}
</script>

<template>
    <div class="profile-screen">
        <h1 class="border-bottom">Profile</h1>
        <error-msg v-if="errorMsg" :msg="errorMsg"/>
        <div v-else class="d-flex">
            <div class="card profile-card">
                <div class="card-header d-flex align-items-center">
                    <h2>{{ this.username }}</h2>
                    <button class="btn btn-sm btn-outline-primary" v-if="isCurrentUser"
                            @click="showUsernameModal = true">
                        Edit username
                        <svg-icon icon="edit-2"/>
                    </button>
                </div>
                <loading-spinner :loading="loadingStates.profileCard">
                    <div class="card-body d-flex flex-column">
                        <profile
                            :profile="profile"
                            :isCurrentUser="isCurrentUser"
                            :followers="followers"
                            :following="following"
                            :bannedList="bannedList"
                            :isFollowing="isFollowing"
                            :hasBannedUser="hasBannedUser"
                            :isBannedByProfileUser="isBannedByProfileUser"
                            @toggleFollow="toggleFollow"
                            @toggleBan="toggleBan"
                            @toggleFollowers="toggleFollowers"
                            @toggleFollowing="toggleFollowing"
                            @toggleBanned="toggleBanned"
                        />
                    </div>
                </loading-spinner>
            </div>
            <div v-if="!isBannedByProfileUser" class="card posts-card">
                <div class="card-header d-flex align-items-center">
                    <h2 class="user-posts-title">{{ this.username }}'s posts</h2>
                    <button class="btn btn-sm btn-outline-primary"
                            v-if="isCurrentUser"
                            @click="showPostUploadModal = true">
                        Create new post
                        <svg-icon icon="file-plus"/>
                    </button>
                </div>
                <loading-spinner :loading="loadingStates.postsCard">
                    <div class="card-body">
                        <div class="d-flex flex-column list-group posts-container">
                            <div v-if="postsList.length === 0" class="list-group-item"> No posts to show yet</div>
                            <post v-else v-for="(post, index) in this.postsList" :key="index"
                                  :post="post"
                                  :index="index"
                                  :openCommentCardIndex="openCommentCardIndex"
                                  :isCurrentUser="isCurrentUser"
                                  :isStream="false"
                                  @toggleComments="toggleComments"
                                  @toggleLike="toggleLike"
                                  @deletePost="deletePost"
                            />
                        </div>
                    </div>
                </loading-spinner>
            </div>
            <div v-if="!isBannedByProfileUser && openCommentCardIndex !== null" class="card comment-card">
                <div class="card-header d-flex align-items-center">
                    <h2 class="user-comments-title">Comments</h2>
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
                                     :isCurrentUser="isCurrentUser"
                                     :isCommentOwner="isCommentOwner"
                                     @deleteComment="deleteComment"
                            />
                        </div>
                    </div>
                </loading-spinner>
            </div>
        </div>
    </div>

    <set-text-modal v-if="showUsernameModal"
                    header="Change your username"
                    placeholder="New username"
                    pattern="^[A-Za-z0-9_\-]{3,32}$"
                    :minlength=3
                    :maxlength=32
                    title="3 to 32 alphanumeric characters, allowing _ and -"
                    :rows=1
                    :allowEnter=false
                    @confirm="setUsername"
                    @cancel="this.showUsernameModal = false"
    />
    <image-upload-modal v-if="showPostUploadModal"
                        :question="'Upload new post'"
                        @confirm="uploadPost"
                        @cancel="this.showPostUploadModal = false"
    />
    <set-text-modal v-if="showCommentModal"
                    header="Add a comment"
                    placeholder="Comment text"
                    pattern="^[\p{L}\p{N}\p{M}\p{P}\p{S} \n]{1,256}$"
                    :minlength=1
                    :maxlength=256
                    title="1 to 256 characters (UNICODE supported)"
                    :rows=3
                    :allowEnter=true
                    @confirm="uploadComment($event, this.openCommentCardIndex)"
                    @cancel="this.showCommentModal = false"
    />
</template>

<style scoped>
.profile-screen {
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

/* Profile card rules */
.profile-card {
    min-width: 275px;
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
