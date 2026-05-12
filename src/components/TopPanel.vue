
<template>
        <!-- <Menubar  class="notification_bar">
            
        </Menubar> 

        <Menubar :model="items" class="navigation_bar">
        </Menubar>-->
    
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from 'vue-router';
import Menubar from "primevue/menubar";
import Image  from "primevue/image";

const router = useRouter();

// primevue Menubar doesn't automatically turn `to` into router-links; we
// provide a `command` handler that uses the router to navigate.
const items = ref([
    {
        label: 'Home',
        command: () => router.push('/')
    },
    {
        label: 'About',
        command: () => router.push('/about')
    }
]);
</script>
<style scoped>
.notification_bar {
    --p-menubar-border-color: Transparent;
    --p-menubar-border-radius: 0px; 
    border-bottom: 1px solid #212121;
    font-size: 12px;
    
    height: 2.5rem;
}
.navigation_bar {
    display: flex;
    justify-content : center;
    --p-menubar-border-color: Transparent;
    --p-menubar-border-radius: 0px; 
    font-size: 10px;
    
    height: 2rem;
    font-weight: 700;
}
</style>