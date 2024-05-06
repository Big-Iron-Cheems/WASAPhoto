<script setup>
const props = defineProps(['msg'])
</script>

<template>
    <div class="alert alert-danger" role="alert">
        {{ msg.message || msg }}
    </div>
</template>

<style scoped>
</style>
