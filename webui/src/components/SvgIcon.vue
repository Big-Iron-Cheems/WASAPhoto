<script setup>
import {computed} from "vue";

const props = defineProps(['icon'])
const iconName = computed(() => `/feather-sprite-v4.29.0.svg#${props.icon}`)
</script>

<template>
    <svg class="feather" aria-hidden="true">
        <use :href="iconName"/>
    </svg>
</template>

<style scoped>
</style>
