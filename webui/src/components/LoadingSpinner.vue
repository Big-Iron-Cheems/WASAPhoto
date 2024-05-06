<script setup>
const props = defineProps(["loading"])
</script>

<template>
    <div v-if="loading">
        <div class="spinner-container">
            <div class="spinner-border">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>
    </div>
    <div v-if="!loading">
        <slot/>
    </div>
</template>

<style>
.spinner-container {
    text-align: center;
    margin-top: 5px;
    margin-bottom: 5px;
}
</style>
