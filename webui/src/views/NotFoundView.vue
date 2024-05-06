<script>
import SvgIcon from "@/components/SvgIcon.vue";

export default {
    components: {SvgIcon}
}
</script>

<template>
    <div class="not-found-screen d-flex">
        <div class="card text-center">
            <div class="card-body d-flex flex-column">
                <h1 class="card-title">404</h1>
                <p class="card-text">It's dangerous to go alone! Take this.</p>
                <button class="btn btn-primary btn-sm" @click="$router.push('/')">
                    Back to safety
                    <svg-icon icon="shield"/>
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.not-found-screen {
    padding-top: 16px;
}

.card {
    min-width: max-content;
}
</style>
