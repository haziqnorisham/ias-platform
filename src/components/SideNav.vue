<script setup>
import { useRouter } from 'vue-router'
import Card from 'primevue/card'

const router = useRouter()

const navItems = [
  { label: 'Home', icon: 'pi pi-home', to: '/' },
  { label: 'IAS AI (Preview)', icon: 'pi pi-sparkles', to: '/ai' },
  { label: 'Dashboards', icon: 'pi pi-chart-bar', to: '/dashboards' },
  { label: 'Devices', icon: 'pi pi-microchip', to: '/devices' },
  { label: 'Device Profiles', icon: 'pi pi-wrench', to: '/device-profiles' },
  { label: 'Ingest Logs', icon: 'pi pi-chevron-circle-down', to: '/ingest-logs' },
  { label: 'Settings', icon: 'pi pi-cog', to: '/settings' },
  { label: 'Diagnostics', icon: 'pi pi-wave-pulse', to: '/diagnostics' },
  { label: 'Grid Playground', icon: 'pi pi-th-large', to: '/grid-playground' },
  { label: 'About', icon: 'pi pi-info-circle', to: '/about' },

]

const navigate = (to) => {
  router.push(to)
}
</script>

<template>
    <div class="sidenav">
        <div class="logo-container">
            <img src="/bitmap.png" alt="IAS Logo" class="logo" />
        </div>

        <div v-for="item in navItems" :key="item.to" class="nav-card" @click="navigate(item.to)">
            <Card class="nav-card-card" :class="{ 'active': $route.path === item.to }">
                <template #content>
                    <div class="card-content">
                        <div class="icon-container">
                            <i :class="item.icon"></i>
                        </div>
                        <div class="text-container">
                            {{ item.label }}
                        </div>
                    </div>
                </template>
            </Card>
        </div>
    </div>
</template>

<style scoped>
.sidenav {
  font-family: "Space Grotesk", sans-serif;
    font-size: 14px;
    height: 100%;
    width: 16rem;
    position: fixed;
    z-index: 1;
    top: 0;
    left: 0;
    background-color: #18181B;
    border-right: 1px solid #212121;
    overflow-x: hidden;
    display: flex;
    flex-direction: column;
    border-radius: 0px;

}

.logo-container {
    display: flex;
    justify-content: left;
    align-items: center;
    padding-top: 20px;
    padding-bottom: 20px;
    padding-left: 40px;
    padding-right: 40px;
    border-bottom: #212121 1px solid;
}


.logo {
    width: 100%;
    object-fit: contain;
}

.menubar{
    border: none;
    padding: 0;
}

.nav-card {
    border-bottom: 1px solid #212121;
    border-radius: 0px;
    cursor: pointer;
    transition: background-color 0.3s;
}

.nav-card:hover {
    background-color: #3a3a3e;
}

.nav-card-card {
    border: none;
    background: transparent;
}

.nav-card-card.active {
    background-color: rgba(255, 255, 255, 0.08);
    border-radius: 0px;
}

.card-content {
    display: flex;
    align-items: center;
    
}

.icon-container {
    margin-right: 10px;
    color: #ffffff;
}

.text-container {
    color: #ffffff;
}
</style>