
<template>
    <Card class="card">
        <template #title>{{ title }}</template>
        <template #content>
            <div class="card-value">
                {{ value }}
            </div>
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
  value: {
    type: [String, Number],
    default: "0"
  }
});
</script>

<style scoped>
.card {
  width: 18vw;
  min-width: 250px;
  margin: 10px;
  border-color: #202024;
  border-width: 1px;
  border-style: solid;
}

.card-value {
  font-size: 3rem;
  font-weight: bold;
  color: #6d6d6d;
  padding: 20px 0;
}
</style>