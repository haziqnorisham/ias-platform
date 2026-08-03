<script setup>
// App.vue is now just a layout wrapper; content is rendered by router-view.
import SideNav from "./components/SideNav.vue";
import Breadcrumb from 'primevue/breadcrumb';
import ScrollPanel from 'primevue/scrollpanel';
import Footer from "./components/Footer.vue";
import Toast from 'primevue/toast';
import Button from 'primevue/button';
import { ref, watch, computed } from "vue";
import { useRoute } from 'vue-router';
import { useSidebar } from './composables/useSidebar';

const route = useRoute();
const { sidebarVisible, toggleSidebar } = useSidebar();

const isLoginPage = computed(() => route.name === 'Login')

const home = ref({
  icon: 'pi pi-home'
});

// Map route names to breadcrumb labels
const routeLabelMap = {
  Home: 'Home',
  Ai: 'IAS AI (Preview)',
  Dashboards: 'Dashboards',
  DashboardEdit: 'Edit Dashboard',
  DashboardView: 'View Dashboard',
  Devices: 'Devices',
  DeviceProfiles: 'Device Profiles',
  IngestLogs: 'Ingest Logs',
  Diagnostics: 'Diagnostics',
  Extensions: 'Extensions',
  ExtensionDetail: 'Extension',
  Settings: 'Settings',
  About: 'About',
};

const items = ref([
  { label: 'Home', to: '/' }
]);

// Update breadcrumbs whenever the route changes.
// The :home prop already renders the home icon, so items should only
// contain the subsequent path segments (no "Home" entry).
watch(
  () => route.path,
  (newPath) => {
    const pathSegments = newPath.split('/').filter(Boolean);
    
    if (pathSegments.length === 0) {
      items.value = [];
      return;
    }

    const crumbs = [];
    let accumulatedPath = '';
    
    for (const segment of pathSegments) {
      accumulatedPath += '/' + segment;
      const routeName = route.matched.find(r => r.path === accumulatedPath)?.name;
      const label = routeLabelMap[routeName] || segment.charAt(0).toUpperCase() + segment.slice(1).replace(/-/g, ' ');
      
      crumbs.push({
        label: label,
        to: accumulatedPath
      });
    }

    items.value = crumbs;
  },
  { immediate: true }
);

</script>

<template>
  <Toast position="top-center"/>

  <template v-if="!isLoginPage">
    <SideNav />
    <Button
      :icon="sidebarVisible ? 'pi pi-chevron-left' : 'pi pi-chevron-right'"
      class="sidebar-toggle"
      :class="{ 'sidebar-toggle--collapsed': !sidebarVisible }"
      text
      rounded
      size="small"
      :title="sidebarVisible ? 'Hide sidebar' : 'Show sidebar'"
      @click="toggleSidebar"
    />
    <main class="main" :class="{ 'main--expanded': !sidebarVisible }">
      <Breadcrumb :home="home" :model="items" class="breadcrumb" />
      <div class="routerView">
        <router-view />
      </div>
      <ScrollPanel class="content">

      </ScrollPanel>
      <div class="footer_position">
        <Footer />
      </div>
    </main>
  </template>
  <template v-else>
    <router-view />
  </template>
</template>

<style scoped>
.main {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 16rem;
  display: flex;
  flex-direction: column;
  padding-left: 40px;
  padding-right: 40px;
  box-sizing: border-box;
  overflow-y: auto;
  overflow-x: hidden;
  background-color: #0e0e10;
  background-image: radial-gradient(#28282c 1px, transparent 1px);
  background-size: 24px 24px;
  transition: left 0.2s ease;
}

.main--expanded {
  left: 0;
}

.sidebar-toggle {
  position: fixed;
  top: 12px;
  left: calc(16rem + 8px);
  z-index: 100;
  transition: left 0.2s ease;
}

.sidebar-toggle--collapsed {
  left: 8px;
}

.content {
  flex: 1;
  min-height: 0;
  border-radius: 12px;
  border: 0px solid #212121;
}
</style>

<style scoped>
.cards_container {
  margin: 10px;
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  justify-content: space-evenly;
  flex-wrap: wrap;
}

.main_container {
  margin-left: 20px;
  margin-right: 20px;
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  justify-content: space-evenly;
}

.breadcrumb {
  border-radius: 12px;
  margin-bottom: 20px;
  background-color: transparent;
}

</style>
