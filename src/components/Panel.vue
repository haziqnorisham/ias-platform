<template>
    <Card class="card">
        <template #title>{{ title }}</template>
        <template #subtitle>{{ subtitle }}</template>
        <template #content>
            <slot> </slot>
        </template>
    </Card>
</template>

<script setup>
import Card from "primevue/card";

defineProps({
    title: {
        type: String,
        default: "Card"
    },
    subtitle: {
        type: String,
        default: "Card subtitle"
    }
});
</script>

<style scoped>
.card {
    border-color: #202024;
    border-width: 1px;
    border-style: solid;
}
</style>