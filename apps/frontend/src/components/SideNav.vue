<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Card from 'primevue/card'
import Menu from 'primevue/menu'
import { useAuth } from '../composables/useAuth'

const router = useRouter()
const { currentUser, logout } = useAuth()
const menu = ref()

const menuItems = [
  {
    label: 'Logout',
    icon: 'pi pi-sign-out',
    command: handleLogout
  }
]

const toggleMenu = (event) => {
  menu.value.toggle(event)
}

async function handleLogout() {
  await logout()
  router.push('/login')
}

const navItems = [
  { label: 'Home', icon: 'pi pi-home', to: '/' },
  { label: 'IAS AI (Preview)', icon: 'pi pi-sparkles', to: '/ai' },
  { label: 'Dashboards', icon: 'pi pi-chart-bar', to: '/dashboards' },
  { label: 'Devices', icon: 'pi pi-microchip', to: '/devices' },
  { label: 'Device Profiles', icon: 'pi pi-wrench', to: '/device-profiles' },
  { label: 'Data Browser', icon: 'pi pi-database', to: '/data-browser' },
  { label: 'Ingest Logs', icon: 'pi pi-chevron-circle-down', to: '/ingest-logs' },
  { label: 'Settings', icon: 'pi pi-cog', to: '/settings' },
  { label: 'Extensions', icon: 'pi pi-bolt', to: '/extensions' },
  { label: 'Diagnostics', icon: 'pi pi-wave-pulse', to: '/diagnostics' },
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

        <div class="profile-section" @click="toggleMenu">
            <i class="pi pi-user profile-icon"></i>
            <span class="profile-name">{{ currentUser?.username || 'User' }}</span>
            <span class="profile-role">({{ currentUser?.role }})</span>
            <i class="pi pi-chevron-down profile-arrow"></i>
        </div>
        <Menu ref="menu" :model="menuItems" :popup="true" />
        <div class="ribbon"><span>DEMO BUILD</span></div>
    </div>
</template>

<style scoped>
.sidenav {
  font-family: "Space Grotesk", sans-serif;
    font-size: var(--font-size-md);
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

.profile-section {
    margin-top: auto;
    display: flex;
    align-items: center;
    padding: 14px 20px;
    border-top: 1px solid #212121;
    cursor: pointer;
    transition: background-color 0.3s;
}

.profile-section:hover {
    background-color: #3a3a3e;
}

.profile-icon {
    font-size: 1rem;
    color: #ffffff;
    margin-right: 10px;
}

.profile-name {
    color: #e0e0e0;
    font-size: var(--font-size-sm);
    font-weight: 500;
}

.profile-role {
    color: #888;
    font-size: var(--font-size-xs);
}

.profile-arrow {
    font-size: var(--font-size-xs);
    color: #888;
    margin-left: auto;
}

.ribbon span {
    display: block;
    width: 100%;
    background: rgba(240, 173, 78, 0.85);
    color: #18181B;
    text-align: center;
    font-size: var(--font-size-md);
    font-weight: 700;
    padding: 14px 0;
    letter-spacing: 1px;
}
</style>