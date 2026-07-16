
<template>
    <Card class="card">
        <template #title>{{ title }}</template>
        <template #subtitle><div class="subtitle">{{ subtitle }}</div></template>
        <template #content>
            
                <slot></slot>
            
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
        default: "Subtitle"
    },
  value: {
    type: [String, Number],
    default: "0"
  },
  width: {
    type: String,
    default: "50vw"
  },
  height: {
    type: String,
    default: "auto"
  }
});
</script>

<style scoped>
.card {
  width: v-bind(width);
  height: v-bind(height);
  min-width: 250px;
  margin: 10px;
}
.subtitle {
    font-style: italic;
    font-size: 0.7rem;
}
</style>